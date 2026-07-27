BEGIN;

DROP INDEX IF EXISTS idx_site_leads_org;
DROP INDEX IF EXISTS idx_site_pages_org;

DROP TABLE IF EXISTS site_leads;
DROP TABLE IF EXISTS site_settings;
DROP TABLE IF EXISTS site_pages;

COMMIT;
