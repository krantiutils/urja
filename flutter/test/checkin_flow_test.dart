import 'package:flutter_test/flutter_test.dart';

import 'package:urja_flutter/models/attendance.dart';
import 'package:urja_flutter/screens/checkin/qr_scanner_screen.dart';

import 'helpers/test_helpers.dart';

void main() {
  group('Check-in Flow', () {
    // Note: We cannot test the actual QR scanner in widget tests because
    // MobileScanner requires platform channels. Instead, we test the check-in
    // success view and verify the widget structure.

    testWidgets('QrScannerScreen renders appbar with Check In title',
        (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const QrScannerScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      // Pump once to see the initial frame (MobileScanner may error but
      // the scaffold should render)
      await tester.pump();

      // The AppBar title should appear
      expect(find.text('Check In'), findsWidgets);
    });

    testWidgets('CheckInResult model creates correctly', (tester) async {
      final result = CheckInResult.fromJson({
        'check_in': {
          'id': 'att-1',
          'user_id': 'u-1',
          'org_id': 'org-1',
          'check_in_at': '2024-06-15T08:30:00Z',
          'method': 'qr',
        },
        'streak': {
          'id': 's-1',
          'member_id': 'u-1',
          'org_id': 'org-1',
          'current_streak': 5,
          'longest_streak': 10,
          'last_check_in': '2024-06-15',
          'updated_at': '2024-06-15T10:00:00Z',
        },
      });

      expect(result.checkIn.method, 'qr');
      expect(result.streak, isNotNull);
      expect(result.streak!.currentStreak, 5);
    });

    testWidgets('CheckInResult without streak', (tester) async {
      final result = CheckInResult.fromJson({
        'check_in': {
          'id': 'att-1',
          'user_id': 'u-1',
          'org_id': 'org-1',
          'check_in_at': '2024-06-15T08:30:00Z',
          'method': 'nfc',
        },
      });

      expect(result.checkIn.method, 'nfc');
      expect(result.streak, isNull);
    });
  });
}
