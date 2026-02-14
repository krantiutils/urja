import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:urja_flutter/screens/profile/profile_screen.dart';

import 'helpers/test_helpers.dart';

void main() {
  group('Profile Edit', () {
    testWidgets('renders profile view with member data', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const ProfileScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Ram Sharma'), findsOneWidget);
      expect(find.text('9841234567'), findsOneWidget);
      expect(find.text('ram@example.com'), findsOneWidget);
    });

    testWidgets('shows avatar with first letter of name', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const ProfileScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('R'), findsOneWidget);
    });

    testWidgets('shows organization membership', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const ProfileScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Urja Fitness'), findsOneWidget);
      expect(find.text('member - active'), findsOneWidget);
    });

    testWidgets('edit button opens edit form', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const ProfileScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.edit));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('name_field')), findsOneWidget);
      expect(find.byKey(const Key('email_field')), findsOneWidget);
      expect(find.byKey(const Key('save_profile_button')), findsOneWidget);
    });

    testWidgets('edit form pre-fills current data', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const ProfileScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.edit));
      await tester.pumpAndSettle();

      // Name field should contain 'Ram Sharma' (pre-filled from profile)
      final editableTexts = tester
          .widgetList<EditableText>(find.byType(EditableText))
          .toList();
      final nameText = editableTexts.firstWhere(
        (e) => e.controller.text == 'Ram Sharma',
        orElse: () => throw StateError(
          'No EditableText with "Ram Sharma" found. '
          'Values: ${editableTexts.map((e) => e.controller.text).toList()}',
        ),
      );
      expect(nameText.controller.text, 'Ram Sharma');
    });

    testWidgets('cancel button returns to view mode', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const ProfileScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      // Enter edit mode
      await tester.tap(find.byIcon(Icons.edit));
      await tester.pumpAndSettle();
      expect(find.byKey(const Key('name_field')), findsOneWidget);

      // Cancel
      await tester.tap(find.text('Cancel'));
      await tester.pumpAndSettle();
      expect(find.byKey(const Key('name_field')), findsNothing);
      expect(find.text('Ram Sharma'), findsOneWidget);
    });

    testWidgets('logout button is visible after scrolling', (tester) async {
      await tester.pumpWidget(
        buildTestApp(
          const ProfileScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      // Logout button is at the bottom — scroll to it
      await tester.scrollUntilVisible(
        find.byKey(const Key('logout_button')),
        200.0,
      );
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('logout_button')), findsOneWidget);
    });

    testWidgets('golden — profile view', (tester) async {
      await setGoldenSurfaceSize(tester);
      await tester.pumpWidget(
        buildTestApp(
          const ProfileScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      await expectLater(
        find.byType(ProfileScreen),
        matchesGoldenFile('goldens/profile_view.png'),
      );
    });

    testWidgets('golden — profile edit form', (tester) async {
      await setGoldenSurfaceSize(tester);
      await tester.pumpWidget(
        buildTestApp(
          const ProfileScreen(),
          overrides: authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.edit));
      await tester.pumpAndSettle();

      await expectLater(
        find.byType(ProfileScreen),
        matchesGoldenFile('goldens/profile_edit.png'),
      );
    });
  });
}
