-- Seed 5 workout programs with day-by-day structure using program_days.

-- ============================================================
-- Program 1: 30-Day HIIT Challenge
-- ============================================================
INSERT INTO workout_programs (name, name_ne, description, program_type, difficulty, duration_weeks, equipment, goal)
VALUES (
    '30-Day HIIT Challenge',
    '३०-दिने HIIT चुनौती',
    'An intense 4-week high-intensity interval training program designed to torch calories and improve cardiovascular fitness. No equipment needed — just determination.',
    'hiit', 'advanced', 4, 'none', 'lose_weight'
) ON CONFLICT DO NOTHING;

-- Week 1
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 1, 1, 'HIIT Blast A', false,
    '[{"exercise": "Jumping Jacks", "sets": 3, "reps": 30, "rest_seconds": 20},
      {"exercise": "Burpees", "sets": 3, "reps": 10, "rest_seconds": 30},
      {"exercise": "Mountain Climbers", "sets": 3, "reps": 20, "rest_seconds": 20},
      {"exercise": "Jump Squats", "sets": 3, "reps": 15, "rest_seconds": 30},
      {"exercise": "High Knees", "sets": 3, "reps": 30, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 1, 2, 'HIIT Blast B', false,
    '[{"exercise": "Burpees", "sets": 4, "reps": 10, "rest_seconds": 30},
      {"exercise": "Tuck Jumps", "sets": 3, "reps": 10, "rest_seconds": 30},
      {"exercise": "Skater Jumps", "sets": 3, "reps": 20, "rest_seconds": 20},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 45, "rest_seconds": 15},
      {"exercise": "Butt Kicks", "sets": 3, "reps": 30, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 1, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 1, 4, 'HIIT Core Burn', false,
    '[{"exercise": "Mountain Climbers", "sets": 4, "reps": 20, "rest_seconds": 20},
      {"exercise": "Bicycle Crunches", "sets": 3, "reps": 20, "rest_seconds": 20},
      {"exercise": "Russian Twists", "sets": 3, "reps": 20, "rest_seconds": 20},
      {"exercise": "Flutter Kicks", "sets": 3, "reps": 20, "rest_seconds": 20},
      {"exercise": "Burpees", "sets": 3, "reps": 12, "rest_seconds": 30}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 1, 5, 'HIIT Blast A', false,
    '[{"exercise": "Jumping Jacks", "sets": 3, "reps": 30, "rest_seconds": 20},
      {"exercise": "Burpees", "sets": 3, "reps": 10, "rest_seconds": 30},
      {"exercise": "Mountain Climbers", "sets": 3, "reps": 20, "rest_seconds": 20},
      {"exercise": "Jump Squats", "sets": 3, "reps": 15, "rest_seconds": 30},
      {"exercise": "High Knees", "sets": 3, "reps": 30, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 1, 6, 'HIIT Tabata', false,
    '[{"exercise": "Burpees", "sets": 4, "reps": 12, "rest_seconds": 20},
      {"exercise": "Jump Squats", "sets": 4, "reps": 15, "rest_seconds": 20},
      {"exercise": "High Knees", "sets": 4, "reps": 30, "rest_seconds": 20},
      {"exercise": "Skater Jumps", "sets": 4, "reps": 20, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 1, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;

-- Week 2 (increased intensity)
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 2, 1, 'HIIT Blast A+', false,
    '[{"exercise": "Jumping Jacks", "sets": 4, "reps": 35, "rest_seconds": 15},
      {"exercise": "Burpees", "sets": 4, "reps": 12, "rest_seconds": 25},
      {"exercise": "Mountain Climbers", "sets": 4, "reps": 25, "rest_seconds": 15},
      {"exercise": "Jump Squats", "sets": 4, "reps": 18, "rest_seconds": 25},
      {"exercise": "High Knees", "sets": 4, "reps": 35, "rest_seconds": 15}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 2, 2, 'HIIT Blast B+', false,
    '[{"exercise": "Burpees", "sets": 4, "reps": 12, "rest_seconds": 25},
      {"exercise": "Tuck Jumps", "sets": 4, "reps": 12, "rest_seconds": 25},
      {"exercise": "Skater Jumps", "sets": 4, "reps": 24, "rest_seconds": 15},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 60, "rest_seconds": 15},
      {"exercise": "Bear Crawl", "sets": 3, "duration_seconds": 30, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 2, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 2, 4, 'HIIT Core Burn+', false,
    '[{"exercise": "Mountain Climbers", "sets": 4, "reps": 25, "rest_seconds": 15},
      {"exercise": "Bicycle Crunches", "sets": 4, "reps": 25, "rest_seconds": 15},
      {"exercise": "V-ups", "sets": 3, "reps": 12, "rest_seconds": 20},
      {"exercise": "Flutter Kicks", "sets": 4, "reps": 25, "rest_seconds": 15},
      {"exercise": "Burpees", "sets": 4, "reps": 14, "rest_seconds": 25}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 2, 5, 'HIIT Blast A+', false,
    '[{"exercise": "Jumping Jacks", "sets": 4, "reps": 35, "rest_seconds": 15},
      {"exercise": "Burpees", "sets": 4, "reps": 12, "rest_seconds": 25},
      {"exercise": "Mountain Climbers", "sets": 4, "reps": 25, "rest_seconds": 15},
      {"exercise": "Jump Squats", "sets": 4, "reps": 18, "rest_seconds": 25},
      {"exercise": "High Knees", "sets": 4, "reps": 35, "rest_seconds": 15}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 2, 6, 'HIIT Tabata+', false,
    '[{"exercise": "Burpees", "sets": 5, "reps": 14, "rest_seconds": 15},
      {"exercise": "Jump Squats", "sets": 5, "reps": 18, "rest_seconds": 15},
      {"exercise": "High Knees", "sets": 5, "reps": 35, "rest_seconds": 15},
      {"exercise": "Tuck Jumps", "sets": 4, "reps": 12, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 2, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;

-- Week 3 (peak week)
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 3, 1, 'Peak HIIT A', false,
    '[{"exercise": "Burpees", "sets": 5, "reps": 15, "rest_seconds": 20},
      {"exercise": "Jump Squats", "sets": 4, "reps": 20, "rest_seconds": 20},
      {"exercise": "Mountain Climbers", "sets": 4, "reps": 30, "rest_seconds": 15},
      {"exercise": "Tuck Jumps", "sets": 4, "reps": 12, "rest_seconds": 20},
      {"exercise": "Skater Jumps", "sets": 4, "reps": 24, "rest_seconds": 15}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 3, 2, 'Peak HIIT B', false,
    '[{"exercise": "High Knees", "sets": 4, "reps": 40, "rest_seconds": 15},
      {"exercise": "Burpees", "sets": 5, "reps": 15, "rest_seconds": 20},
      {"exercise": "Box Jumps", "sets": 4, "reps": 12, "rest_seconds": 25},
      {"exercise": "Bear Crawl", "sets": 3, "duration_seconds": 40, "rest_seconds": 20},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 60, "rest_seconds": 15}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 3, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 3, 4, 'Peak Core', false,
    '[{"exercise": "V-ups", "sets": 4, "reps": 15, "rest_seconds": 20},
      {"exercise": "Russian Twists", "sets": 4, "reps": 25, "rest_seconds": 15},
      {"exercise": "Leg Raises", "sets": 4, "reps": 15, "rest_seconds": 20},
      {"exercise": "Mountain Climbers", "sets": 4, "reps": 30, "rest_seconds": 15},
      {"exercise": "Bicycle Crunches", "sets": 4, "reps": 30, "rest_seconds": 15}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 3, 5, 'Peak HIIT A', false,
    '[{"exercise": "Burpees", "sets": 5, "reps": 15, "rest_seconds": 20},
      {"exercise": "Jump Squats", "sets": 4, "reps": 20, "rest_seconds": 20},
      {"exercise": "Mountain Climbers", "sets": 4, "reps": 30, "rest_seconds": 15},
      {"exercise": "Tuck Jumps", "sets": 4, "reps": 12, "rest_seconds": 20},
      {"exercise": "Skater Jumps", "sets": 4, "reps": 24, "rest_seconds": 15}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 3, 6, 'Peak Tabata', false,
    '[{"exercise": "Burpees", "sets": 5, "reps": 15, "rest_seconds": 15},
      {"exercise": "Jump Squats", "sets": 5, "reps": 20, "rest_seconds": 15},
      {"exercise": "High Knees", "sets": 5, "reps": 40, "rest_seconds": 15},
      {"exercise": "Tuck Jumps", "sets": 5, "reps": 12, "rest_seconds": 15}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 3, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;

-- Week 4 (deload + finisher)
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 4, 1, 'Deload HIIT A', false,
    '[{"exercise": "Jumping Jacks", "sets": 3, "reps": 30, "rest_seconds": 25},
      {"exercise": "Burpees", "sets": 3, "reps": 10, "rest_seconds": 30},
      {"exercise": "Mountain Climbers", "sets": 3, "reps": 20, "rest_seconds": 25},
      {"exercise": "Jump Squats", "sets": 3, "reps": 12, "rest_seconds": 30}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 4, 2, 'Deload HIIT B', false,
    '[{"exercise": "High Knees", "sets": 3, "reps": 30, "rest_seconds": 25},
      {"exercise": "Skater Jumps", "sets": 3, "reps": 20, "rest_seconds": 25},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 45, "rest_seconds": 20},
      {"exercise": "Butt Kicks", "sets": 3, "reps": 30, "rest_seconds": 25}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 4, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 4, 4, 'Deload Core', false,
    '[{"exercise": "Bicycle Crunches", "sets": 3, "reps": 20, "rest_seconds": 20},
      {"exercise": "Russian Twists", "sets": 3, "reps": 20, "rest_seconds": 20},
      {"exercise": "Leg Raises", "sets": 3, "reps": 12, "rest_seconds": 20},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 45, "rest_seconds": 15}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 4, 5, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 4, 6, 'Final HIIT Blitz', false,
    '[{"exercise": "Burpees", "sets": 5, "reps": 15, "rest_seconds": 15},
      {"exercise": "Jump Squats", "sets": 5, "reps": 20, "rest_seconds": 15},
      {"exercise": "Mountain Climbers", "sets": 5, "reps": 30, "rest_seconds": 15},
      {"exercise": "Tuck Jumps", "sets": 4, "reps": 12, "rest_seconds": 20},
      {"exercise": "High Knees", "sets": 4, "reps": 40, "rest_seconds": 15},
      {"exercise": "Plank", "sets": 1, "duration_seconds": 120, "rest_seconds": 0}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = '30-Day HIIT Challenge'), 4, 7, 'Rest & Celebrate', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;


-- ============================================================
-- Program 2: Beginner Bodyweight
-- ============================================================
INSERT INTO workout_programs (name, name_ne, description, program_type, difficulty, duration_weeks, equipment, goal)
VALUES (
    'Beginner Bodyweight',
    'शुरुआती बडीवेट',
    'A gentle 4-week introduction to bodyweight training. Perfect for those new to exercise. Build a foundation of strength with simple, effective movements.',
    'bodyweight', 'beginner', 4, 'none', 'stay_fit'
) ON CONFLICT DO NOTHING;

-- Week 1
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 1, 1, 'Upper Body Basics', false,
    '[{"exercise": "Push-ups", "sets": 3, "reps": 8, "rest_seconds": 45},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 20, "rest_seconds": 30},
      {"exercise": "Tricep Dips (Bodyweight)", "sets": 2, "reps": 8, "rest_seconds": 45},
      {"exercise": "Superman", "sets": 3, "reps": 10, "rest_seconds": 30}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 1, 2, 'Lower Body Basics', false,
    '[{"exercise": "Squats", "sets": 3, "reps": 12, "rest_seconds": 45},
      {"exercise": "Lunges", "sets": 3, "reps": 10, "rest_seconds": 45},
      {"exercise": "Glute Bridges", "sets": 3, "reps": 12, "rest_seconds": 30},
      {"exercise": "Wall Sit", "sets": 2, "duration_seconds": 20, "rest_seconds": 30}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 1, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 1, 4, 'Core Foundations', false,
    '[{"exercise": "Plank", "sets": 3, "duration_seconds": 20, "rest_seconds": 30},
      {"exercise": "Bicycle Crunches", "sets": 3, "reps": 12, "rest_seconds": 30},
      {"exercise": "Leg Raises", "sets": 2, "reps": 10, "rest_seconds": 30},
      {"exercise": "Superman", "sets": 3, "reps": 10, "rest_seconds": 30}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 1, 5, 'Full Body Light', false,
    '[{"exercise": "Jumping Jacks", "sets": 3, "reps": 20, "rest_seconds": 30},
      {"exercise": "Push-ups", "sets": 2, "reps": 8, "rest_seconds": 45},
      {"exercise": "Squats", "sets": 3, "reps": 12, "rest_seconds": 45},
      {"exercise": "Plank", "sets": 2, "duration_seconds": 20, "rest_seconds": 30}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 1, 6, 'Rest Day', true, '[]'::jsonb),
((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 1, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;

-- Week 2
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 2, 1, 'Upper Body Basics', false,
    '[{"exercise": "Push-ups", "sets": 3, "reps": 10, "rest_seconds": 45},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 25, "rest_seconds": 30},
      {"exercise": "Tricep Dips (Bodyweight)", "sets": 3, "reps": 8, "rest_seconds": 45},
      {"exercise": "Inchworm", "sets": 3, "reps": 6, "rest_seconds": 30}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 2, 2, 'Lower Body Basics', false,
    '[{"exercise": "Squats", "sets": 3, "reps": 15, "rest_seconds": 45},
      {"exercise": "Reverse Lunges", "sets": 3, "reps": 10, "rest_seconds": 45},
      {"exercise": "Glute Bridges", "sets": 3, "reps": 15, "rest_seconds": 30},
      {"exercise": "Wall Sit", "sets": 3, "duration_seconds": 25, "rest_seconds": 30}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 2, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 2, 4, 'Core Foundations', false,
    '[{"exercise": "Plank", "sets": 3, "duration_seconds": 30, "rest_seconds": 25},
      {"exercise": "Bicycle Crunches", "sets": 3, "reps": 15, "rest_seconds": 25},
      {"exercise": "Leg Raises", "sets": 3, "reps": 10, "rest_seconds": 30},
      {"exercise": "Side Plank", "sets": 2, "duration_seconds": 15, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 2, 5, 'Full Body Light', false,
    '[{"exercise": "Jumping Jacks", "sets": 3, "reps": 25, "rest_seconds": 25},
      {"exercise": "Push-ups", "sets": 3, "reps": 10, "rest_seconds": 45},
      {"exercise": "Squats", "sets": 3, "reps": 15, "rest_seconds": 45},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 25, "rest_seconds": 25}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 2, 6, 'Rest Day', true, '[]'::jsonb),
((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 2, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;

-- Week 3
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 3, 1, 'Upper Body Progress', false,
    '[{"exercise": "Wide Push-ups", "sets": 3, "reps": 8, "rest_seconds": 45},
      {"exercise": "Push-ups", "sets": 3, "reps": 10, "rest_seconds": 40},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 30, "rest_seconds": 25},
      {"exercise": "Tricep Dips (Bodyweight)", "sets": 3, "reps": 10, "rest_seconds": 40}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 3, 2, 'Lower Body Progress', false,
    '[{"exercise": "Squats", "sets": 3, "reps": 18, "rest_seconds": 40},
      {"exercise": "Lunges", "sets": 3, "reps": 12, "rest_seconds": 40},
      {"exercise": "Glute Bridges", "sets": 3, "reps": 18, "rest_seconds": 25},
      {"exercise": "Wall Sit", "sets": 3, "duration_seconds": 30, "rest_seconds": 25}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 3, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 3, 4, 'Core Progress', false,
    '[{"exercise": "Plank", "sets": 3, "duration_seconds": 35, "rest_seconds": 25},
      {"exercise": "Bicycle Crunches", "sets": 3, "reps": 18, "rest_seconds": 25},
      {"exercise": "Russian Twists", "sets": 3, "reps": 15, "rest_seconds": 25},
      {"exercise": "Side Plank", "sets": 3, "duration_seconds": 20, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 3, 5, 'Full Body Progress', false,
    '[{"exercise": "Jumping Jacks", "sets": 3, "reps": 30, "rest_seconds": 20},
      {"exercise": "Push-ups", "sets": 3, "reps": 12, "rest_seconds": 40},
      {"exercise": "Squats", "sets": 3, "reps": 18, "rest_seconds": 40},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 30, "rest_seconds": 20},
      {"exercise": "Glute Bridges", "sets": 3, "reps": 15, "rest_seconds": 25}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 3, 6, 'Rest Day', true, '[]'::jsonb),
((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 3, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;

-- Week 4
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 4, 1, 'Upper Body Challenge', false,
    '[{"exercise": "Wide Push-ups", "sets": 3, "reps": 10, "rest_seconds": 40},
      {"exercise": "Push-ups", "sets": 3, "reps": 12, "rest_seconds": 40},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 40, "rest_seconds": 20},
      {"exercise": "Tricep Dips (Bodyweight)", "sets": 3, "reps": 12, "rest_seconds": 40},
      {"exercise": "Inchworm", "sets": 3, "reps": 8, "rest_seconds": 30}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 4, 2, 'Lower Body Challenge', false,
    '[{"exercise": "Squats", "sets": 3, "reps": 20, "rest_seconds": 40},
      {"exercise": "Lunges", "sets": 3, "reps": 14, "rest_seconds": 40},
      {"exercise": "Reverse Lunges", "sets": 3, "reps": 12, "rest_seconds": 40},
      {"exercise": "Glute Bridges", "sets": 3, "reps": 20, "rest_seconds": 25},
      {"exercise": "Wall Sit", "sets": 3, "duration_seconds": 35, "rest_seconds": 25}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 4, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 4, 4, 'Core Challenge', false,
    '[{"exercise": "Plank", "sets": 3, "duration_seconds": 40, "rest_seconds": 20},
      {"exercise": "Bicycle Crunches", "sets": 3, "reps": 20, "rest_seconds": 20},
      {"exercise": "Russian Twists", "sets": 3, "reps": 20, "rest_seconds": 20},
      {"exercise": "Leg Raises", "sets": 3, "reps": 12, "rest_seconds": 25},
      {"exercise": "Side Plank", "sets": 3, "duration_seconds": 25, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 4, 5, 'Full Body Finisher', false,
    '[{"exercise": "Jumping Jacks", "sets": 3, "reps": 30, "rest_seconds": 20},
      {"exercise": "Push-ups", "sets": 3, "reps": 15, "rest_seconds": 40},
      {"exercise": "Squats", "sets": 3, "reps": 20, "rest_seconds": 40},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 40, "rest_seconds": 20},
      {"exercise": "Glute Bridges", "sets": 3, "reps": 20, "rest_seconds": 25},
      {"exercise": "Burpees", "sets": 2, "reps": 5, "rest_seconds": 45}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 4, 6, 'Rest Day', true, '[]'::jsonb),
((SELECT id FROM workout_programs WHERE name = 'Beginner Bodyweight'), 4, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;


-- ============================================================
-- Program 3: Core Crusher (2 weeks)
-- ============================================================
INSERT INTO workout_programs (name, name_ne, description, program_type, difficulty, duration_weeks, equipment, goal)
VALUES (
    'Core Crusher',
    'कोर क्रशर',
    'A focused 2-week core strengthening program. Target your abs, obliques, and lower back with progressive overload. No equipment needed.',
    'bodyweight', 'intermediate', 2, 'none', 'stay_fit'
) ON CONFLICT DO NOTHING;

-- Week 1
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = 'Core Crusher'), 1, 1, 'Core Blast', false,
    '[{"exercise": "Plank", "sets": 3, "duration_seconds": 40, "rest_seconds": 20},
      {"exercise": "Bicycle Crunches", "sets": 3, "reps": 20, "rest_seconds": 20},
      {"exercise": "Russian Twists", "sets": 3, "reps": 20, "rest_seconds": 20},
      {"exercise": "Leg Raises", "sets": 3, "reps": 12, "rest_seconds": 25},
      {"exercise": "Flutter Kicks", "sets": 3, "reps": 20, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Core Crusher'), 1, 2, 'Oblique Focus', false,
    '[{"exercise": "Side Plank", "sets": 3, "duration_seconds": 25, "rest_seconds": 15},
      {"exercise": "Russian Twists", "sets": 4, "reps": 25, "rest_seconds": 20},
      {"exercise": "Bicycle Crunches", "sets": 3, "reps": 20, "rest_seconds": 20},
      {"exercise": "Mountain Climbers", "sets": 3, "reps": 20, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Core Crusher'), 1, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Core Crusher'), 1, 4, 'Lower Abs', false,
    '[{"exercise": "Leg Raises", "sets": 4, "reps": 15, "rest_seconds": 20},
      {"exercise": "Flutter Kicks", "sets": 3, "reps": 25, "rest_seconds": 20},
      {"exercise": "V-ups", "sets": 3, "reps": 10, "rest_seconds": 25},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 45, "rest_seconds": 15}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Core Crusher'), 1, 5, 'Core Endurance', false,
    '[{"exercise": "Plank", "sets": 4, "duration_seconds": 45, "rest_seconds": 15},
      {"exercise": "Side Plank", "sets": 3, "duration_seconds": 25, "rest_seconds": 15},
      {"exercise": "Superman", "sets": 3, "reps": 15, "rest_seconds": 20},
      {"exercise": "Glute Bridges", "sets": 3, "reps": 15, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Core Crusher'), 1, 6, 'Rest Day', true, '[]'::jsonb),
((SELECT id FROM workout_programs WHERE name = 'Core Crusher'), 1, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;

-- Week 2
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = 'Core Crusher'), 2, 1, 'Core Blast+', false,
    '[{"exercise": "Plank", "sets": 4, "duration_seconds": 50, "rest_seconds": 15},
      {"exercise": "Bicycle Crunches", "sets": 4, "reps": 25, "rest_seconds": 15},
      {"exercise": "V-ups", "sets": 3, "reps": 12, "rest_seconds": 20},
      {"exercise": "Leg Raises", "sets": 4, "reps": 15, "rest_seconds": 20},
      {"exercise": "Flutter Kicks", "sets": 4, "reps": 25, "rest_seconds": 15}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Core Crusher'), 2, 2, 'Oblique Focus+', false,
    '[{"exercise": "Side Plank", "sets": 4, "duration_seconds": 30, "rest_seconds": 15},
      {"exercise": "Russian Twists", "sets": 4, "reps": 30, "rest_seconds": 15},
      {"exercise": "Bicycle Crunches", "sets": 4, "reps": 25, "rest_seconds": 15},
      {"exercise": "Mountain Climbers", "sets": 4, "reps": 25, "rest_seconds": 15}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Core Crusher'), 2, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Core Crusher'), 2, 4, 'Lower Abs+', false,
    '[{"exercise": "Leg Raises", "sets": 4, "reps": 18, "rest_seconds": 15},
      {"exercise": "Flutter Kicks", "sets": 4, "reps": 30, "rest_seconds": 15},
      {"exercise": "V-ups", "sets": 4, "reps": 15, "rest_seconds": 20},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 60, "rest_seconds": 15}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Core Crusher'), 2, 5, 'Core Finisher', false,
    '[{"exercise": "Plank", "sets": 3, "duration_seconds": 60, "rest_seconds": 15},
      {"exercise": "Side Plank", "sets": 3, "duration_seconds": 35, "rest_seconds": 15},
      {"exercise": "V-ups", "sets": 4, "reps": 15, "rest_seconds": 20},
      {"exercise": "Russian Twists", "sets": 4, "reps": 30, "rest_seconds": 15},
      {"exercise": "Bicycle Crunches", "sets": 4, "reps": 30, "rest_seconds": 15},
      {"exercise": "Superman", "sets": 3, "reps": 15, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Core Crusher'), 2, 6, 'Rest Day', true, '[]'::jsonb),
((SELECT id FROM workout_programs WHERE name = 'Core Crusher'), 2, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;


-- ============================================================
-- Program 4: Full Body Strength (dumbbells)
-- ============================================================
INSERT INTO workout_programs (name, name_ne, description, program_type, difficulty, duration_weeks, equipment, goal)
VALUES (
    'Full Body Strength',
    'पूर्ण शरीर शक्ति',
    'A 4-week dumbbell-based strength program targeting all major muscle groups. Requires a pair of dumbbells. Progressive overload built in week over week.',
    'strength', 'intermediate', 4, 'dumbbells', 'build_muscle'
) ON CONFLICT DO NOTHING;

-- Week 1
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 1, 1, 'Push Day', false,
    '[{"exercise": "Push-ups", "sets": 3, "reps": 12, "rest_seconds": 45},
      {"exercise": "Dumbbell Bench Press", "sets": 3, "reps": 10, "rest_seconds": 60},
      {"exercise": "Dumbbell Shoulder Press", "sets": 3, "reps": 10, "rest_seconds": 60},
      {"exercise": "Lateral Raises", "sets": 3, "reps": 12, "rest_seconds": 45},
      {"exercise": "Tricep Extensions", "sets": 3, "reps": 12, "rest_seconds": 45}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 1, 2, 'Pull Day', false,
    '[{"exercise": "Dumbbell Rows", "sets": 3, "reps": 10, "rest_seconds": 60},
      {"exercise": "Bicep Curls", "sets": 3, "reps": 12, "rest_seconds": 45},
      {"exercise": "Dumbbell Deadlift", "sets": 3, "reps": 10, "rest_seconds": 60},
      {"exercise": "Superman", "sets": 3, "reps": 15, "rest_seconds": 30}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 1, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 1, 4, 'Legs Day', false,
    '[{"exercise": "Goblet Squats", "sets": 3, "reps": 12, "rest_seconds": 60},
      {"exercise": "Dumbbell Lunges", "sets": 3, "reps": 10, "rest_seconds": 60},
      {"exercise": "Dumbbell Deadlift", "sets": 3, "reps": 10, "rest_seconds": 60},
      {"exercise": "Glute Bridges", "sets": 3, "reps": 15, "rest_seconds": 45},
      {"exercise": "Wall Sit", "sets": 3, "duration_seconds": 30, "rest_seconds": 30}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 1, 5, 'Full Body', false,
    '[{"exercise": "Dumbbell Bench Press", "sets": 3, "reps": 10, "rest_seconds": 60},
      {"exercise": "Dumbbell Rows", "sets": 3, "reps": 10, "rest_seconds": 60},
      {"exercise": "Goblet Squats", "sets": 3, "reps": 12, "rest_seconds": 60},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 40, "rest_seconds": 20},
      {"exercise": "Bicep Curls", "sets": 3, "reps": 12, "rest_seconds": 45}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 1, 6, 'Rest Day', true, '[]'::jsonb),
((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 1, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;

-- Week 2 (increased reps)
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 2, 1, 'Push Day', false,
    '[{"exercise": "Push-ups", "sets": 3, "reps": 15, "rest_seconds": 40},
      {"exercise": "Dumbbell Bench Press", "sets": 4, "reps": 10, "rest_seconds": 60},
      {"exercise": "Dumbbell Shoulder Press", "sets": 4, "reps": 10, "rest_seconds": 60},
      {"exercise": "Lateral Raises", "sets": 3, "reps": 15, "rest_seconds": 40},
      {"exercise": "Tricep Extensions", "sets": 3, "reps": 15, "rest_seconds": 40}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 2, 2, 'Pull Day', false,
    '[{"exercise": "Dumbbell Rows", "sets": 4, "reps": 10, "rest_seconds": 60},
      {"exercise": "Bicep Curls", "sets": 3, "reps": 15, "rest_seconds": 40},
      {"exercise": "Dumbbell Deadlift", "sets": 4, "reps": 10, "rest_seconds": 60},
      {"exercise": "Superman", "sets": 3, "reps": 18, "rest_seconds": 25}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 2, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 2, 4, 'Legs Day', false,
    '[{"exercise": "Goblet Squats", "sets": 4, "reps": 12, "rest_seconds": 60},
      {"exercise": "Dumbbell Lunges", "sets": 3, "reps": 12, "rest_seconds": 60},
      {"exercise": "Dumbbell Deadlift", "sets": 3, "reps": 12, "rest_seconds": 60},
      {"exercise": "Glute Bridges", "sets": 3, "reps": 18, "rest_seconds": 40},
      {"exercise": "Wall Sit", "sets": 3, "duration_seconds": 40, "rest_seconds": 25}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 2, 5, 'Full Body', false,
    '[{"exercise": "Dumbbell Bench Press", "sets": 4, "reps": 10, "rest_seconds": 60},
      {"exercise": "Dumbbell Rows", "sets": 4, "reps": 10, "rest_seconds": 60},
      {"exercise": "Goblet Squats", "sets": 4, "reps": 12, "rest_seconds": 60},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 45, "rest_seconds": 15},
      {"exercise": "Lateral Raises", "sets": 3, "reps": 12, "rest_seconds": 40}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 2, 6, 'Rest Day', true, '[]'::jsonb),
((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 2, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;

-- Week 3 (increased sets)
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 3, 1, 'Push Day', false,
    '[{"exercise": "Diamond Push-ups", "sets": 3, "reps": 10, "rest_seconds": 45},
      {"exercise": "Dumbbell Bench Press", "sets": 4, "reps": 12, "rest_seconds": 60},
      {"exercise": "Dumbbell Shoulder Press", "sets": 4, "reps": 12, "rest_seconds": 60},
      {"exercise": "Lateral Raises", "sets": 4, "reps": 12, "rest_seconds": 40},
      {"exercise": "Tricep Extensions", "sets": 4, "reps": 12, "rest_seconds": 40}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 3, 2, 'Pull Day', false,
    '[{"exercise": "Dumbbell Rows", "sets": 4, "reps": 12, "rest_seconds": 60},
      {"exercise": "Bicep Curls", "sets": 4, "reps": 12, "rest_seconds": 40},
      {"exercise": "Dumbbell Deadlift", "sets": 4, "reps": 12, "rest_seconds": 60},
      {"exercise": "Superman", "sets": 3, "reps": 20, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 3, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 3, 4, 'Legs Day', false,
    '[{"exercise": "Goblet Squats", "sets": 4, "reps": 15, "rest_seconds": 60},
      {"exercise": "Dumbbell Lunges", "sets": 4, "reps": 12, "rest_seconds": 60},
      {"exercise": "Dumbbell Deadlift", "sets": 4, "reps": 12, "rest_seconds": 60},
      {"exercise": "Glute Bridges", "sets": 4, "reps": 20, "rest_seconds": 30},
      {"exercise": "Wall Sit", "sets": 3, "duration_seconds": 45, "rest_seconds": 25}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 3, 5, 'Full Body', false,
    '[{"exercise": "Dumbbell Bench Press", "sets": 4, "reps": 12, "rest_seconds": 55},
      {"exercise": "Dumbbell Rows", "sets": 4, "reps": 12, "rest_seconds": 55},
      {"exercise": "Goblet Squats", "sets": 4, "reps": 15, "rest_seconds": 55},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 50, "rest_seconds": 15},
      {"exercise": "Dumbbell Shoulder Press", "sets": 3, "reps": 12, "rest_seconds": 50}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 3, 6, 'Rest Day', true, '[]'::jsonb),
((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 3, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;

-- Week 4 (peak + test)
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 4, 1, 'Push Day Peak', false,
    '[{"exercise": "Decline Push-ups", "sets": 3, "reps": 10, "rest_seconds": 45},
      {"exercise": "Dumbbell Bench Press", "sets": 4, "reps": 12, "rest_seconds": 55},
      {"exercise": "Dumbbell Shoulder Press", "sets": 4, "reps": 12, "rest_seconds": 55},
      {"exercise": "Lateral Raises", "sets": 4, "reps": 15, "rest_seconds": 35},
      {"exercise": "Tricep Extensions", "sets": 4, "reps": 15, "rest_seconds": 35}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 4, 2, 'Pull Day Peak', false,
    '[{"exercise": "Dumbbell Rows", "sets": 4, "reps": 12, "rest_seconds": 55},
      {"exercise": "Bicep Curls", "sets": 4, "reps": 15, "rest_seconds": 35},
      {"exercise": "Dumbbell Deadlift", "sets": 4, "reps": 12, "rest_seconds": 55},
      {"exercise": "Superman", "sets": 3, "reps": 20, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 4, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 4, 4, 'Legs Day Peak', false,
    '[{"exercise": "Goblet Squats", "sets": 4, "reps": 15, "rest_seconds": 55},
      {"exercise": "Dumbbell Lunges", "sets": 4, "reps": 14, "rest_seconds": 55},
      {"exercise": "Dumbbell Deadlift", "sets": 4, "reps": 14, "rest_seconds": 55},
      {"exercise": "Glute Bridges", "sets": 4, "reps": 20, "rest_seconds": 30},
      {"exercise": "Wall Sit", "sets": 3, "duration_seconds": 50, "rest_seconds": 20}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 4, 5, 'Full Body Finisher', false,
    '[{"exercise": "Dumbbell Bench Press", "sets": 4, "reps": 12, "rest_seconds": 50},
      {"exercise": "Dumbbell Rows", "sets": 4, "reps": 12, "rest_seconds": 50},
      {"exercise": "Goblet Squats", "sets": 4, "reps": 15, "rest_seconds": 50},
      {"exercise": "Dumbbell Shoulder Press", "sets": 4, "reps": 12, "rest_seconds": 50},
      {"exercise": "Plank", "sets": 3, "duration_seconds": 60, "rest_seconds": 15},
      {"exercise": "Bicep Curls", "sets": 3, "reps": 15, "rest_seconds": 35}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 4, 6, 'Rest Day', true, '[]'::jsonb),
((SELECT id FROM workout_programs WHERE name = 'Full Body Strength'), 4, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;


-- ============================================================
-- Program 5: Flexibility & Mobility (3 weeks)
-- ============================================================
INSERT INTO workout_programs (name, name_ne, description, program_type, difficulty, duration_weeks, equipment, goal)
VALUES (
    'Flexibility & Mobility',
    'लचिलोपन र गतिशीलता',
    'A gentle 3-week program focused on improving flexibility, mobility, and recovery. Perfect as a standalone routine or complement to strength training. No equipment needed.',
    'flexibility', 'beginner', 3, 'none', 'stay_fit'
) ON CONFLICT DO NOTHING;

-- Week 1
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 1, 1, 'Full Body Stretch', false,
    '[{"exercise": "Cat-Cow Stretch", "sets": 3, "reps": 10, "rest_seconds": 10},
      {"exercise": "Downward Dog", "sets": 3, "duration_seconds": 30, "rest_seconds": 10},
      {"exercise": "Child''s Pose", "sets": 2, "duration_seconds": 30, "rest_seconds": 10},
      {"exercise": "Forward Fold", "sets": 3, "duration_seconds": 30, "rest_seconds": 10},
      {"exercise": "Cobra Stretch", "sets": 3, "duration_seconds": 20, "rest_seconds": 10}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 1, 2, 'Lower Body Mobility', false,
    '[{"exercise": "Forward Fold", "sets": 3, "duration_seconds": 30, "rest_seconds": 10},
      {"exercise": "Pigeon Stretch", "sets": 3, "duration_seconds": 30, "rest_seconds": 10},
      {"exercise": "Glute Bridges", "sets": 3, "reps": 12, "rest_seconds": 20},
      {"exercise": "Squats", "sets": 2, "reps": 10, "rest_seconds": 20},
      {"exercise": "Downward Dog", "sets": 2, "duration_seconds": 30, "rest_seconds": 10}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 1, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 1, 4, 'Upper Body & Spine', false,
    '[{"exercise": "Cat-Cow Stretch", "sets": 3, "reps": 10, "rest_seconds": 10},
      {"exercise": "Cobra Stretch", "sets": 3, "duration_seconds": 25, "rest_seconds": 10},
      {"exercise": "Child''s Pose", "sets": 3, "duration_seconds": 30, "rest_seconds": 10},
      {"exercise": "Downward Dog", "sets": 3, "duration_seconds": 30, "rest_seconds": 10}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 1, 5, 'Full Body Flow', false,
    '[{"exercise": "Cat-Cow Stretch", "sets": 2, "reps": 10, "rest_seconds": 10},
      {"exercise": "Downward Dog", "sets": 3, "duration_seconds": 30, "rest_seconds": 10},
      {"exercise": "Pigeon Stretch", "sets": 2, "duration_seconds": 30, "rest_seconds": 10},
      {"exercise": "Forward Fold", "sets": 3, "duration_seconds": 30, "rest_seconds": 10},
      {"exercise": "Cobra Stretch", "sets": 2, "duration_seconds": 25, "rest_seconds": 10},
      {"exercise": "Child''s Pose", "sets": 2, "duration_seconds": 30, "rest_seconds": 10}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 1, 6, 'Rest Day', true, '[]'::jsonb),
((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 1, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;

-- Week 2 (longer holds)
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 2, 1, 'Full Body Stretch', false,
    '[{"exercise": "Cat-Cow Stretch", "sets": 3, "reps": 12, "rest_seconds": 10},
      {"exercise": "Downward Dog", "sets": 3, "duration_seconds": 40, "rest_seconds": 10},
      {"exercise": "Child''s Pose", "sets": 3, "duration_seconds": 40, "rest_seconds": 10},
      {"exercise": "Forward Fold", "sets": 3, "duration_seconds": 40, "rest_seconds": 10},
      {"exercise": "Cobra Stretch", "sets": 3, "duration_seconds": 30, "rest_seconds": 10}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 2, 2, 'Lower Body Mobility', false,
    '[{"exercise": "Forward Fold", "sets": 3, "duration_seconds": 40, "rest_seconds": 10},
      {"exercise": "Pigeon Stretch", "sets": 3, "duration_seconds": 40, "rest_seconds": 10},
      {"exercise": "Glute Bridges", "sets": 3, "reps": 15, "rest_seconds": 15},
      {"exercise": "Lunges", "sets": 2, "reps": 10, "rest_seconds": 20},
      {"exercise": "Downward Dog", "sets": 3, "duration_seconds": 40, "rest_seconds": 10}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 2, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 2, 4, 'Upper Body & Spine', false,
    '[{"exercise": "Cat-Cow Stretch", "sets": 3, "reps": 12, "rest_seconds": 10},
      {"exercise": "Cobra Stretch", "sets": 3, "duration_seconds": 30, "rest_seconds": 10},
      {"exercise": "Side Plank", "sets": 2, "duration_seconds": 20, "rest_seconds": 15},
      {"exercise": "Child''s Pose", "sets": 3, "duration_seconds": 40, "rest_seconds": 10},
      {"exercise": "Inchworm", "sets": 3, "reps": 6, "rest_seconds": 15}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 2, 5, 'Full Body Flow', false,
    '[{"exercise": "Cat-Cow Stretch", "sets": 3, "reps": 12, "rest_seconds": 10},
      {"exercise": "Downward Dog", "sets": 3, "duration_seconds": 40, "rest_seconds": 10},
      {"exercise": "Pigeon Stretch", "sets": 3, "duration_seconds": 40, "rest_seconds": 10},
      {"exercise": "Forward Fold", "sets": 3, "duration_seconds": 40, "rest_seconds": 10},
      {"exercise": "Cobra Stretch", "sets": 3, "duration_seconds": 30, "rest_seconds": 10},
      {"exercise": "Child''s Pose", "sets": 3, "duration_seconds": 40, "rest_seconds": 10}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 2, 6, 'Rest Day', true, '[]'::jsonb),
((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 2, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;

-- Week 3 (deeper stretches + active recovery)
INSERT INTO program_days (program_id, week_number, day_number, day_name, rest_day, exercises) VALUES
((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 3, 1, 'Deep Stretch A', false,
    '[{"exercise": "Cat-Cow Stretch", "sets": 3, "reps": 15, "rest_seconds": 10},
      {"exercise": "Downward Dog", "sets": 3, "duration_seconds": 45, "rest_seconds": 10},
      {"exercise": "Pigeon Stretch", "sets": 3, "duration_seconds": 45, "rest_seconds": 10},
      {"exercise": "Forward Fold", "sets": 3, "duration_seconds": 45, "rest_seconds": 10},
      {"exercise": "Child''s Pose", "sets": 3, "duration_seconds": 45, "rest_seconds": 10}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 3, 2, 'Active Recovery', false,
    '[{"exercise": "Inchworm", "sets": 3, "reps": 8, "rest_seconds": 15},
      {"exercise": "Glute Bridges", "sets": 3, "reps": 15, "rest_seconds": 15},
      {"exercise": "Squats", "sets": 2, "reps": 12, "rest_seconds": 20},
      {"exercise": "Cat-Cow Stretch", "sets": 3, "reps": 10, "rest_seconds": 10},
      {"exercise": "Cobra Stretch", "sets": 3, "duration_seconds": 30, "rest_seconds": 10}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 3, 3, 'Rest Day', true, '[]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 3, 4, 'Deep Stretch B', false,
    '[{"exercise": "Forward Fold", "sets": 3, "duration_seconds": 45, "rest_seconds": 10},
      {"exercise": "Pigeon Stretch", "sets": 3, "duration_seconds": 50, "rest_seconds": 10},
      {"exercise": "Cobra Stretch", "sets": 3, "duration_seconds": 35, "rest_seconds": 10},
      {"exercise": "Downward Dog", "sets": 3, "duration_seconds": 45, "rest_seconds": 10},
      {"exercise": "Side Plank", "sets": 2, "duration_seconds": 25, "rest_seconds": 15}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 3, 5, 'Final Flow', false,
    '[{"exercise": "Cat-Cow Stretch", "sets": 3, "reps": 15, "rest_seconds": 10},
      {"exercise": "Downward Dog", "sets": 3, "duration_seconds": 45, "rest_seconds": 10},
      {"exercise": "Pigeon Stretch", "sets": 3, "duration_seconds": 50, "rest_seconds": 10},
      {"exercise": "Forward Fold", "sets": 3, "duration_seconds": 50, "rest_seconds": 10},
      {"exercise": "Cobra Stretch", "sets": 3, "duration_seconds": 35, "rest_seconds": 10},
      {"exercise": "Child''s Pose", "sets": 3, "duration_seconds": 60, "rest_seconds": 0}]'::jsonb),

((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 3, 6, 'Rest Day', true, '[]'::jsonb),
((SELECT id FROM workout_programs WHERE name = 'Flexibility & Mobility'), 3, 7, 'Rest Day', true, '[]'::jsonb)

ON CONFLICT DO NOTHING;
