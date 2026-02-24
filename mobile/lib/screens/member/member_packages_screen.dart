import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../models/package.dart';
import '../../providers/auth_provider.dart';
import '../../services/member_service.dart';
import '../shared/widgets/empty_state.dart';
import '../shared/widgets/loading_indicator.dart';
import '../shared/widgets/status_badge.dart';

class MemberPackagesScreen extends ConsumerStatefulWidget {
  const MemberPackagesScreen({super.key});

  @override
  ConsumerState<MemberPackagesScreen> createState() =>
      _MemberPackagesScreenState();
}

class _MemberPackagesScreenState extends ConsumerState<MemberPackagesScreen>
    with SingleTickerProviderStateMixin {
  late final MemberService _memberService;
  late final TabController _tabController;

  List<MemberPackage> _packages = [];
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _memberService = MemberService(ref.read(apiClientProvider));
    _tabController = TabController(length: 2, vsync: this);
    _loadData();
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final packages = await _memberService.getPackages(limit: 50);
      setState(() {
        _packages = packages;
        _loading = false;
      });
    } catch (_) {
      // Independent users may have no packages — treat as empty, not error
      setState(() {
        _packages = [];
        _loading = false;
      });
    }
  }

  List<MemberPackage> get _activePackages =>
      _packages.where((p) => p.status == 'active' && !p.isExpired).toList();

  List<MemberPackage> get _pastPackages =>
      _packages.where((p) => p.status != 'active' || p.isExpired).toList();

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      backgroundColor: AppTheme.background,
      appBar: AppBar(
        title: Text(l10n.myPackages),
        bottom: TabBar(
          controller: _tabController,
          indicatorColor: AppTheme.primary,
          labelColor: AppTheme.primary,
          unselectedLabelColor: AppTheme.textSecondary,
          tabs: [
            Tab(text: l10n.active),
            Tab(text: l10n.past),
          ],
        ),
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
                        onPressed: _loadData,
                        child: Text(l10n.retry),
                      ),
                    ],
                  ),
                )
              : TabBarView(
                  controller: _tabController,
                  children: [
                    _buildPackageList(_activePackages, l10n, isActive: true),
                    _buildPackageList(_pastPackages, l10n, isActive: false),
                  ],
                ),
    );
  }

  Widget _buildPackageList(
    List<MemberPackage> packages,
    AppLocalizations l10n, {
    required bool isActive,
  }) {
    if (packages.isEmpty) {
      return EmptyState(
        icon: Icons.card_membership,
        message: isActive ? l10n.noActivePackages : l10n.noPastPackages,
      );
    }

    return RefreshIndicator(
      color: AppTheme.primary,
      onRefresh: _loadData,
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: packages.length,
        itemBuilder: (context, index) =>
            _buildPackageCard(packages[index], l10n),
      ),
    );
  }

  Widget _buildPackageCard(MemberPackage pkg, AppLocalizations l10n) {
    final daysLeft = pkg.daysRemaining;
    final isExpired = pkg.isExpired;
    final statusLabel = isExpired ? l10n.expired : pkg.status;

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Header row: name + status badge
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Expanded(
                  child: Text(
                    pkg.packageName ?? l10n.package,
                    style: const TextStyle(
                      fontSize: 17,
                      fontWeight: FontWeight.bold,
                      color: AppTheme.textPrimary,
                    ),
                  ),
                ),
                StatusBadge(label: statusLabel),
              ],
            ),
            if (pkg.orgName != null) ...[
              const SizedBox(height: 4),
              Text(
                pkg.orgName!,
                style: const TextStyle(
                    color: AppTheme.textSecondary, fontSize: 13),
              ),
            ],
            const SizedBox(height: 12),
            const Divider(height: 1),
            const SizedBox(height: 12),
            // Amount
            Row(
              children: [
                const Icon(Icons.payment,
                    size: 16, color: AppTheme.textSecondary),
                const SizedBox(width: 6),
                Text(
                  'NPR ${pkg.amountPaid}',
                  style: const TextStyle(
                    color: AppTheme.textPrimary,
                    fontWeight: FontWeight.w600,
                    fontSize: 15,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            // Dates
            Row(
              children: [
                const Icon(Icons.calendar_today,
                    size: 16, color: AppTheme.textSecondary),
                const SizedBox(width: 6),
                Text(
                  '${_formatDate(pkg.startDate)} - ${_formatDate(pkg.endDate)}',
                  style: const TextStyle(
                      color: AppTheme.textSecondary, fontSize: 13),
                ),
              ],
            ),
            const SizedBox(height: 8),
            // Days remaining
            if (!isExpired)
              Row(
                children: [
                  Icon(
                    Icons.timer_outlined,
                    size: 16,
                    color: daysLeft <= 7 ? AppTheme.warning : AppTheme.primary,
                  ),
                  const SizedBox(width: 6),
                  Text(
                    '$daysLeft ${l10n.daysRemaining}',
                    style: TextStyle(
                      color:
                          daysLeft <= 7 ? AppTheme.warning : AppTheme.primary,
                      fontWeight: FontWeight.w600,
                      fontSize: 14,
                    ),
                  ),
                ],
              ),
          ],
        ),
      ),
    );
  }

  String _formatDate(String dateStr) {
    final dt = DateTime.tryParse(dateStr);
    if (dt == null) return dateStr;
    return '${dt.year}-${dt.month.toString().padLeft(2, '0')}-${dt.day.toString().padLeft(2, '0')}';
  }
}
