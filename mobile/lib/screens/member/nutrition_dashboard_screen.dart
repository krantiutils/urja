import 'dart:math' as math;

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';

import '../../config/theme.dart';
import '../../models/nutrition.dart';
import '../../models/water.dart';
import '../../models/weight.dart';
import '../../providers/auth_provider.dart';
import '../../services/member_service.dart';
import '../../services/nutrition_service.dart';
import '../../services/water_service.dart';
import '../shared/widgets/loading_indicator.dart';

// ─── Colors ──────────────────────────────────────────────────────────────────

const _kCardColor = Color(0xFF1E293B);
const _kCardBorder = Color(0xFF334155);
const _kGreen = Color(0xFF4ADE80);
const _kBlue = Color(0xFF60A5FA);
const _kOrange = Color(0xFFFB923C);
const _kPink = Color(0xFFF472B6);
const _kPurple = Color(0xFFA78BFA);
const _kYellow = Color(0xFFFBBF24);
const _kCyan = Color(0xFF22D3EE);
const _kTextMuted = Color(0xFF94A3B8);
const _kTextDim = Color(0xFF64748B);

// ─── Screen ──────────────────────────────────────────────────────────────────

class NutritionDashboardScreen extends ConsumerStatefulWidget {
  const NutritionDashboardScreen({super.key});

  @override
  ConsumerState<NutritionDashboardScreen> createState() =>
      _NutritionDashboardScreenState();
}

class _NutritionDashboardScreenState
    extends ConsumerState<NutritionDashboardScreen>
    with TickerProviderStateMixin {
  late final NutritionService _nutritionService;
  late final WaterService _waterService;
  late final MemberService _memberService;

  DateTime _selectedDate = DateTime.now();
  DailySummary? _dailySummary;
  NutritionGoal? _goal;
  List<WeeklySummaryDay> _weeklySummary = [];
  NutritionStreak? _streak;
  WaterDailySummary? _waterSummary;
  WeightTrend? _weightTrend;
  bool _loading = true;
  String? _error;

  late AnimationController _ringAnimController;
  late Animation<double> _ringAnim;

  @override
  void initState() {
    super.initState();
    final apiClient = ref.read(apiClientProvider);
    _nutritionService = NutritionService(apiClient);
    _waterService = WaterService(apiClient);
    _memberService = MemberService(apiClient);

    _ringAnimController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 900),
    );
    _ringAnim = CurvedAnimation(
      parent: _ringAnimController,
      curve: Curves.easeOutCubic,
    );

    _loadData();
  }

  @override
  void dispose() {
    _ringAnimController.dispose();
    super.dispose();
  }

  String? get _orgId => ref.read(authProvider).user?.orgId;

  String _fmtDate(DateTime dt) =>
      '${dt.year}-${dt.month.toString().padLeft(2, '0')}-${dt.day.toString().padLeft(2, '0')}';

  bool get _isToday {
    final now = DateTime.now();
    return _selectedDate.year == now.year &&
        _selectedDate.month == now.month &&
        _selectedDate.day == now.day;
  }

  Future<void> _loadData() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final orgId = _orgId ?? '';
      final dateStr = _fmtDate(_selectedDate);
      final weekStart =
          _selectedDate.subtract(Duration(days: _selectedDate.weekday - 1));

      final results = await Future.wait([
        _nutritionService.getDailySummary(orgId, dateStr),
        _nutritionService.getNutritionGoal(orgId).catchError((_) =>
            NutritionGoal(
              id: '',
              userId: '',
              organizationId: '',
              calorieGoal: 2000,
              proteinGoalG: 150,
              carbsGoalG: 250,
              fatGoalG: 65,
              weightKg: 0,
              heightCm: 0,
              age: 0,
              gender: '',
              activityLevel: '',
              goalType: '',
            )),
        _nutritionService.getWeeklySummary(orgId, _fmtDate(weekStart)),
        _nutritionService.getNutritionStreak().catchError((_) =>
            NutritionStreak(
              userId: '',
              currentStreak: 0,
              longestStreak: 0,
              updatedAt: '',
            )),
        _waterService.getDailySummary(dateStr).catchError((_) =>
            WaterDailySummary(
              date: dateStr,
              totalMl: 0,
              goalMl: 2500,
              entries: [],
            )),
        _memberService.getWeightTrend(days: 10).catchError((_) => WeightTrend(
              entries: [],
              summary: WeightTrendSummary(
                currentKg: 0,
                startKg: 0,
                changeKg: 0,
                bmi: 0,
                entriesCount: 0,
              ),
            )),
      ]);

      if (!mounted) return;
      setState(() {
        _dailySummary = results[0] as DailySummary;
        _goal = results[1] as NutritionGoal;
        _weeklySummary = results[2] as List<WeeklySummaryDay>;
        _streak = results[3] as NutritionStreak;
        _waterSummary = results[4] as WaterDailySummary;
        _weightTrend = results[5] as WeightTrend;
        _loading = false;
      });
      _ringAnimController.forward(from: 0);
    } catch (e) {
      if (!mounted) return;
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  void _changeDate(int days) {
    setState(() => _selectedDate = _selectedDate.add(Duration(days: days)));
    _loadData();
  }

  Future<void> _pickDate() async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: DateTime(2020),
      lastDate: DateTime.now(),
      builder: (ctx, child) => Theme(
        data: Theme.of(ctx).copyWith(
          colorScheme: const ColorScheme.dark(
            primary: _kGreen,
            onPrimary: Colors.white,
            surface: _kCardColor,
            onSurface: Colors.white,
          ),
        ),
        child: child!,
      ),
    );
    if (picked != null && mounted) {
      setState(() => _selectedDate = picked);
      _loadData();
    }
  }

  Future<void> _addWater(int ml) async {
    try {
      await _waterService.logWater(
        amountMl: ml,
        date: _fmtDate(_selectedDate),
      );
      _loadData();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to log water: $e'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    }
  }

  // ─── Build ─────────────────────────────────────────────────────────────────

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppTheme.background,
      body: _loading
          ? const LoadingIndicator()
          : _error != null
              ? _buildError()
              : _buildContent(),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {
          final loc = GoRouterState.of(context).uri.toString();
          if (loc.startsWith('/tracker')) {
            context.go('/tracker/nutrition/log');
          } else {
            context.go('/member/nutrition/log');
          }
        },
        backgroundColor: _kGreen,
        icon: const Icon(Icons.add, color: Colors.black),
        label: const Text(
          'Log Food',
          style: TextStyle(
            color: Colors.black,
            fontWeight: FontWeight.w600,
          ),
        ),
      ),
    );
  }

  Widget _buildError() {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline, color: AppTheme.error, size: 48),
            const SizedBox(height: 16),
            Text(
              'Something went wrong',
              style: const TextStyle(
                color: Colors.white,
                fontSize: 18,
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              _error ?? '',
              style: const TextStyle(color: _kTextMuted, fontSize: 14),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            ElevatedButton.icon(
              onPressed: _loadData,
              icon: const Icon(Icons.refresh),
              label: const Text('Retry'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildContent() {
    return RefreshIndicator(
      color: _kGreen,
      onRefresh: _loadData,
      child: CustomScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        slivers: [
          // Safe area top padding
          SliverToBoxAdapter(
            child: SizedBox(height: MediaQuery.of(context).padding.top + 8),
          ),
          // Date nav bar
          SliverToBoxAdapter(child: _buildDateNav()),
          // Calories card
          SliverToBoxAdapter(child: _buildCaloriesCard()),
          // Macros card
          SliverToBoxAdapter(child: _buildMacrosCard()),
          // Streak card
          SliverToBoxAdapter(child: _buildStreakCard()),
          // Weight card
          SliverToBoxAdapter(child: _buildWeightCard()),
          // Water card
          SliverToBoxAdapter(child: _buildWaterCard()),
          // Bottom padding for FAB
          const SliverToBoxAdapter(child: SizedBox(height: 100)),
        ],
      ),
    );
  }

  // ─── Date Navigation ───────────────────────────────────────────────────────

  Widget _buildDateNav() {
    final dayName = DateFormat('EEE').format(_selectedDate);
    final monthDay = DateFormat('MMM d').format(_selectedDate);
    final displayDate = _isToday ? 'Today' : '$dayName, $monthDay';

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Container(
        decoration: BoxDecoration(
          color: _kCardColor,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: _kCardBorder.withAlpha(80)),
        ),
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            IconButton(
              onPressed: () => _changeDate(-1),
              icon: const Icon(Icons.chevron_left, color: Colors.white),
              splashRadius: 20,
            ),
            GestureDetector(
              onTap: _pickDate,
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    displayDate,
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 17,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(width: 8),
                  const Icon(
                    Icons.calendar_today,
                    color: _kGreen,
                    size: 18,
                  ),
                ],
              ),
            ),
            IconButton(
              onPressed: _isToday ? null : () => _changeDate(1),
              icon: Icon(
                Icons.chevron_right,
                color: _isToday ? _kTextDim : Colors.white,
              ),
              splashRadius: 20,
            ),
          ],
        ),
      ),
    );
  }

  // ─── Calories Card ─────────────────────────────────────────────────────────

  Widget _buildCaloriesCard() {
    final consumed = _dailySummary?.totalCalories ?? 0;
    final goal = _goal?.calorieGoal ?? 2000;
    final remaining = goal - consumed;

    return _DashboardCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const _SectionTitle(icon: Icons.local_fire_department, title: 'Calories'),
          const SizedBox(height: 16),
          Row(
            children: [
              // Donut chart
              Expanded(
                flex: 4,
                child: AnimatedBuilder(
                  animation: _ringAnim,
                  builder: (context, _) {
                    return SizedBox(
                      height: 140,
                      child: CustomPaint(
                        painter: _CalorieRingPainter(
                          consumed: consumed,
                          goal: goal,
                          progress: _ringAnim.value,
                        ),
                        child: Center(
                          child: Column(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Text(
                                consumed.toInt().toString(),
                                style: const TextStyle(
                                  color: Colors.white,
                                  fontSize: 28,
                                  fontWeight: FontWeight.w700,
                                ),
                              ),
                              const Text(
                                'eaten',
                                style: TextStyle(
                                  color: _kTextMuted,
                                  fontSize: 12,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    );
                  },
                ),
              ),
              const SizedBox(width: 12),
              // Stats column
              Expanded(
                flex: 3,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _CalorieStat(
                      label: 'Goal',
                      value: goal.toInt().toString(),
                      color: _kTextMuted,
                    ),
                    const SizedBox(height: 12),
                    _CalorieStat(
                      label: 'Eaten',
                      value: consumed.toInt().toString(),
                      color: _kGreen,
                    ),
                    const SizedBox(height: 12),
                    _CalorieStat(
                      label: remaining >= 0 ? 'Remaining' : 'Over',
                      value: remaining.abs().toInt().toString(),
                      color: remaining >= 0 ? _kBlue : AppTheme.error,
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 20),
          // Weekly bar chart
          _buildWeeklyCalorieChart(),
        ],
      ),
    );
  }

  Widget _buildWeeklyCalorieChart() {
    final days = ['M', 'T', 'W', 'T', 'F', 'S', 'S'];
    final goal = _goal?.calorieGoal ?? 2000;

    // Ensure we have 7 days of data
    final weekData = List.generate(7, (i) {
      if (i < _weeklySummary.length) return _weeklySummary[i].totalCalories;
      return 0.0;
    });
    final maxVal = [goal, ...weekData].reduce(math.max);

    // Determine which bar represents the selected date
    final selectedDayIndex = _selectedDate.weekday - 1; // 0=Mon, 6=Sun

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text(
              'Weekly Overview',
              style: TextStyle(
                color: _kTextMuted,
                fontSize: 13,
                fontWeight: FontWeight.w500,
              ),
            ),
            Text(
              'Goal: ${goal.toInt()} cal',
              style: const TextStyle(
                color: _kTextDim,
                fontSize: 12,
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        SizedBox(
          height: 80,
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: List.generate(7, (i) {
              final val = weekData[i];
              final fraction = maxVal > 0 ? (val / maxVal) : 0.0;
              final isSelected = i == selectedDayIndex;
              final isOverGoal = val > goal;

              return Expanded(
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 3),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.end,
                    children: [
                      Flexible(
                        child: AnimatedContainer(
                          duration: const Duration(milliseconds: 400),
                          curve: Curves.easeOutCubic,
                          width: double.infinity,
                          height: math.max(4, 60 * fraction),
                          decoration: BoxDecoration(
                            borderRadius: BorderRadius.circular(4),
                            gradient: LinearGradient(
                              begin: Alignment.bottomCenter,
                              end: Alignment.topCenter,
                              colors: isOverGoal
                                  ? [
                                      AppTheme.error.withAlpha(200),
                                      AppTheme.error,
                                    ]
                                  : isSelected
                                      ? [
                                          _kGreen.withAlpha(180),
                                          _kGreen,
                                        ]
                                      : [
                                          _kGreen.withAlpha(60),
                                          _kGreen.withAlpha(120),
                                        ],
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(height: 6),
                      Text(
                        days[i],
                        style: TextStyle(
                          color: isSelected ? Colors.white : _kTextDim,
                          fontSize: 11,
                          fontWeight:
                              isSelected ? FontWeight.w600 : FontWeight.w400,
                        ),
                      ),
                    ],
                  ),
                ),
              );
            }),
          ),
        ),
      ],
    );
  }

  // ─── Macros Card ───────────────────────────────────────────────────────────

  Widget _buildMacrosCard() {
    final protein = _dailySummary?.totalProtein ?? 0;
    final carbs = _dailySummary?.totalCarbs ?? 0;
    final fat = _dailySummary?.totalFat ?? 0;
    final total = protein + carbs + fat;

    final proteinGoal = _goal?.proteinGoalG ?? 150;
    final carbsGoal = _goal?.carbsGoalG ?? 250;
    final fatGoal = _goal?.fatGoalG ?? 65;

    final proteinPct = total > 0 ? (protein / total * 100) : 0.0;
    final carbsPct = total > 0 ? (carbs / total * 100) : 0.0;
    final fatPct = total > 0 ? (fat / total * 100) : 0.0;

    // Weekly averages
    double avgProtein = 0, avgCarbs = 0, avgFat = 0;
    if (_weeklySummary.isNotEmpty) {
      final logged = _weeklySummary.where((d) =>
          d.totalCalories > 0 || d.totalProtein > 0 || d.totalCarbs > 0);
      final count = logged.isEmpty ? 1 : logged.length;
      avgProtein = _weeklySummary.fold(0.0, (s, d) => s + d.totalProtein) / count;
      avgCarbs = _weeklySummary.fold(0.0, (s, d) => s + d.totalCarbs) / count;
      avgFat = _weeklySummary.fold(0.0, (s, d) => s + d.totalFat) / count;
    }

    return _DashboardCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const _SectionTitle(icon: Icons.pie_chart, title: 'Macronutrients'),
          const SizedBox(height: 16),
          Row(
            children: [
              // Macro donut
              Expanded(
                flex: 4,
                child: AnimatedBuilder(
                  animation: _ringAnim,
                  builder: (context, _) {
                    return SizedBox(
                      height: 130,
                      child: CustomPaint(
                        painter: _MacroDonutPainter(
                          protein: protein,
                          carbs: carbs,
                          fat: fat,
                          progress: _ringAnim.value,
                        ),
                        child: Center(
                          child: Column(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Text(
                                '${total.toInt()}g',
                                style: const TextStyle(
                                  color: Colors.white,
                                  fontSize: 20,
                                  fontWeight: FontWeight.w700,
                                ),
                              ),
                              const Text(
                                'total',
                                style: TextStyle(
                                  color: _kTextMuted,
                                  fontSize: 11,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    );
                  },
                ),
              ),
              const SizedBox(width: 12),
              // Legend
              Expanded(
                flex: 3,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _MacroLegendItem(
                      color: _kBlue,
                      label: 'Protein',
                      value: '${protein.toInt()}g',
                      goal: '${proteinGoal.toInt()}g',
                      pct: proteinPct,
                    ),
                    const SizedBox(height: 10),
                    _MacroLegendItem(
                      color: _kOrange,
                      label: 'Carbs',
                      value: '${carbs.toInt()}g',
                      goal: '${carbsGoal.toInt()}g',
                      pct: carbsPct,
                    ),
                    const SizedBox(height: 10),
                    _MacroLegendItem(
                      color: _kPink,
                      label: 'Fat',
                      value: '${fat.toInt()}g',
                      goal: '${fatGoal.toInt()}g',
                      pct: fatPct,
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 20),
          // Weekly macro bar chart
          _buildWeeklyMacroChart(),
          const SizedBox(height: 16),
          // Average chips
          Row(
            children: [
              const Text(
                'Avg: ',
                style: TextStyle(color: _kTextDim, fontSize: 12),
              ),
              _MacroChip(
                color: _kBlue,
                label: '${avgProtein.toInt()}g P',
              ),
              const SizedBox(width: 6),
              _MacroChip(
                color: _kOrange,
                label: '${avgCarbs.toInt()}g C',
              ),
              const SizedBox(width: 6),
              _MacroChip(
                color: _kPink,
                label: '${avgFat.toInt()}g F',
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildWeeklyMacroChart() {
    final days = ['M', 'T', 'W', 'T', 'F', 'S', 'S'];
    final weekData = List.generate(7, (i) {
      if (i < _weeklySummary.length) {
        return _weeklySummary[i];
      }
      return WeeklySummaryDay(
        date: '',
        totalCalories: 0,
        totalProtein: 0,
        totalCarbs: 0,
        totalFat: 0,
      );
    });

    final maxVal = weekData.fold<double>(
        1, (m, d) => math.max(m, d.totalProtein + d.totalCarbs + d.totalFat));

    final selectedDayIndex = _selectedDate.weekday - 1;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Weekly Macros',
          style: TextStyle(
            color: _kTextMuted,
            fontSize: 13,
            fontWeight: FontWeight.w500,
          ),
        ),
        const SizedBox(height: 12),
        SizedBox(
          height: 70,
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: List.generate(7, (i) {
              final d = weekData[i];
              final total = d.totalProtein + d.totalCarbs + d.totalFat;
              final h = maxVal > 0 ? (total / maxVal * 50) : 0.0;
              final isSelected = i == selectedDayIndex;

              // Stacked proportions
              final pFrac = total > 0 ? d.totalProtein / total : 0.33;
              final cFrac = total > 0 ? d.totalCarbs / total : 0.34;

              return Expanded(
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 3),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.end,
                    children: [
                      Flexible(
                        child: SizedBox(
                          width: double.infinity,
                          height: math.max(4, h),
                          child: ClipRRect(
                            borderRadius: BorderRadius.circular(3),
                            child: total > 0
                                ? Row(
                                    children: [
                                      Expanded(
                                        flex: (pFrac * 100).toInt().clamp(1, 100),
                                        child: Container(
                                          color: isSelected
                                              ? _kBlue
                                              : _kBlue.withAlpha(100),
                                        ),
                                      ),
                                      Expanded(
                                        flex: (cFrac * 100).toInt().clamp(1, 100),
                                        child: Container(
                                          color: isSelected
                                              ? _kOrange
                                              : _kOrange.withAlpha(100),
                                        ),
                                      ),
                                      Expanded(
                                        flex: ((1 - pFrac - cFrac) * 100)
                                            .toInt()
                                            .clamp(1, 100),
                                        child: Container(
                                          color: isSelected
                                              ? _kPink
                                              : _kPink.withAlpha(100),
                                        ),
                                      ),
                                    ],
                                  )
                                : Container(
                                    decoration: BoxDecoration(
                                      color: _kCardBorder.withAlpha(60),
                                      borderRadius: BorderRadius.circular(3),
                                    ),
                                  ),
                          ),
                        ),
                      ),
                      const SizedBox(height: 6),
                      Text(
                        days[i],
                        style: TextStyle(
                          color: isSelected ? Colors.white : _kTextDim,
                          fontSize: 11,
                          fontWeight:
                              isSelected ? FontWeight.w600 : FontWeight.w400,
                        ),
                      ),
                    ],
                  ),
                ),
              );
            }),
          ),
        ),
      ],
    );
  }

  // ─── Streak Card ───────────────────────────────────────────────────────────

  Widget _buildStreakCard() {
    final currentStreak = _streak?.currentStreak ?? 0;
    final longestStreak = _streak?.longestStreak ?? 0;

    // Build the last 7 days streak indicators
    // We show dots for the week: filled if the user logged food
    final today = DateTime.now();
    final weekStart = today.subtract(Duration(days: today.weekday - 1));
    final dayLabels = ['M', 'T', 'W', 'T', 'F', 'S', 'S'];

    return _DashboardCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const _SectionTitle(
                  icon: Icons.local_fire_department, title: 'Streak'),
              Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                decoration: BoxDecoration(
                  color: _kGreen.withAlpha(25),
                  borderRadius: BorderRadius.circular(20),
                  border: Border.all(color: _kGreen.withAlpha(60)),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const Icon(
                      Icons.local_fire_department,
                      color: _kOrange,
                      size: 18,
                    ),
                    const SizedBox(width: 4),
                    Text(
                      '$currentStreak Days',
                      style: const TextStyle(
                        color: _kOrange,
                        fontSize: 14,
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 20),
          // Weekly dot indicators
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: List.generate(7, (i) {
              final dayDate = weekStart.add(Duration(days: i));
              final isToday = dayDate.year == today.year &&
                  dayDate.month == today.month &&
                  dayDate.day == today.day;
              final isPast = dayDate.isBefore(today) || isToday;

              // Check if this day has data in weekly summary
              bool hasLogged = false;
              if (i < _weeklySummary.length) {
                hasLogged = _weeklySummary[i].totalCalories > 0;
              }

              return Column(
                children: [
                  Container(
                    width: 36,
                    height: 36,
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      color: hasLogged
                          ? _kGreen
                          : isPast
                              ? _kCardBorder.withAlpha(80)
                              : Colors.transparent,
                      border: Border.all(
                        color: isToday
                            ? _kGreen
                            : hasLogged
                                ? _kGreen
                                : _kCardBorder.withAlpha(100),
                        width: isToday ? 2 : 1,
                      ),
                    ),
                    child: Center(
                      child: hasLogged
                          ? const Icon(Icons.check, color: Colors.black, size: 18)
                          : null,
                    ),
                  ),
                  const SizedBox(height: 6),
                  Text(
                    dayLabels[i],
                    style: TextStyle(
                      color: isToday ? Colors.white : _kTextDim,
                      fontSize: 12,
                      fontWeight:
                          isToday ? FontWeight.w600 : FontWeight.w400,
                    ),
                  ),
                ],
              );
            }),
          ),
          const SizedBox(height: 16),
          // Best streak
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.emoji_events, color: _kYellow, size: 16),
              const SizedBox(width: 6),
              Text(
                'Best streak: $longestStreak days',
                style: const TextStyle(
                  color: _kTextMuted,
                  fontSize: 13,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  // ─── Weight Card ───────────────────────────────────────────────────────────

  Widget _buildWeightCard() {
    final trend = _weightTrend;
    if (trend == null || trend.summary.entriesCount == 0) {
      return _DashboardCard(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const _SectionTitle(icon: Icons.monitor_weight, title: 'Weight'),
            const SizedBox(height: 20),
            Center(
              child: Column(
                children: [
                  Icon(Icons.monitor_weight_outlined,
                      color: _kTextDim, size: 40),
                  const SizedBox(height: 8),
                  const Text(
                    'No weight data yet',
                    style: TextStyle(color: _kTextMuted, fontSize: 14),
                  ),
                  const SizedBox(height: 4),
                  const Text(
                    'Log your weight in the Progress tab',
                    style: TextStyle(color: _kTextDim, fontSize: 12),
                  ),
                ],
              ),
            ),
          ],
        ),
      );
    }

    final current = trend.summary.currentKg;
    final change = trend.summary.changeKg;
    final start = trend.summary.startKg;
    final isLoss = change < 0;
    final entries = trend.entries;

    // Find the 10-day low
    double lowWeight = current;
    for (final e in entries) {
      if (e.weightKg < lowWeight) lowWeight = e.weightKg;
    }

    return _DashboardCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const _SectionTitle(icon: Icons.monitor_weight, title: 'Weight'),
              if (entries.length >= 2)
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                  decoration: BoxDecoration(
                    color: _kPurple.withAlpha(20),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: _kPurple.withAlpha(50)),
                  ),
                  child: Text(
                    '10 Day Low',
                    style: TextStyle(
                      color: _kPurple,
                      fontSize: 11,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                current.toStringAsFixed(1),
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 36,
                  fontWeight: FontWeight.w700,
                ),
              ),
              const SizedBox(width: 4),
              const Padding(
                padding: EdgeInsets.only(bottom: 6),
                child: Text(
                  'kg',
                  style: TextStyle(
                    color: _kTextMuted,
                    fontSize: 16,
                  ),
                ),
              ),
              const Spacer(),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(
                        isLoss ? Icons.arrow_downward : Icons.arrow_upward,
                        color: isLoss ? _kGreen : AppTheme.error,
                        size: 16,
                      ),
                      const SizedBox(width: 2),
                      Text(
                        '${change.abs().toStringAsFixed(1)} kg',
                        style: TextStyle(
                          color: isLoss ? _kGreen : AppTheme.error,
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 2),
                  Text(
                    'Since start (${start.toStringAsFixed(1)} kg)',
                    style: const TextStyle(
                      color: _kTextDim,
                      fontSize: 11,
                    ),
                  ),
                ],
              ),
            ],
          ),
          if (entries.length >= 2) ...[
            const SizedBox(height: 16),
            // Mini weight sparkline
            SizedBox(
              height: 50,
              child: CustomPaint(
                size: const Size(double.infinity, 50),
                painter: _WeightSparklinePainter(
                  entries: entries,
                  lowWeight: lowWeight,
                ),
              ),
            ),
          ],
        ],
      ),
    );
  }

  // ─── Water Card ────────────────────────────────────────────────────────────

  Widget _buildWaterCard() {
    final water = _waterSummary;
    final totalMl = water?.totalMl ?? 0;
    final goalMl = water?.goalMl ?? 2500;
    final progress = goalMl > 0 ? (totalMl / goalMl).clamp(0.0, 1.0) : 0.0;
    final glasses = (totalMl / 250).floor(); // 250ml per glass
    final goalGlasses = (goalMl / 250).ceil();

    return _DashboardCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const _SectionTitle(icon: Icons.water_drop, title: 'Water'),
              Text(
                '${(progress * 100).toInt()}%',
                style: TextStyle(
                  color: progress >= 1.0 ? _kGreen : _kCyan,
                  fontSize: 14,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          // Progress bar
          ClipRRect(
            borderRadius: BorderRadius.circular(6),
            child: SizedBox(
              height: 12,
              child: Stack(
                children: [
                  Container(
                    decoration: BoxDecoration(
                      color: _kCardBorder.withAlpha(80),
                      borderRadius: BorderRadius.circular(6),
                    ),
                  ),
                  AnimatedFractionallySizedBox(
                    duration: const Duration(milliseconds: 600),
                    curve: Curves.easeOutCubic,
                    widthFactor: progress,
                    alignment: Alignment.centerLeft,
                    child: Container(
                      decoration: BoxDecoration(
                        borderRadius: BorderRadius.circular(6),
                        gradient: LinearGradient(
                          colors: progress >= 1.0
                              ? [_kGreen.withAlpha(200), _kGreen]
                              : [_kCyan.withAlpha(180), _kBlue],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 12),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              RichText(
                text: TextSpan(
                  style: const TextStyle(fontSize: 14),
                  children: [
                    TextSpan(
                      text: '${(totalMl / 1000).toStringAsFixed(1)}L',
                      style: const TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    TextSpan(
                      text: ' / ${(goalMl / 1000).toStringAsFixed(1)}L',
                      style: const TextStyle(color: _kTextMuted),
                    ),
                  ],
                ),
              ),
              Text(
                '$glasses / $goalGlasses glasses',
                style: const TextStyle(color: _kTextDim, fontSize: 12),
              ),
            ],
          ),
          const SizedBox(height: 16),
          // Quick add buttons
          Row(
            children: [
              _WaterButton(
                label: '+250ml',
                icon: Icons.local_drink,
                onTap: () => _addWater(250),
              ),
              const SizedBox(width: 8),
              _WaterButton(
                label: '+500ml',
                icon: Icons.water_drop,
                onTap: () => _addWater(500),
              ),
              const SizedBox(width: 8),
              _WaterButton(
                label: '+1L',
                icon: Icons.water,
                onTap: () => _addWater(1000),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

// ═════════════════════════════════════════════════════════════════════════════
// Reusable Widgets
// ═════════════════════════════════════════════════════════════════════════════

class _DashboardCard extends StatelessWidget {
  final Widget child;
  const _DashboardCard({required this.child});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
      child: Container(
        decoration: BoxDecoration(
          color: _kCardColor,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: _kCardBorder.withAlpha(50)),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withAlpha(30),
              blurRadius: 12,
              offset: const Offset(0, 4),
            ),
          ],
        ),
        padding: const EdgeInsets.all(20),
        child: child,
      ),
    );
  }
}

class _SectionTitle extends StatelessWidget {
  final IconData icon;
  final String title;
  const _SectionTitle({required this.icon, required this.title});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(icon, color: _kGreen, size: 20),
        const SizedBox(width: 8),
        Text(
          title,
          style: const TextStyle(
            color: Colors.white,
            fontSize: 16,
            fontWeight: FontWeight.w600,
          ),
        ),
      ],
    );
  }
}

class _CalorieStat extends StatelessWidget {
  final String label;
  final String value;
  final Color color;
  const _CalorieStat({
    required this.label,
    required this.value,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(color: _kTextMuted, fontSize: 12),
        ),
        const SizedBox(height: 2),
        Text(
          value,
          style: TextStyle(
            color: color,
            fontSize: 18,
            fontWeight: FontWeight.w700,
          ),
        ),
      ],
    );
  }
}

class _MacroLegendItem extends StatelessWidget {
  final Color color;
  final String label;
  final String value;
  final String goal;
  final double pct;

  const _MacroLegendItem({
    required this.color,
    required this.label,
    required this.value,
    required this.goal,
    required this.pct,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Container(
          width: 10,
          height: 10,
          decoration: BoxDecoration(
            color: color,
            shape: BoxShape.circle,
          ),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                label,
                style: const TextStyle(
                  color: _kTextMuted,
                  fontSize: 11,
                ),
              ),
              Row(
                children: [
                  Text(
                    value,
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  Text(
                    ' / $goal',
                    style: const TextStyle(
                      color: _kTextDim,
                      fontSize: 11,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
        Text(
          '${pct.toInt()}%',
          style: TextStyle(
            color: color,
            fontSize: 12,
            fontWeight: FontWeight.w600,
          ),
        ),
      ],
    );
  }
}

class _MacroChip extends StatelessWidget {
  final Color color;
  final String label;
  const _MacroChip({required this.color, required this.label});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: color.withAlpha(20),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: color.withAlpha(50)),
      ),
      child: Text(
        label,
        style: TextStyle(
          color: color,
          fontSize: 12,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }
}

class _WaterButton extends StatelessWidget {
  final String label;
  final IconData icon;
  final VoidCallback onTap;
  const _WaterButton({
    required this.label,
    required this.icon,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(12),
          child: Container(
            padding: const EdgeInsets.symmetric(vertical: 10),
            decoration: BoxDecoration(
              color: _kCyan.withAlpha(15),
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: _kCyan.withAlpha(40)),
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(icon, color: _kCyan, size: 16),
                const SizedBox(width: 4),
                Text(
                  label,
                  style: const TextStyle(
                    color: _kCyan,
                    fontSize: 13,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

// ═════════════════════════════════════════════════════════════════════════════
// CustomPainters
// ═════════════════════════════════════════════════════════════════════════════

class _CalorieRingPainter extends CustomPainter {
  final double consumed;
  final double goal;
  final double progress;

  _CalorieRingPainter({
    required this.consumed,
    required this.goal,
    required this.progress,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final radius = math.min(size.width, size.height) / 2 - 8;
    const strokeWidth = 12.0;

    // Background ring
    final bgPaint = Paint()
      ..color = _kCardBorder.withAlpha(60)
      ..style = PaintingStyle.stroke
      ..strokeWidth = strokeWidth
      ..strokeCap = StrokeCap.round;
    canvas.drawCircle(center, radius, bgPaint);

    // Progress ring
    final fraction = goal > 0 ? (consumed / goal).clamp(0.0, 1.0) : 0.0;
    final sweepAngle = 2 * math.pi * fraction * progress;

    if (sweepAngle > 0) {
      final rect = Rect.fromCircle(center: center, radius: radius);
      final gradient = SweepGradient(
        startAngle: -math.pi / 2,
        endAngle: -math.pi / 2 + 2 * math.pi,
        colors: consumed > goal
            ? [AppTheme.error.withAlpha(200), AppTheme.error]
            : [_kGreen.withAlpha(180), _kGreen],
      );

      final progressPaint = Paint()
        ..shader = gradient.createShader(rect)
        ..style = PaintingStyle.stroke
        ..strokeWidth = strokeWidth
        ..strokeCap = StrokeCap.round;

      canvas.drawArc(
        rect,
        -math.pi / 2,
        sweepAngle,
        false,
        progressPaint,
      );
    }
  }

  @override
  bool shouldRepaint(covariant _CalorieRingPainter oldDelegate) =>
      oldDelegate.consumed != consumed ||
      oldDelegate.goal != goal ||
      oldDelegate.progress != progress;
}

class _MacroDonutPainter extends CustomPainter {
  final double protein;
  final double carbs;
  final double fat;
  final double progress;

  _MacroDonutPainter({
    required this.protein,
    required this.carbs,
    required this.fat,
    required this.progress,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final radius = math.min(size.width, size.height) / 2 - 8;
    const strokeWidth = 14.0;
    final rect = Rect.fromCircle(center: center, radius: radius);

    // Background ring
    final bgPaint = Paint()
      ..color = _kCardBorder.withAlpha(60)
      ..style = PaintingStyle.stroke
      ..strokeWidth = strokeWidth
      ..strokeCap = StrokeCap.round;
    canvas.drawCircle(center, radius, bgPaint);

    final total = protein + carbs + fat;
    if (total <= 0) return;

    final segments = [
      (protein / total, _kBlue),
      (carbs / total, _kOrange),
      (fat / total, _kPink),
    ];

    double startAngle = -math.pi / 2;
    const gap = 0.04; // Small gap between segments

    for (final seg in segments) {
      final sweepAngle = (2 * math.pi * seg.$1 - gap) * progress;
      if (sweepAngle <= 0) continue;

      final paint = Paint()
        ..color = seg.$2
        ..style = PaintingStyle.stroke
        ..strokeWidth = strokeWidth
        ..strokeCap = StrokeCap.round;

      canvas.drawArc(rect, startAngle, sweepAngle, false, paint);
      startAngle += sweepAngle + gap;
    }
  }

  @override
  bool shouldRepaint(covariant _MacroDonutPainter oldDelegate) =>
      oldDelegate.protein != protein ||
      oldDelegate.carbs != carbs ||
      oldDelegate.fat != fat ||
      oldDelegate.progress != progress;
}

class _WeightSparklinePainter extends CustomPainter {
  final List<WeightEntry> entries;
  final double lowWeight;

  _WeightSparklinePainter({required this.entries, required this.lowWeight});

  @override
  void paint(Canvas canvas, Size size) {
    if (entries.length < 2) return;

    final weights = entries.map((e) => e.weightKg).toList();
    final minW = weights.reduce(math.min) - 0.5;
    final maxW = weights.reduce(math.max) + 0.5;
    final range = maxW - minW;

    final points = <Offset>[];
    for (int i = 0; i < entries.length; i++) {
      final x = size.width * i / (entries.length - 1);
      final y = size.height - (range > 0
          ? ((entries[i].weightKg - minW) / range * size.height)
          : size.height / 2);
      points.add(Offset(x, y));
    }

    // Draw gradient fill
    final path = Path()..moveTo(points.first.dx, size.height);
    for (final p in points) {
      path.lineTo(p.dx, p.dy);
    }
    path.lineTo(points.last.dx, size.height);
    path.close();

    final fillPaint = Paint()
      ..shader = LinearGradient(
        begin: Alignment.topCenter,
        end: Alignment.bottomCenter,
        colors: [_kPurple.withAlpha(40), _kPurple.withAlpha(5)],
      ).createShader(Rect.fromLTWH(0, 0, size.width, size.height));
    canvas.drawPath(path, fillPaint);

    // Draw line
    final linePaint = Paint()
      ..color = _kPurple
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2
      ..strokeCap = StrokeCap.round;

    final linePath = Path()..moveTo(points.first.dx, points.first.dy);
    for (int i = 1; i < points.length; i++) {
      linePath.lineTo(points[i].dx, points[i].dy);
    }
    canvas.drawPath(linePath, linePaint);

    // Draw the last point (current)
    final dotPaint = Paint()..color = _kPurple;
    canvas.drawCircle(points.last, 4, dotPaint);
    final dotOutlinePaint = Paint()
      ..color = _kCardColor
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2;
    canvas.drawCircle(points.last, 4, dotOutlinePaint);
  }

  @override
  bool shouldRepaint(covariant _WeightSparklinePainter oldDelegate) =>
      oldDelegate.entries != entries;
}
