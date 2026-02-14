import 'package:flutter_test/flutter_test.dart';

import 'package:urja_flutter/models/user.dart';
import 'package:urja_flutter/models/attendance.dart';
import 'package:urja_flutter/models/package.dart';

void main() {
  group('Member model', () {
    final json = {
      'id': 'u-1',
      'phone': '9841234567',
      'name': 'Ram',
      'name_ne': 'राम',
      'email': 'ram@example.com',
      'date_of_birth': '1990-05-20',
      'gender': 'male',
      'avatar_url': null,
      'emergency_contact_name': 'Sita',
      'emergency_contact_phone': '9841234568',
      'privacy_settings': {
        'show_email': true,
        'show_phone': false,
        'show_profile': true,
        'show_attendance': true,
        'show_on_leaderboard': true,
      },
      'organizations': [
        {
          'org_id': 'org-1',
          'org_name': 'Urja',
          'role': 'member',
          'status': 'active',
          'joined_at': '2024-01-01T00:00:00Z',
        }
      ],
      'created_at': '2024-01-01T00:00:00Z',
      'updated_at': '2024-06-01T00:00:00Z',
    };

    test('fromJson deserializes all fields', () {
      final member = Member.fromJson(json);
      expect(member.id, 'u-1');
      expect(member.phone, '9841234567');
      expect(member.name, 'Ram');
      expect(member.nameNe, 'राम');
      expect(member.email, 'ram@example.com');
      expect(member.gender, 'male');
      expect(member.dateOfBirth, '1990-05-20');
      expect(member.emergencyContactName, 'Sita');
      expect(member.emergencyContactPhone, '9841234568');
      expect(member.organizations, hasLength(1));
      expect(member.organizations.first.orgName, 'Urja');
    });

    test('displayName returns name when available', () {
      final member = Member.fromJson(json);
      expect(member.displayName, 'Ram');
    });

    test('displayName returns phone when name is empty', () {
      final member = Member.fromJson({...json, 'name': ''});
      expect(member.displayName, '9841234567');
    });

    test('primaryOrg returns first org', () {
      final member = Member.fromJson(json);
      expect(member.primaryOrg, isNotNull);
      expect(member.primaryOrg!.orgId, 'org-1');
    });

    test('primaryOrg returns null when no orgs', () {
      final member = Member.fromJson({...json, 'organizations': <dynamic>[]});
      expect(member.primaryOrg, isNull);
    });

    test('handles null/empty email as null', () {
      final member = Member.fromJson({...json, 'email': ''});
      expect(member.email, isNull);
    });

    test('handles null/empty avatar_url as null', () {
      final member = Member.fromJson({...json, 'avatar_url': ''});
      expect(member.avatarUrl, isNull);
    });
  });

  group('PrivacySettings model', () {
    test('fromJson with all fields', () {
      final settings = PrivacySettings.fromJson({
        'show_email': true,
        'show_phone': true,
        'show_profile': false,
        'show_attendance': false,
        'show_on_leaderboard': false,
      });
      expect(settings.showEmail, true);
      expect(settings.showPhone, true);
      expect(settings.showProfile, false);
      expect(settings.showAttendance, false);
      expect(settings.showOnLeaderboard, false);
    });

    test('toJson round-trips correctly', () {
      const settings = PrivacySettings(
        showEmail: true,
        showPhone: false,
      );
      final json = settings.toJson();
      expect(json['show_email'], true);
      expect(json['show_phone'], false);
      expect(json['show_profile'], true);
    });

    test('defaults when json values are null', () {
      final settings = PrivacySettings.fromJson({});
      expect(settings.showEmail, false);
      expect(settings.showPhone, false);
      expect(settings.showProfile, true);
      expect(settings.showAttendance, true);
      expect(settings.showOnLeaderboard, true);
    });
  });

  group('ProfileUpdate model', () {
    test('toJson includes only non-null fields', () {
      const update = ProfileUpdate(name: 'New Name', email: 'test@test.com');
      final json = update.toJson();
      expect(json, containsPair('name', 'New Name'));
      expect(json, containsPair('email', 'test@test.com'));
      expect(json.containsKey('gender'), false);
      expect(json.containsKey('date_of_birth'), false);
    });

    test('toJson is empty when all fields null', () {
      const update = ProfileUpdate();
      expect(update.toJson(), isEmpty);
    });

    test('toJson includes emergency contact fields', () {
      const update = ProfileUpdate(
        emergencyContactName: 'Jane',
        emergencyContactPhone: '1234567890',
      );
      final json = update.toJson();
      expect(json['emergency_contact_name'], 'Jane');
      expect(json['emergency_contact_phone'], '1234567890');
    });
  });

  group('OrgMembership model', () {
    test('fromJson deserializes correctly', () {
      final org = OrgMembership.fromJson({
        'org_id': 'org-1',
        'org_name': 'Test Org',
        'role': 'admin',
        'status': 'active',
        'joined_at': '2024-03-15T12:00:00Z',
      });
      expect(org.orgId, 'org-1');
      expect(org.orgName, 'Test Org');
      expect(org.role, 'admin');
      expect(org.status, 'active');
      expect(org.joinedAt.year, 2024);
    });
  });

  group('AttendanceRecord model', () {
    test('fromJson deserializes correctly', () {
      final record = AttendanceRecord.fromJson({
        'id': 'att-1',
        'user_id': 'u-1',
        'org_id': 'org-1',
        'check_in_at': '2024-06-15T08:30:00Z',
        'method': 'qr',
      });
      expect(record.id, 'att-1');
      expect(record.userId, 'u-1');
      expect(record.orgId, 'org-1');
      expect(record.checkInAt.hour, 8);
      expect(record.method, 'qr');
    });
  });

  group('Streak model', () {
    test('fromJson deserializes correctly', () {
      final streak = Streak.fromJson({
        'id': 's-1',
        'member_id': 'u-1',
        'org_id': 'org-1',
        'current_streak': 7,
        'longest_streak': 30,
        'last_check_in': '2024-06-15',
        'updated_at': '2024-06-15T10:00:00Z',
      });
      expect(streak.id, 's-1');
      expect(streak.currentStreak, 7);
      expect(streak.longestStreak, 30);
      expect(streak.lastCheckIn, '2024-06-15');
    });

    test('lastCheckIn can be null', () {
      final streak = Streak.fromJson({
        'id': 's-1',
        'member_id': 'u-1',
        'org_id': 'org-1',
        'current_streak': 0,
        'longest_streak': 0,
        'last_check_in': null,
        'updated_at': '2024-06-15T10:00:00Z',
      });
      expect(streak.lastCheckIn, isNull);
    });
  });

  group('CheckInResult model', () {
    test('fromJson with streak', () {
      final result = CheckInResult.fromJson({
        'check_in': {
          'id': 'att-1',
          'user_id': 'u-1',
          'org_id': 'org-1',
          'check_in_at': '2024-06-15T08:30:00Z',
          'method': 'qr',
        },
        'streak': {
          'id': 's-1',
          'member_id': 'u-1',
          'org_id': 'org-1',
          'current_streak': 5,
          'longest_streak': 10,
          'last_check_in': '2024-06-15',
          'updated_at': '2024-06-15T10:00:00Z',
        },
      });
      expect(result.checkIn.id, 'att-1');
      expect(result.streak, isNotNull);
      expect(result.streak!.currentStreak, 5);
    });

    test('fromJson without streak', () {
      final result = CheckInResult.fromJson({
        'check_in': {
          'id': 'att-1',
          'user_id': 'u-1',
          'org_id': 'org-1',
          'check_in_at': '2024-06-15T08:30:00Z',
          'method': 'nfc',
        },
      });
      expect(result.checkIn.method, 'nfc');
      expect(result.streak, isNull);
    });
  });

  group('GymPackage model', () {
    test('fromJson deserializes correctly', () {
      final pkg = GymPackage.fromJson({
        'id': 'pkg-1',
        'organization_id': 'org-1',
        'name': 'Monthly',
        'name_ne': 'मासिक',
        'description': 'One month access',
        'description_ne': null,
        'duration_days': 30,
        'price': '3000.00',
        'currency': 'NPR',
        'max_members': 100,
        'features': ['Gym', 'Cardio'],
        'is_active': true,
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-01T00:00:00Z',
      });
      expect(pkg.id, 'pkg-1');
      expect(pkg.name, 'Monthly');
      expect(pkg.durationDays, 30);
      expect(pkg.price, '3000.00');
      expect(pkg.priceAsDouble, 3000.0);
      expect(pkg.features, hasLength(2));
      expect(pkg.isActive, true);
    });

    test('priceAsDouble handles non-numeric price', () {
      final pkg = GymPackage.fromJson({
        'id': 'p-1',
        'organization_id': 'o-1',
        'name': 'Test',
        'duration_days': 7,
        'price': 'free',
        'is_active': true,
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-01T00:00:00Z',
      });
      expect(pkg.priceAsDouble, 0);
    });
  });

  group('MemberPackage model', () {
    test('fromJson deserializes correctly', () {
      final pkg = MemberPackage.fromJson({
        'id': 'mp-1',
        'user_id': 'u-1',
        'package_id': 'pkg-1',
        'organization_id': 'org-1',
        'start_date': '2024-06-01T00:00:00Z',
        'end_date': '2024-07-01T00:00:00Z',
        'payment_method': 'khalti',
        'payment_reference': 'ref-123',
        'amount_paid': '3000',
        'status': 'active',
        'created_at': '2024-06-01T00:00:00Z',
        'package_name': 'Monthly',
        'org_name': 'Urja',
      });
      expect(pkg.id, 'mp-1');
      expect(pkg.isActive, true);
      expect(pkg.packageName, 'Monthly');
    });

    test('isExpired returns true for past end dates', () {
      final pkg = MemberPackage(
        id: 'mp-1',
        userId: 'u-1',
        packageId: 'pkg-1',
        organizationId: 'org-1',
        startDate: DateTime(2020, 1, 1),
        endDate: DateTime(2020, 2, 1),
        amountPaid: '1000',
        status: 'active',
        createdAt: DateTime(2020, 1, 1),
      );
      expect(pkg.isExpired, true);
      expect(pkg.daysRemaining, 0);
    });

    test('progressFraction is between 0 and 1', () {
      final now = DateTime.now();
      final pkg = MemberPackage(
        id: 'mp-1',
        userId: 'u-1',
        packageId: 'pkg-1',
        organizationId: 'org-1',
        startDate: now.subtract(const Duration(days: 15)),
        endDate: now.add(const Duration(days: 15)),
        amountPaid: '1000',
        status: 'active',
        createdAt: now.subtract(const Duration(days: 15)),
      );
      expect(pkg.progressFraction, greaterThanOrEqualTo(0.0));
      expect(pkg.progressFraction, lessThanOrEqualTo(1.0));
    });

    test('progressFraction handles zero-duration package', () {
      final now = DateTime.now();
      final pkg = MemberPackage(
        id: 'mp-1',
        userId: 'u-1',
        packageId: 'pkg-1',
        organizationId: 'org-1',
        startDate: now,
        endDate: now,
        amountPaid: '0',
        status: 'active',
        createdAt: now,
      );
      expect(pkg.progressFraction, 0);
    });
  });

  group('SubscribeResult model', () {
    test('fromJson deserializes correctly', () {
      final result = SubscribeResult.fromJson({
        'member_package_id': 'mp-1',
        'payment_url': 'https://pay.khalti.com/abc',
        'pidx': 'pidx-123',
      });
      expect(result.memberPackageId, 'mp-1');
      expect(result.paymentUrl, 'https://pay.khalti.com/abc');
      expect(result.pidx, 'pidx-123');
    });
  });
}
