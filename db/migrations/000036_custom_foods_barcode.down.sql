DROP INDEX IF EXISTS idx_food_items_created_by;
DROP INDEX IF EXISTS idx_food_items_barcode;
ALTER TABLE food_items DROP COLUMN IF EXISTS created_by;
ALTER TABLE food_items DROP COLUMN IF EXISTS barcode;
