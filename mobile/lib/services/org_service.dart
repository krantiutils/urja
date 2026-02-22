import '../models/attendance.dart';
import '../models/member.dart';
import '../models/package.dart';
import 'api_client.dart';

class OrgService {
  final ApiClient _api;

  OrgService(this._api);

  Future<Organization> getOrg(String orgId) async {
    final response = await _api.get('/gyms/$orgId');
    return Organization.fromJson(response.data);
  }

  Future<Organization> updateOrg(
      String orgId, Map<String, dynamic> data) async {
    final response = await _api.put('/orgs/$orgId', data: data);
    return Organization.fromJson(response.data);
  }

  // Members
  Future<({List<OrgMember> data, int total})> getMembers(
    String orgId, {
    int limit = 20,
    int offset = 0,
    String? search,
  }) async {
    final response = await _api.get('/orgs/$orgId/members', queryParameters: {
      'limit': limit,
      'offset': offset,
      if (search != null && search.isNotEmpty) 'search': search,
    });
    final list = (response.data['data'] as List<dynamic>? ?? [])
        .map((e) => OrgMember.fromJson(e as Map<String, dynamic>))
        .toList();
    return (data: list, total: response.data['total'] as int? ?? list.length);
  }

  Future<OrgMember> addMember(String orgId, Map<String, dynamic> data) async {
    final response = await _api.post('/orgs/$orgId/members', data: data);
    return OrgMember.fromJson(response.data);
  }

  Future<void> updateMember(
      String orgId, String memberId, Map<String, dynamic> data) async {
    await _api.put('/orgs/$orgId/members/$memberId', data: data);
  }

  Future<void> removeMember(String orgId, String memberId) async {
    await _api.deleteRequest('/orgs/$orgId/members/$memberId');
  }

  // Attendance
  Future<List<OrgAttendance>> getAttendance(
    String orgId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final response =
        await _api.get('/orgs/$orgId/attendance', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => OrgAttendance.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<CheckInResponse> manualCheckIn(
      String orgId, String memberId) async {
    final response = await _api.post('/orgs/$orgId/attendance/check-in',
        data: {'member_id': memberId});
    return CheckInResponse.fromJson(response.data);
  }

  // Packages
  Future<List<OrgPackage>> getPackages(
    String orgId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final response =
        await _api.get('/orgs/$orgId/packages', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data
        .map((e) => OrgPackage.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<OrgPackage> addPackage(
      String orgId, Map<String, dynamic> data) async {
    final response = await _api.post('/orgs/$orgId/packages', data: data);
    return OrgPackage.fromJson(response.data);
  }

  Future<OrgPackage> updatePackage(
      String orgId, String packageId, Map<String, dynamic> data) async {
    final response =
        await _api.put('/orgs/$orgId/packages/$packageId', data: data);
    return OrgPackage.fromJson(response.data);
  }

  Future<void> deletePackage(String orgId, String packageId) async {
    await _api.deleteRequest('/orgs/$orgId/packages/$packageId');
  }

  // QR Code URL
  String getQrCodeUrl(String orgId) =>
      '${_api.dio.options.baseUrl}/orgs/$orgId/qr-code';
}
