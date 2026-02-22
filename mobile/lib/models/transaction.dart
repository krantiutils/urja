class Transaction {
  final String id;
  final String organizationId;
  final String category;
  final String description;
  final String transactionDate;
  final String transactionType;
  final String amount;
  final String paymentType;
  final String reference;
  final String entryBy;
  final String entryByName;
  final String createdAt;
  final String updatedAt;

  Transaction({
    required this.id,
    required this.organizationId,
    required this.category,
    required this.description,
    required this.transactionDate,
    required this.transactionType,
    required this.amount,
    required this.paymentType,
    required this.reference,
    required this.entryBy,
    required this.entryByName,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Transaction.fromJson(Map<String, dynamic> json) => Transaction(
        id: json['id'] as String,
        organizationId: json['organization_id'] as String,
        category: json['category'] as String? ?? '',
        description: json['description'] as String? ?? '',
        transactionDate: json['transaction_date'] as String,
        transactionType: json['transaction_type'] as String,
        amount: json['amount']?.toString() ?? '0',
        paymentType: json['payment_type'] as String? ?? '',
        reference: json['reference'] as String? ?? '',
        entryBy: json['entry_by'] as String? ?? '',
        entryByName: json['entry_by_name'] as String? ?? '',
        createdAt: json['created_at'] as String,
        updatedAt: json['updated_at'] as String,
      );

  Map<String, dynamic> toJson() => {
        'category': category,
        'description': description,
        'transaction_date': transactionDate,
        'transaction_type': transactionType,
        'amount': amount,
        'payment_type': paymentType,
        'reference': reference,
      };
}

class AccountsSummary {
  final String totalIncome;
  final String totalExpenses;
  final String grossProfit;
  final double profitPercent;

  AccountsSummary({
    required this.totalIncome,
    required this.totalExpenses,
    required this.grossProfit,
    required this.profitPercent,
  });

  factory AccountsSummary.fromJson(Map<String, dynamic> json) =>
      AccountsSummary(
        totalIncome: json['total_income']?.toString() ?? '0',
        totalExpenses: json['total_expenses']?.toString() ?? '0',
        grossProfit: json['gross_profit']?.toString() ?? '0',
        profitPercent: (json['profit_percent'] as num?)?.toDouble() ?? 0,
      );
}
