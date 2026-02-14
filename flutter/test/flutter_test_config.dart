import 'dart:async';

import 'helpers/test_helpers.dart';

Future<void> testExecutable(FutureOr<void> Function() testMain) async {
  await loadInterFont();
  await testMain();
}
