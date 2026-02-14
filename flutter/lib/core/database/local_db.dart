import 'package:path/path.dart';
import 'package:sqflite/sqflite.dart';

import '../../models/attendance.dart';

class LocalDatabase {
  Database? _db;

  Future<Database> get database async {
    _db ??= await _initDb();
    return _db!;
  }

  Future<Database> _initDb() async {
    final path = join(await getDatabasesPath(), 'urja.db');
    return openDatabase(
      path,
      version: 1,
      onCreate: (db, version) async {
        await db.execute('''
          CREATE TABLE attendance (
            id TEXT PRIMARY KEY,
            user_id TEXT NOT NULL,
            org_id TEXT NOT NULL,
            check_in_at TEXT NOT NULL,
            method TEXT NOT NULL,
            synced INTEGER DEFAULT 1
          )
        ''');
        await db.execute('''
          CREATE TABLE pending_checkins (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            method TEXT NOT NULL,
            qr_token TEXT,
            card_uid TEXT,
            org_id TEXT,
            created_at TEXT NOT NULL,
            synced INTEGER DEFAULT 0
          )
        ''');
        await db.execute('''
          CREATE TABLE streaks (
            id TEXT PRIMARY KEY,
            member_id TEXT NOT NULL,
            org_id TEXT NOT NULL,
            current_streak INTEGER DEFAULT 0,
            longest_streak INTEGER DEFAULT 0,
            last_check_in TEXT,
            updated_at TEXT NOT NULL
          )
        ''');
      },
    );
  }

  // --- Attendance ---

  Future<void> cacheAttendanceRecords(
      List<AttendanceRecord> records) async {
    final db = await database;
    final batch = db.batch();
    for (final r in records) {
      batch.insert(
        'attendance',
        {
          'id': r.id,
          'user_id': r.userId,
          'org_id': r.orgId,
          'check_in_at': r.checkInAt.toIso8601String(),
          'method': r.method,
          'synced': 1,
        },
        conflictAlgorithm: ConflictAlgorithm.replace,
      );
    }
    await batch.commit(noResult: true);
  }

  Future<List<AttendanceRecord>> getAttendanceRecords() async {
    final db = await database;
    final rows = await db.query('attendance',
        orderBy: 'check_in_at DESC', limit: 100);
    return rows.map((row) {
      return AttendanceRecord(
        id: row['id'] as String,
        userId: row['user_id'] as String,
        orgId: row['org_id'] as String,
        checkInAt: DateTime.parse(row['check_in_at'] as String),
        method: row['method'] as String,
      );
    }).toList();
  }

  // --- Pending Check-ins (offline queue) ---

  Future<int> queueCheckIn({
    required String method,
    String? qrToken,
    String? cardUid,
    String? orgId,
  }) async {
    final db = await database;
    return db.insert('pending_checkins', {
      'method': method,
      'qr_token': qrToken,
      'card_uid': cardUid,
      'org_id': orgId,
      'created_at': DateTime.now().toIso8601String(),
      'synced': 0,
    });
  }

  Future<List<Map<String, dynamic>>> getPendingCheckIns() async {
    final db = await database;
    return db.query('pending_checkins',
        where: 'synced = 0', orderBy: 'created_at ASC');
  }

  Future<void> markCheckInSynced(int id) async {
    final db = await database;
    await db.update('pending_checkins', {'synced': 1},
        where: 'id = ?', whereArgs: [id]);
  }

  Future<void> clearSyncedCheckIns() async {
    final db = await database;
    await db.delete('pending_checkins', where: 'synced = 1');
  }

  Future<void> close() async {
    final db = _db;
    if (db != null) {
      await db.close();
      _db = null;
    }
  }
}
