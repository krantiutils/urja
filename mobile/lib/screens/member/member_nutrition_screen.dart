import 'dart:async';
import 'dart:math';

import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../models/nutrition.dart';
import '../../models/water.dart';
import '../../providers/auth_provider.dart';
import '../../services/nutrition_service.dart';
import '../../services/water_service.dart';
import '../shared/widgets/loading_indicator.dart';

class MemberNutritionScreen extends ConsumerStatefulWidget {
  const MemberNutritionScreen({super.key});

  @override
  ConsumerState<MemberNutritionScreen> createState() =>
      _MemberNutritionScreenState();
}

class _MemberNutritionScreenState
    extends ConsumerState<MemberNutritionScreen> {
  late final NutritionService _nutritionService;
  late final WaterService _waterService;

  DateTime _selectedDate = DateTime.now();
  DailySummary? _dailySummary;
  NutritionGoal? _goal;
  List<WeeklySummaryDay> _weeklySummary = [];
  List<MealTemplate> _mealTemplates = [];
  NutritionStreak? _nutritionStreak;
  WaterDailySummary? _waterSummary;
  bool _loading = true;
  String? _error;
  bool _goalNotFound = false;

  @override
  void initState() {
    super.initState();
    final apiClient = ref.read(apiClientProvider);
    _nutritionService = NutritionService(apiClient);
    _waterService = WaterService(apiClient);
    _loadData();
  }

  String get _orgId {
    final authState = ref.read(authProvider);
    return authState.user?.orgId ?? '';
  }

  String _formatDate(DateTime dt) {
    return '${dt.year}-${dt.month.toString().padLeft(2, '0')}-${dt.day.toString().padLeft(2, '0')}';
  }

  Future<void> _loadData() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final orgId = _orgId;
      final dateStr = _formatDate(_selectedDate);

      // Calculate the start of the week (Monday)
      final weekStart =
          _selectedDate.subtract(Duration(days: _selectedDate.weekday - 1));

      final results = await Future.wait([
        _nutritionService.getDailySummary(orgId, dateStr),
        _nutritionService.getNutritionGoal(orgId).catchError((e) {
          _goalNotFound = true;
          return NutritionGoal(
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
          );
        }),
        _nutritionService.getWeeklySummary(orgId, _formatDate(weekStart)),
        _nutritionService.getMyMealTemplates(orgId),
        _nutritionService.getNutritionStreak().catchError((_) =>
            NutritionStreak(
                userId: '', currentStreak: 0, longestStreak: 0, updatedAt: '')),
        _waterService.getDailySummary(dateStr).catchError((_) =>
            WaterDailySummary(date: dateStr, totalMl: 0, goalMl: 2500, entries: [])),
      ]);

      setState(() {
        _dailySummary = results[0] as DailySummary;
        _goal = results[1] as NutritionGoal;
        _weeklySummary = results[2] as List<WeeklySummaryDay>;
        _mealTemplates = results[3] as List<MealTemplate>;
        _nutritionStreak = results[4] as NutritionStreak;
        _waterSummary = results[5] as WaterDailySummary;
        _loading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  void _changeDate(int days) {
    setState(() {
      _selectedDate = _selectedDate.add(Duration(days: days));
    });
    _loadData();
  }

  void _goToToday() {
    setState(() {
      _selectedDate = DateTime.now();
    });
    _loadData();
  }

  bool get _isToday {
    final now = DateTime.now();
    return _selectedDate.year == now.year &&
        _selectedDate.month == now.month &&
        _selectedDate.day == now.day;
  }

  Future<void> _addWater(int amountMl) async {
    try {
      await _waterService.logWater(
          amountMl: amountMl, date: _formatDate(_selectedDate));
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

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    return Scaffold(
      backgroundColor: AppTheme.background,
      appBar: AppBar(
        title: Text(l10n.nutrition),
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
                          style:
                              const TextStyle(color: AppTheme.textSecondary)),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: _loadData,
                        child: Text(l10n.retry),
                      ),
                    ],
                  ),
                )
              : RefreshIndicator(
                  color: AppTheme.primary,
                  onRefresh: _loadData,
                  child: _buildContent(),
                ),
    );
  }

  Widget _buildContent() {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        _buildDateNavigation(),
        const SizedBox(height: 16),
        _buildStreakBadge(),
        const SizedBox(height: 20),
        _buildCalorieRing(),
        const SizedBox(height: 20),
        _buildMacroProgressBars(),
        const SizedBox(height: 20),
        _buildWaterSection(),
        const SizedBox(height: 24),
        _buildMealSections(),
        const SizedBox(height: 24),
        _buildWeeklyProgress(),
        const SizedBox(height: 24),
        _buildNutritionGoalCard(),
        const SizedBox(height: 24),
        _buildMealTemplatesSection(),
        const SizedBox(height: 32),
      ],
    );
  }

  // --- Nutrition Streak Badge ---

  Widget _buildStreakBadge() {
    final streak = _nutritionStreak;
    if (streak == null || (streak.currentStreak == 0 && streak.longestStreak == 0)) {
      return const SizedBox.shrink();
    }
    final l10n = AppLocalizations.of(context)!;
    return Card(
      color: Colors.orange.withAlpha(25),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(color: Colors.orange.withAlpha(80)),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        child: Row(
          children: [
            const Icon(Icons.local_fire_department,
                color: Colors.orange, size: 28),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    '${streak.currentStreak} ${l10n.days} ${l10n.streak}',
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                      color: Colors.orange,
                    ),
                  ),
                  Text(
                    '${l10n.longestStreak}: ${streak.longestStreak} ${l10n.days}',
                    style: TextStyle(
                      fontSize: 12,
                      color: Colors.orange.withAlpha(180),
                    ),
                  ),
                ],
              ),
            ),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
              decoration: BoxDecoration(
                color: Colors.orange.withAlpha(40),
                borderRadius: BorderRadius.circular(20),
              ),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Icon(Icons.local_fire_department,
                      color: Colors.orange, size: 16),
                  const SizedBox(width: 4),
                  Text(
                    '${streak.currentStreak}',
                    style: const TextStyle(
                      fontSize: 14,
                      fontWeight: FontWeight.bold,
                      color: Colors.orange,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  // --- Water Section ---

  Widget _buildWaterSection() {
    final l10n = AppLocalizations.of(context)!;
    final water = _waterSummary;
    final totalMl = water?.totalMl ?? 0;
    final goalMl = water?.goalMl ?? 2500;
    final progress = goalMl > 0 ? (totalMl / goalMl).clamp(0.0, 1.0) : 0.0;
    final goalReached = totalMl >= goalMl;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              l10n.waterTitle,
              style: const TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: AppTheme.textPrimary,
              ),
            ),
            if (goalReached)
              Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: AppTheme.primary.withAlpha(30),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Text(
                  l10n.waterCompleted,
                  style: const TextStyle(
                    color: AppTheme.primary,
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
          ],
        ),
        const SizedBox(height: 12),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Row(
                      children: [
                        Icon(Icons.water_drop,
                            color: Colors.blue[400], size: 20),
                        const SizedBox(width: 8),
                        Text(
                          '${totalMl}ml / ${goalMl}ml',
                          style: const TextStyle(
                            color: AppTheme.textPrimary,
                            fontSize: 15,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ],
                    ),
                    Text(
                      '${(progress * 100).round()}%',
                      style: TextStyle(
                        color: goalReached
                            ? AppTheme.primary
                            : Colors.blue[400],
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
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
                    color: goalReached
                        ? AppTheme.primary
                        : Colors.blue[400],
                    minHeight: 8,
                  ),
                ),
                const SizedBox(height: 16),
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton.icon(
                        onPressed: () => _addWater(250),
                        icon: Icon(Icons.water_drop_outlined,
                            size: 16, color: Colors.blue[400]),
                        label: Text(l10n.waterGlassSize),
                        style: OutlinedButton.styleFrom(
                          foregroundColor: Colors.blue[400],
                          side: BorderSide(
                              color: Colors.blue[400]!.withAlpha(80)),
                          padding: const EdgeInsets.symmetric(
                              vertical: 10),
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: OutlinedButton.icon(
                        onPressed: () => _addWater(500),
                        icon: Icon(Icons.water_drop,
                            size: 16, color: Colors.blue[400]),
                        label: Text(l10n.waterBottleSize),
                        style: OutlinedButton.styleFrom(
                          foregroundColor: Colors.blue[400],
                          side: BorderSide(
                              color: Colors.blue[400]!.withAlpha(80)),
                          padding: const EdgeInsets.symmetric(
                              vertical: 10),
                        ),
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  // --- 1. Date Navigation ---

  Widget _buildDateNavigation() {
    final dateStr = _formatDateDisplay(_selectedDate);
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        IconButton(
          onPressed: () => _changeDate(-1),
          icon:
              const Icon(Icons.chevron_left, color: AppTheme.textPrimary, size: 28),
        ),
        GestureDetector(
          onTap: _isToday ? null : _goToToday,
          child: Column(
            children: [
              Text(
                dateStr,
                style: const TextStyle(
                  fontSize: 17,
                  fontWeight: FontWeight.bold,
                  color: AppTheme.textPrimary,
                ),
              ),
              if (!_isToday)
                const Padding(
                  padding: EdgeInsets.only(top: 2),
                  child: Text(
                    'Tap for today',
                    style: TextStyle(
                      fontSize: 11,
                      color: AppTheme.primary,
                    ),
                  ),
                ),
            ],
          ),
        ),
        IconButton(
          onPressed: () => _changeDate(1),
          icon: const Icon(Icons.chevron_right,
              color: AppTheme.textPrimary, size: 28),
        ),
        if (!_isToday)
          ActionChip(
            label: const Text('Today',
                style: TextStyle(color: AppTheme.primary, fontSize: 12)),
            backgroundColor: AppTheme.primary.withAlpha(30),
            side: BorderSide.none,
            onPressed: _goToToday,
          ),
      ],
    );
  }

  String _formatDateDisplay(DateTime dt) {
    final months = [
      'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
      'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec',
    ];
    final weekdays = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
    if (_isToday) return 'Today, ${months[dt.month - 1]} ${dt.day}';
    return '${weekdays[dt.weekday - 1]}, ${months[dt.month - 1]} ${dt.day}';
  }

  // --- 2. Calorie Ring ---

  Widget _buildCalorieRing() {
    final consumed = _dailySummary?.totalCalories ?? 0;
    final goalCal = _goal?.calorieGoal ?? 2000;
    final remaining = (goalCal - consumed).clamp(0, double.infinity);
    final ratio = goalCal > 0 ? consumed / goalCal : 0.0;

    Color ringColor;
    if (ratio > 1.0) {
      ringColor = AppTheme.error;
    } else if (ratio >= 0.8) {
      ringColor = AppTheme.warning;
    } else {
      ringColor = AppTheme.primary;
    }

    return Center(
      child: Column(
        children: [
          SizedBox(
            width: 180,
            height: 180,
            child: CustomPaint(
              painter: _CalorieRingPainter(
                ratio: ratio.clamp(0.0, 1.5),
                color: ringColor,
                backgroundColor: AppTheme.surfaceLight,
              ),
              child: Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      consumed.round().toString(),
                      style: TextStyle(
                        fontSize: 36,
                        fontWeight: FontWeight.bold,
                        color: ringColor,
                      ),
                    ),
                    Text(
                      '/ ${goalCal.round()} kcal',
                      style: const TextStyle(
                        fontSize: 13,
                        color: AppTheme.textSecondary,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
          const SizedBox(height: 8),
          Text(
            ratio > 1.0
                ? '${(consumed - goalCal).round()} over goal'
                : '${remaining.round()} remaining',
            style: TextStyle(
              fontSize: 15,
              fontWeight: FontWeight.w600,
              color: ringColor,
            ),
          ),
        ],
      ),
    );
  }

  // --- 3. Macro Progress Bars ---

  Widget _buildMacroProgressBars() {
    final summary = _dailySummary;
    final goal = _goal;

    return Row(
      children: [
        Expanded(
          child: _macroBar(
            label: 'Protein',
            consumed: summary?.totalProtein ?? 0,
            goal: goal?.proteinGoalG ?? 150,
            color: Colors.blue,
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: _macroBar(
            label: 'Carbs',
            consumed: summary?.totalCarbs ?? 0,
            goal: goal?.carbsGoalG ?? 250,
            color: Colors.amber,
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: _macroBar(
            label: 'Fat',
            consumed: summary?.totalFat ?? 0,
            goal: goal?.fatGoalG ?? 65,
            color: Colors.pink[300]!,
          ),
        ),
      ],
    );
  }

  Widget _macroBar({
    required String label,
    required double consumed,
    required double goal,
    required Color color,
  }) {
    final ratio = goal > 0 ? (consumed / goal).clamp(0.0, 1.0) : 0.0;
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(label,
                style: TextStyle(
                    color: color, fontSize: 12, fontWeight: FontWeight.w600)),
            const SizedBox(height: 6),
            ClipRRect(
              borderRadius: BorderRadius.circular(4),
              child: LinearProgressIndicator(
                value: ratio,
                backgroundColor: AppTheme.surfaceLight,
                color: color,
                minHeight: 6,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              '${consumed.round()}/${goal.round()}g',
              style:
                  const TextStyle(color: AppTheme.textSecondary, fontSize: 11),
            ),
          ],
        ),
      ),
    );
  }

  // --- 4. Meal Sections ---

  Widget _buildMealSections() {
    final mealTypes = ['breakfast', 'lunch', 'dinner', 'snack'];
    final mealLabels = {
      'breakfast': 'Breakfast',
      'lunch': 'Lunch',
      'dinner': 'Dinner',
      'snack': 'Snack',
    };
    final mealIcons = {
      'breakfast': Icons.wb_sunny_outlined,
      'lunch': Icons.light_mode_outlined,
      'dinner': Icons.nights_stay_outlined,
      'snack': Icons.cookie_outlined,
    };

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Meals',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppTheme.textPrimary,
          ),
        ),
        const SizedBox(height: 12),
        ...mealTypes.map((mealType) {
          final mealSummary = _dailySummary?.meals
              .where((m) => m.mealType == mealType)
              .toList();
          final items = mealSummary?.expand((m) => m.items).toList() ?? [];
          final totalCal =
              mealSummary?.fold<double>(0, (sum, m) => sum + m.calories) ?? 0;

          return Card(
            margin: const EdgeInsets.only(bottom: 8),
            clipBehavior: Clip.antiAlias,
            child: ExpansionTile(
              tilePadding: const EdgeInsets.symmetric(horizontal: 16),
              leading: Icon(mealIcons[mealType],
                  color: AppTheme.textSecondary, size: 22),
              title: Text(
                mealLabels[mealType]!,
                style: const TextStyle(
                  color: AppTheme.textPrimary,
                  fontWeight: FontWeight.w600,
                  fontSize: 15,
                ),
              ),
              subtitle: Text(
                '${totalCal.round()} kcal',
                style: const TextStyle(
                    color: AppTheme.textSecondary, fontSize: 12),
              ),
              trailing: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  IconButton(
                    icon: const Icon(Icons.add_circle_outline,
                        color: AppTheme.primary, size: 22),
                    onPressed: () => _showFoodSearchSheet(mealType),
                    padding: EdgeInsets.zero,
                    constraints: const BoxConstraints(),
                  ),
                  const SizedBox(width: 8),
                  const Icon(Icons.expand_more,
                      color: AppTheme.textSecondary, size: 20),
                ],
              ),
              children: [
                if (items.isEmpty)
                  Padding(
                    padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
                    child: Text(
                      'No foods logged yet',
                      style: TextStyle(
                        color: AppTheme.textSecondary.withAlpha(150),
                        fontSize: 13,
                      ),
                    ),
                  )
                else
                  ...items.map((log) => _buildFoodLogItem(log)),
              ],
            ),
          );
        }),
      ],
    );
  }

  Widget _buildFoodLogItem(FoodLog log) {
    return ListTile(
      dense: true,
      contentPadding: const EdgeInsets.symmetric(horizontal: 16),
      title: Text(
        log.foodItem?.name ?? 'Food item',
        style: const TextStyle(
            color: AppTheme.textPrimary, fontSize: 14),
      ),
      subtitle: Text(
        '${log.quantityGrams.round()}g',
        style: const TextStyle(
            color: AppTheme.textSecondary, fontSize: 12),
      ),
      trailing: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            '${log.calories.round()} kcal',
            style: const TextStyle(
                color: AppTheme.textPrimary,
                fontSize: 13,
                fontWeight: FontWeight.w500),
          ),
          const SizedBox(width: 4),
          IconButton(
            icon: Icon(Icons.delete_outline,
                color: AppTheme.error.withAlpha(180), size: 18),
            onPressed: () => _deleteFoodLog(log.id),
            padding: EdgeInsets.zero,
            constraints: const BoxConstraints(),
          ),
        ],
      ),
    );
  }

  Future<void> _deleteFoodLog(String id) async {
    try {
      await _nutritionService.deleteFoodLog(id);
      _loadData();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to delete: $e'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    }
  }

  // --- 5. Food Search Bottom Sheet ---

  void _showFoodSearchSheet(String mealType) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: AppTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (ctx) => _FoodSearchSheet(
        nutritionService: _nutritionService,
        orgId: _orgId,
        mealType: mealType,
        loggedDate: _formatDate(_selectedDate),
        onLogged: () {
          Navigator.pop(ctx);
          _loadData();
        },
      ),
    );
  }

  // --- 6. Weekly Progress ---

  Widget _buildWeeklyProgress() {
    final goalCal = _goal?.calorieGoal ?? 2000;
    final dayLabels = ['M', 'T', 'W', 'T', 'F', 'S', 'S'];
    final maxCal =
        _weeklySummary.fold<double>(goalCal, (m, d) => max(m, d.totalCalories));
    final chartHeight = 120.0;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Weekly Progress',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppTheme.textPrimary,
          ),
        ),
        const SizedBox(height: 12),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: SizedBox(
              height: chartHeight + 30,
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: List.generate(7, (index) {
                  final dayCal = index < _weeklySummary.length
                      ? _weeklySummary[index].totalCalories
                      : 0.0;
                  final barHeight =
                      maxCal > 0 ? (dayCal / maxCal * chartHeight) : 0.0;
                  final goalLineY =
                      maxCal > 0 ? (goalCal / maxCal * chartHeight) : 0.0;

                  return Expanded(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.end,
                      children: [
                        SizedBox(
                          height: chartHeight,
                          child: Stack(
                            alignment: Alignment.bottomCenter,
                            children: [
                              // Goal line (dashed)
                              Positioned(
                                bottom: goalLineY,
                                left: 0,
                                right: 0,
                                child: Container(
                                  height: 1,
                                  color: AppTheme.primary.withAlpha(80),
                                ),
                              ),
                              // Bar
                              Container(
                                width: 24,
                                height: barHeight.clamp(0, chartHeight),
                                decoration: BoxDecoration(
                                  color: dayCal > goalCal
                                      ? AppTheme.error.withAlpha(180)
                                      : AppTheme.primary.withAlpha(180),
                                  borderRadius: const BorderRadius.vertical(
                                      top: Radius.circular(4)),
                                ),
                              ),
                            ],
                          ),
                        ),
                        const SizedBox(height: 6),
                        Text(
                          dayLabels[index],
                          style: const TextStyle(
                            color: AppTheme.textSecondary,
                            fontSize: 12,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ],
                    ),
                  );
                }),
              ),
            ),
          ),
        ),
      ],
    );
  }

  // --- 7. Nutrition Goal Card ---

  Widget _buildNutritionGoalCard() {
    final goal = _goal;
    final hasGoal = goal != null && !_goalNotFound && goal.id.isNotEmpty;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Nutrition Goal',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppTheme.textPrimary,
          ),
        ),
        const SizedBox(height: 12),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: hasGoal
                ? Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Container(
                            padding: const EdgeInsets.all(8),
                            decoration: BoxDecoration(
                              color: AppTheme.primary.withAlpha(30),
                              borderRadius: BorderRadius.circular(8),
                            ),
                            child: const Icon(Icons.track_changes,
                                color: AppTheme.primary, size: 20),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  '${goal.calorieGoal.round()} kcal/day',
                                  style: const TextStyle(
                                    fontSize: 17,
                                    fontWeight: FontWeight.bold,
                                    color: AppTheme.textPrimary,
                                  ),
                                ),
                                Text(
                                  _goalTypeLabel(goal.goalType),
                                  style: const TextStyle(
                                      color: AppTheme.textSecondary,
                                      fontSize: 13),
                                ),
                              ],
                            ),
                          ),
                          TextButton(
                            onPressed: () => _showGoalSetupSheet(),
                            child: const Text('Edit'),
                          ),
                        ],
                      ),
                      const SizedBox(height: 12),
                      const Divider(height: 1),
                      const SizedBox(height: 12),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceAround,
                        children: [
                          _goalMacroChip('Protein', '${goal.proteinGoalG.round()}g',
                              Colors.blue),
                          _goalMacroChip('Carbs', '${goal.carbsGoalG.round()}g',
                              Colors.amber),
                          _goalMacroChip(
                              'Fat', '${goal.fatGoalG.round()}g', Colors.pink[300]!),
                        ],
                      ),
                    ],
                  )
                : Column(
                    children: [
                      const Icon(Icons.track_changes,
                          size: 48, color: AppTheme.textSecondary),
                      const SizedBox(height: 12),
                      const Text(
                        'Set up your nutrition goals',
                        style: TextStyle(
                          color: AppTheme.textSecondary,
                          fontSize: 15,
                        ),
                      ),
                      const SizedBox(height: 8),
                      const Text(
                        'Get personalized calorie and macro targets based on your body and goals.',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          color: AppTheme.textSecondary,
                          fontSize: 12,
                        ),
                      ),
                      const SizedBox(height: 16),
                      ElevatedButton.icon(
                        onPressed: () => _showGoalSetupSheet(),
                        icon: const Icon(Icons.calculate, size: 18),
                        label: const Text('Calculate My Goals'),
                      ),
                    ],
                  ),
          ),
        ),
      ],
    );
  }

  Widget _goalMacroChip(String label, String value, Color color) {
    return Column(
      children: [
        Container(
          width: 8,
          height: 8,
          decoration: BoxDecoration(color: color, shape: BoxShape.circle),
        ),
        const SizedBox(height: 4),
        Text(value,
            style: const TextStyle(
                color: AppTheme.textPrimary,
                fontWeight: FontWeight.w600,
                fontSize: 14)),
        Text(label,
            style: const TextStyle(
                color: AppTheme.textSecondary, fontSize: 11)),
      ],
    );
  }

  String _goalTypeLabel(String goalType) {
    switch (goalType) {
      case 'lose_weight':
        return 'Weight Loss';
      case 'gain_weight':
        return 'Weight Gain';
      case 'maintain':
        return 'Maintenance';
      case 'build_muscle':
        return 'Muscle Building';
      default:
        return goalType;
    }
  }

  void _showGoalSetupSheet() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: AppTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (ctx) => _GoalSetupSheet(
        nutritionService: _nutritionService,
        orgId: _orgId,
        existingGoal: (_goalNotFound || (_goal?.id.isEmpty ?? true)) ? null : _goal,
        onSaved: () {
          Navigator.pop(ctx);
          _goalNotFound = false;
          _loadData();
        },
      ),
    );
  }

  // --- 8. Meal Templates Section ---

  Widget _buildMealTemplatesSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text(
              'Meal Templates',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: AppTheme.textPrimary,
              ),
            ),
            TextButton.icon(
              onPressed: () => _showCreateTemplateSheet(),
              icon: const Icon(Icons.add, size: 18),
              label: const Text('Create'),
            ),
          ],
        ),
        const SizedBox(height: 8),
        if (_mealTemplates.isEmpty)
          Card(
            child: Padding(
              padding: const EdgeInsets.all(24),
              child: Center(
                child: Column(
                  children: [
                    Icon(Icons.restaurant_menu,
                        size: 40,
                        color: AppTheme.textSecondary.withAlpha(100)),
                    const SizedBox(height: 8),
                    const Text(
                      'No meal templates yet',
                      style: TextStyle(
                          color: AppTheme.textSecondary, fontSize: 14),
                    ),
                    const SizedBox(height: 4),
                    const Text(
                      'Create templates for meals you eat often',
                      style: TextStyle(
                          color: AppTheme.textSecondary, fontSize: 12),
                    ),
                  ],
                ),
              ),
            ),
          )
        else
          SizedBox(
            height: 130,
            child: ListView.builder(
              scrollDirection: Axis.horizontal,
              itemCount: _mealTemplates.length,
              itemBuilder: (context, index) {
                final template = _mealTemplates[index];
                return _buildTemplateCard(template);
              },
            ),
          ),
      ],
    );
  }

  Widget _buildTemplateCard(MealTemplate template) {
    return Container(
      width: 180,
      margin: const EdgeInsets.only(right: 12),
      child: Card(
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Expanded(
                    child: Text(
                      template.name,
                      style: const TextStyle(
                        color: AppTheme.textPrimary,
                        fontWeight: FontWeight.w600,
                        fontSize: 14,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                  InkWell(
                    onTap: () => _deleteMealTemplate(template.id),
                    child: Icon(Icons.close,
                        size: 16,
                        color: AppTheme.textSecondary.withAlpha(120)),
                  ),
                ],
              ),
              const SizedBox(height: 4),
              Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                decoration: BoxDecoration(
                  color: AppTheme.primary.withAlpha(30),
                  borderRadius: BorderRadius.circular(4),
                ),
                child: Text(
                  template.mealType,
                  style: const TextStyle(
                      color: AppTheme.primary, fontSize: 10),
                ),
              ),
              const Spacer(),
              SizedBox(
                width: double.infinity,
                height: 32,
                child: ElevatedButton(
                  onPressed: () => _quickLogTemplate(template),
                  style: ElevatedButton.styleFrom(
                    padding: EdgeInsets.zero,
                    textStyle: const TextStyle(fontSize: 12),
                  ),
                  child: const Text('Quick Log'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _quickLogTemplate(MealTemplate template) async {
    try {
      await _nutritionService.logMealTemplate(
          template.id, _orgId, _formatDate(_selectedDate));
      _loadData();
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('${template.name} logged'),
            backgroundColor: AppTheme.primary,
            behavior: SnackBarBehavior.floating,
            shape:
                RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to log template: $e'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    }
  }

  Future<void> _deleteMealTemplate(String id) async {
    try {
      await _nutritionService.deleteMealTemplate(id);
      _loadData();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to delete template: $e'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    }
  }

  void _showCreateTemplateSheet() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: AppTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (ctx) => _CreateTemplateSheet(
        nutritionService: _nutritionService,
        orgId: _orgId,
        onCreated: () {
          Navigator.pop(ctx);
          _loadData();
        },
      ),
    );
  }
}

// =============================================================================
// Custom Painter for Calorie Ring
// =============================================================================

class _CalorieRingPainter extends CustomPainter {
  final double ratio;
  final Color color;
  final Color backgroundColor;

  _CalorieRingPainter({
    required this.ratio,
    required this.color,
    required this.backgroundColor,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final radius = min(size.width, size.height) / 2 - 10;
    const strokeWidth = 12.0;

    // Background arc
    final bgPaint = Paint()
      ..color = backgroundColor
      ..style = PaintingStyle.stroke
      ..strokeWidth = strokeWidth
      ..strokeCap = StrokeCap.round;
    canvas.drawCircle(center, radius, bgPaint);

    // Progress arc
    final progressPaint = Paint()
      ..color = color
      ..style = PaintingStyle.stroke
      ..strokeWidth = strokeWidth
      ..strokeCap = StrokeCap.round;
    final sweepAngle = 2 * pi * ratio.clamp(0.0, 1.0);
    canvas.drawArc(
      Rect.fromCircle(center: center, radius: radius),
      -pi / 2,
      sweepAngle,
      false,
      progressPaint,
    );
  }

  @override
  bool shouldRepaint(covariant _CalorieRingPainter oldDelegate) =>
      oldDelegate.ratio != ratio || oldDelegate.color != color;
}

// =============================================================================
// Food Search Bottom Sheet (with Custom Food & Barcode Lookup)
// =============================================================================

class _FoodSearchSheet extends StatefulWidget {
  final NutritionService nutritionService;
  final String orgId;
  final String mealType;
  final String loggedDate;
  final VoidCallback onLogged;

  const _FoodSearchSheet({
    required this.nutritionService,
    required this.orgId,
    required this.mealType,
    required this.loggedDate,
    required this.onLogged,
  });

  @override
  State<_FoodSearchSheet> createState() => _FoodSearchSheetState();
}

class _FoodSearchSheetState extends State<_FoodSearchSheet> {
  final _searchController = TextEditingController();
  final _barcodeController = TextEditingController();
  Timer? _debounce;
  List<FoodItem> _results = [];
  bool _searching = false;
  FoodItem? _selectedFood;
  final _quantityController = TextEditingController(text: '100');
  bool _logging = false;
  bool _lookingUpBarcode = false;

  @override
  void initState() {
    super.initState();
    _doSearch('');
  }

  @override
  void dispose() {
    _searchController.dispose();
    _quantityController.dispose();
    _barcodeController.dispose();
    _debounce?.cancel();
    super.dispose();
  }

  void _onSearchChanged(String query) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 400), () {
      _doSearch(query);
    });
  }

  Future<void> _doSearch(String query) async {
    setState(() => _searching = true);
    try {
      final results = await widget.nutritionService
          .searchFoods(widget.orgId, query: query.isEmpty ? null : query);
      if (mounted) {
        setState(() {
          _results = results;
          _searching = false;
        });
      }
    } catch (e) {
      if (mounted) setState(() => _searching = false);
    }
  }

  Future<void> _logFood() async {
    if (_selectedFood == null) return;
    final qty = double.tryParse(_quantityController.text) ?? 100;
    setState(() => _logging = true);
    try {
      await widget.nutritionService.logFood(
        orgId: widget.orgId,
        foodItemId: _selectedFood!.id,
        mealType: widget.mealType,
        quantityGrams: qty,
        loggedDate: widget.loggedDate,
      );
      widget.onLogged();
    } catch (e) {
      setState(() => _logging = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to log food: $e'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    }
  }

  Future<void> _lookupBarcode() async {
    final barcode = _barcodeController.text.trim();
    if (barcode.isEmpty) return;
    setState(() => _lookingUpBarcode = true);
    try {
      final food = await widget.nutritionService.lookupBarcode(barcode);
      if (mounted) {
        setState(() {
          _selectedFood = food;
          _quantityController.text =
              (food.servingSizeG ?? 100).round().toString();
          _lookingUpBarcode = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _lookingUpBarcode = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Barcode not found: $e'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    }
  }

  void _showCreateCustomFoodDialog() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: AppTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (ctx) => _CreateCustomFoodSheet(
        nutritionService: widget.nutritionService,
        onCreated: (food) {
          Navigator.pop(ctx);
          setState(() {
            _results.insert(0, food);
            _selectedFood = food;
            _quantityController.text =
                (food.servingSizeG ?? 100).round().toString();
          });
        },
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final bottomPadding = MediaQuery.of(context).viewInsets.bottom;

    return Padding(
      padding: EdgeInsets.fromLTRB(16, 16, 16, bottomPadding + 16),
      child: SizedBox(
        height: MediaQuery.of(context).size.height * 0.7,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Handle bar
            Center(
              child: Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: AppTheme.textSecondary.withAlpha(80),
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
            ),
            const SizedBox(height: 12),
            Text(
              'Add Food - ${_mealLabel(widget.mealType)}',
              style: const TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: AppTheme.textPrimary,
              ),
            ),
            const SizedBox(height: 12),
            // Search field
            TextField(
              controller: _searchController,
              style: const TextStyle(color: AppTheme.textPrimary),
              decoration: InputDecoration(
                hintText: 'Search foods...',
                hintStyle: const TextStyle(color: AppTheme.textSecondary),
                prefixIcon:
                    const Icon(Icons.search, color: AppTheme.textSecondary),
                suffixIcon: _searchController.text.isNotEmpty
                    ? IconButton(
                        icon: const Icon(Icons.clear,
                            color: AppTheme.textSecondary, size: 18),
                        onPressed: () {
                          _searchController.clear();
                          _doSearch('');
                          setState(() => _selectedFood = null);
                        },
                      )
                    : null,
              ),
              onChanged: (v) {
                _onSearchChanged(v);
                setState(() => _selectedFood = null);
              },
            ),
            const SizedBox(height: 8),
            // Barcode lookup row
            Row(
              children: [
                Expanded(
                  child: SizedBox(
                    height: 44,
                    child: TextField(
                      controller: _barcodeController,
                      style: const TextStyle(
                          color: AppTheme.textPrimary, fontSize: 13),
                      decoration: InputDecoration(
                        hintText: 'Enter barcode...',
                        hintStyle: const TextStyle(
                            color: AppTheme.textSecondary, fontSize: 13),
                        prefixIcon: const Icon(Icons.qr_code,
                            color: AppTheme.textSecondary, size: 18),
                        isDense: true,
                        contentPadding: const EdgeInsets.symmetric(
                            horizontal: 12, vertical: 10),
                      ),
                      onSubmitted: (_) => _lookupBarcode(),
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                SizedBox(
                  height: 44,
                  child: ElevatedButton(
                    onPressed: _lookingUpBarcode ? null : _lookupBarcode,
                    style: ElevatedButton.styleFrom(
                      padding: const EdgeInsets.symmetric(horizontal: 12),
                      textStyle: const TextStyle(fontSize: 13),
                    ),
                    child: _lookingUpBarcode
                        ? const SizedBox(
                            width: 18,
                            height: 18,
                            child: CircularProgressIndicator(
                                color: Colors.white, strokeWidth: 2))
                        : const Text('Lookup'),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            // Create Custom Food button
            SizedBox(
              width: double.infinity,
              child: OutlinedButton.icon(
                onPressed: _showCreateCustomFoodDialog,
                icon: const Icon(Icons.add, size: 18),
                label: const Text('Create Custom Food'),
                style: OutlinedButton.styleFrom(
                  foregroundColor: AppTheme.primary,
                  side: BorderSide(color: AppTheme.primary.withAlpha(80)),
                  padding: const EdgeInsets.symmetric(vertical: 10),
                ),
              ),
            ),
            const SizedBox(height: 8),
            // Selected food detail or results list
            if (_selectedFood != null) ...[
              _buildSelectedFoodDetail(),
            ] else
              Expanded(
                child: _searching
                    ? const Center(
                        child: CircularProgressIndicator(
                            color: AppTheme.primary))
                    : _results.isEmpty
                        ? const Center(
                            child: Text('No foods found',
                                style: TextStyle(
                                    color: AppTheme.textSecondary)),
                          )
                        : ListView.builder(
                            itemCount: _results.length,
                            itemBuilder: (context, index) {
                              final food = _results[index];
                              return _buildFoodResultItem(food);
                            },
                          ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildFoodResultItem(FoodItem food) {
    return ListTile(
      contentPadding: EdgeInsets.zero,
      title: Text(
        food.name,
        style: const TextStyle(color: AppTheme.textPrimary, fontSize: 14),
      ),
      subtitle: Text(
        '${food.caloriesPer100g.round()} kcal/100g',
        style:
            const TextStyle(color: AppTheme.textSecondary, fontSize: 12),
      ),
      trailing: Container(
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
        decoration: BoxDecoration(
          color: AppTheme.surfaceLight,
          borderRadius: BorderRadius.circular(4),
        ),
        child: Text(
          food.category,
          style: const TextStyle(
              color: AppTheme.textSecondary, fontSize: 10),
        ),
      ),
      onTap: () {
        setState(() {
          _selectedFood = food;
          _quantityController.text =
              (food.servingSizeG ?? 100).round().toString();
        });
      },
    );
  }

  Widget _buildSelectedFoodDetail() {
    final food = _selectedFood!;
    final qty = double.tryParse(_quantityController.text) ?? 100;
    final factor = qty / 100;
    final cal = (food.caloriesPer100g * factor).round();

    return Expanded(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Back button
          InkWell(
            onTap: () => setState(() => _selectedFood = null),
            child: const Row(
              children: [
                Icon(Icons.arrow_back, color: AppTheme.primary, size: 18),
                SizedBox(width: 4),
                Text('Back to results',
                    style: TextStyle(color: AppTheme.primary, fontSize: 13)),
              ],
            ),
          ),
          const SizedBox(height: 12),
          Text(
            food.name,
            style: const TextStyle(
              fontSize: 17,
              fontWeight: FontWeight.bold,
              color: AppTheme.textPrimary,
            ),
          ),
          if (food.servingLabel != null)
            Text(
              'Serving: ${food.servingLabel} (${food.servingSizeG?.round() ?? 100}g)',
              style: const TextStyle(
                  color: AppTheme.textSecondary, fontSize: 12),
            ),
          const SizedBox(height: 16),
          // Quantity input
          TextField(
            controller: _quantityController,
            keyboardType: TextInputType.number,
            style: const TextStyle(color: AppTheme.textPrimary),
            decoration: const InputDecoration(
              labelText: 'Quantity (grams)',
              labelStyle: TextStyle(color: AppTheme.textSecondary),
              prefixIcon: Icon(Icons.scale,
                  color: AppTheme.textSecondary, size: 20),
            ),
            onChanged: (_) => setState(() {}),
          ),
          const SizedBox(height: 12),
          // Nutrition preview
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: [
              _nutrientPreview('Cal', '$cal', Colors.green),
              _nutrientPreview('P',
                  '${(food.proteinPer100g * factor).round()}g', Colors.blue),
              _nutrientPreview('C',
                  '${(food.carbsPer100g * factor).round()}g', Colors.amber),
              _nutrientPreview('F',
                  '${(food.fatPer100g * factor).round()}g', Colors.pink[300]!),
            ],
          ),
          const Spacer(),
          // Log button
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: _logging ? null : _logFood,
              child: _logging
                  ? const SizedBox(
                      width: 20,
                      height: 20,
                      child: CircularProgressIndicator(
                          color: Colors.white, strokeWidth: 2))
                  : Text('Log $cal kcal'),
            ),
          ),
        ],
      ),
    );
  }

  Widget _nutrientPreview(String label, String value, Color color) {
    return Column(
      children: [
        Text(value,
            style: TextStyle(
                color: color, fontWeight: FontWeight.bold, fontSize: 16)),
        Text(label,
            style: const TextStyle(
                color: AppTheme.textSecondary, fontSize: 11)),
      ],
    );
  }

  String _mealLabel(String mealType) {
    switch (mealType) {
      case 'breakfast':
        return 'Breakfast';
      case 'lunch':
        return 'Lunch';
      case 'dinner':
        return 'Dinner';
      case 'snack':
        return 'Snack';
      default:
        return mealType;
    }
  }
}

// =============================================================================
// Create Custom Food Bottom Sheet
// =============================================================================

class _CreateCustomFoodSheet extends StatefulWidget {
  final NutritionService nutritionService;
  final ValueChanged<FoodItem> onCreated;

  const _CreateCustomFoodSheet({
    required this.nutritionService,
    required this.onCreated,
  });

  @override
  State<_CreateCustomFoodSheet> createState() => _CreateCustomFoodSheetState();
}

class _CreateCustomFoodSheetState extends State<_CreateCustomFoodSheet> {
  final _nameCtrl = TextEditingController();
  final _caloriesCtrl = TextEditingController();
  final _proteinCtrl = TextEditingController();
  final _carbsCtrl = TextEditingController();
  final _fatCtrl = TextEditingController();
  bool _saving = false;

  @override
  void dispose() {
    _nameCtrl.dispose();
    _caloriesCtrl.dispose();
    _proteinCtrl.dispose();
    _carbsCtrl.dispose();
    _fatCtrl.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    final name = _nameCtrl.text.trim();
    final calories = double.tryParse(_caloriesCtrl.text);
    final protein = double.tryParse(_proteinCtrl.text);
    final carbs = double.tryParse(_carbsCtrl.text);
    final fat = double.tryParse(_fatCtrl.text);

    if (name.isEmpty ||
        calories == null ||
        protein == null ||
        carbs == null ||
        fat == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Please fill in all required fields'),
          backgroundColor: AppTheme.error,
        ),
      );
      return;
    }

    setState(() => _saving = true);
    try {
      final input = CreateCustomFoodInput(
        name: name,
        caloriesPer100g: calories,
        proteinPer100g: protein,
        carbsPer100g: carbs,
        fatPer100g: fat,
      );
      final food = await widget.nutritionService.createCustomFood(input);
      widget.onCreated(food);
    } catch (e) {
      setState(() => _saving = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to create food: $e'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final bottomPadding = MediaQuery.of(context).viewInsets.bottom;
    return Padding(
      padding: EdgeInsets.fromLTRB(16, 24, 16, bottomPadding + 24),
      child: SingleChildScrollView(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Handle bar
            Center(
              child: Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: AppTheme.textSecondary.withAlpha(80),
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
            ),
            const SizedBox(height: 16),
            const Text(
              'Create Custom Food',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: AppTheme.textPrimary,
              ),
            ),
            const SizedBox(height: 4),
            const Text(
              'Enter nutrition info per 100g',
              style: TextStyle(
                color: AppTheme.textSecondary,
                fontSize: 13,
              ),
            ),
            const SizedBox(height: 20),
            TextField(
              controller: _nameCtrl,
              style: const TextStyle(color: AppTheme.textPrimary),
              decoration: const InputDecoration(
                labelText: 'Food Name *',
                labelStyle: TextStyle(color: AppTheme.textSecondary),
                prefixIcon: Icon(Icons.restaurant,
                    color: AppTheme.textSecondary, size: 20),
              ),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _caloriesCtrl,
              keyboardType: TextInputType.number,
              style: const TextStyle(color: AppTheme.textPrimary),
              decoration: const InputDecoration(
                labelText: 'Calories per 100g *',
                labelStyle: TextStyle(color: AppTheme.textSecondary),
                prefixIcon: Icon(Icons.local_fire_department,
                    color: AppTheme.textSecondary, size: 20),
              ),
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _proteinCtrl,
                    keyboardType: TextInputType.number,
                    style: const TextStyle(color: AppTheme.textPrimary),
                    decoration: const InputDecoration(
                      labelText: 'Protein (g) *',
                      labelStyle: TextStyle(color: AppTheme.textSecondary),
                      isDense: true,
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: TextField(
                    controller: _carbsCtrl,
                    keyboardType: TextInputType.number,
                    style: const TextStyle(color: AppTheme.textPrimary),
                    decoration: const InputDecoration(
                      labelText: 'Carbs (g) *',
                      labelStyle: TextStyle(color: AppTheme.textSecondary),
                      isDense: true,
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: TextField(
                    controller: _fatCtrl,
                    keyboardType: TextInputType.number,
                    style: const TextStyle(color: AppTheme.textPrimary),
                    decoration: const InputDecoration(
                      labelText: 'Fat (g) *',
                      labelStyle: TextStyle(color: AppTheme.textSecondary),
                      isDense: true,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 24),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton.icon(
                onPressed: _saving ? null : _submit,
                icon: const Icon(Icons.add, size: 18),
                label: _saving
                    ? const SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(
                            color: Colors.white, strokeWidth: 2))
                    : const Text('Create Food'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

// =============================================================================
// Goal Setup Bottom Sheet
// =============================================================================

class _GoalSetupSheet extends StatefulWidget {
  final NutritionService nutritionService;
  final String orgId;
  final NutritionGoal? existingGoal;
  final VoidCallback onSaved;

  const _GoalSetupSheet({
    required this.nutritionService,
    required this.orgId,
    this.existingGoal,
    required this.onSaved,
  });

  @override
  State<_GoalSetupSheet> createState() => _GoalSetupSheetState();
}

class _GoalSetupSheetState extends State<_GoalSetupSheet> {
  late final TextEditingController _weightCtrl;
  late final TextEditingController _heightCtrl;
  late final TextEditingController _ageCtrl;
  String _gender = 'male';
  String _activityLevel = 'moderate';
  String _goalType = 'maintain';
  bool _saving = false;

  final _activityLevels = [
    'sedentary',
    'light',
    'moderate',
    'active',
    'very_active',
  ];
  final _activityLabels = {
    'sedentary': 'Sedentary',
    'light': 'Lightly Active',
    'moderate': 'Moderately Active',
    'active': 'Active',
    'very_active': 'Very Active',
  };
  final _goalTypes = ['lose_weight', 'maintain', 'gain_weight', 'build_muscle'];
  final _goalLabels = {
    'lose_weight': 'Lose Weight',
    'maintain': 'Maintain',
    'gain_weight': 'Gain Weight',
    'build_muscle': 'Build Muscle',
  };

  @override
  void initState() {
    super.initState();
    final g = widget.existingGoal;
    _weightCtrl =
        TextEditingController(text: g != null && g.weightKg > 0 ? g.weightKg.toString() : '');
    _heightCtrl =
        TextEditingController(text: g != null && g.heightCm > 0 ? g.heightCm.toString() : '');
    _ageCtrl =
        TextEditingController(text: g != null && g.age > 0 ? g.age.toString() : '');
    if (g != null) {
      _gender = g.gender.isNotEmpty ? g.gender : 'male';
      _activityLevel =
          g.activityLevel.isNotEmpty ? g.activityLevel : 'moderate';
      _goalType = g.goalType.isNotEmpty ? g.goalType : 'maintain';
    }
  }

  @override
  void dispose() {
    _weightCtrl.dispose();
    _heightCtrl.dispose();
    _ageCtrl.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    final weight = double.tryParse(_weightCtrl.text);
    final height = double.tryParse(_heightCtrl.text);
    final age = int.tryParse(_ageCtrl.text);
    if (weight == null || height == null || age == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Please fill all fields'),
          backgroundColor: AppTheme.error,
        ),
      );
      return;
    }
    setState(() => _saving = true);
    try {
      await widget.nutritionService.setNutritionGoal(
        orgId: widget.orgId,
        weightKg: weight,
        heightCm: height,
        age: age,
        gender: _gender,
        activityLevel: _activityLevel,
        goalType: _goalType,
      );
      widget.onSaved();
    } catch (e) {
      setState(() => _saving = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to save goal: $e'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final bottomPadding = MediaQuery.of(context).viewInsets.bottom;
    return Padding(
      padding: EdgeInsets.fromLTRB(16, 24, 16, bottomPadding + 24),
      child: SingleChildScrollView(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Nutrition Goal',
              style: TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                  color: AppTheme.textPrimary),
            ),
            const SizedBox(height: 20),
            // Weight
            TextField(
              controller: _weightCtrl,
              keyboardType: TextInputType.number,
              style: const TextStyle(color: AppTheme.textPrimary),
              decoration: const InputDecoration(
                labelText: 'Weight (kg)',
                labelStyle: TextStyle(color: AppTheme.textSecondary),
                prefixIcon: Icon(Icons.monitor_weight_outlined,
                    color: AppTheme.textSecondary, size: 20),
              ),
            ),
            const SizedBox(height: 12),
            // Height
            TextField(
              controller: _heightCtrl,
              keyboardType: TextInputType.number,
              style: const TextStyle(color: AppTheme.textPrimary),
              decoration: const InputDecoration(
                labelText: 'Height (cm)',
                labelStyle: TextStyle(color: AppTheme.textSecondary),
                prefixIcon: Icon(Icons.height,
                    color: AppTheme.textSecondary, size: 20),
              ),
            ),
            const SizedBox(height: 12),
            // Age
            TextField(
              controller: _ageCtrl,
              keyboardType: TextInputType.number,
              style: const TextStyle(color: AppTheme.textPrimary),
              decoration: const InputDecoration(
                labelText: 'Age',
                labelStyle: TextStyle(color: AppTheme.textSecondary),
                prefixIcon: Icon(Icons.cake_outlined,
                    color: AppTheme.textSecondary, size: 20),
              ),
            ),
            const SizedBox(height: 16),
            // Gender
            const Text('Gender',
                style: TextStyle(color: AppTheme.textSecondary, fontSize: 13)),
            const SizedBox(height: 8),
            SegmentedButton<String>(
              segments: const [
                ButtonSegment(value: 'male', label: Text('Male')),
                ButtonSegment(value: 'female', label: Text('Female')),
                ButtonSegment(value: 'other', label: Text('Other')),
              ],
              selected: {_gender},
              onSelectionChanged: (v) => setState(() => _gender = v.first),
              style: SegmentedButton.styleFrom(
                backgroundColor: AppTheme.surfaceLight,
                foregroundColor: AppTheme.textSecondary,
                selectedForegroundColor: Colors.white,
                selectedBackgroundColor: AppTheme.primary,
              ),
            ),
            const SizedBox(height: 16),
            // Activity Level
            const Text('Activity Level',
                style: TextStyle(color: AppTheme.textSecondary, fontSize: 13)),
            const SizedBox(height: 8),
            DropdownButtonFormField<String>(
              value: _activityLevel,
              dropdownColor: AppTheme.surface,
              style: const TextStyle(color: AppTheme.textPrimary),
              decoration: const InputDecoration(
                prefixIcon: Icon(Icons.directions_run,
                    color: AppTheme.textSecondary, size: 20),
              ),
              items: _activityLevels
                  .map((a) => DropdownMenuItem(
                      value: a,
                      child: Text(_activityLabels[a] ?? a)))
                  .toList(),
              onChanged: (v) => setState(() => _activityLevel = v!),
            ),
            const SizedBox(height: 16),
            // Goal Type
            const Text('Goal',
                style: TextStyle(color: AppTheme.textSecondary, fontSize: 13)),
            const SizedBox(height: 8),
            Wrap(
              spacing: 8,
              children: _goalTypes.map((g) {
                final isSelected = _goalType == g;
                return ChoiceChip(
                  label: Text(_goalLabels[g] ?? g),
                  selected: isSelected,
                  onSelected: (_) => setState(() => _goalType = g),
                  selectedColor: AppTheme.primary.withAlpha(50),
                  labelStyle: TextStyle(
                    color: isSelected ? AppTheme.primary : AppTheme.textSecondary,
                    fontWeight:
                        isSelected ? FontWeight.w600 : FontWeight.normal,
                  ),
                );
              }).toList(),
            ),
            const SizedBox(height: 24),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton.icon(
                onPressed: _saving ? null : _save,
                icon: const Icon(Icons.calculate, size: 18),
                label: _saving
                    ? const SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(
                            color: Colors.white, strokeWidth: 2))
                    : const Text('Calculate & Save'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

// =============================================================================
// Create Meal Template Bottom Sheet
// =============================================================================

class _CreateTemplateSheet extends StatefulWidget {
  final NutritionService nutritionService;
  final String orgId;
  final VoidCallback onCreated;

  const _CreateTemplateSheet({
    required this.nutritionService,
    required this.orgId,
    required this.onCreated,
  });

  @override
  State<_CreateTemplateSheet> createState() => _CreateTemplateSheetState();
}

class _CreateTemplateSheetState extends State<_CreateTemplateSheet> {
  final _nameCtrl = TextEditingController();
  String _mealType = 'breakfast';
  bool _saving = false;

  // Simple template: just a list of food item IDs with quantities
  final List<_TemplateItem> _items = [];
  List<FoodItem> _searchResults = [];
  bool _searchingFood = false;
  final _foodSearchCtrl = TextEditingController();
  Timer? _debounce;

  @override
  void dispose() {
    _nameCtrl.dispose();
    _foodSearchCtrl.dispose();
    _debounce?.cancel();
    super.dispose();
  }

  void _searchFoods(String query) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 400), () async {
      if (!mounted) return;
      setState(() => _searchingFood = true);
      try {
        final results = await widget.nutritionService
            .searchFoods(widget.orgId, query: query.isEmpty ? null : query);
        if (mounted) {
          setState(() {
            _searchResults = results;
            _searchingFood = false;
          });
        }
      } catch (_) {
        if (mounted) setState(() => _searchingFood = false);
      }
    });
  }

  Future<void> _create() async {
    if (_nameCtrl.text.trim().isEmpty || _items.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Name and at least one food item required'),
          backgroundColor: AppTheme.error,
        ),
      );
      return;
    }
    setState(() => _saving = true);
    try {
      await widget.nutritionService.createMealTemplate(
        orgId: widget.orgId,
        name: _nameCtrl.text.trim(),
        mealType: _mealType,
        items: _items
            .map((i) => {
                  'food_item_id': i.foodItem.id,
                  'quantity_grams': i.quantityGrams,
                })
            .toList(),
      );
      widget.onCreated();
    } catch (e) {
      setState(() => _saving = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to create template: $e'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final bottomPadding = MediaQuery.of(context).viewInsets.bottom;
    return Padding(
      padding: EdgeInsets.fromLTRB(16, 24, 16, bottomPadding + 24),
      child: SizedBox(
        height: MediaQuery.of(context).size.height * 0.7,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Create Meal Template',
              style: TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                  color: AppTheme.textPrimary),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _nameCtrl,
              style: const TextStyle(color: AppTheme.textPrimary),
              decoration: const InputDecoration(
                labelText: 'Template Name',
                labelStyle: TextStyle(color: AppTheme.textSecondary),
              ),
            ),
            const SizedBox(height: 12),
            DropdownButtonFormField<String>(
              value: _mealType,
              dropdownColor: AppTheme.surface,
              style: const TextStyle(color: AppTheme.textPrimary),
              decoration: const InputDecoration(
                labelText: 'Meal Type',
                labelStyle: TextStyle(color: AppTheme.textSecondary),
              ),
              items: const [
                DropdownMenuItem(value: 'breakfast', child: Text('Breakfast')),
                DropdownMenuItem(value: 'lunch', child: Text('Lunch')),
                DropdownMenuItem(value: 'dinner', child: Text('Dinner')),
                DropdownMenuItem(value: 'snack', child: Text('Snack')),
              ],
              onChanged: (v) => setState(() => _mealType = v!),
            ),
            const SizedBox(height: 12),
            // Added items
            if (_items.isNotEmpty) ...[
              const Text('Items:',
                  style: TextStyle(
                      color: AppTheme.textSecondary, fontSize: 13)),
              const SizedBox(height: 4),
              ..._items.asMap().entries.map((e) {
                final idx = e.key;
                final item = e.value;
                return Padding(
                  padding: const EdgeInsets.only(bottom: 4),
                  child: Row(
                    children: [
                      Expanded(
                        child: Text(
                          '${item.foodItem.name} (${item.quantityGrams.round()}g)',
                          style: const TextStyle(
                              color: AppTheme.textPrimary, fontSize: 13),
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.remove_circle_outline,
                            color: AppTheme.error, size: 18),
                        onPressed: () =>
                            setState(() => _items.removeAt(idx)),
                        padding: EdgeInsets.zero,
                        constraints: const BoxConstraints(),
                      ),
                    ],
                  ),
                );
              }),
              const SizedBox(height: 8),
            ],
            // Food search for adding items
            TextField(
              controller: _foodSearchCtrl,
              style: const TextStyle(color: AppTheme.textPrimary),
              decoration: const InputDecoration(
                hintText: 'Search food to add...',
                hintStyle: TextStyle(color: AppTheme.textSecondary),
                prefixIcon:
                    Icon(Icons.search, color: AppTheme.textSecondary),
                isDense: true,
              ),
              onChanged: _searchFoods,
            ),
            const SizedBox(height: 4),
            Expanded(
              child: _searchingFood
                  ? const Center(
                      child: CircularProgressIndicator(
                          color: AppTheme.primary))
                  : ListView.builder(
                      itemCount: _searchResults.length,
                      itemBuilder: (context, index) {
                        final food = _searchResults[index];
                        return ListTile(
                          dense: true,
                          contentPadding: EdgeInsets.zero,
                          title: Text(food.name,
                              style: const TextStyle(
                                  color: AppTheme.textPrimary,
                                  fontSize: 13)),
                          subtitle: Text(
                              '${food.caloriesPer100g.round()} kcal/100g',
                              style: const TextStyle(
                                  color: AppTheme.textSecondary,
                                  fontSize: 11)),
                          trailing: IconButton(
                            icon: const Icon(Icons.add_circle,
                                color: AppTheme.primary, size: 22),
                            onPressed: () {
                              setState(() {
                                _items.add(_TemplateItem(
                                  foodItem: food,
                                  quantityGrams:
                                      food.servingSizeG ?? 100,
                                ));
                              });
                            },
                          ),
                        );
                      },
                    ),
            ),
            const SizedBox(height: 8),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: _saving ? null : _create,
                child: _saving
                    ? const SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(
                            color: Colors.white, strokeWidth: 2))
                    : const Text('Create Template'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _TemplateItem {
  final FoodItem foodItem;
  final double quantityGrams;

  _TemplateItem({required this.foodItem, required this.quantityGrams});
}
