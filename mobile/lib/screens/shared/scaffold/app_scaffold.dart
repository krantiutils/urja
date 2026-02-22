import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../config/theme.dart';

class AppScaffold extends ConsumerWidget {
  final String role;
  final Widget child;

  const AppScaffold({super.key, required this.role, required this.child});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      body: child,
      bottomNavigationBar: _buildBottomNav(context),
    );
  }

  Widget _buildBottomNav(BuildContext context) {
    final loc = GoRouterState.of(context).uri.toString();

    if (role == 'admin') {
      return _AdminBottomNav(currentLocation: loc);
    } else if (role == 'tracker') {
      return _TrackerBottomNav(currentLocation: loc);
    } else {
      return _MemberBottomNav(currentLocation: loc);
    }
  }
}

// ─── Member Bottom Nav ───────────────────────────────────────────────────────

class _MemberBottomNav extends StatelessWidget {
  final String currentLocation;

  const _MemberBottomNav({required this.currentLocation});

  @override
  Widget build(BuildContext context) {
    int currentIndex;
    if (currentLocation.startsWith('/member/attendance')) {
      currentIndex = 1;
    } else if (currentLocation.startsWith('/member/nutrition')) {
      currentIndex = 2;
    } else if (currentLocation.startsWith('/member/workouts')) {
      currentIndex = 3;
    } else if (currentLocation.startsWith('/member/settings') ||
        currentLocation.startsWith('/member/packages') ||
        currentLocation.startsWith('/member/health')) {
      currentIndex = 4;
    } else {
      currentIndex = 0;
    }

    return BottomNavigationBar(
      currentIndex: currentIndex,
      backgroundColor: AppTheme.surface,
      selectedItemColor: AppTheme.primary,
      unselectedItemColor: AppTheme.textSecondary,
      type: BottomNavigationBarType.fixed,
      selectedFontSize: 12,
      unselectedFontSize: 12,
      onTap: (index) {
        switch (index) {
          case 0:
            context.go('/member');
            break;
          case 1:
            context.go('/member/attendance');
            break;
          case 2:
            context.go('/member/nutrition');
            break;
          case 3:
            context.go('/member/workouts');
            break;
          case 4:
            context.go('/member/settings');
            break;
        }
      },
      items: const [
        BottomNavigationBarItem(
          icon: Icon(Icons.home_outlined),
          activeIcon: Icon(Icons.home),
          label: 'Home',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.calendar_today_outlined),
          activeIcon: Icon(Icons.calendar_today),
          label: 'Attendance',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.restaurant_outlined),
          activeIcon: Icon(Icons.restaurant),
          label: 'Nutrition',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.fitness_center_outlined),
          activeIcon: Icon(Icons.fitness_center),
          label: 'Workouts',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.more_horiz_outlined),
          activeIcon: Icon(Icons.more_horiz),
          label: 'More',
        ),
      ],
    );
  }
}

// ─── Tracker Bottom Nav ──────────────────────────────────────────────────────

class _TrackerBottomNav extends StatelessWidget {
  final String currentLocation;

  const _TrackerBottomNav({required this.currentLocation});

  @override
  Widget build(BuildContext context) {
    int currentIndex;
    if (currentLocation.startsWith('/tracker/nutrition')) {
      currentIndex = 1;
    } else if (currentLocation.startsWith('/tracker/workouts')) {
      currentIndex = 2;
    } else if (currentLocation.startsWith('/tracker/progress')) {
      currentIndex = 3;
    } else if (currentLocation.startsWith('/tracker/settings')) {
      currentIndex = 4;
    } else {
      currentIndex = 0;
    }

    return BottomNavigationBar(
      currentIndex: currentIndex,
      backgroundColor: AppTheme.surface,
      selectedItemColor: AppTheme.primary,
      unselectedItemColor: AppTheme.textSecondary,
      type: BottomNavigationBarType.fixed,
      selectedFontSize: 12,
      unselectedFontSize: 12,
      onTap: (index) {
        switch (index) {
          case 0:
            context.go('/tracker');
            break;
          case 1:
            context.go('/tracker/nutrition');
            break;
          case 2:
            context.go('/tracker/workouts');
            break;
          case 3:
            context.go('/tracker/progress');
            break;
          case 4:
            context.go('/tracker/settings');
            break;
        }
      },
      items: const [
        BottomNavigationBarItem(
          icon: Icon(Icons.home_outlined),
          activeIcon: Icon(Icons.home),
          label: 'Home',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.restaurant_outlined),
          activeIcon: Icon(Icons.restaurant),
          label: 'Nutrition',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.fitness_center_outlined),
          activeIcon: Icon(Icons.fitness_center),
          label: 'Workouts',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.trending_up_outlined),
          activeIcon: Icon(Icons.trending_up),
          label: 'Progress',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.person_outline),
          activeIcon: Icon(Icons.person),
          label: 'Account',
        ),
      ],
    );
  }
}

// ─── Admin Bottom Nav ────────────────────────────────────────────────────────

class _AdminBottomNav extends StatelessWidget {
  final String currentLocation;

  const _AdminBottomNav({required this.currentLocation});

  @override
  Widget build(BuildContext context) {
    int currentIndex;
    if (currentLocation.startsWith('/dashboard/members')) {
      currentIndex = 1;
    } else if (currentLocation.startsWith('/dashboard/attendance')) {
      currentIndex = 2;
    } else if (currentLocation.startsWith('/dashboard/packages')) {
      currentIndex = 3;
    } else if (currentLocation.startsWith('/dashboard/staff') ||
        currentLocation.startsWith('/dashboard/dues') ||
        currentLocation.startsWith('/dashboard/accounts') ||
        currentLocation.startsWith('/dashboard/notices') ||
        currentLocation.startsWith('/dashboard/feedbacks') ||
        currentLocation.startsWith('/dashboard/nfc') ||
        currentLocation.startsWith('/dashboard/sms') ||
        currentLocation.startsWith('/dashboard/settings')) {
      currentIndex = 4;
    } else {
      currentIndex = 0;
    }

    return BottomNavigationBar(
      currentIndex: currentIndex,
      backgroundColor: AppTheme.surface,
      selectedItemColor: AppTheme.primary,
      unselectedItemColor: AppTheme.textSecondary,
      type: BottomNavigationBarType.fixed,
      selectedFontSize: 12,
      unselectedFontSize: 12,
      onTap: (index) {
        switch (index) {
          case 0:
            context.go('/dashboard');
            break;
          case 1:
            context.go('/dashboard/members');
            break;
          case 2:
            context.go('/dashboard/attendance');
            break;
          case 3:
            context.go('/dashboard/packages');
            break;
          case 4:
            context.go('/dashboard/settings');
            break;
        }
      },
      items: const [
        BottomNavigationBarItem(
          icon: Icon(Icons.dashboard_outlined),
          activeIcon: Icon(Icons.dashboard),
          label: 'Dashboard',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.people_outline),
          activeIcon: Icon(Icons.people),
          label: 'Members',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.calendar_today_outlined),
          activeIcon: Icon(Icons.calendar_today),
          label: 'Attendance',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.card_membership_outlined),
          activeIcon: Icon(Icons.card_membership),
          label: 'Packages',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.more_horiz_outlined),
          activeIcon: Icon(Icons.more_horiz),
          label: 'More',
        ),
      ],
    );
  }
}
