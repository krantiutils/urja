import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:urja_flutter/l10n/app_localizations.dart';

import '../../config/theme.dart';
import '../../core/api/api_exceptions.dart';
import '../../providers/providers.dart';

class OtpScreen extends ConsumerStatefulWidget {
  final String phone;

  const OtpScreen({super.key, required this.phone});

  @override
  ConsumerState<OtpScreen> createState() => _OtpScreenState();
}

class _OtpScreenState extends ConsumerState<OtpScreen> {
  final _otpController = TextEditingController();
  final _formKey = GlobalKey<FormState>();
  bool _loading = false;
  String? _error;

  @override
  void dispose() {
    _otpController.dispose();
    super.dispose();
  }

  Future<void> _verifyOtp() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _loading = true;
      _error = null;
    });

    try {
      final api = ref.read(apiClientProvider);
      await api.verifyOtp(widget.phone, _otpController.text.trim());
      ref.read(authStateProvider.notifier).state = AuthState.authenticated;
      if (mounted) context.go('/');
    } on ApiException catch (e) {
      setState(() => _error = e.message);
    } on NetworkException catch (e) {
      setState(() => _error = e.message);
    } catch (e) {
      setState(() => _error = 'Verification failed');
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _resendOtp() async {
    try {
      final api = ref.read(apiClientProvider);
      await api.requestOtp(widget.phone);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
              content: Text(AppLocalizations.of(context)!.otpSent)),
        );
      }
    } catch (e) {
      // Ignore resend errors
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            child: Form(
              key: _formKey,
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Text(
                    l10n.verifyOtp,
                    style: Theme.of(context).textTheme.headlineMedium,
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 12),
                  Text(
                    l10n.enterOtp(widget.phone),
                    style: Theme.of(context).textTheme.bodyLarge,
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 40),

                  TextFormField(
                    key: const Key('otp_field'),
                    controller: _otpController,
                    keyboardType: TextInputType.number,
                    maxLength: 6,
                    textAlign: TextAlign.center,
                    style: const TextStyle(
                        fontSize: 24, letterSpacing: 12),
                    decoration: InputDecoration(
                      hintText: l10n.otpHint,
                      counterText: '',
                    ),
                    validator: (value) {
                      if (value == null ||
                          value.length != 6 ||
                          int.tryParse(value) == null) {
                        return l10n.invalidOtp;
                      }
                      return null;
                    },
                  ),
                  const SizedBox(height: 16),

                  if (_error != null) ...[
                    Container(
                      padding: const EdgeInsets.all(12),
                      decoration: BoxDecoration(
                        color: UrjaColors.error.withAlpha(25),
                        borderRadius: BorderRadius.circular(8),
                        border: Border.all(
                            color: UrjaColors.error.withAlpha(76)),
                      ),
                      child: Text(
                        _error!,
                        style: const TextStyle(
                            color: UrjaColors.error, fontSize: 14),
                        textAlign: TextAlign.center,
                      ),
                    ),
                    const SizedBox(height: 16),
                  ],

                  SizedBox(
                    height: 50,
                    child: ElevatedButton(
                      key: const Key('verify_button'),
                      onPressed: _loading ? null : _verifyOtp,
                      child: _loading
                          ? const SizedBox(
                              height: 20,
                              width: 20,
                              child: CircularProgressIndicator(
                                strokeWidth: 2,
                                color: Colors.white,
                              ),
                            )
                          : Text(l10n.verify),
                    ),
                  ),
                  const SizedBox(height: 16),
                  TextButton(
                    key: const Key('resend_button'),
                    onPressed: _resendOtp,
                    child: Text(l10n.resendOtp),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
