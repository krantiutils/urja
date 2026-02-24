import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:mobile_scanner/mobile_scanner.dart';

import '../../config/theme.dart';
import '../../models/member.dart';
import '../../providers/auth_provider.dart';
import '../../providers/locale_provider.dart';
import '../../services/member_service.dart';
import '../../services/onboarding_service.dart';
import '../shared/widgets/loading_indicator.dart';

class MemberSettingsScreen extends ConsumerStatefulWidget {
  const MemberSettingsScreen({super.key});

  @override
  ConsumerState<MemberSettingsScreen> createState() =>
      _MemberSettingsScreenState();
}

class _MemberSettingsScreenState extends ConsumerState<MemberSettingsScreen> {
  late final MemberService _memberService;

  MemberProfile? _profile;
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _memberService = MemberService(ref.read(apiClientProvider));
    _loadProfile();
  }

  Future<void> _loadProfile() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final profile = await _memberService.getProfile();
      setState(() {
        _profile = profile;
        _loading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  void _showSuccess(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: AppTheme.primary,
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      ),
    );
  }

  void _showError(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: AppTheme.error,
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      ),
    );
  }

  Future<void> _logout() async {
    final l10n = AppLocalizations.of(context)!;
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: AppTheme.surface,
        title: Text(l10n.logout),
        content: Text(l10n.logoutConfirmation),
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

    if (confirmed == true) {
      await ref.read(authProvider.notifier).logout();
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    if (_loading) {
      return Scaffold(
        backgroundColor: AppTheme.background,
        body: const LoadingIndicator(),
      );
    }

    if (_error != null) {
      return Scaffold(
        backgroundColor: AppTheme.background,
        body: Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Icon(Icons.error_outline, size: 48, color: AppTheme.error),
              const SizedBox(height: 16),
              Text(l10n.errorLoadingData,
                  style: const TextStyle(color: AppTheme.textSecondary)),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: _loadProfile,
                child: Text(l10n.retry),
              ),
            ],
          ),
        ),
      );
    }

    final authState = ref.watch(authProvider);
    final userType = authState.user?.userType;
    final isNonGymUser =
        userType == 'fitness_tracker' || userType == 'calorie_tracker';

    return Scaffold(
      backgroundColor: AppTheme.background,
      body: ListView(
        padding: EdgeInsets.zero,
        children: [
          _buildProfileHero(l10n),
          const SizedBox(height: 16),
          if (isNonGymUser) ...[
            _buildJoinGymCard(l10n),
            const SizedBox(height: 16),
          ],
          _buildMenuList(l10n),
          const SizedBox(height: 16),
          _buildLogoutTile(l10n),
          const SizedBox(height: 32),
        ],
      ),
    );
  }

  Widget _buildProfileHero(AppLocalizations l10n) {
    final profile = _profile!;
    final orgName = profile.organizations.isNotEmpty
        ? profile.organizations.first.orgName
        : '';
    final initial = profile.name.isNotEmpty ? profile.name[0].toUpperCase() : '?';

    return Container(
      width: double.infinity,
      decoration: const BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [
            Color(0xFF1A2A1A),
            AppTheme.background,
          ],
        ),
      ),
      child: SafeArea(
        bottom: false,
        child: Padding(
          padding: const EdgeInsets.fromLTRB(16, 16, 16, 24),
          child: Column(
            children: [
              // Avatar with camera button
              Stack(
                children: [
                  CircleAvatar(
                    radius: 52,
                    backgroundColor: AppTheme.primary.withAlpha(40),
                    backgroundImage: profile.avatarUrl != null
                        ? NetworkImage(profile.avatarUrl!)
                        : null,
                    child: profile.avatarUrl == null
                        ? Text(initial,
                            style: const TextStyle(
                                fontSize: 40,
                                fontWeight: FontWeight.bold,
                                color: AppTheme.primary))
                        : null,
                  ),
                  Positioned(
                    bottom: 0,
                    right: 0,
                    child: GestureDetector(
                      onTap: () => _showPhotoOptions(l10n),
                      child: Container(
                        width: 34,
                        height: 34,
                        decoration: BoxDecoration(
                          color: AppTheme.primary,
                          shape: BoxShape.circle,
                          border: Border.all(color: AppTheme.background, width: 3),
                        ),
                        child: const Icon(Icons.camera_alt,
                            color: Colors.white, size: 16),
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              // Name
              Text(
                profile.name,
                style: const TextStyle(
                  fontSize: 22,
                  fontWeight: FontWeight.bold,
                  color: AppTheme.textPrimary,
                ),
              ),
              if (orgName.isNotEmpty) ...[
                const SizedBox(height: 4),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Icon(Icons.fitness_center,
                        size: 14, color: AppTheme.primary),
                    const SizedBox(width: 4),
                    Text(orgName,
                        style: const TextStyle(
                            color: AppTheme.textSecondary, fontSize: 14)),
                  ],
                ),
              ],
              const SizedBox(height: 4),
              Text(profile.phone,
                  style: const TextStyle(
                      color: AppTheme.textSecondary, fontSize: 14)),
            ],
          ),
        ),
      ),
    );
  }

  void _showPhotoOptions(AppLocalizations l10n) {
    showModalBottomSheet(
      context: context,
      backgroundColor: AppTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (context) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 16),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              ListTile(
                leading: const Icon(Icons.camera_alt, color: AppTheme.primary),
                title: Text(l10n.camera,
                    style: const TextStyle(color: AppTheme.textPrimary)),
                onTap: () {
                  Navigator.pop(context);
                  // TODO: image_picker camera
                },
              ),
              ListTile(
                leading:
                    const Icon(Icons.photo_library, color: AppTheme.primary),
                title: Text(l10n.gallery,
                    style: const TextStyle(color: AppTheme.textPrimary)),
                onTap: () {
                  Navigator.pop(context);
                  // TODO: image_picker gallery
                },
              ),
              if (_profile?.avatarUrl != null)
                ListTile(
                  leading: const Icon(Icons.delete_outline, color: AppTheme.error),
                  title: Text(l10n.removePhoto,
                      style: const TextStyle(color: AppTheme.error)),
                  onTap: () {
                    Navigator.pop(context);
                    // TODO: remove avatar
                  },
                ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildMenuList(AppLocalizations l10n) {
    final locale = ref.watch(localeProvider);

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      decoration: BoxDecoration(
        color: AppTheme.surface,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        children: [
          _menuItem(
            icon: Icons.person_outline,
            label: l10n.editMyProfile,
            onTap: () => _showEditProfileSheet(l10n),
          ),
          _divider(),
          _menuItem(
            icon: Icons.straighten,
            label: l10n.measurements,
            onTap: () => context.go('/member/health'),
          ),
          _divider(),
          _menuItem(
            icon: Icons.restaurant_menu,
            label: 'Nutrition',
            onTap: () => context.go('/member/nutrition'),
          ),
          _divider(),
          _menuItem(
            icon: Icons.card_membership,
            label: l10n.subscriptionHistory,
            onTap: () => context.go('/member/packages'),
          ),
          _divider(),
          _menuItem(
            icon: Icons.tune,
            label: l10n.participationSettings,
            onTap: () => _showPrivacySheet(l10n),
          ),
          _divider(),
          _menuItem(
            icon: Icons.emergency,
            label: l10n.emergencyContact,
            onTap: () => _showEmergencySheet(l10n),
          ),
          _divider(),
          _menuItem(
            icon: Icons.rate_review_outlined,
            label: l10n.reviewsFeedbacks,
            onTap: () => _showFeedbackSheet(l10n),
          ),
          _divider(),
          _menuItem(
            icon: Icons.language,
            label: l10n.language,
            trailing: Text(
              locale.languageCode == 'en' ? 'English' : 'नेपाली',
              style: const TextStyle(color: AppTheme.primary, fontSize: 14),
            ),
            onTap: () => ref.read(localeProvider.notifier).toggle(),
          ),
        ],
      ),
    );
  }

  Widget _menuItem({
    required IconData icon,
    required String label,
    required VoidCallback onTap,
    Widget? trailing,
  }) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(8),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        child: Row(
          children: [
            Icon(icon, color: AppTheme.textSecondary, size: 22),
            const SizedBox(width: 14),
            Expanded(
              child: Text(label,
                  style: const TextStyle(
                      color: AppTheme.textPrimary, fontSize: 15)),
            ),
            trailing ??
                const Icon(Icons.chevron_right,
                    color: AppTheme.textSecondary, size: 20),
          ],
        ),
      ),
    );
  }

  Widget _divider() {
    return const Divider(height: 1, indent: 52, endIndent: 16, color: Color(0xFF2D2D44));
  }

  Widget _buildLogoutTile(AppLocalizations l10n) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      decoration: BoxDecoration(
        color: AppTheme.surface,
        borderRadius: BorderRadius.circular(16),
      ),
      child: InkWell(
        onTap: _logout,
        borderRadius: BorderRadius.circular(16),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
          child: Row(
            children: [
              const Icon(Icons.logout, color: AppTheme.error, size: 22),
              const SizedBox(width: 14),
              Text(l10n.logout,
                  style: const TextStyle(color: AppTheme.error, fontSize: 15)),
            ],
          ),
        ),
      ),
    );
  }

  // --- Bottom Sheet: Edit Profile ---
  void _showEditProfileSheet(AppLocalizations l10n) {
    final nameCtrl = TextEditingController(text: _profile?.name ?? '');
    final emailCtrl = TextEditingController(text: _profile?.email ?? '');
    final dobCtrl = TextEditingController(text: _profile?.dateOfBirth ?? '');
    String? gender = _profile?.gender;
    bool saving = false;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: AppTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setSheetState) => Padding(
          padding: EdgeInsets.fromLTRB(
              16, 24, 16, MediaQuery.of(ctx).viewInsets.bottom + 24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(l10n.editMyProfile,
                  style: const TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                      color: AppTheme.textPrimary)),
              const SizedBox(height: 20),
              TextField(
                controller: nameCtrl,
                style: const TextStyle(color: AppTheme.textPrimary),
                decoration: InputDecoration(
                  labelText: l10n.name,
                  labelStyle: const TextStyle(color: AppTheme.textSecondary),
                  prefixIcon: const Icon(Icons.person_outline,
                      color: AppTheme.textSecondary, size: 20),
                ),
              ),
              const SizedBox(height: 12),
              TextField(
                controller: emailCtrl,
                keyboardType: TextInputType.emailAddress,
                style: const TextStyle(color: AppTheme.textPrimary),
                decoration: InputDecoration(
                  labelText: l10n.email,
                  labelStyle: const TextStyle(color: AppTheme.textSecondary),
                  prefixIcon: const Icon(Icons.email_outlined,
                      color: AppTheme.textSecondary, size: 20),
                ),
              ),
              const SizedBox(height: 12),
              TextField(
                controller: dobCtrl,
                keyboardType: TextInputType.datetime,
                style: const TextStyle(color: AppTheme.textPrimary),
                decoration: InputDecoration(
                  labelText: l10n.dateOfBirth,
                  labelStyle: const TextStyle(color: AppTheme.textSecondary),
                  hintText: 'YYYY-MM-DD',
                  hintStyle:
                      TextStyle(color: AppTheme.textSecondary.withAlpha(100)),
                  prefixIcon: const Icon(Icons.cake_outlined,
                      color: AppTheme.textSecondary, size: 20),
                ),
              ),
              const SizedBox(height: 12),
              DropdownButtonFormField<String>(
                value: gender,
                dropdownColor: AppTheme.surface,
                style: const TextStyle(color: AppTheme.textPrimary),
                decoration: InputDecoration(
                  labelText: l10n.gender,
                  labelStyle: const TextStyle(color: AppTheme.textSecondary),
                  prefixIcon: const Icon(Icons.wc,
                      color: AppTheme.textSecondary, size: 20),
                ),
                items: [
                  DropdownMenuItem(value: 'male', child: Text(l10n.male)),
                  DropdownMenuItem(value: 'female', child: Text(l10n.female)),
                  DropdownMenuItem(value: 'other', child: Text(l10n.other)),
                ],
                onChanged: (v) => setSheetState(() => gender = v),
              ),
              const SizedBox(height: 20),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: saving
                      ? null
                      : () async {
                          setSheetState(() => saving = true);
                          try {
                            await _memberService.updateProfile({
                              'name': nameCtrl.text.trim(),
                              'email': emailCtrl.text.trim().isEmpty
                                  ? null
                                  : emailCtrl.text.trim(),
                              'date_of_birth': dobCtrl.text.trim().isEmpty
                                  ? null
                                  : dobCtrl.text.trim(),
                              'gender': gender,
                            });
                            if (ctx.mounted) Navigator.pop(ctx);
                            _showSuccess(l10n.saved);
                            _loadProfile();
                          } catch (e) {
                            setSheetState(() => saving = false);
                            _showError(l10n.error);
                          }
                        },
                  child: saving
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(
                              color: Colors.white, strokeWidth: 2))
                      : Text(l10n.saveProfile),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  // --- Bottom Sheet: Emergency Contact ---
  void _showEmergencySheet(AppLocalizations l10n) {
    final nameCtrl =
        TextEditingController(text: _profile?.emergencyContactName ?? '');
    final phoneCtrl =
        TextEditingController(text: _profile?.emergencyContactPhone ?? '');
    bool saving = false;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: AppTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setSheetState) => Padding(
          padding: EdgeInsets.fromLTRB(
              16, 24, 16, MediaQuery.of(ctx).viewInsets.bottom + 24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(l10n.emergencyContact,
                  style: const TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                      color: AppTheme.textPrimary)),
              const SizedBox(height: 20),
              TextField(
                controller: nameCtrl,
                style: const TextStyle(color: AppTheme.textPrimary),
                decoration: InputDecoration(
                  labelText: l10n.contactName,
                  labelStyle: const TextStyle(color: AppTheme.textSecondary),
                  prefixIcon: const Icon(Icons.person_outline,
                      color: AppTheme.textSecondary, size: 20),
                ),
              ),
              const SizedBox(height: 12),
              TextField(
                controller: phoneCtrl,
                keyboardType: kIsWeb ? TextInputType.text : TextInputType.phone,
                style: const TextStyle(color: AppTheme.textPrimary),
                decoration: InputDecoration(
                  labelText: l10n.contactPhone,
                  labelStyle: const TextStyle(color: AppTheme.textSecondary),
                  prefixIcon: const Icon(Icons.phone_outlined,
                      color: AppTheme.textSecondary, size: 20),
                ),
              ),
              const SizedBox(height: 20),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: saving
                      ? null
                      : () async {
                          setSheetState(() => saving = true);
                          try {
                            await _memberService.updateProfile({
                              'emergency_contact_name':
                                  nameCtrl.text.trim().isEmpty
                                      ? null
                                      : nameCtrl.text.trim(),
                              'emergency_contact_phone':
                                  phoneCtrl.text.trim().isEmpty
                                      ? null
                                      : phoneCtrl.text.trim(),
                            });
                            if (ctx.mounted) Navigator.pop(ctx);
                            _showSuccess(l10n.saved);
                            _loadProfile();
                          } catch (e) {
                            setSheetState(() => saving = false);
                            _showError(l10n.error);
                          }
                        },
                  child: saving
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(
                              color: Colors.white, strokeWidth: 2))
                      : Text(l10n.saveEmergencyContact),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  // --- Join a Gym Card (for non-gym users) ---
  Widget _buildJoinGymCard(AppLocalizations l10n) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [
            AppTheme.primary.withAlpha(30),
            AppTheme.primary.withAlpha(10),
          ],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: AppTheme.primary.withAlpha(60)),
      ),
      child: InkWell(
        onTap: () => _showJoinGymSheet(l10n),
        borderRadius: BorderRadius.circular(16),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              Container(
                width: 48,
                height: 48,
                decoration: BoxDecoration(
                  color: AppTheme.primary.withAlpha(30),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: const Icon(Icons.qr_code_scanner,
                    color: AppTheme.primary, size: 24),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      l10n.joinAGym,
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                        color: AppTheme.textPrimary,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      l10n.joinAGymDesc,
                      style: const TextStyle(
                        fontSize: 13,
                        color: AppTheme.textSecondary,
                      ),
                    ),
                  ],
                ),
              ),
              const Icon(Icons.chevron_right,
                  color: AppTheme.primary, size: 24),
            ],
          ),
        ),
      ),
    );
  }

  // --- Bottom Sheet: Join a Gym ---
  void _showJoinGymSheet(AppLocalizations l10n) {
    final orgIdCtrl = TextEditingController();
    bool joining = false;
    String? joinError;
    bool joined = false;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: AppTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setSheetState) => Padding(
          padding: EdgeInsets.fromLTRB(
              16, 24, 16, MediaQuery.of(ctx).viewInsets.bottom + 24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(l10n.joinAGym,
                  style: const TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                      color: AppTheme.textPrimary)),
              const SizedBox(height: 8),
              Text(l10n.joinAGymDesc,
                  style: const TextStyle(
                      color: AppTheme.textSecondary, fontSize: 14)),
              const SizedBox(height: 20),
              if (joined) ...[
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: AppTheme.primary.withAlpha(15),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: AppTheme.primary),
                  ),
                  child: Row(
                    children: [
                      const Icon(Icons.check_circle,
                          color: AppTheme.primary, size: 28),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Text(
                          l10n.successfullyJoinedGym,
                          style: const TextStyle(
                              color: AppTheme.textPrimary,
                              fontWeight: FontWeight.w600),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 16),
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: () {
                      Navigator.pop(ctx);
                      // Reload profile and update auth state
                      _loadProfile();
                      ref
                          .read(authProvider.notifier)
                          .updateUserType('gym_member');
                    },
                    child: Text(l10n.done),
                  ),
                ),
              ] else ...[
                if (!kIsWeb) ...[
                  SizedBox(
                    width: double.infinity,
                    child: OutlinedButton.icon(
                      onPressed: () {
                        _openGymQRScanner(ctx, orgIdCtrl, setSheetState,
                            (id) async {
                          setSheetState(() {
                            joining = true;
                            joinError = null;
                          });
                          try {
                            final service =
                                OnboardingService(ref.read(apiClientProvider));
                            await service.joinGym(id);
                            setSheetState(() {
                              joining = false;
                              joined = true;
                            });
                          } catch (e) {
                            setSheetState(() {
                              joining = false;
                              joinError = l10n.couldNotJoinGym;
                            });
                          }
                        });
                      },
                      icon: const Icon(Icons.qr_code_scanner, size: 20),
                      label: Text(l10n.scanGymQrCode),
                      style: OutlinedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        side: const BorderSide(color: AppTheme.primary),
                        foregroundColor: AppTheme.primary,
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      const Expanded(
                          child: Divider(color: AppTheme.surfaceLight)),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 16),
                        child: Text(l10n.orEnterCodeManually,
                            style: TextStyle(
                                color: AppTheme.textSecondary.withAlpha(150),
                                fontSize: 13)),
                      ),
                      const Expanded(
                          child: Divider(color: AppTheme.surfaceLight)),
                    ],
                  ),
                  const SizedBox(height: 16),
                ],
                TextField(
                  controller: orgIdCtrl,
                  style: const TextStyle(
                      color: AppTheme.textPrimary, fontSize: 16),
                  decoration: InputDecoration(
                    labelText: l10n.gymCode,
                    labelStyle:
                        const TextStyle(color: AppTheme.textSecondary),
                    hintText: 'e.g. kathmandu-fitness-hub',
                    hintStyle: TextStyle(
                        color: AppTheme.textSecondary.withAlpha(100)),
                    prefixIcon: const Icon(Icons.business,
                        color: AppTheme.textSecondary, size: 20),
                    errorText: joinError,
                  ),
                ),
                const SizedBox(height: 20),
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: joining
                        ? null
                        : () async {
                            final orgId = orgIdCtrl.text.trim();
                            if (orgId.isEmpty) {
                              setSheetState(
                                  () => joinError = l10n.enterGymCode);
                              return;
                            }
                            setSheetState(() {
                              joining = true;
                              joinError = null;
                            });
                            try {
                              final service = OnboardingService(
                                  ref.read(apiClientProvider));
                              await service.joinGym(orgId);
                              setSheetState(() {
                                joining = false;
                                joined = true;
                              });
                            } catch (e) {
                              setSheetState(() {
                                joining = false;
                                joinError = l10n.couldNotJoinGym;
                              });
                            }
                          },
                    child: joining
                        ? const SizedBox(
                            width: 20,
                            height: 20,
                            child: CircularProgressIndicator(
                                color: Colors.white, strokeWidth: 2))
                        : Text(l10n.joinGym),
                  ),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }

  void _openGymQRScanner(
    BuildContext parentCtx,
    TextEditingController orgIdCtrl,
    StateSetter setParentState,
    Future<void> Function(String orgId) onScanned,
  ) {
    showModalBottomSheet(
      context: parentCtx,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (scanCtx) => Container(
        height: MediaQuery.of(scanCtx).size.height * 0.6,
        decoration: const BoxDecoration(
          color: AppTheme.surface,
          borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
        ),
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    AppLocalizations.of(context)!.scanGymQrCode,
                    style: const TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                      color: AppTheme.textPrimary,
                    ),
                  ),
                  IconButton(
                    onPressed: () => Navigator.pop(scanCtx),
                    icon: const Icon(Icons.close, color: AppTheme.textSecondary),
                  ),
                ],
              ),
            ),
            Expanded(
              child: ClipRRect(
                borderRadius: BorderRadius.circular(12),
                child: MobileScanner(
                  onDetect: (capture) {
                    final barcode = capture.barcodes.firstOrNull;
                    if (barcode?.rawValue != null) {
                      Navigator.pop(scanCtx);
                      orgIdCtrl.text = barcode!.rawValue!;
                      onScanned(barcode.rawValue!);
                    }
                  },
                ),
              ),
            ),
            const SizedBox(height: 16),
          ],
        ),
      ),
    );
  }

  // --- Bottom Sheet: Feedback ---
  void _showFeedbackSheet(AppLocalizations l10n) {
    final commentCtrl = TextEditingController();
    int selectedRating = 0;
    bool submitting = false;

    // Only show feedback for users that belong to an org
    if (_profile == null || _profile!.organizations.isEmpty) {
      _showError(l10n.error);
      return;
    }

    final orgId = _profile!.organizations.first.orgId;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: AppTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setSheetState) => Padding(
          padding: EdgeInsets.fromLTRB(
              16, 24, 16, MediaQuery.of(ctx).viewInsets.bottom + 24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(l10n.submitFeedback,
                  style: const TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                      color: AppTheme.textPrimary)),
              const SizedBox(height: 20),
              // Star rating
              Text(l10n.rateYourExperience,
                  style: const TextStyle(
                      color: AppTheme.textSecondary, fontSize: 14)),
              const SizedBox(height: 8),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: List.generate(5, (index) {
                  final starIndex = index + 1;
                  return GestureDetector(
                    onTap: () =>
                        setSheetState(() => selectedRating = starIndex),
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 4),
                      child: Icon(
                        starIndex <= selectedRating
                            ? Icons.star_rounded
                            : Icons.star_outline_rounded,
                        color: starIndex <= selectedRating
                            ? Colors.amber
                            : AppTheme.textSecondary,
                        size: 40,
                      ),
                    ),
                  );
                }),
              ),
              const SizedBox(height: 16),
              // Comment text field
              TextField(
                controller: commentCtrl,
                maxLines: 4,
                onChanged: (_) => setSheetState(() {}),
                style: const TextStyle(color: AppTheme.textPrimary),
                decoration: InputDecoration(
                  hintText: l10n.feedbackHint,
                  hintStyle:
                      TextStyle(color: AppTheme.textSecondary.withAlpha(150)),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide:
                        const BorderSide(color: AppTheme.surfaceLight),
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide:
                        const BorderSide(color: AppTheme.surfaceLight),
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: const BorderSide(color: AppTheme.primary),
                  ),
                ),
              ),
              const SizedBox(height: 20),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: submitting ||
                          (selectedRating == 0 &&
                              commentCtrl.text.trim().isEmpty)
                      ? null
                      : () async {
                          setSheetState(() => submitting = true);
                          try {
                            final api = ref.read(apiClientProvider);
                            await api.post(
                              '/orgs/$orgId/feedbacks',
                              data: {
                                'message': commentCtrl.text.trim(),
                                'rating': selectedRating > 0
                                    ? selectedRating
                                    : null,
                              },
                            );
                            if (ctx.mounted) Navigator.pop(ctx);
                            _showSuccess(l10n.feedbackSubmitted);
                          } catch (e) {
                            setSheetState(() => submitting = false);
                            _showError(l10n.error);
                          }
                        },
                  child: submitting
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(
                              color: Colors.white, strokeWidth: 2))
                      : Text(l10n.submit),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  // --- Bottom Sheet: Privacy / Participation Settings ---
  void _showPrivacySheet(AppLocalizations l10n) {
    bool showEmail = _profile?.privacySettings.showEmail ?? false;
    bool showPhone = _profile?.privacySettings.showPhone ?? false;
    bool showProfile = _profile?.privacySettings.showProfile ?? false;
    bool showAttendance = _profile?.privacySettings.showAttendance ?? false;
    bool showOnLeaderboard =
        _profile?.privacySettings.showOnLeaderboard ?? false;
    bool saving = false;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: AppTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setSheetState) => Padding(
          padding: const EdgeInsets.fromLTRB(0, 24, 0, 16),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                child: Text(l10n.participationSettings,
                    style: const TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                        color: AppTheme.textPrimary)),
              ),
              const SizedBox(height: 12),
              SwitchListTile(
                title: Text(l10n.showEmail,
                    style: const TextStyle(
                        color: AppTheme.textPrimary, fontSize: 15)),
                subtitle: Text(l10n.showEmailDesc,
                    style: const TextStyle(
                        color: AppTheme.textSecondary, fontSize: 12)),
                value: showEmail,
                onChanged: (v) => setSheetState(() => showEmail = v),
                activeThumbColor: AppTheme.primary,
              ),
              SwitchListTile(
                title: Text(l10n.showPhone,
                    style: const TextStyle(
                        color: AppTheme.textPrimary, fontSize: 15)),
                subtitle: Text(l10n.showPhoneDesc,
                    style: const TextStyle(
                        color: AppTheme.textSecondary, fontSize: 12)),
                value: showPhone,
                onChanged: (v) => setSheetState(() => showPhone = v),
                activeThumbColor: AppTheme.primary,
              ),
              SwitchListTile(
                title: Text(l10n.showProfile,
                    style: const TextStyle(
                        color: AppTheme.textPrimary, fontSize: 15)),
                subtitle: Text(l10n.showProfileDesc,
                    style: const TextStyle(
                        color: AppTheme.textSecondary, fontSize: 12)),
                value: showProfile,
                onChanged: (v) => setSheetState(() => showProfile = v),
                activeThumbColor: AppTheme.primary,
              ),
              SwitchListTile(
                title: Text(l10n.showAttendance,
                    style: const TextStyle(
                        color: AppTheme.textPrimary, fontSize: 15)),
                subtitle: Text(l10n.showAttendanceDesc,
                    style: const TextStyle(
                        color: AppTheme.textSecondary, fontSize: 12)),
                value: showAttendance,
                onChanged: (v) => setSheetState(() => showAttendance = v),
                activeThumbColor: AppTheme.primary,
              ),
              SwitchListTile(
                title: Text(l10n.showOnLeaderboard,
                    style: const TextStyle(
                        color: AppTheme.textPrimary, fontSize: 15)),
                subtitle: Text(l10n.showOnLeaderboardDesc,
                    style: const TextStyle(
                        color: AppTheme.textSecondary, fontSize: 12)),
                value: showOnLeaderboard,
                onChanged: (v) => setSheetState(() => showOnLeaderboard = v),
                activeThumbColor: AppTheme.primary,
              ),
              Padding(
                padding: const EdgeInsets.fromLTRB(16, 8, 16, 0),
                child: SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: saving
                        ? null
                        : () async {
                            setSheetState(() => saving = true);
                            try {
                              await _memberService.updatePrivacy({
                                'show_email': showEmail,
                                'show_phone': showPhone,
                                'show_profile': showProfile,
                                'show_attendance': showAttendance,
                                'show_on_leaderboard': showOnLeaderboard,
                              });
                              if (ctx.mounted) Navigator.pop(ctx);
                              _showSuccess(l10n.saved);
                              _loadProfile();
                            } catch (e) {
                              setSheetState(() => saving = false);
                              _showError(l10n.error);
                            }
                          },
                    child: saving
                        ? const SizedBox(
                            width: 20,
                            height: 20,
                            child: CircularProgressIndicator(
                                color: Colors.white, strokeWidth: 2))
                        : Text(l10n.savePrivacy),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
