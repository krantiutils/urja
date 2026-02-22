import '../models/staff.dart';
import 'api_client.dart';

class StaffService {
  final ApiClient _api;

  StaffService(this._api);

  Future<({List<StaffMember> data, int total})> getStaff(
    String orgId, {
    String? search,
    int limit = 20,
    int offset = 0,
  }) async {
    final response = await _api.get('/orgs/$orgId/staff', queryParameters: {
      'limit': limit,
      'offset': offset,
      if (search != null && search.isNotEmpty) 'search': search,
    });
    final list = (response.data['data'] as List<dynamic>? ?? [])
        .map((e) => StaffMember.fromJson(e as Map<String, dynamic>))
        .toList();
    return (data: list, total: response.data['total'] as int? ?? list.length);
  }

  Future<StaffMember> addStaff(
      String orgId, Map<String, dynamic> data) async {
    final response = await _api.post('/orgs/$orgId/staff', data: data);
    return StaffMember.fromJson(response.data);
  }

  Future<StaffMember> updateStaff(
      String orgId, String staffId, Map<String, dynamic> data) async {
    final response =
        await _api.put('/orgs/$orgId/staff/$staffId', data: data);
    return StaffMember.fromJson(response.data);
  }

  Future<void> deleteStaff(String orgId, String staffId) async {
    await _api.deleteRequest('/orgs/$orgId/staff/$staffId');
  }
}
