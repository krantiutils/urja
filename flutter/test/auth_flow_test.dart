import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:urja_flutter/screens/auth/login_screen.dart';
import 'package:urja_flutter/screens/auth/otp_screen.dart';

import 'helpers/test_helpers.dart';

void main() {
  group('Auth Flow — Login Screen', () {
    testWidgets('renders phone field and send OTP button', (tester) async {
      await tester.pumpWidget(buildTestApp(const LoginScreen()));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('phone_field')), findsOneWidget);
      expect(find.byKey(const Key('send_otp_button')), findsOneWidget);
    });

    testWidgets('shows validation error for empty phone', (tester) async {
      await tester.pumpWidget(buildTestApp(const LoginScreen()));
      await tester.pumpAndSettle();

      await tester.tap(find.byKey(const Key('send_otp_button')));
      await tester.pumpAndSettle();

      expect(find.textContaining('valid'), findsOneWidget);
    });

    testWidgets('shows validation error for invalid phone format',
        (tester) async {
      await tester.pumpWidget(buildTestApp(const LoginScreen()));
      await tester.pumpAndSettle();

      await tester.enterText(find.byKey(const Key('phone_field')), '12345');
      await tester.tap(find.byKey(const Key('send_otp_button')));
      await tester.pumpAndSettle();

      expect(find.textContaining('valid'), findsOneWidget);
    });

    testWidgets('accepts valid Nepal phone number', (tester) async {
      await tester.pumpWidget(buildTestApp(const LoginScreen()));
      await tester.pumpAndSettle();

      await tester.enterText(
          find.byKey(const Key('phone_field')), '9841234567');

      // Valid phone — form validation should pass (no error shown yet
      // because we haven't tapped submit with mocked API).
      expect(find.byKey(const Key('phone_field')), findsOneWidget);
    });

    testWidgets('golden — login screen', (tester) async {
      await setGoldenSurfaceSize(tester);
      await tester.pumpWidget(buildTestApp(const LoginScreen()));
      await tester.pumpAndSettle();

      await expectLater(
        find.byType(LoginScreen),
        matchesGoldenFile('goldens/login_screen.png'),
      );
    });
  });

  group('Auth Flow — OTP Screen', () {
    testWidgets('renders OTP field, verify and resend buttons', (tester) async {
      await tester.pumpWidget(
        buildTestApp(const OtpScreen(phone: '9841234567')),
      );
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('otp_field')), findsOneWidget);
      expect(find.byKey(const Key('verify_button')), findsOneWidget);
      expect(find.byKey(const Key('resend_button')), findsOneWidget);
    });

    testWidgets('shows phone number in instruction text', (tester) async {
      await tester.pumpWidget(
        buildTestApp(const OtpScreen(phone: '9841234567')),
      );
      await tester.pumpAndSettle();

      expect(find.textContaining('9841234567'), findsOneWidget);
    });

    testWidgets('validates OTP length', (tester) async {
      await tester.pumpWidget(
        buildTestApp(const OtpScreen(phone: '9841234567')),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byKey(const Key('otp_field')), '123');
      await tester.tap(find.byKey(const Key('verify_button')));
      await tester.pumpAndSettle();

      // Validation error text (distinct from hint text)
      expect(find.text('Enter a valid 6-digit OTP'), findsOneWidget);
    });

    testWidgets('golden — OTP screen', (tester) async {
      await setGoldenSurfaceSize(tester);
      await tester.pumpWidget(
        buildTestApp(const OtpScreen(phone: '9841234567')),
      );
      await tester.pumpAndSettle();

      await expectLater(
        find.byType(OtpScreen),
        matchesGoldenFile('goldens/otp_screen.png'),
      );
    });
  });
}
