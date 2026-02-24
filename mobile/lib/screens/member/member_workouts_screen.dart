import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../models/workout.dart';
import '../../providers/auth_provider.dart';
import '../../services/member_service.dart';
import '../../services/workout_service.dart';
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

// ---------------------------------------------------------------------------
// BodyMapPainter – simple front-view human silhouette
// ---------------------------------------------------------------------------
class BodyMapPainter extends CustomPainter {
  final Set<String> activeMuscles;

  BodyMapPainter({required this.activeMuscles});

  // Muscle group approximate positions on a 200x300 canvas
  static const Map<String, Offset> _positions = {
    'chest': Offset(100, 90),
    'back': Offset(100, 115),
    'shoulders': Offset(60, 72),
    'traps': Offset(100, 65),
    'biceps': Offset(48, 110),
    'triceps': Offset(152, 110),
    'forearms': Offset(40, 145),
    'core': Offset(100, 130),
    'obliques': Offset(78, 130),
    'lats': Offset(122, 105),
    'lower_back': Offset(100, 155),
    'quadriceps': Offset(82, 200),
    'hamstrings': Offset(118, 200),
    'glutes': Offset(100, 175),
    'calves': Offset(85, 255),
    'hip_flexors': Offset(100, 185),
  };

  // Duplicate shoulder on right side
  static const Map<String, Offset> _mirrorPositions = {
    'shoulders': Offset(140, 72),
    'biceps': Offset(152, 110),
    'triceps': Offset(48, 110),
    'forearms': Offset(160, 145),
    'quadriceps': Offset(118, 200),
    'hamstrings': Offset(82, 200),
    'calves': Offset(115, 255),
    'obliques': Offset(122, 130),
  };

  @override
  void paint(Canvas canvas, Size size) {
    final sx = size.width / 200;
    final sy = size.height / 300;
    canvas.save();
    canvas.scale(sx, sy);

    final bodyPaint = Paint()
      ..color = const Color(0xFF2D2D44)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2;

    final bodyFill = Paint()
      ..color = const Color(0xFF1A1A2E)
      ..style = PaintingStyle.fill;

    // Head
    canvas.drawCircle(const Offset(100, 35), 18, bodyFill);
    canvas.drawCircle(const Offset(100, 35), 18, bodyPaint);

    // Neck
    canvas.drawRect(const Rect.fromLTWH(93, 53, 14, 12), bodyFill);
    canvas.drawRect(const Rect.fromLTWH(93, 53, 14, 12), bodyPaint);

    // Torso
    final torsoPath = Path()
      ..moveTo(65, 65)
      ..lineTo(135, 65)
      ..lineTo(130, 170)
      ..lineTo(70, 170)
      ..close();
    canvas.drawPath(torsoPath, bodyFill);
    canvas.drawPath(torsoPath, bodyPaint);

    // Left arm
    final leftArm = Path()
      ..moveTo(65, 65)
      ..lineTo(50, 70)
      ..lineTo(35, 150)
      ..lineTo(48, 155)
      ..lineTo(58, 100)
      ..lineTo(65, 85);
    canvas.drawPath(leftArm, bodyFill);
    canvas.drawPath(leftArm, bodyPaint);

    // Right arm
    final rightArm = Path()
      ..moveTo(135, 65)
      ..lineTo(150, 70)
      ..lineTo(165, 150)
      ..lineTo(152, 155)
      ..lineTo(142, 100)
      ..lineTo(135, 85);
    canvas.drawPath(rightArm, bodyFill);
    canvas.drawPath(rightArm, bodyPaint);

    // Left leg
    final leftLeg = Path()
      ..moveTo(70, 170)
      ..lineTo(68, 270)
      ..lineTo(82, 275)
      ..lineTo(95, 170);
    canvas.drawPath(leftLeg, bodyFill);
    canvas.drawPath(leftLeg, bodyPaint);

    // Right leg
    final rightLeg = Path()
      ..moveTo(105, 170)
      ..lineTo(118, 275)
      ..lineTo(132, 270)
      ..lineTo(130, 170);
    canvas.drawPath(rightLeg, bodyFill);
    canvas.drawPath(rightLeg, bodyPaint);

    // Draw muscle group dots
    for (final entry in _positions.entries) {
      final group = entry.key;
      final pos = entry.value;
      final isActive = activeMuscles.contains(group);

      final dotPaint = Paint()
        ..color = isActive
            ? _muscleGroupColor(group).withAlpha(200)
            : const Color(0xFF3D3D55)
        ..style = PaintingStyle.fill;

      canvas.drawCircle(pos, isActive ? 8 : 5, dotPaint);

      if (isActive) {
        // Glow ring
        final glowPaint = Paint()
          ..color = _muscleGroupColor(group).withAlpha(60)
          ..style = PaintingStyle.fill;
        canvas.drawCircle(pos, 12, glowPaint);
        // Re-draw inner dot on top of glow
        canvas.drawCircle(pos, 8, dotPaint);
      }
    }

    // Draw mirrored dots
    for (final entry in _mirrorPositions.entries) {
      final group = entry.key;
      final pos = entry.value;
      final isActive = activeMuscles.contains(group);

      final dotPaint = Paint()
        ..color = isActive
            ? _muscleGroupColor(group).withAlpha(200)
            : const Color(0xFF3D3D55)
        ..style = PaintingStyle.fill;

      canvas.drawCircle(pos, isActive ? 8 : 5, dotPaint);

      if (isActive) {
        final glowPaint = Paint()
          ..color = _muscleGroupColor(group).withAlpha(60)
          ..style = PaintingStyle.fill;
        canvas.drawCircle(pos, 12, glowPaint);
        canvas.drawCircle(pos, 8, dotPaint);
      }
    }

    canvas.restore();
  }

  @override
  bool shouldRepaint(covariant BodyMapPainter oldDelegate) =>
      oldDelegate.activeMuscles != activeMuscles;
}

// ---------------------------------------------------------------------------
// Main screen
// ---------------------------------------------------------------------------
class MemberWorkoutsScreen extends ConsumerStatefulWidget {
  const MemberWorkoutsScreen({super.key});

  @override
  ConsumerState<MemberWorkoutsScreen> createState() =>
      _MemberWorkoutsScreenState();
}

class _MemberWorkoutsScreenState extends ConsumerState<MemberWorkoutsScreen>
    with SingleTickerProviderStateMixin {
  late final TabController _tabController;
  late final MemberService _memberService;
  late final WorkoutService _workoutService;

  // My Plan state
  WorkoutPlan? _plan;
  bool _planLoading = true;
  String? _planError;

  // Browse state
  List<WorkoutTemplate> _templates = [];
  bool _browseLoading = false;
  String? _browseError;
  String? _goalFilter;
  String? _difficultyFilter;

  // Find Plan (questionnaire) state
  String? _questionGoal;
  String? _questionExperience;
  String? _questionDays;
  int _stepIndex = 0;
  WorkoutPlan? _recommendedPlan;
  bool _findLoading = false;

  String? get _orgId => ref.read(authProvider).user?.orgId;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _tabController.addListener(() {
      if (_tabController.index == 1 && _templates.isEmpty && !_browseLoading) {
        _loadBrowse();
      }
    });
    _memberService = MemberService(ref.read(apiClientProvider));
    _workoutService = WorkoutService(ref.read(apiClientProvider));
    _loadPlan();
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  // ---- data loading -------------------------------------------------------

  Future<void> _loadPlan() async {
    setState(() {
      _planLoading = true;
      _planError = null;
    });
    try {
      final plan =
          await _memberService.getWorkoutPlan(organizationId: _orgId ?? '');
      setState(() {
        _plan = plan;
        _planLoading = false;
      });
    } catch (_) {
      // No plan assigned — show empty state, not error
      setState(() {
        _plan = null;
        _planLoading = false;
      });
    }
  }

  Future<void> _loadBrowse() async {
    setState(() {
      _browseLoading = true;
      _browseError = null;
    });
    try {
      final templates = await _workoutService.browseTemplates(
        _orgId ?? '',
        goal: _goalFilter,
        difficulty: _difficultyFilter,
      );
      setState(() {
        _templates = templates;
        _browseLoading = false;
      });
    } catch (_) {
      // Independent users see empty browse — show empty state, not error
      setState(() {
        _templates = [];
        _browseLoading = false;
      });
    }
  }

  Future<void> _choosePlan(String templateId) async {
    try {
      final plan = await _workoutService.selfAssignPlan(_orgId ?? '', templateId);
      setState(() {
        _plan = plan;
      });
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Plan assigned successfully!')),
        );
        _tabController.animateTo(0);
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to assign plan: $e')),
        );
      }
    }
  }

  Future<void> _submitQuestionnaire() async {
    if (_questionGoal == null ||
        _questionExperience == null ||
        _questionDays == null) return;
    setState(() => _findLoading = true);
    try {
      final plan = await _workoutService.recommendPlan(
        _orgId ?? '',
        _questionGoal!,
        _questionExperience!,
        _questionDays!,
      );
      setState(() {
        _recommendedPlan = plan;
        _findLoading = false;
      });
    } catch (e) {
      setState(() => _findLoading = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to get recommendation: $e')),
        );
      }
    }
  }

  Future<void> _acceptRecommendation() async {
    if (_recommendedPlan == null) return;
    try {
      if (_recommendedPlan!.template != null) {
        await _workoutService.selfAssignPlan(
            _orgId ?? '', _recommendedPlan!.template!.id);
      }
    } catch (_) {
      // Even if persistence fails, still accept locally
    }
    setState(() {
      _plan = _recommendedPlan;
      _recommendedPlan = null;
      _stepIndex = 0;
      _questionGoal = null;
      _questionExperience = null;
      _questionDays = null;
    });
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Plan accepted!')),
      );
      _tabController.animateTo(0);
    }
  }

  // ---- build --------------------------------------------------------------

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      backgroundColor: AppTheme.background,
      appBar: AppBar(
        title: Text(l10n.myWorkouts),
        bottom: TabBar(
          controller: _tabController,
          indicatorColor: AppTheme.primary,
          labelColor: AppTheme.primary,
          unselectedLabelColor: AppTheme.textSecondary,
          onTap: (index) {
            if (index == 1 && _templates.isEmpty && !_browseLoading) {
              _loadBrowse();
            }
          },
          tabs: const [
            Tab(text: 'My Plan'),
            Tab(text: 'Browse'),
            Tab(text: 'Find Plan'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildMyPlanTab(l10n),
          _buildBrowseTab(l10n),
          _buildFindPlanTab(l10n),
        ],
      ),
    );
  }

  // ===================== TAB 1 – My Plan ===================================

  Widget _buildMyPlanTab(AppLocalizations l10n) {
    if (_planLoading) return const LoadingIndicator();

    if (_planError != null) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline, size: 48, color: AppTheme.error),
            const SizedBox(height: 16),
            Text(l10n.errorLoadingData,
                style: const TextStyle(color: AppTheme.textSecondary)),
            const SizedBox(height: 16),
            ElevatedButton(onPressed: _loadPlan, child: Text(l10n.retry)),
          ],
        ),
      );
    }

    if (_plan == null || _plan!.template == null) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            EmptyState(
              icon: Icons.fitness_center,
              message: l10n.noAssignedPlan,
            ),
            const SizedBox(height: 24),
            ElevatedButton.icon(
              onPressed: () => _tabController.animateTo(1),
              icon: const Icon(Icons.search),
              label: const Text('Browse Plans'),
            ),
          ],
        ),
      );
    }

    final template = _plan!.template!;
    final allMuscles = <String>{};
    for (final ex in template.exercises) {
      allMuscles.addAll(ex.muscleGroups);
    }

    return RefreshIndicator(
      color: AppTheme.primary,
      onRefresh: _loadPlan,
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // Body map
          Center(
            child: SizedBox(
              width: 200,
              height: 200,
              child: CustomPaint(
                painter: BodyMapPainter(activeMuscles: allMuscles),
                size: const Size(200, 300),
              ),
            ),
          ),
          const SizedBox(height: 8),
          // Legend row
          if (allMuscles.isNotEmpty)
            Wrap(
              spacing: 6,
              runSpacing: 4,
              alignment: WrapAlignment.center,
              children: allMuscles.map((m) {
                return Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                  decoration: BoxDecoration(
                    color: _muscleGroupColor(m).withAlpha(40),
                    borderRadius: BorderRadius.circular(6),
                  ),
                  child: Text(
                    _muscleGroupLabel(m),
                    style: TextStyle(
                      color: _muscleGroupColor(m),
                      fontSize: 11,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                );
              }).toList(),
            ),
          const SizedBox(height: 16),

          // Plan name card
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          template.name,
                          style: const TextStyle(
                            fontSize: 18,
                            fontWeight: FontWeight.bold,
                            color: AppTheme.textPrimary,
                          ),
                        ),
                      ),
                      if (template.difficulty != null)
                        _difficultyChip(template.difficulty!),
                    ],
                  ),
                  const SizedBox(height: 6),
                  Row(
                    children: [
                      if (template.durationMinutes != null) ...[
                        const Icon(Icons.timer_outlined,
                            size: 14, color: AppTheme.textSecondary),
                        const SizedBox(width: 4),
                        Text('${template.durationMinutes} min',
                            style: const TextStyle(
                                color: AppTheme.textSecondary, fontSize: 13)),
                        const SizedBox(width: 12),
                      ],
                      if (_plan!.assignmentMethod != null)
                        _assignmentBadge(_plan!.assignmentMethod!),
                    ],
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 12),

          // Exercises
          ...template.exercises.asMap().entries.map((entry) {
            final idx = entry.key + 1;
            final exercise = entry.value;
            return _buildExerciseTile(idx, exercise);
          }),
        ],
      ),
    );
  }

  Widget _buildExerciseTile(int index, Exercise exercise) {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  width: 30,
                  height: 30,
                  decoration: BoxDecoration(
                    color: AppTheme.primary.withAlpha(30),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Center(
                    child: Text(
                      '$index',
                      style: const TextStyle(
                        color: AppTheme.primary,
                        fontWeight: FontWeight.bold,
                        fontSize: 14,
                      ),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Text(
                    exercise.name,
                    style: const TextStyle(
                      color: AppTheme.textPrimary,
                      fontSize: 15,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                const SizedBox(width: 42), // align with name
                Text(
                  '${exercise.sets} sets x ${exercise.reps} reps',
                  style: const TextStyle(
                      color: AppTheme.textSecondary, fontSize: 13),
                ),
                const SizedBox(width: 12),
                if (exercise.restSeconds > 0)
                  Row(
                    children: [
                      const Icon(Icons.hourglass_bottom,
                          size: 13, color: AppTheme.textSecondary),
                      const SizedBox(width: 2),
                      Text(
                        '${exercise.restSeconds}s rest',
                        style: const TextStyle(
                            color: AppTheme.textSecondary, fontSize: 13),
                      ),
                    ],
                  ),
              ],
            ),
            if (exercise.muscleGroups.isNotEmpty) ...[
              const SizedBox(height: 8),
              Padding(
                padding: const EdgeInsets.only(left: 42),
                child: Wrap(
                  spacing: 6,
                  runSpacing: 4,
                  children: exercise.muscleGroups.map((m) {
                    return Container(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 8, vertical: 2),
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
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _difficultyChip(String difficulty) {
    Color chipColor;
    switch (difficulty.toLowerCase()) {
      case 'beginner':
        chipColor = Colors.green;
        break;
      case 'intermediate':
        chipColor = Colors.orange;
        break;
      case 'advanced':
        chipColor = Colors.red;
        break;
      default:
        chipColor = AppTheme.textSecondary;
    }
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: chipColor.withAlpha(30),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(
        difficulty,
        style:
            TextStyle(color: chipColor, fontSize: 12, fontWeight: FontWeight.w600),
      ),
    );
  }

  Widget _assignmentBadge(String method) {
    IconData icon;
    String label;
    switch (method.toLowerCase()) {
      case 'staff':
        icon = Icons.person;
        label = 'Staff Assigned';
        break;
      case 'self':
        icon = Icons.touch_app;
        label = 'Self Assigned';
        break;
      case 'questionnaire':
        icon = Icons.quiz;
        label = 'Via Questionnaire';
        break;
      default:
        icon = Icons.info_outline;
        label = method;
    }
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
      decoration: BoxDecoration(
        color: AppTheme.info.withAlpha(25),
        borderRadius: BorderRadius.circular(6),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 13, color: AppTheme.info),
          const SizedBox(width: 4),
          Text(label,
              style:
                  const TextStyle(color: AppTheme.info, fontSize: 12)),
        ],
      ),
    );
  }

  // ===================== TAB 2 – Browse ====================================

  Widget _buildBrowseTab(AppLocalizations l10n) {
    return Column(
      children: [
        // Filter chips
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 12, 16, 4),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Goal filter
              SingleChildScrollView(
                scrollDirection: Axis.horizontal,
                child: Row(
                  children: [
                    const Text('Goal: ',
                        style: TextStyle(
                            color: AppTheme.textSecondary, fontSize: 13)),
                    const SizedBox(width: 4),
                    _filterChip('All', _goalFilter == null,
                        () => _applyGoalFilter(null)),
                    _filterChip('Lose Weight', _goalFilter == 'lose_weight',
                        () => _applyGoalFilter('lose_weight')),
                    _filterChip('Build Muscle', _goalFilter == 'build_muscle',
                        () => _applyGoalFilter('build_muscle')),
                    _filterChip('Stay Fit', _goalFilter == 'stay_fit',
                        () => _applyGoalFilter('stay_fit')),
                  ],
                ),
              ),
              const SizedBox(height: 6),
              // Difficulty filter
              SingleChildScrollView(
                scrollDirection: Axis.horizontal,
                child: Row(
                  children: [
                    const Text('Difficulty: ',
                        style: TextStyle(
                            color: AppTheme.textSecondary, fontSize: 13)),
                    const SizedBox(width: 4),
                    _filterChip('All', _difficultyFilter == null,
                        () => _applyDifficultyFilter(null)),
                    _filterChip('Beginner', _difficultyFilter == 'beginner',
                        () => _applyDifficultyFilter('beginner')),
                    _filterChip(
                        'Intermediate',
                        _difficultyFilter == 'intermediate',
                        () => _applyDifficultyFilter('intermediate')),
                    _filterChip('Advanced', _difficultyFilter == 'advanced',
                        () => _applyDifficultyFilter('advanced')),
                  ],
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: 4),
        // Template list
        Expanded(
          child: _browseLoading
              ? const LoadingIndicator()
              : _browseError != null
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
                              onPressed: _loadBrowse,
                              child: Text(l10n.retry)),
                        ],
                      ),
                    )
                  : _templates.isEmpty
                      ? Center(
                          child: EmptyState(
                            icon: Icons.fitness_center,
                            message: l10n.noResults,
                          ),
                        )
                      : RefreshIndicator(
                          color: AppTheme.primary,
                          onRefresh: _loadBrowse,
                          child: ListView.builder(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 16, vertical: 8),
                            itemCount: _templates.length,
                            itemBuilder: (context, index) =>
                                _buildTemplateCard(_templates[index]),
                          ),
                        ),
        ),
      ],
    );
  }

  void _applyGoalFilter(String? goal) {
    setState(() => _goalFilter = goal);
    _loadBrowse();
  }

  void _applyDifficultyFilter(String? difficulty) {
    setState(() => _difficultyFilter = difficulty);
    _loadBrowse();
  }

  Widget _filterChip(String label, bool selected, VoidCallback onTap) {
    return Padding(
      padding: const EdgeInsets.only(right: 6),
      child: ChoiceChip(
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

  Widget _buildTemplateCard(WorkoutTemplate template) {
    return Card(
      margin: const EdgeInsets.only(bottom: 10),
      child: Padding(
        padding: const EdgeInsets.all(14),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Text(
                    template.name,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                      color: AppTheme.textPrimary,
                    ),
                  ),
                ),
                if (template.difficulty != null)
                  _difficultyChip(template.difficulty!),
              ],
            ),
            if (template.description != null) ...[
              const SizedBox(height: 6),
              Text(
                template.description!,
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
                style: const TextStyle(
                    color: AppTheme.textSecondary, fontSize: 13),
              ),
            ],
            const SizedBox(height: 8),
            Row(
              children: [
                if (template.durationMinutes != null) ...[
                  const Icon(Icons.timer_outlined,
                      size: 14, color: AppTheme.textSecondary),
                  const SizedBox(width: 3),
                  Text('${template.durationMinutes} min',
                      style: const TextStyle(
                          color: AppTheme.textSecondary, fontSize: 12)),
                  const SizedBox(width: 12),
                ],
                const Icon(Icons.fitness_center,
                    size: 14, color: AppTheme.textSecondary),
                const SizedBox(width: 3),
                Text('${template.exercises.length} exercises',
                    style: const TextStyle(
                        color: AppTheme.textSecondary, fontSize: 12)),
                const Spacer(),
                SizedBox(
                  height: 32,
                  child: ElevatedButton(
                    onPressed: () => _choosePlan(template.id),
                    style: ElevatedButton.styleFrom(
                      padding: const EdgeInsets.symmetric(horizontal: 14),
                      textStyle: const TextStyle(fontSize: 13),
                    ),
                    child: const Text('Choose Plan'),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  // ===================== TAB 3 – Find Plan =================================

  Widget _buildFindPlanTab(AppLocalizations l10n) {
    if (_findLoading) {
      return const Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            CircularProgressIndicator(color: AppTheme.primary),
            SizedBox(height: 16),
            Text('Finding your perfect plan...',
                style: TextStyle(color: AppTheme.textSecondary)),
          ],
        ),
      );
    }

    if (_recommendedPlan != null) {
      return _buildRecommendedResult();
    }

    return Stepper(
      currentStep: _stepIndex,
      onStepContinue: () {
        if (_stepIndex == 0 && _questionGoal == null) return;
        if (_stepIndex == 1 && _questionExperience == null) return;
        if (_stepIndex < 2) {
          setState(() => _stepIndex++);
        } else {
          if (_questionDays != null) _submitQuestionnaire();
        }
      },
      onStepCancel: () {
        if (_stepIndex > 0) setState(() => _stepIndex--);
      },
      controlsBuilder: (context, details) {
        final isLast = _stepIndex == 2;
        return Padding(
          padding: const EdgeInsets.only(top: 12),
          child: Row(
            children: [
              ElevatedButton(
                onPressed: details.onStepContinue,
                child: Text(isLast ? 'Find My Plan' : 'Continue'),
              ),
              if (_stepIndex > 0) ...[
                const SizedBox(width: 12),
                TextButton(
                  onPressed: details.onStepCancel,
                  child: const Text('Back'),
                ),
              ],
            ],
          ),
        );
      },
      steps: [
        Step(
          title: const Text('What is your goal?',
              style: TextStyle(color: AppTheme.textPrimary)),
          isActive: _stepIndex >= 0,
          state: _stepIndex > 0 ? StepState.complete : StepState.indexed,
          content: Column(
            children: [
              _radioTile(
                icon: Icons.local_fire_department,
                iconColor: Colors.orange,
                title: 'Lose Weight',
                subtitle: 'Burn fat and get lean',
                value: 'lose_weight',
                groupValue: _questionGoal,
                onChanged: (v) => setState(() => _questionGoal = v),
              ),
              _radioTile(
                icon: Icons.fitness_center,
                iconColor: Colors.blue,
                title: 'Build Muscle',
                subtitle: 'Gain strength and size',
                value: 'build_muscle',
                groupValue: _questionGoal,
                onChanged: (v) => setState(() => _questionGoal = v),
              ),
              _radioTile(
                icon: Icons.favorite,
                iconColor: Colors.red,
                title: 'Stay Fit',
                subtitle: 'General fitness and health',
                value: 'stay_fit',
                groupValue: _questionGoal,
                onChanged: (v) => setState(() => _questionGoal = v),
              ),
            ],
          ),
        ),
        Step(
          title: const Text('Your experience level?',
              style: TextStyle(color: AppTheme.textPrimary)),
          isActive: _stepIndex >= 1,
          state: _stepIndex > 1 ? StepState.complete : StepState.indexed,
          content: Column(
            children: [
              _radioTile(
                icon: Icons.emoji_people,
                iconColor: Colors.green,
                title: 'Beginner',
                subtitle: 'New to working out',
                value: 'beginner',
                groupValue: _questionExperience,
                onChanged: (v) => setState(() => _questionExperience = v),
              ),
              _radioTile(
                icon: Icons.directions_run,
                iconColor: Colors.orange,
                title: 'Intermediate',
                subtitle: 'Some gym experience',
                value: 'intermediate',
                groupValue: _questionExperience,
                onChanged: (v) => setState(() => _questionExperience = v),
              ),
              _radioTile(
                icon: Icons.sports_gymnastics,
                iconColor: Colors.red,
                title: 'Advanced',
                subtitle: 'Experienced lifter',
                value: 'advanced',
                groupValue: _questionExperience,
                onChanged: (v) => setState(() => _questionExperience = v),
              ),
            ],
          ),
        ),
        Step(
          title: const Text('Days per week?',
              style: TextStyle(color: AppTheme.textPrimary)),
          isActive: _stepIndex >= 2,
          state: StepState.indexed,
          content: Column(
            children: [
              _radioTile(
                icon: Icons.calendar_today,
                iconColor: Colors.green,
                title: '2-3 days',
                subtitle: 'Light schedule',
                value: '2-3',
                groupValue: _questionDays,
                onChanged: (v) => setState(() => _questionDays = v),
              ),
              _radioTile(
                icon: Icons.calendar_view_week,
                iconColor: Colors.orange,
                title: '4-5 days',
                subtitle: 'Moderate schedule',
                value: '4-5',
                groupValue: _questionDays,
                onChanged: (v) => setState(() => _questionDays = v),
              ),
              _radioTile(
                icon: Icons.calendar_month,
                iconColor: Colors.red,
                title: '6-7 days',
                subtitle: 'Intensive schedule',
                value: '6-7',
                groupValue: _questionDays,
                onChanged: (v) => setState(() => _questionDays = v),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _radioTile({
    required IconData icon,
    required Color iconColor,
    required String title,
    required String subtitle,
    required String value,
    required String? groupValue,
    required ValueChanged<String?> onChanged,
  }) {
    final selected = value == groupValue;
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: InkWell(
        borderRadius: BorderRadius.circular(12),
        onTap: () => onChanged(value),
        child: Container(
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: selected
                ? AppTheme.primary.withAlpha(15)
                : AppTheme.surfaceLight,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: selected ? AppTheme.primary : Colors.transparent,
              width: 1.5,
            ),
          ),
          child: Row(
            children: [
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: iconColor.withAlpha(25),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Icon(icon, color: iconColor, size: 22),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(title,
                        style: const TextStyle(
                            color: AppTheme.textPrimary,
                            fontWeight: FontWeight.w600,
                            fontSize: 15)),
                    Text(subtitle,
                        style: const TextStyle(
                            color: AppTheme.textSecondary, fontSize: 12)),
                  ],
                ),
              ),
              Radio<String>(
                value: value,
                groupValue: groupValue,
                activeColor: AppTheme.primary,
                onChanged: onChanged,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildRecommendedResult() {
    final template = _recommendedPlan!.template;
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          const Icon(Icons.check_circle, color: AppTheme.primary, size: 56),
          const SizedBox(height: 12),
          const Text(
            'We found a plan for you!',
            style: TextStyle(
              color: AppTheme.textPrimary,
              fontSize: 20,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 16),
          if (template != null)
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(template.name,
                        style: const TextStyle(
                            fontSize: 18,
                            fontWeight: FontWeight.bold,
                            color: AppTheme.textPrimary)),
                    if (template.description != null) ...[
                      const SizedBox(height: 6),
                      Text(template.description!,
                          style: const TextStyle(
                              color: AppTheme.textSecondary, fontSize: 14)),
                    ],
                    const SizedBox(height: 8),
                    Row(
                      children: [
                        if (template.difficulty != null)
                          _difficultyChip(template.difficulty!),
                        if (template.difficulty != null &&
                            template.durationMinutes != null)
                          const SizedBox(width: 8),
                        if (template.durationMinutes != null)
                          Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              const Icon(Icons.timer_outlined,
                                  size: 14, color: AppTheme.textSecondary),
                              const SizedBox(width: 3),
                              Text('${template.durationMinutes} min',
                                  style: const TextStyle(
                                      color: AppTheme.textSecondary,
                                      fontSize: 13)),
                            ],
                          ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Text('${template.exercises.length} exercises',
                        style: const TextStyle(
                            color: AppTheme.textSecondary, fontSize: 13)),
                  ],
                ),
              ),
            ),
          const SizedBox(height: 20),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton.icon(
              onPressed: _acceptRecommendation,
              icon: const Icon(Icons.check),
              label: const Text('Accept Plan'),
            ),
          ),
          const SizedBox(height: 8),
          TextButton(
            onPressed: () {
              setState(() {
                _recommendedPlan = null;
                _stepIndex = 0;
                _questionGoal = null;
                _questionExperience = null;
                _questionDays = null;
              });
            },
            child: const Text('Try Again'),
          ),
        ],
      ),
    );
  }
}
