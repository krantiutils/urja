import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../services/finance_service.dart';
import '../../models/due.dart';
import '../shared/widgets/search_field.dart';
import '../shared/widgets/status_badge.dart';
import '../shared/widgets/empty_state.dart';
import '../shared/widgets/loading_indicator.dart';

class DuePaymentsScreen extends ConsumerStatefulWidget {
  const DuePaymentsScreen({super.key});

  @override
  ConsumerState<DuePaymentsScreen> createState() => _DuePaymentsScreenState();
}

class _DuePaymentsScreenState extends ConsumerState<DuePaymentsScreen> {
  bool _loading = true;
  List<Due> _dues = [];
  String _filter = 'all'; // all, unpaid, paid, waived
  String _search = '';

  @override
  void initState() {
    super.initState();
    _loadDues();
  }

  Future<void> _loadDues() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) return;

    setState(() => _loading = true);
    final financeService = FinanceService(ref.read(apiClientProvider));

    try {
      final result = await financeService.getDues(
        orgId,
        status: _filter != 'all' ? _filter : null,
        search: _search.isNotEmpty ? _search : null,
      );
      if (mounted) {
        setState(() {
          _dues = result.data;
          _loading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _loading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load dues: $e')),
        );
      }
    }
  }

  void _onSearchChanged(String value) {
    _search = value;
    _loadDues();
  }

  Future<void> _showPaymentDialog(Due due) async {
    final l10n = AppLocalizations.of(context)!;
    final amountCtrl = TextEditingController(text: due.amount);
    final referenceCtrl = TextEditingController();
    String paymentMethod = 'cash';

    final result = await showDialog<bool>(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setDialogState) => AlertDialog(
          title: Text(l10n.recordPayment),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  '${due.memberName} - ${due.description}',
                  style: const TextStyle(color: AppTheme.textSecondary),
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: amountCtrl,
                  decoration: InputDecoration(
                    labelText: l10n.amount,
                    prefixText: 'Rs. ',
                  ),
                  keyboardType: TextInputType.number,
                ),
                const SizedBox(height: 12),
                DropdownButtonFormField<String>(
                  initialValue: paymentMethod,
                  decoration: InputDecoration(labelText: l10n.paymentMethod),
                  items: ['cash', 'esewa', 'bank']
                      .map((m) => DropdownMenuItem(value: m, child: Text(m)))
                      .toList(),
                  onChanged: (v) {
                    if (v != null) {
                      setDialogState(() => paymentMethod = v);
                    }
                  },
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: referenceCtrl,
                  decoration: InputDecoration(labelText: l10n.reference),
                ),
              ],
            ),
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context, false),
              child: Text(l10n.cancel),
            ),
            ElevatedButton(
              onPressed: () => Navigator.pop(context, true),
              child: Text(l10n.recordPayment),
            ),
          ],
        ),
      ),
    );

    if (result == true && mounted) {
      final orgId = ref.read(authProvider).user?.orgId;
      if (orgId == null) return;

      final financeService = FinanceService(ref.read(apiClientProvider));
      try {
        await financeService.payDue(
          orgId,
          due.id,
          amount: amountCtrl.text.trim(),
          paymentMethod: paymentMethod,
          paymentReference: referenceCtrl.text.trim().isNotEmpty
              ? referenceCtrl.text.trim()
              : null,
        );
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(l10n.paymentRecorded),
              backgroundColor: AppTheme.primary,
            ),
          );
          _loadDues();
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to record payment: $e')),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.duePayments),
      ),
      body: Column(
        children: [
          // Filter chips
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: Row(
              children: [
                _buildFilterChip(l10n.all, 'all'),
                const SizedBox(width: 8),
                _buildFilterChip(l10n.unpaid, 'unpaid'),
                const SizedBox(width: 8),
                _buildFilterChip(l10n.paid, 'paid'),
                const SizedBox(width: 8),
                _buildFilterChip(l10n.waived, 'waived'),
              ],
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: SearchField(onChanged: _onSearchChanged),
          ),
          const SizedBox(height: 8),
          if (_loading)
            const Expanded(child: LoadingIndicator())
          else if (_dues.isEmpty)
            Expanded(
              child: EmptyState(
                icon: Icons.payments_outlined,
                message: l10n.noDues,
              ),
            )
          else
            Expanded(
              child: RefreshIndicator(
                onRefresh: _loadDues,
                child: ListView.builder(
                  itemCount: _dues.length,
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  itemBuilder: (context, index) {
                    final due = _dues[index];
                    final dueDate = DateTime.tryParse(due.dueDate);
                    final dueDateStr = dueDate != null
                        ? '${dueDate.year}-${dueDate.month.toString().padLeft(2, '0')}-${dueDate.day.toString().padLeft(2, '0')}'
                        : due.dueDate;

                    return Card(
                      margin: const EdgeInsets.only(bottom: 8),
                      child: ListTile(
                        title: Text(due.memberName),
                        subtitle: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              'Rs. ${due.amount}',
                              style: const TextStyle(
                                color: AppTheme.primary,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                            Text(
                              '${l10n.dueDate}: $dueDateStr',
                              style: const TextStyle(
                                color: AppTheme.textSecondary,
                                fontSize: 12,
                              ),
                            ),
                            if (due.description.isNotEmpty)
                              Text(
                                due.description,
                                style: const TextStyle(
                                  color: AppTheme.textSecondary,
                                  fontSize: 12,
                                ),
                                maxLines: 1,
                                overflow: TextOverflow.ellipsis,
                              ),
                          ],
                        ),
                        trailing: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            StatusBadge(label: due.status),
                            if (due.status == 'unpaid') ...[
                              const SizedBox(width: 8),
                              IconButton(
                                icon: const Icon(Icons.payment,
                                    color: AppTheme.primary),
                                onPressed: () => _showPaymentDialog(due),
                              ),
                            ],
                          ],
                        ),
                        isThreeLine: true,
                      ),
                    );
                  },
                ),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildFilterChip(String label, String value) {
    final isSelected = _filter == value;
    return ChoiceChip(
      label: Text(label),
      selected: isSelected,
      onSelected: (_) {
        setState(() => _filter = value);
        _loadDues();
      },
      selectedColor: AppTheme.primary.withAlpha(50),
      labelStyle: TextStyle(
        color: isSelected ? AppTheme.primary : AppTheme.textSecondary,
      ),
    );
  }
}
