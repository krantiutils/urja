class WaterLog {
  final String id;
  final String userId;
  final int amountMl;
  final String loggedDate;
  final String loggedAt;

  WaterLog({
    required this.id,
    required this.userId,
    required this.amountMl,
    required this.loggedDate,
    required this.loggedAt,
  });

  factory WaterLog.fromJson(Map<String, dynamic> json) => WaterLog(
        id: json['id'] as String,
        userId: json['user_id'] as String,
        amountMl: json['amount_ml'] as int,
        loggedDate: json['logged_date'] as String,
        loggedAt: json['logged_at'] as String,
      );
}

class WaterDailySummary {
  final String date;
  final int totalMl;
  final int goalMl;
  final List<WaterLog> entries;

  WaterDailySummary({
    required this.date,
    required this.totalMl,
    required this.goalMl,
    required this.entries,
  });

  factory WaterDailySummary.fromJson(Map<String, dynamic> json) =>
      WaterDailySummary(
        date: json['date'] as String,
        totalMl: json['total_ml'] as int,
        goalMl: json['goal_ml'] as int,
        entries: (json['entries'] as List? ?? [])
            .map((e) => WaterLog.fromJson(e as Map<String, dynamic>))
            .toList(),
      );

  double get progressPercent =>
      goalMl > 0 ? (totalMl / goalMl).clamp(0.0, 1.0) : 0.0;
  int get remainingMl => (goalMl - totalMl).clamp(0, goalMl);
}
