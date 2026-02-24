import 'package:flutter/material.dart';
import '../../../config/theme.dart';

class MethodBadge extends StatelessWidget {
  final String method;

  const MethodBadge({super.key, required this.method});

  IconData get _icon {
    switch (method.toLowerCase()) {
      case 'qr':
        return Icons.qr_code;
      case 'nfc':
        return Icons.nfc;
      case 'manual':
        return Icons.touch_app;
      default:
        return Icons.check;
    }
  }

  Color get _color {
    switch (method.toLowerCase()) {
      case 'qr':
        return AppTheme.info;
      case 'nfc':
        return AppTheme.primary;
      case 'manual':
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
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(_icon, size: 14, color: _color),
          const SizedBox(width: 4),
          Text(method.toUpperCase(),
              style: TextStyle(
                  color: _color, fontSize: 11, fontWeight: FontWeight.w600)),
        ],
      ),
    );
  }
}
