import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/auth_provider.dart';
import '../screens/auth/login_screen.dart';
import '../screens/onboarding/onboarding_screen.dart';
import '../screens/member/member_home_screen.dart';
import '../screens/member/member_attendance_screen.dart';
import '../screens/member/member_packages_screen.dart';
import '../screens/member/member_health_screen.dart';
import '../screens/member/member_workouts_screen.dart';
import '../screens/member/member_nutrition_screen.dart';
import '../screens/member/nutrition_dashboard_screen.dart';
import '../screens/member/member_settings_screen.dart';
import '../screens/member/member_progress_screen.dart';
import '../screens/tracker/tracker_home_screen.dart';
import '../screens/tracker/tracker_progress_screen.dart';
import '../screens/training/exercise_library_screen.dart';
import '../screens/training/program_list_screen.dart';
import '../screens/training/program_detail_screen.dart';
import '../screens/training/hiit_timer_screen.dart';
import '../screens/dashboard/dashboard_home_screen.dart';
import '../screens/dashboard/members_screen.dart';
import '../screens/dashboard/attendance_screen.dart';
import '../screens/dashboard/packages_screen.dart';
import '../screens/dashboard/staff_screen.dart';
import '../screens/dashboard/due_payments_screen.dart';
import '../screens/dashboard/accounts_screen.dart';
import '../screens/dashboard/notices_screen.dart';
import '../screens/dashboard/feedbacks_screen.dart';
import '../screens/dashboard/nfc_screen.dart';
import '../screens/dashboard/sms_screen.dart';
import '../screens/dashboard/settings_screen.dart';
import '../screens/check_in/qr_scan_screen.dart';
import '../screens/check_in/nfc_tap_screen.dart';
import '../screens/shared/scaffold/app_scaffold.dart';

final _rootNavigatorKey = GlobalKey<NavigatorState>();
final _memberShellKey = GlobalKey<NavigatorState>();
final _trackerShellKey = GlobalKey<NavigatorState>();
final _dashboardShellKey = GlobalKey<NavigatorState>();

String _homeRouteForUser(AuthState auth) {
  final role = auth.user?.role ?? 'member';
  final userType = auth.user?.userType;

  // Admin/staff go to dashboard
  if (role != 'member') return '/dashboard';

  // Fitness tracker and calorie tracker go to tracker home
  if (userType == 'fitness_tracker' || userType == 'calorie_tracker') {
    return '/tracker';
  }

  // Default: gym member
  return '/member';
}

final routerProvider = Provider<GoRouter>((ref) {
  final auth = ref.watch(authProvider);

  return GoRouter(
    navigatorKey: _rootNavigatorKey,
    initialLocation: '/login',
    redirect: (context, state) {
      // Don't redirect while auth state is loading
      if (auth.status == AuthStatus.initial) return null;

      final isLoggedIn = auth.status == AuthStatus.authenticated;
      final isLoginPage = state.matchedLocation == '/login';
      final isOnboardingPage = state.matchedLocation == '/onboarding';

      // Not logged in -> go to login (unless already there)
      if (!isLoggedIn && !isLoginPage) return '/login';

      // Logged in but onboarding not completed -> go to onboarding
      if (isLoggedIn && !auth.onboardingCompleted && !isOnboardingPage) {
        return '/onboarding';
      }

      // Logged in + onboarding done + on login or onboarding page -> go home
      if (isLoggedIn && auth.onboardingCompleted && (isLoginPage || isOnboardingPage)) {
        return _homeRouteForUser(auth);
      }

      return null;
    },
    routes: [
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginScreen(),
      ),
      // Onboarding (full screen, no shell)
      GoRoute(
        path: '/onboarding',
        builder: (context, state) => const OnboardingScreen(),
      ),
      // Member routes
      ShellRoute(
        navigatorKey: _memberShellKey,
        builder: (context, state, child) =>
            AppScaffold(role: 'member', child: child),
        routes: [
          GoRoute(
            path: '/member',
            builder: (context, state) => const MemberHomeScreen(),
          ),
          GoRoute(
            path: '/member/attendance',
            builder: (context, state) => const MemberAttendanceScreen(),
          ),
          GoRoute(
            path: '/member/packages',
            builder: (context, state) => const MemberPackagesScreen(),
          ),
          GoRoute(
            path: '/member/health',
            builder: (context, state) => const MemberHealthScreen(),
          ),
          GoRoute(
            path: '/member/workouts',
            builder: (context, state) => const MemberWorkoutsScreen(),
          ),
          GoRoute(
            path: '/member/nutrition',
            builder: (context, state) => const NutritionDashboardScreen(),
          ),
          GoRoute(
            path: '/member/nutrition/log',
            builder: (context, state) => const MemberNutritionScreen(),
          ),
          GoRoute(
            path: '/member/progress',
            builder: (context, state) => const MemberProgressScreen(),
          ),
          GoRoute(
            path: '/member/exercises',
            builder: (context, state) => const ExerciseLibraryScreen(),
          ),
          GoRoute(
            path: '/member/programs',
            builder: (context, state) => const ProgramListScreen(),
          ),
          GoRoute(
            path: '/member/programs/:id',
            builder: (context, state) => ProgramDetailScreen(
              programId: state.pathParameters['id']!,
            ),
          ),
          GoRoute(
            path: '/member/hiit-timer',
            builder: (context, state) => const HiitTimerScreen(),
          ),
          GoRoute(
            path: '/member/settings',
            builder: (context, state) => const MemberSettingsScreen(),
          ),
        ],
      ),
      // Tracker routes (fitness_tracker / calorie_tracker)
      ShellRoute(
        navigatorKey: _trackerShellKey,
        builder: (context, state, child) =>
            AppScaffold(role: 'tracker', child: child),
        routes: [
          GoRoute(
            path: '/tracker',
            builder: (context, state) => const TrackerHomeScreen(),
          ),
          GoRoute(
            path: '/tracker/nutrition',
            builder: (context, state) => const NutritionDashboardScreen(),
          ),
          GoRoute(
            path: '/tracker/nutrition/log',
            builder: (context, state) => const MemberNutritionScreen(),
          ),
          GoRoute(
            path: '/tracker/workouts',
            builder: (context, state) => const MemberWorkoutsScreen(),
          ),
          GoRoute(
            path: '/tracker/progress',
            builder: (context, state) => const TrackerProgressScreen(),
          ),
          GoRoute(
            path: '/tracker/exercises',
            builder: (context, state) => const ExerciseLibraryScreen(),
          ),
          GoRoute(
            path: '/tracker/programs',
            builder: (context, state) => const ProgramListScreen(),
          ),
          GoRoute(
            path: '/tracker/programs/:id',
            builder: (context, state) => ProgramDetailScreen(
              programId: state.pathParameters['id']!,
            ),
          ),
          GoRoute(
            path: '/tracker/hiit-timer',
            builder: (context, state) => const HiitTimerScreen(),
          ),
          GoRoute(
            path: '/tracker/settings',
            builder: (context, state) => const MemberSettingsScreen(),
          ),
        ],
      ),
      // Dashboard routes (admin/staff)
      ShellRoute(
        navigatorKey: _dashboardShellKey,
        builder: (context, state, child) =>
            AppScaffold(role: 'admin', child: child),
        routes: [
          GoRoute(
            path: '/dashboard',
            builder: (context, state) => const DashboardHomeScreen(),
          ),
          GoRoute(
            path: '/dashboard/members',
            builder: (context, state) => const MembersScreen(),
          ),
          GoRoute(
            path: '/dashboard/attendance',
            builder: (context, state) => const AttendanceScreen(),
          ),
          GoRoute(
            path: '/dashboard/packages',
            builder: (context, state) => const PackagesScreen(),
          ),
          GoRoute(
            path: '/dashboard/staff',
            builder: (context, state) => const StaffScreen(),
          ),
          GoRoute(
            path: '/dashboard/dues',
            builder: (context, state) => const DuePaymentsScreen(),
          ),
          GoRoute(
            path: '/dashboard/accounts',
            builder: (context, state) => const AccountsScreen(),
          ),
          GoRoute(
            path: '/dashboard/notices',
            builder: (context, state) => const NoticesScreen(),
          ),
          GoRoute(
            path: '/dashboard/feedbacks',
            builder: (context, state) => const FeedbacksScreen(),
          ),
          GoRoute(
            path: '/dashboard/nfc',
            builder: (context, state) => const NfcScreen(),
          ),
          GoRoute(
            path: '/dashboard/sms',
            builder: (context, state) => const SmsScreen(),
          ),
          GoRoute(
            path: '/dashboard/settings',
            builder: (context, state) => const SettingsScreen(),
          ),
        ],
      ),
      // Check-in routes (full screen, no shell)
      GoRoute(
        path: '/check-in/qr',
        builder: (context, state) => const QrScanScreen(),
      ),
      GoRoute(
        path: '/check-in/nfc',
        builder: (context, state) => const NfcTapScreen(),
      ),
    ],
  );
});
