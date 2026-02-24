-- Seed additional food items commonly searched by users.
-- Covers South Indian, International fast food, more Nepali, and fitness foods.
-- All nutritional values are per 100g.
-- Uses ON CONFLICT DO NOTHING to avoid duplicates with existing seed data.

INSERT INTO food_items (organization_id, name, name_ne, category, calories_per_100g, protein_per_100g, carbs_per_100g, fat_per_100g, fiber_per_100g, serving_size_g, serving_label, serving_label_ne, is_verified) VALUES

-- ============================================================
-- South Indian Foods
-- ============================================================

(NULL, 'Idli (Steamed Rice Cake)', 'इड्ली', 'grain',
    130, 2, 26, 0.5, 0.8,
    80, '2 pieces', '२ वटा', true),

(NULL, 'Dosa (Plain)', 'डोसा (सादा)', 'grain',
    120, 3, 18, 4, 0.5,
    100, '1 piece', '१ वटा', true),

(NULL, 'Masala Dosa', 'मसाला डोसा', 'grain',
    165, 4, 22, 7, 1.2,
    150, '1 piece', '१ वटा', true),

(NULL, 'Sambar', 'साम्बार', 'legume',
    65, 3, 10, 1.5, 2.5,
    200, '1 bowl', '१ कचौरा', true),

(NULL, 'Vada (Medu)', 'मेदु वडा', 'snack',
    290, 6, 25, 18, 2,
    50, '2 pieces', '२ वटा', true),

(NULL, 'Uttapam', 'उत्तपम', 'grain',
    150, 4, 22, 5, 1,
    120, '1 piece', '१ वटा', true),

(NULL, 'Upma', 'उप्मा', 'grain',
    120, 3, 18, 4, 1.5,
    200, '1 bowl', '१ कचौरा', true),

-- ============================================================
-- International Fast Food / Popular
-- ============================================================

(NULL, 'Pizza (Cheese)', 'पिज्जा (चिज)', 'snack',
    266, 11, 33, 10, 2.3,
    107, '1 slice', '१ स्लाइस', true),

(NULL, 'Pizza (Pepperoni)', 'पिज्जा (पेपरोनी)', 'snack',
    300, 13, 33, 13, 2,
    107, '1 slice', '१ स्लाइस', true),

(NULL, 'Burger (Beef)', 'बर्गर (बीफ)', 'protein',
    250, 15, 24, 10, 1.5,
    200, '1 burger', '१ बर्गर', true),

(NULL, 'French Fries', 'फ्रेन्च फ्राइज', 'snack',
    312, 3, 41, 15, 3.8,
    120, '1 serving', '१ सर्भिङ', true),

(NULL, 'Fried Chicken', 'फ्राइड चिकेन', 'protein',
    260, 20, 10, 16, 0.5,
    100, '1 piece', '१ वटा', true),

(NULL, 'Sandwich (Club)', 'स्यान्डविच (क्लब)', 'snack',
    220, 12, 25, 8, 1.5,
    150, '1 sandwich', '१ स्यान्डविच', true),

-- ============================================================
-- More Nepali Traditional Foods
-- Short names added as separate entries for better search matching.
-- Existing entries have longer descriptive names like "Thukpa (Noodle Soup)".
-- ============================================================

(NULL, 'Thukpa', 'थुक्पा', 'nepali_traditional',
    80, 4, 12, 2, 1,
    350, '1 bowl', '१ कचौरा', true),

(NULL, 'Chatamari', 'चतामरी', 'nepali_traditional',
    180, 6, 25, 6, 1,
    120, '1 piece', '१ वटा', true),

(NULL, 'Yomari', 'योमरी', 'nepali_traditional',
    200, 4, 35, 5, 0.5,
    80, '2 pieces', '२ वटा', true),

(NULL, 'Kwati (Bean Soup)', 'क्वाटी (मिश्रित दालको झोल)', 'nepali_traditional',
    85, 5, 14, 1, 4,
    250, '1 bowl', '१ कचौरा', true),

(NULL, 'Gundruk', 'गुन्द्रुक', 'nepali_traditional',
    20, 2, 3, 0.3, 2,
    30, '1 serving', '१ सर्भिङ', true),

(NULL, 'Dhido', 'ढिडो', 'nepali_traditional',
    140, 2, 32, 0.5, 2,
    250, '1 serving', '१ सर्भिङ', true),

(NULL, 'Sel Roti (Ring Bread)', 'सेल रोटी (रिंग ब्रेड)', 'nepali_traditional',
    350, 4, 55, 12, 1,
    60, '1 piece', '१ वटा', true),

(NULL, 'Bara (Fried Lentil Patty)', 'बरा (तलेको)', 'nepali_traditional',
    250, 8, 28, 12, 2.5,
    80, '1 piece', '१ वटा', true),

-- ============================================================
-- Healthy / Fitness Foods
-- ============================================================

(NULL, 'Protein Shake (Whey)', 'प्रोटिन शेक (व्हे)', 'supplement',
    120, 24, 3, 1.5, 0,
    300, '1 shake', '१ शेक', true),

(NULL, 'Greek Yogurt', 'ग्रीक दही', 'dairy',
    59, 10, 4, 0.7, 0,
    150, '1 cup', '१ कप', true),

(NULL, 'Oatmeal (Cooked)', 'ओटमिल (पकाएको)', 'grain',
    68, 2.4, 12, 1.4, 1.7,
    240, '1 bowl', '१ कचौरा', true)

ON CONFLICT DO NOTHING;
