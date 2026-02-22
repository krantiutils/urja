import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../services/sms_service.dart';
import '../../services/org_service.dart';
import '../../models/sms.dart';
import '../../models/member.dart';
import '../shared/widgets/stat_card.dart';
import '../shared/widgets/loading_indicator.dart';

class SmsScreen extends ConsumerStatefulWidget {
  const SmsScreen({super.key});

  @override
  ConsumerState<SmsScreen> createState() => _SmsScreenState();
}

class _SmsScreenState extends ConsumerState<SmsScreen> {
  bool _loading = true;
  SmsBalance? _balance;
  List<dynamic> _history = [];
  List<OrgMember> _members = [];
  final Set<String> _selectedMemberIds = {};
  final _messageCtrl = TextEditingController();
  bool _sending = false;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  @override
  void dispose() {
    _messageCtrl.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) return;

    setState(() => _loading = true);
    final smsService = SmsService(ref.read(apiClientProvider));
    final orgService = OrgService(ref.read(apiClientProvider));

    try {
      final results = await Future.wait([
        smsService.getBalance(orgId),
        smsService.getHistory(orgId),
        orgService.getMembers(orgId, limit: 100),
      ]);

      if (mounted) {
        setState(() {
          _balance = results[0] as SmsBalance;
          _history = results[1] as List<dynamic>;
          _members = (results[2] as ({List<OrgMember> data, int total})).data;
          _loading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _loading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load SMS data: $e')),
        );
      }
    }
  }

  Future<void> _sendSms() async {
    final l10n = AppLocalizations.of(context)!;
    if (_messageCtrl.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(l10n.enterMessage)),
      );
      return;
    }
    if (_selectedMemberIds.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(l10n.selectMembers)),
      );
      return;
    }

    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) return;

    setState(() => _sending = true);
    final smsService = SmsService(ref.read(apiClientProvider));

    try {
      await smsService.send(
        orgId,
        _messageCtrl.text.trim(),
        _selectedMemberIds.toList(),
      );
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(l10n.smsSent),
            backgroundColor: AppTheme.primary,
          ),
        );
        _messageCtrl.clear();
        _selectedMemberIds.clear();
        _loadData();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to send SMS: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _sending = false);
    }
  }

  void _showMemberSelectionDialog() {
    final l10n = AppLocalizations.of(context)!;
    showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setDialogState) => AlertDialog(
          title: Text(l10n.selectMembers),
          content: SizedBox(
            width: double.maxFinite,
            height: 400,
            child: Column(
              children: [
                Row(
                  children: [
                    TextButton(
                      onPressed: () {
                        setDialogState(() {
                          _selectedMemberIds
                              .addAll(_members.map((m) => m.id));
                        });
                      },
                      child: Text(l10n.selectAll),
                    ),
                    TextButton(
                      onPressed: () {
                        setDialogState(() => _selectedMemberIds.clear());
                      },
                      child: Text(l10n.clearAll),
                    ),
                  ],
                ),
                Expanded(
                  child: ListView.builder(
                    itemCount: _members.length,
                    itemBuilder: (context, index) {
                      final member = _members[index];
                      final isSelected =
                          _selectedMemberIds.contains(member.id);
                      return CheckboxListTile(
                        title: Text(member.name),
                        subtitle: Text(member.phone),
                        value: isSelected,
                        activeColor: AppTheme.primary,
                        onChanged: (v) {
                          setDialogState(() {
                            if (v == true) {
                              _selectedMemberIds.add(member.id);
                            } else {
                              _selectedMemberIds.remove(member.id);
                            }
                          });
                        },
                      );
                    },
                  ),
                ),
              ],
            ),
          ),
          actions: [
            ElevatedButton(
              onPressed: () {
                Navigator.pop(context);
                setState(() {});
              },
              child: Text(
                  '${l10n.done} (${_selectedMemberIds.length})'),
            ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.sms),
      ),
      body: _loading
          ? const LoadingIndicator()
          : RefreshIndicator(
              onRefresh: _loadData,
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  // Balance card
                  if (_balance != null) _buildBalanceCard(l10n),
                  const SizedBox(height: 24),

                  // Send SMS section
                  _buildSendSection(l10n),
                  const SizedBox(height: 24),

                  // History section
                  _buildHistorySection(l10n),
                ],
              ),
            ),
    );
  }

  Widget _buildBalanceCard(AppLocalizations l10n) {
    final b = _balance!;
    return GridView.count(
      crossAxisCount: 2,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      mainAxisSpacing: 12,
      crossAxisSpacing: 12,
      childAspectRatio: 1.4,
      children: [
        StatCard(
          title: l10n.remaining,
          value: '${b.balance}',
          icon: Icons.sms,
          color: AppTheme.primary,
        ),
        StatCard(
          title: l10n.totalPurchased,
          value: '${b.totalPurchased}',
          icon: Icons.shopping_cart,
          color: AppTheme.info,
        ),
        StatCard(
          title: l10n.totalUsed,
          value: '${b.totalUsed}',
          icon: Icons.send,
          color: AppTheme.warning,
        ),
        StatCard(
          title: l10n.campaigns,
          value: '${b.totalCampaigns}',
          icon: Icons.campaign,
          color: AppTheme.textSecondary,
        ),
      ],
    );
  }

  Widget _buildSendSection(AppLocalizations l10n) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          l10n.sendSms,
          style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 12),
        TextField(
          controller: _messageCtrl,
          decoration: InputDecoration(
            labelText: l10n.message,
            hintText: l10n.enterMessage,
            alignLabelWithHint: true,
          ),
          maxLines: 4,
        ),
        const SizedBox(height: 12),
        OutlinedButton.icon(
          onPressed: _showMemberSelectionDialog,
          icon: const Icon(Icons.people),
          label: Text(
            _selectedMemberIds.isEmpty
                ? l10n.selectMembers
                : l10n.membersSelected(_selectedMemberIds.length),
          ),
          style: OutlinedButton.styleFrom(
            foregroundColor: AppTheme.textSecondary,
            side: const BorderSide(color: AppTheme.textSecondary),
            minimumSize: const Size(double.infinity, 48),
          ),
        ),
        const SizedBox(height: 12),
        SizedBox(
          width: double.infinity,
          child: ElevatedButton.icon(
            onPressed: _sending ? null : _sendSms,
            icon: _sending
                ? const SizedBox(
                    width: 18,
                    height: 18,
                    child: CircularProgressIndicator(
                      strokeWidth: 2,
                      color: Colors.white,
                    ),
                  )
                : const Icon(Icons.send),
            label: Text(l10n.send),
          ),
        ),
      ],
    );
  }

  Widget _buildHistorySection(AppLocalizations l10n) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          l10n.history,
          style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 12),
        if (_history.isEmpty)
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Center(
                child: Text(
                  l10n.noHistory,
                  style: const TextStyle(color: AppTheme.textSecondary),
                ),
              ),
            ),
          )
        else
          ...List.generate(_history.length, (index) {
            final item = _history[index];
            final message = item is Map ? (item['message'] ?? '') : '$item';
            final date = item is Map ? (item['created_at'] ?? '') : '';

            return Card(
              margin: const EdgeInsets.only(bottom: 8),
              child: ListTile(
                leading: const Icon(Icons.sms, color: AppTheme.info),
                title: Text(
                  message.toString(),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
                subtitle: Text(
                  date.toString(),
                  style: const TextStyle(
                      color: AppTheme.textSecondary, fontSize: 12),
                ),
              ),
            );
          }),
      ],
    );
  }
}
