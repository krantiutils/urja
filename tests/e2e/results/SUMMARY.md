# E2E Test Results Summary

**Date:** 2026-02-14
**Go version:** 1.24.0
**Total tests:** 96 (110 including subtests)
**Passed:** 96/96
**Coverage:** 49.3% of internal/pkg statements

## Domains Covered

| Domain | Tests | File |
|--------|-------|------|
| Health | 1 | health_test.go |
| Auth | 12 | auth_test.go |
| Organization | 9 | org_test.go |
| Member | 11 | member_test.go |
| Packages | 7 | packages_test.go |
| Staff | 7 | staff_test.go |
| Attendance | 5 | attendance_test.go |
| Dues | 6 | dues_test.go |
| Accounts | 6 | accounts_test.go |
| Notices | 8 | notice_test.go |
| Feedback | 4 | feedback_test.go |
| SMS API | 6 | smsapi_test.go |
| Absentees | 5 | absentee_test.go |
| NFC Cards | 5 | nfc_test.go |
| Activity Logs | 3 | activitylog_test.go |
| QR Code | 3 | qrcode_test.go |
| Subscription | 7 | subscription_test.go |
| Billing | 4 | billing_test.go |

## Known Application Bugs Found

1. **Absentee List (500):** `absentee/repository.go` references `om.created_at` but `organization_members` table uses `om.joined_at` (migration 000004).
2. **Org Update route unreachable:** `r.Put("/orgs/{orgId}")` is shadowed by `r.Route("/orgs/{orgId}")` subrouter in chi.
3. **Member Delete route unreachable:** `r.Delete("/{id}")` in member routes shadowed by `r.Route("/{memberId}")` for subscription routes.
4. **NFC Unassign (constraint violation):** `assigned_at` has NOT NULL constraint but `UnassignCard` sets it to NULL.
