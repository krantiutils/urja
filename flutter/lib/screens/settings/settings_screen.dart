import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:urja_flutter/l10n/app_localizations.dart';

import '../../config/theme.dart';
import '../../providers/providers.dart';

class SettingsScreen extends ConsumerWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = AppLocalizations.of(context)!;
    final currentLocale = ref.watch(localeProvider);

    return Scaffold(
      appBar: AppBar(title: Text(l10n.settings)),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // Language toggle
          Card(
            child: ListTile(
              leading:
                  const Icon(Icons.language, color: UrjaColors.accent),
              title: Text(l10n.language),
              trailing: SegmentedButton<String>(
                segments: [
                  ButtonSegment(value: 'en', label: Text(l10n.english)),
                  ButtonSegment(value: 'ne', label: Text(l10n.nepali)),
                ],
                selected: {currentLocale},
                onSelectionChanged: (selected) {
                  ref.read(localeProvider.notifier).state =
                      selected.first;
                },
                style: ButtonStyle(
                  backgroundColor: WidgetStateProperty.resolveWith((states) {
                    if (states.contains(WidgetState.selected)) {
                      return UrjaColors.accent;
                    }
                    return UrjaColors.surface;
                  }),
                ),
              ),
            ),
          ),
          const SizedBox(height: 8),

          // Privacy settings
          Card(
            child: ListTile(
              leading: const Icon(Icons.privacy_tip_outlined,
                  color: UrjaColors.accent),
              title: Text(l10n.privacySettings),
              trailing: const Icon(Icons.chevron_right,
                  color: UrjaColors.foregroundMuted),
              onTap: () {
                // TODO: Navigate to privacy settings
              },
            ),
          ),
        ],
      ),
    );
  }
}
