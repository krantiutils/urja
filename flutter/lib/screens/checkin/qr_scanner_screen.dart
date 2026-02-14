import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile_scanner/mobile_scanner.dart';
import 'package:urja_flutter/l10n/app_localizations.dart';

import '../../config/theme.dart';
import '../../core/api/api_exceptions.dart';
import '../../models/attendance.dart';
import '../../providers/providers.dart';

class QrScannerScreen extends ConsumerStatefulWidget {
  const QrScannerScreen({super.key});

  @override
  ConsumerState<QrScannerScreen> createState() => _QrScannerScreenState();
}

class _QrScannerScreenState extends ConsumerState<QrScannerScreen> {
  MobileScannerController? _scannerController;
  bool _processing = false;
  CheckInResult? _result;
  String? _error;

  @override
  void initState() {
    super.initState();
    _scannerController = MobileScannerController();
  }

  @override
  void dispose() {
    _scannerController?.dispose();
    super.dispose();
  }

  Future<void> _handleBarcode(BarcodeCapture capture) async {
    if (_processing || _result != null) return;

    final barcodes = capture.barcodes;
    if (barcodes.isEmpty) return;

    final qrToken = barcodes.first.rawValue;
    if (qrToken == null || qrToken.isEmpty) return;

    setState(() {
      _processing = true;
      _error = null;
    });

    try {
      final isOnline = ref.read(isOnlineProvider);

      if (!isOnline) {
        // Queue for offline sync
        final db = ref.read(localDbProvider);
        await db.queueCheckIn(method: 'qr', qrToken: qrToken);
        if (mounted) {
          setState(() {
            _processing = false;
          });
          _showOfflineSnackbar();
        }
        return;
      }

      final api = ref.read(apiClientProvider);
      final response = await api.selfCheckIn({
        'method': 'qr',
        'qr_token': qrToken,
      });
      final result =
          CheckInResult.fromJson(response.data as Map<String, dynamic>);

      setState(() {
        _result = result;
        _processing = false;
      });

      // Refresh attendance and streaks
      ref.invalidate(attendanceProvider);
      ref.invalidate(streaksProvider);
    } on ApiException catch (e) {
      setState(() {
        _error = e.message;
        _processing = false;
      });
    } on NetworkException {
      // Queue for offline
      final db = ref.read(localDbProvider);
      await db.queueCheckIn(method: 'qr', qrToken: qrToken);
      if (mounted) {
        setState(() => _processing = false);
        _showOfflineSnackbar();
      }
    } catch (e) {
      setState(() {
        _error = 'Check-in failed';
        _processing = false;
      });
    }
  }

  void _showOfflineSnackbar() {
    final l10n = AppLocalizations.of(context)!;
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(l10n.offlineMode),
        backgroundColor: UrjaColors.warning,
      ),
    );
  }

  void _reset() {
    setState(() {
      _result = null;
      _error = null;
      _processing = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    if (_result != null) {
      return _CheckInSuccessView(result: _result!, l10n: l10n, onReset: _reset);
    }

    return Scaffold(
      appBar: AppBar(title: Text(l10n.checkIn)),
      body: Column(
        children: [
          Expanded(
            child: Stack(
              children: [
                MobileScanner(
                  key: const Key('qr_scanner'),
                  controller: _scannerController,
                  onDetect: _handleBarcode,
                ),
                // Overlay
                Center(
                  child: Container(
                    width: 250,
                    height: 250,
                    decoration: BoxDecoration(
                      border: Border.all(
                          color: UrjaColors.accent, width: 2),
                      borderRadius: BorderRadius.circular(16),
                    ),
                  ),
                ),
                if (_processing)
                  const Center(child: CircularProgressIndicator()),
              ],
            ),
          ),
          if (_error != null)
            Container(
              width: double.infinity,
              padding: const EdgeInsets.all(16),
              color: UrjaColors.error.withAlpha(25),
              child: Column(
                children: [
                  Text(_error!,
                      style: const TextStyle(color: UrjaColors.error)),
                  const SizedBox(height: 8),
                  TextButton(
                    onPressed: _reset,
                    child: Text(l10n.retry),
                  ),
                ],
              ),
            ),
          Padding(
            padding: const EdgeInsets.all(24),
            child: Text(
              l10n.scanGymQr,
              style: Theme.of(context).textTheme.bodyLarge,
              textAlign: TextAlign.center,
            ),
          ),
        ],
      ),
    );
  }
}

class _CheckInSuccessView extends StatelessWidget {
  final CheckInResult result;
  final AppLocalizations l10n;
  final VoidCallback onReset;

  const _CheckInSuccessView({
    required this.result,
    required this.l10n,
    required this.onReset,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(l10n.checkIn)),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(32),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.check_circle,
                  size: 80, color: UrjaColors.success),
              const SizedBox(height: 24),
              Text(
                l10n.checkInSuccess,
                style: Theme.of(context).textTheme.headlineMedium,
              ),
              if (result.streak != null) ...[
                const SizedBox(height: 16),
                Text(
                  l10n.streakUpdated(result.streak!.currentStreak),
                  style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        color: UrjaColors.warning,
                      ),
                ),
              ],
              const SizedBox(height: 32),
              OutlinedButton(
                onPressed: onReset,
                child: Text(l10n.checkIn),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
