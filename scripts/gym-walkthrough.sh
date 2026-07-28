#!/bin/bash
# Drive the sequence a gym actually performs, in order, against a live API.
# Unit tests check pieces; this checks that the pieces join up.
set -uo pipefail
API=http://localhost:8080
PASS=0; FAIL=0

step() { printf "\n\033[1m%s\033[0m\n" "$1"; }
check() { # check <label> <actual> <expected>
  if [ "$2" = "$3" ]; then printf "  ok    %-46s %s\n" "$1" "$2"; PASS=$((PASS+1))
  else printf "  FAIL  %-46s got %s want %s\n" "$1" "$2" "$3"
       printf "        %s\n" "$(head -c 160 /tmp/wt.out 2>/dev/null)"; FAIL=$((FAIL+1)); fi
}
code() { curl -s -o /tmp/wt.out -w "%{http_code}" "$@"; }

PSQL="psql -h 127.0.0.1 -U urja -d urja -tAc"
export PGPASSWORD=urja_secret
set -a; . "$(dirname "$0")/../.env"; set +a

SUF=$(printf "%05d" $((RANDOM % 90000 + 10000)))  # fixed width: phones must be 10 digits
ORG=$($PSQL "INSERT INTO organizations (name, slug) VALUES ('Walk Gym $SUF','walk-gym-$SUF') RETURNING id" | head -1 | tr -d '[:space:]')
OWNER=$($PSQL "INSERT INTO users (phone, name, user_type) VALUES ('98007$SUF','Owner','gym_member') RETURNING id" | head -1 | tr -d '[:space:]')
$PSQL "INSERT INTO organization_members (user_id, organization_id, role, status) VALUES ('$OWNER','$ORG','admin','active')" >/dev/null
ADMIN=$(python3 "$(dirname "$0")/mint-token.py" "$OWNER" "98007$SUF" admin)
H_ADMIN=(-H "Authorization: Bearer $ADMIN" -H "Content-Type: application/json")

step "1. Owner creates a membership package"
c=$(code -X POST "${H_ADMIN[@]}" -d '{"name":"Monthly","duration_days":30,"price":"3000"}' "$API/api/v1/orgs/$ORG/packages")
check "create package" "$c" 201
PKG=$(python3 -c "import json;print(json.load(open('/tmp/wt.out')).get('id',''))" 2>/dev/null)

step "2. Owner adds a member at the desk"
c=$(code -X POST "${H_ADMIN[@]}" -d "{\"phone\":\"98008$SUF\",\"name\":\"Walk In\"}" "$API/api/v1/orgs/$ORG/members")
check "create member" "$c" 201
MEMBER=$(python3 -c "import json;d=json.load(open('/tmp/wt.out'));print(d.get('id') or d.get('user_id',''))" 2>/dev/null)

step "3. Owner assigns the package, member pays 1000 of 3000"
c=$(code -X POST "${H_ADMIN[@]}" -d "{\"package_id\":\"$PKG\",\"start_date\":\"$(date +%F)\",\"payment_method\":\"cash\",\"amount_paid\":1000}" \
  "$API/api/v1/orgs/$ORG/members/$MEMBER/packages/assign")
check "assign package" "$c" 201

step "4. The 2000 balance shows up as money owed"
code "${H_ADMIN[@]}" "$API/api/v1/orgs/$ORG/dues" >/dev/null
DUES=$(python3 -c "import json;print(len(json.load(open('/tmp/wt.out')).get('data') or []))")
check "dues raised" "$DUES" 1
DUE_AMT=$(python3 -c "import json;d=json.load(open('/tmp/wt.out'))['data'];print(int(d[0]['amount']) if d else 0)")
check "balance amount" "$DUE_AMT" 2000
DUE_ID=$(python3 -c "import json;d=json.load(open('/tmp/wt.out'))['data'];print(d[0]['id'] if d else '')")

step "5. Staff check the member in"
c=$(code -X POST "${H_ADMIN[@]}" -d "{\"member_id\":\"$MEMBER\"}" "$API/api/v1/orgs/$ORG/attendance/check-in")
check "manual check-in" "$c" 201
code "${H_ADMIN[@]}" "$API/api/v1/orgs/$ORG/attendance" >/dev/null
ATT=$(python3 -c "import json;print(len(json.load(open('/tmp/wt.out')).get('data') or []))")
check "attendance recorded" "$ATT" 1

step "6. Member pays off the balance"
c=$(code -X POST "${H_ADMIN[@]}" -d '{"amount":2000,"payment_method":"cash"}' "$API/api/v1/orgs/$ORG/dues/$DUE_ID/pay")
check "record payment" "$c" 200
code "${H_ADMIN[@]}" "$API/api/v1/orgs/$ORG/dues?status=unpaid" >/dev/null
LEFT=$(python3 -c "import json;print(len(json.load(open('/tmp/wt.out')).get('data') or []))")
check "no unpaid dues left" "$LEFT" 0

step "7. Does the payment reach the books?"
code "${H_ADMIN[@]}" "$API/api/v1/orgs/$ORG/accounts/summary" >/dev/null
INCOME=$(python3 -c "
import json
d=json.load(open('/tmp/wt.out'))
print(int(float(d.get('total_income') or 0)))" 2>/dev/null || echo 0)
check "income recorded in accounts" "$INCOME" 3000

step "8. Member sees their own membership"
MTOK=$(python3 "$(dirname "$0")/mint-token.py" "$MEMBER" "98008$SUF" member)
c=$(code -H "Authorization: Bearer $MTOK" "$API/api/v1/members/me/packages")
check "member sees own packages" "$c" 200
c=$(code -H "Authorization: Bearer $MTOK" "$API/api/v1/members/me/attendance")
check "member sees own attendance" "$c" 200

printf "\n\033[1m%d passed, %d failed\033[0m\n" "$PASS" "$FAIL"
$PSQL "DELETE FROM organizations WHERE id='$ORG'" >/dev/null 2>&1
exit $FAIL
