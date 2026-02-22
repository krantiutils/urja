-- Expand food database with ~100 additional items: Nepali staples + international fitness foods.
-- Uses ON CONFLICT DO NOTHING to avoid duplicates with existing seed data.

INSERT INTO food_items (organization_id, name, name_ne, category, calories_per_100g, protein_per_100g, carbs_per_100g, fat_per_100g, fiber_per_100g, serving_size_g, serving_label, serving_label_ne, is_verified) VALUES

-- ============================================================
-- Nepali Traditional Foods
-- ============================================================

(NULL, 'Kwati (Mixed Bean Soup)', 'क्वाटी', 'nepali_traditional',
    120, 8, 18, 1.5, 5,
    250, '1 bowl', '१ कचौरा', true),

(NULL, 'Yomari (Rice Flour Dumpling)', 'योमरी', 'nepali_traditional',
    280, 4, 48, 8, 1,
    80, '2 pieces', '२ वटा', true),

(NULL, 'Sisno (Nettle Soup)', 'सिस्नो', 'nepali_traditional',
    42, 2.7, 7, 0.1, 6.9,
    200, '1 bowl', '१ कचौरा', true),

(NULL, 'Saag (Mixed Greens)', 'साग', 'nepali_traditional',
    35, 3, 4, 0.5, 3,
    150, '1 serving', '१ सर्भिङ', true),

(NULL, 'Bhuja (Fried Rice Snack)', 'भुजा', 'nepali_traditional',
    400, 7, 55, 18, 2,
    50, '1 cup', '१ कप', true),

(NULL, 'Papad (Lentil Crackers)', 'पापड', 'nepali_traditional',
    350, 20, 45, 8, 5,
    15, '1 piece', '१ वटा', true),

(NULL, 'Lapsi (Nepali Hog Plum)', 'लप्सी', 'nepali_traditional',
    55, 0.5, 13, 0.2, 2,
    100, '3-4 fruits', '३-४ वटा', true),

(NULL, 'Titaura (Dried Fruit Candy)', 'तितौरा', 'nepali_traditional',
    300, 1, 72, 0.5, 3,
    20, '2 pieces', '२ वटा', true),

(NULL, 'Aalu Tama (Potato Bamboo Shoot Curry)', 'आलु तामा', 'nepali_traditional',
    80, 2, 12, 3, 2,
    200, '1 bowl', '१ कचौरा', true),

(NULL, 'Wo (Newari Lentil Cake)', 'वो', 'nepali_traditional',
    210, 9, 20, 10, 3,
    80, '1 piece', '१ वटा', true),

(NULL, 'Tongba (Millet Beer)', 'तोंग्बा', 'nepali_traditional',
    45, 0.5, 5, 0, 0,
    300, '1 serving', '१ सर्भिङ', true),

(NULL, 'Chhurpi (Hard Cheese)', 'छुर्पी', 'nepali_traditional',
    350, 50, 5, 15, 0,
    20, '1 piece', '१ टुक्रा', true),

(NULL, 'Kinema (Fermented Soybeans)', 'किनेमा', 'nepali_traditional',
    200, 18, 10, 10, 6,
    50, '1 serving', '१ सर्भिङ', true),

(NULL, 'Sekuwa (Grilled Meat)', 'सेकुवा', 'nepali_traditional',
    190, 28, 2, 8, 0,
    150, '1 serving', '१ सर्भिङ', true),

(NULL, 'Choila (Spiced Grilled Meat)', 'चोइला', 'nepali_traditional',
    180, 25, 3, 8, 0.5,
    100, '1 serving', '१ सर्भिङ', true),

(NULL, 'Bhakka (Steamed Rice Cake)', 'भक्का', 'nepali_traditional',
    180, 3, 38, 1, 0.5,
    100, '2 pieces', '२ वटा', true),

(NULL, 'Maas ko Dal (Black Lentil Soup)', 'मासको दाल', 'nepali_traditional',
    105, 7, 18, 0.5, 4,
    250, '1 bowl', '१ कचौरा', true),

(NULL, 'Masyaura (Dried Lentil Nuggets)', 'मस्यौरा', 'nepali_traditional',
    320, 22, 50, 3, 5,
    30, '5 pieces', '५ वटा', true),

(NULL, 'Aaloo Chop (Potato Cutlet)', 'आलुचप', 'nepali_traditional',
    200, 3, 22, 12, 2,
    60, '1 piece', '१ वटा', true),

(NULL, 'Sikarni (Sweet Yogurt)', 'सिकार्नी', 'nepali_traditional',
    150, 5, 20, 6, 0,
    100, '1 cup', '१ कप', true),

-- ============================================================
-- Grains & Carbs (International)
-- ============================================================

(NULL, 'Quinoa (Cooked)', 'क्विनोआ (पकाएको)', 'grains',
    120, 4.4, 21, 1.9, 2.8,
    185, '1 cup cooked', '१ कप पकाएको', true),

(NULL, 'Couscous (Cooked)', 'कुसकुस (पकाएको)', 'grains',
    112, 3.8, 23, 0.2, 1.4,
    160, '1 cup cooked', '१ कप पकाएको', true),

(NULL, 'Basmati Rice (Cooked)', 'बासमती चामल (पकाएको)', 'grains',
    130, 2.7, 28, 0.3, 0.4,
    200, '1 cup cooked', '१ कप पकाएको', true),

(NULL, 'Granola', 'ग्रानोला', 'grains',
    471, 10, 64, 20, 7,
    45, '1/2 cup', '१/२ कप', true),

(NULL, 'Rice Cakes', 'राइस केक', 'grains',
    387, 7, 82, 3, 1.3,
    9, '1 cake', '१ केक', true),

(NULL, 'Buckwheat (Cooked)', 'फापर (पकाएको)', 'grains',
    92, 3.4, 20, 0.6, 2.7,
    170, '1 cup cooked', '१ कप पकाएको', true),

(NULL, 'Millet (Cooked)', 'कोदो (पकाएको)', 'grains',
    119, 3.5, 23, 1.0, 1.3,
    175, '1 cup cooked', '१ कप पकाएको', true),

-- ============================================================
-- Protein Sources (International)
-- ============================================================

(NULL, 'Chicken Breast (Plain, Boiled)', 'चिकेन ब्रेस्ट (उसिनेको)', 'protein',
    165, 31, 0, 3.6, 0,
    120, '1 breast', '१ ब्रेस्ट', true),

(NULL, 'Chicken Thigh (Skinless)', 'चिकेन थाई (छालाविना)', 'protein',
    177, 24, 0, 8.5, 0,
    100, '1 thigh', '१ थाई', true),

(NULL, 'Turkey Breast', 'टर्की ब्रेस्ट', 'protein',
    135, 30, 0, 1, 0,
    120, '1 serving', '१ सर्भिङ', true),

(NULL, 'Salmon Fillet (Baked)', 'स्याल्मन फिलेट (बेक्ड)', 'protein',
    208, 20, 0, 13, 0,
    125, '1 fillet', '१ फिलेट', true),

(NULL, 'Tuna (Canned in Water)', 'ट्युना (पानीमा क्यान्ड)', 'protein',
    116, 26, 0, 0.8, 0,
    100, '1 can', '१ क्यान', true),

(NULL, 'Shrimp (Cooked)', 'झिँगा (पकाएको)', 'protein',
    99, 24, 0.2, 0.3, 0,
    85, '1 serving', '१ सर्भिङ', true),

(NULL, 'Tempeh', 'टेम्पे', 'protein',
    192, 20, 7.6, 11, 0,
    85, '1/2 cup', '१/२ कप', true),

(NULL, 'Lamb (Lean)', 'खसीको मासु (मसिनो)', 'protein',
    250, 26, 0, 16, 0,
    100, '1 serving', '१ सर्भिङ', true),

(NULL, 'Buff/Buffalo Meat (Lean)', 'बफ/रांगोको मासु (मसिनो)', 'protein',
    131, 27, 0, 2, 0,
    100, '1 serving', '१ सर्भिङ', true),

(NULL, 'Egg Whites', 'अण्डाको सेतो भाग', 'protein',
    52, 11, 0.7, 0.2, 0,
    100, '3 egg whites', '३ अण्डा सेतो', true),

(NULL, 'Pork Tenderloin', 'बंगुरको फिलेट', 'protein',
    143, 26, 0, 3.5, 0,
    100, '1 serving', '१ सर्भिङ', true),

-- ============================================================
-- Dairy (International)
-- ============================================================

(NULL, 'Greek Yogurt (Plain)', 'ग्रीक दही (सादा)', 'dairy',
    59, 10, 3.6, 0.7, 0,
    150, '1 cup', '१ कप', true),

(NULL, 'Milk (Skim)', 'दूध (स्किम)', 'dairy',
    34, 3.4, 5, 0.1, 0,
    250, '1 glass', '१ गिलास', true),

(NULL, 'Cheese (Cheddar)', 'चिज (चेडर)', 'dairy',
    403, 25, 1.3, 33, 0,
    28, '1 slice', '१ स्लाइस', true),

(NULL, 'Mozzarella Cheese', 'मोजारेला चिज', 'dairy',
    280, 22, 2.2, 17, 0,
    28, '1 slice', '१ स्लाइस', true),

(NULL, 'Cream Cheese', 'क्रीम चिज', 'dairy',
    342, 6, 4, 34, 0,
    30, '2 tablespoons', '२ चम्चा', true),

(NULL, 'Butter', 'मक्खन', 'dairy',
    717, 0.9, 0.1, 81, 0,
    14, '1 tablespoon', '१ चम्चा', true),

(NULL, 'Ghee', 'घिउ', 'dairy',
    900, 0, 0, 100, 0,
    14, '1 tablespoon', '१ चम्चा', true),

-- ============================================================
-- Vegetables
-- ============================================================

(NULL, 'Kale', 'केल', 'vegetables',
    49, 4.3, 9, 0.9, 3.6,
    67, '1 cup chopped', '१ कप काटेको', true),

(NULL, 'Cauliflower', 'काउली', 'vegetables',
    25, 1.9, 5, 0.3, 2,
    100, '1 cup', '१ कप', true),

(NULL, 'Green Beans', 'सिमी', 'vegetables',
    31, 1.8, 7, 0.2, 2.7,
    100, '1 cup', '१ कप', true),

(NULL, 'Bell Pepper', 'भेडे खुर्सानी', 'vegetables',
    31, 1, 6, 0.3, 2.1,
    120, '1 medium', '१ मध्यम', true),

(NULL, 'Zucchini', 'जुकिनी', 'vegetables',
    17, 1.2, 3.1, 0.3, 1,
    200, '1 medium', '१ मध्यम', true),

(NULL, 'Carrot', 'गाजर', 'vegetables',
    41, 0.9, 10, 0.2, 2.8,
    60, '1 medium', '१ मध्यम', true),

(NULL, 'Cabbage', 'बन्दा', 'vegetables',
    25, 1.3, 6, 0.1, 2.5,
    90, '1 cup shredded', '१ कप काटेको', true),

(NULL, 'Asparagus', 'एस्पारागस', 'vegetables',
    20, 2.2, 3.9, 0.1, 2.1,
    100, '5 spears', '५ वटा', true),

(NULL, 'Lettuce (Romaine)', 'लेटुस', 'vegetables',
    17, 1.2, 3.3, 0.3, 2.1,
    50, '1 cup shredded', '१ कप काटेको', true),

(NULL, 'Radish (Mula)', 'मुला', 'vegetables',
    16, 0.7, 3.4, 0.1, 1.6,
    100, '1 medium', '१ मध्यम', true),

(NULL, 'Pumpkin', 'फर्सी', 'vegetables',
    26, 1, 6.5, 0.1, 0.5,
    120, '1 cup cubed', '१ कप काटेको', true),

(NULL, 'Bitter Gourd (Karela)', 'करेला', 'vegetables',
    17, 1, 3.7, 0.2, 2.8,
    100, '1 medium', '१ मध्यम', true),

-- ============================================================
-- Fruits
-- ============================================================

(NULL, 'Papaya', 'मेवा', 'fruits',
    43, 0.5, 11, 0.3, 1.7,
    150, '1 cup cubed', '१ कप काटेको', true),

(NULL, 'Guava', 'अमरुद', 'fruits',
    68, 2.6, 14, 1, 5.4,
    55, '1 medium', '१ मध्यम', true),

(NULL, 'Pomegranate', 'अनार', 'fruits',
    83, 1.7, 19, 1.2, 4,
    175, '1 medium', '१ मध्यम', true),

(NULL, 'Grapes', 'अंगुर', 'fruits',
    69, 0.7, 18, 0.2, 0.9,
    150, '1 cup', '१ कप', true),

(NULL, 'Pineapple', 'भुइँकटहर', 'fruits',
    50, 0.5, 13, 0.1, 1.4,
    165, '1 cup cubed', '१ कप काटेको', true),

(NULL, 'Strawberries', 'स्ट्रबेरी', 'fruits',
    32, 0.7, 7.7, 0.3, 2,
    150, '1 cup', '१ कप', true),

(NULL, 'Blueberries', 'ब्लुबेरी', 'fruits',
    57, 0.7, 14, 0.3, 2.4,
    150, '1 cup', '१ कप', true),

(NULL, 'Lychee', 'लिची', 'fruits',
    66, 0.8, 17, 0.4, 1.3,
    100, '5-6 fruits', '५-६ वटा', true),

(NULL, 'Jackfruit', 'कटहर', 'fruits',
    95, 1.7, 23, 0.6, 1.5,
    150, '1 cup cubed', '१ कप काटेको', true),

-- ============================================================
-- Nuts & Seeds
-- ============================================================

(NULL, 'Walnuts', 'ओखर', 'nuts_seeds',
    654, 15, 14, 65, 6.7,
    28, '1 handful', '१ मुठ्ठी', true),

(NULL, 'Cashews', 'काजु', 'nuts_seeds',
    553, 18, 30, 44, 3.3,
    28, '1 handful', '१ मुठ्ठी', true),

(NULL, 'Pistachios', 'पिस्ता', 'nuts_seeds',
    560, 20, 28, 45, 10,
    28, '1 handful', '१ मुठ्ठी', true),

(NULL, 'Chia Seeds', 'चिया बीउ', 'nuts_seeds',
    486, 17, 42, 31, 34,
    15, '1 tablespoon', '१ चम्चा', true),

(NULL, 'Flax Seeds', 'अलसी', 'nuts_seeds',
    534, 18, 29, 42, 27,
    10, '1 tablespoon', '१ चम्चा', true),

(NULL, 'Sunflower Seeds', 'सूर्यमुखी बीउ', 'nuts_seeds',
    584, 21, 20, 51, 8.6,
    28, '1 handful', '१ मुठ्ठी', true),

(NULL, 'Pumpkin Seeds', 'फर्सीको बीउ', 'nuts_seeds',
    559, 30, 11, 49, 6,
    28, '1 handful', '१ मुठ्ठी', true),

(NULL, 'Peanuts (Roasted)', 'बदाम (भुटेको)', 'nuts_seeds',
    567, 26, 16, 49, 8.5,
    28, '1 handful', '१ मुठ्ठी', true),

-- ============================================================
-- Oils & Fats
-- ============================================================

(NULL, 'Coconut Oil', 'नरिवल तेल', 'oils',
    862, 0, 0, 100, 0,
    14, '1 tablespoon', '१ चम्चा', true),

(NULL, 'Mustard Oil', 'तोरीको तेल', 'oils',
    884, 0, 0, 100, 0,
    14, '1 tablespoon', '१ चम्चा', true),

(NULL, 'Sesame Oil', 'तिलको तेल', 'oils',
    884, 0, 0, 100, 0,
    14, '1 tablespoon', '१ चम्चा', true),

-- ============================================================
-- Beverages
-- ============================================================

(NULL, 'Black Coffee', 'ब्ल्याक कफी', 'beverages',
    2, 0.3, 0, 0, 0,
    240, '1 cup', '१ कप', true),

(NULL, 'Green Tea', 'हरियो चिया', 'beverages',
    1, 0.2, 0.3, 0, 0,
    240, '1 cup', '१ कप', true),

(NULL, 'Coconut Water', 'नरिवल पानी', 'beverages',
    19, 0.7, 3.7, 0.2, 1.1,
    240, '1 cup', '१ कप', true),

(NULL, 'Orange Juice (Fresh)', 'सुन्तला जुस (ताजा)', 'beverages',
    45, 0.7, 10, 0.2, 0.2,
    250, '1 glass', '१ गिलास', true),

(NULL, 'Protein Shake (Mixed)', 'प्रोटिन शेक (मिक्स्ड)', 'beverages',
    120, 24, 5, 1.5, 0,
    350, '1 shake', '१ शेक', true),

-- ============================================================
-- Legumes & Beans
-- ============================================================

(NULL, 'Chana (Chickpeas, Dry Roasted)', 'चना (भुटेको)', 'protein',
    164, 9, 27, 2.6, 7.6,
    50, '1/2 cup', '१/२ कप', true),

(NULL, 'Rajma (Kidney Beans, Cooked)', 'राजमा (पकाएको)', 'protein',
    127, 8.7, 22, 0.5, 6.4,
    175, '1 cup cooked', '१ कप पकाएको', true),

(NULL, 'Soybean (Bhatmas, Cooked)', 'भटमास (पकाएको)', 'protein',
    173, 17, 10, 9, 6,
    175, '1 cup cooked', '१ कप पकाएको', true),

(NULL, 'Moong Dal (Cooked)', 'मुँगको दाल (पकाएको)', 'protein',
    106, 7, 19, 0.4, 7.6,
    200, '1 bowl', '१ कचौरा', true),

(NULL, 'Hummus', 'हम्मस', 'protein',
    166, 8, 14, 10, 6,
    30, '2 tablespoons', '२ चम्चा', true),

(NULL, 'Lentil Soup (Mixed)', 'मिश्रित दालको झोल', 'protein',
    80, 5, 13, 0.4, 4,
    250, '1 bowl', '१ कचौरा', true),

-- ============================================================
-- Snacks
-- ============================================================

(NULL, 'Protein Bar', 'प्रोटिन बार', 'snacks',
    350, 20, 40, 12, 5,
    60, '1 bar', '१ बार', true),

(NULL, 'Trail Mix', 'ट्रेल मिक्स', 'snacks',
    462, 14, 44, 29, 5,
    40, '1/4 cup', '१/४ कप', true),

(NULL, 'Roasted Chickpeas', 'भुटेको चना', 'snacks',
    390, 20, 56, 8, 12,
    40, '1/4 cup', '१/४ कप', true),

(NULL, 'Popcorn (Air-popped)', 'पप्कर्न (तेलबिना)', 'snacks',
    375, 12, 74, 4.3, 14.5,
    28, '3 cups popped', '३ कप पप्कर्न', true),

-- ============================================================
-- Supplements
-- ============================================================

(NULL, 'Casein Protein Powder', 'क्यासिन प्रोटिन पाउडर', 'supplements',
    360, 75, 6, 4, 0,
    33, '1 scoop', '१ स्कूप', true),

(NULL, 'Plant Protein Powder', 'प्लान्ट प्रोटिन पाउडर', 'supplements',
    350, 70, 10, 5, 3,
    35, '1 scoop', '१ स्कूप', true),

(NULL, 'BCAA Powder', 'बीसीएए पाउडर', 'supplements',
    0, 0, 0, 0, 0,
    7, '1 scoop', '१ स्कूप', true),

(NULL, 'Creatine Monohydrate', 'क्रिएटिन', 'supplements',
    0, 0, 0, 0, 0,
    5, '1 scoop', '१ स्कूप', true),

(NULL, 'Mass Gainer', 'मास गेनर', 'supplements',
    400, 25, 65, 5, 2,
    100, '1 serving', '१ सर्भिङ', true)

ON CONFLICT DO NOTHING;
