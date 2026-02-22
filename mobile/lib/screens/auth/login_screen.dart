import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';

class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});

  @override
  ConsumerState<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends ConsumerState<LoginScreen> {
  final _phoneController = TextEditingController();
  final _otpController = TextEditingController();
  bool _otpSent = false;
  bool _loading = false;
  int _resendCountdown = 0;
  Timer? _resendTimer;

  @override
  void dispose() {
    _phoneController.dispose();
    _otpController.dispose();
    _resendTimer?.cancel();
    super.dispose();
  }

  void _startResendTimer() {
    _resendCountdown = 60;
    _resendTimer?.cancel();
    _resendTimer = Timer.periodic(const Duration(seconds: 1), (timer) {
      setState(() {
        _resendCountdown--;
        if (_resendCountdown <= 0) {
          timer.cancel();
        }
      });
    });
  }

  Future<void> _sendOtp() async {
    final phone = _phoneController.text.trim();
    if (phone.isEmpty || phone.length < 10) {
      _showError('Please enter a valid phone number');
      return;
    }

    setState(() => _loading = true);
    try {
      await ref.read(authProvider.notifier).login(phone);
      setState(() {
        _otpSent = true;
        _loading = false;
      });
      _startResendTimer();
    } catch (e) {
      setState(() => _loading = false);
      _showError('Failed to send OTP. Please try again.');
    }
  }

  Future<void> _verifyOtp() async {
    final phone = _phoneController.text.trim();
    final otp = _otpController.text.trim();
    if (otp.length != 6) {
      _showError('Please enter the 6-digit code');
      return;
    }

    setState(() => _loading = true);
    try {
      await ref.read(authProvider.notifier).verifyOtp(phone, otp);
      // Navigation handled automatically by router redirect
    } catch (e) {
      setState(() => _loading = false);
      _showError('Invalid verification code. Please try again.');
    }
  }

  Future<void> _resendOtp() async {
    if (_resendCountdown > 0) return;
    await _sendOtp();
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

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      backgroundColor: AppTheme.background,
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                // Logo / brand icon
                Icon(
                  Icons.fitness_center,
                  size: 64,
                  color: AppTheme.primary,
                ),
                const SizedBox(height: 24),
                // Welcome text
                Text(
                  l10n.welcomeToUrja,
                  textAlign: TextAlign.center,
                  style: const TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.bold,
                    color: AppTheme.textPrimary,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  l10n.gymManagementSimplified,
                  textAlign: TextAlign.center,
                  style: const TextStyle(
                    fontSize: 16,
                    color: AppTheme.textSecondary,
                  ),
                ),
                const SizedBox(height: 48),
                // Phone input
                if (!_otpSent) ...[
                  _buildPhoneInput(l10n),
                  const SizedBox(height: 20),
                  _buildSendOtpButton(l10n),
                ] else ...[
                  // Show phone number with change option
                  _buildPhoneDisplay(),
                  const SizedBox(height: 24),
                  _buildOtpInput(l10n),
                  const SizedBox(height: 20),
                  _buildVerifyButton(l10n),
                  const SizedBox(height: 16),
                  _buildResendRow(l10n),
                ],
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildPhoneInput(AppLocalizations l10n) {
    return TextField(
      controller: _phoneController,
      keyboardType: TextInputType.phone,
      inputFormatters: [
        FilteringTextInputFormatter.digitsOnly,
        LengthLimitingTextInputFormatter(10),
      ],
      style: const TextStyle(color: AppTheme.textPrimary, fontSize: 16),
      decoration: InputDecoration(
        labelText: l10n.phoneNumber,
        labelStyle: const TextStyle(color: AppTheme.textSecondary),
        prefixIcon: Container(
          padding: const EdgeInsets.symmetric(horizontal: 12),
          child: const Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                '+977',
                style: TextStyle(
                  color: AppTheme.textPrimary,
                  fontSize: 16,
                  fontWeight: FontWeight.w500,
                ),
              ),
              SizedBox(width: 8),
              SizedBox(
                height: 24,
                child: VerticalDivider(
                  color: AppTheme.textSecondary,
                  thickness: 1,
                ),
              ),
            ],
          ),
        ),
        prefixIconConstraints: const BoxConstraints(minWidth: 0, minHeight: 0),
        hintText: '98XXXXXXXX',
        hintStyle: TextStyle(color: AppTheme.textSecondary.withAlpha(100)),
      ),
    );
  }

  Widget _buildSendOtpButton(AppLocalizations l10n) {
    return SizedBox(
      height: 52,
      child: ElevatedButton(
        onPressed: _loading ? null : _sendOtp,
        child: _loading
            ? const SizedBox(
                width: 24,
                height: 24,
                child: CircularProgressIndicator(
                  color: Colors.white,
                  strokeWidth: 2.5,
                ),
              )
            : Text(l10n.sendOtp),
      ),
    );
  }

  Widget _buildPhoneDisplay() {
    return Row(
      children: [
        Expanded(
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
            decoration: BoxDecoration(
              color: AppTheme.surfaceLight,
              borderRadius: BorderRadius.circular(12),
            ),
            child: Text(
              '+977 ${_phoneController.text}',
              style: const TextStyle(
                color: AppTheme.textPrimary,
                fontSize: 16,
              ),
            ),
          ),
        ),
        const SizedBox(width: 12),
        TextButton(
          onPressed: () {
            setState(() {
              _otpSent = false;
              _otpController.clear();
              _resendTimer?.cancel();
              _resendCountdown = 0;
            });
          },
          child: const Text('Change'),
        ),
      ],
    );
  }

  Widget _buildOtpInput(AppLocalizations l10n) {
    return TextField(
      controller: _otpController,
      keyboardType: TextInputType.number,
      textAlign: TextAlign.center,
      inputFormatters: [
        FilteringTextInputFormatter.digitsOnly,
        LengthLimitingTextInputFormatter(6),
      ],
      style: const TextStyle(
        color: AppTheme.textPrimary,
        fontSize: 24,
        fontWeight: FontWeight.bold,
        letterSpacing: 12,
      ),
      decoration: InputDecoration(
        labelText: l10n.enterOtp,
        labelStyle: const TextStyle(color: AppTheme.textSecondary),
        hintText: '------',
        hintStyle: TextStyle(
          color: AppTheme.textSecondary.withAlpha(100),
          letterSpacing: 12,
        ),
      ),
      onChanged: (value) {
        if (value.length == 6) {
          _verifyOtp();
        }
      },
    );
  }

  Widget _buildVerifyButton(AppLocalizations l10n) {
    return SizedBox(
      height: 52,
      child: ElevatedButton(
        onPressed: _loading ? null : _verifyOtp,
        child: _loading
            ? const SizedBox(
                width: 24,
                height: 24,
                child: CircularProgressIndicator(
                  color: Colors.white,
                  strokeWidth: 2.5,
                ),
              )
            : Text(l10n.verifyOtp),
      ),
    );
  }

  Widget _buildResendRow(AppLocalizations l10n) {
    return Wrap(
      alignment: WrapAlignment.center,
      crossAxisAlignment: WrapCrossAlignment.center,
      children: [
        Text(
          l10n.didNotReceiveCode,
          style: const TextStyle(color: AppTheme.textSecondary),
        ),
        const SizedBox(width: 4),
        if (_resendCountdown > 0)
          Text(
            l10n.resendIn(_resendCountdown),
            style: const TextStyle(color: AppTheme.textSecondary),
          )
        else
          TextButton(
            onPressed: _resendOtp,
            child: Text(l10n.resend),
          ),
      ],
    );
  }
}
