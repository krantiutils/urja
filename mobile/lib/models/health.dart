class HealthMetric {
  final String id;
  final String memberId;
  final String metricType;
  final Map<String, dynamic> value;
  final String recordedAt;

  HealthMetric({
    required this.id,
    required this.memberId,
    required this.metricType,
    required this.value,
    required this.recordedAt,
  });

  factory HealthMetric.fromJson(Map<String, dynamic> json) => HealthMetric(
        id: json['id'] as String,
        memberId: json['member_id'] as String,
        metricType: json['metric_type'] as String,
        value: json['value'] as Map<String, dynamic>? ?? {},
        recordedAt: json['recorded_at'] as String,
      );
}
