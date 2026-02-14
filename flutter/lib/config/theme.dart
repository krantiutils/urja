import 'package:flutter/material.dart';

/// Urja design system — Linear/Modern dark theme
/// Colors from UI.md design system
class UrjaColors {
  static const backgroundDeep = Color(0xFF020203);
  static const backgroundBase = Color(0xFF050506);
  static const backgroundElevated = Color(0xFF0A0A0C);
  static const surface = Color(0x0DFFFFFF); // rgba(255,255,255,0.05)
  static const surfaceHover = Color(0x14FFFFFF); // rgba(255,255,255,0.08)
  static const foreground = Color(0xFFEDEDEF);
  static const foregroundMuted = Color(0xFF8A8F98);
  static const foregroundSubtle = Color(0x99FFFFFF); // 60%
  static const accent = Color(0xFF5E6AD2);
  static const accentBright = Color(0xFF6872D9);
  static const accentGlow = Color(0x4D5E6AD2); // 30%
  static const borderDefault = Color(0x0FFFFFFF); // 6%
  static const borderHover = Color(0x1AFFFFFF); // 10%
  static const borderAccent = Color(0x4D5E6AD2); // 30%
  static const inputBg = Color(0xFF0F0F12);
  static const error = Color(0xFFE5484D);
  static const success = Color(0xFF30A46C);
  static const warning = Color(0xFFF5A623);
}

ThemeData urjaTheme() {
  return ThemeData(
    brightness: Brightness.dark,
    scaffoldBackgroundColor: UrjaColors.backgroundBase,
    primaryColor: UrjaColors.accent,
    colorScheme: const ColorScheme.dark(
      primary: UrjaColors.accent,
      secondary: UrjaColors.accentBright,
      surface: UrjaColors.backgroundElevated,
      error: UrjaColors.error,
      onPrimary: Colors.white,
      onSecondary: Colors.white,
      onSurface: UrjaColors.foreground,
      onError: Colors.white,
    ),
    fontFamily: 'Inter',
    textTheme: const TextTheme(
      displayLarge: TextStyle(
        fontSize: 48,
        fontWeight: FontWeight.w600,
        letterSpacing: -1.5,
        color: UrjaColors.foreground,
      ),
      headlineLarge: TextStyle(
        fontSize: 32,
        fontWeight: FontWeight.w600,
        letterSpacing: -0.5,
        color: UrjaColors.foreground,
      ),
      headlineMedium: TextStyle(
        fontSize: 24,
        fontWeight: FontWeight.w600,
        letterSpacing: -0.3,
        color: UrjaColors.foreground,
      ),
      headlineSmall: TextStyle(
        fontSize: 20,
        fontWeight: FontWeight.w600,
        color: UrjaColors.foreground,
      ),
      titleLarge: TextStyle(
        fontSize: 18,
        fontWeight: FontWeight.w600,
        color: UrjaColors.foreground,
      ),
      titleMedium: TextStyle(
        fontSize: 16,
        fontWeight: FontWeight.w500,
        color: UrjaColors.foreground,
      ),
      bodyLarge: TextStyle(
        fontSize: 16,
        fontWeight: FontWeight.w400,
        color: UrjaColors.foregroundMuted,
        height: 1.6,
      ),
      bodyMedium: TextStyle(
        fontSize: 14,
        fontWeight: FontWeight.w400,
        color: UrjaColors.foregroundMuted,
        height: 1.5,
      ),
      bodySmall: TextStyle(
        fontSize: 12,
        fontWeight: FontWeight.w400,
        color: UrjaColors.foregroundMuted,
      ),
      labelSmall: TextStyle(
        fontSize: 11,
        fontWeight: FontWeight.w500,
        letterSpacing: 1.2,
        color: UrjaColors.foregroundMuted,
      ),
    ),
    appBarTheme: const AppBarTheme(
      backgroundColor: UrjaColors.backgroundBase,
      foregroundColor: UrjaColors.foreground,
      elevation: 0,
      centerTitle: true,
      titleTextStyle: TextStyle(
        fontSize: 18,
        fontWeight: FontWeight.w600,
        color: UrjaColors.foreground,
      ),
    ),
    cardTheme: CardThemeData(
      color: UrjaColors.backgroundElevated,
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
        side: const BorderSide(color: UrjaColors.borderDefault),
      ),
    ),
    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: UrjaColors.accent,
        foregroundColor: Colors.white,
        elevation: 0,
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 14),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        textStyle: const TextStyle(
          fontSize: 15,
          fontWeight: FontWeight.w600,
        ),
      ),
    ),
    outlinedButtonTheme: OutlinedButtonThemeData(
      style: OutlinedButton.styleFrom(
        foregroundColor: UrjaColors.foreground,
        side: const BorderSide(color: UrjaColors.borderDefault),
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 14),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
      ),
    ),
    textButtonTheme: TextButtonThemeData(
      style: TextButton.styleFrom(
        foregroundColor: UrjaColors.accent,
      ),
    ),
    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: UrjaColors.inputBg,
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
        borderSide: const BorderSide(color: UrjaColors.borderDefault),
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
        borderSide: const BorderSide(color: UrjaColors.borderDefault),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
        borderSide: const BorderSide(color: UrjaColors.accent, width: 2),
      ),
      errorBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
        borderSide: const BorderSide(color: UrjaColors.error),
      ),
      hintStyle: const TextStyle(color: UrjaColors.foregroundSubtle),
      labelStyle: const TextStyle(color: UrjaColors.foregroundMuted),
      contentPadding:
          const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
    ),
    bottomNavigationBarTheme: const BottomNavigationBarThemeData(
      backgroundColor: UrjaColors.backgroundDeep,
      selectedItemColor: UrjaColors.accent,
      unselectedItemColor: UrjaColors.foregroundMuted,
      type: BottomNavigationBarType.fixed,
      elevation: 0,
    ),
    dividerTheme: const DividerThemeData(
      color: UrjaColors.borderDefault,
      thickness: 1,
    ),
    snackBarTheme: SnackBarThemeData(
      backgroundColor: UrjaColors.backgroundElevated,
      contentTextStyle: const TextStyle(color: UrjaColors.foreground),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8),
        side: const BorderSide(color: UrjaColors.borderDefault),
      ),
      behavior: SnackBarBehavior.floating,
    ),
  );
}
