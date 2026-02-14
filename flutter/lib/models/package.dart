class GymPackage {
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
  final List<dynamic> features;
  final bool isActive;
  final DateTime createdAt;
  final DateTime updatedAt;

  const GymPackage({
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
    this.features = const [],
    required this.isActive,
    required this.createdAt,
    required this.updatedAt,
  });

  factory GymPackage.fromJson(Map<String, dynamic> json) {
    return GymPackage(
      id: json['id'] as String,
      organizationId: json['organization_id'] as String,
      name: json['name'] as String,
      nameNe: json['name_ne'] as String?,
      description: json['description'] as String?,
      descriptionNe: json['description_ne'] as String?,
      durationDays: json['duration_days'] as int,
      price: json['price'].toString(),
      currency: json['currency'] as String? ?? 'NPR',
      maxMembers: json['max_members'] as int?,
      features: json['features'] as List<dynamic>? ?? [],
      isActive: json['is_active'] as bool,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  double get priceAsDouble => double.tryParse(price) ?? 0;
}

class MemberPackage {
  final String id;
  final String userId;
  final String packageId;
  final String organizationId;
  final DateTime startDate;
  final DateTime endDate;
  final String? paymentMethod;
  final String? paymentReference;
  final String amountPaid;
  final String status;
  final DateTime createdAt;
  final String? packageName;
  final String? orgName;

  const MemberPackage({
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

  factory MemberPackage.fromJson(Map<String, dynamic> json) {
    return MemberPackage(
      id: json['id'] as String,
      userId: json['user_id'] as String,
      packageId: json['package_id'] as String,
      organizationId: json['organization_id'] as String,
      startDate: DateTime.parse(json['start_date'] as String),
      endDate: DateTime.parse(json['end_date'] as String),
      paymentMethod: json['payment_method'] as String?,
      paymentReference: json['payment_reference'] as String?,
      amountPaid: json['amount_paid'].toString(),
      status: json['status'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
      packageName: json['package_name'] as String?,
      orgName: json['org_name'] as String?,
    );
  }

  bool get isActive => status == 'active';
  bool get isExpired => endDate.isBefore(DateTime.now());

  int get daysRemaining {
    final now = DateTime.now();
    if (endDate.isBefore(now)) return 0;
    return endDate.difference(now).inDays;
  }

  double get progressFraction {
    final total = endDate.difference(startDate).inDays;
    if (total <= 0) return 0;
    final elapsed = DateTime.now().difference(startDate).inDays;
    return (elapsed / total).clamp(0.0, 1.0);
  }
}

class SubscribeResult {
  final String memberPackageId;
  final String paymentUrl;
  final String pidx;

  const SubscribeResult({
    required this.memberPackageId,
    required this.paymentUrl,
    required this.pidx,
  });

  factory SubscribeResult.fromJson(Map<String, dynamic> json) {
    return SubscribeResult(
      memberPackageId: json['member_package_id'] as String,
      paymentUrl: json['payment_url'] as String,
      pidx: json['pidx'] as String,
    );
  }
}
