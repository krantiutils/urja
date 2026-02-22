import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../services/staff_service.dart';
import '../../models/staff.dart';
import '../shared/widgets/search_field.dart';
import '../shared/widgets/status_badge.dart';
import '../shared/widgets/empty_state.dart';
import '../shared/widgets/loading_indicator.dart';
import '../shared/widgets/confirm_dialog.dart';

class StaffScreen extends ConsumerStatefulWidget {
  const StaffScreen({super.key});

  @override
  ConsumerState<StaffScreen> createState() => _StaffScreenState();
}

class _StaffScreenState extends ConsumerState<StaffScreen> {
  bool _loading = true;
  List<StaffMember> _staff = [];
  String _search = '';

  static const _staffRoles = ['owner', 'manager', 'trainer', 'receptionist'];

  @override
  void initState() {
    super.initState();
    _loadStaff();
  }

  Future<void> _loadStaff() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) return;

    setState(() => _loading = true);
    final staffService = StaffService(ref.read(apiClientProvider));

    try {
      final result = await staffService.getStaff(
        orgId,
        search: _search.isNotEmpty ? _search : null,
      );
      if (mounted) {
        setState(() {
          _staff = result.data;
          _loading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _loading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load staff: $e')),
        );
      }
    }
  }

  void _onSearchChanged(String value) {
    _search = value;
    _loadStaff();
  }

  Future<void> _showStaffDialog({StaffMember? existing}) async {
    final l10n = AppLocalizations.of(context)!;
    final nameCtrl = TextEditingController(text: existing?.name ?? '');
    final phoneCtrl = TextEditingController(text: existing?.phone ?? '');
    final emailCtrl = TextEditingController(text: existing?.email ?? '');
    String selectedRole = existing?.staffRole ?? 'trainer';

    final result = await showDialog<bool>(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setDialogState) => AlertDialog(
          title:
              Text(existing != null ? l10n.editStaff : l10n.addStaff),
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
                  controller: phoneCtrl,
                  decoration: InputDecoration(labelText: l10n.phone),
                  keyboardType: TextInputType.phone,
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: emailCtrl,
                  decoration: InputDecoration(labelText: l10n.email),
                  keyboardType: TextInputType.emailAddress,
                ),
                const SizedBox(height: 12),
                DropdownButtonFormField<String>(
                  initialValue: selectedRole,
                  decoration: InputDecoration(labelText: l10n.staffRole),
                  items: _staffRoles
                      .map((r) => DropdownMenuItem(value: r, child: Text(r)))
                      .toList(),
                  onChanged: (v) {
                    if (v != null) {
                      setDialogState(() => selectedRole = v);
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

      final staffService = StaffService(ref.read(apiClientProvider));
      final data = {
        'name': nameCtrl.text.trim(),
        'phone': phoneCtrl.text.trim(),
        if (emailCtrl.text.trim().isNotEmpty) 'email': emailCtrl.text.trim(),
        'staff_role': selectedRole,
      };

      try {
        if (existing != null) {
          await staffService.updateStaff(orgId, existing.id, data);
        } else {
          await staffService.addStaff(orgId, data);
        }
        _loadStaff();
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to save staff: $e')),
          );
        }
      }
    }
  }

  Future<void> _deleteStaff(StaffMember staff) async {
    final confirmed = await showConfirmDialog(
      context,
      'Remove staff member "${staff.name}"?',
    );

    if (confirmed && mounted) {
      final orgId = ref.read(authProvider).user?.orgId;
      if (orgId == null) return;

      final staffService = StaffService(ref.read(apiClientProvider));
      try {
        await staffService.deleteStaff(orgId, staff.id);
        _loadStaff();
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to remove staff: $e')),
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
        title: Text(l10n.staff),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => _showStaffDialog(),
          ),
        ],
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: SearchField(onChanged: _onSearchChanged),
          ),
          if (_loading)
            const Expanded(child: LoadingIndicator())
          else if (_staff.isEmpty)
            Expanded(
              child: EmptyState(
                icon: Icons.badge_outlined,
                message: l10n.noStaff,
              ),
            )
          else
            Expanded(
              child: RefreshIndicator(
                onRefresh: _loadStaff,
                child: ListView.builder(
                  itemCount: _staff.length,
                  itemBuilder: (context, index) {
                    final staff = _staff[index];
                    return Card(
                      margin: const EdgeInsets.symmetric(
                          horizontal: 16, vertical: 4),
                      child: ListTile(
                        leading: CircleAvatar(
                          backgroundColor: AppTheme.info,
                          backgroundImage: staff.avatarUrl != null
                              ? NetworkImage(staff.avatarUrl!)
                              : null,
                          child: staff.avatarUrl == null
                              ? Text(
                                  staff.name.isNotEmpty
                                      ? staff.name[0].toUpperCase()
                                      : '?',
                                  style: const TextStyle(color: Colors.white),
                                )
                              : null,
                        ),
                        title: Text(staff.name),
                        subtitle: Text(
                          '${staff.phone} - ${staff.staffRole}',
                          style:
                              const TextStyle(color: AppTheme.textSecondary),
                        ),
                        trailing: StatusBadge(label: staff.status),
                        onTap: () => _showStaffDialog(existing: staff),
                        onLongPress: () => _deleteStaff(staff),
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
}
