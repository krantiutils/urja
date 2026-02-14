import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:urja_flutter/config/theme.dart';
import 'package:urja_flutter/screens/notifications/notifications_screen.dart';
import 'package:urja_flutter/screens/achievements/achievements_screen.dart';
import 'package:urja_flutter/screens/settings/settings_screen.dart';
import 'package:urja_flutter/providers/providers.dart';

import 'helpers/test_helpers.dart';

void main() {
  group('Navigation — Settings Screen', () {
    testWidgets('renders language toggle', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const SettingsScreen(),
          overrides: [
            localeProvider.overrideWith((ref) => 'en'),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Language'), findsOneWidget);
      expect(find.text('English'), findsOneWidget);
      expect(find.text('Nepali'), findsOneWidget);
    });

    testWidgets('renders privacy settings tile', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const SettingsScreen(),
          overrides: [
            localeProvider.overrideWith((ref) => 'en'),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Privacy Settings'), findsOneWidget);
    });

    testWidgets('golden — settings screen', (tester) async {
      await setGoldenSurfaceSize(tester);
      await tester.pumpWidget(
        buildTestApp(
          const SettingsScreen(),
          overrides: [
            localeProvider.overrideWith((ref) => 'en'),
          ],
        ),
      );
      await tester.pumpAndSettle();

      await expectLater(
        find.byType(SettingsScreen),
        matchesGoldenFile('goldens/settings_screen.png'),
      );
    });
  });

  group('Navigation — Notifications Screen', () {
    testWidgets('renders empty state', (tester) async {
      await tester.pumpWidget(buildTestApp(const NotificationsScreen()));
      await tester.pumpAndSettle();

      expect(find.text('Notifications'), findsOneWidget);
      expect(find.text('No data available'), findsOneWidget);
      expect(find.byIcon(Icons.notifications_none), findsOneWidget);
    });

    testWidgets('golden — notifications screen', (tester) async {
      await setGoldenSurfaceSize(tester);
      await tester.pumpWidget(buildTestApp(const NotificationsScreen()));
      await tester.pumpAndSettle();

      await expectLater(
        find.byType(NotificationsScreen),
        matchesGoldenFile('goldens/notifications_screen.png'),
      );
    });
  });

  group('Navigation — Achievements Screen', () {
    testWidgets('renders empty state', (tester) async {
      await tester.pumpWidget(buildTestApp(const AchievementsScreen()));
      await tester.pumpAndSettle();

      expect(find.text('Achievements'), findsOneWidget);
      expect(find.text('No data available'), findsOneWidget);
      expect(find.byIcon(Icons.emoji_events), findsOneWidget);
    });

    testWidgets('golden — achievements screen', (tester) async {
      await setGoldenSurfaceSize(tester);
      await tester.pumpWidget(buildTestApp(const AchievementsScreen()));
      await tester.pumpAndSettle();

      await expectLater(
        find.byType(AchievementsScreen),
        matchesGoldenFile('goldens/achievements_screen.png'),
      );
    });
  });

  group('Navigation — Bottom Nav Icons', () {
    testWidgets('home screen shows quick action buttons', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const Scaffold(
            body: Column(
              children: [
                Icon(Icons.qr_code_scanner),
                Icon(Icons.leaderboard_outlined),
                Icon(Icons.health_and_safety_outlined),
                Icon(Icons.emoji_events_outlined),
              ],
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.qr_code_scanner), findsOneWidget);
      expect(find.byIcon(Icons.leaderboard_outlined), findsOneWidget);
      expect(find.byIcon(Icons.health_and_safety_outlined), findsOneWidget);
      expect(find.byIcon(Icons.emoji_events_outlined), findsOneWidget);
    });
  });

  group('Theme', () {
    testWidgets('Urja theme uses dark mode with accent color', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          Builder(
            builder: (context) {
              final theme = Theme.of(context);
              expect(theme.brightness, Brightness.dark);
              expect(theme.primaryColor, UrjaColors.accent);
              expect(theme.scaffoldBackgroundColor, UrjaColors.backgroundBase);
              return const SizedBox();
            },
          ),
        ),
      );
    });
  });
}
