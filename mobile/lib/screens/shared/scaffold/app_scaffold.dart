import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../config/theme.dart';
import '../../../providers/auth_provider.dart';

class AppScaffold extends ConsumerWidget {
  final String role;
  final Widget child;

  const AppScaffold({super.key, required this.role, required this.child});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      body: child,
      bottomNavigationBar: _buildBottomNav(context, ref),
    );
  }

  Widget _buildBottomNav(BuildContext context, WidgetRef ref) {
    final loc = GoRouterState.of(context).uri.toString();

    if (role == 'admin') {
      return _AdminBottomNav(currentLocation: loc);
    } else if (role == 'tracker') {
      final userType = ref.watch(authProvider).user?.userType;
      return _TrackerBottomNav(
        currentLocation: loc,
        isCalorieTracker: userType == 'calorie_tracker',
      );
    } else {
      return _MemberBottomNav(currentLocation: loc);
    }
  }
}

// ─── Member Bottom Nav ───────────────────────────────────────────────────────

class _MemberBottomNav extends StatelessWidget {
  final String currentLocation;

  const _MemberBottomNav({required this.currentLocation});

  int get _currentIndex {
    if (currentLocation.startsWith('/member/nutrition')) {
      return 1;
    } else if (currentLocation.startsWith('/member/workouts')) {
      return 3;
    } else if (currentLocation.startsWith('/member/settings') ||
        currentLocation.startsWith('/member/packages') ||
        currentLocation.startsWith('/member/health')) {
      return 4;
    } else if (currentLocation.startsWith('/member')) {
      return 0;
    }
    return 0;
  }

  @override
  Widget build(BuildContext context) {
    final int selected = _currentIndex;
    final bottomPadding = MediaQuery.of(context).padding.bottom;

    return Container(
      decoration: BoxDecoration(
        color: AppTheme.surface,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withAlpha(40),
            blurRadius: 8,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: Padding(
        padding: EdgeInsets.only(bottom: bottomPadding),
        child: SizedBox(
          height: 60,
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: [
              _NavItem(
                icon: Icons.home_outlined,
                activeIcon: Icons.home,
                label: 'Home',
                isSelected: selected == 0,
                onTap: () => context.go('/member'),
              ),
              _NavItem(
                icon: Icons.restaurant_outlined,
                activeIcon: Icons.restaurant,
                label: 'Nutrition',
                isSelected: selected == 1,
                onTap: () => context.go('/member/nutrition'),
              ),
              _QrScanNavItem(
                onTap: () => context.push('/check-in/qr'),
              ),
              _NavItem(
                icon: Icons.fitness_center_outlined,
                activeIcon: Icons.fitness_center,
                label: 'Workouts',
                isSelected: selected == 3,
                onTap: () => context.go('/member/workouts'),
              ),
              _NavItem(
                icon: Icons.person_outline,
                activeIcon: Icons.person,
                label: 'Account',
                isSelected: selected == 4,
                onTap: () => context.go('/member/settings'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _NavItem extends StatelessWidget {
  final IconData icon;
  final IconData activeIcon;
  final String label;
  final bool isSelected;
  final VoidCallback onTap;

  const _NavItem({
    required this.icon,
    required this.activeIcon,
    required this.label,
    required this.isSelected,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: GestureDetector(
        onTap: onTap,
        behavior: HitTestBehavior.opaque,
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              isSelected ? activeIcon : icon,
              color: isSelected ? AppTheme.primary : AppTheme.textSecondary,
              size: 24,
            ),
            const SizedBox(height: 4),
            Text(
              label,
              style: TextStyle(
                fontSize: 12,
                color: isSelected ? AppTheme.primary : AppTheme.textSecondary,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _QrScanNavItem extends StatelessWidget {
  final VoidCallback onTap;

  const _QrScanNavItem({required this.onTap});

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: GestureDetector(
        onTap: onTap,
        behavior: HitTestBehavior.opaque,
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                color: AppTheme.primary,
                shape: BoxShape.circle,
                boxShadow: [
                  BoxShadow(
                    color: AppTheme.primary.withAlpha(80),
                    blurRadius: 8,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              child: const Icon(
                Icons.qr_code_scanner,
                color: Colors.white,
                size: 26,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

// ─── Tracker Bottom Nav ──────────────────────────────────────────────────────

class _TrackerBottomNav extends StatelessWidget {
  final String currentLocation;
  final bool isCalorieTracker;

  const _TrackerBottomNav({
    required this.currentLocation,
    this.isCalorieTracker = false,
  });

  @override
  Widget build(BuildContext context) {
    if (isCalorieTracker) {
      return _buildCalorieTrackerNav(context);
    }
    return _buildFitnessTrackerNav(context);
  }

  Widget _buildCalorieTrackerNav(BuildContext context) {
    int currentIndex;
    if (currentLocation.startsWith('/tracker/nutrition')) {
      currentIndex = 1;
    } else if (currentLocation.startsWith('/tracker/progress')) {
      currentIndex = 2;
    } else if (currentLocation.startsWith('/tracker/settings')) {
      currentIndex = 3;
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
            context.go('/tracker/progress');
            break;
          case 3:
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

  Widget _buildFitnessTrackerNav(BuildContext context) {
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
