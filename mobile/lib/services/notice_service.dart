import '../models/feedback.dart';
import '../models/notice.dart';
import 'api_client.dart';

class NoticeService {
  final ApiClient _api;

  NoticeService(this._api);

  // Notices
  Future<({List<Notice> data, int total})> getNotices(
    String orgId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final response = await _api.get('/orgs/$orgId/notices', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    final list = (response.data['data'] as List<dynamic>? ?? [])
        .map((e) => Notice.fromJson(e as Map<String, dynamic>))
        .toList();
    return (data: list, total: response.data['total'] as int? ?? list.length);
  }

  Future<Notice> addNotice(String orgId, Map<String, dynamic> data) async {
    final response = await _api.post('/orgs/$orgId/notices', data: data);
    return Notice.fromJson(response.data);
  }

  Future<Notice> updateNotice(
      String orgId, String id, Map<String, dynamic> data) async {
    final response = await _api.put('/orgs/$orgId/notices/$id', data: data);
    return Notice.fromJson(response.data);
  }

  Future<void> deleteNotice(String orgId, String id) async {
    await _api.deleteRequest('/orgs/$orgId/notices/$id');
  }

  // Feedbacks
  Future<({List<MemberFeedback> data, int total})> getFeedbacks(
    String orgId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final response =
        await _api.get('/orgs/$orgId/feedbacks', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    final list = (response.data['data'] as List<dynamic>? ?? [])
        .map((e) => MemberFeedback.fromJson(e as Map<String, dynamic>))
        .toList();
    return (data: list, total: response.data['total'] as int? ?? list.length);
  }

  Future<void> deleteFeedback(String orgId, String id) async {
    await _api.deleteRequest('/orgs/$orgId/feedbacks/$id');
  }
}
