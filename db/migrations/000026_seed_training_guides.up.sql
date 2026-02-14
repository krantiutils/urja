-- Seed data: 5 sample bilingual training guides
INSERT INTO training_guides (id, title, title_ne, content, content_ne, category, is_published) VALUES
    ('e1b2c3d4-e5f6-7890-abcd-ef1234567801',
     'Beginner''s Guide to Gym Training',
     'जिम प्रशिक्षणको सुरुवात गाइड',
     'Starting your fitness journey can feel overwhelming, but with the right approach, you''ll build lasting habits. Begin with 3 days per week, focusing on compound movements like squats, deadlifts, and bench press. Start with light weights to master form before progressing. Rest 60-90 seconds between sets. Always warm up for 5-10 minutes with light cardio. Cool down with stretching to prevent injury and improve flexibility.',
     'फिटनेस यात्रा सुरु गर्नु भारी लाग्न सक्छ, तर सही दृष्टिकोणले तपाईंलाई स्थायी बानी बनाउन मद्दत गर्नेछ। हप्ताको ३ दिनबाट सुरु गर्नुहोस्, स्क्वाट, डेडलिफ्ट र बेन्च प्रेस जस्ता कम्पाउन्ड मुभमेन्टहरूमा ध्यान दिनुहोस्। प्रगति गर्नुअघि फर्म सिक्न हल्का तौलबाट सुरु गर्नुहोस्। सेटहरू बीच ६०-९० सेकेन्ड आराम गर्नुहोस्। सधैं ५-१० मिनेट हल्का कार्डियोसँग वार्म अप गर्नुहोस्। चोटपटक रोक्न र लचिलोपन सुधार गर्न स्ट्रेचिङसँग कूल डाउन गर्नुहोस्।',
     'beginner', true),

    ('e1b2c3d4-e5f6-7890-abcd-ef1234567802',
     'Building Strength: Progressive Overload',
     'शक्ति निर्माण: प्रगतिशील ओभरलोड',
     'Progressive overload is the foundation of strength training. To get stronger, you must gradually increase the demands on your muscles. This can be done by adding weight, increasing reps, adding sets, or reducing rest time. Track your lifts in a journal. Aim to increase weight by 2.5-5kg when you can complete all prescribed reps with good form. Focus on the big five: squat, deadlift, bench press, overhead press, and barbell row. Eat adequate protein (1.6-2.2g per kg of body weight) to support muscle recovery.',
     'प्रगतिशील ओभरलोड शक्ति प्रशिक्षणको आधार हो। बलियो हुन, तपाईंले आफ्नो मांसपेशीहरूमा बिस्तारै माग बढाउनुपर्छ। यो तौल थप्ने, रेप्स बढाउने, सेटहरू थप्ने वा आराम समय घटाएर गर्न सकिन्छ। जर्नलमा आफ्नो लिफ्टहरू ट्र्याक गर्नुहोस्। राम्रो फर्ममा सबै निर्धारित रेप्स पूरा गर्न सक्नुहुँदा २.५-५ किलो तौल बढाउने लक्ष्य राख्नुहोस्। पाँच ठूला व्यायाममा ध्यान दिनुहोस्: स्क्वाट, डेडलिफ्ट, बेन्च प्रेस, ओभरहेड प्रेस र बारबेल रो। मांसपेशी रिकभरी सपोर्ट गर्न पर्याप्त प्रोटीन (शरीरको तौलको प्रति किलो १.६-२.२ ग्राम) खानुहोस्।',
     'strength', true),

    ('e1b2c3d4-e5f6-7890-abcd-ef1234567803',
     'Cardio Training for Fat Loss',
     'चर्बी घटाउनको लागि कार्डियो प्रशिक्षण',
     'Effective cardio training combines both steady-state and high-intensity interval training (HIIT). For steady-state, aim for 30-45 minutes at 60-70% of your max heart rate, 2-3 times per week. For HIIT, alternate between 30 seconds of all-out effort and 60-90 seconds of rest for 15-20 minutes. Walking is underrated — aim for 8,000-10,000 steps daily. Remember: you cannot out-train a bad diet. Maintain a moderate caloric deficit of 300-500 calories for sustainable fat loss.',
     'प्रभावकारी कार्डियो प्रशिक्षणले स्थिर-अवस्था र उच्च-तीव्रता अन्तराल प्रशिक्षण (HIIT) दुवै जोड्दछ। स्थिर-अवस्थाको लागि, हप्ताको २-३ पटक आफ्नो अधिकतम हृदय गतिको ६०-७०% मा ३०-४५ मिनेट लक्ष्य राख्नुहोस्। HIIT को लागि, १५-२० मिनेटको लागि ३० सेकेन्ड पूर्ण प्रयास र ६०-९० सेकेन्ड आरामको बीचमा वैकल्पिक गर्नुहोस्। हिँड्नु कम मूल्याङ्कन गरिएको छ — दैनिक ८,०००-१०,००० कदमको लक्ष्य राख्नुहोस्। याद गर्नुहोस्: खराब आहारलाई प्रशिक्षणले जित्न सक्दैन। दिगो चर्बी घटाउनको लागि ३००-५०० क्यालोरीको मध्यम क्यालोरिक घाटा कायम राख्नुहोस्।',
     'cardio', true),

    ('e1b2c3d4-e5f6-7890-abcd-ef1234567804',
     'Nutrition Basics for Gym-Goers',
     'जिम जानेहरूको लागि पोषण आधारभूत कुराहरू',
     'Nutrition fuels your training and recovery. Focus on whole foods: lean proteins (chicken, fish, eggs, dal), complex carbs (rice, roti, oats), and healthy fats (nuts, ghee, olive oil). Eat protein with every meal — aim for 4-5 meals spread throughout the day. Pre-workout: eat a balanced meal 2-3 hours before, or a light snack 30-60 minutes before. Post-workout: consume protein and carbs within 2 hours. Stay hydrated — drink at least 2-3 liters of water daily. Avoid processed foods and sugary drinks.',
     'पोषणले तपाईंको प्रशिक्षण र रिकभरीलाई इन्धन दिन्छ। सम्पूर्ण खानामा ध्यान दिनुहोस्: पातलो प्रोटीनहरू (कुखुरा, माछा, अण्डा, दाल), जटिल कार्ब्स (भात, रोटी, ओट्स), र स्वस्थ बोसो (नट, घ्यू, जैतून तेल)। हरेक खानामा प्रोटीन खानुहोस् — दिनभर फैलिएको ४-५ खाना लक्ष्य राख्नुहोस्। प्रि-वर्कआउट: २-३ घण्टा अघि सन्तुलित खाना, वा ३०-६० मिनेट अघि हल्का नास्ता खानुहोस्। पोस्ट-वर्कआउट: २ घण्टा भित्र प्रोटीन र कार्ब्स खानुहोस्। हाइड्रेटेड रहनुहोस् — दैनिक कम्तिमा २-३ लिटर पानी पिउनुहोस्। प्रशोधित खाना र चिनीयुक्त पेयहरूबाट टाढा रहनुहोस्।',
     'nutrition', true),

    ('e1b2c3d4-e5f6-7890-abcd-ef1234567805',
     'Recovery and Rest: The Missing Piece',
     'रिकभरी र आराम: हराइरहेको टुक्रा',
     'Recovery is when your muscles actually grow. Without proper rest, you risk overtraining and injury. Sleep 7-9 hours per night — this is when growth hormone peaks. Take 1-2 rest days per week. Active recovery (light walking, stretching, yoga) on rest days promotes blood flow and healing. Foam rolling helps release muscle tension. If a muscle group is still sore, do not train it — wait until soreness subsides. Listen to your body. Deload every 4-6 weeks by reducing weight or volume by 40-50% to prevent burnout.',
     'रिकभरी भनेको तपाईंको मांसपेशीहरू वास्तवमा बढ्ने समय हो। उचित आरामविना, तपाईं ओभरट्रेनिंग र चोटपटकको जोखिममा हुनुहुन्छ। प्रति रात ७-९ घण्टा सुत्नुहोस् — यो समयमा ग्रोथ हर्मोन चरम हुन्छ। हप्ताको १-२ आराम दिन लिनुहोस्। आराम दिनहरूमा सक्रिय रिकभरी (हल्का हिँडाइ, स्ट्रेचिङ, योग) ले रक्त प्रवाह र उपचारलाई बढावा दिन्छ। फोम रोलिंगले मांसपेशीको तनाव मुक्त गर्न मद्दत गर्छ। यदि मांसपेशी समूह अझै दुखेको छ भने, प्रशिक्षण नगर्नुहोस् — दुखाइ कम नभएसम्म पर्खनुहोस्। आफ्नो शरीरको कुरा सुन्नुहोस्। बर्नआउट रोक्न हरेक ४-६ हप्तामा तौल वा भोल्युम ४०-५०% ले घटाएर डिलोड गर्नुहोस्।',
     'recovery', true);
