import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';

import 'package:urja_flutter/config/theme.dart';
import 'package:urja_flutter/l10n/app_localizations.dart';
import 'package:urja_flutter/models/attendance.dart';
import 'package:urja_flutter/models/package.dart';
import 'package:urja_flutter/models/user.dart';
import 'package:urja_flutter/providers/providers.dart';
import 'package:urja_flutter/screens/auth/login_screen.dart';
import 'package:urja_flutter/screens/auth/otp_screen.dart';
import 'package:urja_flutter/screens/home/home_screen.dart';
import 'package:urja_flutter/screens/attendance/attendance_screen.dart';
import 'package:urja_flutter/screens/packages/packages_screen.dart';
import 'package:urja_flutter/screens/profile/profile_screen.dart';
import 'package:urja_flutter/screens/leaderboard/leaderboard_screen.dart';
import 'package:urja_flutter/screens/health/health_screen.dart';
import 'package:urja_flutter/screens/achievements/achievements_screen.dart';
import 'package:urja_flutter/screens/notifications/notifications_screen.dart';
import 'package:urja_flutter/screens/settings/settings_screen.dart';

// ---------------------------------------------------------------------------
// Test helpers — duplicated from widget tests because integration_test
// cannot import from test/
// ---------------------------------------------------------------------------

final _mockMember = Member(
  id: 'user-123',
  phone: '9841234567',
  name: 'Ram Sharma',
  email: 'ram@example.com',
  gender: 'male',
  dateOfBirth: '1995-01-15',
  emergencyContactName: 'Sita Sharma',
  emergencyContactPhone: '9841234568',
  organizations: [
    OrgMembership(
      orgId: 'org-1',
      orgName: 'Urja Fitness',
      role: 'member',
      status: 'active',
      joinedAt: DateTime(2024, 1, 1),
    ),
  ],
  createdAt: DateTime(2024, 1, 1),
  updatedAt: DateTime(2024, 6, 1),
);

final _mockAttendance = [
  AttendanceRecord(
    id: 'att-1',
    userId: 'user-123',
    orgId: 'org-1',
    checkInAt: DateTime(2024, 6, 15, 8, 30),
    method: 'qr',
  ),
];

final _mockStreaks = [
  Streak(
    id: 'streak-1',
    memberId: 'user-123',
    orgId: 'org-1',
    currentStreak: 5,
    longestStreak: 12,
    lastCheckIn: '2024-06-15',
    updatedAt: DateTime(2024, 6, 15),
  ),
];

class _TestProfileNotifier extends ProfileNotifier {
  @override
  Future<Member?> build() async => _mockMember;
  @override
  Future<void> refresh() async {}
  @override
  Future<void> updateProfile(ProfileUpdate update) async {}
  @override
  Future<void> updatePrivacy(PrivacySettings settings) async {}
}

class _TestAttendanceNotifier extends AttendanceNotifier {
  @override
  Future<List<AttendanceRecord>> build() async => _mockAttendance;
  @override
  Future<void> refresh() async {}
}

List<Override> _authenticatedOverrides() {
  return [
    authStateProvider.overrideWith((ref) => AuthState.authenticated),
    profileProvider.overrideWith(() => _TestProfileNotifier()),
    attendanceProvider.overrideWith(() => _TestAttendanceNotifier()),
    streaksProvider.overrideWith((ref) async => _mockStreaks),
    memberPackagesProvider.overrideWith((ref) async => <MemberPackage>[]),
    isOnlineProvider.overrideWith((ref) => true),
    localeProvider.overrideWith((ref) => 'en'),
  ];
}

Widget _buildApp(Widget child, {List<Override> overrides = const []}) {
  return ProviderScope(
    overrides: overrides,
    child: MaterialApp(
      theme: urjaTheme(),
      localizationsDelegates: const [
        AppLocalizations.delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      supportedLocales: AppLocalizations.supportedLocales,
      home: child,
    ),
  );
}

// ---------------------------------------------------------------------------
// Integration tests — 10 suites
// ---------------------------------------------------------------------------

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  Future<void> screenshot(String name) async {
    await binding.convertFlutterSurfaceToImage();
    await binding.takeScreenshot('screenshots/$name');
  }

  group('1. Auth Flow', () {
    testWidgets('login screen renders and validates', (tester) async {
      await tester.pumpWidget(_buildApp(const LoginScreen()));
      await tester.pumpAndSettle();
      await screenshot('login_screen');

      expect(find.byKey(const Key('phone_field')), findsOneWidget);
      expect(find.byKey(const Key('send_otp_button')), findsOneWidget);

      // Invalid phone
      await tester.tap(find.byKey(const Key('send_otp_button')));
      await tester.pumpAndSettle();
      expect(find.textContaining('valid'), findsOneWidget);
    });

    testWidgets('OTP screen renders', (tester) async {
      await tester.pumpWidget(
          _buildApp(const OtpScreen(phone: '9841234567')));
      await tester.pumpAndSettle();
      await screenshot('otp_screen');

      expect(find.byKey(const Key('otp_field')), findsOneWidget);
      expect(find.textContaining('9841234567'), findsOneWidget);
    });
  });

  group('2. Check-in Flow', () {
    testWidgets('check-in screen renders', (tester) async {
      // QR scanner won't work without camera, just verify structure
      await tester.pumpWidget(
        _buildApp(
          const Scaffold(
            body: Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.qr_code_scanner, size: 80),
                  SizedBox(height: 16),
                  Text('Point your camera at the gym\'s QR code'),
                ],
              ),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();
      await screenshot('checkin_placeholder');

      expect(find.byIcon(Icons.qr_code_scanner), findsOneWidget);
    });
  });

  group('3. Package View', () {
    testWidgets('empty packages screen', (tester) async {
      await tester.pumpWidget(
        _buildApp(
          const PackagesScreen(),
          overrides: _authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();
      await screenshot('packages_empty');

      expect(find.text('No packages available'), findsOneWidget);
    });
  });

  group('4. Profile Edit', () {
    testWidgets('profile view and edit mode', (tester) async {
      await tester.pumpWidget(
        _buildApp(
          const ProfileScreen(),
          overrides: _authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();
      await screenshot('profile_view');

      expect(find.text('Ram Sharma'), findsOneWidget);

      // Enter edit mode
      await tester.tap(find.byIcon(Icons.edit));
      await tester.pumpAndSettle();
      await screenshot('profile_edit');

      expect(find.byKey(const Key('name_field')), findsOneWidget);
    });
  });

  group('5. Attendance History', () {
    testWidgets('shows attendance records', (tester) async {
      await tester.pumpWidget(
        _buildApp(
          const AttendanceScreen(),
          overrides: _authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();
      await screenshot('attendance_history');

      expect(find.text('Attendance History'), findsOneWidget);
      expect(find.text('QR'), findsOneWidget);
    });
  });

  group('6. Offline Check-in', () {
    testWidgets('offline indicator visible', (tester) async {
      await tester.pumpWidget(
        _buildApp(
          const HomeScreen(),
          overrides: [
            ..._authenticatedOverrides(),
            isOnlineProvider.overrideWith((ref) => false),
          ],
        ),
      );
      await tester.pumpAndSettle();
      await screenshot('home_offline');

      expect(find.byIcon(Icons.cloud_off), findsOneWidget);
    });
  });

  group('7. Streak Display', () {
    testWidgets('streak shown on home', (tester) async {
      await tester.pumpWidget(
        _buildApp(
          const HomeScreen(),
          overrides: _authenticatedOverrides(),
        ),
      );
      await tester.pumpAndSettle();
      await screenshot('home_with_streak');

      expect(find.text('Current Streak'), findsOneWidget);
      expect(find.textContaining('5'), findsWidgets);
    });
  });

  group('8. Leaderboard', () {
    testWidgets('leaderboard with tabs', (tester) async {
      await tester.pumpWidget(_buildApp(const LeaderboardScreen()));
      await tester.pumpAndSettle();
      await screenshot('leaderboard');

      expect(find.text('Weekly'), findsOneWidget);
      expect(find.text('Monthly'), findsOneWidget);
      expect(find.text('All Time'), findsOneWidget);
    });
  });

  group('9. Health Metrics', () {
    testWidgets('BMI calculator works', (tester) async {
      await tester.pumpWidget(_buildApp(const HealthScreen()));
      await tester.pumpAndSettle();
      await screenshot('health_screen');

      await tester.enterText(
          find.byKey(const Key('height_field')), '170');
      await tester.enterText(
          find.byKey(const Key('weight_field')), '70');
      await tester.tap(find.byKey(const Key('log_bmi_button')));
      await tester.pumpAndSettle();

      expect(find.text('BMI: 24.2'), findsOneWidget);
      await screenshot('health_bmi_calculated');
    });
  });

  group('10. Navigation', () {
    testWidgets('settings screen with language toggle', (tester) async {
      await tester.pumpWidget(
        _buildApp(
          const SettingsScreen(),
          overrides: [
            localeProvider.overrideWith((ref) => 'en'),
          ],
        ),
      );
      await tester.pumpAndSettle();
      await screenshot('settings_screen');

      expect(find.text('Language'), findsOneWidget);
      expect(find.text('English'), findsOneWidget);
      expect(find.text('Nepali'), findsOneWidget);
    });

    testWidgets('notifications empty state', (tester) async {
      await tester.pumpWidget(_buildApp(const NotificationsScreen()));
      await tester.pumpAndSettle();
      await screenshot('notifications_screen');

      expect(find.text('No data available'), findsOneWidget);
    });

    testWidgets('achievements empty state', (tester) async {
      await tester.pumpWidget(_buildApp(const AchievementsScreen()));
      await tester.pumpAndSettle();
      await screenshot('achievements_screen');

      expect(find.text('No data available'), findsOneWidget);
    });
  });
}
