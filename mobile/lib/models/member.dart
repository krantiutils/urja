import 'user.dart';

class OrgMembership {
  final String orgId;
  final String orgName;
  final String role;
  final String status;
  final String joinedAt;

  OrgMembership({
    required this.orgId,
    required this.orgName,
    required this.role,
    required this.status,
    required this.joinedAt,
  });

  factory OrgMembership.fromJson(Map<String, dynamic> json) => OrgMembership(
        orgId: json['org_id'] as String,
        orgName: json['org_name'] as String,
        role: json['role'] as String,
        status: json['status'] as String,
        joinedAt: json['joined_at'] as String,
      );
}

class MemberProfile {
  final String id;
  final String phone;
  final String name;
  final String? nameNe;
  final String? email;
  final String? dateOfBirth;
  final String? gender;
  final String? avatarUrl;
  final String? emergencyContactName;
  final String? emergencyContactPhone;
  final PrivacySettings privacySettings;
  final List<OrgMembership> organizations;
  final String createdAt;
  final String updatedAt;

  MemberProfile({
    required this.id,
    required this.phone,
    required this.name,
    this.nameNe,
    this.email,
    this.dateOfBirth,
    this.gender,
    this.avatarUrl,
    this.emergencyContactName,
    this.emergencyContactPhone,
    required this.privacySettings,
    required this.organizations,
    required this.createdAt,
    required this.updatedAt,
  });

  factory MemberProfile.fromJson(Map<String, dynamic> json) => MemberProfile(
        id: json['id'] as String,
        phone: json['phone'] as String,
        name: json['name'] as String,
        nameNe: json['name_ne'] as String?,
        email: json['email'] as String?,
        dateOfBirth: json['date_of_birth'] as String?,
        gender: json['gender'] as String?,
        avatarUrl: json['avatar_url'] as String?,
        emergencyContactName: json['emergency_contact_name'] as String?,
        emergencyContactPhone: json['emergency_contact_phone'] as String?,
        privacySettings: PrivacySettings.fromJson(
            json['privacy_settings'] as Map<String, dynamic>? ?? {}),
        organizations: (json['organizations'] as List<dynamic>?)
                ?.map((e) =>
                    OrgMembership.fromJson(e as Map<String, dynamic>))
                .toList() ??
            [],
        createdAt: json['created_at'] as String,
        updatedAt: json['updated_at'] as String,
      );
}

class OrgMember {
  final String id;
  final String phone;
  final String name;
  final String? nameNe;
  final String? email;
  final String? avatarUrl;
  final String role;
  final String status;
  final String joinedAt;

  OrgMember({
    required this.id,
    required this.phone,
    required this.name,
    this.nameNe,
    this.email,
    this.avatarUrl,
    required this.role,
    required this.status,
    required this.joinedAt,
  });

  factory OrgMember.fromJson(Map<String, dynamic> json) => OrgMember(
        id: json['id'] as String,
        phone: json['phone'] as String,
        name: json['name'] as String,
        nameNe: json['name_ne'] as String?,
        email: json['email'] as String?,
        avatarUrl: json['avatar_url'] as String?,
        role: json['role'] as String,
        status: json['status'] as String,
        joinedAt: json['joined_at'] as String,
      );
}

class Organization {
  final String id;
  final String name;
  final String nameNe;
  final String slug;
  final String description;
  final String descriptionNe;
  final String logoUrl;
  final String address;
  final String addressNe;
  final String phone;
  final String email;
  final double? latitude;
  final double? longitude;
  final String timezone;
  final bool isActive;
  final String createdAt;
  final String updatedAt;

  Organization({
    required this.id,
    required this.name,
    required this.nameNe,
    required this.slug,
    required this.description,
    required this.descriptionNe,
    required this.logoUrl,
    required this.address,
    required this.addressNe,
    required this.phone,
    required this.email,
    this.latitude,
    this.longitude,
    required this.timezone,
    required this.isActive,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Organization.fromJson(Map<String, dynamic> json) => Organization(
        id: json['id'] as String,
        name: json['name'] as String? ?? '',
        nameNe: json['name_ne'] as String? ?? '',
        slug: json['slug'] as String? ?? '',
        description: json['description'] as String? ?? '',
        descriptionNe: json['description_ne'] as String? ?? '',
        logoUrl: json['logo_url'] as String? ?? '',
        address: json['address'] as String? ?? '',
        addressNe: json['address_ne'] as String? ?? '',
        phone: json['phone'] as String? ?? '',
        email: json['email'] as String? ?? '',
        latitude: (json['latitude'] as num?)?.toDouble(),
        longitude: (json['longitude'] as num?)?.toDouble(),
        timezone: json['timezone'] as String? ?? 'Asia/Kathmandu',
        isActive: json['is_active'] as bool? ?? true,
        createdAt: json['created_at'] as String? ?? '',
        updatedAt: json['updated_at'] as String? ?? '',
      );

  Map<String, dynamic> toJson() => {
        'name': name,
        'name_ne': nameNe,
        'description': description,
        'description_ne': descriptionNe,
        'logo_url': logoUrl,
        'address': address,
        'address_ne': addressNe,
        'phone': phone,
        'email': email,
        'latitude': latitude,
        'longitude': longitude,
      };
}
