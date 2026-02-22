import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../providers/locale_provider.dart';
import '../../services/org_service.dart';
import '../../services/member_service.dart';
import '../../models/member.dart';
import '../shared/widgets/loading_indicator.dart';

class SettingsScreen extends ConsumerStatefulWidget {
  const SettingsScreen({super.key});

  @override
  ConsumerState<SettingsScreen> createState() => _SettingsScreenState();
}

class _SettingsScreenState extends ConsumerState<SettingsScreen> {
  bool _loading = true;
  MemberProfile? _profile;

  // Org editing controllers
  final _orgNameCtrl = TextEditingController();
  final _orgNameNeCtrl = TextEditingController();
  final _orgDescCtrl = TextEditingController();
  final _orgAddressCtrl = TextEditingController();
  final _orgPhoneCtrl = TextEditingController();
  final _orgEmailCtrl = TextEditingController();

  bool _savingOrg = false;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  @override
  void dispose() {
    _orgNameCtrl.dispose();
    _orgNameNeCtrl.dispose();
    _orgDescCtrl.dispose();
    _orgAddressCtrl.dispose();
    _orgPhoneCtrl.dispose();
    _orgEmailCtrl.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) return;

    setState(() => _loading = true);
    final orgService = OrgService(ref.read(apiClientProvider));
    final memberService = MemberService(ref.read(apiClientProvider));

    try {
      final results = await Future.wait([
        orgService.getOrg(orgId),
        memberService.getProfile(),
      ]);

      final org = results[0] as Organization;
      final profile = results[1] as MemberProfile;

      if (mounted) {
        setState(() {
          _profile = profile;
          _orgNameCtrl.text = org.name;
          _orgNameNeCtrl.text = org.nameNe;
          _orgDescCtrl.text = org.description;
          _orgAddressCtrl.text = org.address;
          _orgPhoneCtrl.text = org.phone;
          _orgEmailCtrl.text = org.email;
          _loading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _loading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load settings: $e')),
        );
      }
    }
  }

  Future<void> _saveOrgSettings() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) return;

    setState(() => _savingOrg = true);
    final orgService = OrgService(ref.read(apiClientProvider));

    try {
      await orgService.updateOrg(orgId, {
        'name': _orgNameCtrl.text.trim(),
        'name_ne': _orgNameNeCtrl.text.trim(),
        'description': _orgDescCtrl.text.trim(),
        'address': _orgAddressCtrl.text.trim(),
        'phone': _orgPhoneCtrl.text.trim(),
        'email': _orgEmailCtrl.text.trim(),
      });
      if (mounted) {
        final l10n = AppLocalizations.of(context)!;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(l10n.settingsSaved),
            backgroundColor: AppTheme.primary,
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to save settings: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _savingOrg = false);
    }
  }

  Future<void> _logout() async {
    final l10n = AppLocalizations.of(context)!;
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.logout),
        content: Text(l10n.logoutConfirm),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: Text(l10n.cancel),
          ),
          TextButton(
            onPressed: () => Navigator.pop(context, true),
            style: TextButton.styleFrom(foregroundColor: AppTheme.error),
            child: Text(l10n.logout),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      await ref.read(authProvider.notifier).logout();
      if (mounted) context.go('/login');
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final locale = ref.watch(localeProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.settings),
      ),
      body: _loading
          ? const LoadingIndicator()
          : ListView(
              padding: const EdgeInsets.all(16),
              children: [
                // Organization profile section
                _buildSectionHeader(l10n.organizationProfile),
                const SizedBox(height: 12),
                _buildOrgSection(l10n),
                const SizedBox(height: 24),

                // Personal profile section
                _buildSectionHeader(l10n.personalProfile),
                const SizedBox(height: 12),
                _buildPersonalSection(l10n),
                const SizedBox(height: 24),

                // Language toggle
                _buildSectionHeader(l10n.language),
                const SizedBox(height: 12),
                Card(
                  child: ListTile(
                    leading: const Icon(Icons.language),
                    title: Text(l10n.language),
                    subtitle: Text(
                        locale.languageCode == 'en' ? 'English' : 'Nepali'),
                    trailing: Switch(
                      value: locale.languageCode == 'ne',
                      onChanged: (_) =>
                          ref.read(localeProvider.notifier).toggle(),
                      activeThumbColor: AppTheme.primary,
                    ),
                  ),
                ),
                const SizedBox(height: 24),

                // Logout button
                SizedBox(
                  width: double.infinity,
                  child: OutlinedButton.icon(
                    onPressed: _logout,
                    icon: const Icon(Icons.logout),
                    label: Text(l10n.logout),
                    style: OutlinedButton.styleFrom(
                      foregroundColor: AppTheme.error,
                      side: const BorderSide(color: AppTheme.error),
                      padding: const EdgeInsets.symmetric(vertical: 14),
                    ),
                  ),
                ),
                const SizedBox(height: 32),
              ],
            ),
    );
  }

  Widget _buildSectionHeader(String title) {
    return Text(
      title,
      style: const TextStyle(
        fontSize: 18,
        fontWeight: FontWeight.bold,
      ),
    );
  }

  Widget _buildOrgSection(AppLocalizations l10n) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            TextField(
              controller: _orgNameCtrl,
              decoration: InputDecoration(labelText: l10n.name),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _orgNameNeCtrl,
              decoration: InputDecoration(labelText: l10n.nameNe),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _orgDescCtrl,
              decoration: InputDecoration(labelText: l10n.description),
              maxLines: 2,
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _orgAddressCtrl,
              decoration: InputDecoration(labelText: l10n.address),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _orgPhoneCtrl,
              decoration: InputDecoration(labelText: l10n.phone),
              keyboardType: TextInputType.phone,
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _orgEmailCtrl,
              decoration: InputDecoration(labelText: l10n.email),
              keyboardType: TextInputType.emailAddress,
            ),
            const SizedBox(height: 16),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: _savingOrg ? null : _saveOrgSettings,
                child: _savingOrg
                    ? const SizedBox(
                        width: 18,
                        height: 18,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          color: Colors.white,
                        ),
                      )
                    : Text(l10n.save),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildPersonalSection(AppLocalizations l10n) {
    if (_profile == null) {
      return Card(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Center(
            child: Text(
              l10n.noProfile,
              style: const TextStyle(color: AppTheme.textSecondary),
            ),
          ),
        ),
      );
    }

    final p = _profile!;
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            CircleAvatar(
              radius: 40,
              backgroundColor: AppTheme.primary,
              backgroundImage:
                  p.avatarUrl != null ? NetworkImage(p.avatarUrl!) : null,
              child: p.avatarUrl == null
                  ? Text(
                      p.name.isNotEmpty ? p.name[0].toUpperCase() : '?',
                      style:
                          const TextStyle(fontSize: 32, color: Colors.white),
                    )
                  : null,
            ),
            const SizedBox(height: 16),
            _buildProfileRow(l10n.name, p.name),
            _buildProfileRow(l10n.phone, p.phone),
            if (p.email != null) _buildProfileRow(l10n.email, p.email!),
            if (p.gender != null) _buildProfileRow(l10n.gender, p.gender!),
            if (p.dateOfBirth != null)
              _buildProfileRow(l10n.dateOfBirth, p.dateOfBirth!),
          ],
        ),
      ),
    );
  }

  Widget _buildProfileRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6),
      child: Row(
        children: [
          SizedBox(
            width: 100,
            child: Text(
              label,
              style: const TextStyle(color: AppTheme.textSecondary),
            ),
          ),
          Expanded(
            child: Text(value, style: const TextStyle(fontSize: 15)),
          ),
        ],
      ),
    );
  }
}
