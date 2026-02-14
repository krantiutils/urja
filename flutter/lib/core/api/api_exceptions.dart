class ApiException implements Exception {
  final String message;
  final int? statusCode;

  const ApiException(this.message, {this.statusCode});

  @override
  String toString() => 'ApiException($statusCode): $message';

  bool get isUnauthorized => statusCode == 401;
  bool get isForbidden => statusCode == 403;
  bool get isNotFound => statusCode == 404;
  bool get isRateLimited => statusCode == 429;
  bool get isServerError => statusCode != null && statusCode! >= 500;
}

class NetworkException implements Exception {
  final String message;
  const NetworkException([this.message = 'Network error']);

  @override
  String toString() => 'NetworkException: $message';
}

class SessionExpiredException implements Exception {
  @override
  String toString() => 'SessionExpiredException: Session has expired';
}
