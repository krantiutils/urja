import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:urja_flutter/config/theme.dart';

void main() {
  testWidgets('Urja theme smoke test', (WidgetTester tester) async {
    final theme = urjaTheme();
    expect(theme.brightness, Brightness.dark);
    expect(theme.primaryColor, UrjaColors.accent);

    await tester.pumpWidget(
      ProviderScope(
        child: MaterialApp(
          theme: theme,
          home: const Scaffold(
            body: Center(child: Text('Urja')),
          ),
        ),
      ),
    );

    expect(find.text('Urja'), findsOneWidget);
  });
}
