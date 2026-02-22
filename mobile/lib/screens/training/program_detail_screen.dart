import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../models/exercise.dart';
import '../../providers/auth_provider.dart';
import '../../services/exercise_service.dart';
import '../shared/widgets/loading_indicator.dart';

class ProgramDetailScreen extends ConsumerStatefulWidget {
  final String programId;

  const ProgramDetailScreen({super.key, required this.programId});

  @override
  ConsumerState<ProgramDetailScreen> createState() =>
      _ProgramDetailScreenState();
}

class _ProgramDetailScreenState extends ConsumerState<ProgramDetailScreen> {
  late final ExerciseService _exerciseService;

  WorkoutProgram? _program;
  UserProgramEnrollment? _enrollment;
  bool _loading = true;
  bool _enrolling = false;
  bool _completing = false;
  String? _error;
  int _selectedWeek = 1;

  @override
  void initState() {
    super.initState();
    _exerciseService = ExerciseService(ref.read(apiClientProvider));
    _loadData();
  }

  Future<void> _loadData() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final results = await Future.wait([
        _exerciseService.getProgram(widget.programId),
        _loadEnrollment(),
      ]);
      setState(() {
        _program = results[0] as WorkoutProgram;
        _loading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  Future<Object?> _loadEnrollment() async {
    try {
      final enrollment = await _exerciseService.getCurrentProgram();
      if (enrollment.programId == widget.programId) {
        setState(() {
          _enrollment = enrollment;
          _selectedWeek = enrollment.currentWeek;
        });
      }
      return enrollment;
    } catch (_) {
      return null;
    }
  }

  Future<void> _enroll() async {
    setState(() => _enrolling = true);
    try {
      final enrollment =
          await _exerciseService.enrollInProgram(widget.programId);
      setState(() {
        _enrollment = enrollment;
        _enrolling = false;
        _selectedWeek = 1;
      });
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Enrolled successfully!')),
        );
      }
    } catch (e) {
      setState(() => _enrolling = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to enroll: $e')),
        );
      }
    }
  }

  Future<void> _completeDay() async {
    setState(() => _completing = true);
    try {
      final enrollment = await _exerciseService.completeProgramDay();
      setState(() {
        _enrollment = enrollment;
        _completing = false;
      });
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Day completed! Great work!')),
        );
      }
    } catch (e) {
      setState(() => _completing = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to complete day: $e')),
        );
      }
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

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      backgroundColor: AppTheme.background,
      appBar: AppBar(
        title: Text(_program?.name ?? 'Program Details'),
      ),
      body: _loading
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
                          onPressed: _loadData, child: Text(l10n.retry)),
                    ],
                  ),
                )
              : _buildContent(l10n),
    );
  }

  Widget _buildContent(AppLocalizations l10n) {
    final program = _program!;
    final isEnrolled = _enrollment != null;
    final totalDays = program.durationWeeks * 7;
    final completedDays = isEnrolled
        ? ((_enrollment!.currentWeek - 1) * 7) + _enrollment!.currentDay - 1
        : 0;
    final progress =
        totalDays > 0 ? (completedDays / totalDays).clamp(0.0, 1.0) : 0.0;

    return Column(
      children: [
        Expanded(
          child: RefreshIndicator(
            color: AppTheme.primary,
            onRefresh: _loadData,
            child: ListView(
              padding: const EdgeInsets.all(16),
              children: [
                // Program info card
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          program.name,
                          style: const TextStyle(
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                            color: AppTheme.textPrimary,
                          ),
                        ),
                        if (program.description != null &&
                            program.description!.isNotEmpty) ...[
                          const SizedBox(height: 8),
                          Text(
                            program.description!,
                            style: const TextStyle(
                              color: AppTheme.textSecondary,
                              fontSize: 14,
                              height: 1.5,
                            ),
                          ),
                        ],
                        const SizedBox(height: 14),
                        // Info chips row
                        Wrap(
                          spacing: 10,
                          runSpacing: 8,
                          children: [
                            _infoChip(
                              Icons.calendar_month,
                              '${program.durationWeeks} weeks',
                            ),
                            _infoChip(
                              Icons.fitness_center,
                              _equipmentLabel(program.equipment),
                            ),
                            _difficultyChip(program.difficulty),
                            if (program.programType != null)
                              _infoChip(
                                Icons.category,
                                _capitalize(program.programType!),
                              ),
                            if (program.goal != null)
                              _infoChip(
                                Icons.track_changes,
                                _capitalize(
                                    program.goal!.replaceAll('_', ' ')),
                              ),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
                const SizedBox(height: 16),

                // Progress section (if enrolled)
                if (isEnrolled) ...[
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            children: [
                              const Text(
                                'Your Progress',
                                style: TextStyle(
                                  fontSize: 16,
                                  fontWeight: FontWeight.bold,
                                  color: AppTheme.textPrimary,
                                ),
                              ),
                              const Spacer(),
                              Text(
                                '${(progress * 100).toInt()}%',
                                style: const TextStyle(
                                  color: AppTheme.primary,
                                  fontSize: 16,
                                  fontWeight: FontWeight.bold,
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 12),
                          ClipRRect(
                            borderRadius: BorderRadius.circular(6),
                            child: LinearProgressIndicator(
                              value: progress,
                              backgroundColor: AppTheme.surfaceLight,
                              valueColor:
                                  const AlwaysStoppedAnimation<Color>(
                                      AppTheme.primary),
                              minHeight: 10,
                            ),
                          ),
                          const SizedBox(height: 10),
                          Text(
                            'Week ${_enrollment!.currentWeek}, Day ${_enrollment!.currentDay}',
                            style: const TextStyle(
                              color: AppTheme.textSecondary,
                              fontSize: 13,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),
                ],

                // Week selector
                const Text(
                  'Schedule',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: AppTheme.textPrimary,
                  ),
                ),
                const SizedBox(height: 12),
                SizedBox(
                  height: 44,
                  child: ListView.builder(
                    scrollDirection: Axis.horizontal,
                    itemCount: program.durationWeeks,
                    itemBuilder: (context, index) {
                      final week = index + 1;
                      final isSelected = week == _selectedWeek;
                      final isCurrentWeek =
                          isEnrolled && week == _enrollment!.currentWeek;

                      return Padding(
                        padding: const EdgeInsets.only(right: 8),
                        child: ChoiceChip(
                          label: Text(
                            'Week $week',
                            style: TextStyle(
                              fontSize: 13,
                              color: isSelected
                                  ? AppTheme.primary
                                  : AppTheme.textSecondary,
                              fontWeight: isSelected
                                  ? FontWeight.w600
                                  : FontWeight.normal,
                            ),
                          ),
                          selected: isSelected,
                          selectedColor: AppTheme.primary.withAlpha(40),
                          backgroundColor: AppTheme.surfaceLight,
                          side: BorderSide(
                            color: isCurrentWeek && !isSelected
                                ? AppTheme.primary.withAlpha(100)
                                : isSelected
                                    ? AppTheme.primary
                                    : Colors.transparent,
                            width: 1,
                          ),
                          onSelected: (_) {
                            setState(() => _selectedWeek = week);
                          },
                        ),
                      );
                    },
                  ),
                ),
                const SizedBox(height: 12),

                // Days for selected week
                _buildWeekDays(program),
              ],
            ),
          ),
        ),

        // Bottom action bar
        _buildBottomBar(isEnrolled),
      ],
    );
  }

  Widget _buildWeekDays(WorkoutProgram program) {
    final days = program.days
            ?.where((d) => d.weekNumber == _selectedWeek)
            .toList() ??
        [];

    if (days.isEmpty) {
      // Generate default 7-day layout if no specific days data
      return Column(
        children: List.generate(7, (index) {
          final dayNum = index + 1;
          final isCurrentDay = _enrollment != null &&
              _enrollment!.currentWeek == _selectedWeek &&
              _enrollment!.currentDay == dayNum;
          final isPastDay = _enrollment != null &&
              (_enrollment!.currentWeek > _selectedWeek ||
                  (_enrollment!.currentWeek == _selectedWeek &&
                      _enrollment!.currentDay > dayNum));

          return Card(
            margin: const EdgeInsets.only(bottom: 6),
            child: Padding(
              padding: const EdgeInsets.all(12),
              child: Row(
                children: [
                  Container(
                    width: 36,
                    height: 36,
                    decoration: BoxDecoration(
                      color: isCurrentDay
                          ? AppTheme.primary.withAlpha(30)
                          : isPastDay
                              ? AppTheme.primary.withAlpha(15)
                              : AppTheme.surfaceLight,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Center(
                      child: isPastDay
                          ? const Icon(Icons.check,
                              color: AppTheme.primary, size: 18)
                          : Text(
                              '$dayNum',
                              style: TextStyle(
                                color: isCurrentDay
                                    ? AppTheme.primary
                                    : AppTheme.textSecondary,
                                fontWeight: FontWeight.bold,
                                fontSize: 14,
                              ),
                            ),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Text(
                    _dayName(dayNum),
                    style: TextStyle(
                      color: isCurrentDay
                          ? AppTheme.textPrimary
                          : AppTheme.textSecondary,
                      fontSize: 14,
                      fontWeight:
                          isCurrentDay ? FontWeight.w600 : FontWeight.normal,
                    ),
                  ),
                  if (isCurrentDay) ...[
                    const SizedBox(width: 8),
                    Container(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 8, vertical: 2),
                      decoration: BoxDecoration(
                        color: AppTheme.primary.withAlpha(25),
                        borderRadius: BorderRadius.circular(6),
                      ),
                      child: const Text(
                        'TODAY',
                        style: TextStyle(
                          color: AppTheme.primary,
                          fontSize: 10,
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                    ),
                  ],
                ],
              ),
            ),
          );
        }),
      );
    }

    // Render actual day data
    return Column(
      children: days.map((day) {
        final isCurrentDay = _enrollment != null &&
            _enrollment!.currentWeek == _selectedWeek &&
            _enrollment!.currentDay == day.dayNumber;
        final isPastDay = _enrollment != null &&
            (_enrollment!.currentWeek > _selectedWeek ||
                (_enrollment!.currentWeek == _selectedWeek &&
                    _enrollment!.currentDay > day.dayNumber));

        return Card(
          margin: const EdgeInsets.only(bottom: 6),
          child: Padding(
            padding: const EdgeInsets.all(12),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Container(
                      width: 36,
                      height: 36,
                      decoration: BoxDecoration(
                        color: isCurrentDay
                            ? AppTheme.primary.withAlpha(30)
                            : isPastDay
                                ? AppTheme.primary.withAlpha(15)
                                : day.restDay
                                    ? AppTheme.warning.withAlpha(20)
                                    : AppTheme.surfaceLight,
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Center(
                        child: isPastDay
                            ? const Icon(Icons.check,
                                color: AppTheme.primary, size: 18)
                            : day.restDay
                                ? const Icon(Icons.bed,
                                    color: AppTheme.warning, size: 18)
                                : Text(
                                    '${day.dayNumber}',
                                    style: TextStyle(
                                      color: isCurrentDay
                                          ? AppTheme.primary
                                          : AppTheme.textSecondary,
                                      fontWeight: FontWeight.bold,
                                      fontSize: 14,
                                    ),
                                  ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            day.dayName ??
                                (day.restDay
                                    ? 'Rest Day'
                                    : 'Day ${day.dayNumber}'),
                            style: TextStyle(
                              color: isCurrentDay
                                  ? AppTheme.textPrimary
                                  : AppTheme.textSecondary,
                              fontSize: 14,
                              fontWeight: isCurrentDay
                                  ? FontWeight.w600
                                  : FontWeight.normal,
                            ),
                          ),
                          if (!day.restDay && day.exercises.isNotEmpty)
                            Text(
                              '${day.exercises.length} exercises',
                              style: const TextStyle(
                                color: AppTheme.textSecondary,
                                fontSize: 12,
                              ),
                            ),
                        ],
                      ),
                    ),
                    if (isCurrentDay)
                      Container(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 8, vertical: 2),
                        decoration: BoxDecoration(
                          color: AppTheme.primary.withAlpha(25),
                          borderRadius: BorderRadius.circular(6),
                        ),
                        child: const Text(
                          'TODAY',
                          style: TextStyle(
                            color: AppTheme.primary,
                            fontSize: 10,
                            fontWeight: FontWeight.w700,
                          ),
                        ),
                      ),
                    if (day.restDay && !isCurrentDay)
                      Container(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 8, vertical: 2),
                        decoration: BoxDecoration(
                          color: AppTheme.warning.withAlpha(25),
                          borderRadius: BorderRadius.circular(6),
                        ),
                        child: const Text(
                          'REST',
                          style: TextStyle(
                            color: AppTheme.warning,
                            fontSize: 10,
                            fontWeight: FontWeight.w700,
                          ),
                        ),
                      ),
                  ],
                ),
              ],
            ),
          ),
        );
      }).toList(),
    );
  }

  Widget _buildBottomBar(bool isEnrolled) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: const BoxDecoration(
        color: AppTheme.surface,
        border: Border(
          top: BorderSide(color: Color(0xFF2D2D44), width: 1),
        ),
      ),
      child: SafeArea(
        child: isEnrolled
            ? SizedBox(
                width: double.infinity,
                child: ElevatedButton.icon(
                  onPressed: _completing ? null : _completeDay,
                  icon: _completing
                      ? const SizedBox(
                          width: 18,
                          height: 18,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            color: Colors.white,
                          ),
                        )
                      : const Icon(Icons.check_circle_outline),
                  label: Text(
                      _completing ? 'Completing...' : 'Complete Today\'s Workout'),
                ),
              )
            : SizedBox(
                width: double.infinity,
                child: ElevatedButton.icon(
                  onPressed: _enrolling ? null : _enroll,
                  icon: _enrolling
                      ? const SizedBox(
                          width: 18,
                          height: 18,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            color: Colors.white,
                          ),
                        )
                      : const Icon(Icons.play_arrow),
                  label: Text(
                      _enrolling ? 'Enrolling...' : 'Enroll in Program'),
                ),
              ),
      ),
    );
  }

  Widget _infoChip(IconData icon, String label) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
      decoration: BoxDecoration(
        color: AppTheme.surfaceLight,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 14, color: AppTheme.textSecondary),
          const SizedBox(width: 5),
          Text(
            label,
            style: const TextStyle(
              color: AppTheme.textSecondary,
              fontSize: 12,
            ),
          ),
        ],
      ),
    );
  }

  Widget _difficultyChip(String difficulty) {
    final color = _difficultyColor(difficulty);
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
      decoration: BoxDecoration(
        color: color.withAlpha(25),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(
        _capitalize(difficulty),
        style: TextStyle(
          color: color,
          fontSize: 12,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }

  String _dayName(int dayNumber) {
    const days = [
      'Monday',
      'Tuesday',
      'Wednesday',
      'Thursday',
      'Friday',
      'Saturday',
      'Sunday',
    ];
    return days[(dayNumber - 1) % 7];
  }
}
