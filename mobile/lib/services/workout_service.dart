import '../models/workout.dart';
import 'api_client.dart';

class WorkoutService {
  final ApiClient _api;

  WorkoutService(this._api);

  Future<List<WorkoutTemplate>> getTemplates(String orgId) async {
    final response = await _api.get('/orgs/$orgId/workout-templates');
    final data = response.data['data'] as List<dynamic>? ??
        (response.data is List ? response.data as List<dynamic> : []);
    return data
        .map((e) => WorkoutTemplate.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<WorkoutTemplate> addTemplate(
      String orgId, Map<String, dynamic> data) async {
    final response =
        await _api.post('/orgs/$orgId/workout-templates', data: data);
    return WorkoutTemplate.fromJson(response.data);
  }

  Future<List<WorkoutTemplate>> browseTemplates(String orgId,
      {String? goal, String? difficulty}) async {
    final params = <String, dynamic>{'organization_id': orgId};
    if (goal != null) params['goal'] = goal;
    if (difficulty != null) params['difficulty'] = difficulty;
    final response = await _api.get('/members/me/workout-templates',
        queryParameters: params);
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => WorkoutTemplate.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<WorkoutPlan> selfAssignPlan(String orgId, String templateId) async {
    final response = await _api.post('/members/me/workout-plan', data: {
      'organization_id': orgId,
      'workout_template_id': templateId,
    });
    return WorkoutPlan.fromJson(response.data as Map<String, dynamic>);
  }

  Future<WorkoutPlan> recommendPlan(
      String orgId, String goal, String experience, String daysPerWeek) async {
    final response =
        await _api.post('/members/me/workout-plan/recommend', data: {
      'organization_id': orgId,
      'goal': goal,
      'experience': experience,
      'days_per_week': daysPerWeek,
    });
    return WorkoutPlan.fromJson(response.data as Map<String, dynamic>);
  }
}
