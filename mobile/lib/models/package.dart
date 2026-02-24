class OrgPackage {
  final String id;
  final String organizationId;
  final String name;
  final String? nameNe;
  final String? description;
  final String? descriptionNe;
  final int durationDays;
  final String price;
  final String currency;
  final int? maxMembers;
  final List<String>? features;
  final bool isActive;
  final String createdAt;
  final String updatedAt;

  OrgPackage({
    required this.id,
    required this.organizationId,
    required this.name,
    this.nameNe,
    this.description,
    this.descriptionNe,
    required this.durationDays,
    required this.price,
    this.currency = 'NPR',
    this.maxMembers,
    this.features,
    required this.isActive,
    required this.createdAt,
    required this.updatedAt,
  });

  factory OrgPackage.fromJson(Map<String, dynamic> json) => OrgPackage(
        id: json['id'] as String,
        organizationId: json['organization_id'] as String,
        name: json['name'] as String,
        nameNe: json['name_ne'] as String?,
        description: json['description'] as String?,
        descriptionNe: json['description_ne'] as String?,
        durationDays: (json['duration_days'] as num).toInt(),
        price: json['price'] as String,
        currency: json['currency'] as String? ?? 'NPR',
        maxMembers: json['max_members'] as int?,
        features: (json['features'] as List<dynamic>?)
            ?.map((e) => e as String)
            .toList(),
        isActive: json['is_active'] as bool,
        createdAt: json['created_at'] as String,
        updatedAt: json['updated_at'] as String,
      );

  Map<String, dynamic> toJson() => {
        'name': name,
        'name_ne': nameNe,
        'description': description,
        'description_ne': descriptionNe,
        'duration_days': durationDays,
        'price': price,
        'max_members': maxMembers,
        'features': features,
        'is_active': isActive,
      };
}

class MemberPackage {
  final String id;
  final String userId;
  final String packageId;
  final String organizationId;
  final String startDate;
  final String endDate;
  final String? paymentMethod;
  final String? paymentReference;
  final String amountPaid;
  final String status;
  final String createdAt;
  final String? packageName;
  final String? orgName;

  MemberPackage({
    required this.id,
    required this.userId,
    required this.packageId,
    required this.organizationId,
    required this.startDate,
    required this.endDate,
    this.paymentMethod,
    this.paymentReference,
    required this.amountPaid,
    required this.status,
    required this.createdAt,
    this.packageName,
    this.orgName,
  });

  factory MemberPackage.fromJson(Map<String, dynamic> json) => MemberPackage(
        id: json['id'] as String,
        userId: json['user_id'] as String,
        packageId: json['package_id'] as String,
        organizationId: json['organization_id'] as String,
        startDate: json['start_date'] as String,
        endDate: json['end_date'] as String,
        paymentMethod: json['payment_method'] as String?,
        paymentReference: json['payment_reference'] as String?,
        amountPaid: json['amount_paid'] as String,
        status: json['status'] as String,
        createdAt: json['created_at'] as String,
        packageName: json['package_name'] as String?,
        orgName: json['org_name'] as String?,
      );

  int get daysRemaining {
    final end = DateTime.tryParse(endDate);
    if (end == null) return 0;
    final today = DateTime.now();
    final endOfDay = DateTime(end.year, end.month, end.day, 23, 59, 59);
    final diff = endOfDay.difference(today).inDays;
    return diff < 0 ? 0 : diff;
  }

  bool get isExpired => daysRemaining <= 0;
}
