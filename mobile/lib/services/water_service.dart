import '../models/water.dart';
import 'api_client.dart';

class WaterService {
  final ApiClient _api;
  WaterService(this._api);

  Future<WaterLog> logWater({required int amountMl, String? date}) async {
    final data = <String, dynamic>{'amount_ml': amountMl};
    if (date != null) data['date'] = date;
    final response = await _api.post('/members/me/water-logs', data: data);
    return WaterLog.fromJson(response.data);
  }

  Future<List<WaterLog>> getWaterLogs(String date) async {
    final response = await _api.get('/members/me/water-logs',
        queryParameters: {'date': date});
    final list = response.data['data'] as List? ?? [];
    return list
        .map((e) => WaterLog.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<void> deleteWaterLog(String id) async {
    await _api.deleteRequest('/members/me/water-logs/$id');
  }

  Future<WaterDailySummary> getDailySummary(String date) async {
    final response = await _api.get('/members/me/water-logs/summary',
        queryParameters: {'date': date});
    return WaterDailySummary.fromJson(response.data);
  }
}
