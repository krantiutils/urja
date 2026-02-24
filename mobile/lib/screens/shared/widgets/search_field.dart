import 'package:flutter/material.dart';
import 'package:mobile/l10n/app_localizations.dart';

class SearchField extends StatelessWidget {
  final ValueChanged<String> onChanged;
  final String? hint;

  const SearchField({super.key, required this.onChanged, this.hint});

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    return TextField(
      decoration: InputDecoration(
        hintText: hint ?? l10n.search,
        prefixIcon: const Icon(Icons.search),
      ),
      onChanged: onChanged,
    );
  }
}
