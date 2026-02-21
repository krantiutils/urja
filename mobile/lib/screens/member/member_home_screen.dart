import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../models/calendar.dart';
import '../../models/leaderboard.dart';
import '../../models/member.dart';
import '../../models/package.dart';
import '../../models/streak.dart';
import '../../providers/auth_provider.dart';
import '../../services/member_service.dart';
import '../shared/widgets/loading_indicator.dart';

class MemberHomeScreen extends ConsumerStatefulWidget {
  const MemberHomeScreen({super.key});

  @override
  ConsumerState<MemberHomeScreen> createState() => _MemberHomeScreenState();
}

class _MemberHomeScreenState extends ConsumerState<MemberHomeScreen> {
  late final MemberService _memberService;
  MemberProfile? _profile;
  List<MemberStreak> _streaks = [];
  List<MemberPackage> _packages = [];
  AttendanceCalendar? _calendar;
  LeaderboardResponse? _leaderboard;
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _memberService = MemberService(ref.read(apiClientProvider));
    _loadData();
  }

  Future<void> _loadData() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final results = await Future.wait([
        _memberService.getProfile(),
        _memberService.getStreaks(),
        _memberService.getPackages(limit: 5),
        _memberService.getAttendanceCalendar(),
        _memberService.getLeaderboard(),
      ]);
      setState(() {
        _profile = results[0] as MemberProfile;
        _streaks = results[1] as List<MemberStreak>;
        _packages = results[2] as List<MemberPackage>;
        _calendar = results[3] as AttendanceCalendar;
        _leaderboard = results[4] as LeaderboardResponse;
        _loading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  MemberStreak? get _primaryStreak =>
      _streaks.isNotEmpty ? _streaks.first : null;

  MemberPackage? get _activePackage {
    try {
      return _packages.firstWhere(
        (p) => p.status == 'active' && !p.isExpired,
      );
    } catch (_) {
      return null;
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    if (_loading) {
      return Scaffold(
        backgroundColor: AppTheme.background,
        body: const LoadingIndicator(),
      );
    }

    if (_error != null) {
      return Scaffold(
        backgroundColor: AppTheme.background,
        body: Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Icon(Icons.error_outline,
                  size: 48, color: AppTheme.error),
              const SizedBox(height: 16),
              Text(l10n.errorLoadingData,
                  style: const TextStyle(color: AppTheme.textSecondary)),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: _loadData,
                child: Text(l10n.retry),
              ),
            ],
          ),
        ),
      );
    }

    return Scaffold(
      backgroundColor: AppTheme.background,
      body: RefreshIndicator(
        color: AppTheme.primary,
        onRefresh: _loadData,
        child: ListView(
          padding: const EdgeInsets.fromLTRB(16, 60, 16, 24),
          children: [
            _buildHeader(l10n),
            const SizedBox(height: 24),
            _buildMonthlyStreakCard(l10n),
            const SizedBox(height: 16),
            _buildLeaderboardPodium(l10n),
            const SizedBox(height: 16),
            _buildSubscriptionCard(l10n),
            const SizedBox(height: 16),
            _buildGymInfoCard(l10n),
            const SizedBox(height: 16),
            _buildFeedbackButton(l10n),
            const SizedBox(height: 24),
          ],
        ),
      ),
    );
  }

  Widget _buildHeader(AppLocalizations l10n) {
    final name = _profile?.name ?? '';
    final firstName = name.contains(' ') ? name.split(' ').first : name;
    return Row(
      children: [
        CircleAvatar(
          radius: 24,
          backgroundColor: AppTheme.primary.withAlpha(40),
          backgroundImage: _profile?.avatarUrl != null
              ? NetworkImage(_profile!.avatarUrl!)
              : null,
          child: _profile?.avatarUrl == null
              ? Text(
                  firstName.isNotEmpty ? firstName[0].toUpperCase() : '?',
                  style: const TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: AppTheme.primary,
                  ),
                )
              : null,
        ),
        const SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                '${l10n.hello}, $firstName',
                style: const TextStyle(
                  fontSize: 22,
                  fontWeight: FontWeight.bold,
                  color: AppTheme.textPrimary,
                ),
              ),
              Text(
                l10n.letsGetMoving,
                style: const TextStyle(
                  fontSize: 14,
                  color: AppTheme.textSecondary,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildMonthlyStreakCard(AppLocalizations l10n) {
    final streak = _primaryStreak;
    final cal = _calendar;
    final now = DateTime.now();
    final daysInMonth = cal?.daysInMonth ?? DateTime(now.year, now.month + 1, 0).day;
    final checkInDays = cal?.checkInDays ?? [];
    final totalCheckIns = cal?.totalCheckIns ?? 0;
    final currentStreak = streak?.currentStreak ?? 0;
    final pct = daysInMonth > 0 ? (totalCheckIns / daysInMonth * 100).round() : 0;

    // Month label
    final months = ['January', 'February', 'March', 'April', 'May', 'June',
        'July', 'August', 'September', 'October', 'November', 'December'];
    final monthLabel = months[now.month - 1];

    // Find first day of month (0=Mon, 6=Sun) - we want Sun=0
    final firstDayOfMonth = DateTime(now.year, now.month, 1);
    final firstWeekday = firstDayOfMonth.weekday % 7; // Sun=0

    final dayLabels = [l10n.sun, l10n.mon, l10n.tue, l10n.wed, l10n.thu, l10n.fri, l10n.sat];

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Header row
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  monthLabel,
                  style: const TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: AppTheme.textPrimary,
                  ),
                ),
                Row(
                  children: [
                    const Icon(Icons.local_fire_department,
                        color: AppTheme.warning, size: 20),
                    const SizedBox(width: 4),
                    Text(
                      '$currentStreak ${l10n.days} ${l10n.streak}',
                      style: const TextStyle(
                        color: AppTheme.warning,
                        fontWeight: FontWeight.w600,
                        fontSize: 14,
                      ),
                    ),
                  ],
                ),
              ],
            ),
            const SizedBox(height: 16),
            // Day labels
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: dayLabels.map((d) => SizedBox(
                width: 36,
                child: Center(
                  child: Text(d, style: const TextStyle(
                    color: AppTheme.textSecondary, fontSize: 12, fontWeight: FontWeight.w500)),
                ),
              )).toList(),
            ),
            const SizedBox(height: 8),
            // Calendar grid
            ..._buildCalendarRows(firstWeekday, daysInMonth, checkInDays, now.day),
            const SizedBox(height: 16),
            // Progress bar
            ClipRRect(
              borderRadius: BorderRadius.circular(4),
              child: LinearProgressIndicator(
                value: daysInMonth > 0 ? totalCheckIns / daysInMonth : 0,
                backgroundColor: AppTheme.surface,
                color: AppTheme.primary,
                minHeight: 6,
              ),
            ),
            const SizedBox(height: 8),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  '$totalCheckIns / $daysInMonth ${l10n.days}',
                  style: const TextStyle(color: AppTheme.textSecondary, fontSize: 13),
                ),
                Text(
                  '$pct%',
                  style: const TextStyle(
                    color: AppTheme.primary,
                    fontWeight: FontWeight.w600,
                    fontSize: 13,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  List<Widget> _buildCalendarRows(int firstWeekday, int daysInMonth, List<int> checkInDays, int today) {
    final rows = <Widget>[];
    var dayNum = 1;
    final checkInSet = checkInDays.toSet();

    // Calculate how many rows we need
    final totalCells = firstWeekday + daysInMonth;
    final numRows = (totalCells / 7).ceil();

    for (var row = 0; row < numRows; row++) {
      final cells = <Widget>[];
      for (var col = 0; col < 7; col++) {
        final cellIndex = row * 7 + col;
        if (cellIndex < firstWeekday || dayNum > daysInMonth) {
          cells.add(const SizedBox(width: 36, height: 36));
        } else {
          final d = dayNum;
          final isCheckedIn = checkInSet.contains(d);
          final isToday = d == today;
          final isPast = d < today;

          Widget indicator;
          if (isCheckedIn) {
            indicator = Container(
              width: 36, height: 36,
              decoration: BoxDecoration(
                color: AppTheme.primary.withAlpha(30),
                shape: BoxShape.circle,
              ),
              child: const Icon(Icons.check, color: AppTheme.primary, size: 18),
            );
          } else if (isToday) {
            indicator = Container(
              width: 36, height: 36,
              decoration: BoxDecoration(
                border: Border.all(color: AppTheme.primary, width: 2),
                shape: BoxShape.circle,
              ),
              child: Center(
                child: Text('$d', style: const TextStyle(
                  color: AppTheme.primary, fontWeight: FontWeight.bold, fontSize: 13)),
              ),
            );
          } else if (isPast) {
            indicator = Container(
              width: 36, height: 36,
              decoration: BoxDecoration(
                color: AppTheme.error.withAlpha(20),
                shape: BoxShape.circle,
              ),
              child: Center(
                child: Text('$d', style: TextStyle(
                  color: AppTheme.error.withAlpha(150), fontSize: 13)),
              ),
            );
          } else {
            indicator = SizedBox(
              width: 36, height: 36,
              child: Center(
                child: Text('$d', style: const TextStyle(
                  color: AppTheme.textSecondary, fontSize: 13)),
              ),
            );
          }

          cells.add(indicator);
          dayNum++;
        }
      }
      rows.add(Padding(
        padding: const EdgeInsets.only(bottom: 4),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: cells,
        ),
      ));
    }
    return rows;
  }

  Widget _buildLeaderboardPodium(AppLocalizations l10n) {
    final rankings = _leaderboard?.rankings ?? [];
    if (rankings.isEmpty) return const SizedBox.shrink();

    // Reorder for podium: #2, #1, #3
    LeaderboardEntry? first, second, third;
    for (final r in rankings) {
      if (r.rank == 1) first = r;
      if (r.rank == 2) second = r;
      if (r.rank == 3) third = r;
    }

    final now = DateTime.now();
    final months = ['January', 'February', 'March', 'April', 'May', 'June',
        'July', 'August', 'September', 'October', 'November', 'December'];

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.emoji_events, color: AppTheme.warning, size: 20),
                const SizedBox(width: 8),
                Text(
                  '${l10n.top3CheckIns} \u2022 ${months[now.month - 1]}',
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: AppTheme.textPrimary,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 20),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                if (second != null) _buildPodiumEntry(second, 60, false),
                if (first != null) _buildPodiumEntry(first, 80, true),
                if (third != null) _buildPodiumEntry(third, 50, false),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildPodiumEntry(LeaderboardEntry entry, double height, bool isFirst) {
    final firstName = entry.name.contains(' ')
        ? entry.name.split(' ').first
        : entry.name;
    return Column(
      children: [
        if (isFirst)
          const Icon(Icons.workspace_premium, color: AppTheme.warning, size: 24),
        if (isFirst) const SizedBox(height: 4),
        CircleAvatar(
          radius: isFirst ? 28 : 22,
          backgroundColor: isFirst
              ? AppTheme.warning.withAlpha(40)
              : AppTheme.primary.withAlpha(30),
          backgroundImage: entry.avatarUrl != null
              ? NetworkImage(entry.avatarUrl!)
              : null,
          child: entry.avatarUrl == null
              ? Text(
                  firstName.isNotEmpty ? firstName[0].toUpperCase() : '?',
                  style: TextStyle(
                    fontSize: isFirst ? 20 : 16,
                    fontWeight: FontWeight.bold,
                    color: isFirst ? AppTheme.warning : AppTheme.primary,
                  ),
                )
              : null,
        ),
        const SizedBox(height: 6),
        Text(
          firstName,
          style: TextStyle(
            color: AppTheme.textPrimary,
            fontWeight: isFirst ? FontWeight.bold : FontWeight.w500,
            fontSize: isFirst ? 14 : 12,
          ),
          overflow: TextOverflow.ellipsis,
        ),
        Text(
          '${entry.value}',
          style: TextStyle(
            color: isFirst ? AppTheme.warning : AppTheme.textSecondary,
            fontWeight: FontWeight.bold,
            fontSize: isFirst ? 18 : 15,
          ),
        ),
        Container(
          width: isFirst ? 64 : 52,
          height: height,
          margin: const EdgeInsets.only(top: 6),
          decoration: BoxDecoration(
            color: isFirst
                ? AppTheme.warning.withAlpha(40)
                : AppTheme.primary.withAlpha(20),
            borderRadius: const BorderRadius.vertical(top: Radius.circular(8)),
          ),
          child: Center(
            child: Text(
              '#${entry.rank}',
              style: TextStyle(
                color: isFirst ? AppTheme.warning : AppTheme.textSecondary,
                fontWeight: FontWeight.bold,
                fontSize: 16,
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildSubscriptionCard(AppLocalizations l10n) {
    final pkg = _activePackage;
    if (pkg == null) {
      return Card(
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            children: [
              const Icon(Icons.card_membership,
                  size: 40, color: AppTheme.textSecondary),
              const SizedBox(height: 8),
              Text(
                l10n.noActivePackage,
                style: const TextStyle(color: AppTheme.textSecondary),
              ),
            ],
          ),
        ),
      );
    }

    final daysLeft = pkg.daysRemaining;
    final totalDays = DateTime.parse(pkg.endDate).difference(DateTime.parse(pkg.startDate)).inDays;
    final progress = totalDays > 0 ? ((totalDays - daysLeft) / totalDays).clamp(0.0, 1.0) : 0.0;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Row(
          children: [
            // Circular progress
            SizedBox(
              width: 80,
              height: 80,
              child: Stack(
                alignment: Alignment.center,
                children: [
                  SizedBox(
                    width: 80,
                    height: 80,
                    child: CircularProgressIndicator(
                      value: progress,
                      strokeWidth: 6,
                      backgroundColor: AppTheme.surface,
                      color: daysLeft <= 7 ? AppTheme.warning : AppTheme.primary,
                    ),
                  ),
                  Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        '$daysLeft',
                        style: TextStyle(
                          fontSize: 22,
                          fontWeight: FontWeight.bold,
                          color: daysLeft <= 7 ? AppTheme.warning : AppTheme.primary,
                        ),
                      ),
                      Text(
                        l10n.days,
                        style: const TextStyle(
                          fontSize: 11,
                          color: AppTheme.textSecondary,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
            const SizedBox(width: 20),
            // Package details
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          l10n.mySubscription,
                          style: const TextStyle(
                            fontSize: 13,
                            color: AppTheme.textSecondary,
                          ),
                        ),
                      ),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                        decoration: BoxDecoration(
                          color: AppTheme.primary.withAlpha(30),
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: Text(
                          l10n.subscriptionActive,
                          style: const TextStyle(
                            fontSize: 10,
                            fontWeight: FontWeight.bold,
                            color: AppTheme.primary,
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 6),
                  Text(
                    pkg.packageName ?? l10n.package,
                    style: const TextStyle(
                      fontSize: 17,
                      fontWeight: FontWeight.bold,
                      color: AppTheme.textPrimary,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    'NPR ${pkg.amountPaid}',
                    style: const TextStyle(
                      fontSize: 14,
                      color: AppTheme.textSecondary,
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

  Widget _buildGymInfoCard(AppLocalizations l10n) {
    final org = _profile?.organizations.isNotEmpty == true
        ? _profile!.organizations.first
        : null;
    if (org == null) return const SizedBox.shrink();

    return Card(
      color: const Color(0xFF1A1A2E),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Row(
          children: [
            Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                color: AppTheme.primary.withAlpha(30),
                borderRadius: BorderRadius.circular(12),
              ),
              child: const Icon(Icons.fitness_center,
                  color: AppTheme.primary, size: 24),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    l10n.myHomeClub,
                    style: const TextStyle(
                      fontSize: 12,
                      color: AppTheme.textSecondary,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    org.orgName,
                    style: const TextStyle(
                      fontSize: 17,
                      fontWeight: FontWeight.bold,
                      color: AppTheme.textPrimary,
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

  Widget _buildFeedbackButton(AppLocalizations l10n) {
    return SizedBox(
      width: double.infinity,
      height: 48,
      child: OutlinedButton.icon(
        onPressed: _showFeedbackDialog,
        icon: const Icon(Icons.chat_bubble_outline, size: 18),
        label: Text(l10n.sendFeedback),
        style: OutlinedButton.styleFrom(
          foregroundColor: AppTheme.primary,
          side: BorderSide(color: AppTheme.primary.withAlpha(80)),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
        ),
      ),
    );
  }

  void _showFeedbackDialog() {
    final l10n = AppLocalizations.of(context)!;
    final controller = TextEditingController();
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: AppTheme.surface,
        title: Text(l10n.sendFeedback,
            style: const TextStyle(color: AppTheme.textPrimary)),
        content: TextField(
          controller: controller,
          maxLines: 4,
          style: const TextStyle(color: AppTheme.textPrimary),
          decoration: InputDecoration(
            hintText: l10n.feedbackHint,
            hintStyle: const TextStyle(color: AppTheme.textSecondary),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(),
            child: Text(l10n.cancel),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.of(ctx).pop();
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text(l10n.feedbackSent)),
              );
            },
            child: Text(l10n.send),
          ),
        ],
      ),
    );
  }
}
