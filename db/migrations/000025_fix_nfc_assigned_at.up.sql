-- Fix nfc_cards.assigned_at: allow NULL for unassigned cards.
-- The original migration (000007) set assigned_at as NOT NULL, but
-- migration 000024 made user_id nullable (allowing unassigned cards).
-- UnassignCard sets assigned_at = NULL, which violates the constraint.
ALTER TABLE nfc_cards ALTER COLUMN assigned_at DROP NOT NULL;
