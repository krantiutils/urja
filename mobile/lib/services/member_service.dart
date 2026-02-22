import '../models/attendance.dart';
import '../models/calendar.dart';
import '../models/health.dart';
import '../models/leaderboard.dart';
import '../models/member.dart';
import '../models/package.dart';
import '../models/streak.dart';
import '../models/weight.dart';
import '../models/workout.dart';
import 'api_client.dart';

class MemberService {
  final ApiClient _api;

  MemberService(this._api);

  Future<MemberProfile> getProfile() async {
    final response = await _api.get('/members/me');
    return MemberProfile.fromJson(response.data);
  }

  Future<void> updateProfile(Map<String, dynamic> data) async {
    await _api.put('/members/me', data: data);
  }

  Future<void> updatePrivacy(Map<String, dynamic> data) async {
    await _api.put('/members/me/privacy', data: data);
  }

  Future<List<MemberAttendanceRecord>> getAttendance({
    int limit = 20,
    int offset = 0,
  }) async {
    final response = await _api.get('/members/me/attendance', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) =>
            MemberAttendanceRecord.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<List<MemberStreak>> getStreaks() async {
    final response = await _api.get('/members/me/streaks');
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => MemberStreak.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<List<MemberPackage>> getPackages({
    int limit = 20,
    int offset = 0,
  }) async {
    final response = await _api.get('/members/me/packages', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => MemberPackage.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<List<HealthMetric>> getHealth({
    String type = 'bmi',
    int limit = 20,
  }) async {
    final response = await _api.get('/members/me/health', queryParameters: {
      'type': type,
      'limit': limit,
    });
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => HealthMetric.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<List<WorkoutLog>> getWorkoutLogs({
    required String organizationId,
    int limit = 20,
    int offset = 0,
  }) async {
    final response =
        await _api.get('/members/me/workout-logs', queryParameters: {
      'organization_id': organizationId,
      'limit': limit,
      'offset': offset,
    });
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => WorkoutLog.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<WorkoutPlan?> getWorkoutPlan({
    required String organizationId,
  }) async {
    try {
      final response =
          await _api.get('/members/me/workout-plan', queryParameters: {
        'organization_id': organizationId,
      });
      return WorkoutPlan.fromJson(response.data);
    } catch (_) {
      return null;
    }
  }

  Future<AttendanceCalendar> getAttendanceCalendar({String? month}) async {
    final queryParams = <String, dynamic>{};
    if (month != null) queryParams['month'] = month;
    final response = await _api.get('/members/me/attendance/calendar',
        queryParameters: queryParams);
    return AttendanceCalendar.fromJson(response.data);
  }

  Future<LeaderboardResponse> getLeaderboard({
    String period = 'monthly',
    int limit = 3,
  }) async {
    final response =
        await _api.get('/members/me/leaderboard', queryParameters: {
      'period': period,
      'limit': limit.toString(),
    });
    return LeaderboardResponse.fromJson(response.data);
  }

  Future<HealthMetric> quickLogWeight(double weightKg) async {
    final response = await _api.post('/members/me/health/weight', data: {
      'weight_kg': weightKg,
    });
    return HealthMetric.fromJson(response.data as Map<String, dynamic>);
  }

  Future<WeightTrend> getWeightTrend({int days = 30}) async {
    final response = await _api
        .get('/members/me/health/weight-trend', queryParameters: {
      'days': days,
    });
    return WeightTrend.fromJson(response.data as Map<String, dynamic>);
  }
}
