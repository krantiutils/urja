import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../services/water_service.dart';
import 'auth_provider.dart';

final waterServiceProvider = Provider<WaterService>(
    (ref) => WaterService(ref.read(apiClientProvider)));
