import '../models/sms.dart';
import 'api_client.dart';

class SmsService {
  final ApiClient _api;

  SmsService(this._api);

  Future<SmsBalance> getBalance(String orgId) async {
    final response = await _api.get('/orgs/$orgId/sms/balance');
    return SmsBalance.fromJson(response.data);
  }

  Future<void> send(
      String orgId, String message, List<String> memberIds) async {
    await _api.post('/orgs/$orgId/sms/send', data: {
      'message': message,
      'member_ids': memberIds,
    });
  }

  Future<List<dynamic>> getHistory(
    String orgId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final response =
        await _api.get('/orgs/$orgId/sms/history', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    return response.data['data'] as List<dynamic>? ?? [];
  }
}
