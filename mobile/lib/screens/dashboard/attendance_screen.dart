import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../services/org_service.dart';
import '../../models/attendance.dart';
import '../shared/widgets/stat_card.dart';
import '../shared/widgets/method_badge.dart';
import '../shared/widgets/empty_state.dart';
import '../shared/widgets/loading_indicator.dart';

class AttendanceScreen extends ConsumerStatefulWidget {
  const AttendanceScreen({super.key});

  @override
  ConsumerState<AttendanceScreen> createState() => _AttendanceScreenState();
}

class _AttendanceScreenState extends ConsumerState<AttendanceScreen> {
  bool _loading = true;
  List<OrgAttendance> _records = [];
  DateTime _selectedDate = DateTime.now();

  @override
  void initState() {
    super.initState();
    _loadAttendance();
  }

  Future<void> _loadAttendance() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) return;

    setState(() => _loading = true);
    final orgService = OrgService(ref.read(apiClientProvider));

    try {
      final records = await orgService.getAttendance(orgId, limit: 100);
      if (mounted) {
        setState(() {
          _records = records;
          _loading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _loading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load attendance: $e')),
        );
      }
    }
  }

  List<OrgAttendance> get _filteredRecords {
    return _records.where((a) {
      final date = DateTime.tryParse(a.checkInAt);
      if (date == null) return false;
      return date.year == _selectedDate.year &&
          date.month == _selectedDate.month &&
          date.day == _selectedDate.day;
    }).toList();
  }

  Future<void> _selectDate() async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: DateTime(2020),
      lastDate: DateTime.now(),
    );
    if (picked != null) {
      setState(() => _selectedDate = picked);
    }
  }

  Future<void> _showManualCheckInDialog() async {
    final l10n = AppLocalizations.of(context)!;
    final memberIdCtrl = TextEditingController();

    final result = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.manualCheckIn),
        content: TextField(
          controller: memberIdCtrl,
          decoration: InputDecoration(
            labelText: l10n.memberId,
            hintText: l10n.enterMemberId,
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: Text(l10n.cancel),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context, true),
            child: Text(l10n.checkIn),
          ),
        ],
      ),
    );

    if (result == true && memberIdCtrl.text.trim().isNotEmpty && mounted) {
      final orgId = ref.read(authProvider).user?.orgId;
      if (orgId == null) return;

      final orgService = OrgService(ref.read(apiClientProvider));
      try {
        final response =
            await orgService.manualCheckIn(orgId, memberIdCtrl.text.trim());
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(
                '${l10n.checkInSuccess} - ${l10n.streak}: ${response.streak.currentStreak}',
              ),
              backgroundColor: AppTheme.primary,
            ),
          );
          _loadAttendance();
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('${l10n.checkInFailed}: $e')),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final filtered = _filteredRecords;
    final dateStr =
        '${_selectedDate.year}-${_selectedDate.month.toString().padLeft(2, '0')}-${_selectedDate.day.toString().padLeft(2, '0')}';

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.attendance),
      ),
      body: _loading
          ? const LoadingIndicator()
          : RefreshIndicator(
              onRefresh: _loadAttendance,
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  // Stats row
                  StatCard(
                    title: l10n.todayAttendance,
                    value: '${filtered.length}',
                    icon: Icons.how_to_reg,
                    color: AppTheme.primary,
                  ),
                  const SizedBox(height: 16),

                  // Manual check-in and date filter
                  Row(
                    children: [
                      Expanded(
                        child: ElevatedButton.icon(
                          onPressed: _showManualCheckInDialog,
                          icon: const Icon(Icons.person_add),
                          label: Text(l10n.manualCheckIn),
                        ),
                      ),
                      const SizedBox(width: 12),
                      OutlinedButton.icon(
                        onPressed: _selectDate,
                        icon: const Icon(Icons.calendar_today, size: 18),
                        label: Text(dateStr),
                        style: OutlinedButton.styleFrom(
                          foregroundColor: AppTheme.textSecondary,
                          side: const BorderSide(color: AppTheme.textSecondary),
                          padding: const EdgeInsets.symmetric(
                              horizontal: 12, vertical: 14),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),

                  // Attendance list
                  if (filtered.isEmpty)
                    EmptyState(
                      icon: Icons.event_busy,
                      message: l10n.noAttendance,
                    )
                  else
                    ...List.generate(filtered.length, (index) {
                      final record = filtered[index];
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
              ),
            ),
    );
  }
}
