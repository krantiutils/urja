-- Backfill workout rows whose organization_id was never set.
--
-- The member self-service routes took organization_id as an optional query
-- parameter and stored NULL when it was absent. A NULL row belongs to no gym,
-- so every org-scoped query skips it: the member's own history silently loses
-- entries, and staff never see them at all.
--
-- The handlers now resolve the organization from real membership, so no new
-- NULL rows are written. This repairs the ones already stored.
--
-- Only rows whose owner has exactly one active membership are attributable.
-- A member of two gyms genuinely could have logged the session at either, and
-- guessing would file their work under the wrong gym — those rows are left as
-- they are rather than assigned wrongly.

UPDATE workout_logs wl
SET organization_id = om.organization_id
FROM organization_members om
WHERE wl.organization_id IS NULL
  AND om.user_id = wl.user_id
  AND om.status = 'active'
  AND (
    SELECT count(*) FROM organization_members m
    WHERE m.user_id = wl.user_id AND m.status = 'active'
  ) = 1;

UPDATE member_workout_plans mp
SET organization_id = om.organization_id
FROM organization_members om
WHERE mp.organization_id IS NULL
  AND om.user_id = mp.user_id
  AND om.status = 'active'
  AND (
    SELECT count(*) FROM organization_members m
    WHERE m.user_id = mp.user_id AND m.status = 'active'
  ) = 1;
