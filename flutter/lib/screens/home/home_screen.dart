import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:percent_indicator/circular_percent_indicator.dart';
import 'package:urja_flutter/l10n/app_localizations.dart';

import '../../config/theme.dart';
import '../../providers/providers.dart';

class HomeScreen extends ConsumerWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = AppLocalizations.of(context)!;
    final profileAsync = ref.watch(profileProvider);
    final streaksAsync = ref.watch(streaksProvider);
    final activePackage = ref.watch(activePackageProvider);
    final isOnline = ref.watch(isOnlineProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.appTitle),
        actions: [
          if (!isOnline)
            const Padding(
              padding: EdgeInsets.only(right: 8),
              child: Icon(Icons.cloud_off, color: UrjaColors.warning),
            ),
          IconButton(
            icon: const Icon(Icons.notifications_outlined),
            onPressed: () => context.push('/notifications'),
          ),
          IconButton(
            icon: const Icon(Icons.settings_outlined),
            onPressed: () => context.push('/settings'),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          ref.invalidate(profileProvider);
          ref.invalidate(streaksProvider);
          ref.invalidate(memberPackagesProvider);
        },
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            // Greeting
            profileAsync.when(
              data: (member) => Text(
                member != null
                    ? l10n.greeting(member.displayName)
                    : l10n.greeting(''),
                style: Theme.of(context).textTheme.headlineMedium,
              ),
              loading: () => const SizedBox(height: 32),
              error: (_, _) => Text(l10n.greeting('')),
            ),
            const SizedBox(height: 24),

            // Package status ring
            _PackageStatusCard(l10n: l10n, activePackage: activePackage),
            const SizedBox(height: 16),

            // Today's check-in status + streak
            Row(
              children: [
                Expanded(
                  child: _StatusCard(
                    title: l10n.todayStatus,
                    icon: Icons.check_circle_outline,
                    value: l10n.notCheckedIn,
                    color: UrjaColors.foregroundMuted,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: streaksAsync.when(
                    data: (streaks) {
                      final current = streaks.isNotEmpty
                          ? streaks.first.currentStreak
                          : 0;
                      return _StatusCard(
                        title: l10n.currentStreak,
                        icon: Icons.local_fire_department,
                        value: '$current ${l10n.days}',
                        color: current > 0
                            ? UrjaColors.warning
                            : UrjaColors.foregroundMuted,
                      );
                    },
                    loading: () => _StatusCard(
                      title: l10n.currentStreak,
                      icon: Icons.local_fire_department,
                      value: '...',
                      color: UrjaColors.foregroundMuted,
                    ),
                    error: (_, _) => _StatusCard(
                      title: l10n.currentStreak,
                      icon: Icons.local_fire_department,
                      value: '-',
                      color: UrjaColors.foregroundMuted,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 24),

            // Quick actions
            Text(
              'Quick Actions',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                _QuickAction(
                  icon: Icons.qr_code_scanner,
                  label: l10n.checkIn,
                  onTap: () => context.go('/checkin'),
                ),
                const SizedBox(width: 12),
                _QuickAction(
                  icon: Icons.leaderboard_outlined,
                  label: l10n.leaderboard,
                  onTap: () => context.push('/leaderboard'),
                ),
                const SizedBox(width: 12),
                _QuickAction(
                  icon: Icons.health_and_safety_outlined,
                  label: l10n.health,
                  onTap: () => context.push('/health'),
                ),
                const SizedBox(width: 12),
                _QuickAction(
                  icon: Icons.emoji_events_outlined,
                  label: l10n.achievements,
                  onTap: () => context.push('/achievements'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class _PackageStatusCard extends StatelessWidget {
  final AppLocalizations l10n;
  final dynamic activePackage;

  const _PackageStatusCard({required this.l10n, required this.activePackage});

  @override
  Widget build(BuildContext context) {
    if (activePackage == null) {
      return Card(
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            children: [
              const Icon(Icons.card_membership,
                  size: 48, color: UrjaColors.foregroundMuted),
              const SizedBox(height: 12),
              Text(l10n.noActivePackage,
                  style: Theme.of(context).textTheme.titleMedium),
            ],
          ),
        ),
      );
    }

    final pkg = activePackage;
    final daysLeft = pkg.daysRemaining;
    final progress = pkg.progressFraction;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Row(
          children: [
            CircularPercentIndicator(
              radius: 45,
              lineWidth: 6,
              percent: 1.0 - progress,
              center: Text(
                '$daysLeft',
                style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                      color: UrjaColors.accent,
                    ),
              ),
              progressColor: UrjaColors.accent,
              backgroundColor: UrjaColors.borderDefault,
              circularStrokeCap: CircularStrokeCap.round,
            ),
            const SizedBox(width: 20),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    pkg.packageName ?? l10n.currentPackage,
                    style: Theme.of(context).textTheme.titleMedium,
                  ),
                  const SizedBox(height: 4),
                  Text(
                    l10n.daysLeft(daysLeft),
                    style: Theme.of(context).textTheme.bodyMedium,
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _StatusCard extends StatelessWidget {
  final String title;
  final IconData icon;
  final String value;
  final Color color;

  const _StatusCard({
    required this.title,
    required this.icon,
    required this.value,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Icon(icon, color: color, size: 24),
            const SizedBox(height: 8),
            Text(title, style: Theme.of(context).textTheme.bodySmall),
            const SizedBox(height: 4),
            Text(
              value,
              style: Theme.of(context)
                  .textTheme
                  .titleMedium
                  ?.copyWith(color: color),
            ),
          ],
        ),
      ),
    );
  }
}

class _QuickAction extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;

  const _QuickAction({
    required this.icon,
    required this.label,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 16),
          decoration: BoxDecoration(
            color: UrjaColors.surface,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(color: UrjaColors.borderDefault),
          ),
          child: Column(
            children: [
              Icon(icon, color: UrjaColors.accent, size: 28),
              const SizedBox(height: 6),
              Text(
                label,
                style: Theme.of(context).textTheme.bodySmall,
                textAlign: TextAlign.center,
                overflow: TextOverflow.ellipsis,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
