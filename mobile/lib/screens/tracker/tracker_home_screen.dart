import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';

class TrackerHomeScreen extends ConsumerStatefulWidget {
  const TrackerHomeScreen({super.key});

  @override
  ConsumerState<TrackerHomeScreen> createState() => _TrackerHomeScreenState();
}

class _TrackerHomeScreenState extends ConsumerState<TrackerHomeScreen> {
  @override
  Widget build(BuildContext context) {
    final auth = ref.watch(authProvider);
    final userName = auth.user?.phone ?? 'there';
    final userType = auth.user?.userType;
    final isCalorieTracker = userType == 'calorie_tracker';

    final hour = DateTime.now().hour;
    String greeting;
    if (hour < 12) {
      greeting = 'Good Morning';
    } else if (hour < 17) {
      greeting = 'Good Afternoon';
    } else {
      greeting = 'Good Evening';
    }

    return Scaffold(
      backgroundColor: AppTheme.background,
      body: SafeArea(
        child: RefreshIndicator(
          color: AppTheme.primary,
          onRefresh: () async {
            // Placeholder - will refresh data when integrated
          },
          child: ListView(
            padding: const EdgeInsets.fromLTRB(16, 24, 16, 24),
            children: [
              // Header
              Row(
                children: [
                  CircleAvatar(
                    radius: 24,
                    backgroundColor: AppTheme.primary.withAlpha(40),
                    child: const Icon(Icons.person,
                        color: AppTheme.primary, size: 24),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          '$greeting!',
                          style: const TextStyle(
                            fontSize: 22,
                            fontWeight: FontWeight.bold,
                            color: AppTheme.textPrimary,
                          ),
                        ),
                        Text(
                          isCalorieTracker
                              ? 'Track your nutrition today'
                              : "Let's crush your goals",
                          style: const TextStyle(
                            fontSize: 14,
                            color: AppTheme.textSecondary,
                          ),
                        ),
                      ],
                    ),
                  ),
                  IconButton(
                    onPressed: () => context.go('/tracker/settings'),
                    icon: const Icon(Icons.settings_outlined,
                        color: AppTheme.textSecondary),
                  ),
                ],
              ),
              const SizedBox(height: 24),

              // Calorie Summary Card
              _buildCalorieSummaryCard(),
              const SizedBox(height: 16),

              // Water Tracker Card
              _buildWaterTrackerCard(),
              const SizedBox(height: 16),

              // Workout Plan Card (hide for calorie_tracker)
              if (!isCalorieTracker) ...[
                _buildWorkoutPlanCard(),
                const SizedBox(height: 16),
              ],

              // Quick Actions
              _buildQuickActions(isCalorieTracker),
              const SizedBox(height: 24),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildCalorieSummaryCard() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
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
                  child: const Icon(Icons.local_fire_department,
                      color: AppTheme.primary, size: 20),
                ),
                const SizedBox(width: 12),
                const Text(
                  "Today's Calories",
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: AppTheme.textPrimary,
                  ),
                ),
                const Spacer(),
                TextButton(
                  onPressed: () => context.go('/tracker/nutrition'),
                  child: const Text('View'),
                ),
              ],
            ),
            const SizedBox(height: 20),
            // Calorie progress (placeholder)
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                _calorieColumn('Eaten', '0', AppTheme.primary),
                Container(
                  width: 1,
                  height: 40,
                  color: AppTheme.surfaceLight,
                ),
                _calorieColumn('Remaining', '2000', AppTheme.textSecondary),
                Container(
                  width: 1,
                  height: 40,
                  color: AppTheme.surfaceLight,
                ),
                _calorieColumn('Goal', '2000', AppTheme.info),
              ],
            ),
            const SizedBox(height: 16),
            ClipRRect(
              borderRadius: BorderRadius.circular(4),
              child: const LinearProgressIndicator(
                value: 0.0,
                backgroundColor: Color(0xFF16213E),
                color: AppTheme.primary,
                minHeight: 8,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _calorieColumn(String label, String value, Color color) {
    return Column(
      children: [
        Text(
          value,
          style: TextStyle(
            fontSize: 22,
            fontWeight: FontWeight.bold,
            color: color,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: const TextStyle(
            fontSize: 12,
            color: AppTheme.textSecondary,
          ),
        ),
      ],
    );
  }

  Widget _buildWaterTrackerCard() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: AppTheme.info.withAlpha(30),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: const Icon(Icons.water_drop,
                      color: AppTheme.info, size: 20),
                ),
                const SizedBox(width: 12),
                const Text(
                  'Water Intake',
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: AppTheme.textPrimary,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            // Water glasses placeholder
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: List.generate(8, (index) {
                return Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 4),
                  child: Icon(
                    Icons.water_drop,
                    color: AppTheme.info.withAlpha(40),
                    size: 28,
                  ),
                );
              }),
            ),
            const SizedBox(height: 12),
            Center(
              child: Text(
                '0 / 8 glasses',
                style: TextStyle(
                  fontSize: 14,
                  color: AppTheme.textSecondary,
                ),
              ),
            ),
            const SizedBox(height: 12),
            Center(
              child: ElevatedButton.icon(
                onPressed: () {
                  // Placeholder - will be implemented in water tracking step
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(
                      content: Text('Water tracking coming soon!'),
                      backgroundColor: AppTheme.info,
                    ),
                  );
                },
                icon: const Icon(Icons.add, size: 18),
                label: const Text('Add Glass'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppTheme.info,
                  padding:
                      const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildWorkoutPlanCard() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: AppTheme.warning.withAlpha(30),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: const Icon(Icons.fitness_center,
                      color: AppTheme.warning, size: 20),
                ),
                const SizedBox(width: 12),
                const Text(
                  "Today's Workout",
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: AppTheme.textPrimary,
                  ),
                ),
                const Spacer(),
                TextButton(
                  onPressed: () => context.go('/tracker/workouts'),
                  child: const Text('View'),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Container(
              width: double.infinity,
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: AppTheme.surfaceLight,
                borderRadius: BorderRadius.circular(12),
              ),
              child: Column(
                children: [
                  Icon(Icons.fitness_center,
                      size: 40,
                      color: AppTheme.textSecondary.withAlpha(100)),
                  const SizedBox(height: 8),
                  const Text(
                    'No workout plan yet',
                    style: TextStyle(
                      color: AppTheme.textSecondary,
                      fontSize: 14,
                    ),
                  ),
                  const SizedBox(height: 4),
                  const Text(
                    'Browse plans or create your own',
                    style: TextStyle(
                      color: AppTheme.textSecondary,
                      fontSize: 12,
                    ),
                  ),
                  const SizedBox(height: 12),
                  SizedBox(
                    height: 36,
                    child: ElevatedButton(
                      onPressed: () => context.go('/tracker/workouts'),
                      style: ElevatedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(horizontal: 16),
                        textStyle: const TextStyle(fontSize: 13),
                      ),
                      child: const Text('Browse Plans'),
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

  Widget _buildQuickActions(bool isCalorieTracker) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Quick Actions',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppTheme.textPrimary,
          ),
        ),
        const SizedBox(height: 12),
        Row(
          children: [
            Expanded(
              child: _quickActionCard(
                icon: Icons.restaurant,
                iconColor: AppTheme.primary,
                label: 'Log Food',
                onTap: () => context.go('/tracker/nutrition'),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: _quickActionCard(
                icon: Icons.water_drop,
                iconColor: AppTheme.info,
                label: 'Log Water',
                onTap: () {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(
                      content: Text('Water tracking coming soon!'),
                      backgroundColor: AppTheme.info,
                    ),
                  );
                },
              ),
            ),
            if (!isCalorieTracker) ...[
              const SizedBox(width: 12),
              Expanded(
                child: _quickActionCard(
                  icon: Icons.fitness_center,
                  iconColor: AppTheme.warning,
                  label: 'Workout',
                  onTap: () => context.go('/tracker/workouts'),
                ),
              ),
            ],
          ],
        ),
      ],
    );
  }

  Widget _quickActionCard({
    required IconData icon,
    required Color iconColor,
    required String label,
    required VoidCallback onTap,
  }) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 20),
        decoration: BoxDecoration(
          color: AppTheme.surface,
          borderRadius: BorderRadius.circular(12),
        ),
        child: Column(
          children: [
            Container(
              width: 44,
              height: 44,
              decoration: BoxDecoration(
                color: iconColor.withAlpha(25),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(icon, color: iconColor, size: 22),
            ),
            const SizedBox(height: 8),
            Text(
              label,
              style: const TextStyle(
                color: AppTheme.textPrimary,
                fontSize: 13,
                fontWeight: FontWeight.w500,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
