class PrivacySettings {
  final bool showEmail;
  final bool showPhone;
  final bool showProfile;
  final bool showAttendance;
  final bool showOnLeaderboard;

  const PrivacySettings({
    this.showEmail = false,
    this.showPhone = false,
    this.showProfile = true,
    this.showAttendance = true,
    this.showOnLeaderboard = true,
  });

  factory PrivacySettings.fromJson(Map<String, dynamic> json) {
    return PrivacySettings(
      showEmail: json['show_email'] as bool? ?? false,
      showPhone: json['show_phone'] as bool? ?? false,
      showProfile: json['show_profile'] as bool? ?? true,
      showAttendance: json['show_attendance'] as bool? ?? true,
      showOnLeaderboard: json['show_on_leaderboard'] as bool? ?? true,
    );
  }

  Map<String, dynamic> toJson() => {
        'show_email': showEmail,
        'show_phone': showPhone,
        'show_profile': showProfile,
        'show_attendance': showAttendance,
        'show_on_leaderboard': showOnLeaderboard,
      };
}

class OrgMembership {
  final String orgId;
  final String orgName;
  final String role;
  final String status;
  final DateTime joinedAt;

  const OrgMembership({
    required this.orgId,
    required this.orgName,
    required this.role,
    required this.status,
    required this.joinedAt,
  });

  factory OrgMembership.fromJson(Map<String, dynamic> json) {
    return OrgMembership(
      orgId: json['org_id'] as String,
      orgName: json['org_name'] as String,
      role: json['role'] as String,
      status: json['status'] as String,
      joinedAt: DateTime.parse(json['joined_at'] as String),
    );
  }
}

class Member {
  final String id;
  final String phone;
  final String name;
  final String? nameNe;
  final String? email;
  final String? dateOfBirth;
  final String? gender;
  final String? avatarUrl;
  final String? emergencyContactName;
  final String? emergencyContactPhone;
  final PrivacySettings privacySettings;
  final List<OrgMembership> organizations;
  final DateTime createdAt;
  final DateTime updatedAt;

  const Member({
    required this.id,
    required this.phone,
    required this.name,
    this.nameNe,
    this.email,
    this.dateOfBirth,
    this.gender,
    this.avatarUrl,
    this.emergencyContactName,
    this.emergencyContactPhone,
    this.privacySettings = const PrivacySettings(),
    this.organizations = const [],
    required this.createdAt,
    required this.updatedAt,
  });

  factory Member.fromJson(Map<String, dynamic> json) {
    return Member(
      id: json['id'] as String,
      phone: json['phone'] as String,
      name: json['name'] as String? ?? '',
      nameNe: json['name_ne'] as String?,
      email: _nonEmpty(json['email'] as String?),
      dateOfBirth: json['date_of_birth'] as String?,
      gender: _nonEmpty(json['gender'] as String?),
      avatarUrl: _nonEmpty(json['avatar_url'] as String?),
      emergencyContactName:
          _nonEmpty(json['emergency_contact_name'] as String?),
      emergencyContactPhone:
          _nonEmpty(json['emergency_contact_phone'] as String?),
      privacySettings: json['privacy_settings'] != null
          ? PrivacySettings.fromJson(
              json['privacy_settings'] as Map<String, dynamic>)
          : const PrivacySettings(),
      organizations: (json['organizations'] as List<dynamic>?)
              ?.map(
                  (e) => OrgMembership.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  String get displayName => name.isNotEmpty ? name : phone;

  /// Returns the first active org membership, if any.
  OrgMembership? get primaryOrg {
    if (organizations.isEmpty) return null;
    return organizations.first;
  }
}

class ProfileUpdate {
  final String? name;
  final String? nameNe;
  final String? email;
  final String? dateOfBirth;
  final String? gender;
  final String? avatarUrl;
  final String? emergencyContactName;
  final String? emergencyContactPhone;

  const ProfileUpdate({
    this.name,
    this.nameNe,
    this.email,
    this.dateOfBirth,
    this.gender,
    this.avatarUrl,
    this.emergencyContactName,
    this.emergencyContactPhone,
  });

  Map<String, dynamic> toJson() {
    final map = <String, dynamic>{};
    if (name != null) map['name'] = name;
    if (nameNe != null) map['name_ne'] = nameNe;
    if (email != null) map['email'] = email;
    if (dateOfBirth != null) map['date_of_birth'] = dateOfBirth;
    if (gender != null) map['gender'] = gender;
    if (avatarUrl != null) map['avatar_url'] = avatarUrl;
    if (emergencyContactName != null) {
      map['emergency_contact_name'] = emergencyContactName;
    }
    if (emergencyContactPhone != null) {
      map['emergency_contact_phone'] = emergencyContactPhone;
    }
    return map;
  }
}

String? _nonEmpty(String? s) => (s != null && s.isNotEmpty) ? s : null;
