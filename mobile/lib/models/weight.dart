class WeightEntry {
  final String date;
  final double weightKg;

  WeightEntry({required this.date, required this.weightKg});

  factory WeightEntry.fromJson(Map<String, dynamic> json) => WeightEntry(
    date: json['date'] as String,
    weightKg: (json['weight_kg'] as num).toDouble(),
  );
}

class WeightTrendSummary {
  final double currentKg;
  final double startKg;
  final double changeKg;
  final double bmi;
  final int entriesCount;

  WeightTrendSummary({
    required this.currentKg,
    required this.startKg,
    required this.changeKg,
    required this.bmi,
    required this.entriesCount,
  });

  factory WeightTrendSummary.fromJson(Map<String, dynamic> json) =>
      WeightTrendSummary(
        currentKg: (json['current_kg'] as num?)?.toDouble() ?? 0,
        startKg: (json['start_kg'] as num?)?.toDouble() ?? 0,
        changeKg: (json['change_kg'] as num?)?.toDouble() ?? 0,
        bmi: (json['bmi'] as num?)?.toDouble() ?? 0,
        entriesCount: json['entries_count'] as int? ?? 0,
      );
}

class WeightTrend {
  final List<WeightEntry> entries;
  final WeightTrendSummary summary;

  WeightTrend({required this.entries, required this.summary});

  factory WeightTrend.fromJson(Map<String, dynamic> json) => WeightTrend(
        entries: (json['entries'] as List<dynamic>?)
                ?.map((e) =>
                    WeightEntry.fromJson(e as Map<String, dynamic>))
                .toList() ??
            [],
        summary: WeightTrendSummary.fromJson(
            json['summary'] as Map<String, dynamic>? ?? {}),
      );
}
