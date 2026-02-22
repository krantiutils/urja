class Due {
  final String id;
  final String organizationId;
  final String userId;
  final String memberName;
  final String memberPhone;
  final String amount;
  final String dueDate;
  final String description;
  final String status;
  final String? paidAt;
  final String? paidAmount;
  final String? paymentMethod;
  final String? paymentReference;
  final String createdAt;

  Due({
    required this.id,
    required this.organizationId,
    required this.userId,
    required this.memberName,
    required this.memberPhone,
    required this.amount,
    required this.dueDate,
    required this.description,
    required this.status,
    this.paidAt,
    this.paidAmount,
    this.paymentMethod,
    this.paymentReference,
    required this.createdAt,
  });

  factory Due.fromJson(Map<String, dynamic> json) => Due(
        id: json['id'] as String,
        organizationId: json['organization_id'] as String,
        userId: json['user_id'] as String,
        memberName: json['member_name'] as String,
        memberPhone: json['member_phone'] as String,
        amount: json['amount'] as String,
        dueDate: json['due_date'] as String,
        description: json['description'] as String? ?? '',
        status: json['status'] as String,
        paidAt: json['paid_at'] as String?,
        paidAmount: json['paid_amount'] as String?,
        paymentMethod: json['payment_method'] as String?,
        paymentReference: json['payment_reference'] as String?,
        createdAt: json['created_at'] as String,
      );
}
