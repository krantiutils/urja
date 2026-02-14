import 'dart:async';

import 'package:connectivity_plus/connectivity_plus.dart';

import '../api/api_client.dart';
import '../api/api_exceptions.dart';
import '../database/local_db.dart';

class SyncManager {
  final ApiClient api;
  final LocalDatabase db;
  StreamSubscription<List<ConnectivityResult>>? _connectivitySub;
  bool _syncing = false;

  SyncManager({required this.api, required this.db});

  void startListening() {
    _connectivitySub =
        Connectivity().onConnectivityChanged.listen((results) {
      final hasConnection = results.any((r) =>
          r == ConnectivityResult.wifi ||
          r == ConnectivityResult.mobile ||
          r == ConnectivityResult.ethernet);
      if (hasConnection) {
        syncPendingCheckIns();
      }
    });
  }

  void stopListening() {
    _connectivitySub?.cancel();
  }

  Future<void> syncPendingCheckIns() async {
    if (_syncing) return;
    _syncing = true;

    try {
      final pending = await db.getPendingCheckIns();
      for (final item in pending) {
        try {
          await api.selfCheckIn({
            'method': item['method'] as String,
            if (item['qr_token'] != null) 'qr_token': item['qr_token'],
            if (item['card_uid'] != null) 'card_uid': item['card_uid'],
            if (item['org_id'] != null) 'org_id': item['org_id'],
          });
          await db.markCheckInSynced(item['id'] as int);
        } on ApiException {
          // If the server rejects it (e.g. expired QR), mark it synced
          // to prevent retry loops.
          await db.markCheckInSynced(item['id'] as int);
        } on NetworkException {
          // Still offline, stop trying
          break;
        }
      }
      await db.clearSyncedCheckIns();
    } finally {
      _syncing = false;
    }
  }
}
