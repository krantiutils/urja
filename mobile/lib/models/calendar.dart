class AttendanceCalendar {
  final String month;
  final int daysInMonth;
  final List<int> checkInDays;
  final int totalCheckIns;

  AttendanceCalendar({
    required this.month,
    required this.daysInMonth,
    required this.checkInDays,
    required this.totalCheckIns,
  });

  factory AttendanceCalendar.fromJson(Map<String, dynamic> json) {
    return AttendanceCalendar(
      month: json['month'] as String,
      daysInMonth: json['days_in_month'] as int,
      checkInDays: (json['check_in_days'] as List<dynamic>).map((e) => e as int).toList(),
      totalCheckIns: json['total_check_ins'] as int,
    );
  }
}
