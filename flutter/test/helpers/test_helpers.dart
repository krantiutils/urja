import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:urja_flutter/config/theme.dart';
import 'package:urja_flutter/l10n/app_localizations.dart';
import 'package:urja_flutter/models/attendance.dart';
import 'package:urja_flutter/models/package.dart';
import 'package:urja_flutter/models/user.dart';
import 'package:urja_flutter/providers/providers.dart';

// ---------------------------------------------------------------------------
// Mock data fixtures — all dates are deterministic, no DateTime.now() usage
// ---------------------------------------------------------------------------

final mockMember = Member(
  id: 'user-123',
  phone: '9841234567',
  name: 'Ram Sharma',
  nameNe: '\u0930\u093e\u092e \u0936\u0930\u094d\u092e\u093e',
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

final mockAttendanceRecords = [
  AttendanceRecord(
    id: 'att-1',
    userId: 'user-123',
    orgId: 'org-1',
    checkInAt: DateTime(2024, 6, 15, 8, 30),
    method: 'qr',
  ),
  AttendanceRecord(
    id: 'att-2',
    userId: 'user-123',
    orgId: 'org-1',
    checkInAt: DateTime(2024, 6, 14, 9, 0),
    method: 'nfc',
  ),
  AttendanceRecord(
    id: 'att-3',
    userId: 'user-123',
    orgId: 'org-1',
    checkInAt: DateTime(2024, 6, 13, 7, 45),
    method: 'manual',
  ),
];

final mockStreaks = [
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

final mockGymPackage = GymPackage(
  id: 'pkg-1',
  organizationId: 'org-1',
  name: 'Monthly Plan',
  description: 'Full access for one month',
  durationDays: 30,
  price: '3000',
  currency: 'NPR',
  features: ['Gym access', 'Cardio', 'Weights'],
  isActive: true,
  createdAt: DateTime(2024, 1, 1),
  updatedAt: DateTime(2024, 1, 1),
);

// ---------------------------------------------------------------------------
// Test Notifier overrides
// ---------------------------------------------------------------------------

class TestProfileNotifier extends ProfileNotifier {
  final Member? _mockMember;
  TestProfileNotifier(this._mockMember);

  @override
  Future<Member?> build() async => _mockMember;

  @override
  Future<void> refresh() async {}

  @override
  Future<void> updateProfile(ProfileUpdate update) async {}

  @override
  Future<void> updatePrivacy(PrivacySettings settings) async {}
}

class TestAttendanceNotifier extends AttendanceNotifier {
  final List<AttendanceRecord> _mockRecords;
  TestAttendanceNotifier(this._mockRecords);

  @override
  Future<List<AttendanceRecord>> build() async => _mockRecords;

  @override
  Future<void> refresh() async {}
}

// ---------------------------------------------------------------------------
// Test app wrapper
// ---------------------------------------------------------------------------

/// Wraps [child] in a [ProviderScope] with [MaterialApp], theme, and l10n.
Widget buildTestApp(
  Widget child, {
  List<Override> overrides = const [],
}) {
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

/// Standard overrides for screens that display member data.
List<Override> authenticatedOverrides({
  Member? member,
  List<AttendanceRecord>? attendance,
  List<Streak>? streaks,
  List<MemberPackage>? packages,
  bool online = true,
}) {
  return [
    authStateProvider.overrideWith((ref) => AuthState.authenticated),
    profileProvider.overrideWith(
        () => TestProfileNotifier(member ?? mockMember)),
    attendanceProvider.overrideWith(
        () => TestAttendanceNotifier(attendance ?? mockAttendanceRecords)),
    streaksProvider
        .overrideWith((ref) async => streaks ?? mockStreaks),
    memberPackagesProvider
        .overrideWith((ref) async => packages ?? []),
    isOnlineProvider.overrideWith((ref) => online),
    localeProvider.overrideWith((ref) => 'en'),
  ];
}

/// Sets a deterministic surface size for golden file tests.
Future<void> setGoldenSurfaceSize(WidgetTester tester) async {
  await tester.binding.setSurfaceSize(const Size(400, 800));
  addTearDown(() => tester.binding.setSurfaceSize(null));
}
