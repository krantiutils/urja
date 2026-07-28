-- Irreversible by design. The up migration recovers which gym a row belongs to;
-- setting those columns back to NULL would re-orphan rows that may since have
-- been edited, and there is no record of which ones this migration touched.
SELECT 1;
