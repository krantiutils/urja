import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../config/theme.dart';
import '../../models/attendance.dart';
import '../../models/streak.dart';
import '../../providers/auth_provider.dart';
import '../../services/member_service.dart';
import '../shared/widgets/empty_state.dart';
import '../shared/widgets/loading_indicator.dart';
import '../shared/widgets/method_badge.dart';
import '../shared/widgets/stat_card.dart';

class MemberAttendanceScreen extends ConsumerStatefulWidget {
  const MemberAttendanceScreen({super.key});

  @override
  ConsumerState<MemberAttendanceScreen> createState() =>
      _MemberAttendanceScreenState();
}

class _MemberAttendanceScreenState
    extends ConsumerState<MemberAttendanceScreen> {
  late final MemberService _memberService;
  final ScrollController _scrollController = ScrollController();

  List<MemberAttendanceRecord> _records = [];
  List<MemberStreak> _streaks = [];
  bool _loading = true;
  bool _loadingMore = false;
  bool _hasMore = true;
  String? _error;

  static const int _pageSize = 20;

  @override
  void initState() {
    super.initState();
    _memberService = MemberService(ref.read(apiClientProvider));
    _scrollController.addListener(_onScroll);
    _loadData();
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollController.position.pixels >=
            _scrollController.position.maxScrollExtent - 200 &&
        !_loadingMore &&
        _hasMore) {
      _loadMore();
    }
  }

  Future<void> _loadData() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final results = await Future.wait([
        _memberService.getAttendance(limit: _pageSize, offset: 0),
        _memberService.getStreaks(),
      ]);
      setState(() {
        _records = results[0] as List<MemberAttendanceRecord>;
        _streaks = results[1] as List<MemberStreak>;
        _hasMore = _records.length >= _pageSize;
        _loading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  Future<void> _loadMore() async {
    if (_loadingMore) return;
    setState(() => _loadingMore = true);
    try {
      final more = await _memberService.getAttendance(
        limit: _pageSize,
        offset: _records.length,
      );
      setState(() {
        _records.addAll(more);
        _hasMore = more.length >= _pageSize;
        _loadingMore = false;
      });
    } catch (_) {
      setState(() => _loadingMore = false);
    }
  }

  MemberStreak? get _primaryStreak =>
      _streaks.isNotEmpty ? _streaks.first : null;

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      backgroundColor: AppTheme.background,
      appBar: AppBar(title: Text(l10n.myAttendance)),
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
              : RefreshIndicator(
                  color: AppTheme.primary,
                  onRefresh: _loadData,
                  child: _buildContent(l10n),
                ),
    );
  }

  Widget _buildContent(AppLocalizations l10n) {
    return CustomScrollView(
      controller: _scrollController,
      slivers: [
        SliverToBoxAdapter(child: _buildStreakCards(l10n)),
        SliverToBoxAdapter(child: _buildCheckInButtons(l10n)),
        SliverToBoxAdapter(
          child: Padding(
            padding: const EdgeInsets.fromLTRB(16, 24, 16, 12),
            child: Text(
              l10n.attendanceHistory,
              style: const TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: AppTheme.textPrimary,
              ),
            ),
          ),
        ),
        if (_records.isEmpty)
          SliverFillRemaining(
            hasScrollBody: false,
            child: EmptyState(
              icon: Icons.check_circle_outline,
              message: l10n.noAttendanceRecords,
            ),
          )
        else
          SliverList(
            delegate: SliverChildBuilderDelegate(
              (context, index) {
                if (index == _records.length) {
                  return _loadingMore
                      ? const Padding(
                          padding: EdgeInsets.all(16),
                          child: Center(
                            child: CircularProgressIndicator(
                                color: AppTheme.primary),
                          ),
                        )
                      : const SizedBox.shrink();
                }
                return _buildAttendanceItem(_records[index]);
              },
              childCount: _records.length + (_hasMore ? 1 : 0),
            ),
          ),
      ],
    );
  }

  Widget _buildStreakCards(AppLocalizations l10n) {
    final streak = _primaryStreak;
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 0),
      child: Row(
        children: [
          Expanded(
            child: StatCard(
              title: l10n.currentStreak,
              value: '${streak?.currentStreak ?? 0}',
              icon: Icons.local_fire_department,
              color: AppTheme.warning,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: StatCard(
              title: l10n.longestStreak,
              value: '${streak?.longestStreak ?? 0}',
              icon: Icons.emoji_events,
              color: AppTheme.primary,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCheckInButtons(AppLocalizations l10n) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 0),
      child: Row(
        children: [
          Expanded(
            child: SizedBox(
              height: 48,
              child: ElevatedButton.icon(
                onPressed: () => context.push('/check-in/qr'),
                icon: const Icon(Icons.qr_code_scanner, size: 20),
                label: Text(l10n.scanQr),
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppTheme.primary,
                ),
              ),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: SizedBox(
              height: 48,
              child: ElevatedButton.icon(
                onPressed: () => context.push('/check-in/nfc'),
                icon: const Icon(Icons.nfc, size: 20),
                label: Text(l10n.tapNfc),
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppTheme.info,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAttendanceItem(MemberAttendanceRecord record) {
    final dt = DateTime.tryParse(record.checkInAt);
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: AppTheme.primary.withAlpha(30),
          child: const Icon(Icons.check, color: AppTheme.primary, size: 20),
        ),
        title: Text(
          dt != null ? _formatDate(record.checkInAt) : record.checkInAt,
          style: const TextStyle(
            color: AppTheme.textPrimary,
            fontWeight: FontWeight.w500,
          ),
        ),
        subtitle: Text(
          dt != null ? _formatTime(record.checkInAt) : '',
          style: const TextStyle(color: AppTheme.textSecondary, fontSize: 13),
        ),
        trailing: MethodBadge(method: record.method),
      ),
    );
  }

  String _formatDate(String dateStr) {
    final dt = DateTime.tryParse(dateStr);
    if (dt == null) return dateStr;
    return '${dt.year}-${dt.month.toString().padLeft(2, '0')}-${dt.day.toString().padLeft(2, '0')}';
  }

  String _formatTime(String dateStr) {
    final dt = DateTime.tryParse(dateStr);
    if (dt == null) return '';
    final hour = dt.hour > 12 ? dt.hour - 12 : (dt.hour == 0 ? 12 : dt.hour);
    final amPm = dt.hour >= 12 ? 'PM' : 'AM';
    return '${hour.toString().padLeft(2, '0')}:${dt.minute.toString().padLeft(2, '0')} $amPm';
  }
}
