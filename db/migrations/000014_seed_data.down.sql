-- Remove seed data in reverse order of dependencies
DELETE FROM saas_plans WHERE id IN (
    'd1b2c3d4-e5f6-7890-abcd-ef1234567801',
    'd1b2c3d4-e5f6-7890-abcd-ef1234567802',
    'd1b2c3d4-e5f6-7890-abcd-ef1234567803'
);

DELETE FROM achievements WHERE id IN (
    'c1b2c3d4-e5f6-7890-abcd-ef1234567801',
    'c1b2c3d4-e5f6-7890-abcd-ef1234567802',
    'c1b2c3d4-e5f6-7890-abcd-ef1234567803',
    'c1b2c3d4-e5f6-7890-abcd-ef1234567804',
    'c1b2c3d4-e5f6-7890-abcd-ef1234567805',
    'c1b2c3d4-e5f6-7890-abcd-ef1234567806',
    'c1b2c3d4-e5f6-7890-abcd-ef1234567807',
    'c1b2c3d4-e5f6-7890-abcd-ef1234567808'
);

DELETE FROM workout_templates WHERE id IN (
    'b1b2c3d4-e5f6-7890-abcd-ef1234567801',
    'b1b2c3d4-e5f6-7890-abcd-ef1234567802',
    'b1b2c3d4-e5f6-7890-abcd-ef1234567803',
    'b1b2c3d4-e5f6-7890-abcd-ef1234567804',
    'b1b2c3d4-e5f6-7890-abcd-ef1234567805'
);

DELETE FROM organizations WHERE id IN (
    'a1b2c3d4-e5f6-7890-abcd-ef1234567801',
    'a1b2c3d4-e5f6-7890-abcd-ef1234567802',
    'a1b2c3d4-e5f6-7890-abcd-ef1234567803'
);
