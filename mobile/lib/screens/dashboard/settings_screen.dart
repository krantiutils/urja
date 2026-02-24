import 'dart:typed_data';

import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
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

  // Personal profile controllers
  final _profileNameCtrl = TextEditingController();
  final _profileEmailCtrl = TextEditingController();
  final _profileGenderCtrl = TextEditingController();
  final _profileDobCtrl = TextEditingController();
  bool _savingProfile = false;

  // Org slug
  String _orgSlug = '';

  // QR code state
  Uint8List? _qrBytes;
  bool _qrLoading = true;
  String? _qrError;

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
    _profileNameCtrl.dispose();
    _profileEmailCtrl.dispose();
    _profileGenderCtrl.dispose();
    _profileDobCtrl.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) {
      setState(() => _loading = false);
      return;
    }

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
          _orgSlug = org.slug;
          _orgNameCtrl.text = org.name;
          _orgNameNeCtrl.text = org.nameNe;
          _orgDescCtrl.text = org.description;
          _orgAddressCtrl.text = org.address;
          _orgPhoneCtrl.text = org.phone;
          _orgEmailCtrl.text = org.email;
          // Populate personal profile controllers
          _profileNameCtrl.text = profile.name;
          _profileEmailCtrl.text = profile.email ?? '';
          _profileGenderCtrl.text = profile.gender ?? '';
          _profileDobCtrl.text = profile.dateOfBirth ?? '';
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

    // Load QR code in parallel (non-blocking)
    _loadQrCode();
  }

  Future<void> _loadQrCode() async {
    final orgId = ref.read(authProvider).user?.orgId;
    if (orgId == null) {
      setState(() {
        _qrLoading = false;
        _qrError = 'No organization';
      });
      return;
    }

    setState(() {
      _qrLoading = true;
      _qrError = null;
    });

    try {
      final apiClient = ref.read(apiClientProvider);
      // Fetch the QR code PNG directly using the org endpoint path
      final response = await apiClient.dio.get<List<int>>(
        '/orgs/$orgId/qr-code',
        options: Options(responseType: ResponseType.bytes),
      );

      if (mounted) {
        setState(() {
          _qrBytes = Uint8List.fromList(response.data!);
          _qrLoading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _qrLoading = false;
          _qrError = e.toString();
        });
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

  Future<void> _saveProfile() async {
    setState(() => _savingProfile = true);
    final memberService = MemberService(ref.read(apiClientProvider));

    try {
      final data = <String, dynamic>{
        'name': _profileNameCtrl.text.trim(),
      };
      final email = _profileEmailCtrl.text.trim();
      if (email.isNotEmpty) data['email'] = email;
      final gender = _profileGenderCtrl.text.trim();
      if (gender.isNotEmpty) data['gender'] = gender;
      final dob = _profileDobCtrl.text.trim();
      if (dob.isNotEmpty) data['date_of_birth'] = dob;

      await memberService.updateProfile(data);
      if (mounted) {
        final l10n = AppLocalizations.of(context)!;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(l10n.saved),
            backgroundColor: AppTheme.primary,
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to save profile: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _savingProfile = false);
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

                // Gym QR Code section
                _buildSectionHeader(l10n.gymCode),
                const SizedBox(height: 12),
                _buildQrCodeSection(l10n),
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

  Widget _buildQrCodeSection(AppLocalizations l10n) {
    final gymName = ref.read(authProvider).user?.orgName ?? _orgNameCtrl.text;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          children: [
            // Gym icon badge
            Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                gradient: const LinearGradient(
                  colors: [AppTheme.primary, AppTheme.primaryDark],
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                ),
                borderRadius: BorderRadius.circular(14),
              ),
              child: const Icon(
                Icons.fitness_center,
                color: Colors.white,
                size: 26,
              ),
            ),
            const SizedBox(height: 12),

            // Gym name
            if (gymName.isNotEmpty)
              Text(
                gymName,
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                ),
                textAlign: TextAlign.center,
              ),
            const SizedBox(height: 16),

            // QR code image
            if (_qrLoading)
              Container(
                width: 200,
                height: 200,
                decoration: BoxDecoration(
                  color: AppTheme.surfaceLight,
                  borderRadius: BorderRadius.circular(16),
                ),
                child: const Center(
                  child: SizedBox(
                    width: 32,
                    height: 32,
                    child: CircularProgressIndicator(
                      strokeWidth: 3,
                      color: AppTheme.primary,
                    ),
                  ),
                ),
              )
            else if (_qrError != null)
              Container(
                width: 200,
                height: 200,
                decoration: BoxDecoration(
                  color: AppTheme.surfaceLight,
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Icon(
                      Icons.error_outline,
                      color: AppTheme.error,
                      size: 36,
                    ),
                    const SizedBox(height: 8),
                    Text(
                      l10n.gymQrCodeError,
                      style: const TextStyle(
                        color: AppTheme.textSecondary,
                        fontSize: 13,
                      ),
                      textAlign: TextAlign.center,
                    ),
                    const SizedBox(height: 12),
                    TextButton.icon(
                      onPressed: _loadQrCode,
                      icon: const Icon(Icons.refresh, size: 16),
                      label: Text(l10n.retry),
                      style: TextButton.styleFrom(
                        foregroundColor: AppTheme.primary,
                      ),
                    ),
                  ],
                ),
              )
            else if (_qrBytes != null)
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Image.memory(
                  _qrBytes!,
                  width: 180,
                  height: 180,
                  fit: BoxFit.contain,
                ),
              ),

            const SizedBox(height: 16),

            // Gym code slug (copyable)
            if (_orgSlug.isNotEmpty) ...[
              Text(
                l10n.gymCode.toUpperCase(),
                style: const TextStyle(
                  fontSize: 11,
                  fontWeight: FontWeight.w600,
                  color: AppTheme.textSecondary,
                  letterSpacing: 1.5,
                ),
              ),
              const SizedBox(height: 6),
              GestureDetector(
                onTap: () {
                  Clipboard.setData(ClipboardData(text: _orgSlug));
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text('${l10n.gymCode} copied!'),
                      duration: const Duration(seconds: 2),
                    ),
                  );
                },
                child: Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 16,
                    vertical: 10,
                  ),
                  decoration: BoxDecoration(
                    color: AppTheme.surfaceLight,
                    borderRadius: BorderRadius.circular(10),
                    border: Border.all(
                      color: AppTheme.primary.withValues(alpha: 0.3),
                    ),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        _orgSlug,
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w700,
                          color: AppTheme.primary,
                          letterSpacing: 0.5,
                        ),
                      ),
                      const SizedBox(width: 8),
                      Icon(
                        Icons.copy,
                        size: 16,
                        color: AppTheme.primary.withValues(alpha: 0.6),
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 16),
            ],

            // "Scan to check in" label
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(
                  Icons.smartphone,
                  size: 16,
                  color: AppTheme.textSecondary,
                ),
                const SizedBox(width: 6),
                Text(
                  l10n.scanToCheckIn.toUpperCase(),
                  style: const TextStyle(
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                    color: AppTheme.textSecondary,
                    letterSpacing: 1.2,
                  ),
                ),
              ],
            ),

            const SizedBox(height: 8),

            // Description
            Text(
              l10n.gymQrCodeDesc,
              style: const TextStyle(
                fontSize: 12,
                color: AppTheme.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
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
              keyboardType: kIsWeb ? TextInputType.text : TextInputType.phone,
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
            const SizedBox(height: 8),
            // Phone (read-only)
            Text(
              p.phone,
              style: const TextStyle(
                fontSize: 13,
                color: AppTheme.textSecondary,
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _profileNameCtrl,
              decoration: InputDecoration(labelText: l10n.name),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _profileEmailCtrl,
              decoration: InputDecoration(labelText: l10n.email),
              keyboardType: TextInputType.emailAddress,
            ),
            const SizedBox(height: 12),
            DropdownButtonFormField<String>(
              value: _profileGenderCtrl.text.isEmpty
                  ? null
                  : _profileGenderCtrl.text,
              decoration: InputDecoration(labelText: l10n.gender),
              items: [
                DropdownMenuItem(value: 'male', child: Text(l10n.male)),
                DropdownMenuItem(value: 'female', child: Text(l10n.female)),
                DropdownMenuItem(value: 'other', child: Text(l10n.other)),
              ],
              onChanged: (v) {
                if (v != null) _profileGenderCtrl.text = v;
              },
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _profileDobCtrl,
              decoration: InputDecoration(
                labelText: l10n.dateOfBirth,
                hintText: 'YYYY-MM-DD',
              ),
              readOnly: true,
              onTap: () async {
                final now = DateTime.now();
                final initial = DateTime.tryParse(_profileDobCtrl.text) ??
                    DateTime(now.year - 25);
                final picked = await showDatePicker(
                  context: context,
                  initialDate: initial,
                  firstDate: DateTime(1940),
                  lastDate: now,
                );
                if (picked != null) {
                  _profileDobCtrl.text =
                      '${picked.year}-${picked.month.toString().padLeft(2, '0')}-${picked.day.toString().padLeft(2, '0')}';
                }
              },
            ),
            const SizedBox(height: 16),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: _savingProfile ? null : _saveProfile,
                child: _savingProfile
                    ? const SizedBox(
                        width: 18,
                        height: 18,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          color: Colors.white,
                        ),
                      )
                    : Text(l10n.saveProfile),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
