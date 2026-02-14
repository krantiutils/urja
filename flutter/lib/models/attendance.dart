class AttendanceRecord {
  final String id;
  final String userId;
  final String orgId;
  final DateTime checkInAt;
  final String method;

  const AttendanceRecord({
    required this.id,
    required this.userId,
    required this.orgId,
    required this.checkInAt,
    required this.method,
  });

  factory AttendanceRecord.fromJson(Map<String, dynamic> json) {
    return AttendanceRecord(
      id: json['id'] as String,
      userId: json['user_id'] as String,
      orgId: json['org_id'] as String,
      checkInAt: DateTime.parse(json['check_in_at'] as String),
      method: json['method'] as String,
    );
  }
}

class Streak {
  final String id;
  final String memberId;
  final String orgId;
  final int currentStreak;
  final int longestStreak;
  final String? lastCheckIn;
  final DateTime updatedAt;

  const Streak({
    required this.id,
    required this.memberId,
    required this.orgId,
    required this.currentStreak,
    required this.longestStreak,
    this.lastCheckIn,
    required this.updatedAt,
  });

  factory Streak.fromJson(Map<String, dynamic> json) {
    return Streak(
      id: json['id'] as String,
      memberId: json['member_id'] as String,
      orgId: json['org_id'] as String,
      currentStreak: json['current_streak'] as int,
      longestStreak: json['longest_streak'] as int,
      lastCheckIn: json['last_check_in'] as String?,
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }
}

class CheckInResult {
  final AttendanceRecord checkIn;
  final Streak? streak;

  const CheckInResult({required this.checkIn, this.streak});

  factory CheckInResult.fromJson(Map<String, dynamic> json) {
    return CheckInResult(
      checkIn: AttendanceRecord.fromJson(
          json['check_in'] as Map<String, dynamic>),
      streak: json['streak'] != null
          ? Streak.fromJson(json['streak'] as Map<String, dynamic>)
          : null,
    );
  }
}
