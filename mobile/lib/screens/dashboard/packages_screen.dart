import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../services/org_service.dart';
import '../../models/package.dart';
import '../shared/widgets/status_badge.dart';
import '../shared/widgets/empty_state.dart';
import '../shared/widgets/loading_indicator.dart';
import '../shared/widgets/confirm_dialog.dart';

class PackagesScreen extends ConsumerStatefulWidget {
  const PackagesScreen({super.key});

  @override
  ConsumerState<PackagesScreen> createState() => _PackagesScreenState();
}

class _PackagesScreenState extends ConsumerState<PackagesScreen> {
  bool _loading = true;
  List<OrgPackage> _packages = [];
  String _filter = 'all'; // all, active, inactive

  @override
  void initState() {
    super.initState();
    _loadPackages();
  }

  Future<void> _loadPackages() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) return;

    setState(() => _loading = true);
    final orgService = OrgService(ref.read(apiClientProvider));

    try {
      final packages = await orgService.getPackages(orgId, limit: 100);
      if (mounted) {
        setState(() {
          _packages = packages;
          _loading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _loading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load packages: $e')),
        );
      }
    }
  }

  List<OrgPackage> get _filteredPackages {
    if (_filter == 'active') {
      return _packages.where((p) => p.isActive).toList();
    } else if (_filter == 'inactive') {
      return _packages.where((p) => !p.isActive).toList();
    }
    return _packages;
  }

  Future<void> _showPackageDialog({OrgPackage? existing}) async {
    final l10n = AppLocalizations.of(context)!;
    final nameCtrl = TextEditingController(text: existing?.name ?? '');
    final nameNeCtrl = TextEditingController(text: existing?.nameNe ?? '');
    final descCtrl = TextEditingController(text: existing?.description ?? '');
    final priceCtrl = TextEditingController(text: existing?.price ?? '');
    final durationCtrl = TextEditingController(
        text: existing?.durationDays.toString() ?? '');
    final featuresCtrl = TextEditingController(
        text: existing?.features?.join(', ') ?? '');
    bool isActive = existing?.isActive ?? true;

    final result = await showDialog<bool>(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setDialogState) => AlertDialog(
          title: Text(existing != null ? l10n.editPackage : l10n.addPackage),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                TextField(
                  controller: nameCtrl,
                  decoration: InputDecoration(labelText: l10n.name),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: nameNeCtrl,
                  decoration: InputDecoration(labelText: l10n.nameNe),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: descCtrl,
                  decoration: InputDecoration(labelText: l10n.description),
                  maxLines: 2,
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: priceCtrl,
                  decoration: InputDecoration(
                    labelText: l10n.price,
                    prefixText: 'Rs. ',
                  ),
                  keyboardType: TextInputType.number,
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: durationCtrl,
                  decoration: InputDecoration(labelText: l10n.durationDays),
                  keyboardType: TextInputType.number,
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: featuresCtrl,
                  decoration: InputDecoration(
                    labelText: l10n.features,
                    hintText: l10n.featuresHint,
                  ),
                  maxLines: 2,
                ),
                const SizedBox(height: 12),
                SwitchListTile(
                  title: Text(l10n.active),
                  value: isActive,
                  onChanged: (v) => setDialogState(() => isActive = v),
                  contentPadding: EdgeInsets.zero,
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

      final orgService = OrgService(ref.read(apiClientProvider));
      final data = {
        'name': nameCtrl.text.trim(),
        'name_ne': nameNeCtrl.text.trim(),
        'description': descCtrl.text.trim(),
        'price': priceCtrl.text.trim(),
        'duration_days': int.tryParse(durationCtrl.text.trim()) ?? 30,
        'features': featuresCtrl.text
            .split(',')
            .map((f) => f.trim())
            .where((f) => f.isNotEmpty)
            .toList(),
        'is_active': isActive,
      };

      try {
        if (existing != null) {
          await orgService.updatePackage(orgId, existing.id, data);
        } else {
          await orgService.addPackage(orgId, data);
        }
        _loadPackages();
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to save package: $e')),
          );
        }
      }
    }
  }

  Future<void> _deletePackage(OrgPackage pkg) async {
    final confirmed = await showConfirmDialog(
      context,
      'Delete package "${pkg.name}"?',
    );

    if (confirmed && mounted) {
      final orgId = ref.read(authProvider).user?.orgId;
      if (orgId == null) return;

      final orgService = OrgService(ref.read(apiClientProvider));
      try {
        await orgService.deletePackage(orgId, pkg.id);
        _loadPackages();
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to delete package: $e')),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final filtered = _filteredPackages;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.packages),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => _showPackageDialog(),
          ),
        ],
      ),
      body: _loading
          ? const LoadingIndicator()
          : Column(
              children: [
                // Filter chips
                Padding(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                  child: Row(
                    children: [
                      _buildFilterChip(l10n.all, 'all'),
                      const SizedBox(width: 8),
                      _buildFilterChip(l10n.active, 'active'),
                      const SizedBox(width: 8),
                      _buildFilterChip(l10n.inactive, 'inactive'),
                    ],
                  ),
                ),
                if (filtered.isEmpty)
                  Expanded(
                    child: EmptyState(
                      icon: Icons.inventory_2_outlined,
                      message: l10n.noPackages,
                    ),
                  )
                else
                  Expanded(
                    child: RefreshIndicator(
                      onRefresh: _loadPackages,
                      child: GridView.builder(
                        padding: const EdgeInsets.all(16),
                        gridDelegate:
                            const SliverGridDelegateWithFixedCrossAxisCount(
                          crossAxisCount: 2,
                          mainAxisSpacing: 12,
                          crossAxisSpacing: 12,
                          childAspectRatio: 0.75,
                        ),
                        itemCount: filtered.length,
                        itemBuilder: (context, index) {
                          final pkg = filtered[index];
                          return _buildPackageCard(pkg, l10n);
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
      onSelected: (_) => setState(() => _filter = value),
      selectedColor: AppTheme.primary.withAlpha(50),
      labelStyle: TextStyle(
        color: isSelected ? AppTheme.primary : AppTheme.textSecondary,
      ),
    );
  }

  Widget _buildPackageCard(OrgPackage pkg, AppLocalizations l10n) {
    return GestureDetector(
      onTap: () => _showPackageDialog(existing: pkg),
      onLongPress: () => _deletePackage(pkg),
      child: Card(
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Expanded(
                    child: Text(
                      pkg.name,
                      style: const TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 15,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                  StatusBadge(label: pkg.isActive ? l10n.active : l10n.inactive),
                ],
              ),
              const SizedBox(height: 8),
              Text(
                'Rs. ${pkg.price}',
                style: const TextStyle(
                  color: AppTheme.primary,
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                '${pkg.durationDays} ${l10n.days}',
                style: const TextStyle(color: AppTheme.textSecondary),
              ),
              const SizedBox(height: 8),
              if (pkg.features != null && pkg.features!.isNotEmpty)
                Expanded(
                  child: ListView(
                    shrinkWrap: true,
                    physics: const NeverScrollableScrollPhysics(),
                    children: pkg.features!
                        .take(3)
                        .map((f) => Padding(
                              padding: const EdgeInsets.only(bottom: 2),
                              child: Row(
                                children: [
                                  const Icon(Icons.check,
                                      size: 14, color: AppTheme.primary),
                                  const SizedBox(width: 4),
                                  Expanded(
                                    child: Text(
                                      f,
                                      style: const TextStyle(fontSize: 12),
                                      maxLines: 1,
                                      overflow: TextOverflow.ellipsis,
                                    ),
                                  ),
                                ],
                              ),
                            ))
                        .toList(),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
