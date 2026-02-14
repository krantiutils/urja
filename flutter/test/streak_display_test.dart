import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:urja_flutter/models/attendance.dart';
import 'package:urja_flutter/screens/home/home_screen.dart';

import 'helpers/test_helpers.dart';

void main() {
  group('Streak Display', () {
    testWidgets('shows streak count on home screen', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const HomeScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      // mockStreaks has currentStreak=5
      expect(find.textContaining('5'), findsWidgets);
      expect(find.text('Current Streak'), findsOneWidget);
    });

    testWidgets('shows fire icon for active streak', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const HomeScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.local_fire_department), findsOneWidget);
    });

    testWidgets('shows zero streak when no streaks', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const HomeScreen(),
          overrides: authenticatedOverrides(streaks: []),
        ),
      );
      await tester.pumpAndSettle();

      // With empty streaks, current shows "0 days"
      expect(find.text('0 days'), findsOneWidget);
    });

    testWidgets('shows high streak value correctly', (tester) async {
      final bigStreak = [
        Streak(
          id: 's-1',
          memberId: 'user-123',
          orgId: 'org-1',
          currentStreak: 100,
          longestStreak: 200,
          lastCheckIn: '2024-06-15',
          updatedAt: DateTime(2024, 6, 15),
        ),
      ];

      await tester.pumpWidget(
        buildTestApp(
          const HomeScreen(),
          overrides: authenticatedOverrides(streaks: bigStreak),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('100 days'), findsOneWidget);
    });

    testWidgets('golden — home with streak', (tester) async {
      await setGoldenSurfaceSize(tester);
      await tester.pumpWidget(
        buildTestApp(
          const HomeScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      await expectLater(
        find.byType(HomeScreen),
        matchesGoldenFile('goldens/home_with_streak.png'),
      );
    });
  });
}
