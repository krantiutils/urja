import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../models/exercise.dart';
import '../../providers/auth_provider.dart';
import '../../services/exercise_service.dart';
import '../shared/widgets/empty_state.dart';
import '../shared/widgets/loading_indicator.dart';

// ---------------------------------------------------------------------------
// Muscle group color mapping
// ---------------------------------------------------------------------------
Color _muscleGroupColor(String group) {
  switch (group.toLowerCase()) {
    case 'chest':
    case 'back':
      return Colors.blue;
    case 'shoulders':
    case 'traps':
      return Colors.indigo;
    case 'biceps':
    case 'triceps':
    case 'forearms':
      return Colors.purple;
    case 'core':
    case 'obliques':
      return Colors.orange;
    case 'quadriceps':
    case 'hamstrings':
    case 'glutes':
      return Colors.teal;
    case 'calves':
      return Colors.cyan;
    case 'hip_flexors':
    case 'lower_back':
    case 'lats':
      return Colors.green;
    default:
      return Colors.grey;
  }
}

String _muscleGroupLabel(String group) {
  return group
      .replaceAll('_', ' ')
      .split(' ')
      .map((w) => w.isNotEmpty ? '${w[0].toUpperCase()}${w.substring(1)}' : '')
      .join(' ');
}

IconData _equipmentIcon(String equipment) {
  switch (equipment.toLowerCase()) {
    case 'barbell':
      return Icons.fitness_center;
    case 'dumbbell':
      return Icons.fitness_center;
    case 'machine':
      return Icons.precision_manufacturing;
    case 'cable':
      return Icons.linear_scale;
    case 'bodyweight':
    case 'none':
      return Icons.accessibility_new;
    case 'kettlebell':
      return Icons.sports_martial_arts;
    case 'resistance_band':
    case 'band':
      return Icons.gesture;
    default:
      return Icons.fitness_center;
  }
}

IconData _typeIcon(String type) {
  switch (type.toLowerCase()) {
    case 'cardio':
      return Icons.directions_run;
    case 'strength':
      return Icons.fitness_center;
    case 'flexibility':
      return Icons.self_improvement;
    case 'hiit':
      return Icons.timer;
    default:
      return Icons.fitness_center;
  }
}

Color _difficultyColor(String difficulty) {
  switch (difficulty.toLowerCase()) {
    case 'beginner':
      return Colors.green;
    case 'intermediate':
      return Colors.orange;
    case 'advanced':
      return Colors.red;
    default:
      return AppTheme.textSecondary;
  }
}

// ---------------------------------------------------------------------------
// Main screen
// ---------------------------------------------------------------------------
class ExerciseLibraryScreen extends ConsumerStatefulWidget {
  const ExerciseLibraryScreen({super.key});

  @override
  ConsumerState<ExerciseLibraryScreen> createState() =>
      _ExerciseLibraryScreenState();
}

class _ExerciseLibraryScreenState extends ConsumerState<ExerciseLibraryScreen> {
  late final ExerciseService _exerciseService;

  List<ExerciseItem> _exercises = [];
  bool _loading = true;
  String? _error;

  // Filters
  String? _typeFilter;
  String? _equipmentFilter;

  // Expanded exercise id
  String? _expandedId;

  static const _types = ['cardio', 'strength', 'flexibility', 'hiit'];
  static const _equipment = [
    'none',
    'barbell',
    'dumbbell',
    'machine',
    'cable',
    'kettlebell',
    'resistance_band',
  ];

  @override
  void initState() {
    super.initState();
    _exerciseService = ExerciseService(ref.read(apiClientProvider));
    _loadExercises();
  }

  Future<void> _loadExercises() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final exercises = await _exerciseService.listExercises(
        type: _typeFilter,
        equipment: _equipmentFilter,
      );
      setState(() {
        _exercises = exercises;
        _loading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  void _applyTypeFilter(String? type) {
    setState(() => _typeFilter = type);
    _loadExercises();
  }

  void _applyEquipmentFilter(String? equipment) {
    setState(() => _equipmentFilter = equipment);
    _loadExercises();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      backgroundColor: AppTheme.background,
      appBar: AppBar(
        title: const Text('Exercise Library'),
      ),
      body: Column(
        children: [
          // Filter chips
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 12, 16, 4),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Type filter
                SingleChildScrollView(
                  scrollDirection: Axis.horizontal,
                  child: Row(
                    children: [
                      const Text(
                        'Type: ',
                        style: TextStyle(
                          color: AppTheme.textSecondary,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(width: 4),
                      _filterChip('All', _typeFilter == null,
                          () => _applyTypeFilter(null)),
                      ..._types.map((type) => _filterChip(
                            _capitalize(type),
                            _typeFilter == type,
                            () => _applyTypeFilter(type),
                            icon: _typeIcon(type),
                          )),
                    ],
                  ),
                ),
                const SizedBox(height: 6),
                // Equipment filter
                SingleChildScrollView(
                  scrollDirection: Axis.horizontal,
                  child: Row(
                    children: [
                      const Text(
                        'Equipment: ',
                        style: TextStyle(
                          color: AppTheme.textSecondary,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(width: 4),
                      _filterChip('All', _equipmentFilter == null,
                          () => _applyEquipmentFilter(null)),
                      ..._equipment.map((eq) => _filterChip(
                            _equipmentLabel(eq),
                            _equipmentFilter == eq,
                            () => _applyEquipmentFilter(eq),
                          )),
                    ],
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 4),
          // Exercise list
          Expanded(
            child: _loading
                ? const LoadingIndicator()
                : _error != null
                    ? Center(
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            const Icon(Icons.error_outline,
                                size: 48, color: AppTheme.error),
                            const SizedBox(height: 16),
                            Text(l10n.errorLoadingData,
                                style: const TextStyle(
                                    color: AppTheme.textSecondary)),
                            const SizedBox(height: 16),
                            ElevatedButton(
                                onPressed: _loadExercises,
                                child: Text(l10n.retry)),
                          ],
                        ),
                      )
                    : _exercises.isEmpty
                        ? Center(
                            child: EmptyState(
                              icon: Icons.fitness_center,
                              message: l10n.noResults,
                            ),
                          )
                        : RefreshIndicator(
                            color: AppTheme.primary,
                            onRefresh: _loadExercises,
                            child: ListView.builder(
                              padding: const EdgeInsets.symmetric(
                                  horizontal: 16, vertical: 8),
                              itemCount: _exercises.length,
                              itemBuilder: (context, index) =>
                                  _buildExerciseCard(_exercises[index]),
                            ),
                          ),
          ),
        ],
      ),
    );
  }

  Widget _filterChip(String label, bool selected, VoidCallback onTap,
      {IconData? icon}) {
    return Padding(
      padding: const EdgeInsets.only(right: 6),
      child: ChoiceChip(
        avatar: icon != null && !selected
            ? Icon(icon, size: 16, color: AppTheme.textSecondary)
            : icon != null && selected
                ? Icon(icon, size: 16, color: AppTheme.primary)
                : null,
        label: Text(label, style: const TextStyle(fontSize: 12)),
        selected: selected,
        selectedColor: AppTheme.primary.withAlpha(40),
        backgroundColor: AppTheme.surfaceLight,
        labelStyle: TextStyle(
          color: selected ? AppTheme.primary : AppTheme.textSecondary,
          fontWeight: selected ? FontWeight.w600 : FontWeight.normal,
        ),
        side: BorderSide(
          color: selected ? AppTheme.primary : Colors.transparent,
          width: 1,
        ),
        onSelected: (_) => onTap(),
        visualDensity: VisualDensity.compact,
      ),
    );
  }

  Widget _buildExerciseCard(ExerciseItem exercise) {
    final isExpanded = _expandedId == exercise.id;

    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: InkWell(
        borderRadius: BorderRadius.circular(12),
        onTap: () {
          setState(() {
            _expandedId = isExpanded ? null : exercise.id;
          });
        },
        child: Padding(
          padding: const EdgeInsets.all(14),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Header row
              Row(
                children: [
                  // Type icon
                  Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: AppTheme.primary.withAlpha(25),
                      borderRadius: BorderRadius.circular(10),
                    ),
                    child: Icon(
                      _typeIcon(exercise.exerciseType),
                      color: AppTheme.primary,
                      size: 20,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          exercise.name,
                          style: const TextStyle(
                            color: AppTheme.textPrimary,
                            fontSize: 15,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                        const SizedBox(height: 2),
                        Row(
                          children: [
                            Icon(
                              _equipmentIcon(exercise.equipment),
                              size: 13,
                              color: AppTheme.textSecondary,
                            ),
                            const SizedBox(width: 4),
                            Text(
                              _equipmentLabel(exercise.equipment),
                              style: const TextStyle(
                                color: AppTheme.textSecondary,
                                fontSize: 12,
                              ),
                            ),
                            const SizedBox(width: 10),
                            if (exercise.caloriesPerMinute > 0) ...[
                              const Icon(Icons.local_fire_department,
                                  size: 13, color: AppTheme.warning),
                              const SizedBox(width: 2),
                              Text(
                                '${exercise.caloriesPerMinute.toStringAsFixed(1)} cal/min',
                                style: const TextStyle(
                                  color: AppTheme.textSecondary,
                                  fontSize: 12,
                                ),
                              ),
                            ],
                          ],
                        ),
                      ],
                    ),
                  ),
                  // Difficulty chip
                  Container(
                    padding:
                        const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: _difficultyColor(exercise.difficulty).withAlpha(30),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Text(
                      _capitalize(exercise.difficulty),
                      style: TextStyle(
                        color: _difficultyColor(exercise.difficulty),
                        fontSize: 11,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                  const SizedBox(width: 4),
                  Icon(
                    isExpanded
                        ? Icons.keyboard_arrow_up
                        : Icons.keyboard_arrow_down,
                    color: AppTheme.textSecondary,
                    size: 20,
                  ),
                ],
              ),
              // Muscle group chips
              if (exercise.muscleGroups.isNotEmpty) ...[
                const SizedBox(height: 10),
                Wrap(
                  spacing: 6,
                  runSpacing: 4,
                  children: exercise.muscleGroups.map((m) {
                    return Container(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 8, vertical: 3),
                      decoration: BoxDecoration(
                        color: _muscleGroupColor(m).withAlpha(30),
                        borderRadius: BorderRadius.circular(6),
                      ),
                      child: Text(
                        _muscleGroupLabel(m),
                        style: TextStyle(
                          color: _muscleGroupColor(m),
                          fontSize: 11,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    );
                  }).toList(),
                ),
              ],
              // Expanded content
              if (isExpanded) ...[
                const SizedBox(height: 12),
                const Divider(height: 1),
                const SizedBox(height: 12),
                if (exercise.description != null &&
                    exercise.description!.isNotEmpty) ...[
                  Text(
                    exercise.description!,
                    style: const TextStyle(
                      color: AppTheme.textSecondary,
                      fontSize: 13,
                      height: 1.5,
                    ),
                  ),
                  const SizedBox(height: 12),
                ],
                if (exercise.instructions != null &&
                    exercise.instructions!.isNotEmpty) ...[
                  const Text(
                    'Instructions',
                    style: TextStyle(
                      color: AppTheme.textPrimary,
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 6),
                  Container(
                    width: double.infinity,
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: AppTheme.surfaceLight,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Text(
                      exercise.instructions!,
                      style: const TextStyle(
                        color: AppTheme.textSecondary,
                        fontSize: 13,
                        height: 1.6,
                      ),
                    ),
                  ),
                ],
                if (exercise.instructions == null &&
                    exercise.description == null)
                  const Text(
                    'No additional information available.',
                    style: TextStyle(
                      color: AppTheme.textSecondary,
                      fontSize: 13,
                      fontStyle: FontStyle.italic,
                    ),
                  ),
              ],
            ],
          ),
        ),
      ),
    );
  }

  String _capitalize(String s) {
    if (s.isEmpty) return s;
    return '${s[0].toUpperCase()}${s.substring(1)}';
  }

  String _equipmentLabel(String equipment) {
    return equipment
        .replaceAll('_', ' ')
        .split(' ')
        .map((w) => w.isNotEmpty ? '${w[0].toUpperCase()}${w.substring(1)}' : '')
        .join(' ');
  }
}
