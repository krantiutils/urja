import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:urja_flutter/screens/health/health_screen.dart';

import 'helpers/test_helpers.dart';

void main() {
  group('Health Metrics', () {
    testWidgets('renders BMI calculator form', (tester) async {
      await tester.pumpWidget(buildTestApp(const HealthScreen()));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('height_field')), findsOneWidget);
      expect(find.byKey(const Key('weight_field')), findsOneWidget);
      expect(find.byKey(const Key('log_bmi_button')), findsOneWidget);
    });

    testWidgets('renders health tiles', (tester) async {
      await tester.pumpWidget(buildTestApp(const HealthScreen()));
      await tester.pumpAndSettle();

      expect(find.text('Body Measurements'), findsOneWidget);
      expect(find.text('Progress Photos'), findsOneWidget);
      expect(find.text('Weight'), findsOneWidget);
    });

    testWidgets('calculates BMI correctly', (tester) async {
      await tester.pumpWidget(buildTestApp(const HealthScreen()));
      await tester.pumpAndSettle();

      // Enter height 170cm
      await tester.enterText(
          find.byKey(const Key('height_field')), '170');
      // Enter weight 70kg
      await tester.enterText(
          find.byKey(const Key('weight_field')), '70');

      await tester.tap(find.byKey(const Key('log_bmi_button')));
      await tester.pumpAndSettle();

      // BMI = 70 / (1.7 * 1.7) = 24.2
      expect(find.text('BMI: 24.2'), findsOneWidget);
    });

    testWidgets('shows health tiles with chevrons', (tester) async {
      await tester.pumpWidget(buildTestApp(const HealthScreen()));
      await tester.pumpAndSettle();

      // 3 health tiles, each with a chevron
      expect(find.byIcon(Icons.chevron_right), findsNWidgets(3));
    });

    testWidgets('golden — health screen', (tester) async {
      await setGoldenSurfaceSize(tester);
      await tester.pumpWidget(buildTestApp(const HealthScreen()));
      await tester.pumpAndSettle();

      await expectLater(
        find.byType(HealthScreen),
        matchesGoldenFile('goldens/health_screen.png'),
      );
    });
  });
}
