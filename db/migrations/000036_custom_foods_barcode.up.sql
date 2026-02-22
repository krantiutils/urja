ALTER TABLE food_items ADD COLUMN IF NOT EXISTS barcode VARCHAR(50);
ALTER TABLE food_items ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_food_items_barcode ON food_items(barcode) WHERE barcode IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_food_items_created_by ON food_items(created_by) WHERE created_by IS NOT NULL;
