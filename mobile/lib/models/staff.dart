class StaffMember {
  final String id;
  final String userId;
  final String name;
  final String phone;
  final String? email;
  final String? avatarUrl;
  final String role;
  final String staffRole;
  final String status;
  final String joinedAt;

  StaffMember({
    required this.id,
    required this.userId,
    required this.name,
    required this.phone,
    this.email,
    this.avatarUrl,
    required this.role,
    required this.staffRole,
    required this.status,
    required this.joinedAt,
  });

  factory StaffMember.fromJson(Map<String, dynamic> json) => StaffMember(
        id: json['id'] as String,
        userId: json['user_id'] as String,
        name: json['name'] as String,
        phone: json['phone'] as String,
        email: json['email'] as String?,
        avatarUrl: json['avatar_url'] as String?,
        role: json['role'] as String,
        staffRole: json['staff_role'] as String,
        status: json['status'] as String,
        joinedAt: json['joined_at'] as String,
      );

  Map<String, dynamic> toJson() => {
        'name': name,
        'email': email,
        'avatar_url': avatarUrl,
        'staff_role': staffRole,
      };
}
