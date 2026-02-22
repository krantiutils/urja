class MemberStreak {
  final String id;
  final String memberId;
  final String orgId;
  final int currentStreak;
  final int longestStreak;
  final String? lastCheckIn;
  final String updatedAt;

  MemberStreak({
    required this.id,
    required this.memberId,
    required this.orgId,
    required this.currentStreak,
    required this.longestStreak,
    this.lastCheckIn,
    required this.updatedAt,
  });

  factory MemberStreak.fromJson(Map<String, dynamic> json) => MemberStreak(
        id: json['id'] as String,
        memberId: json['member_id'] as String,
        orgId: json['org_id'] as String,
        currentStreak: json['current_streak'] as int? ?? 0,
        longestStreak: json['longest_streak'] as int? ?? 0,
        lastCheckIn: json['last_check_in'] as String?,
        updatedAt: json['updated_at'] as String,
      );
}
