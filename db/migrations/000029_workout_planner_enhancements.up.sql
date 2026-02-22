-- Workout planner enhancements: add goal, days_per_week, and assignment_method columns.

ALTER TABLE workout_templates
    ADD COLUMN goal VARCHAR(50)
        CHECK (goal IN ('lose_weight', 'build_muscle', 'stay_fit'));

ALTER TABLE workout_templates
    ADD COLUMN days_per_week VARCHAR(10)
        CHECK (days_per_week IN ('2-3', '4-5', '6-7'));

ALTER TABLE member_workout_plans
    ADD COLUMN assignment_method VARCHAR(20) NOT NULL DEFAULT 'staff'
        CHECK (assignment_method IN ('staff', 'self', 'questionnaire'));
