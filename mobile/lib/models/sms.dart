class SmsBalance {
  final String organizationId;
  final int balance;
  final int totalPurchased;
  final int totalUsed;
  final int totalCampaigns;

  SmsBalance({
    required this.organizationId,
    required this.balance,
    required this.totalPurchased,
    required this.totalUsed,
    required this.totalCampaigns,
  });

  factory SmsBalance.fromJson(Map<String, dynamic> json) => SmsBalance(
        organizationId: json['organization_id'] as String,
        balance: json['balance'] as int? ?? 0,
        totalPurchased: json['total_purchased'] as int? ?? 0,
        totalUsed: json['total_used'] as int? ?? 0,
        totalCampaigns: json['total_campaigns'] as int? ?? 0,
      );
}
