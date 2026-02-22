-- Remove all global seed food items.
DELETE FROM food_items WHERE organization_id IS NULL;
