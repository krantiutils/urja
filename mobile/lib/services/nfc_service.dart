import '../models/nfc.dart';
import 'api_client.dart';

class NfcApiService {
  final ApiClient _api;

  NfcApiService(this._api);

  // Cards
  Future<({List<NfcCard> data, int total})> getCards(
    String orgId, {
    String? status,
    int limit = 20,
    int offset = 0,
  }) async {
    final response =
        await _api.get('/orgs/$orgId/nfc-cards', queryParameters: {
      'limit': limit,
      'offset': offset,
      if (status != null && status.isNotEmpty) 'status': status,
    });
    final list = (response.data['data'] as List<dynamic>? ?? [])
        .map((e) => NfcCard.fromJson(e as Map<String, dynamic>))
        .toList();
    return (data: list, total: response.data['total'] as int? ?? list.length);
  }

  Future<NfcCard> registerCard(String orgId, String cardNumber) async {
    final response = await _api.post('/orgs/$orgId/nfc-cards', data: {
      'card_number': cardNumber,
    });
    return NfcCard.fromJson(response.data);
  }

  Future<NfcCard> assignCard(
      String orgId, String cardId, String userId) async {
    final response = await _api.put('/orgs/$orgId/nfc-cards/$cardId/assign',
        data: {'user_id': userId});
    return NfcCard.fromJson(response.data);
  }

  Future<NfcCard> unassignCard(String orgId, String cardId) async {
    final response =
        await _api.put('/orgs/$orgId/nfc-cards/$cardId/unassign');
    return NfcCard.fromJson(response.data);
  }

  // Devices
  Future<List<NfcDevice>> getDevices(String orgId) async {
    final response = await _api.get('/orgs/$orgId/nfc-devices');
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => NfcDevice.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<NfcDevice> addDevice(String orgId, Map<String, dynamic> data) async {
    final response = await _api.post('/orgs/$orgId/nfc-devices', data: data);
    return NfcDevice.fromJson(response.data);
  }
}
