import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:urja_flutter/l10n/app_localizations.dart';

import '../../config/theme.dart';
import '../../models/user.dart';
import '../../providers/providers.dart';

class ProfileScreen extends ConsumerStatefulWidget {
  const ProfileScreen({super.key});

  @override
  ConsumerState<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends ConsumerState<ProfileScreen> {
  final _nameController = TextEditingController();
  final _emailController = TextEditingController();
  final _emergencyNameController = TextEditingController();
  final _emergencyPhoneController = TextEditingController();
  String? _selectedGender;
  bool _editing = false;
  bool _saving = false;

  @override
  void dispose() {
    _nameController.dispose();
    _emailController.dispose();
    _emergencyNameController.dispose();
    _emergencyPhoneController.dispose();
    super.dispose();
  }

  void _loadProfile(Member member) {
    _nameController.text = member.name;
    _emailController.text = member.email ?? '';
    _emergencyNameController.text = member.emergencyContactName ?? '';
    _emergencyPhoneController.text = member.emergencyContactPhone ?? '';
    _selectedGender = member.gender;
  }

  Future<void> _saveProfile() async {
    setState(() => _saving = true);
    try {
      final update = ProfileUpdate(
        name: _nameController.text.trim().isNotEmpty
            ? _nameController.text.trim()
            : null,
        email: _emailController.text.trim().isNotEmpty
            ? _emailController.text.trim()
            : null,
        gender: _selectedGender,
        emergencyContactName:
            _emergencyNameController.text.trim().isNotEmpty
                ? _emergencyNameController.text.trim()
                : null,
        emergencyContactPhone:
            _emergencyPhoneController.text.trim().isNotEmpty
                ? _emergencyPhoneController.text.trim()
                : null,
      );

      await ref.read(profileProvider.notifier).updateProfile(update);

      if (mounted) {
        setState(() => _editing = false);
        final l10n = AppLocalizations.of(context)!;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(l10n.profileUpdated)),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  Future<void> _logout() async {
    final l10n = AppLocalizations.of(context)!;
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: UrjaColors.backgroundElevated,
        title: Text(l10n.logout),
        content: Text(l10n.logoutConfirm),
        actions: [
          TextButton(
              onPressed: () => Navigator.pop(context, false),
              child: Text(l10n.cancel)),
          ElevatedButton(
            onPressed: () => Navigator.pop(context, true),
            style: ElevatedButton.styleFrom(
                backgroundColor: UrjaColors.error),
            child: Text(l10n.logout),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      final api = ref.read(apiClientProvider);
      await api.logout();
      ref.read(authStateProvider.notifier).state =
          AuthState.unauthenticated;
      if (mounted) context.go('/login');
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final profileAsync = ref.watch(profileProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.profile),
        actions: [
          if (!_editing)
            IconButton(
              icon: const Icon(Icons.edit),
              onPressed: () {
                final member = profileAsync.valueOrNull;
                if (member != null) {
                  _loadProfile(member);
                  setState(() => _editing = true);
                }
              },
            ),
        ],
      ),
      body: profileAsync.when(
        data: (member) {
          if (member == null) {
            return const Center(child: CircularProgressIndicator());
          }

          if (_editing) {
            return _buildEditForm(context, l10n);
          }

          return _buildProfileView(context, member, l10n);
        },
        loading: () =>
            const Center(child: CircularProgressIndicator()),
        error: (_, _) => Center(
          child: Text(l10n.error),
        ),
      ),
    );
  }

  Widget _buildProfileView(
      BuildContext context, Member member, AppLocalizations l10n) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        // Avatar + Name
        Center(
          child: Column(
            children: [
              CircleAvatar(
                radius: 48,
                backgroundColor: UrjaColors.accent,
                backgroundImage: member.avatarUrl != null
                    ? NetworkImage(member.avatarUrl!)
                    : null,
                child: member.avatarUrl == null
                    ? Text(
                        member.name.isNotEmpty
                            ? member.name[0].toUpperCase()
                            : '?',
                        style: const TextStyle(
                            fontSize: 32, color: Colors.white),
                      )
                    : null,
              ),
              const SizedBox(height: 12),
              Text(member.displayName,
                  style: Theme.of(context).textTheme.headlineSmall),
              const SizedBox(height: 4),
              Text(member.phone,
                  style: Theme.of(context).textTheme.bodyMedium),
            ],
          ),
        ),
        const SizedBox(height: 24),

        _ProfileRow(label: l10n.email, value: member.email ?? '-'),
        _ProfileRow(label: l10n.gender, value: member.gender ?? '-'),
        _ProfileRow(
            label: l10n.dateOfBirth, value: member.dateOfBirth ?? '-'),
        _ProfileRow(
            label: l10n.emergencyContact,
            value: member.emergencyContactName ?? '-'),
        _ProfileRow(
            label: l10n.emergencyPhone,
            value: member.emergencyContactPhone ?? '-'),

        const SizedBox(height: 24),

        // Org memberships
        if (member.organizations.isNotEmpty) ...[
          Text('Organizations',
              style: Theme.of(context).textTheme.titleMedium),
          const SizedBox(height: 8),
          ...member.organizations.map((org) => Card(
                margin: const EdgeInsets.only(bottom: 8),
                child: ListTile(
                  leading: const Icon(Icons.fitness_center,
                      color: UrjaColors.accent),
                  title: Text(org.orgName),
                  subtitle: Text('${org.role} - ${org.status}'),
                ),
              )),
          const SizedBox(height: 16),
        ],

        // Logout
        OutlinedButton.icon(
          key: const Key('logout_button'),
          onPressed: _logout,
          icon: const Icon(Icons.logout, color: UrjaColors.error),
          label: Text(l10n.logout,
              style: const TextStyle(color: UrjaColors.error)),
          style: OutlinedButton.styleFrom(
            side: const BorderSide(color: UrjaColors.error),
          ),
        ),
      ],
    );
  }

  Widget _buildEditForm(BuildContext context, AppLocalizations l10n) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        TextFormField(
          key: const Key('name_field'),
          controller: _nameController,
          decoration: InputDecoration(labelText: l10n.name),
        ),
        const SizedBox(height: 12),
        TextFormField(
          key: const Key('email_field'),
          controller: _emailController,
          decoration: InputDecoration(labelText: l10n.email),
          keyboardType: TextInputType.emailAddress,
        ),
        const SizedBox(height: 12),
        DropdownButtonFormField<String>(
          initialValue: _selectedGender,
          decoration: InputDecoration(labelText: l10n.gender),
          dropdownColor: UrjaColors.backgroundElevated,
          items: [
            DropdownMenuItem(value: 'male', child: Text(l10n.male)),
            DropdownMenuItem(value: 'female', child: Text(l10n.female)),
            DropdownMenuItem(value: 'other', child: Text(l10n.other)),
            DropdownMenuItem(
                value: 'prefer_not_say',
                child: Text(l10n.preferNotSay)),
          ],
          onChanged: (v) => setState(() => _selectedGender = v),
        ),
        const SizedBox(height: 12),
        TextFormField(
          controller: _emergencyNameController,
          decoration: InputDecoration(labelText: l10n.emergencyContact),
        ),
        const SizedBox(height: 12),
        TextFormField(
          controller: _emergencyPhoneController,
          decoration: InputDecoration(labelText: l10n.emergencyPhone),
          keyboardType: TextInputType.phone,
        ),
        const SizedBox(height: 24),
        Row(
          children: [
            Expanded(
              child: OutlinedButton(
                onPressed: () => setState(() => _editing = false),
                child: Text(l10n.cancel),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: ElevatedButton(
                key: const Key('save_profile_button'),
                onPressed: _saving ? null : _saveProfile,
                child: _saving
                    ? const SizedBox(
                        height: 20,
                        width: 20,
                        child: CircularProgressIndicator(
                            strokeWidth: 2, color: Colors.white),
                      )
                    : Text(l10n.saveChanges),
              ),
            ),
          ],
        ),
      ],
    );
  }
}

class _ProfileRow extends StatelessWidget {
  final String label;
  final String value;

  const _ProfileRow({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 140,
            child: Text(label,
                style: Theme.of(context).textTheme.bodyMedium),
          ),
          Expanded(
            child: Text(value,
                style: Theme.of(context)
                    .textTheme
                    .bodyMedium
                    ?.copyWith(color: UrjaColors.foreground)),
          ),
        ],
      ),
    );
  }
}
