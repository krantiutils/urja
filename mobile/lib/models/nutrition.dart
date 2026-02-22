class FoodItem {
  final String id;
  final String? organizationId;
  final String name;
  final String? nameNe;
  final String category;
  final double caloriesPer100g;
  final double proteinPer100g;
  final double carbsPer100g;
  final double fatPer100g;
  final double fiberPer100g;
  final double? servingSizeG;
  final String? servingLabel;
  final String? servingLabelNe;
  final bool isVerified;

  FoodItem({
    required this.id,
    this.organizationId,
    required this.name,
    this.nameNe,
    required this.category,
    required this.caloriesPer100g,
    required this.proteinPer100g,
    required this.carbsPer100g,
    required this.fatPer100g,
    required this.fiberPer100g,
    this.servingSizeG,
    this.servingLabel,
    this.servingLabelNe,
    this.isVerified = false,
  });

  factory FoodItem.fromJson(Map<String, dynamic> json) => FoodItem(
        id: json['id'] as String,
        organizationId: json['organization_id'] as String?,
        name: json['name'] as String,
        nameNe: json['name_ne'] as String?,
        category: json['category'] as String? ?? '',
        caloriesPer100g: (json['calories_per_100g'] as num?)?.toDouble() ?? 0,
        proteinPer100g: (json['protein_per_100g'] as num?)?.toDouble() ?? 0,
        carbsPer100g: (json['carbs_per_100g'] as num?)?.toDouble() ?? 0,
        fatPer100g: (json['fat_per_100g'] as num?)?.toDouble() ?? 0,
        fiberPer100g: (json['fiber_per_100g'] as num?)?.toDouble() ?? 0,
        servingSizeG: (json['serving_size_g'] as num?)?.toDouble(),
        servingLabel: json['serving_label'] as String?,
        servingLabelNe: json['serving_label_ne'] as String?,
        isVerified: json['is_verified'] as bool? ?? false,
      );
}

class FoodLog {
  final String id;
  final String userId;
  final String organizationId;
  final String foodItemId;
  final FoodItem? foodItem;
  final String mealType;
  final double quantityGrams;
  final double calories;
  final double protein;
  final double carbs;
  final double fat;
  final String loggedDate;
  final String loggedAt;
  final String? notes;

  FoodLog({
    required this.id,
    required this.userId,
    required this.organizationId,
    required this.foodItemId,
    this.foodItem,
    required this.mealType,
    required this.quantityGrams,
    required this.calories,
    required this.protein,
    required this.carbs,
    required this.fat,
    required this.loggedDate,
    required this.loggedAt,
    this.notes,
  });

  factory FoodLog.fromJson(Map<String, dynamic> json) => FoodLog(
        id: json['id'] as String,
        userId: json['user_id'] as String,
        organizationId: json['organization_id'] as String,
        foodItemId: json['food_item_id'] as String,
        foodItem: json['food_item'] != null
            ? FoodItem.fromJson(json['food_item'] as Map<String, dynamic>)
            : null,
        mealType: json['meal_type'] as String? ?? '',
        quantityGrams: (json['quantity_grams'] as num?)?.toDouble() ?? 0,
        calories: (json['calories'] as num?)?.toDouble() ?? 0,
        protein: (json['protein'] as num?)?.toDouble() ?? 0,
        carbs: (json['carbs'] as num?)?.toDouble() ?? 0,
        fat: (json['fat'] as num?)?.toDouble() ?? 0,
        loggedDate: json['logged_date'] as String? ?? '',
        loggedAt: json['logged_at'] as String? ?? '',
        notes: json['notes'] as String?,
      );
}

class MealTemplate {
  final String id;
  final String userId;
  final String organizationId;
  final String name;
  final String? nameNe;
  final String mealType;
  final List<dynamic> items;
  final String createdAt;

  MealTemplate({
    required this.id,
    required this.userId,
    required this.organizationId,
    required this.name,
    this.nameNe,
    required this.mealType,
    required this.items,
    required this.createdAt,
  });

  factory MealTemplate.fromJson(Map<String, dynamic> json) => MealTemplate(
        id: json['id'] as String,
        userId: json['user_id'] as String,
        organizationId: json['organization_id'] as String,
        name: json['name'] as String,
        nameNe: json['name_ne'] as String?,
        mealType: json['meal_type'] as String? ?? '',
        items: json['items'] as List<dynamic>? ?? [],
        createdAt: json['created_at'] as String? ?? '',
      );
}

class NutritionGoal {
  final String id;
  final String userId;
  final String organizationId;
  final double calorieGoal;
  final double proteinGoalG;
  final double carbsGoalG;
  final double fatGoalG;
  final double weightKg;
  final double heightCm;
  final int age;
  final String gender;
  final String activityLevel;
  final String goalType;

  NutritionGoal({
    required this.id,
    required this.userId,
    required this.organizationId,
    required this.calorieGoal,
    required this.proteinGoalG,
    required this.carbsGoalG,
    required this.fatGoalG,
    required this.weightKg,
    required this.heightCm,
    required this.age,
    required this.gender,
    required this.activityLevel,
    required this.goalType,
  });

  factory NutritionGoal.fromJson(Map<String, dynamic> json) => NutritionGoal(
        id: json['id'] as String,
        userId: json['user_id'] as String,
        organizationId: json['organization_id'] as String,
        calorieGoal: (json['calorie_goal'] as num?)?.toDouble() ?? 0,
        proteinGoalG: (json['protein_goal_g'] as num?)?.toDouble() ?? 0,
        carbsGoalG: (json['carbs_goal_g'] as num?)?.toDouble() ?? 0,
        fatGoalG: (json['fat_goal_g'] as num?)?.toDouble() ?? 0,
        weightKg: (json['weight_kg'] as num?)?.toDouble() ?? 0,
        heightCm: (json['height_cm'] as num?)?.toDouble() ?? 0,
        age: json['age'] as int? ?? 0,
        gender: json['gender'] as String? ?? '',
        activityLevel: json['activity_level'] as String? ?? '',
        goalType: json['goal_type'] as String? ?? '',
      );
}

class MealSummary {
  final String mealType;
  final double calories;
  final double protein;
  final double carbs;
  final double fat;
  final List<FoodLog> items;

  MealSummary({
    required this.mealType,
    required this.calories,
    required this.protein,
    required this.carbs,
    required this.fat,
    required this.items,
  });

  factory MealSummary.fromJson(Map<String, dynamic> json) => MealSummary(
        mealType: json['meal_type'] as String? ?? '',
        calories: (json['calories'] as num?)?.toDouble() ?? 0,
        protein: (json['protein'] as num?)?.toDouble() ?? 0,
        carbs: (json['carbs'] as num?)?.toDouble() ?? 0,
        fat: (json['fat'] as num?)?.toDouble() ?? 0,
        items: (json['items'] as List<dynamic>?)
                ?.map((e) => FoodLog.fromJson(e as Map<String, dynamic>))
                .toList() ??
            [],
      );
}

class DailySummary {
  final String date;
  final double totalCalories;
  final double totalProtein;
  final double totalCarbs;
  final double totalFat;
  final List<MealSummary> meals;

  DailySummary({
    required this.date,
    required this.totalCalories,
    required this.totalProtein,
    required this.totalCarbs,
    required this.totalFat,
    required this.meals,
  });

  factory DailySummary.fromJson(Map<String, dynamic> json) => DailySummary(
        date: json['date'] as String? ?? '',
        totalCalories: (json['total_calories'] as num?)?.toDouble() ?? 0,
        totalProtein: (json['total_protein'] as num?)?.toDouble() ?? 0,
        totalCarbs: (json['total_carbs'] as num?)?.toDouble() ?? 0,
        totalFat: (json['total_fat'] as num?)?.toDouble() ?? 0,
        meals: (json['meals'] as List<dynamic>?)
                ?.map((e) => MealSummary.fromJson(e as Map<String, dynamic>))
                .toList() ??
            [],
      );
}

class WeeklySummaryDay {
  final String date;
  final double totalCalories;
  final double totalProtein;
  final double totalCarbs;
  final double totalFat;

  WeeklySummaryDay({
    required this.date,
    required this.totalCalories,
    required this.totalProtein,
    required this.totalCarbs,
    required this.totalFat,
  });

  factory WeeklySummaryDay.fromJson(Map<String, dynamic> json) =>
      WeeklySummaryDay(
        date: json['date'] as String? ?? '',
        totalCalories: (json['total_calories'] as num?)?.toDouble() ?? 0,
        totalProtein: (json['total_protein'] as num?)?.toDouble() ?? 0,
        totalCarbs: (json['total_carbs'] as num?)?.toDouble() ?? 0,
        totalFat: (json['total_fat'] as num?)?.toDouble() ?? 0,
      );
}

class NutritionStreak {
  final String userId;
  final int currentStreak;
  final int longestStreak;
  final String? lastLogDate;
  final String updatedAt;

  NutritionStreak({
    required this.userId,
    required this.currentStreak,
    required this.longestStreak,
    this.lastLogDate,
    required this.updatedAt,
  });

  factory NutritionStreak.fromJson(Map<String, dynamic> json) => NutritionStreak(
    userId: json['user_id'] as String,
    currentStreak: json['current_streak'] as int? ?? 0,
    longestStreak: json['longest_streak'] as int? ?? 0,
    lastLogDate: json['last_log_date'] as String?,
    updatedAt: json['updated_at'] as String? ?? '',
  );
}

class CreateCustomFoodInput {
  final String name;
  final String? nameNe;
  final String? category;
  final double caloriesPer100g;
  final double proteinPer100g;
  final double carbsPer100g;
  final double fatPer100g;
  final double? fiberPer100g;
  final double? servingSizeG;
  final String? servingLabel;
  final String? barcode;

  CreateCustomFoodInput({
    required this.name,
    this.nameNe,
    this.category,
    required this.caloriesPer100g,
    required this.proteinPer100g,
    required this.carbsPer100g,
    required this.fatPer100g,
    this.fiberPer100g,
    this.servingSizeG,
    this.servingLabel,
    this.barcode,
  });

  Map<String, dynamic> toJson() => {
    'name': name,
    if (nameNe != null) 'name_ne': nameNe,
    if (category != null) 'category': category,
    'calories_per_100g': caloriesPer100g,
    'protein_per_100g': proteinPer100g,
    'carbs_per_100g': carbsPer100g,
    'fat_per_100g': fatPer100g,
    if (fiberPer100g != null) 'fiber_per_100g': fiberPer100g,
    if (servingSizeG != null) 'serving_size_g': servingSizeG,
    if (servingLabel != null) 'serving_label': servingLabel,
    if (barcode != null) 'barcode': barcode,
  };
}
