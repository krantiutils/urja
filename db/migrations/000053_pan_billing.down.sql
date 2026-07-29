DROP TRIGGER IF EXISTS invoice_items_immutable ON invoice_items;
DROP TRIGGER IF EXISTS invoices_immutable ON invoices;
DROP FUNCTION IF EXISTS invoice_items_guard_immutable();
DROP FUNCTION IF EXISTS invoices_guard_immutable();

DROP TABLE IF EXISTS invoice_prints;
DROP TABLE IF EXISTS invoice_items;
DROP TABLE IF EXISTS invoice_counters;
DROP TABLE IF EXISTS invoices;

ALTER TABLE organizations
    DROP COLUMN IF EXISTS tax_address,
    DROP COLUMN IF EXISTS tax_legal_name,
    DROP COLUMN IF EXISTS is_vat_registered,
    DROP COLUMN IF EXISTS pan_number;
