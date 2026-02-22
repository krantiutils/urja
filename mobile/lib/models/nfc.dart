class NfcCard {
  final String id;
  final String orgId;
  final String cardNumber;
  final String? assignedMember;
  final String? fullName;
  final String status;
  final String? assignedAt;
  final String lastUpdated;

  NfcCard({
    required this.id,
    required this.orgId,
    required this.cardNumber,
    this.assignedMember,
    this.fullName,
    required this.status,
    this.assignedAt,
    required this.lastUpdated,
  });

  factory NfcCard.fromJson(Map<String, dynamic> json) => NfcCard(
        id: json['id'] as String,
        orgId: json['org_id'] as String,
        cardNumber: json['card_number'] as String,
        assignedMember: json['assigned_member'] as String?,
        fullName: json['full_name'] as String?,
        status: json['status'] as String,
        assignedAt: json['assigned_at'] as String?,
        lastUpdated: json['last_updated'] as String,
      );
}

class NfcDevice {
  final String id;
  final String orgId;
  final String name;
  final String deviceIdentifier;
  final String status;
  final String doorState;
  final String? lastSyncAt;
  final int uptimeSeconds;
  final String createdAt;
  final String updatedAt;

  NfcDevice({
    required this.id,
    required this.orgId,
    required this.name,
    required this.deviceIdentifier,
    required this.status,
    required this.doorState,
    this.lastSyncAt,
    required this.uptimeSeconds,
    required this.createdAt,
    required this.updatedAt,
  });

  factory NfcDevice.fromJson(Map<String, dynamic> json) => NfcDevice(
        id: json['id'] as String,
        orgId: json['org_id'] as String,
        name: json['name'] as String,
        deviceIdentifier: json['device_identifier'] as String,
        status: json['status'] as String? ?? 'offline',
        doorState: json['door_state'] as String? ?? 'locked',
        lastSyncAt: json['last_sync_at'] as String?,
        uptimeSeconds: json['uptime_seconds'] as int? ?? 0,
        createdAt: json['created_at'] as String,
        updatedAt: json['updated_at'] as String,
      );
}
