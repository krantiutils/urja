import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:urja_flutter/l10n/app_localizations.dart';

import '../../config/theme.dart';

class LeaderboardScreen extends ConsumerStatefulWidget {
  const LeaderboardScreen({super.key});

  @override
  ConsumerState<LeaderboardScreen> createState() => _LeaderboardScreenState();
}

class _LeaderboardScreenState extends ConsumerState<LeaderboardScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.leaderboard),
        bottom: TabBar(
          controller: _tabController,
          indicatorColor: UrjaColors.accent,
          labelColor: UrjaColors.accent,
          unselectedLabelColor: UrjaColors.foregroundMuted,
          tabs: [
            Tab(text: l10n.weekly),
            Tab(text: l10n.monthly),
            Tab(text: l10n.allTime),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _LeaderboardList(period: 'weekly', l10n: l10n),
          _LeaderboardList(period: 'monthly', l10n: l10n),
          _LeaderboardList(period: 'all_time', l10n: l10n),
        ],
      ),
    );
  }
}

class _LeaderboardList extends StatelessWidget {
  final String period;
  final AppLocalizations l10n;

  const _LeaderboardList({required this.period, required this.l10n});

  @override
  Widget build(BuildContext context) {
    // Placeholder — leaderboard endpoint is TODO in the Go API
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.leaderboard,
              size: 64, color: UrjaColors.foregroundMuted),
          const SizedBox(height: 16),
          Text(l10n.noData,
              style: Theme.of(context).textTheme.bodyLarge),
        ],
      ),
    );
  }
}
