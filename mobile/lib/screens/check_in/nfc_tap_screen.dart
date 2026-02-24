import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:nfc_manager/nfc_manager.dart';
import 'package:nfc_manager/nfc_manager_android.dart';

import '../../config/theme.dart';

class NfcTapScreen extends ConsumerStatefulWidget {
  const NfcTapScreen({super.key});

  @override
  ConsumerState<NfcTapScreen> createState() => _NfcTapScreenState();
}

class _NfcTapScreenState extends ConsumerState<NfcTapScreen>
    with SingleTickerProviderStateMixin {
  late AnimationController _pulseController;
  late Animation<double> _pulseAnimation;
  bool _nfcAvailable = false;
  String? _scannedUid;
  String? _error;

  @override
  void initState() {
    super.initState();
    _pulseController = AnimationController(
      duration: const Duration(milliseconds: 1500),
      vsync: this,
    )..repeat(reverse: true);

    _pulseAnimation = Tween<double>(begin: 0.8, end: 1.2).animate(
      CurvedAnimation(parent: _pulseController, curve: Curves.easeInOut),
    );

    _checkNfc();
  }

  @override
  void dispose() {
    _pulseController.dispose();
    NfcManager.instance.stopSession();
    super.dispose();
  }

  Future<void> _checkNfc() async {
    try {
      final availability = await NfcManager.instance.checkAvailability();
      final available = availability == NfcAvailability.enabled;
      if (mounted) {
        setState(() => _nfcAvailable = available);
        if (available) {
          _startNfcSession();
        }
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _nfcAvailable = false;
          _error = e.toString();
        });
      }
    }
  }

  void _startNfcSession() {
    setState(() {
      _scannedUid = null;
      _error = null;
    });

    NfcManager.instance.startSession(
      pollingOptions: {NfcPollingOption.iso14443, NfcPollingOption.iso15693},
      onDiscovered: (NfcTag tag) async {
        String uid = '';

        if (!kIsWeb) {
          final androidTag = NfcTagAndroid.from(tag);
          if (androidTag != null) {
            uid = _bytesToHex(androidTag.id);
          }
        }

        if (mounted) {
          setState(() {
            _scannedUid = uid.isNotEmpty ? uid : 'Unknown';
          });
        }

        NfcManager.instance.stopSession();
      },
    );
  }

  String _bytesToHex(List<int> bytes) {
    return bytes
        .map((b) => b.toRadixString(16).padLeft(2, '0').toUpperCase())
        .join(':');
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
        title: Text(l10n.tapNfcCard),
        backgroundColor: AppTheme.surface,
      ),
      body: Container(
        width: double.infinity,
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            if (!_nfcAvailable && _error == null) ...[
              const Icon(Icons.nfc, size: 80, color: AppTheme.textSecondary),
              const SizedBox(height: 24),
              Text(l10n.nfcNotAvailable,
                  style: const TextStyle(
                      fontSize: 18, color: AppTheme.textSecondary),
                  textAlign: TextAlign.center),
            ] else if (_scannedUid != null) ...[
              const Icon(Icons.check_circle, size: 80, color: AppTheme.primary),
              const SizedBox(height: 24),
              Text(l10n.nfcCardDetected,
                  style: const TextStyle(
                      fontSize: 20, fontWeight: FontWeight.bold)),
              const SizedBox(height: 16),
              Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: AppTheme.surfaceLight,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Column(
                  children: [
                    Text(l10n.cardUid,
                        style: const TextStyle(
                            color: AppTheme.textSecondary, fontSize: 13)),
                    const SizedBox(height: 8),
                    Text(_scannedUid!,
                        style: const TextStyle(
                          fontFamily: 'monospace',
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                          color: AppTheme.primary,
                        )),
                  ],
                ),
              ),
              const SizedBox(height: 32),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  OutlinedButton(
                      onPressed: _startNfcSession,
                      child: Text(l10n.scanAgain)),
                  const SizedBox(width: 16),
                  ElevatedButton(
                      onPressed: () => context.pop(),
                      child: Text(l10n.done)),
                ],
              ),
            ] else if (_error != null) ...[
              const Icon(Icons.error_outline, size: 80, color: AppTheme.error),
              const SizedBox(height: 24),
              Text(l10n.nfcError,
                  style:
                      const TextStyle(fontSize: 18, color: AppTheme.error)),
              const SizedBox(height: 8),
              Text(_error!,
                  style: const TextStyle(color: AppTheme.textSecondary),
                  textAlign: TextAlign.center),
              const SizedBox(height: 24),
              ElevatedButton(
                  onPressed: _startNfcSession, child: Text(l10n.tryAgain)),
            ] else ...[
              AnimatedBuilder(
                animation: _pulseAnimation,
                builder: (context, child) {
                  return Transform.scale(
                    scale: _pulseAnimation.value,
                    child: Container(
                      width: 120,
                      height: 120,
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        color: AppTheme.primary.withAlpha(20),
                        border: Border.all(
                            color: AppTheme.primary.withAlpha(80), width: 2),
                      ),
                      child: const Icon(Icons.nfc,
                          size: 60, color: AppTheme.primary),
                    ),
                  );
                },
              ),
              const SizedBox(height: 40),
              Text(l10n.holdPhoneNearCard,
                  style: const TextStyle(
                      fontSize: 20, fontWeight: FontWeight.bold),
                  textAlign: TextAlign.center),
              const SizedBox(height: 12),
              Text(l10n.nfcInstructions,
                  style: const TextStyle(
                      color: AppTheme.textSecondary, fontSize: 14),
                  textAlign: TextAlign.center),
              const SizedBox(height: 32),
              const CircularProgressIndicator(
                  color: AppTheme.primary, strokeWidth: 2),
            ],
          ],
        ),
      ),
    );
  }
}
