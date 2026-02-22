-- Remove seeded workout programs and their associated program_days (CASCADE handles days).
DELETE FROM program_days WHERE program_id IN (
    SELECT id FROM workout_programs WHERE name IN (
        '30-Day HIIT Challenge',
        'Beginner Bodyweight',
        'Core Crusher',
        'Full Body Strength',
        'Flexibility & Mobility'
    )
);

DELETE FROM workout_programs WHERE name IN (
    '30-Day HIIT Challenge',
    'Beginner Bodyweight',
    'Core Crusher',
    'Full Body Strength',
    'Flexibility & Mobility'
);
