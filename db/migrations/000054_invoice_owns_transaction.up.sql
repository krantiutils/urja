-- Whether this invoice created the transactions row it links, and must
-- therefore reverse it when cancelled. A package sale's income row is linked,
-- not owned: reversing it here would undo money the package flow recorded.
ALTER TABLE invoices
    ADD COLUMN owns_transaction BOOLEAN NOT NULL DEFAULT false;

-- Issued invoices are immutable. This is the guarantee the whole feature rests
-- on, so it lives in the database rather than only in the service layer.
--
-- Redefines the function from 000053_pan_billing with owns_transaction added
-- to the protected tuple below: the comparison is positional, so a column
-- left out of it is silently mutable, and that is exactly the defect this
-- plan already had to fix once for id.
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
        NEW.issued_by, NEW.created_at, NEW.owns_transaction)
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
        OLD.issued_by, OLD.created_at, OLD.owns_transaction)
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

-- The trigger already points at this function by name; replacing the
-- function body is enough, no need to recreate the trigger itself.
