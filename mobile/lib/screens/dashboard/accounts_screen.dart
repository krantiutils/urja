import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../services/finance_service.dart';
import '../../models/transaction.dart';
import '../shared/widgets/stat_card.dart';
import '../shared/widgets/empty_state.dart';
import '../shared/widgets/loading_indicator.dart';
import '../shared/widgets/confirm_dialog.dart';

class AccountsScreen extends ConsumerStatefulWidget {
  const AccountsScreen({super.key});

  @override
  ConsumerState<AccountsScreen> createState() => _AccountsScreenState();
}

class _AccountsScreenState extends ConsumerState<AccountsScreen> {
  bool _loading = true;
  List<Transaction> _transactions = [];
  AccountsSummary? _summary;
  String _filter = 'all'; // all, income, expense

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) {
      setState(() => _loading = false);
      return;
    }

    setState(() => _loading = true);
    final financeService = FinanceService(ref.read(apiClientProvider));

    try {
      final results = await Future.wait([
        financeService.getTransactions(
          orgId,
          transactionType: _filter != 'all' ? _filter : null,
        ),
        financeService.getSummary(orgId),
      ]);

      final txResult = results[0] as ({List<Transaction> data, int total});
      final summary = results[1] as AccountsSummary;

      if (mounted) {
        setState(() {
          _transactions = txResult.data;
          _summary = summary;
          _loading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _loading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load accounts: $e')),
        );
      }
    }
  }

  Future<void> _showTransactionDialog({Transaction? existing}) async {
    final l10n = AppLocalizations.of(context)!;
    final categoryCtrl =
        TextEditingController(text: existing?.category ?? '');
    final descCtrl =
        TextEditingController(text: existing?.description ?? '');
    final amountCtrl = TextEditingController(text: existing?.amount ?? '');
    final referenceCtrl =
        TextEditingController(text: existing?.reference ?? '');
    String txType = existing?.transactionType ?? 'income';
    String paymentType = existing?.paymentType ?? 'cash';
    DateTime txDate = existing != null
        ? (DateTime.tryParse(existing.transactionDate) ?? DateTime.now())
        : DateTime.now();

    final result = await showDialog<bool>(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setDialogState) => AlertDialog(
          title: Text(existing != null
              ? l10n.editTransaction
              : l10n.addTransaction),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                DropdownButtonFormField<String>(
                  initialValue: txType,
                  decoration: InputDecoration(labelText: l10n.type),
                  items: ['income', 'expense']
                      .map((t) => DropdownMenuItem(value: t, child: Text(t)))
                      .toList(),
                  onChanged: (v) {
                    if (v != null) setDialogState(() => txType = v);
                  },
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: categoryCtrl,
                  decoration: InputDecoration(labelText: l10n.category),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: descCtrl,
                  decoration: InputDecoration(labelText: l10n.description),
                  maxLines: 2,
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: amountCtrl,
                  decoration: InputDecoration(
                    labelText: l10n.amount,
                    prefixText: 'Rs. ',
                  ),
                  keyboardType: kIsWeb ? TextInputType.text : TextInputType.number,
                ),
                const SizedBox(height: 12),
                DropdownButtonFormField<String>(
                  initialValue: paymentType,
                  decoration: InputDecoration(labelText: l10n.paymentMethod),
                  items: ['cash', 'esewa', 'bank', 'other']
                      .map((m) => DropdownMenuItem(value: m, child: Text(m)))
                      .toList(),
                  onChanged: (v) {
                    if (v != null) setDialogState(() => paymentType = v);
                  },
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: referenceCtrl,
                  decoration: InputDecoration(labelText: l10n.reference),
                ),
                const SizedBox(height: 12),
                ListTile(
                  contentPadding: EdgeInsets.zero,
                  title: Text(l10n.date),
                  subtitle: Text(
                    '${txDate.year}-${txDate.month.toString().padLeft(2, '0')}-${txDate.day.toString().padLeft(2, '0')}',
                  ),
                  trailing: const Icon(Icons.calendar_today),
                  onTap: () async {
                    final picked = await showDatePicker(
                      context: context,
                      initialDate: txDate,
                      firstDate: DateTime(2020),
                      lastDate: DateTime.now(),
                    );
                    if (picked != null) {
                      setDialogState(() => txDate = picked);
                    }
                  },
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
              child: Text(l10n.save),
            ),
          ],
        ),
      ),
    );

    if (result == true && mounted) {
      final orgId = ref.read(authProvider).user?.orgId;
      if (orgId == null) return;

      final financeService = FinanceService(ref.read(apiClientProvider));
      final data = {
        'category': categoryCtrl.text.trim(),
        'description': descCtrl.text.trim(),
        'amount': amountCtrl.text.trim(),
        'transaction_type': txType,
        'payment_type': paymentType,
        'reference': referenceCtrl.text.trim(),
        'transaction_date':
            '${txDate.year}-${txDate.month.toString().padLeft(2, '0')}-${txDate.day.toString().padLeft(2, '0')}',
      };

      try {
        if (existing != null) {
          await financeService.updateTransaction(orgId, existing.id, data);
        } else {
          await financeService.addTransaction(orgId, data);
        }
        _loadData();
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to save transaction: $e')),
          );
        }
      }
    }
  }

  Future<void> _deleteTransaction(Transaction tx) async {
    final confirmed = await showConfirmDialog(
      context,
      'Delete transaction "${tx.description}"?',
    );

    if (confirmed && mounted) {
      final orgId = ref.read(authProvider).user?.orgId;
      if (orgId == null) return;

      final financeService = FinanceService(ref.read(apiClientProvider));
      try {
        await financeService.deleteTransaction(orgId, tx.id);
        _loadData();
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to delete transaction: $e')),
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
        title: Text(l10n.accounts),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => _showTransactionDialog(),
          ),
        ],
      ),
      body: _loading
          ? const LoadingIndicator()
          : RefreshIndicator(
              onRefresh: _loadData,
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  // Summary cards
                  if (_summary != null) _buildSummary(l10n),
                  const SizedBox(height: 16),

                  // Filter chips
                  Row(
                    children: [
                      _buildFilterChip(l10n.all, 'all'),
                      const SizedBox(width: 8),
                      _buildFilterChip(l10n.income, 'income'),
                      const SizedBox(width: 8),
                      _buildFilterChip(l10n.expense, 'expense'),
                    ],
                  ),
                  const SizedBox(height: 16),

                  // Transactions list
                  if (_transactions.isEmpty)
                    EmptyState(
                      icon: Icons.account_balance_wallet_outlined,
                      message: l10n.noTransactions,
                    )
                  else
                    ...List.generate(_transactions.length, (index) {
                      final tx = _transactions[index];
                      final isIncome = tx.transactionType == 'income';
                      final date = DateTime.tryParse(tx.transactionDate);
                      final dateStr = date != null
                          ? '${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')}'
                          : tx.transactionDate;

                      return Card(
                        margin: const EdgeInsets.only(bottom: 8),
                        child: ListTile(
                          leading: CircleAvatar(
                            backgroundColor: isIncome
                                ? AppTheme.primary.withAlpha(30)
                                : AppTheme.error.withAlpha(30),
                            child: Icon(
                              isIncome
                                  ? Icons.arrow_downward
                                  : Icons.arrow_upward,
                              color:
                                  isIncome ? AppTheme.primary : AppTheme.error,
                            ),
                          ),
                          title: Text(tx.category.isNotEmpty
                              ? tx.category
                              : tx.description),
                          subtitle: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              if (tx.description.isNotEmpty &&
                                  tx.category.isNotEmpty)
                                Text(
                                  tx.description,
                                  style: const TextStyle(
                                      color: AppTheme.textSecondary,
                                      fontSize: 12),
                                  maxLines: 1,
                                  overflow: TextOverflow.ellipsis,
                                ),
                              Text(
                                dateStr,
                                style: const TextStyle(
                                    color: AppTheme.textSecondary,
                                    fontSize: 12),
                              ),
                            ],
                          ),
                          trailing: Text(
                            '${isIncome ? '+' : '-'} Rs. ${tx.amount}',
                            style: TextStyle(
                              color:
                                  isIncome ? AppTheme.primary : AppTheme.error,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          isThreeLine:
                              tx.description.isNotEmpty &&
                              tx.category.isNotEmpty,
                          onTap: () =>
                              _showTransactionDialog(existing: tx),
                          onLongPress: () => _deleteTransaction(tx),
                        ),
                      );
                    }),
                ],
              ),
            ),
    );
  }

  Widget _buildSummary(AppLocalizations l10n) {
    final s = _summary!;
    return GridView.count(
      crossAxisCount: 2,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      mainAxisSpacing: 12,
      crossAxisSpacing: 12,
      childAspectRatio: 1.4,
      children: [
        StatCard(
          title: l10n.totalIncome,
          value: 'Rs. ${s.totalIncome}',
          icon: Icons.trending_up,
          color: AppTheme.primary,
        ),
        StatCard(
          title: l10n.totalExpenses,
          value: 'Rs. ${s.totalExpenses}',
          icon: Icons.trending_down,
          color: AppTheme.error,
        ),
        StatCard(
          title: l10n.grossProfit,
          value: 'Rs. ${s.grossProfit}',
          icon: Icons.account_balance,
          color: AppTheme.info,
        ),
        StatCard(
          title: l10n.profitPercent,
          value: '${s.profitPercent.toStringAsFixed(1)}%',
          icon: Icons.pie_chart,
          color: AppTheme.warning,
        ),
      ],
    );
  }

  Widget _buildFilterChip(String label, String value) {
    final isSelected = _filter == value;
    return ChoiceChip(
      label: Text(label),
      selected: isSelected,
      onSelected: (_) {
        setState(() => _filter = value);
        _loadData();
      },
      selectedColor: AppTheme.primary.withAlpha(50),
      labelStyle: TextStyle(
        color: isSelected ? AppTheme.primary : AppTheme.textSecondary,
      ),
    );
  }
}
