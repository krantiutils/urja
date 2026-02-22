import 'dart:math' as math;

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/l10n/app_localizations.dart';

import '../../config/theme.dart';
import '../../models/nutrition.dart';
import '../../models/weight.dart';
import '../../providers/auth_provider.dart';
import '../../services/member_service.dart';
import '../../services/nutrition_service.dart';
import '../shared/widgets/loading_indicator.dart';

class MemberProgressScreen extends ConsumerStatefulWidget {
  const MemberProgressScreen({super.key});

  @override
  ConsumerState<MemberProgressScreen> createState() =>
      _MemberProgressScreenState();
}

class _MemberProgressScreenState extends ConsumerState<MemberProgressScreen> {
  late final MemberService _memberService;
  late final NutritionService _nutritionService;

  final _weightController = TextEditingController();
  final _formKey = GlobalKey<FormState>();

  WeightTrend? _weightTrend;
  NutritionStreak? _nutritionStreak;
  bool _loading = true;
  bool _logging = false;
  String? _error;
  int _selectedDays = 30;

  @override
  void initState() {
    super.initState();
    _memberService = MemberService(ref.read(apiClientProvider));
    _nutritionService = NutritionService(ref.read(apiClientProvider));
    _loadData();
  }

  @override
  void dispose() {
    _weightController.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final results = await Future.wait([
        _memberService.getWeightTrend(days: _selectedDays),
        _nutritionService.getNutritionStreak(),
      ]);
      setState(() {
        _weightTrend = results[0] as WeightTrend;
        _nutritionStreak = results[1] as NutritionStreak;
        _loading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  Future<void> _logWeight() async {
    if (!_formKey.currentState!.validate()) return;

    final weight = double.tryParse(_weightController.text.trim());
    if (weight == null) return;

    setState(() => _logging = true);
    try {
      await _memberService.quickLogWeight(weight);
      _weightController.clear();
      if (mounted) {
        final l10n = AppLocalizations.of(context)!;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(l10n.progressWeightLogged),
            backgroundColor: AppTheme.primary,
          ),
        );
      }
      await _loadData();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(e.toString()),
            backgroundColor: AppTheme.error,
          ),
        );
      }
    } finally {
      if (mounted) setState(() => _logging = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      backgroundColor: AppTheme.background,
      appBar: AppBar(title: Text(l10n.progressTitle)),
      body: _loading
          ? const LoadingIndicator()
          : _error != null
              ? Center(
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      const Icon(Icons.error_outline,
                          size: 48, color: AppTheme.error),
                      const SizedBox(height: 16),
                      Text(l10n.errorLoadingData,
                          style: const TextStyle(
                              color: AppTheme.textSecondary)),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: _loadData,
                        child: Text(l10n.retry),
                      ),
                    ],
                  ),
                )
              : RefreshIndicator(
                  color: AppTheme.primary,
                  onRefresh: _loadData,
                  child: _buildContent(l10n),
                ),
    );
  }

  Widget _buildContent(AppLocalizations l10n) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        _buildWeightLogForm(l10n),
        const SizedBox(height: 20),
        if (_weightTrend != null &&
            _weightTrend!.summary.entriesCount > 0) ...[
          _buildSummaryStats(l10n),
          const SizedBox(height: 20),
          _buildPeriodSelector(l10n),
          const SizedBox(height: 12),
          _buildWeightChart(l10n),
        ] else
          _buildEmptyChart(l10n),
        const SizedBox(height: 20),
        _buildNutritionStreakCard(l10n),
        const SizedBox(height: 24),
      ],
    );
  }

  // -- Weight Log Form --
  Widget _buildWeightLogForm(AppLocalizations l10n) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: AppTheme.primary.withAlpha(30),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: const Icon(Icons.monitor_weight_outlined,
                        color: AppTheme.primary, size: 20),
                  ),
                  const SizedBox(width: 12),
                  Text(
                    l10n.progressLogWeight,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: AppTheme.textPrimary,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Expanded(
                    child: TextFormField(
                      controller: _weightController,
                      keyboardType:
                          const TextInputType.numberWithOptions(decimal: true),
                      style: const TextStyle(color: AppTheme.textPrimary),
                      decoration: InputDecoration(
                        hintText: l10n.progressWeightKg,
                        hintStyle:
                            const TextStyle(color: AppTheme.textSecondary),
                        suffixText: l10n.progressKg,
                        suffixStyle:
                            const TextStyle(color: AppTheme.textSecondary),
                      ),
                      validator: (value) {
                        if (value == null || value.trim().isEmpty) {
                          return l10n.progressWeightKg;
                        }
                        final w = double.tryParse(value.trim());
                        if (w == null || w <= 0 || w > 500) {
                          return l10n.progressWeightKg;
                        }
                        return null;
                      },
                    ),
                  ),
                  const SizedBox(width: 12),
                  SizedBox(
                    height: 48,
                    child: ElevatedButton(
                      onPressed: _logging ? null : _logWeight,
                      child: _logging
                          ? const SizedBox(
                              width: 20,
                              height: 20,
                              child: CircularProgressIndicator(
                                strokeWidth: 2,
                                color: Colors.white,
                              ),
                            )
                          : Text(l10n.progressLogWeight),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  // -- Summary Stats --
  Widget _buildSummaryStats(AppLocalizations l10n) {
    final summary = _weightTrend!.summary;
    final isLoss = summary.changeKg < 0;
    final changeColor = isLoss ? AppTheme.primary : AppTheme.error;
    final changePrefix = summary.changeKg > 0 ? '+' : '';

    return Row(
      children: [
        Expanded(
          child: _statCard(
            icon: Icons.monitor_weight,
            iconColor: AppTheme.primary,
            value: '${summary.currentKg.toStringAsFixed(1)}',
            label: l10n.progressCurrentWeight,
          ),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: _statCard(
            icon: Icons.flag_outlined,
            iconColor: AppTheme.info,
            value: '${summary.startKg.toStringAsFixed(1)}',
            label: l10n.progressStartWeight,
          ),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: _statCard(
            icon: Icons.trending_down,
            iconColor: changeColor,
            value: '$changePrefix${summary.changeKg.toStringAsFixed(1)}',
            label: l10n.progressChange,
            valueColor: changeColor,
          ),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: _statCard(
            icon: Icons.accessibility_new,
            iconColor: AppTheme.warning,
            value: summary.bmi.toStringAsFixed(1),
            label: l10n.progressBmi,
          ),
        ),
      ],
    );
  }

  Widget _statCard({
    required IconData icon,
    required Color iconColor,
    required String value,
    required String label,
    Color? valueColor,
  }) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 8),
        child: Column(
          children: [
            Icon(icon, size: 20, color: iconColor),
            const SizedBox(height: 8),
            Text(
              value,
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
                color: valueColor ?? AppTheme.textPrimary,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              label,
              style: const TextStyle(
                fontSize: 11,
                color: AppTheme.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  // -- Period Selector --
  Widget _buildPeriodSelector(AppLocalizations l10n) {
    final periods = [
      (30, l10n.progressPeriod30),
      (90, l10n.progressPeriod90),
      (180, l10n.progressPeriod180),
    ];

    return Row(
      children: periods.map((p) {
        final (days, label) = p;
        final selected = _selectedDays == days;
        return Padding(
          padding: const EdgeInsets.only(right: 8),
          child: ChoiceChip(
            label: Text(label),
            selected: selected,
            selectedColor: AppTheme.primary.withAlpha(50),
            labelStyle: TextStyle(
              color: selected ? AppTheme.primary : AppTheme.textSecondary,
              fontWeight: selected ? FontWeight.w600 : FontWeight.normal,
              fontSize: 13,
            ),
            side: selected
                ? const BorderSide(color: AppTheme.primary, width: 1)
                : BorderSide.none,
            onSelected: (val) {
              if (val && _selectedDays != days) {
                setState(() => _selectedDays = days);
                _loadData();
              }
            },
          ),
        );
      }).toList(),
    );
  }

  // -- Weight Chart --
  Widget _buildWeightChart(AppLocalizations l10n) {
    final entries = _weightTrend!.entries;
    if (entries.isEmpty) return _buildEmptyChart(l10n);

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.show_chart,
                    color: AppTheme.primary, size: 20),
                const SizedBox(width: 8),
                Text(
                  l10n.progressWeightTrend,
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: AppTheme.textPrimary,
                  ),
                ),
                const Spacer(),
                Text(
                  '${entries.length} ${l10n.progressEntries}',
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppTheme.textSecondary,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 20),
            SizedBox(
              height: 200,
              child: CustomPaint(
                size: Size.infinite,
                painter: _WeightChartPainter(entries: entries),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyChart(AppLocalizations l10n) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          children: [
            Icon(Icons.show_chart,
                size: 48, color: AppTheme.textSecondary.withAlpha(100)),
            const SizedBox(height: 12),
            Text(
              l10n.progressNoData,
              textAlign: TextAlign.center,
              style: const TextStyle(
                color: AppTheme.textSecondary,
                fontSize: 14,
              ),
            ),
          ],
        ),
      ),
    );
  }

  // -- Nutrition Streak --
  Widget _buildNutritionStreakCard(AppLocalizations l10n) {
    final streak = _nutritionStreak;
    if (streak == null) return const SizedBox.shrink();

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: AppTheme.warning.withAlpha(30),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: const Icon(Icons.local_fire_department,
                      color: AppTheme.warning, size: 20),
                ),
                const SizedBox(width: 12),
                Text(
                  l10n.progressNutritionStreak,
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: AppTheme.textPrimary,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: AppTheme.surfaceLight,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Column(
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            const Icon(Icons.local_fire_department,
                                color: AppTheme.warning, size: 24),
                            const SizedBox(width: 4),
                            Text(
                              '${streak.currentStreak}',
                              style: const TextStyle(
                                fontSize: 28,
                                fontWeight: FontWeight.bold,
                                color: AppTheme.warning,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 4),
                        Text(
                          l10n.currentStreak,
                          style: const TextStyle(
                            fontSize: 12,
                            color: AppTheme.textSecondary,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: AppTheme.surfaceLight,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Column(
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            const Icon(Icons.emoji_events,
                                color: AppTheme.primary, size: 24),
                            const SizedBox(width: 4),
                            Text(
                              '${streak.longestStreak}',
                              style: const TextStyle(
                                fontSize: 28,
                                fontWeight: FontWeight.bold,
                                color: AppTheme.primary,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 4),
                        Text(
                          l10n.longestStreak,
                          style: const TextStyle(
                            fontSize: 12,
                            color: AppTheme.textSecondary,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

// -- Custom Weight Trend Chart Painter --
class _WeightChartPainter extends CustomPainter {
  final List<WeightEntry> entries;

  _WeightChartPainter({required this.entries});

  @override
  void paint(Canvas canvas, Size size) {
    if (entries.isEmpty) return;

    const leftPadding = 44.0;
    const bottomPadding = 28.0;
    const topPadding = 8.0;
    const rightPadding = 12.0;

    final chartWidth = size.width - leftPadding - rightPadding;
    final chartHeight = size.height - topPadding - bottomPadding;

    final weights = entries.map((e) => e.weightKg).toList();
    final minWeight = weights.reduce(math.min) - 1;
    final maxWeight = weights.reduce(math.max) + 1;
    final weightRange = maxWeight - minWeight;

    // Calculate points
    final points = <Offset>[];
    for (var i = 0; i < entries.length; i++) {
      final x = entries.length == 1
          ? leftPadding + chartWidth / 2
          : leftPadding + (i / (entries.length - 1)) * chartWidth;
      final y = topPadding +
          chartHeight -
          ((entries[i].weightKg - minWeight) / weightRange) * chartHeight;
      points.add(Offset(x, y));
    }

    // Grid lines
    final gridPaint = Paint()
      ..color = const Color(0xFF2D2D44)
      ..strokeWidth = 0.5;

    final textStyle = TextStyle(
      color: const Color(0xFF9CA3AF),
      fontSize: 10,
    );

    const gridLines = 4;
    for (var i = 0; i <= gridLines; i++) {
      final y = topPadding + (i / gridLines) * chartHeight;
      canvas.drawLine(
        Offset(leftPadding, y),
        Offset(size.width - rightPadding, y),
        gridPaint,
      );

      // Y-axis labels
      final weightVal = maxWeight - (i / gridLines) * weightRange;
      final tp = TextPainter(
        text: TextSpan(text: weightVal.toStringAsFixed(1), style: textStyle),
        textDirection: TextDirection.ltr,
      );
      tp.layout();
      tp.paint(canvas, Offset(leftPadding - tp.width - 6, y - tp.height / 2));
    }

    // X-axis date labels (show a few evenly spaced)
    final labelCount = math.min(entries.length, 5);
    for (var i = 0; i < labelCount; i++) {
      final idx = entries.length == 1
          ? 0
          : (i * (entries.length - 1) / (labelCount - 1)).round();
      final entry = entries[idx];
      final dateStr = _shortDate(entry.date);
      final tp = TextPainter(
        text: TextSpan(text: dateStr, style: textStyle),
        textDirection: TextDirection.ltr,
      );
      tp.layout();
      final x = entries.length == 1
          ? leftPadding + chartWidth / 2
          : leftPadding + (idx / (entries.length - 1)) * chartWidth;
      tp.paint(canvas,
          Offset(x - tp.width / 2, size.height - bottomPadding + 8));
    }

    if (points.length < 2) {
      // Single dot
      final dotPaint = Paint()..color = const Color(0xFF22C55E);
      canvas.drawCircle(points.first, 6, dotPaint);
      final innerPaint = Paint()..color = Colors.white;
      canvas.drawCircle(points.first, 3, innerPaint);
      return;
    }

    // Gradient fill below line
    final fillPath = Path();
    fillPath.moveTo(points.first.dx, topPadding + chartHeight);
    for (final p in points) {
      fillPath.lineTo(p.dx, p.dy);
    }
    fillPath.lineTo(points.last.dx, topPadding + chartHeight);
    fillPath.close();

    final fillPaint = Paint()
      ..shader = LinearGradient(
        begin: Alignment.topCenter,
        end: Alignment.bottomCenter,
        colors: [
          const Color(0xFF22C55E).withAlpha(60),
          const Color(0xFF22C55E).withAlpha(5),
        ],
      ).createShader(
          Rect.fromLTWH(leftPadding, topPadding, chartWidth, chartHeight));
    canvas.drawPath(fillPath, fillPaint);

    // Line
    final linePaint = Paint()
      ..color = const Color(0xFF22C55E)
      ..strokeWidth = 2.5
      ..style = PaintingStyle.stroke
      ..strokeCap = StrokeCap.round
      ..strokeJoin = StrokeJoin.round;

    final linePath = Path();
    linePath.moveTo(points.first.dx, points.first.dy);
    for (var i = 1; i < points.length; i++) {
      linePath.lineTo(points[i].dx, points[i].dy);
    }
    canvas.drawPath(linePath, linePaint);

    // Dots
    final dotPaint = Paint()..color = const Color(0xFF22C55E);
    final dotBorderPaint = Paint()
      ..color = Colors.white
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2;

    for (final p in points) {
      canvas.drawCircle(p, 4, dotPaint);
      canvas.drawCircle(p, 4, dotBorderPaint);
    }
  }

  String _shortDate(String date) {
    // Expects "2026-02-22" or similar
    final parts = date.split('-');
    if (parts.length >= 3) {
      return '${parts[1]}/${parts[2]}';
    }
    return date;
  }

  @override
  bool shouldRepaint(covariant _WeightChartPainter oldDelegate) {
    return oldDelegate.entries != entries;
  }
}
