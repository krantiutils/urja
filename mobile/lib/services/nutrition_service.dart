import '../models/nutrition.dart';
import 'api_client.dart';

class NutritionService {
  final ApiClient _api;

  NutritionService(this._api);

  Future<List<FoodItem>> searchFoods(String orgId,
      {String? query, String? category, int limit = 20}) async {
    final params = <String, dynamic>{
      'organization_id': orgId,
      'limit': limit,
    };
    if (query != null && query.isNotEmpty) params['q'] = query;
    if (category != null) params['category'] = category;
    final response =
        await _api.get('/members/me/foods', queryParameters: params);
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => FoodItem.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<FoodItem> getFoodItem(String id) async {
    final response = await _api.get('/members/me/foods/$id');
    return FoodItem.fromJson(response.data as Map<String, dynamic>);
  }

  Future<FoodLog> logFood({
    required String orgId,
    required String foodItemId,
    required String mealType,
    required double quantityGrams,
    required String loggedDate,
    String? notes,
  }) async {
    final response = await _api.post('/members/me/food-logs', data: {
      'organization_id': orgId,
      'food_item_id': foodItemId,
      'meal_type': mealType,
      'quantity_grams': quantityGrams,
      'logged_date': loggedDate,
      if (notes != null) 'notes': notes,
    });
    return FoodLog.fromJson(response.data as Map<String, dynamic>);
  }

  Future<List<FoodLog>> getMyFoodLogs(String orgId,
      {String? date, String? mealType}) async {
    final params = <String, dynamic>{'organization_id': orgId};
    if (date != null) params['date'] = date;
    if (mealType != null) params['meal_type'] = mealType;
    final response =
        await _api.get('/members/me/food-logs', queryParameters: params);
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => FoodLog.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<void> deleteFoodLog(String id) async {
    await _api.deleteRequest('/members/me/food-logs/$id');
  }

  Future<DailySummary> getDailySummary(String orgId, String date) async {
    final response =
        await _api.get('/members/me/nutrition/summary', queryParameters: {
      'organization_id': orgId,
      'date': date,
    });
    return DailySummary.fromJson(response.data as Map<String, dynamic>);
  }

  Future<List<WeeklySummaryDay>> getWeeklySummary(
      String orgId, String from) async {
    final response =
        await _api.get('/members/me/nutrition/weekly', queryParameters: {
      'organization_id': orgId,
      'from': from,
    });
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => WeeklySummaryDay.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<NutritionGoal> setNutritionGoal({
    required String orgId,
    required double weightKg,
    required double heightCm,
    required int age,
    required String gender,
    required String activityLevel,
    required String goalType,
  }) async {
    final response =
        await _api.post('/members/me/nutrition/goal', data: {
      'organization_id': orgId,
      'weight_kg': weightKg,
      'height_cm': heightCm,
      'age': age,
      'gender': gender,
      'activity_level': activityLevel,
      'goal_type': goalType,
    });
    return NutritionGoal.fromJson(response.data as Map<String, dynamic>);
  }

  Future<NutritionGoal> getNutritionGoal(String orgId) async {
    final response =
        await _api.get('/members/me/nutrition/goal', queryParameters: {
      'organization_id': orgId,
    });
    return NutritionGoal.fromJson(response.data as Map<String, dynamic>);
  }

  Future<MealTemplate> createMealTemplate({
    required String orgId,
    required String name,
    String? nameNe,
    required String mealType,
    required List<Map<String, dynamic>> items,
  }) async {
    final response = await _api.post('/members/me/meal-templates', data: {
      'organization_id': orgId,
      'name': name,
      if (nameNe != null) 'name_ne': nameNe,
      'meal_type': mealType,
      'items': items,
    });
    return MealTemplate.fromJson(response.data as Map<String, dynamic>);
  }

  Future<List<MealTemplate>> getMyMealTemplates(String orgId) async {
    final response =
        await _api.get('/members/me/meal-templates', queryParameters: {
      'organization_id': orgId,
    });
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => MealTemplate.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<void> deleteMealTemplate(String id) async {
    await _api.deleteRequest('/members/me/meal-templates/$id');
  }

  // Custom food creation
  Future<FoodItem> createCustomFood(CreateCustomFoodInput input) async {
    final response = await _api.post('/members/me/foods', data: input.toJson());
    return FoodItem.fromJson(response.data as Map<String, dynamic>);
  }

  // Barcode lookup
  Future<FoodItem> lookupBarcode(String barcode) async {
    final response = await _api.get('/members/me/foods/barcode/$barcode');
    return FoodItem.fromJson(response.data as Map<String, dynamic>);
  }

  // Nutrition streak
  Future<NutritionStreak> getNutritionStreak() async {
    final response = await _api.get('/members/me/nutrition/streak');
    return NutritionStreak.fromJson(response.data as Map<String, dynamic>);
  }

  Future<List<FoodLog>> logMealTemplate(
      String templateId, String orgId, String loggedDate) async {
    final response =
        await _api.post('/members/me/meal-templates/$templateId/log', data: {
      'organization_id': orgId,
      'logged_date': loggedDate,
    });
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => FoodLog.fromJson(e as Map<String, dynamic>))
        .toList();
  }
}
