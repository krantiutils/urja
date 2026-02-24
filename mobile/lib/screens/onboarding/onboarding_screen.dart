import 'dart:math';

import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:mobile_scanner/mobile_scanner.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../services/notification_service.dart';
import '../../services/onboarding_service.dart';

class OnboardingScreen extends ConsumerStatefulWidget {
  const OnboardingScreen({super.key});

  @override
  ConsumerState<OnboardingScreen> createState() => _OnboardingScreenState();
}

class _OnboardingScreenState extends ConsumerState<OnboardingScreen> {
  final PageController _pageController = PageController();
  int _currentPage = 0;
  bool _submitting = false;

  // Step 1: User type
  String? _userType;

  // Step 2: Gym org ID (gym_member only)
  final _orgIdController = TextEditingController();
  bool _joiningGym = false;
  String? _joinGymError;
  bool _gymJoined = false;

  // Step 3: Goal
  String? _goal;

  // Step 4: About you
  final _nameController = TextEditingController();
  final _weightController = TextEditingController();
  final _heightController = TextEditingController();
  final _ageController = TextEditingController();
  String _gender = 'male';
  String _activityLevel = 'moderate';

  // Step 5: Results
  Map<String, dynamic>? _results;

  int get _totalPages {
    if (_userType == 'gym_member') return 5;
    return 4; // Skip gym join step
  }

  int _adjustedPage(int page) {
    // If not gym_member, skip step 2 (gym join)
    if (_userType != 'gym_member' && page >= 1) {
      return page + 1;
    }
    return page;
  }

  @override
  void dispose() {
    _pageController.dispose();
    _orgIdController.dispose();
    _nameController.dispose();
    _weightController.dispose();
    _heightController.dispose();
    _ageController.dispose();
    super.dispose();
  }

  void _nextPage() {
    if (_currentPage < _totalPages - 1) {
      _pageController.nextPage(
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeInOut,
      );
    }
  }

  void _previousPage() {
    if (_currentPage > 0) {
      _pageController.previousPage(
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeInOut,
      );
    }
  }

  Future<void> _joinGym() async {
    final orgId = _orgIdController.text.trim();
    if (orgId.isEmpty) {
      setState(() => _joinGymError = 'Please enter your gym code');
      return;
    }

    setState(() {
      _joiningGym = true;
      _joinGymError = null;
    });

    try {
      final service = OnboardingService(ref.read(apiClientProvider));
      await service.joinGym(orgId);
      setState(() {
        _joiningGym = false;
        _gymJoined = true;
      });
    } catch (e) {
      setState(() {
        _joiningGym = false;
        _joinGymError = 'Could not join gym. Please check the code.';
      });
    }
  }

  Future<void> _submitOnboarding() async {
    setState(() => _submitting = true);

    try {
      final service = OnboardingService(ref.read(apiClientProvider));

      final data = <String, dynamic>{
        'user_type': _userType,
        'goal_type': _goal,
        'name': _nameController.text.trim(),
        'weight_kg': double.tryParse(_weightController.text) ?? 0,
        'height_cm': double.tryParse(_heightController.text) ?? 0,
        'age': int.tryParse(_ageController.text) ?? 0,
        'gender': _gender,
        'activity_level': _activityLevel,
      };

      final response = await service.submitOnboarding(data);
      setState(() {
        _results = response;
        _submitting = false;
      });

      // Move to results page
      _nextPage();
    } catch (e) {
      setState(() => _submitting = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to complete setup: $e'),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    }
  }

  void _completeOnboarding() {
    ref.read(authProvider.notifier).updateUserType(_userType ?? 'gym_member');
    ref.read(authProvider.notifier).completeOnboarding();

    // Schedule periodic water & food intake reminders
    NotificationService().scheduleAllReminders();

    // Navigate based on user type
    if (_userType == 'gym_member') {
      context.go('/member');
    } else {
      context.go('/tracker');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // Progress indicator
            Padding(
              padding: const EdgeInsets.fromLTRB(24, 16, 24, 0),
              child: Row(
                children: List.generate(_totalPages, (index) {
                  return Expanded(
                    child: Container(
                      height: 4,
                      margin: const EdgeInsets.symmetric(horizontal: 2),
                      decoration: BoxDecoration(
                        color: index <= _currentPage
                            ? AppTheme.primary
                            : AppTheme.surfaceLight,
                        borderRadius: BorderRadius.circular(2),
                      ),
                    ),
                  );
                }),
              ),
            ),
            // Page content
            Expanded(
              child: PageView.builder(
                controller: _pageController,
                physics: const NeverScrollableScrollPhysics(),
                onPageChanged: (page) {
                  setState(() => _currentPage = page);
                },
                itemCount: _totalPages,
                itemBuilder: (context, index) {
                  final logical = _adjustedPage(index);
                  switch (logical) {
                    case 0:
                      return _buildUserTypePage();
                    case 1:
                      return _buildGymJoinPage();
                    case 2:
                      return _buildGoalPage();
                    case 3:
                      return _buildAboutYouPage();
                    case 4:
                      return _buildResultsPage();
                    default:
                      return const SizedBox.shrink();
                  }
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  // ─── Step 1: User Type ──────────────────────────────────────────────────────

  Widget _buildUserTypePage() {
    return LayoutBuilder(
      builder: (context, constraints) => SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: ConstrainedBox(
          constraints: BoxConstraints(minHeight: constraints.maxHeight),
          child: IntrinsicHeight(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const SizedBox(height: 32),
                const Text(
                  'What brings you here?',
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.bold,
                    color: AppTheme.textPrimary,
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  'Choose how you want to use the app',
                  style: TextStyle(
                    fontSize: 16,
                    color: AppTheme.textSecondary,
                  ),
                ),
                const SizedBox(height: 40),
                _buildUserTypeCard(
                  value: 'gym_member',
                  icon: Icons.fitness_center,
                  iconColor: AppTheme.primary,
                  title: 'Gym Member',
                  subtitle: 'I go to a gym and want to track my membership, attendance, workouts and nutrition',
                ),
                const SizedBox(height: 12),
                _buildUserTypeCard(
                  value: 'fitness_tracker',
                  icon: Icons.directions_run,
                  iconColor: AppTheme.info,
                  title: 'Fitness Tracker',
                  subtitle: 'I work out on my own and want to track workouts and nutrition',
                ),
                const SizedBox(height: 12),
                _buildUserTypeCard(
                  value: 'calorie_tracker',
                  icon: Icons.restaurant,
                  iconColor: AppTheme.warning,
                  title: 'Calorie Tracker',
                  subtitle: 'I just want to track what I eat and manage my calories',
                ),
                const Spacer(),
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: _userType != null ? _nextPage : null,
                    child: const Text('Continue'),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildUserTypeCard({
    required String value,
    required IconData icon,
    required Color iconColor,
    required String title,
    required String subtitle,
  }) {
    final isSelected = _userType == value;
    return InkWell(
      onTap: () => setState(() => _userType = value),
      borderRadius: BorderRadius.circular(16),
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: isSelected
              ? AppTheme.primary.withAlpha(15)
              : AppTheme.surface,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(
            color: isSelected ? AppTheme.primary : Colors.transparent,
            width: 2,
          ),
        ),
        child: Row(
          children: [
            Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                color: iconColor.withAlpha(25),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(icon, color: iconColor, size: 24),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: AppTheme.textPrimary,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    subtitle,
                    style: const TextStyle(
                      fontSize: 13,
                      color: AppTheme.textSecondary,
                    ),
                  ),
                ],
              ),
            ),
            Radio<String>(
              value: value,
              groupValue: _userType,
              activeColor: AppTheme.primary,
              onChanged: (v) => setState(() => _userType = v),
            ),
          ],
        ),
      ),
    );
  }

  // ─── Step 2: Gym Join (gym_member only) ─────────────────────────────────────

  void _openQRScanner() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        height: MediaQuery.of(context).size.height * 0.6,
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
                  const Text(
                    'Scan Gym QR Code',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                      color: AppTheme.textPrimary,
                    ),
                  ),
                  IconButton(
                    onPressed: () => Navigator.pop(context),
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
                      Navigator.pop(context);
                      _orgIdController.text = barcode!.rawValue!;
                      _joinGym();
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

  Widget _buildGymJoinPage() {
    return LayoutBuilder(
      builder: (context, constraints) => SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: ConstrainedBox(
          constraints: BoxConstraints(minHeight: constraints.maxHeight),
          child: IntrinsicHeight(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const SizedBox(height: 32),
                const Text(
                  'Join your gym',
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.bold,
                    color: AppTheme.textPrimary,
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  'Scan the QR code at your gym or enter the gym code manually',
                  style: TextStyle(
                    fontSize: 16,
                    color: AppTheme.textSecondary,
                  ),
                ),
                const SizedBox(height: 40),
                if (_gymJoined) ...[
                  Container(
                    padding: const EdgeInsets.all(20),
                    decoration: BoxDecoration(
                      color: AppTheme.primary.withAlpha(15),
                      borderRadius: BorderRadius.circular(16),
                      border: Border.all(color: AppTheme.primary, width: 2),
                    ),
                    child: Row(
                      children: [
                        Container(
                          width: 48,
                          height: 48,
                          decoration: BoxDecoration(
                            color: AppTheme.primary.withAlpha(30),
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: const Icon(Icons.check_circle,
                              color: AppTheme.primary, size: 28),
                        ),
                        const SizedBox(width: 16),
                        const Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                'Successfully joined!',
                                style: TextStyle(
                                  fontSize: 16,
                                  fontWeight: FontWeight.w600,
                                  color: AppTheme.textPrimary,
                                ),
                              ),
                              SizedBox(height: 4),
                              Text(
                                'You are now connected to your gym',
                                style: TextStyle(
                                  fontSize: 13,
                                  color: AppTheme.textSecondary,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ] else ...[
                  if (!kIsWeb) ...[
                    // QR Scan button
                    SizedBox(
                      width: double.infinity,
                      child: OutlinedButton.icon(
                        onPressed: _openQRScanner,
                        icon: const Icon(Icons.qr_code_scanner, size: 24),
                        label: const Text('Scan QR Code'),
                        style: OutlinedButton.styleFrom(
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          side: const BorderSide(color: AppTheme.primary),
                          foregroundColor: AppTheme.primary,
                        ),
                      ),
                    ),
                    const SizedBox(height: 20),
                    Row(
                      children: [
                        const Expanded(child: Divider(color: AppTheme.surfaceLight)),
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 16),
                          child: Text('or enter code manually',
                              style: TextStyle(
                                  color: AppTheme.textSecondary.withAlpha(150),
                                  fontSize: 13)),
                        ),
                        const Expanded(child: Divider(color: AppTheme.surfaceLight)),
                      ],
                    ),
                    const SizedBox(height: 20),
                  ],
                  TextField(
                    controller: _orgIdController,
                    style:
                        const TextStyle(color: AppTheme.textPrimary, fontSize: 16),
                    decoration: InputDecoration(
                      labelText: 'Gym Code',
                      labelStyle: const TextStyle(color: AppTheme.textSecondary),
                      hintText: 'e.g. GYM-XXXXX',
                      hintStyle:
                          TextStyle(color: AppTheme.textSecondary.withAlpha(100)),
                      prefixIcon: const Icon(Icons.business,
                          color: AppTheme.textSecondary, size: 20),
                      errorText: _joinGymError,
                    ),
                  ),
                  const SizedBox(height: 20),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: _joiningGym ? null : _joinGym,
                      child: _joiningGym
                          ? const SizedBox(
                              width: 24,
                              height: 24,
                              child: CircularProgressIndicator(
                                color: Colors.white,
                                strokeWidth: 2.5,
                              ),
                            )
                          : const Text('Join Gym'),
                    ),
                  ),
                ],
                const Spacer(),
                Row(
                  children: [
                    TextButton(
                      onPressed: _previousPage,
                      child: const Text('Back'),
                    ),
                    const Spacer(),
                    if (_gymJoined)
                      ElevatedButton(
                        onPressed: _nextPage,
                        child: const Text('Continue'),
                      )
                    else
                      TextButton(
                        onPressed: _nextPage,
                        child: const Text('Skip for now'),
                      ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  // ─── Step 3: Goal ───────────────────────────────────────────────────────────

  Widget _buildGoalPage() {
    return LayoutBuilder(
      builder: (context, constraints) => SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: ConstrainedBox(
          constraints: BoxConstraints(minHeight: constraints.maxHeight),
          child: IntrinsicHeight(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const SizedBox(height: 32),
                const Text(
                  "What's your goal?",
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.bold,
                    color: AppTheme.textPrimary,
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  'This helps us personalize your experience',
                  style: TextStyle(
                    fontSize: 16,
                    color: AppTheme.textSecondary,
                  ),
                ),
                const SizedBox(height: 40),
                _buildGoalCard(
                  value: 'lose_weight',
                  icon: Icons.local_fire_department,
                  iconColor: Colors.orange,
                  title: 'Lose Weight',
                  subtitle: 'Burn fat and get lean',
                ),
                const SizedBox(height: 12),
                _buildGoalCard(
                  value: 'build_muscle',
                  icon: Icons.fitness_center,
                  iconColor: Colors.blue,
                  title: 'Build Muscle',
                  subtitle: 'Gain strength and size',
                ),
                const SizedBox(height: 12),
                _buildGoalCard(
                  value: 'stay_fit',
                  icon: Icons.favorite,
                  iconColor: Colors.red,
                  title: 'Stay Fit',
                  subtitle: 'Maintain general fitness',
                ),
                const SizedBox(height: 12),
                _buildGoalCard(
                  value: 'general_health',
                  icon: Icons.spa,
                  iconColor: Colors.green,
                  title: 'General Health',
                  subtitle: 'Eat better and feel better',
                ),
                const Spacer(),
                Row(
                  children: [
                    TextButton(
                      onPressed: _previousPage,
                      child: const Text('Back'),
                    ),
                    const Spacer(),
                    ElevatedButton(
                      onPressed: _goal != null ? _nextPage : null,
                      child: const Text('Continue'),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildGoalCard({
    required String value,
    required IconData icon,
    required Color iconColor,
    required String title,
    required String subtitle,
  }) {
    final isSelected = _goal == value;
    return InkWell(
      onTap: () => setState(() => _goal = value),
      borderRadius: BorderRadius.circular(16),
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: isSelected
              ? AppTheme.primary.withAlpha(15)
              : AppTheme.surface,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(
            color: isSelected ? AppTheme.primary : Colors.transparent,
            width: 2,
          ),
        ),
        child: Row(
          children: [
            Container(
              width: 44,
              height: 44,
              decoration: BoxDecoration(
                color: iconColor.withAlpha(25),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(icon, color: iconColor, size: 22),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: AppTheme.textPrimary,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    subtitle,
                    style: const TextStyle(
                      fontSize: 13,
                      color: AppTheme.textSecondary,
                    ),
                  ),
                ],
              ),
            ),
            Radio<String>(
              value: value,
              groupValue: _goal,
              activeColor: AppTheme.primary,
              onChanged: (v) => setState(() => _goal = v),
            ),
          ],
        ),
      ),
    );
  }

  // ─── Step 4: About You ──────────────────────────────────────────────────────

  Widget _buildAboutYouPage() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SizedBox(height: 32),
          const Text(
            'About You',
            style: TextStyle(
              fontSize: 28,
              fontWeight: FontWeight.bold,
              color: AppTheme.textPrimary,
            ),
          ),
          const SizedBox(height: 8),
          const Text(
            'Tell us about yourself so we can calculate your targets',
            style: TextStyle(
              fontSize: 16,
              color: AppTheme.textSecondary,
            ),
          ),
          const SizedBox(height: 32),
          // Name
          TextField(
            controller: _nameController,
            style: const TextStyle(color: AppTheme.textPrimary),
            decoration: const InputDecoration(
              labelText: 'Your Name',
              labelStyle: TextStyle(color: AppTheme.textSecondary),
              prefixIcon: Icon(Icons.person_outline,
                  color: AppTheme.textSecondary, size: 20),
            ),
          ),
          const SizedBox(height: 16),
          // Weight & Height
          Row(
            children: [
              Expanded(
                child: TextField(
                  controller: _weightController,
                  keyboardType: kIsWeb ? TextInputType.text : TextInputType.number,
                  style: const TextStyle(color: AppTheme.textPrimary),
                  decoration: const InputDecoration(
                    labelText: 'Weight (kg)',
                    labelStyle: TextStyle(color: AppTheme.textSecondary),
                    prefixIcon: Icon(Icons.monitor_weight_outlined,
                        color: AppTheme.textSecondary, size: 20),
                  ),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: TextField(
                  controller: _heightController,
                  keyboardType: kIsWeb ? TextInputType.text : TextInputType.number,
                  style: const TextStyle(color: AppTheme.textPrimary),
                  decoration: const InputDecoration(
                    labelText: 'Height (cm)',
                    labelStyle: TextStyle(color: AppTheme.textSecondary),
                    prefixIcon: Icon(Icons.height,
                        color: AppTheme.textSecondary, size: 20),
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          // Age
          TextField(
            controller: _ageController,
            keyboardType: kIsWeb ? TextInputType.text : TextInputType.number,
            style: const TextStyle(color: AppTheme.textPrimary),
            decoration: const InputDecoration(
              labelText: 'Age',
              labelStyle: TextStyle(color: AppTheme.textSecondary),
              prefixIcon: Icon(Icons.cake_outlined,
                  color: AppTheme.textSecondary, size: 20),
            ),
          ),
          const SizedBox(height: 20),
          // Gender
          const Text('Gender',
              style: TextStyle(color: AppTheme.textSecondary, fontSize: 14)),
          const SizedBox(height: 8),
          SegmentedButton<String>(
            segments: const [
              ButtonSegment(value: 'male', label: Text('Male')),
              ButtonSegment(value: 'female', label: Text('Female')),
              ButtonSegment(value: 'other', label: Text('Other')),
            ],
            selected: {_gender},
            onSelectionChanged: (v) => setState(() => _gender = v.first),
            style: SegmentedButton.styleFrom(
              backgroundColor: AppTheme.surfaceLight,
              foregroundColor: AppTheme.textSecondary,
              selectedForegroundColor: Colors.white,
              selectedBackgroundColor: AppTheme.primary,
            ),
          ),
          const SizedBox(height: 20),
          // Activity Level
          const Text('Activity Level',
              style: TextStyle(color: AppTheme.textSecondary, fontSize: 14)),
          const SizedBox(height: 8),
          DropdownButtonFormField<String>(
            value: _activityLevel,
            dropdownColor: AppTheme.surface,
            style: const TextStyle(color: AppTheme.textPrimary),
            decoration: const InputDecoration(
              prefixIcon: Icon(Icons.directions_run,
                  color: AppTheme.textSecondary, size: 20),
            ),
            items: const [
              DropdownMenuItem(
                  value: 'sedentary', child: Text('Sedentary')),
              DropdownMenuItem(
                  value: 'light', child: Text('Lightly Active')),
              DropdownMenuItem(
                  value: 'moderate', child: Text('Moderately Active')),
              DropdownMenuItem(value: 'active', child: Text('Active')),
              DropdownMenuItem(
                  value: 'very_active', child: Text('Very Active')),
            ],
            onChanged: (v) => setState(() => _activityLevel = v!),
          ),
          const SizedBox(height: 32),
          Row(
            children: [
              TextButton(
                onPressed: _previousPage,
                child: const Text('Back'),
              ),
              const Spacer(),
              ElevatedButton(
                onPressed: _submitting ? null : _submitOnboarding,
                child: _submitting
                    ? const SizedBox(
                        width: 24,
                        height: 24,
                        child: CircularProgressIndicator(
                          color: Colors.white,
                          strokeWidth: 2.5,
                        ),
                      )
                    : const Text('Calculate My Targets'),
              ),
            ],
          ),
          const SizedBox(height: 24),
        ],
      ),
    );
  }

  // ─── Step 5: Results ────────────────────────────────────────────────────────

  Widget _buildResultsPage() {
    final calorieGoal = _results?['calorie_goal'] ?? _calculateTDEE();
    final proteinGoal = _results?['protein_goal_g'] ?? (calorieGoal * 0.3 / 4).round();
    final carbsGoal = _results?['carbs_goal_g'] ?? (calorieGoal * 0.4 / 4).round();
    final fatGoal = _results?['fat_goal_g'] ?? (calorieGoal * 0.3 / 9).round();

    return LayoutBuilder(
      builder: (context, constraints) => SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: ConstrainedBox(
          constraints: BoxConstraints(minHeight: constraints.maxHeight),
          child: IntrinsicHeight(
            child: Column(
              children: [
                const SizedBox(height: 48),
                Container(
                  width: 80,
                  height: 80,
                  decoration: BoxDecoration(
                    color: AppTheme.primary.withAlpha(30),
                    shape: BoxShape.circle,
                  ),
                  child: const Icon(Icons.check_circle,
                      color: AppTheme.primary, size: 48),
                ),
                const SizedBox(height: 24),
                const Text(
                  "You're all set!",
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.bold,
                    color: AppTheme.textPrimary,
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  'Here are your personalized targets',
                  style: TextStyle(
                    fontSize: 16,
                    color: AppTheme.textSecondary,
                  ),
                ),
                const SizedBox(height: 40),
                // Calorie target card
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(24),
                    child: Column(
                      children: [
                        const Text(
                          'Daily Calorie Target',
                          style: TextStyle(
                            fontSize: 14,
                            color: AppTheme.textSecondary,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          '${_toInt(calorieGoal)}',
                          style: const TextStyle(
                            fontSize: 48,
                            fontWeight: FontWeight.bold,
                            color: AppTheme.primary,
                          ),
                        ),
                        const Text(
                          'kcal / day',
                          style: TextStyle(
                            fontSize: 16,
                            color: AppTheme.textSecondary,
                          ),
                        ),
                        const SizedBox(height: 24),
                        const Divider(),
                        const SizedBox(height: 16),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceAround,
                          children: [
                            _macroTarget(
                                'Protein', '${_toInt(proteinGoal)}g', Colors.blue),
                            _macroTarget(
                                'Carbs', '${_toInt(carbsGoal)}g', Colors.amber),
                            _macroTarget(
                                'Fat', '${_toInt(fatGoal)}g', Colors.pink[300]!),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
                const Spacer(),
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: _completeOnboarding,
                    child: const Text('Get Started'),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  int _toInt(dynamic value) {
    if (value is int) return value;
    if (value is double) return value.round();
    if (value is String) return double.tryParse(value)?.round() ?? 0;
    return 0;
  }

  Widget _macroTarget(String label, String value, Color color) {
    return Column(
      children: [
        Container(
          width: 10,
          height: 10,
          decoration: BoxDecoration(color: color, shape: BoxShape.circle),
        ),
        const SizedBox(height: 6),
        Text(
          value,
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: color,
          ),
        ),
        const SizedBox(height: 2),
        Text(
          label,
          style: const TextStyle(
            fontSize: 13,
            color: AppTheme.textSecondary,
          ),
        ),
      ],
    );
  }

  /// Fallback TDEE calculation if the server didn't return results
  int _calculateTDEE() {
    final weight = double.tryParse(_weightController.text) ?? 70;
    final height = double.tryParse(_heightController.text) ?? 170;
    final age = int.tryParse(_ageController.text) ?? 25;

    // Mifflin-St Jeor equation
    double bmr;
    if (_gender == 'female') {
      bmr = 10 * weight + 6.25 * height - 5 * age - 161;
    } else {
      bmr = 10 * weight + 6.25 * height - 5 * age + 5;
    }

    // Activity multiplier
    double multiplier;
    switch (_activityLevel) {
      case 'sedentary':
        multiplier = 1.2;
        break;
      case 'light':
        multiplier = 1.375;
        break;
      case 'moderate':
        multiplier = 1.55;
        break;
      case 'active':
        multiplier = 1.725;
        break;
      case 'very_active':
        multiplier = 1.9;
        break;
      default:
        multiplier = 1.55;
    }

    double tdee = bmr * multiplier;

    // Goal adjustment
    switch (_goal) {
      case 'lose_weight':
        tdee -= 500;
        break;
      case 'build_muscle':
        tdee += 300;
        break;
    }

    return max(1200, tdee.round());
  }
}
