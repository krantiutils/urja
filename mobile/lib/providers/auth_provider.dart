import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/user.dart';
import '../services/api_client.dart';
import '../services/auth_service.dart';

final apiClientProvider = Provider<ApiClient>((ref) => ApiClient());

final authServiceProvider = Provider<AuthService>(
    (ref) => AuthService(ref.read(apiClientProvider)));

enum AuthStatus { initial, unauthenticated, authenticated }

class AuthState {
  final AuthStatus status;
  final User? user;
  final String? error;
  final bool isNewUser;
  final bool onboardingCompleted;

  const AuthState({
    this.status = AuthStatus.initial,
    this.user,
    this.error,
    this.isNewUser = false,
    this.onboardingCompleted = false,
  });

  AuthState copyWith({
    AuthStatus? status,
    User? user,
    String? error,
    bool? isNewUser,
    bool? onboardingCompleted,
  }) =>
      AuthState(
        status: status ?? this.status,
        user: user ?? this.user,
        error: error,
        isNewUser: isNewUser ?? this.isNewUser,
        onboardingCompleted: onboardingCompleted ?? this.onboardingCompleted,
      );
}

class AuthNotifier extends StateNotifier<AuthState> {
  final AuthService _authService;

  AuthNotifier(this._authService) : super(const AuthState()) {
    _init();
  }

  Future<void> _init() async {
    final user = await _authService.getCurrentUser();
    if (user != null) {
      state = AuthState(
        status: AuthStatus.authenticated,
        user: user,
        onboardingCompleted: user.onboardingCompleted,
      );
    } else {
      state = const AuthState(status: AuthStatus.unauthenticated);
    }
  }

  Future<String> login(String phone) async {
    try {
      return await _authService.login(phone);
    } catch (e) {
      throw Exception('Failed to send OTP');
    }
  }

  Future<void> verifyOtp(String phone, String otp) async {
    try {
      final result = await _authService.verifyOtp(phone, otp);
      state = AuthState(
        status: AuthStatus.authenticated,
        user: result.user,
        isNewUser: result.isNewUser,
        onboardingCompleted: result.onboardingCompleted,
      );
    } catch (e) {
      state = state.copyWith(error: 'Invalid verification code');
      rethrow;
    }
  }

  void completeOnboarding() {
    if (state.user != null) {
      state = state.copyWith(
        onboardingCompleted: true,
        user: state.user!.copyWith(onboardingCompleted: true),
      );
    }
  }

  void updateUserType(String userType) {
    if (state.user != null) {
      state = state.copyWith(
        user: state.user!.copyWith(userType: userType),
      );
    }
  }

  Future<void> logout() async {
    await _authService.logout();
    state = const AuthState(status: AuthStatus.unauthenticated);
  }
}

final authProvider = StateNotifierProvider<AuthNotifier, AuthState>(
    (ref) => AuthNotifier(ref.read(authServiceProvider)));
