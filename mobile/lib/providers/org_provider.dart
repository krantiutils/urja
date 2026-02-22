import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/member.dart';
import '../providers/auth_provider.dart';
import '../services/org_service.dart';

final orgServiceProvider = Provider<OrgService>(
    (ref) => OrgService(ref.read(apiClientProvider)));

final currentOrgProvider = FutureProvider<Organization?>((ref) async {
  final auth = ref.watch(authProvider);
  if (auth.user?.orgId == null) return null;
  final orgService = ref.read(orgServiceProvider);
  return orgService.getOrg(auth.user!.orgId!);
});
