import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:urja_flutter/models/attendance.dart';
import 'package:urja_flutter/screens/attendance/attendance_screen.dart';

import 'helpers/test_helpers.dart';

void main() {
  group('Attendance History', () {
    testWidgets('shows "no attendance" when list is empty', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const AttendanceScreen(),
          overrides: authenticatedOverrides(attendance: []),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('No attendance records yet'), findsOneWidget);
      expect(find.byIcon(Icons.calendar_today), findsOneWidget);
    });

    testWidgets('golden — empty attendance', (tester) async {
      await setGoldenSurfaceSize(tester);
      await tester.pumpWidget(
        buildTestApp(
          const AttendanceScreen(),
          overrides: authenticatedOverrides(attendance: [], streaks: []),
        ),
      );
      await tester.pumpAndSettle();

      await expectLater(
        find.byType(AttendanceScreen),
        matchesGoldenFile('goldens/attendance_empty.png'),
      );
    });

    testWidgets('displays attendance records with method badges',
        (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const AttendanceScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      // 3 records in mock data
      expect(find.text('QR'), findsOneWidget);
      expect(find.text('NFC'), findsOneWidget);
      expect(find.text('Manual'), findsOneWidget);
    });

    testWidgets('shows correct method icons', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const AttendanceScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.qr_code), findsOneWidget);
      expect(find.byIcon(Icons.nfc), findsOneWidget);
      expect(find.byIcon(Icons.touch_app), findsOneWidget);
    });

    testWidgets('shows streak summary card', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const AttendanceScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      // mockStreaks has currentStreak=5, longestStreak=12
      expect(find.text('5'), findsOneWidget);
      expect(find.text('12'), findsOneWidget);
    });

    testWidgets('golden — attendance with records', (tester) async {
      await setGoldenSurfaceSize(tester);
      await tester.pumpWidget(
        buildTestApp(
          const AttendanceScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      await expectLater(
        find.byType(AttendanceScreen),
        matchesGoldenFile('goldens/attendance_with_records.png'),
      );
    });

    testWidgets('hides streak card when no streaks', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const AttendanceScreen(),
          overrides: authenticatedOverrides(streaks: []),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.local_fire_department), findsNothing);
      expect(find.byIcon(Icons.emoji_events), findsNothing);
    });

    testWidgets('shows date and time for each record', (tester) async {
      final records = [
        AttendanceRecord(
          id: 'att-x',
          userId: 'user-123',
          orgId: 'org-1',
          checkInAt: DateTime(2024, 12, 25, 10, 30),
          method: 'qr',
        ),
      ];

      await tester.pumpWidget(
        buildTestApp(
          const AttendanceScreen(),
          overrides: authenticatedOverrides(attendance: records),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Dec 25, 2024'), findsOneWidget);
      expect(find.text('10:30 AM'), findsOneWidget);
    });
  });
}
