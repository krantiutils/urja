class MemberFeedback {
  final String id;
  final String organizationId;
  final String userId;
  final String memberName;
  final String? avatarUrl;
  final String message;
  final String createdAt;

  MemberFeedback({
    required this.id,
    required this.organizationId,
    required this.userId,
    required this.memberName,
    this.avatarUrl,
    required this.message,
    required this.createdAt,
  });

  factory MemberFeedback.fromJson(Map<String, dynamic> json) => MemberFeedback(
        id: json['id'] as String,
        organizationId: json['organization_id'] as String,
        userId: json['user_id'] as String,
        memberName: json['member_name'] as String,
        avatarUrl: json['avatar_url'] as String?,
        message: json['message'] as String,
        createdAt: json['created_at'] as String,
      );
}
