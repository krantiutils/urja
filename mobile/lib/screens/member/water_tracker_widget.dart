import 'dart:math';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../config/theme.dart';
import '../../models/water.dart';
import '../../providers/water_provider.dart';
import '../../services/water_service.dart';

class WaterTrackerWidget extends ConsumerStatefulWidget {
  const WaterTrackerWidget({super.key});

  @override
  ConsumerState<WaterTrackerWidget> createState() => _WaterTrackerWidgetState();
}

class _WaterTrackerWidgetState extends ConsumerState<WaterTrackerWidget>
    with SingleTickerProviderStateMixin {
  late final WaterService _waterService;
  WaterDailySummary? _summary;
  bool _loading = true;
  String? _error;
  bool _adding = false;

  late AnimationController _animController;
  late Animation<double> _progressAnim;
  double _prevProgress = 0.0;

  @override
  void initState() {
    super.initState();
    _waterService = ref.read(waterServiceProvider);
    _animController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 600),
    );
    _progressAnim = Tween<double>(begin: 0.0, end: 0.0).animate(
      CurvedAnimation(parent: _animController, curve: Curves.easeInOut),
    );
    _loadSummary();
  }

  @override
  void dispose() {
    _animController.dispose();
    super.dispose();
  }

  String get _today {
    final now = DateTime.now();
    return '${now.year}-${now.month.toString().padLeft(2, '0')}-${now.day.toString().padLeft(2, '0')}';
  }

  Future<void> _loadSummary() async {
    setState(() {
      _loading = _summary == null;
      _error = null;
    });
    try {
      final summary = await _waterService.getDailySummary(_today);
      if (!mounted) return;
      final newProgress = summary.progressPercent;
      _progressAnim = Tween<double>(
        begin: _prevProgress,
        end: newProgress,
      ).animate(
        CurvedAnimation(parent: _animController, curve: Curves.easeInOut),
      );
      _animController.forward(from: 0.0);
      _prevProgress = newProgress;
      setState(() {
        _summary = summary;
        _loading = false;
      });
    } catch (e) {
      if (!mounted) return;
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  Future<void> _addWater(int amountMl) async {
    if (_adding) return;
    setState(() => _adding = true);
    try {
      await _waterService.logWater(amountMl: amountMl, date: _today);
      await _loadSummary();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to log water: $e'),
            backgroundColor: AppTheme.error,
            behavior: SnackBarBehavior.floating,
            shape:
                RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
          ),
        );
      }
    } finally {
      if (mounted) setState(() => _adding = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_loading) {
      return _buildCard(
        child: const SizedBox(
          height: 120,
          child: Center(
            child: CircularProgressIndicator(color: AppTheme.info),
          ),
        ),
      );
    }

    if (_error != null) {
      return _buildCard(
        child: SizedBox(
          height: 120,
          child: Center(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Icon(Icons.error_outline,
                    size: 28, color: AppTheme.error),
                const SizedBox(height: 8),
                const Text(
                  'Could not load water data',
                  style:
                      TextStyle(color: AppTheme.textSecondary, fontSize: 13),
                ),
                const SizedBox(height: 8),
                TextButton(
                  onPressed: _loadSummary,
                  child: const Text('Retry',
                      style: TextStyle(color: AppTheme.info)),
                ),
              ],
            ),
          ),
        ),
      );
    }

    final summary = _summary!;
    return _buildCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header row
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: AppTheme.info.withAlpha(30),
                  borderRadius: BorderRadius.circular(8),
                ),
                child:
                    const Icon(Icons.water_drop, color: AppTheme.info, size: 20),
              ),
              const SizedBox(width: 12),
              const Expanded(
                child: Text(
                  'Water Intake',
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: AppTheme.textPrimary,
                  ),
                ),
              ),
              Text(
                '${summary.totalMl} ml',
                style: const TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  color: AppTheme.info,
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),

          // Circular progress + text
          Center(
            child: SizedBox(
              width: 140,
              height: 140,
              child: AnimatedBuilder(
                animation: _progressAnim,
                builder: (context, child) {
                  return CustomPaint(
                    painter: _WaterRingPainter(
                      progress: _progressAnim.value,
                      color: AppTheme.info,
                      backgroundColor: AppTheme.surfaceLight,
                    ),
                    child: child,
                  );
                },
                child: Center(
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        '${(summary.progressPercent * 100).round()}%',
                        style: const TextStyle(
                          fontSize: 28,
                          fontWeight: FontWeight.bold,
                          color: AppTheme.info,
                        ),
                      ),
                      Text(
                        '${summary.totalMl} / ${summary.goalMl} ml',
                        style: const TextStyle(
                          fontSize: 12,
                          color: AppTheme.textSecondary,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),
          const SizedBox(height: 8),

          // Remaining text
          Center(
            child: Text(
              summary.remainingMl > 0
                  ? '${summary.remainingMl} ml remaining'
                  : 'Goal reached!',
              style: TextStyle(
                fontSize: 13,
                fontWeight: FontWeight.w500,
                color: summary.remainingMl > 0
                    ? AppTheme.textSecondary
                    : AppTheme.primary,
              ),
            ),
          ),
          const SizedBox(height: 16),

          // Quick-add buttons
          Row(
            children: [
              Expanded(
                child: _QuickAddButton(
                  icon: Icons.local_drink_outlined,
                  label: '+250 ml',
                  subtitle: 'Glass',
                  onTap: _adding ? null : () => _addWater(250),
                  isLoading: _adding,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _QuickAddButton(
                  icon: Icons.water_drop_outlined,
                  label: '+500 ml',
                  subtitle: 'Bottle',
                  onTap: _adding ? null : () => _addWater(500),
                  isLoading: _adding,
                ),
              ),
            ],
          ),

          // Recent entries
          if (summary.entries.isNotEmpty) ...[
            const SizedBox(height: 16),
            const Text(
              'Today\'s entries',
              style: TextStyle(
                fontSize: 13,
                fontWeight: FontWeight.w600,
                color: AppTheme.textSecondary,
              ),
            ),
            const SizedBox(height: 8),
            ...summary.entries.reversed.take(5).map((entry) => Padding(
                  padding: const EdgeInsets.only(bottom: 4),
                  child: Row(
                    children: [
                      Icon(Icons.water_drop,
                          size: 14, color: AppTheme.info.withAlpha(150)),
                      const SizedBox(width: 8),
                      Text(
                        '${entry.amountMl} ml',
                        style: const TextStyle(
                          color: AppTheme.textPrimary,
                          fontSize: 13,
                        ),
                      ),
                      const Spacer(),
                      Text(
                        _formatTime(entry.loggedAt),
                        style: const TextStyle(
                          color: AppTheme.textSecondary,
                          fontSize: 12,
                        ),
                      ),
                    ],
                  ),
                )),
          ],
        ],
      ),
    );
  }

  Widget _buildCard({required Widget child}) {
    return Card(
      color: AppTheme.card,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: child,
      ),
    );
  }

  String _formatTime(String loggedAt) {
    try {
      final dt = DateTime.parse(loggedAt).toLocal();
      final hour = dt.hour.toString().padLeft(2, '0');
      final minute = dt.minute.toString().padLeft(2, '0');
      return '$hour:$minute';
    } catch (_) {
      return loggedAt;
    }
  }
}

// ─── Quick Add Button ────────────────────────────────────────────────────────

class _QuickAddButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final String subtitle;
  final VoidCallback? onTap;
  final bool isLoading;

  const _QuickAddButton({
    required this.icon,
    required this.label,
    required this.subtitle,
    required this.onTap,
    this.isLoading = false,
  });

  @override
  Widget build(BuildContext context) {
    return Material(
      color: AppTheme.info.withAlpha(20),
      borderRadius: BorderRadius.circular(12),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 16),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              if (isLoading)
                const SizedBox(
                  width: 18,
                  height: 18,
                  child: CircularProgressIndicator(
                    color: AppTheme.info,
                    strokeWidth: 2,
                  ),
                )
              else
                Icon(icon, color: AppTheme.info, size: 20),
              const SizedBox(width: 8),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    label,
                    style: const TextStyle(
                      color: AppTheme.info,
                      fontWeight: FontWeight.bold,
                      fontSize: 14,
                    ),
                  ),
                  Text(
                    subtitle,
                    style: TextStyle(
                      color: AppTheme.info.withAlpha(150),
                      fontSize: 11,
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
}

// ─── Water Ring Painter ──────────────────────────────────────────────────────

class _WaterRingPainter extends CustomPainter {
  final double progress;
  final Color color;
  final Color backgroundColor;

  _WaterRingPainter({
    required this.progress,
    required this.color,
    required this.backgroundColor,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final radius = min(size.width, size.height) / 2 - 10;
    const strokeWidth = 10.0;

    // Background arc
    final bgPaint = Paint()
      ..color = backgroundColor
      ..style = PaintingStyle.stroke
      ..strokeWidth = strokeWidth
      ..strokeCap = StrokeCap.round;
    canvas.drawCircle(center, radius, bgPaint);

    // Progress arc
    final progressPaint = Paint()
      ..color = color
      ..style = PaintingStyle.stroke
      ..strokeWidth = strokeWidth
      ..strokeCap = StrokeCap.round;
    final sweepAngle = 2 * pi * progress.clamp(0.0, 1.0);
    canvas.drawArc(
      Rect.fromCircle(center: center, radius: radius),
      -pi / 2,
      sweepAngle,
      false,
      progressPaint,
    );
  }

  @override
  bool shouldRepaint(covariant _WaterRingPainter oldDelegate) =>
      oldDelegate.progress != progress || oldDelegate.color != color;
}
