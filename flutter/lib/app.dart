import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:urja_flutter/l10n/app_localizations.dart';

import 'config/routes.dart';
import 'config/theme.dart';
import 'providers/providers.dart';

class UrjaApp extends ConsumerStatefulWidget {
  const UrjaApp({super.key});

  @override
  ConsumerState<UrjaApp> createState() => _UrjaAppState();
}

class _UrjaAppState extends ConsumerState<UrjaApp> {
  @override
  void initState() {
    super.initState();
    ref.read(authInitProvider);
    ref.read(syncManagerProvider).startListening();
  }

  @override
  Widget build(BuildContext context) {
    final router = ref.watch(routerProvider);
    final locale = ref.watch(localeProvider);

    return MaterialApp.router(
      title: 'Urja',
      debugShowCheckedModeBanner: false,
      theme: urjaTheme(),
      locale: Locale(locale),
      supportedLocales: AppLocalizations.supportedLocales,
      localizationsDelegates: const [
        AppLocalizations.delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      routerConfig: router,
    );
  }
}
