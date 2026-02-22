-- Seed muscle groups into preset workout template exercises and set goal/days_per_week.

-- 1. Full Body Beginner
UPDATE workout_templates SET
    exercises = '[
        {"name": "Bodyweight Squats", "name_ne": "बडीवेट स्क्वाट", "sets": 3, "reps": 12, "rest_seconds": 60, "muscle_groups": ["quadriceps", "glutes", "hamstrings"]},
        {"name": "Push-ups", "name_ne": "पुश-अप", "sets": 3, "reps": 10, "rest_seconds": 60, "muscle_groups": ["chest", "triceps", "shoulders"]},
        {"name": "Dumbbell Rows", "name_ne": "डम्बेल रो", "sets": 3, "reps": 10, "rest_seconds": 60, "muscle_groups": ["back", "biceps", "lats"]},
        {"name": "Plank", "name_ne": "प्ल्यांक", "sets": 3, "reps": 1, "rest_seconds": 45, "notes": "Hold for 30 seconds", "muscle_groups": ["core", "obliques"]},
        {"name": "Lunges", "name_ne": "लन्ज", "sets": 3, "reps": 10, "rest_seconds": 60, "muscle_groups": ["quadriceps", "glutes", "hamstrings"]}
    ]'::jsonb,
    goal = 'stay_fit',
    days_per_week = '2-3'
WHERE id = 'b1b2c3d4-e5f6-7890-abcd-ef1234567801';

-- 2. Upper Body Strength
UPDATE workout_templates SET
    exercises = '[
        {"name": "Bench Press", "name_ne": "बेन्च प्रेस", "sets": 4, "reps": 8, "rest_seconds": 90, "muscle_groups": ["chest", "triceps", "shoulders"]},
        {"name": "Overhead Press", "name_ne": "ओभरहेड प्रेस", "sets": 3, "reps": 10, "rest_seconds": 90, "muscle_groups": ["shoulders", "triceps"]},
        {"name": "Barbell Rows", "name_ne": "बारबेल रो", "sets": 4, "reps": 8, "rest_seconds": 90, "muscle_groups": ["back", "biceps", "lats"]},
        {"name": "Dumbbell Curls", "name_ne": "डम्बेल कर्ल", "sets": 3, "reps": 12, "rest_seconds": 60, "muscle_groups": ["biceps", "forearms"]},
        {"name": "Tricep Dips", "name_ne": "ट्राइसेप डिप्स", "sets": 3, "reps": 10, "rest_seconds": 60, "muscle_groups": ["triceps", "chest", "shoulders"]},
        {"name": "Lateral Raises", "name_ne": "ल्याटरल रेज", "sets": 3, "reps": 12, "rest_seconds": 60, "muscle_groups": ["shoulders"]}
    ]'::jsonb,
    goal = 'build_muscle',
    days_per_week = '4-5'
WHERE id = 'b1b2c3d4-e5f6-7890-abcd-ef1234567802';

-- 3. Lower Body Power
UPDATE workout_templates SET
    exercises = '[
        {"name": "Barbell Squats", "name_ne": "बारबेल स्क्वाट", "sets": 4, "reps": 8, "rest_seconds": 120, "muscle_groups": ["quadriceps", "glutes", "hamstrings", "core"]},
        {"name": "Romanian Deadlifts", "name_ne": "रोमानियन डेडलिफ्ट", "sets": 4, "reps": 8, "rest_seconds": 90, "muscle_groups": ["hamstrings", "glutes", "lower_back"]},
        {"name": "Leg Press", "name_ne": "लेग प्रेस", "sets": 3, "reps": 12, "rest_seconds": 90, "muscle_groups": ["quadriceps", "glutes"]},
        {"name": "Walking Lunges", "name_ne": "वाकिंग लन्ज", "sets": 3, "reps": 12, "rest_seconds": 60, "muscle_groups": ["quadriceps", "glutes", "hamstrings"]},
        {"name": "Calf Raises", "name_ne": "काफ रेज", "sets": 4, "reps": 15, "rest_seconds": 45, "muscle_groups": ["calves"]}
    ]'::jsonb,
    goal = 'build_muscle',
    days_per_week = '4-5'
WHERE id = 'b1b2c3d4-e5f6-7890-abcd-ef1234567803';

-- 4. HIIT Cardio Blast
UPDATE workout_templates SET
    exercises = '[
        {"name": "Burpees", "name_ne": "बर्पी", "sets": 4, "reps": 10, "rest_seconds": 30, "muscle_groups": ["chest", "core", "quadriceps", "shoulders"]},
        {"name": "Mountain Climbers", "name_ne": "माउन्टेन क्लाइम्बर", "sets": 4, "reps": 20, "rest_seconds": 30, "muscle_groups": ["core", "hip_flexors", "shoulders"]},
        {"name": "Jump Squats", "name_ne": "जम्प स्क्वाट", "sets": 4, "reps": 15, "rest_seconds": 30, "muscle_groups": ["quadriceps", "glutes", "calves"]},
        {"name": "High Knees", "name_ne": "हाइ नीज", "sets": 4, "reps": 30, "rest_seconds": 30, "muscle_groups": ["core", "hip_flexors", "quadriceps"]},
        {"name": "Box Jumps", "name_ne": "बक्स जम्प", "sets": 4, "reps": 10, "rest_seconds": 30, "muscle_groups": ["quadriceps", "glutes", "calves"]}
    ]'::jsonb,
    goal = 'lose_weight',
    days_per_week = '4-5'
WHERE id = 'b1b2c3d4-e5f6-7890-abcd-ef1234567804';

-- 5. Core & Flexibility
UPDATE workout_templates SET
    exercises = '[
        {"name": "Plank", "name_ne": "प्ल्यांक", "sets": 3, "reps": 1, "rest_seconds": 45, "notes": "Hold for 45 seconds", "muscle_groups": ["core", "obliques"]},
        {"name": "Russian Twists", "name_ne": "रसियन ट्विस्ट", "sets": 3, "reps": 20, "rest_seconds": 45, "muscle_groups": ["core", "obliques"]},
        {"name": "Bicycle Crunches", "name_ne": "बाइसिकल क्रन्च", "sets": 3, "reps": 20, "rest_seconds": 45, "muscle_groups": ["core", "obliques"]},
        {"name": "Dead Bug", "name_ne": "डेड बग", "sets": 3, "reps": 10, "rest_seconds": 45, "muscle_groups": ["core", "hip_flexors"]},
        {"name": "Cat-Cow Stretch", "name_ne": "क्याट-काउ स्ट्रेच", "sets": 3, "reps": 10, "rest_seconds": 30, "muscle_groups": ["core", "lower_back"]},
        {"name": "Seated Forward Fold", "name_ne": "बसेर अगाडि झुकाइ", "sets": 2, "reps": 1, "rest_seconds": 30, "notes": "Hold for 30 seconds", "muscle_groups": ["hamstrings", "lower_back"]}
    ]'::jsonb,
    goal = 'stay_fit',
    days_per_week = '2-3'
WHERE id = 'b1b2c3d4-e5f6-7890-abcd-ef1234567805';
