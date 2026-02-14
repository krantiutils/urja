import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:urja_flutter/screens/home/home_screen.dart';

import 'helpers/test_helpers.dart';

void main() {
  group('Offline Check-in & Indicator', () {
    testWidgets('shows offline icon on home when offline', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const HomeScreen(),
          overrides: authenticatedOverrides(online: false),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.cloud_off), findsOneWidget);
    });

    testWidgets('hides offline icon when online', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const HomeScreen(),
          overrides: authenticatedOverrides(online: true),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.cloud_off), findsNothing);
    });

    testWidgets('golden — home offline state', (tester) async {
      await setGoldenSurfaceSize(tester);
      await tester.pumpWidget(
        buildTestApp(
          const HomeScreen(),
          overrides: authenticatedOverrides(online: false),
        ),
      );
      await tester.pumpAndSettle();

      await expectLater(
        find.byType(HomeScreen),
        matchesGoldenFile('goldens/home_offline.png'),
      );
    });

    testWidgets('home screen renders notification and settings icons',
        (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const HomeScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.notifications_outlined), findsOneWidget);
      expect(find.byIcon(Icons.settings_outlined), findsOneWidget);
    });
  });
}
