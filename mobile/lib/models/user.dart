class AuthTokens {
  final String accessToken;
  final String refreshToken;
  final String tokenType;
  final bool isNewUser;
  final bool onboardingCompleted;

  AuthTokens({
    required this.accessToken,
    required this.refreshToken,
    required this.tokenType,
    this.isNewUser = false,
    this.onboardingCompleted = false,
  });

  factory AuthTokens.fromJson(Map<String, dynamic> json) => AuthTokens(
        accessToken: json['access_token'] as String,
        refreshToken: json['refresh_token'] as String,
        tokenType: json['token_type'] as String,
        isNewUser: json['is_new_user'] as bool? ?? false,
        onboardingCompleted: json['onboarding_completed'] as bool? ?? false,
      );
}

class User {
  final String id;
  final String phone;
  final String role;
  final String? orgId;
  final String? orgName;
  final String? userType;
  final bool onboardingCompleted;

  User({
    required this.id,
    required this.phone,
    required this.role,
    this.orgId,
    this.orgName,
    this.userType,
    this.onboardingCompleted = false,
  });

  User copyWith({
    String? id,
    String? phone,
    String? role,
    String? orgId,
    String? orgName,
    String? userType,
    bool? onboardingCompleted,
  }) =>
      User(
        id: id ?? this.id,
        phone: phone ?? this.phone,
        role: role ?? this.role,
        orgId: orgId ?? this.orgId,
        orgName: orgName ?? this.orgName,
        userType: userType ?? this.userType,
        onboardingCompleted: onboardingCompleted ?? this.onboardingCompleted,
      );
}

class AuthResult {
  final User user;
  final bool isNewUser;
  final bool onboardingCompleted;

  AuthResult({
    required this.user,
    required this.isNewUser,
    required this.onboardingCompleted,
  });
}

class PrivacySettings {
  final bool showEmail;
  final bool showPhone;
  final bool showProfile;
  final bool showAttendance;
  final bool showOnLeaderboard;

  PrivacySettings({
    required this.showEmail,
    required this.showPhone,
    required this.showProfile,
    required this.showAttendance,
    required this.showOnLeaderboard,
  });

  factory PrivacySettings.fromJson(Map<String, dynamic> json) =>
      PrivacySettings(
        showEmail: json['show_email'] as bool? ?? false,
        showPhone: json['show_phone'] as bool? ?? false,
        showProfile: json['show_profile'] as bool? ?? false,
        showAttendance: json['show_attendance'] as bool? ?? false,
        showOnLeaderboard: json['show_on_leaderboard'] as bool? ?? false,
      );

  Map<String, dynamic> toJson() => {
        'show_email': showEmail,
        'show_phone': showPhone,
        'show_profile': showProfile,
        'show_attendance': showAttendance,
        'show_on_leaderboard': showOnLeaderboard,
      };
}
