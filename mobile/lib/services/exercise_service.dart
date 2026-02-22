import '../models/exercise.dart';
import 'api_client.dart';

class ExerciseService {
  final ApiClient _api;
  ExerciseService(this._api);

  Future<List<ExerciseItem>> listExercises({String? type, String? equipment, String? muscleGroup, int limit = 50}) async {
    final params = <String, dynamic>{'limit': limit};
    if (type != null) params['type'] = type;
    if (equipment != null) params['equipment'] = equipment;
    if (muscleGroup != null) params['muscle_group'] = muscleGroup;
    final response = await _api.get('/members/me/exercises', queryParameters: params);
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data.map((e) => ExerciseItem.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<ExerciseItem> getExercise(String id) async {
    final response = await _api.get('/members/me/exercises/$id');
    return ExerciseItem.fromJson(response.data as Map<String, dynamic>);
  }

  Future<List<WorkoutProgram>> listPrograms({String? type, String? difficulty, String? equipment}) async {
    final params = <String, dynamic>{};
    if (type != null) params['type'] = type;
    if (difficulty != null) params['difficulty'] = difficulty;
    if (equipment != null) params['equipment'] = equipment;
    final response = await _api.get('/members/me/programs', queryParameters: params);
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data.map((e) => WorkoutProgram.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<WorkoutProgram> getProgram(String id) async {
    final response = await _api.get('/members/me/programs/$id');
    return WorkoutProgram.fromJson(response.data as Map<String, dynamic>);
  }

  Future<UserProgramEnrollment> enrollInProgram(String programId) async {
    final response = await _api.post('/members/me/programs/$programId/enroll');
    return UserProgramEnrollment.fromJson(response.data as Map<String, dynamic>);
  }

  Future<UserProgramEnrollment> getCurrentProgram() async {
    final response = await _api.get('/members/me/programs/current');
    return UserProgramEnrollment.fromJson(response.data as Map<String, dynamic>);
  }

  Future<UserProgramEnrollment> completeProgramDay() async {
    final response = await _api.post('/members/me/programs/current/complete-day');
    return UserProgramEnrollment.fromJson(response.data as Map<String, dynamic>);
  }
}
