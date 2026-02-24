class OrgAttendance {
  final String id;
  final String userId;
  final String orgId;
  final String checkInAt;
  final String method;
  final String? memberName;

  OrgAttendance({
    required this.id,
    required this.userId,
    required this.orgId,
    required this.checkInAt,
    required this.method,
    this.memberName,
  });

  factory OrgAttendance.fromJson(Map<String, dynamic> json) => OrgAttendance(
        id: json['id'] as String,
        userId: json['user_id'] as String,
        orgId: json['org_id'] as String,
        checkInAt: json['check_in_at'] as String,
        method: json['method'] as String,
        memberName: json['member_name'] as String?,
      );
}

class WeeklyMemberSummary {
  final String userId;
  final String name;
  final int daysCount;
  final int currentStreak;

  WeeklyMemberSummary({
    required this.userId,
    required this.name,
    required this.daysCount,
    required this.currentStreak,
  });

  factory WeeklyMemberSummary.fromJson(Map<String, dynamic> json) =>
      WeeklyMemberSummary(
        userId: json['user_id'] as String,
        name: json['name'] as String? ?? 'Unknown',
        daysCount: (json['days_count'] as num?)?.toInt() ?? 0,
        currentStreak: (json['current_streak'] as num?)?.toInt() ?? 0,
      );
}

class MemberAttendanceRecord {
  final String id;
  final String userId;
  final String orgId;
  final String checkInAt;
  final String method;

  MemberAttendanceRecord({
    required this.id,
    required this.userId,
    required this.orgId,
    required this.checkInAt,
    required this.method,
  });

  factory MemberAttendanceRecord.fromJson(Map<String, dynamic> json) =>
      MemberAttendanceRecord(
        id: json['id'] as String,
        userId: json['user_id'] as String,
        orgId: json['org_id'] as String,
        checkInAt: json['check_in_at'] as String,
        method: json['method'] as String,
      );
}

class CheckInResponse {
  final OrgAttendance checkIn;
  final StreakInfo streak;

  CheckInResponse({required this.checkIn, required this.streak});

  factory CheckInResponse.fromJson(Map<String, dynamic> json) =>
      CheckInResponse(
        checkIn:
            OrgAttendance.fromJson(json['check_in'] as Map<String, dynamic>),
        streak: StreakInfo.fromJson(json['streak'] as Map<String, dynamic>),
      );
}

class StreakInfo {
  final String id;
  final String memberId;
  final String orgId;
  final int currentStreak;
  final int longestStreak;
  final String? lastCheckIn;
  final String updatedAt;

  StreakInfo({
    required this.id,
    required this.memberId,
    required this.orgId,
    required this.currentStreak,
    required this.longestStreak,
    this.lastCheckIn,
    required this.updatedAt,
  });

  factory StreakInfo.fromJson(Map<String, dynamic> json) => StreakInfo(
        id: json['id'] as String,
        memberId: json['member_id'] as String,
        orgId: json['org_id'] as String,
        currentStreak: (json['current_streak'] as num?)?.toInt() ?? 0,
        longestStreak: (json['longest_streak'] as num?)?.toInt() ?? 0,
        lastCheckIn: json['last_check_in'] as String?,
        updatedAt: json['updated_at'] as String,
      );
}
