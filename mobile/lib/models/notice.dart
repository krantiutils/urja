class Notice {
  final String id;
  final String organizationId;
  final String title;
  final String? titleNe;
  final String content;
  final String? contentNe;
  final String createdBy;
  final bool isActive;
  final String createdAt;
  final String updatedAt;

  Notice({
    required this.id,
    required this.organizationId,
    required this.title,
    this.titleNe,
    required this.content,
    this.contentNe,
    required this.createdBy,
    required this.isActive,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Notice.fromJson(Map<String, dynamic> json) => Notice(
        id: json['id'] as String,
        organizationId: json['organization_id'] as String,
        title: json['title'] as String,
        titleNe: json['title_ne'] as String?,
        content: json['content'] as String,
        contentNe: json['content_ne'] as String?,
        createdBy: json['created_by'] as String,
        isActive: json['is_active'] as bool? ?? true,
        createdAt: json['created_at'] as String,
        updatedAt: json['updated_at'] as String,
      );

  Map<String, dynamic> toJson() => {
        'title': title,
        'title_ne': titleNe,
        'content': content,
        'content_ne': contentNe,
        'is_active': isActive,
      };
}
