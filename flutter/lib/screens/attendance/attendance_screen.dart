import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';
import 'package:urja_flutter/l10n/app_localizations.dart';

import '../../config/theme.dart';
import '../../providers/providers.dart';

class AttendanceScreen extends ConsumerWidget {
  const AttendanceScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = AppLocalizations.of(context)!;
    final attendanceAsync = ref.watch(attendanceProvider);
    final streaksAsync = ref.watch(streaksProvider);

    return Scaffold(
      appBar: AppBar(title: Text(l10n.attendanceHistory)),
      body: RefreshIndicator(
        onRefresh: () async {
          ref.invalidate(attendanceProvider);
          ref.invalidate(streaksProvider);
        },
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            // Streak summary
            streaksAsync.when(
              data: (streaks) {
                if (streaks.isEmpty) return const SizedBox.shrink();
                final streak = streaks.first;
                return Card(
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceAround,
                      children: [
                        _StreakStat(
                          label: l10n.currentStreak,
                          value: '${streak.currentStreak}',
                          icon: Icons.local_fire_department,
                          color: UrjaColors.warning,
                        ),
                        Container(
                          width: 1,
                          height: 40,
                          color: UrjaColors.borderDefault,
                        ),
                        _StreakStat(
                          label: l10n.bestStreak,
                          value: '${streak.longestStreak}',
                          icon: Icons.emoji_events,
                          color: UrjaColors.accent,
                        ),
                      ],
                    ),
                  ),
                );
              },
              loading: () => const SizedBox.shrink(),
              error: (_, _) => const SizedBox.shrink(),
            ),
            const SizedBox(height: 16),

            // Attendance list
            attendanceAsync.when(
              data: (records) {
                if (records.isEmpty) {
                  return Center(
                    child: Padding(
                      padding: const EdgeInsets.all(48),
                      child: Column(
                        children: [
                          const Icon(Icons.calendar_today,
                              size: 48, color: UrjaColors.foregroundMuted),
                          const SizedBox(height: 16),
                          Text(l10n.noAttendance,
                              style: Theme.of(context).textTheme.bodyLarge),
                        ],
                      ),
                    ),
                  );
                }

                return Column(
                  children: records.map((record) {
                    final dateStr =
                        DateFormat('MMM dd, yyyy').format(record.checkInAt);
                    final timeStr =
                        DateFormat('hh:mm a').format(record.checkInAt);
                    final methodLabel = _methodLabel(record.method, l10n);

                    return Card(
                      margin: const EdgeInsets.only(bottom: 8),
                      child: ListTile(
                        leading: Icon(
                          _methodIcon(record.method),
                          color: UrjaColors.accent,
                        ),
                        title: Text(dateStr),
                        subtitle: Text(timeStr),
                        trailing: Container(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 8, vertical: 4),
                          decoration: BoxDecoration(
                            color: UrjaColors.accent.withAlpha(25),
                            borderRadius: BorderRadius.circular(12),
                            border: Border.all(
                                color: UrjaColors.borderAccent),
                          ),
                          child: Text(
                            methodLabel,
                            style: const TextStyle(
                              color: UrjaColors.accent,
                              fontSize: 12,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                        ),
                      ),
                    );
                  }).toList(),
                );
              },
              loading: () =>
                  const Center(child: CircularProgressIndicator()),
              error: (err, _) => Center(
                child: Text(l10n.error,
                    style: const TextStyle(color: UrjaColors.error)),
              ),
            ),
          ],
        ),
      ),
    );
  }

  IconData _methodIcon(String method) {
    switch (method) {
      case 'qr':
        return Icons.qr_code;
      case 'nfc':
        return Icons.nfc;
      case 'manual':
        return Icons.touch_app;
      default:
        return Icons.check;
    }
  }

  String _methodLabel(String method, AppLocalizations l10n) {
    switch (method) {
      case 'qr':
        return l10n.qr;
      case 'nfc':
        return l10n.nfc;
      case 'manual':
        return l10n.manual;
      default:
        return method;
    }
  }
}

class _StreakStat extends StatelessWidget {
  final String label;
  final String value;
  final IconData icon;
  final Color color;

  const _StreakStat({
    required this.label,
    required this.value,
    required this.icon,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Icon(icon, color: color, size: 28),
        const SizedBox(height: 4),
        Text(
          value,
          style: Theme.of(context).textTheme.headlineSmall?.copyWith(color: color),
        ),
        const SizedBox(height: 2),
        Text(label, style: Theme.of(context).textTheme.bodySmall),
      ],
    );
  }
}
