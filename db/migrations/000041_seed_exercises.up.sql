-- Seed exercises: 50+ exercises covering bodyweight, dumbbell, barbell, and flexibility.

INSERT INTO exercises (name, name_ne, muscle_groups, equipment, exercise_type, difficulty, calories_per_minute, instructions) VALUES

-- ============================================================
-- Bodyweight / HIIT (no equipment)
-- ============================================================

('Push-ups', 'पुश-अप्स',
    ARRAY['chest', 'triceps', 'shoulders'], 'none', 'strength', 'beginner', 7.0,
    'Start in a high plank position with hands shoulder-width apart. Lower your body until your chest nearly touches the floor. Push back up to the starting position.'),

('Wide Push-ups', 'वाइड पुश-अप्स',
    ARRAY['chest', 'shoulders'], 'none', 'strength', 'beginner', 7.0,
    'Perform a push-up with hands placed wider than shoulder-width apart. This emphasizes the chest more than standard push-ups.'),

('Diamond Push-ups', 'डायमण्ड पुश-अप्स',
    ARRAY['triceps', 'chest'], 'none', 'strength', 'intermediate', 7.5,
    'Place hands close together under your chest forming a diamond shape with thumbs and index fingers. Lower and push back up.'),

('Decline Push-ups', 'डिक्लाइन पुश-अप्स',
    ARRAY['upper_chest', 'shoulders', 'triceps'], 'none', 'strength', 'intermediate', 7.5,
    'Place feet on an elevated surface and hands on the floor. Perform push-ups in this declined position to target upper chest.'),

('Squats', 'स्क्वाट्स',
    ARRAY['quadriceps', 'glutes', 'hamstrings'], 'none', 'strength', 'beginner', 6.0,
    'Stand with feet shoulder-width apart. Bend knees and lower hips as if sitting in a chair. Keep your back straight and return to standing.'),

('Jump Squats', 'जम्प स्क्वाट्स',
    ARRAY['quadriceps', 'glutes', 'calves'], 'none', 'hiit', 'intermediate', 10.0,
    'Perform a squat, then explosively jump upward. Land softly and immediately go into the next squat.'),

('Lunges', 'लन्जेस',
    ARRAY['quadriceps', 'glutes', 'hamstrings'], 'none', 'strength', 'beginner', 6.5,
    'Step forward with one leg, lowering your hips until both knees are bent at 90 degrees. Push back to standing and alternate legs.'),

('Reverse Lunges', 'रिभर्स लन्जेस',
    ARRAY['quadriceps', 'glutes', 'hamstrings'], 'none', 'strength', 'beginner', 6.5,
    'Step backward with one leg, lowering your hips until both knees are at 90 degrees. Return to standing and alternate legs.'),

('Burpees', 'बर्पीज',
    ARRAY['full_body'], 'none', 'hiit', 'intermediate', 12.0,
    'From standing, drop into a squat with hands on the floor. Jump feet back to plank, do a push-up, jump feet forward, then leap up with arms overhead.'),

('Mountain Climbers', 'माउन्टेन क्लाइम्बर्स',
    ARRAY['core', 'shoulders', 'quadriceps'], 'none', 'hiit', 'intermediate', 11.0,
    'Start in a high plank position. Rapidly alternate driving knees toward your chest as if running in place.'),

('Plank', 'प्ल्यांक',
    ARRAY['core', 'shoulders'], 'none', 'strength', 'beginner', 4.0,
    'Hold a push-up position with arms straight or on forearms. Keep your body in a straight line from head to heels. Hold for time.'),

('Side Plank', 'साइड प्ल्यांक',
    ARRAY['obliques', 'core'], 'none', 'strength', 'beginner', 4.0,
    'Lie on your side, propped up on one forearm. Lift hips off the ground forming a straight line. Hold for time on each side.'),

('Jumping Jacks', 'जम्पिंग ज्याक्स',
    ARRAY['full_body', 'calves'], 'none', 'cardio', 'beginner', 8.0,
    'Stand with feet together, arms at sides. Jump feet apart while raising arms overhead. Jump back to start and repeat.'),

('High Knees', 'हाई नीज',
    ARRAY['core', 'quadriceps', 'calves'], 'none', 'hiit', 'beginner', 9.5,
    'Run in place, driving knees as high as possible toward your chest. Pump your arms for momentum.'),

('Butt Kicks', 'बट किक्स',
    ARRAY['hamstrings', 'calves'], 'none', 'cardio', 'beginner', 8.5,
    'Run in place, kicking your heels up toward your glutes with each step. Keep a quick pace.'),

('Bicycle Crunches', 'बाइसिकल क्रन्चेस',
    ARRAY['core', 'obliques'], 'none', 'strength', 'beginner', 5.0,
    'Lie on your back with hands behind your head. Alternate bringing opposite elbow to knee in a cycling motion.'),

('Russian Twists', 'रसियन ट्विस्ट्स',
    ARRAY['obliques', 'core'], 'none', 'strength', 'intermediate', 5.5,
    'Sit on the floor with knees bent, lean back slightly. Rotate your torso side to side, optionally holding a weight.'),

('Leg Raises', 'लेग रेजेस',
    ARRAY['lower_abs', 'hip_flexors'], 'none', 'strength', 'intermediate', 5.0,
    'Lie flat on your back. Keep legs straight and lift them to 90 degrees, then slowly lower back down without touching the floor.'),

('Flutter Kicks', 'फ्लटर किक्स',
    ARRAY['lower_abs', 'hip_flexors'], 'none', 'strength', 'intermediate', 5.5,
    'Lie on your back. Lift both legs slightly off the ground and alternate kicking up and down in small, rapid motions.'),

('V-ups', 'भी-अप्स',
    ARRAY['core', 'hip_flexors'], 'none', 'strength', 'advanced', 6.0,
    'Lie flat on your back with arms extended overhead. Simultaneously lift legs and upper body to touch toes, forming a V shape.'),

('Superman', 'सुपरम्यान',
    ARRAY['lower_back', 'glutes'], 'none', 'strength', 'beginner', 3.5,
    'Lie face down with arms extended forward. Simultaneously lift arms, chest, and legs off the ground. Hold briefly, then lower.'),

('Glute Bridges', 'ग्लुट ब्रिजेस',
    ARRAY['glutes', 'hamstrings'], 'none', 'strength', 'beginner', 4.5,
    'Lie on your back with knees bent and feet flat. Push hips upward squeezing glutes at the top. Lower and repeat.'),

('Wall Sit', 'वाल सिट',
    ARRAY['quadriceps', 'glutes'], 'none', 'strength', 'beginner', 5.0,
    'Lean against a wall and slide down until your thighs are parallel to the ground. Hold the position for time.'),

('Tricep Dips (Bodyweight)', 'ट्राइसेप डिप्स (बडीवेट)',
    ARRAY['triceps', 'chest', 'shoulders'], 'none', 'strength', 'intermediate', 6.5,
    'Place hands on a bench or chair behind you. Lower your body by bending elbows to 90 degrees, then push back up.'),

('Bear Crawl', 'बियर क्रल',
    ARRAY['full_body', 'core', 'shoulders'], 'none', 'hiit', 'intermediate', 9.0,
    'Start on all fours with knees hovering above the ground. Crawl forward by moving opposite hand and foot together.'),

('Inchworm', 'इन्चवर्म',
    ARRAY['core', 'hamstrings', 'shoulders'], 'none', 'strength', 'beginner', 5.5,
    'Stand tall, bend forward touching the floor, walk hands out to plank, then walk feet forward to hands. Stand and repeat.'),

('Tuck Jumps', 'टक जम्प्स',
    ARRAY['quadriceps', 'glutes', 'calves'], 'none', 'hiit', 'advanced', 11.0,
    'Jump as high as possible, pulling knees up toward your chest at the peak. Land softly and immediately repeat.'),

('Box Jumps', 'बक्स जम्प्स',
    ARRAY['quadriceps', 'glutes', 'calves'], 'none', 'hiit', 'intermediate', 10.5,
    'Stand in front of a sturdy elevated surface. Jump onto it with both feet, stand fully, step down, and repeat.'),

('Skater Jumps', 'स्केटर जम्प्स',
    ARRAY['glutes', 'quadriceps', 'calves'], 'none', 'hiit', 'intermediate', 9.5,
    'Jump laterally from one foot to the other, landing softly. Swing arms for balance like a speed skater.'),

-- ============================================================
-- Dumbbell Exercises
-- ============================================================

('Dumbbell Rows', 'डम्बेल रोज',
    ARRAY['back', 'biceps'], 'dumbbells', 'strength', 'beginner', 6.0,
    'Hinge forward at the hips, holding a dumbbell in one hand. Pull the dumbbell toward your hip, squeezing your back. Lower and repeat.'),

('Goblet Squats', 'गब्लेट स्क्वाट्स',
    ARRAY['quadriceps', 'glutes', 'core'], 'dumbbells', 'strength', 'beginner', 7.0,
    'Hold a dumbbell vertically at your chest with both hands. Squat down keeping the weight close to your body, then stand.'),

('Dumbbell Lunges', 'डम्बेल लन्जेस',
    ARRAY['quadriceps', 'glutes', 'hamstrings'], 'dumbbells', 'strength', 'intermediate', 7.0,
    'Hold dumbbells at your sides. Step forward into a lunge, lowering until both knees are at 90 degrees. Push back and alternate.'),

('Dumbbell Bench Press', 'डम्बेल बेन्च प्रेस',
    ARRAY['chest', 'triceps', 'shoulders'], 'dumbbells', 'strength', 'intermediate', 6.5,
    'Lie on a bench holding dumbbells above your chest. Lower them to chest level, then press back up.'),

('Lateral Raises', 'ल्याटरल रेजेस',
    ARRAY['shoulders'], 'dumbbells', 'strength', 'beginner', 5.0,
    'Stand holding dumbbells at your sides. Raise arms out to the sides until parallel with the floor. Lower slowly.'),

('Bicep Curls', 'बाइसेप कर्ल्स',
    ARRAY['biceps'], 'dumbbells', 'strength', 'beginner', 4.5,
    'Stand holding dumbbells with palms facing forward. Curl the weights toward your shoulders, then lower slowly.'),

('Tricep Extensions', 'ट्राइसेप एक्सटेन्सन्स',
    ARRAY['triceps'], 'dumbbells', 'strength', 'beginner', 4.5,
    'Hold a dumbbell overhead with both hands. Lower it behind your head by bending elbows, then extend back up.'),

('Dumbbell Deadlift', 'डम्बेल डेडलिफ्ट',
    ARRAY['hamstrings', 'glutes', 'lower_back'], 'dumbbells', 'strength', 'intermediate', 7.0,
    'Stand holding dumbbells in front of thighs. Hinge at hips, lowering weights along your legs. Drive hips forward to return.'),

('Dumbbell Shoulder Press', 'डम्बेल शोल्डर प्रेस',
    ARRAY['shoulders', 'triceps'], 'dumbbells', 'strength', 'intermediate', 5.5,
    'Sit or stand holding dumbbells at shoulder height. Press weights overhead until arms are fully extended. Lower and repeat.'),

-- ============================================================
-- Barbell Exercises
-- ============================================================

('Barbell Bench Press', 'बार्बेल बेन्च प्रेस',
    ARRAY['chest', 'triceps', 'shoulders'], 'barbell', 'strength', 'intermediate', 7.0,
    'Lie on a bench, grip the barbell slightly wider than shoulder-width. Lower to chest, then press up to full extension.'),

('Barbell Squat', 'बार्बेल स्क्वाट',
    ARRAY['quadriceps', 'glutes', 'hamstrings', 'core'], 'barbell', 'strength', 'intermediate', 8.0,
    'Place barbell on upper back. Squat down until thighs are parallel to the floor, then drive through heels to stand.'),

('Barbell Deadlift', 'बार्बेल डेडलिफ्ट',
    ARRAY['hamstrings', 'glutes', 'lower_back', 'core'], 'barbell', 'strength', 'intermediate', 8.5,
    'Stand with feet hip-width apart. Hinge at hips and grip the bar. Drive through legs and hips to lift the bar to standing.'),

('Overhead Press', 'ओभरहेड प्रेस',
    ARRAY['shoulders', 'triceps', 'core'], 'barbell', 'strength', 'intermediate', 6.0,
    'Stand holding barbell at shoulder height. Press overhead until arms are fully extended. Lower back to shoulders.'),

('Barbell Row', 'बार्बेल रो',
    ARRAY['back', 'biceps', 'core'], 'barbell', 'strength', 'intermediate', 6.5,
    'Hinge forward holding barbell with overhand grip. Pull bar to lower chest, squeezing back muscles. Lower and repeat.'),

('Romanian Deadlift', 'रोमानियन डेडलिफ्ट',
    ARRAY['hamstrings', 'glutes', 'lower_back'], 'barbell', 'strength', 'intermediate', 7.0,
    'Hold barbell in front of thighs. Hinge at hips with slight knee bend, lowering bar along legs. Return to standing.'),

-- ============================================================
-- Flexibility / Stretching
-- ============================================================

('Cat-Cow Stretch', 'क्याट-काउ स्ट्रेच',
    ARRAY['spine', 'core'], 'none', 'flexibility', 'beginner', 2.5,
    'On all fours, alternate between arching your back (cow) and rounding it (cat). Move slowly with your breath.'),

('Downward Dog', 'डाउनवार्ड डग',
    ARRAY['hamstrings', 'calves', 'shoulders'], 'none', 'flexibility', 'beginner', 3.0,
    'From all fours, lift hips up and back forming an inverted V. Press heels toward the floor and straighten arms.'),

('Child''s Pose', 'चाइल्ड्स पोज',
    ARRAY['lower_back', 'hips', 'shoulders'], 'none', 'flexibility', 'beginner', 2.0,
    'Kneel on the floor, sit back on heels, then fold forward extending arms in front. Rest forehead on the ground.'),

('Pigeon Stretch', 'पिजन स्ट्रेच',
    ARRAY['hips', 'glutes'], 'none', 'flexibility', 'intermediate', 2.5,
    'From a plank, bring one knee forward behind the same-side wrist. Extend the other leg back and lower hips toward the floor.'),

('Forward Fold', 'फर्वार्ड फोल्ड',
    ARRAY['hamstrings', 'lower_back'], 'none', 'flexibility', 'beginner', 2.0,
    'Stand with feet hip-width apart. Hinge at hips and fold forward, reaching toward toes. Let head hang heavy.'),

('Cobra Stretch', 'कोब्रा स्ट्रेच',
    ARRAY['core', 'spine', 'chest'], 'none', 'flexibility', 'beginner', 2.5,
    'Lie face down with palms under shoulders. Press up lifting chest off the floor while keeping hips down. Hold and breathe.')

ON CONFLICT DO NOTHING;
