class Exercise {
  final String name;
  final String? nameNe;
  final int sets;
  final int reps;
  final int restSeconds;
  final String? notes;
  final List<String> muscleGroups;

  Exercise({
    required this.name,
    this.nameNe,
    required this.sets,
    required this.reps,
    required this.restSeconds,
    this.notes,
    this.muscleGroups = const [],
  });

  factory Exercise.fromJson(Map<String, dynamic> json) => Exercise(
        name: json['name'] as String? ?? '',
        nameNe: json['name_ne'] as String?,
        sets: json['sets'] as int? ?? 0,
        reps: json['reps'] as int? ?? 0,
        restSeconds: json['rest_seconds'] as int? ?? 0,
        notes: json['notes'] as String?,
        muscleGroups: (json['muscle_groups'] as List<dynamic>?)
                ?.map((e) => e as String)
                .toList() ??
            [],
      );
}

class WorkoutTemplate {
  final String id;
  final String name;
  final String? nameNe;
  final String? description;
  final String? category;
  final String? difficulty;
  final int? durationMinutes;
  final String? goal;
  final String? daysPerWeek;
  final List<Exercise> exercises;

  WorkoutTemplate({
    required this.id,
    required this.name,
    this.nameNe,
    this.description,
    this.category,
    this.difficulty,
    this.durationMinutes,
    this.goal,
    this.daysPerWeek,
    required this.exercises,
  });

  factory WorkoutTemplate.fromJson(Map<String, dynamic> json) =>
      WorkoutTemplate(
        id: json['id'] as String,
        name: json['name'] as String,
        nameNe: json['name_ne'] as String?,
        description: json['description'] as String?,
        category: json['category'] as String?,
        difficulty: json['difficulty'] as String?,
        durationMinutes: json['duration_minutes'] as int?,
        goal: json['goal'] as String?,
        daysPerWeek: json['days_per_week'] as String?,
        exercises: (json['exercises'] as List<dynamic>?)
                ?.map(
                    (e) => Exercise.fromJson(e as Map<String, dynamic>))
                .toList() ??
            [],
      );
}

class WorkoutPlan {
  final String id;
  final String userId;
  final String organizationId;
  final String workoutTemplateId;
  final String? assignedBy;
  final String? assignmentMethod;
  final String assignedAt;
  final WorkoutTemplate? template;

  WorkoutPlan({
    required this.id,
    required this.userId,
    required this.organizationId,
    required this.workoutTemplateId,
    this.assignedBy,
    this.assignmentMethod,
    required this.assignedAt,
    this.template,
  });

  factory WorkoutPlan.fromJson(Map<String, dynamic> json) => WorkoutPlan(
        id: json['id'] as String,
        userId: json['user_id'] as String,
        organizationId: json['organization_id'] as String,
        workoutTemplateId: json['workout_template_id'] as String,
        assignedBy: json['assigned_by'] as String?,
        assignmentMethod: json['assignment_method'] as String?,
        assignedAt: json['assigned_at'] as String,
        template: json['template'] != null
            ? WorkoutTemplate.fromJson(
                json['template'] as Map<String, dynamic>)
            : null,
      );
}

class WorkoutLog {
  final String id;
  final String userId;
  final String? workoutTemplateId;
  final String organizationId;
  final List<dynamic> exercises;
  final int? durationMinutes;
  final String notes;
  final String loggedAt;

  WorkoutLog({
    required this.id,
    required this.userId,
    this.workoutTemplateId,
    required this.organizationId,
    required this.exercises,
    this.durationMinutes,
    required this.notes,
    required this.loggedAt,
  });

  factory WorkoutLog.fromJson(Map<String, dynamic> json) => WorkoutLog(
        id: json['id'] as String,
        userId: json['user_id'] as String,
        workoutTemplateId: json['workout_template_id'] as String?,
        organizationId: json['organization_id'] as String,
        exercises: json['exercises'] as List<dynamic>? ?? [],
        durationMinutes: json['duration_minutes'] as int?,
        notes: json['notes'] as String? ?? '',
        loggedAt: json['logged_at'] as String,
      );
}
