# PAN Billing — Design

**Date:** 2026-07-29
**Status:** Approved for planning

## Goal

Let a gym issue a tax invoice ("PAN bill") to a customer, and correct one
lawfully when it is wrong. Today the product records money in three places —
`payments`, `transactions`, `dues` — and produces no document at all. Nothing a
gym can hand to a customer, and nothing the Inland Revenue Department would
accept as a book of sales.

First tenant: I.B.C Kirtipur, PAN-registered, not VAT-registered.

## Scope

| In | Out |
|----|-----|
| Org tax settings (PAN number, legal name, address) | VAT (13%) invoices — phase 2, see below |
| Issue a multi-line bill | IRD CBMS real-time sync — only binds VAT sellers |
| Cancel a bill, with reason | Splitting one bill across two payment methods |
| Credit note against a bill | Bulk/batch issue, recurring auto-billing |
| List, filter, view, print (with print audit) | Mobile/Flutter clients |
| "Issue bill" pre-fill from a package sale | Customer-facing bill download |

### VAT is deliberately deferred, not designed away

The gym is PAN-only today, so no bill carries tax. But `vat_rate` and
`vat_amount` columns exist from day one, fixed at `0`, and
`organizations.is_vat_registered` exists defaulting to `false` with no UI.

This is not speculative generality. Invoices are immutable financial records: if
tax columns were added in a later migration, every bill issued before that point
would be structurally different from every bill after it, with no lawful way to
backfill them (you cannot rewrite an issued bill). Three always-zero columns buy
one consistent row shape across the whole history. Turning VAT on later becomes
a settings change plus print-format work, not a migration against immutable
rows.

When VAT is switched on for any gym, CBMS sync becomes a live obligation. That
is a separate project against IRD's API. The schema reserves `synced_at` for it.

---

## The correction model — why there is no "edit bill"

The original request asked for "revising a bill". Under IRD rules an issued bill
is immutable and its number is consumed permanently: it cannot be edited,
deleted, or renumbered. There are exactly two lawful corrections, and which one
applies is a question of timing:

| Situation | Mechanism | Effect on numbering |
|-----------|-----------|---------------------|
| Mistake caught before the bill leaves the counter | **Cancel and reissue** | Original keeps its number, marked `cancelled` with a reason. New bill takes the next number. |
| Bill already with the customer, money moved | **Credit note** | A separate document, own number, referencing the original. Reverses part or all of it. |

The UI surfaces a single **Revise** action on a bill and routes to whichever is
lawful for that bill's state. There is no edit path. An edit button would
silently produce non-compliant books, which is worse than not shipping the
feature.

Cancellation is one-way: a cancelled bill can never return to `issued`. It is
permitted at any time, including in a later fiscal year — the number stays
consumed either way, so there is nothing to protect by adding a deadline.

**Credit notes draw from the same number sequence as invoices.** One unbroken
run of numbers per org per fiscal year covers both document types; `doc_type`
distinguishes them. Two parallel sequences would make "no gaps" ambiguous to
audit.

### Effect on the ledger

The two corrections differ here, and the distinction is the whole reason both
exist:

- **Cancel touches `transactions` only to undo what the bill itself wrote.** A
  cancelled bill is a document withdrawn before it counted, so it reverses
  exactly the income it created (see `owns_transaction` below) and nothing
  more — a package sale's income is never this document's to touch. If money
  genuinely moved and there is nothing of the bill's own to reverse, cancelling
  is the wrong instrument — use a credit note.
- **A credit note always writes a reversing ledger row**: a `transactions` row
  of type `expense`, category `refund`, for the credit amount, linked back to
  the credit note. Always, regardless of whether the original bill created its
  income row or merely linked a package sale's. A refund is a real movement of
  money; making the reversal conditional on the original row's provenance would
  leave refunded membership income sitting on the books.

The `expense`/`refund` mapping is a known simplification. A sales return is
properly contra-revenue, but `transactions.transaction_type` only admits
`income` and `expense`, and inventing a third type would ripple through every
existing accounts report. Revisit if the books need it.

**A bill's own income follows the same never-touch-what-you-don't-own rule as
the credit note's refund above, mirrored for issue rather than cancel.**
Ownership cannot be inferred from the row alone — a linked transaction and an
owned one look identical once written — so `invoices.owns_transaction` records
it explicitly at insert, and the immutability trigger protects it from
changing afterwards:

| Bill | At issue | On cancel |
|------|----------|-----------|
| No `transaction_id` supplied (raised from scratch) | Creates an `income`/`Sales` row, links it, and sets `owns_transaction` | Reverses it with an `expense`/`Sales reversal` row |
| `transaction_id` supplied (package sale pre-fill) | Links the existing row, writes nothing, `owns_transaction` stays false | Writes nothing |

A bill only ever reverses income it created. A package sale's income belongs
to the package flow (`recordPackageIncome`) and is not this document's to
undo.

---

## Data model

New Go package `internal/invoice`. Note `internal/billing` already exists and is
**unrelated** — it is the gym's own SaaS subscription to Urja (plans, subscribe).
Do not extend it.

Migration `000053_pan_billing`.

### Org tax settings

```sql
ALTER TABLE organizations
    ADD COLUMN pan_number         VARCHAR(9),
    ADD COLUMN is_vat_registered  BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN tax_legal_name     VARCHAR(255),
    ADD COLUMN tax_address        TEXT;

ALTER TABLE organizations
    ADD CONSTRAINT organizations_pan_number_format
    CHECK (pan_number IS NULL OR pan_number ~ '^[0-9]{9}$');
```

`tax_legal_name` and `tax_address` are separate from `name`/`address` because a
registered name frequently differs from the trading name. Both fall back to the
org fields when blank.

### Invoices

```sql
CREATE TABLE invoices (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id       UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Numbering
    fiscal_year           VARCHAR(7)  NOT NULL,          -- '2082-83'
    sequence              INT         NOT NULL CHECK (sequence > 0),
    invoice_number        VARCHAR(32) NOT NULL,          -- '2082-83/000042'

    doc_type              VARCHAR(20) NOT NULL DEFAULT 'invoice'
                          CHECK (doc_type IN ('invoice', 'credit_note')),
    credit_note_for       UUID REFERENCES invoices(id),

    -- Seller snapshot (see "Why snapshots" below)
    seller_name           VARCHAR(255) NOT NULL,
    seller_pan            VARCHAR(9)   NOT NULL,
    seller_address        TEXT,
    seller_vat_registered BOOLEAN      NOT NULL DEFAULT false,

    -- Customer snapshot
    customer_user_id      UUID REFERENCES users(id) ON DELETE SET NULL,
    customer_name         VARCHAR(255) NOT NULL,
    customer_pan          VARCHAR(9),
    customer_address      TEXT,
    customer_phone        VARCHAR(20),

    issued_date           DATE        NOT NULL,          -- AD
    issued_date_bs        VARCHAR(10) NOT NULL,          -- '2082-04-14'

    subtotal              DECIMAL(12,2) NOT NULL CHECK (subtotal >= 0),
    discount              DECIMAL(12,2) NOT NULL DEFAULT 0 CHECK (discount >= 0),
    taxable_amount        DECIMAL(12,2) NOT NULL CHECK (taxable_amount >= 0),
    vat_rate              DECIMAL(5,2)  NOT NULL DEFAULT 0 CHECK (vat_rate >= 0),
    vat_amount            DECIMAL(12,2) NOT NULL DEFAULT 0 CHECK (vat_amount >= 0),
    total                 DECIMAL(12,2) NOT NULL CHECK (total >= 0),
    amount_in_words       TEXT          NOT NULL,

    payment_method        VARCHAR(50)
                          CHECK (payment_method IS NULL OR payment_method IN
                                 ('cash', 'bank_transfer', 'khalti')),

    status                VARCHAR(20) NOT NULL DEFAULT 'issued'
                          CHECK (status IN ('issued', 'cancelled')),
    cancelled_at          TIMESTAMPTZ,
    cancelled_by          UUID REFERENCES users(id),
    cancellation_reason   TEXT,

    -- Ledger linkage (approach C: link, never duplicate)
    transaction_id        UUID REFERENCES transactions(id) ON DELETE SET NULL,
    member_package_id     UUID REFERENCES member_packages(id) ON DELETE SET NULL,

    synced_at             TIMESTAMPTZ,                   -- reserved for CBMS
    issued_by             UUID NOT NULL REFERENCES users(id),
    print_count           INT  NOT NULL DEFAULT 0 CHECK (print_count >= 0),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (organization_id, fiscal_year, sequence),
    UNIQUE (organization_id, invoice_number),
    CONSTRAINT credit_note_has_parent CHECK (
        (doc_type = 'credit_note') = (credit_note_for IS NOT NULL)
    ),
    CONSTRAINT cancelled_has_reason CHECK (
        status <> 'cancelled'
        OR (cancelled_at IS NOT NULL AND cancellation_reason IS NOT NULL)
    )
);

CREATE INDEX idx_invoices_org_issued   ON invoices(organization_id, issued_date DESC);
CREATE INDEX idx_invoices_org_status   ON invoices(organization_id, status);
CREATE INDEX idx_invoices_org_fy       ON invoices(organization_id, fiscal_year);
CREATE INDEX idx_invoices_customer     ON invoices(customer_user_id) WHERE customer_user_id IS NOT NULL;
CREATE INDEX idx_invoices_credit_for   ON invoices(credit_note_for) WHERE credit_note_for IS NOT NULL;
CREATE INDEX idx_invoices_transaction  ON invoices(transaction_id) WHERE transaction_id IS NOT NULL;

-- One live bill per ledger row: makes double-billing the same payment
-- impossible rather than merely discouraged. Cancelled bills are excluded so a
-- cancel-and-reissue against the same payment still works.
CREATE UNIQUE INDEX idx_invoices_one_per_transaction
    ON invoices(transaction_id)
    WHERE transaction_id IS NOT NULL
      AND status = 'issued'
      AND doc_type = 'invoice';
```

There is deliberately **no `updated_at` and no `set_updated_at` trigger** —
an invoice does not get updated in the ordinary sense, and a mutation trigger
would fight the immutability guard.

#### Why snapshots

`seller_*` and `customer_*` are copied onto the row rather than read through a
FK at print time. If the gym corrects its PAN, renames itself, or a member
changes their name, bills issued last year must not silently change. `vat_rate`
is stored per bill for the same reason: a bill printed today and reprinted after
the gym registers for VAT must show what it originally showed.

### Line items

```sql
CREATE TABLE invoice_items (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id     UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    line_no        INT  NOT NULL CHECK (line_no > 0),
    description    TEXT NOT NULL,
    description_ne TEXT,
    quantity       DECIMAL(10,2) NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price     DECIMAL(12,2) NOT NULL CHECK (unit_price >= 0),
    amount         DECIMAL(12,2) NOT NULL CHECK (amount >= 0),
    UNIQUE (invoice_id, line_no)
);
```

### Print audit

```sql
CREATE TABLE invoice_prints (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    printed_by UUID NOT NULL REFERENCES users(id),
    printed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    copy_label VARCHAR(20) NOT NULL CHECK (copy_label IN ('original', 'copy'))
);

CREATE INDEX idx_invoice_prints_invoice ON invoice_prints(invoice_id);
```

The first print of a bill is `original`; every subsequent print is `copy` and
the printed document is marked as such.

### Gapless numbering

```sql
CREATE TABLE invoice_counters (
    organization_id UUID        NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    fiscal_year     VARCHAR(7)  NOT NULL,
    next_sequence   INT         NOT NULL DEFAULT 1 CHECK (next_sequence > 0),
    PRIMARY KEY (organization_id, fiscal_year)
);
```

**A Postgres sequence is wrong here.** Sequences are non-transactional: a rolled
back insert burns its number and leaves a hole in the book. IRD requires an
unbroken run. So allocation is a locked counter row inside the same transaction
as the insert:

```sql
-- 1. Create-or-lock the counter row for this org and fiscal year.
INSERT INTO invoice_counters (organization_id, fiscal_year, next_sequence)
VALUES ($1, $2, 1)
ON CONFLICT (organization_id, fiscal_year)
DO UPDATE SET next_sequence = invoice_counters.next_sequence   -- no-op write; takes the row lock
RETURNING next_sequence;

-- 2. INSERT the invoice with that sequence.

-- 3. Bump the counter.
UPDATE invoice_counters
SET next_sequence = next_sequence + 1
WHERE organization_id = $1 AND fiscal_year = $2;
```

The `DO UPDATE` no-op write is what takes the row lock; a plain
`ON CONFLICT DO NOTHING` returns no row and would race. Concurrent issues for
the same org serialise on that lock. `UNIQUE (organization_id, fiscal_year,
sequence)` is the backstop if the logic is ever bypassed.

If the whole transaction rolls back, the counter rolls back with it — no gap.

### Immutability, enforced in the database

Service-layer checks are not enough; this is the guarantee the whole feature
rests on.

```sql
CREATE OR REPLACE FUNCTION invoices_guard_immutable() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        RAISE EXCEPTION 'invoice %: issued invoices cannot be deleted', OLD.invoice_number;
    END IF;

    -- The primary key is guarded on its own, outside the positional tuple
    -- below: it isn't a column a future migration would think to add to that
    -- list, so it needs its own explicit check rather than relying on being
    -- remembered.
    IF NEW.id IS DISTINCT FROM OLD.id THEN
        RAISE EXCEPTION 'invoice %: the primary key cannot be changed', OLD.invoice_number;
    END IF;

    -- Only cancellation fields, print_count and synced_at may ever change.
    -- id is guarded separately above, not here.
    -- NOTE: this comparison is positional. A new column added to `invoices`
    -- must be added here too, or it silently becomes mutable.
    IF (NEW.organization_id, NEW.fiscal_year, NEW.sequence, NEW.invoice_number,
        NEW.doc_type, NEW.credit_note_for,
        NEW.seller_name, NEW.seller_pan, NEW.seller_address, NEW.seller_vat_registered,
        NEW.customer_user_id, NEW.customer_name, NEW.customer_pan,
        NEW.customer_address, NEW.customer_phone,
        NEW.issued_date, NEW.issued_date_bs,
        NEW.subtotal, NEW.discount, NEW.taxable_amount,
        NEW.vat_rate, NEW.vat_amount, NEW.total, NEW.amount_in_words,
        NEW.payment_method, NEW.transaction_id, NEW.member_package_id,
        NEW.issued_by, NEW.created_at)
       IS DISTINCT FROM
       (OLD.organization_id, OLD.fiscal_year, OLD.sequence, OLD.invoice_number,
        OLD.doc_type, OLD.credit_note_for,
        OLD.seller_name, OLD.seller_pan, OLD.seller_address, OLD.seller_vat_registered,
        OLD.customer_user_id, OLD.customer_name, OLD.customer_pan,
        OLD.customer_address, OLD.customer_phone,
        OLD.issued_date, OLD.issued_date_bs,
        OLD.subtotal, OLD.discount, OLD.taxable_amount,
        OLD.vat_rate, OLD.vat_amount, OLD.total, OLD.amount_in_words,
        OLD.payment_method, OLD.transaction_id, OLD.member_package_id,
        OLD.issued_by, OLD.created_at)
    THEN
        RAISE EXCEPTION 'invoice %: only cancellation may change an issued invoice',
                        OLD.invoice_number;
    END IF;

    IF OLD.status = 'cancelled' AND NEW.status <> 'cancelled' THEN
        RAISE EXCEPTION 'invoice %: cancellation cannot be undone', OLD.invoice_number;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER invoices_immutable
    BEFORE UPDATE OR DELETE ON invoices
    FOR EACH ROW EXECUTE FUNCTION invoices_guard_immutable();
```

`invoice_items` is insert-only — a trigger raises on any `UPDATE` or `DELETE`.
(Deleting the parent invoice is already impossible, so the `ON DELETE CASCADE`
is unreachable in practice and exists only for org deletion.)

---

## Supporting packages

Both are pure functions with no I/O, which makes them cheap to test exhaustively.

### `pkg/nepalidate`

AD ↔ BS conversion over the standard 2000–2100 BS days-per-month table, plus
fiscal-year derivation. The Nepali fiscal year runs Shrawan 1 → Ashad end, so
`FiscalYear(bsDate)` returns `'2082-83'` for any date from 2082-04-01 through
2083-03-{30,31,32}.

Conversion lives **only in Go**. The API returns both `issued_date` and
`issued_date_bs` as strings, so the web never needs a second implementation that
could drift from the first.

### `pkg/moneywords`

Amount in words using Nepali numbering — crore, lakh, thousand — for the
"Amount in words" line every bill must carry. Handles paisa, zero, and the
lakh/crore boundaries.

---

## API

Mounted at `/api/v1/orgs/{orgId}/invoices`, all behind
`RequireOrgRole("admin", "staff")`.

| Method | Path | Purpose |
|--------|------|---------|
| `GET`  | `/` | List. Filters: `status`, `fiscal_year`, `customer_user_id`, `from`, `to`, `q` (matches invoice number or customer name), `limit`, `offset` |
| `POST` | `/` | Create **and issue** atomically |
| `GET`  | `/{id}` | Detail with line items and credit notes |
| `POST` | `/{id}/cancel` | Body `{reason}` — required |
| `POST` | `/{id}/credit-note` | Create a linked credit note |
| `POST` | `/{id}/print` | Log the print, bump `print_count`, return render payload |
| `GET`  | `/next-number` | Preview the number the next bill will take |

Tax settings ride on the existing org update path rather than a new endpoint.

**There is no draft state and no `PUT`.** A number is consumed only on issue,
and nothing about an issued bill can be rewritten.

### Error contract

| Condition | Status | Code |
|-----------|--------|------|
| PAN not set on the org | 400 | `pan_not_configured` |
| Cancel an already-cancelled bill | 409 | `already_cancelled` |
| Cancel with no reason | 400 | `reason_required` |
| Credit note against a cancelled bill | 409 | `invoice_cancelled` |
| Credit note exceeds uncredited remainder | 400 | `credit_exceeds_balance` |
| Credit note against a credit note | 400 | `invalid_parent` |
| Invoice belongs to another org | 404 | `not_found` |

`pan_not_configured` is a distinct code specifically so the UI can link straight
to settings rather than dead-ending on a message. This is the "they just need to
add their PAN number" requirement made enforceable: no PAN, no bills.

Cross-tenant reads return **404, not 403** — a 403 confirms the row exists.

---

## UI

| Screen | Contents |
|--------|----------|
| `/dashboard/settings` | New **Tax** section: PAN (9 digits, validated client and server), legal name, address |
| `/dashboard/invoices` | List, status chips (issued / cancelled / credit note), filters, **New bill** |
| `/dashboard/invoices/new` | Customer picker (existing member or walk-in name), line items, live totals |
| `/dashboard/invoices/[id]` | The document, with **Print**, **Cancel**, **Revise** |

**Revise** routes to cancel-and-reissue or credit note based on the bill's state,
and explains which one it is doing and why. It never edits.

Print is a dedicated stylesheet, A4 and A5: seller header with PAN, both AD and
BS dates, customer block, line items, totals, amount in words, and an
Original/Copy mark driven by `print_count`.

The package assign flow gains an **Issue bill** action that pre-fills line items
from the sale and sets `transaction_id` and `member_package_id` on the new
invoice — linking the existing ledger row rather than writing a second one. This
is the anti-double-count that decided approach C: `recordPackageIncome` already
fires on assign, renew, and extend, so a bill must never create income of its
own for a package sale.

A bill raised from scratch (gear, admission fee, personal training) writes its
own `transactions` income row and links it the same way.

All strings land in both `en` and `ne`; the repo's i18n parity test fails
otherwise.

---

## Testing

### Unit

- `nepalidate`: known AD/BS pairs; the Ashad→Shrawan fiscal boundary in both
  directions; 32-day Ashad years; year bounds.
- `moneywords`: zero, paisa, 1 lakh, 1 crore, and the boundaries either side.
- Totals: subtotal − discount = taxable; VAT stays 0 throughout PAN-only.

### Integration

These are the ones that carry the design:

1. **Gapless numbering under concurrency** — N goroutines issue simultaneously
   for one org; assert the sequences are exactly 1..N with no gap and no
   duplicate. Then force a mid-transaction rollback and assert the next issue
   reuses the number rather than skipping it.
2. **Immutability** — direct `UPDATE` of an amount is rejected; `DELETE` is
   rejected; `UPDATE` of cancellation fields succeeds; un-cancelling is
   rejected; `UPDATE`/`DELETE` on `invoice_items` is rejected.
3. **Cancel** keeps the number consumed — cancel bill 5, issue again, assert the
   new bill is 6 and 5 still exists as cancelled.
4. **Credit note** links correctly, caps at the uncredited remainder, and is
   refused against cancelled bills and against other credit notes.
5. **Fiscal year rollover** — a bill on Ashad 31 and one on Shrawan 1 land in
   different `fiscal_year` values, and the sequence restarts at 1.
6. **Cross-tenant isolation** — org A cannot list, read, cancel, credit-note or
   print org B's invoices. Not paranoia: packages, training guides, and the
   member roster have each leaked across gyms in this codebase.
7. **PAN guard** — issuing with no PAN returns `pan_not_configured`; setting the
   PAN then issuing succeeds.
8. **Print audit** — first print logs `original`, second logs `copy`,
   `print_count` matches the row count in `invoice_prints`.
9. **No double-billing a payment** — issuing a second bill against a
   `transaction_id` that already has a live bill is rejected; cancelling the
   first then reissuing succeeds.
10. **Ledger effects** — cancelling leaves `transactions` untouched; a credit
    note writes exactly one reversing `expense`/`refund` row for the credit
    amount, in both the from-scratch and package-linked cases.

---

## Risks

**The immutability trigger is positional.** Adding a column to `invoices`
without adding it to the trigger's comparison silently makes that column
mutable. The migration carries a comment saying so, and test 2 should be
extended whenever the table is.

**`ON DELETE SET NULL` on `transaction_id`.** If a ledger row is deleted the
bill survives with a dangling link. That is the correct trade — the bill must
never disappear — but it means the ledger and the book can diverge. Worth a
reconciliation view later; out of scope here.

**No bill-level payment split.** One bill maps to one payment method. A customer
paying part cash, part Khalti needs two bills today. Flagged to the user;
deferred by agreement.
