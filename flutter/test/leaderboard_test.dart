import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:urja_flutter/screens/leaderboard/leaderboard_screen.dart';

import 'helpers/test_helpers.dart';

void main() {
  group('Leaderboard', () {
    testWidgets('renders with three tabs', (tester) async {
      await tester.pumpWidget(buildTestApp(const LeaderboardScreen()));
      await tester.pumpAndSettle();

      expect(find.text('Weekly'), findsOneWidget);
      expect(find.text('Monthly'), findsOneWidget);
      expect(find.text('All Time'), findsOneWidget);
    });

    testWidgets('shows leaderboard title in appbar', (tester) async {
      await tester.pumpWidget(buildTestApp(const LeaderboardScreen()));
      await tester.pumpAndSettle();

      expect(find.text('Leaderboard'), findsOneWidget);
    });

    testWidgets('shows placeholder content', (tester) async {
      await tester.pumpWidget(buildTestApp(const LeaderboardScreen()));
      await tester.pumpAndSettle();

      expect(find.text('No data available'), findsOneWidget);
      expect(find.byIcon(Icons.leaderboard), findsOneWidget);
    });

    testWidgets('can switch tabs', (tester) async {
      await tester.pumpWidget(buildTestApp(const LeaderboardScreen()));
      await tester.pumpAndSettle();

      // Tap Monthly tab
      await tester.tap(find.text('Monthly'));
      await tester.pumpAndSettle();

      // Still shows placeholder (API is TODO)
      expect(find.text('No data available'), findsOneWidget);

      // Tap All Time tab
      await tester.tap(find.text('All Time'));
      await tester.pumpAndSettle();

      expect(find.text('No data available'), findsOneWidget);
    });

    testWidgets('golden — leaderboard', (tester) async {
      await setGoldenSurfaceSize(tester);
      await tester.pumpWidget(buildTestApp(const LeaderboardScreen()));
      await tester.pumpAndSettle();

      await expectLater(
        find.byType(LeaderboardScreen),
        matchesGoldenFile('goldens/leaderboard.png'),
      );
    });
  });
}
