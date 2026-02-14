import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:urja_flutter/l10n/app_localizations.dart';

import '../../config/theme.dart';

class HealthScreen extends ConsumerStatefulWidget {
  const HealthScreen({super.key});

  @override
  ConsumerState<HealthScreen> createState() => _HealthScreenState();
}

class _HealthScreenState extends ConsumerState<HealthScreen> {
  final _heightController = TextEditingController();
  final _weightController = TextEditingController();

  @override
  void dispose() {
    _heightController.dispose();
    _weightController.dispose();
    super.dispose();
  }

  double? _calculateBmi() {
    final height = double.tryParse(_heightController.text);
    final weight = double.tryParse(_weightController.text);
    if (height == null || weight == null || height <= 0) return null;
    final heightM = height / 100;
    return weight / (heightM * heightM);
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(title: Text(l10n.health)),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // BMI Calculator
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(l10n.logBmi,
                      style: Theme.of(context).textTheme.titleMedium),
                  const SizedBox(height: 16),
                  TextFormField(
                    key: const Key('height_field'),
                    controller: _heightController,
                    decoration: InputDecoration(labelText: l10n.heightCm),
                    keyboardType: TextInputType.number,
                  ),
                  const SizedBox(height: 12),
                  TextFormField(
                    key: const Key('weight_field'),
                    controller: _weightController,
                    decoration: InputDecoration(labelText: l10n.weightKg),
                    keyboardType: TextInputType.number,
                  ),
                  const SizedBox(height: 16),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      key: const Key('log_bmi_button'),
                      onPressed: () {
                        final bmi = _calculateBmi();
                        if (bmi != null) {
                          ScaffoldMessenger.of(context).showSnackBar(
                            SnackBar(
                              content: Text(
                                  'BMI: ${bmi.toStringAsFixed(1)}'),
                            ),
                          );
                        }
                      },
                      child: Text(l10n.save),
                    ),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),

          // Quick links
          _HealthTile(
            icon: Icons.straighten,
            title: l10n.bodyMeasurements,
            onTap: () {},
          ),
          _HealthTile(
            icon: Icons.photo_library_outlined,
            title: l10n.progressPhotos,
            onTap: () {},
          ),
          _HealthTile(
            icon: Icons.monitor_weight_outlined,
            title: l10n.weight,
            onTap: () {},
          ),
        ],
      ),
    );
  }
}

class _HealthTile extends StatelessWidget {
  final IconData icon;
  final String title;
  final VoidCallback onTap;

  const _HealthTile({
    required this.icon,
    required this.title,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: Icon(icon, color: UrjaColors.accent),
        title: Text(title),
        trailing: const Icon(Icons.chevron_right,
            color: UrjaColors.foregroundMuted),
        onTap: onTap,
      ),
    );
  }
}
