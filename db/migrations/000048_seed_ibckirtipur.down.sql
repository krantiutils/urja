BEGIN;

-- site_settings and site_pages cascade from organizations, but delete
-- explicitly so a partial rollback leaves nothing orphaned.
DELETE FROM site_pages   WHERE organization_id IN (SELECT id FROM organizations WHERE slug = 'ibckirtipur');
DELETE FROM site_leads   WHERE organization_id IN (SELECT id FROM organizations WHERE slug = 'ibckirtipur');
DELETE FROM site_settings WHERE organization_id IN (SELECT id FROM organizations WHERE slug = 'ibckirtipur');
DELETE FROM organizations WHERE slug = 'ibckirtipur';

COMMIT;
