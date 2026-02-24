import 'package:jwt_decoder/jwt_decoder.dart';
import '../models/user.dart';
import 'api_client.dart';
import 'token_storage.dart';

class AuthService {
  final ApiClient _api;
  final TokenStorage _storage = TokenStorage();

  AuthService(this._api);

  Future<String> login(String phone) async {
    final response = await _api.post('/auth/login', data: {'phone': phone});
    return response.data['message'] as String;
  }

  Future<AuthResult> verifyOtp(String phone, String otp) async {
    final response = await _api.post('/auth/verify-otp', data: {
      'phone': phone,
      'otp': otp,
    });

    final tokens = AuthTokens.fromJson(response.data);
    await _storage.write(key: 'access_token', value: tokens.accessToken);
    await _storage.write(key: 'refresh_token', value: tokens.refreshToken);

    final user = await _fetchUserProfile(tokens.accessToken);

    return AuthResult(
      user: user,
      isNewUser: tokens.isNewUser,
      onboardingCompleted: tokens.onboardingCompleted || user.onboardingCompleted,
    );
  }

  Future<User?> getCurrentUser() async {
    final token = await _storage.read(key: 'access_token');
    if (token == null) return null;

    if (JwtDecoder.isExpired(token)) {
      final refreshToken = await _storage.read(key: 'refresh_token');
      if (refreshToken == null) return null;

      try {
        final response = await _api.post('/auth/refresh', data: {
          'refresh_token': refreshToken,
        });
        final newTokens = AuthTokens.fromJson(response.data);
        await _storage.write(
            key: 'access_token', value: newTokens.accessToken);
        await _storage.write(
            key: 'refresh_token', value: newTokens.refreshToken);
        return _fetchUserProfile(newTokens.accessToken);
      } catch (_) {
        await logout();
        return null;
      }
    }

    // Augment JWT user with stored org_id, user_type, onboarding status
    final user = _userFromToken(token);
    final storedOrgId = await _storage.read(key: 'org_id');
    final storedUserType = await _storage.read(key: 'user_type');
    final storedOnboarding = await _storage.read(key: 'onboarding_completed');

    return User(
      id: user.id,
      phone: user.phone,
      role: user.role,
      orgId: user.orgId ?? storedOrgId,
      orgName: user.orgName,
      userType: user.userType ?? storedUserType,
      onboardingCompleted: user.onboardingCompleted ||
          (storedOnboarding == 'true'),
    );
  }

  Future<void> logout() async {
    final refreshToken = await _storage.read(key: 'refresh_token');
    if (refreshToken != null) {
      try {
        await _api.post('/auth/logout', data: {
          'refresh_token': refreshToken,
        });
      } catch (_) {
        // Ignore logout errors
      }
    }
    await _storage.deleteAll();
  }

  /// Build User from JWT claims, then fetch /members/me for org_id, user_type, onboarding.
  Future<User> _fetchUserProfile(String token) async {
    final payload = JwtDecoder.decode(token);
    final user = User(
      id: payload['sub'] as String,
      phone: payload['phone'] as String,
      role: payload['role'] as String,
      orgId: payload['org_id'] as String?,
      userType: payload['user_type'] as String?,
    );

    // Fetch additional profile data from /members/me
    try {
      final meResponse = await _api.get('/members/me');
      final data = meResponse.data;

      String? orgId = user.orgId;
      String? orgName;
      final orgs = data['organizations'] as List?;
      if (orgs != null && orgs.isNotEmpty) {
        orgId ??= orgs[0]['org_id'] as String?;
        orgName = orgs[0]['org_name'] as String?;
      }

      final userType = data['user_type'] as String? ?? user.userType;
      final onboardingCompleted =
          data['onboarding_completed'] as bool? ?? false;

      // Persist for offline use
      if (orgId != null) {
        await _storage.write(key: 'org_id', value: orgId);
      }
      if (userType != null) {
        await _storage.write(key: 'user_type', value: userType);
      }
      await _storage.write(
          key: 'onboarding_completed',
          value: onboardingCompleted.toString());

      return User(
        id: user.id,
        phone: user.phone,
        role: user.role,
        orgId: orgId,
        orgName: orgName,
        userType: userType,
        onboardingCompleted: onboardingCompleted,
      );
    } catch (_) {
      // If profile fetch fails, continue with JWT-only user
      return user;
    }
  }

  User _userFromToken(String token) {
    final payload = JwtDecoder.decode(token);
    return User(
      id: payload['sub'] as String,
      phone: payload['phone'] as String,
      role: payload['role'] as String,
      orgId: payload['org_id'] as String?,
      userType: payload['user_type'] as String?,
    );
  }
}
