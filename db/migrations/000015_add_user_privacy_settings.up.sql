ALTER TABLE users ADD COLUMN privacy_settings JSONB NOT NULL DEFAULT '{
    "show_email": false,
    "show_phone": false,
    "show_profile": true,
    "show_attendance": true,
    "show_on_leaderboard": true
}';
