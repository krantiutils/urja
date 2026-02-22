class ExerciseItem {
  final String id;
  final String name;
  final String? nameNe;
  final String? description;
  final List<String> muscleGroups;
  final String equipment;
  final String exerciseType;
  final String difficulty;
  final double caloriesPerMinute;
  final String? instructions;
  final String? imageUrl;
  final String createdAt;

  ExerciseItem({
    required this.id,
    required this.name,
    this.nameNe,
    this.description,
    this.muscleGroups = const [],
    required this.equipment,
    required this.exerciseType,
    required this.difficulty,
    this.caloriesPerMinute = 0,
    this.instructions,
    this.imageUrl,
    required this.createdAt,
  });

  factory ExerciseItem.fromJson(Map<String, dynamic> json) => ExerciseItem(
    id: json['id'] as String,
    name: json['name'] as String,
    nameNe: json['name_ne'] as String?,
    description: json['description'] as String?,
    muscleGroups: (json['muscle_groups'] as List<dynamic>?)?.map((e) => e as String).toList() ?? [],
    equipment: json['equipment'] as String? ?? 'none',
    exerciseType: json['exercise_type'] as String? ?? 'strength',
    difficulty: json['difficulty'] as String? ?? 'beginner',
    caloriesPerMinute: (json['calories_per_minute'] as num?)?.toDouble() ?? 0,
    instructions: json['instructions'] as String?,
    imageUrl: json['image_url'] as String?,
    createdAt: json['created_at'] as String? ?? '',
  );
}

class WorkoutProgram {
  final String id;
  final String name;
  final String? nameNe;
  final String? description;
  final String? programType;
  final String difficulty;
  final int durationWeeks;
  final String equipment;
  final String? goal;
  final String? imageUrl;
  final String createdAt;
  final List<ProgramDay>? days;

  WorkoutProgram({
    required this.id,
    required this.name,
    this.nameNe,
    this.description,
    this.programType,
    required this.difficulty,
    required this.durationWeeks,
    required this.equipment,
    this.goal,
    this.imageUrl,
    required this.createdAt,
    this.days,
  });

  factory WorkoutProgram.fromJson(Map<String, dynamic> json) => WorkoutProgram(
    id: json['id'] as String,
    name: json['name'] as String,
    nameNe: json['name_ne'] as String?,
    description: json['description'] as String?,
    programType: json['program_type'] as String?,
    difficulty: json['difficulty'] as String? ?? 'beginner',
    durationWeeks: json['duration_weeks'] as int? ?? 4,
    equipment: json['equipment'] as String? ?? 'none',
    goal: json['goal'] as String?,
    imageUrl: json['image_url'] as String?,
    createdAt: json['created_at'] as String? ?? '',
    days: (json['days'] as List<dynamic>?)?.map((e) => ProgramDay.fromJson(e as Map<String, dynamic>)).toList(),
  );
}

class ProgramDay {
  final String id;
  final String programId;
  final int weekNumber;
  final int dayNumber;
  final String? dayName;
  final List<dynamic> exercises;
  final bool restDay;

  ProgramDay({
    required this.id,
    required this.programId,
    required this.weekNumber,
    required this.dayNumber,
    this.dayName,
    this.exercises = const [],
    this.restDay = false,
  });

  factory ProgramDay.fromJson(Map<String, dynamic> json) => ProgramDay(
    id: json['id'] as String,
    programId: json['program_id'] as String,
    weekNumber: json['week_number'] as int? ?? 1,
    dayNumber: json['day_number'] as int? ?? 1,
    dayName: json['day_name'] as String?,
    exercises: json['exercises'] as List<dynamic>? ?? [],
    restDay: json['rest_day'] as bool? ?? false,
  );
}

class UserProgramEnrollment {
  final String id;
  final String userId;
  final String programId;
  final String startedAt;
  final int currentWeek;
  final int currentDay;
  final String status;
  final WorkoutProgram? program;
  final ProgramDay? todayWorkout;

  UserProgramEnrollment({
    required this.id,
    required this.userId,
    required this.programId,
    required this.startedAt,
    required this.currentWeek,
    required this.currentDay,
    required this.status,
    this.program,
    this.todayWorkout,
  });

  factory UserProgramEnrollment.fromJson(Map<String, dynamic> json) => UserProgramEnrollment(
    id: json['id'] as String,
    userId: json['user_id'] as String,
    programId: json['program_id'] as String,
    startedAt: json['started_at'] as String? ?? '',
    currentWeek: json['current_week'] as int? ?? 1,
    currentDay: json['current_day'] as int? ?? 1,
    status: json['status'] as String? ?? 'active',
    program: json['program'] != null ? WorkoutProgram.fromJson(json['program'] as Map<String, dynamic>) : null,
    todayWorkout: json['today_workout'] != null ? ProgramDay.fromJson(json['today_workout'] as Map<String, dynamic>) : null,
  );
}
