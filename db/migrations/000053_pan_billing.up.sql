-- PAN tax invoicing: org tax settings, invoices and their line items, print
-- audit, and gapless per-org/per-fiscal-year numbering.
--
-- vat_rate/vat_amount/is_vat_registered exist now, always zero/false, because
-- invoices are immutable: adding tax columns after bills already existed would
-- leave earlier rows structurally different from later ones with no lawful
-- way to backfill them.

ALTER TABLE organizations
    ADD COLUMN pan_number         VARCHAR(9),
    ADD COLUMN is_vat_registered  BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN tax_legal_name     VARCHAR(255),
    ADD COLUMN tax_address        TEXT;

ALTER TABLE organizations
    ADD CONSTRAINT organizations_pan_number_format
    CHECK (pan_number IS NULL OR pan_number ~ '^[0-9]{9}$');

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

-- Deliberately no updated_at / trigger_set_updated_at here: an invoice is
-- never updated in the ordinary sense, and a mutation trigger would fight
-- the immutability guard below.

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

CREATE TABLE invoice_prints (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    printed_by UUID NOT NULL REFERENCES users(id),
    printed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    copy_label VARCHAR(20) NOT NULL CHECK (copy_label IN ('original', 'copy'))
);

CREATE INDEX idx_invoice_prints_invoice ON invoice_prints(invoice_id);

-- Gapless numbering: a Postgres sequence is wrong here because sequences are
-- non-transactional and a rolled-back insert would leave a permanent hole in
-- a run of numbers that IRD requires to be unbroken. Allocation instead locks
-- this counter row inside the same transaction as the invoice insert.
CREATE TABLE invoice_counters (
    organization_id UUID        NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    fiscal_year     VARCHAR(7)  NOT NULL,
    next_sequence   INT         NOT NULL DEFAULT 1 CHECK (next_sequence > 0),
    PRIMARY KEY (organization_id, fiscal_year)
);

-- Issued invoices are immutable. This is the guarantee the whole feature rests
-- on, so it lives in the database rather than only in the service layer.
CREATE OR REPLACE FUNCTION invoices_guard_immutable() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        RAISE EXCEPTION 'invoice %: issued invoices cannot be deleted', OLD.invoice_number;
    END IF;

    -- Only cancellation fields, print_count and synced_at may ever change.
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

-- Line items are insert-only.
CREATE OR REPLACE FUNCTION invoice_items_guard_immutable() RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'invoice line items are immutable';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER invoice_items_immutable
    BEFORE UPDATE OR DELETE ON invoice_items
    FOR EACH ROW EXECUTE FUNCTION invoice_items_guard_immutable();
