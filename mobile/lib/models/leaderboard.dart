class LeaderboardEntry {
  final int rank;
  final String memberId;
  final String name;
  final String? avatarUrl;
  final int value;
  final String metric;

  LeaderboardEntry({
    required this.rank,
    required this.memberId,
    required this.name,
    this.avatarUrl,
    required this.value,
    required this.metric,
  });

  factory LeaderboardEntry.fromJson(Map<String, dynamic> json) {
    return LeaderboardEntry(
      rank: (json['rank'] as num).toInt(),
      memberId: json['member_id'] as String,
      name: json['name'] as String,
      avatarUrl: json['avatar_url'] as String?,
      value: (json['value'] as num).toInt(),
      metric: json['metric'] as String,
    );
  }
}

class LeaderboardResponse {
  final String period;
  final List<LeaderboardEntry> rankings;

  LeaderboardResponse({
    required this.period,
    required this.rankings,
  });

  factory LeaderboardResponse.fromJson(Map<String, dynamic> json) {
    return LeaderboardResponse(
      period: json['period'] as String,
      rankings: (json['rankings'] as List<dynamic>? ?? [])
          .map((e) => LeaderboardEntry.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }
}
