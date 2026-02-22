import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../services/org_service.dart';
import '../../models/member.dart';
import '../../models/attendance.dart';
import '../shared/widgets/stat_card.dart';
import '../shared/widgets/method_badge.dart';
import '../shared/widgets/loading_indicator.dart';

class DashboardHomeScreen extends ConsumerStatefulWidget {
  const DashboardHomeScreen({super.key});

  @override
  ConsumerState<DashboardHomeScreen> createState() =>
      _DashboardHomeScreenState();
}

class _DashboardHomeScreenState extends ConsumerState<DashboardHomeScreen> {
  bool _loading = true;
  List<OrgMember> _members = [];
  List<OrgAttendance> _attendance = [];
  int _totalMembers = 0;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) {
      if (mounted) {
        setState(() => _loading = false);
      }
      return;
    }

    final orgService = OrgService(ref.read(apiClientProvider));

    try {
      final membersResult = await orgService.getMembers(orgId, limit: 100);
      final attendance = await orgService.getAttendance(orgId, limit: 20);

      if (mounted) {
        setState(() {
          _members = membersResult.data;
          _totalMembers = membersResult.total;
          _attendance = attendance;
          _loading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _loading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load dashboard: $e')),
        );
      }
    }
  }

  int get _activeMembers =>
      _members.where((m) => m.status == 'active').length;

  List<OrgAttendance> get _todayAttendance {
    final today = DateTime.now();
    return _attendance.where((a) {
      final date = DateTime.tryParse(a.checkInAt);
      if (date == null) return false;
      return date.year == today.year &&
          date.month == today.month &&
          date.day == today.day;
    }).toList();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.dashboard),
      ),
      body: _loading
          ? const LoadingIndicator()
          : RefreshIndicator(
              onRefresh: _loadData,
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  _buildStatsGrid(l10n),
                  const SizedBox(height: 24),
                  _buildExpiringPackages(l10n),
                  const SizedBox(height: 24),
                  _buildRecentActivity(l10n),
                ],
              ),
            ),
    );
  }

  Widget _buildStatsGrid(AppLocalizations l10n) {
    return GridView.count(
      crossAxisCount: 2,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      mainAxisSpacing: 12,
      crossAxisSpacing: 12,
      childAspectRatio: 1.4,
      children: [
        StatCard(
          title: l10n.totalMembers,
          value: '$_totalMembers',
          icon: Icons.people,
          color: AppTheme.info,
        ),
        StatCard(
          title: l10n.activeMembers,
          value: '$_activeMembers',
          icon: Icons.person_outline,
          color: AppTheme.primary,
        ),
        StatCard(
          title: l10n.todayAttendance,
          value: '${_todayAttendance.length}',
          icon: Icons.how_to_reg,
          color: AppTheme.warning,
        ),
        StatCard(
          title: l10n.monthlyRevenue,
          value: '--',
          icon: Icons.attach_money,
          color: AppTheme.primary,
        ),
      ],
    );
  }

  Widget _buildExpiringPackages(AppLocalizations l10n) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          l10n.expiringPackages,
          style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 12),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Center(
              child: Text(
                l10n.noExpiringPackages,
                style: const TextStyle(color: AppTheme.textSecondary),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildRecentActivity(AppLocalizations l10n) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          l10n.recentActivity,
          style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 12),
        if (_todayAttendance.isEmpty)
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Center(
                child: Text(
                  l10n.noActivityToday,
                  style: const TextStyle(color: AppTheme.textSecondary),
                ),
              ),
            ),
          )
        else
          ...List.generate(_todayAttendance.length, (index) {
            final record = _todayAttendance[index];
            final time = DateTime.tryParse(record.checkInAt);
            final timeStr = time != null
                ? '${time.hour.toString().padLeft(2, '0')}:${time.minute.toString().padLeft(2, '0')}'
                : '';
            return Card(
              margin: const EdgeInsets.only(bottom: 8),
              child: ListTile(
                leading: const CircleAvatar(
                  backgroundColor: AppTheme.primary,
                  child: Icon(Icons.person, color: Colors.white),
                ),
                title: Text(record.memberName ?? l10n.unknownMember),
                subtitle: Text(timeStr),
                trailing: MethodBadge(method: record.method),
              ),
            );
          }),
      ],
    );
  }
}
