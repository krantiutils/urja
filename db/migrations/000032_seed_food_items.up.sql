-- Seed global food items with nutritional data per 100g.

INSERT INTO food_items (organization_id, name, name_ne, category, calories_per_100g, protein_per_100g, carbs_per_100g, fat_per_100g, fiber_per_100g, serving_size_g, serving_label, serving_label_ne, is_verified) VALUES

-- ============================================================
-- Nepali Foods
-- ============================================================

-- Grains & Staples
(NULL, 'Bhat (Steamed Rice)', 'भात (उसिनेको चामल)', 'grain',
    130, 2.7, 28, 0.3, 0.4,
    200, '1 plate', '१ प्लेट', true),

(NULL, 'Chiura (Beaten Rice)', 'चिउरा', 'grain',
    345, 6.6, 77, 1.2, 1.0,
    50, '1 cup', '१ कप', true),

(NULL, 'Dhido (Buckwheat Porridge)', 'ढिडो', 'grain',
    340, 12, 72, 2, 10,
    250, '1 serving', '१ सर्भिङ', true),

(NULL, 'Roti (Whole Wheat)', 'रोटी (गहुँको)', 'grain',
    297, 9, 50, 7, 6,
    40, '1 piece', '१ वटा', true),

(NULL, 'Sel Roti', 'सेल रोटी', 'snack',
    350, 5, 55, 12, 1.5,
    60, '1 piece', '१ वटा', true),

-- Protein / Main Dishes
(NULL, 'Dal (Lentil Soup)', 'दाल', 'legume',
    116, 9, 20, 0.4, 4,
    250, '1 bowl', '१ कचौरा', true),

(NULL, 'Tarkari (Mixed Vegetable Curry)', 'तरकारी (मिश्रित)', 'vegetable',
    65, 2, 8, 3, 2.5,
    200, '1 bowl', '१ कचौरा', true),

(NULL, 'Momo (Steamed, Chicken)', 'मोमो (स्टिम, चिकेन)', 'protein',
    180, 12, 18, 7, 1,
    160, '8 pieces', '८ वटा', true),

(NULL, 'Chowmein', 'चाउमिन', 'grain',
    220, 7, 30, 8, 2,
    250, '1 plate', '१ प्लेट', true),

(NULL, 'Sukuti (Dried Meat)', 'सुकुटी', 'protein',
    250, 40, 0, 10, 0,
    50, '1 serving', '१ सर्भिङ', true),

(NULL, 'Thukpa (Noodle Soup)', 'थुक्पा', 'soup',
    120, 6, 15, 4, 1.5,
    350, '1 bowl', '१ कचौरा', true),

(NULL, 'Chatamari (Rice Crepe)', 'चतामरी', 'snack',
    200, 8, 25, 8, 1.5,
    120, '1 piece', '१ वटा', true),

(NULL, 'Bara (Lentil Patty)', 'बरा', 'legume',
    220, 10, 20, 11, 3,
    80, '1 piece', '१ वटा', true),

-- Snacks & Sweets
(NULL, 'Samosa', 'समोसा', 'snack',
    262, 4, 24, 17, 2,
    80, '1 piece', '१ वटा', true),

(NULL, 'Puri', 'पुरी', 'snack',
    297, 5, 36, 15, 1.5,
    30, '1 piece', '१ वटा', true),

(NULL, 'Jeri/Jalebi', 'जेरी/जलेबी', 'sweet',
    370, 2, 60, 14, 0.5,
    50, '2 pieces', '२ वटा', true),

(NULL, 'Kheer (Rice Pudding)', 'खिर', 'sweet',
    140, 3, 20, 5, 0.3,
    150, '1 bowl', '१ कचौरा', true),

-- Fermented & Pickles
(NULL, 'Gundruk (Fermented Greens)', 'गुन्द्रुक', 'vegetable',
    25, 3, 3, 0.2, 2,
    30, '1 serving', '१ सर्भिङ', true),

(NULL, 'Achar (Pickle)', 'अचार', 'condiment',
    40, 1, 5, 2, 1,
    30, '1 tablespoon', '१ चम्चा', true),

(NULL, 'Sadheko (Spiced Salad)', 'सधेको', 'snack',
    90, 5, 8, 5, 1.5,
    100, '1 serving', '१ सर्भिङ', true),

-- Dairy & Beverages
(NULL, 'Dahi (Yogurt)', 'दही', 'dairy',
    60, 3.5, 5, 3.3, 0,
    100, '1 cup', '१ कप', true),

(NULL, 'Lassi (Sweet)', 'लस्सी (मीठो)', 'beverage',
    75, 2, 12, 2, 0,
    250, '1 glass', '१ गिलास', true),

(NULL, 'Chiya (Milk Tea)', 'चिया', 'beverage',
    45, 1, 6, 1.5, 0,
    150, '1 cup', '१ कप', true),

-- ============================================================
-- International Foods — Protein
-- ============================================================

(NULL, 'Chicken Breast (Grilled)', 'चिकेन ब्रेस्ट (ग्रिल्ड)', 'protein',
    165, 31, 0, 3.6, 0,
    120, '1 breast', '१ ब्रेस्ट', true),

(NULL, 'Eggs (Boiled)', 'अण्डा (उसिनेको)', 'protein',
    155, 13, 1.1, 11, 0,
    50, '1 egg', '१ अण्डा', true),

(NULL, 'Tuna (Canned)', 'ट्युना (क्यान्ड)', 'protein',
    116, 25.5, 0, 0.8, 0,
    100, '1 can', '१ क्यान', true),

(NULL, 'Salmon', 'स्याल्मन', 'protein',
    208, 20, 0, 13, 0,
    125, '1 fillet', '१ फिलेट', true),

(NULL, 'Tofu', 'टोफु', 'protein',
    76, 8, 1.9, 4.8, 0.3,
    100, '1/2 block', '१/२ ब्लक', true),

(NULL, 'Paneer', 'पनीर', 'dairy',
    265, 18, 1.2, 21, 0,
    50, '1 piece', '१ टुक्रा', true),

-- ============================================================
-- International Foods — Dairy
-- ============================================================

(NULL, 'Milk (Whole)', 'दूध (पूरा)', 'dairy',
    61, 3.2, 4.8, 3.3, 0,
    250, '1 glass', '१ गिलास', true),

(NULL, 'Yogurt (Plain)', 'दही (सादा)', 'dairy',
    59, 10, 3.6, 0.7, 0,
    150, '1 cup', '१ कप', true),

(NULL, 'Cottage Cheese', 'कटेज चिज', 'dairy',
    98, 11, 3.4, 4.3, 0,
    100, '1/2 cup', '१/२ कप', true),

-- ============================================================
-- International Foods — Grains & Carbs
-- ============================================================

(NULL, 'Oats', 'ओट्स', 'grain',
    389, 17, 66, 7, 10.6,
    40, '1/2 cup dry', '१/२ कप सुक्खा', true),

(NULL, 'Brown Rice', 'ब्राउन राइस', 'grain',
    112, 2.3, 24, 0.8, 1.8,
    200, '1 plate', '१ प्लेट', true),

(NULL, 'Bread (White)', 'पाउरोटी (सेतो)', 'grain',
    265, 9, 49, 3.2, 2.7,
    30, '1 slice', '१ स्लाइस', true),

(NULL, 'Bread (Whole Wheat)', 'पाउरोटी (गहुँको)', 'grain',
    247, 13, 41, 3.4, 7,
    30, '1 slice', '१ स्लाइस', true),

(NULL, 'Pasta (Cooked)', 'पास्ता (पकाएको)', 'grain',
    131, 5, 25, 1.1, 1.8,
    200, '1 plate', '१ प्लेट', true),

(NULL, 'Sweet Potato', 'सखरखण्ड', 'grain',
    86, 1.6, 20, 0.1, 3,
    130, '1 medium', '१ मध्यम', true),

(NULL, 'Potato', 'आलु', 'vegetable',
    77, 2, 17, 0.1, 2.2,
    150, '1 medium', '१ मध्यम', true),

(NULL, 'Corn', 'मकै', 'grain',
    86, 3.3, 19, 1.2, 2.7,
    100, '1 ear', '१ भुट्टा', true),

-- ============================================================
-- International Foods — Legumes
-- ============================================================

(NULL, 'Lentils (Cooked)', 'मसुरो (पकाएको)', 'legume',
    116, 9, 20, 0.4, 7.9,
    200, '1 cup', '१ कप', true),

(NULL, 'Chickpeas (Cooked)', 'चना (पकाएको)', 'legume',
    164, 9, 27, 2.6, 7.6,
    150, '1 cup', '१ कप', true),

-- ============================================================
-- International Foods — Fruits
-- ============================================================

(NULL, 'Banana', 'केरा', 'fruit',
    89, 1.1, 23, 0.3, 2.6,
    120, '1 medium', '१ मध्यम', true),

(NULL, 'Apple', 'स्याउ', 'fruit',
    52, 0.3, 14, 0.2, 2.4,
    180, '1 medium', '१ मध्यम', true),

(NULL, 'Orange', 'सुन्तला', 'fruit',
    47, 0.9, 12, 0.1, 2.4,
    130, '1 medium', '१ मध्यम', true),

(NULL, 'Watermelon', 'तरबूजा', 'fruit',
    30, 0.6, 8, 0.2, 0.4,
    280, '1 slice', '१ स्लाइस', true),

(NULL, 'Mango', 'आँप', 'fruit',
    60, 0.8, 15, 0.4, 1.6,
    200, '1 medium', '१ मध्यम', true),

-- ============================================================
-- International Foods — Vegetables
-- ============================================================

(NULL, 'Broccoli', 'ब्रोकाउली', 'vegetable',
    34, 2.8, 7, 0.4, 2.6,
    90, '1 cup chopped', '१ कप काटेको', true),

(NULL, 'Spinach', 'पालुंगो', 'vegetable',
    23, 2.9, 3.6, 0.4, 2.2,
    30, '1 cup raw', '१ कप कच्चो', true),

(NULL, 'Cucumber', 'काँक्रो', 'vegetable',
    15, 0.7, 3.6, 0.1, 0.5,
    300, '1 medium', '१ मध्यम', true),

(NULL, 'Tomato', 'गोलभेडा', 'vegetable',
    18, 0.9, 3.9, 0.2, 1.2,
    120, '1 medium', '१ मध्यम', true),

(NULL, 'Onion', 'प्याज', 'vegetable',
    40, 1.1, 9, 0.1, 1.7,
    110, '1 medium', '१ मध्यम', true),

(NULL, 'Mushroom', 'च्याउ', 'vegetable',
    22, 3.1, 3.3, 0.3, 1.0,
    70, '1 cup sliced', '१ कप काटेको', true),

-- ============================================================
-- International Foods — Nuts & Seeds
-- ============================================================

(NULL, 'Peanut Butter', 'पीनट बटर', 'nuts_seeds',
    588, 25, 20, 50, 6,
    32, '2 tablespoons', '२ चम्चा', true),

(NULL, 'Almonds', 'बदाम', 'nuts_seeds',
    579, 21, 22, 50, 12.5,
    28, '1 handful', '१ मुठ्ठी', true),

-- ============================================================
-- International Foods — Fats & Oils
-- ============================================================

(NULL, 'Avocado', 'एभोकाडो', 'fruit',
    160, 2, 9, 15, 6.7,
    150, '1 medium', '१ मध्यम', true),

(NULL, 'Olive Oil', 'जैतुन तेल', 'fat_oil',
    884, 0, 0, 100, 0,
    14, '1 tablespoon', '१ चम्चा', true),

-- ============================================================
-- International Foods — Sweets & Condiments
-- ============================================================

(NULL, 'Honey', 'मह', 'sweet',
    304, 0.3, 82, 0, 0.2,
    21, '1 tablespoon', '१ चम्चा', true),

(NULL, 'Dark Chocolate', 'डार्क चकलेट', 'sweet',
    546, 5, 60, 31, 7,
    30, '1 square', '१ टुक्रा', true),

-- ============================================================
-- Supplement
-- ============================================================

(NULL, 'Whey Protein Powder', 'व्हे प्रोटिन पाउडर', 'supplement',
    370, 80, 7, 3.5, 0,
    30, '1 scoop', '१ स्कूप', true);
