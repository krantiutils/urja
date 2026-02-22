import 'package:flutter/material.dart';
import '../../../config/theme.dart';

class EmptyState extends StatelessWidget {
  final IconData icon;
  final String message;

  const EmptyState({
    super.key,
    required this.icon,
    required this.message,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: 48, color: AppTheme.textSecondary),
        const SizedBox(height: 16),
        Text(
          message,
          style: const TextStyle(
            color: AppTheme.textSecondary,
            fontSize: 15,
          ),
          textAlign: TextAlign.center,
        ),
      ],
    );
  }
}
