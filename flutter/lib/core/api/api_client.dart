import 'dart:io';

import 'package:dio/dio.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

import '../auth/secure_storage.dart';
import 'api_exceptions.dart';

class ApiClient {
  late final Dio _dio;
  final SecureTokenStorage _tokenStorage;
  bool _isRefreshing = false;
  final List<void Function(String)> _refreshCallbacks = [];

  /// Called when tokens are fully expired and user must re-authenticate.
  void Function()? onSessionExpired;

  ApiClient({required SecureTokenStorage tokenStorage})
      : _tokenStorage = tokenStorage {
    final baseUrl = dotenv.get('API_BASE_URL',
        fallback: 'http://10.0.2.2:8080/api/v1');

    _dio = Dio(BaseOptions(
      baseUrl: baseUrl,
      connectTimeout: const Duration(seconds: 15),
      receiveTimeout: const Duration(seconds: 15),
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    ));

    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: _onRequest,
      onError: _onError,
    ));
  }

  Dio get dio => _dio;

  Future<void> _onRequest(
      RequestOptions options, RequestInterceptorHandler handler) async {
    // Skip auth header for login/verify/refresh
    final noAuthPaths = ['/auth/login', '/auth/verify-otp', '/auth/refresh'];
    final needsAuth = !noAuthPaths.any((p) => options.path.contains(p));

    if (needsAuth) {
      final token = await _tokenStorage.getAccessToken();
      if (token != null) {
        options.headers['Authorization'] = 'Bearer $token';
      }
    }
    handler.next(options);
  }

  Future<void> _onError(
      DioException error, ErrorInterceptorHandler handler) async {
    if (error.response?.statusCode == 401 &&
        !error.requestOptions.path.contains('/auth/')) {
      // Attempt token refresh
      final refreshed = await _tryRefreshToken(error.requestOptions);
      if (refreshed != null) {
        handler.resolve(refreshed);
        return;
      }
      // Refresh failed — session expired
      onSessionExpired?.call();
      handler.reject(DioException(
        requestOptions: error.requestOptions,
        error: SessionExpiredException(),
      ));
      return;
    }
    handler.next(error);
  }

  Future<Response?> _tryRefreshToken(RequestOptions failedRequest) async {
    if (_isRefreshing) {
      // Wait for the ongoing refresh to complete
      try {
        final newToken = await _waitForRefresh();
        failedRequest.headers['Authorization'] = 'Bearer $newToken';
        return await _dio.fetch(failedRequest);
      } catch (_) {
        return null;
      }
    }

    _isRefreshing = true;
    try {
      final refreshToken = await _tokenStorage.getRefreshToken();
      if (refreshToken == null) return null;

      final baseUrl = _dio.options.baseUrl;
      final refreshDio = Dio(BaseOptions(
        baseUrl: baseUrl,
        headers: {'Content-Type': 'application/json'},
      ));

      final response = await refreshDio.post(
        '/auth/refresh',
        data: {'refresh_token': refreshToken},
      );

      if (response.statusCode == 200) {
        final newAccessToken = response.data['access_token'] as String;
        final newRefreshToken = response.data['refresh_token'] as String;
        await _tokenStorage.saveTokens(
          accessToken: newAccessToken,
          refreshToken: newRefreshToken,
        );

        // Notify waiting requests
        for (final cb in _refreshCallbacks) {
          cb(newAccessToken);
        }
        _refreshCallbacks.clear();

        // Retry the original request
        failedRequest.headers['Authorization'] = 'Bearer $newAccessToken';
        return await _dio.fetch(failedRequest);
      }
      return null;
    } on DioException {
      return null;
    } finally {
      _isRefreshing = false;
    }
  }

  Future<String> _waitForRefresh() {
    return Future.delayed(const Duration(milliseconds: 500), () async {
      final token = await _tokenStorage.getAccessToken();
      return token ?? (throw StateError('No token after refresh'));
    });
  }

  // --- Auth endpoints ---

  Future<void> requestOtp(String phone) async {
    final response = await _request('POST', '/auth/login', data: {'phone': phone});
    // No return value needed — OTP sent
    if (response.statusCode != 200) {
      throw ApiException(
        response.data['error'] as String? ?? 'Failed to send OTP',
        statusCode: response.statusCode,
      );
    }
  }

  Future<void> verifyOtp(String phone, String otp) async {
    final response = await _request(
      'POST',
      '/auth/verify-otp',
      data: {'phone': phone, 'otp': otp},
    );

    if (response.statusCode != 200) {
      throw ApiException(
        response.data['error'] as String? ?? 'OTP verification failed',
        statusCode: response.statusCode,
      );
    }

    await _tokenStorage.saveTokens(
      accessToken: response.data['access_token'] as String,
      refreshToken: response.data['refresh_token'] as String,
    );
  }

  Future<void> logout() async {
    final refreshToken = await _tokenStorage.getRefreshToken();
    if (refreshToken != null) {
      try {
        await _request('POST', '/auth/logout',
            data: {'refresh_token': refreshToken});
      } catch (_) {
        // Best-effort: clear tokens regardless
      }
    }
    await _tokenStorage.clearTokens();
  }

  Future<void> logoutAll() async {
    try {
      await _request('POST', '/auth/logout-all');
    } catch (_) {
      // Best-effort
    }
    await _tokenStorage.clearTokens();
  }

  // --- Generic request helper ---

  Future<Response> _request(
    String method,
    String path, {
    Map<String, dynamic>? data,
    Map<String, dynamic>? queryParams,
  }) async {
    try {
      final response = await _dio.request(
        path,
        data: data,
        queryParameters: queryParams,
        options: Options(method: method),
      );
      return response;
    } on DioException catch (e) {
      if (e.error is SessionExpiredException) {
        rethrow;
      }
      if (e.type == DioExceptionType.connectionTimeout ||
          e.type == DioExceptionType.receiveTimeout ||
          e.type == DioExceptionType.sendTimeout ||
          e.error is SocketException) {
        throw const NetworkException('Connection failed. Check your internet.');
      }
      if (e.response != null) {
        final body = e.response!.data;
        final message = body is Map ? (body['error'] as String?) : null;
        throw ApiException(
          message ?? 'Request failed',
          statusCode: e.response!.statusCode,
        );
      }
      throw ApiException('Unexpected error: ${e.message}');
    }
  }

  // --- Member endpoints ---

  Future<Response> getProfile() => _request('GET', '/members/me');

  Future<Response> updateProfile(Map<String, dynamic> data) =>
      _request('PUT', '/members/me', data: data);

  Future<Response> updatePrivacy(Map<String, dynamic> data) =>
      _request('PUT', '/members/me/privacy', data: data);

  // --- Attendance endpoints ---

  Future<Response> selfCheckIn(Map<String, dynamic> data) =>
      _request('POST', '/members/me/check-in', data: data);

  Future<Response> getMyAttendance({int limit = 20, int offset = 0}) =>
      _request('GET', '/members/me/attendance',
          queryParams: {'limit': limit, 'offset': offset});

  Future<Response> getMyStreaks() =>
      _request('GET', '/members/me/streaks');

  // --- Package endpoints ---

  Future<Response> getMyPackages({int limit = 20, int offset = 0}) =>
      _request('GET', '/members/me/packages',
          queryParams: {'limit': limit, 'offset': offset});

  Future<Response> subscribeToPackage(String packageId) =>
      _request('POST', '/members/me/packages/$packageId/subscribe');

  Future<Response> verifyPayment(String pidx) =>
      _request('POST', '/members/me/packages/verify-payment',
          data: {'pidx': pidx});

  Future<Response> getPublicPackages({int limit = 20, int offset = 0}) =>
      _request('GET', '/packages',
          queryParams: {'limit': limit, 'offset': offset});

  // --- Org endpoints ---

  Future<Response> getLeaderboard(String orgId,
          {String period = 'weekly'}) =>
      _request('GET', '/orgs/$orgId/leaderboard',
          queryParams: {'period': period});
}
