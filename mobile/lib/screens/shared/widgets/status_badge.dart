import 'package:flutter/material.dart';
import '../../../config/theme.dart';

class StatusBadge extends StatelessWidget {
  final String label;
  final Color? color;

  const StatusBadge({super.key, required this.label, this.color});

  Color get _color {
    if (color != null) return color!;
    switch (label.toLowerCase()) {
      case 'active':
      case 'paid':
      case 'online':
      case 'available':
        return AppTheme.primary;
      case 'suspended':
      case 'unpaid':
      case 'expired':
      case 'offline':
        return AppTheme.error;
      case 'waived':
      case 'pending':
        return AppTheme.warning;
      default:
        return AppTheme.textSecondary;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: _color.withAlpha(30),
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        label,
        style: TextStyle(color: _color, fontSize: 12, fontWeight: FontWeight.w600),
      ),
    );
  }
}
