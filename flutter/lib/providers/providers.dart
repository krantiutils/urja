import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../core/api/api_client.dart';
import '../core/auth/secure_storage.dart';
import '../core/database/local_db.dart';
import '../core/sync/sync_manager.dart';
import '../models/attendance.dart';
import '../models/package.dart';
import '../models/user.dart';

// --- Singletons ---

final secureStorageProvider = Provider<SecureTokenStorage>((ref) {
  return SecureTokenStorage();
});

final apiClientProvider = Provider<ApiClient>((ref) {
  final storage = ref.watch(secureStorageProvider);
  final client = ApiClient(tokenStorage: storage);
  client.onSessionExpired = () {
    ref.read(authStateProvider.notifier).state = AuthState.unauthenticated;
  };
  return client;
});

final localDbProvider = Provider<LocalDatabase>((ref) {
  return LocalDatabase();
});

final syncManagerProvider = Provider<SyncManager>((ref) {
  final api = ref.watch(apiClientProvider);
  final db = ref.watch(localDbProvider);
  return SyncManager(api: api, db: db);
});

// --- Auth State ---

enum AuthState { unknown, authenticated, unauthenticated }

final authStateProvider = StateProvider<AuthState>((ref) {
  return AuthState.unknown;
});

final authInitProvider = FutureProvider<AuthState>((ref) async {
  final storage = ref.watch(secureStorageProvider);
  final hasTokens = await storage.hasTokens();
  final state =
      hasTokens ? AuthState.authenticated : AuthState.unauthenticated;
  ref.read(authStateProvider.notifier).state = state;
  return state;
});

// --- Member Profile ---

final profileProvider =
    AsyncNotifierProvider<ProfileNotifier, Member?>(ProfileNotifier.new);

class ProfileNotifier extends AsyncNotifier<Member?> {
  @override
  Future<Member?> build() async {
    final authState = ref.watch(authStateProvider);
    if (authState != AuthState.authenticated) return null;

    final api = ref.read(apiClientProvider);
    try {
      final response = await api.getProfile();
      return Member.fromJson(response.data as Map<String, dynamic>);
    } catch (e) {
      return null;
    }
  }

  Future<void> refresh() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(() async {
      final api = ref.read(apiClientProvider);
      final response = await api.getProfile();
      return Member.fromJson(response.data as Map<String, dynamic>);
    });
  }

  Future<void> updateProfile(ProfileUpdate update) async {
    final api = ref.read(apiClientProvider);
    await api.updateProfile(update.toJson());
    await refresh();
  }

  Future<void> updatePrivacy(PrivacySettings settings) async {
    final api = ref.read(apiClientProvider);
    await api.updatePrivacy(settings.toJson());
    await refresh();
  }
}

// --- Attendance ---

final attendanceProvider =
    AsyncNotifierProvider<AttendanceNotifier, List<AttendanceRecord>>(
        AttendanceNotifier.new);

class AttendanceNotifier extends AsyncNotifier<List<AttendanceRecord>> {
  @override
  Future<List<AttendanceRecord>> build() async {
    final authState = ref.watch(authStateProvider);
    if (authState != AuthState.authenticated) return [];

    return _fetch();
  }

  Future<List<AttendanceRecord>> _fetch(
      {int limit = 50, int offset = 0}) async {
    final api = ref.read(apiClientProvider);
    try {
      final response =
          await api.getMyAttendance(limit: limit, offset: offset);
      final data = response.data['data'] as List<dynamic>? ?? [];
      return data
          .map((e) =>
              AttendanceRecord.fromJson(e as Map<String, dynamic>))
          .toList();
    } catch (e) {
      // Try local DB fallback
      final db = ref.read(localDbProvider);
      return db.getAttendanceRecords();
    }
  }

  Future<void> refresh() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(_fetch);
  }
}

// --- Streaks ---

final streaksProvider = FutureProvider<List<Streak>>((ref) async {
  final authState = ref.watch(authStateProvider);
  if (authState != AuthState.authenticated) return [];

  final api = ref.read(apiClientProvider);
  try {
    final response = await api.getMyStreaks();
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => Streak.fromJson(e as Map<String, dynamic>))
        .toList();
  } catch (e) {
    return [];
  }
});

// --- Member Packages ---

final memberPackagesProvider =
    FutureProvider<List<MemberPackage>>((ref) async {
  final authState = ref.watch(authStateProvider);
  if (authState != AuthState.authenticated) return [];

  final api = ref.read(apiClientProvider);
  try {
    final response = await api.getMyPackages();
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map(
            (e) => MemberPackage.fromJson(e as Map<String, dynamic>))
        .toList();
  } catch (e) {
    return [];
  }
});

/// Active package for the member (first active one).
final activePackageProvider = Provider<MemberPackage?>((ref) {
  final packages = ref.watch(memberPackagesProvider).valueOrNull ?? [];
  for (final pkg in packages) {
    if (pkg.isActive && !pkg.isExpired) return pkg;
  }
  return null;
});

// --- Public Packages ---

final publicPackagesProvider =
    FutureProvider<List<GymPackage>>((ref) async {
  final api = ref.read(apiClientProvider);
  try {
    final response = await api.getPublicPackages();
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => GymPackage.fromJson(e as Map<String, dynamic>))
        .toList();
  } catch (e) {
    return [];
  }
});

// --- Connectivity ---

final isOnlineProvider = StateProvider<bool>((ref) => true);

// --- Locale ---

final localeProvider = StateProvider<String>((ref) => 'en');
