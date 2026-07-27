-- Seed the first tenant gym: ibckirtipur.nepalgym.xyz
--
-- This creates the organization and an empty, not-yet-live site row. The pages
-- themselves are not seeded here: section content lives in Go (internal/site/
-- templates.go) and duplicating it in SQL would guarantee the two drift apart.
-- The admin picks a template in the builder, which calls apply-template and
-- materializes the pages from the single source of truth.

BEGIN;

INSERT INTO organizations (name, name_ne, slug, address, address_ne, is_active)
SELECT 'IBC Kirtipur', 'आईबीसी कीर्तिपुर', 'ibckirtipur',
       'Kirtipur, Kathmandu', 'कीर्तिपुर, काठमाडौं', true
WHERE NOT EXISTS (SELECT 1 FROM organizations WHERE slug = 'ibckirtipur');

INSERT INTO site_settings (organization_id, template, is_live)
SELECT o.id, 'fight_club', false
FROM organizations o
WHERE o.slug = 'ibckirtipur'
  AND NOT EXISTS (SELECT 1 FROM site_settings s WHERE s.organization_id = o.id);

COMMIT;
