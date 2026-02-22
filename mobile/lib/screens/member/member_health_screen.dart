import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../models/health.dart';
import '../../providers/auth_provider.dart';
import '../../services/member_service.dart';
import '../shared/widgets/empty_state.dart';
import '../shared/widgets/loading_indicator.dart';

class MemberHealthScreen extends ConsumerStatefulWidget {
  const MemberHealthScreen({super.key});

  @override
  ConsumerState<MemberHealthScreen> createState() =>
      _MemberHealthScreenState();
}

class _MemberHealthScreenState extends ConsumerState<MemberHealthScreen> {
  late final MemberService _memberService;

  List<HealthMetric> _bmiMetrics = [];
  List<HealthMetric> _bodyMetrics = [];
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _memberService = MemberService(ref.read(apiClientProvider));
    _loadData();
  }

  Future<void> _loadData() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final results = await Future.wait([
        _memberService.getHealth(type: 'bmi', limit: 10),
        _memberService.getHealth(type: 'body_measurement', limit: 10),
      ]);
      setState(() {
        _bmiMetrics = results[0];
        _bodyMetrics = results[1];
        _loading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  String _bmiCategory(double bmi) {
    if (bmi < 18.5) return 'Underweight';
    if (bmi < 25.0) return 'Normal';
    if (bmi < 30.0) return 'Overweight';
    return 'Obese';
  }

  Color _bmiColor(double bmi) {
    if (bmi < 18.5) return AppTheme.info;
    if (bmi < 25.0) return AppTheme.primary;
    if (bmi < 30.0) return AppTheme.warning;
    return AppTheme.error;
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      backgroundColor: AppTheme.background,
      appBar: AppBar(title: Text(l10n.myHealth)),
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
    final hasData = _bmiMetrics.isNotEmpty || _bodyMetrics.isNotEmpty;

    if (!hasData) {
      return ListView(
        children: [
          const SizedBox(height: 120),
          EmptyState(
            icon: Icons.monitor_heart_outlined,
            message: l10n.noHealthData,
          ),
        ],
      );
    }

    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        _buildBmiSection(l10n),
        const SizedBox(height: 24),
        _buildBodyMeasurements(l10n),
      ],
    );
  }

  Widget _buildBmiSection(AppLocalizations l10n) {
    if (_bmiMetrics.isEmpty) {
      return Card(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            children: [
              const Icon(Icons.monitor_weight_outlined,
                  size: 40, color: AppTheme.textSecondary),
              const SizedBox(height: 8),
              Text(l10n.noBmiData,
                  style: const TextStyle(color: AppTheme.textSecondary)),
            ],
          ),
        ),
      );
    }

    final latest = _bmiMetrics.first;
    final bmiValue = (latest.value['bmi'] as num?)?.toDouble() ?? 0.0;
    final category = _bmiCategory(bmiValue);
    final color = _bmiColor(bmiValue);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          l10n.bmi,
          style: const TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppTheme.textPrimary,
          ),
        ),
        const SizedBox(height: 12),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              children: [
                // BMI value circle
                Container(
                  width: 120,
                  height: 120,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    border: Border.all(color: color, width: 4),
                  ),
                  child: Center(
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Text(
                          bmiValue.toStringAsFixed(1),
                          style: TextStyle(
                            fontSize: 32,
                            fontWeight: FontWeight.bold,
                            color: color,
                          ),
                        ),
                        Text(
                          'BMI',
                          style: TextStyle(
                            fontSize: 12,
                            color: color.withAlpha(180),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                // Category badge
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
                  decoration: BoxDecoration(
                    color: color.withAlpha(30),
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: Text(
                    category,
                    style: TextStyle(
                      color: color,
                      fontWeight: FontWeight.w600,
                      fontSize: 15,
                    ),
                  ),
                ),
                const SizedBox(height: 12),
                // Details
                if (latest.value['height'] != null ||
                    latest.value['weight'] != null)
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      if (latest.value['height'] != null)
                        _buildMetricChip(
                          l10n.height,
                          '${latest.value['height']} cm',
                          Icons.height,
                        ),
                      if (latest.value['height'] != null &&
                          latest.value['weight'] != null)
                        const SizedBox(width: 24),
                      if (latest.value['weight'] != null)
                        _buildMetricChip(
                          l10n.weight,
                          '${latest.value['weight']} kg',
                          Icons.monitor_weight_outlined,
                        ),
                    ],
                  ),
                const SizedBox(height: 8),
                Text(
                  '${l10n.recordedOn}: ${_formatDate(latest.recordedAt)}',
                  style: const TextStyle(
                      color: AppTheme.textSecondary, fontSize: 12),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildMetricChip(String label, String value, IconData icon) {
    return Column(
      children: [
        Icon(icon, size: 20, color: AppTheme.textSecondary),
        const SizedBox(height: 4),
        Text(
          value,
          style: const TextStyle(
            color: AppTheme.textPrimary,
            fontWeight: FontWeight.w600,
            fontSize: 15,
          ),
        ),
        Text(
          label,
          style: const TextStyle(color: AppTheme.textSecondary, fontSize: 12),
        ),
      ],
    );
  }

  Widget _buildBodyMeasurements(AppLocalizations l10n) {
    if (_bodyMetrics.isEmpty) {
      return Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            l10n.bodyMeasurements,
            style: const TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
              color: AppTheme.textPrimary,
            ),
          ),
          const SizedBox(height: 12),
          Card(
            child: Padding(
              padding: const EdgeInsets.all(24),
              child: Center(
                child: Text(l10n.noMeasurementData,
                    style: const TextStyle(color: AppTheme.textSecondary)),
              ),
            ),
          ),
        ],
      );
    }

    final latest = _bodyMetrics.first;
    final measurements = latest.value;

    // Common body measurement keys
    final measurementFields = <String, _MeasurementField>{
      'height': _MeasurementField(l10n.height, Icons.height, 'cm'),
      'weight': _MeasurementField(
          l10n.weight, Icons.monitor_weight_outlined, 'kg'),
      'chest': _MeasurementField(l10n.chest, Icons.accessibility, 'cm'),
      'waist': _MeasurementField(l10n.waist, Icons.straighten, 'cm'),
      'hips': _MeasurementField(l10n.hips, Icons.accessibility_new, 'cm'),
      'biceps': _MeasurementField(l10n.biceps, Icons.fitness_center, 'cm'),
      'thighs': _MeasurementField(l10n.thighs, Icons.directions_walk, 'cm'),
      'body_fat': _MeasurementField(l10n.bodyFat, Icons.percent, '%'),
    };

    final availableMeasurements = measurementFields.entries
        .where((e) => measurements.containsKey(e.key) && measurements[e.key] != null)
        .toList();

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          l10n.bodyMeasurements,
          style: const TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppTheme.textPrimary,
          ),
        ),
        const SizedBox(height: 12),
        if (availableMeasurements.isEmpty)
          Card(
            child: Padding(
              padding: const EdgeInsets.all(24),
              child: Center(
                child: Text(l10n.noMeasurementData,
                    style: const TextStyle(color: AppTheme.textSecondary)),
              ),
            ),
          )
        else
          GridView.count(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            crossAxisCount: 2,
            mainAxisSpacing: 12,
            crossAxisSpacing: 12,
            childAspectRatio: 1.6,
            children: availableMeasurements.map((entry) {
              final field = entry.value;
              final value = measurements[entry.key];
              return Card(
                child: Padding(
                  padding: const EdgeInsets.all(14),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(field.icon,
                          size: 20, color: AppTheme.primary),
                      const SizedBox(height: 8),
                      Text(
                        '$value ${field.unit}',
                        style: const TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                          color: AppTheme.textPrimary,
                        ),
                      ),
                      Text(
                        field.label,
                        style: const TextStyle(
                          color: AppTheme.textSecondary,
                          fontSize: 12,
                        ),
                      ),
                    ],
                  ),
                ),
              );
            }).toList(),
          ),
        if (availableMeasurements.isNotEmpty) ...[
          const SizedBox(height: 8),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 4),
            child: Text(
              '${l10n.recordedOn}: ${_formatDate(latest.recordedAt)}',
              style: const TextStyle(
                  color: AppTheme.textSecondary, fontSize: 12),
            ),
          ),
        ],
      ],
    );
  }

  String _formatDate(String dateStr) {
    final dt = DateTime.tryParse(dateStr);
    if (dt == null) return dateStr;
    return '${dt.year}-${dt.month.toString().padLeft(2, '0')}-${dt.day.toString().padLeft(2, '0')}';
  }
}

class _MeasurementField {
  final String label;
  final IconData icon;
  final String unit;

  const _MeasurementField(this.label, this.icon, this.unit);
}
