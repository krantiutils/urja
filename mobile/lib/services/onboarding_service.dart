import 'api_client.dart';

class OnboardingService {
  final ApiClient _api;

  OnboardingService(this._api);

  Future<Map<String, dynamic>> submitOnboarding(
      Map<String, dynamic> data) async {
    final response = await _api.post('/members/me/onboarding', data: data);
    return response.data as Map<String, dynamic>;
  }

  Future<void> joinGym(String orgId) async {
    await _api.post('/members/me/join-gym', data: {'org_id': orgId});
  }

  Future<void> leaveGym(String orgId) async {
    await _api.post('/members/me/leave-gym', data: {'org_id': orgId});
  }
}
