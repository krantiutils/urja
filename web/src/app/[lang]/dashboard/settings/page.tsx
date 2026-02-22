"use client";

import { useCallback, useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api, ApiRequestError } from "@/lib/api";
import type {
  Locale,
  Organization,
  UpdateOrganizationRequest,
  MemberProfile,
  ProfileUpdateRequest,
  PrivacySettingsUpdate,
} from "@/types";
import {
  Settings,
  Building2,
  Phone,
  Shield,
  Loader2,
  MapPin,
  Mail,
  Check,
  User,
  Heart,
  Eye,
  Calendar,
} from "lucide-react";

export default function SettingsPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  // --- Org state ---
  const [org, setOrg] = useState<Organization | null>(null);
  const [orgLoading, setOrgLoading] = useState(true);
  const [orgError, setOrgError] = useState<string | null>(null);

  // Org form
  const [orgName, setOrgName] = useState("");
  const [orgNameNe, setOrgNameNe] = useState("");
  const [description, setDescription] = useState("");
  const [descriptionNe, setDescriptionNe] = useState("");
  const [logoUrl, setLogoUrl] = useState("");
  const [orgPhone, setOrgPhone] = useState("");
  const [orgEmail, setOrgEmail] = useState("");
  const [address, setAddress] = useState("");
  const [addressNe, setAddressNe] = useState("");
  const [latitude, setLatitude] = useState("");
  const [longitude, setLongitude] = useState("");

  // Org save
  const [orgSaving, setOrgSaving] = useState(false);
  const [orgSaved, setOrgSaved] = useState(false);
  const [orgSaveError, setOrgSaveError] = useState<string | null>(null);

  // --- Profile state ---
  const [, setProfile] = useState<MemberProfile | null>(null);
  const [profileLoading, setProfileLoading] = useState(true);

  // Profile form
  const [profileName, setProfileName] = useState("");
  const [profileNameNe, setProfileNameNe] = useState("");
  const [profileEmail, setProfileEmail] = useState("");
  const [dateOfBirth, setDateOfBirth] = useState("");
  const [gender, setGender] = useState("");
  const [avatarUrl, setAvatarUrl] = useState("");
  const [emergencyName, setEmergencyName] = useState("");
  const [emergencyPhone, setEmergencyPhone] = useState("");

  // Privacy form
  const [showEmailPref, setShowEmailPref] = useState(false);
  const [showPhonePref, setShowPhonePref] = useState(false);
  const [showProfilePref, setShowProfilePref] = useState(false);
  const [showAttendancePref, setShowAttendancePref] = useState(false);
  const [showLeaderboardPref, setShowLeaderboardPref] = useState(false);

  // Profile save
  const [profileSaving, setProfileSaving] = useState(false);
  const [profileSaved, setProfileSaved] = useState(false);
  const [profileSaveError, setProfileSaveError] = useState<string | null>(null);

  // Privacy save
  const [privacySaving, setPrivacySaving] = useState(false);
  const [privacySaved, setPrivacySaved] = useState(false);

  const orgId = user?.org_id;
  const canEditOrg = user?.role === "admin" || user?.role === "staff";

  const inputClass =
    "w-full bg-surface border border-white/[0.06] rounded-xl px-4 py-2.5 text-fg text-sm focus:outline-none focus:ring-1 focus:ring-accent/50 transition-colors";
  const labelClass = "block text-xs text-fg-muted mb-1.5";
  const readOnlyValue = (v: string | null | undefined) => v || "\u2014";

  // --- Fetch org ---
  const populateOrg = useCallback((data: Organization) => {
    setOrgName(data.name ?? "");
    setOrgNameNe(data.name_ne ?? "");
    setDescription(data.description ?? "");
    setDescriptionNe(data.description_ne ?? "");
    setLogoUrl(data.logo_url ?? "");
    setOrgPhone(data.phone ?? "");
    setOrgEmail(data.email ?? "");
    setAddress(data.address ?? "");
    setAddressNe(data.address_ne ?? "");
    setLatitude(data.latitude != null ? String(data.latitude) : "");
    setLongitude(data.longitude != null ? String(data.longitude) : "");
  }, []);

  const fetchOrg = useCallback(async () => {
    if (!orgId) {
      setOrgLoading(false);
      return;
    }
    setOrgLoading(true);
    setOrgError(null);
    try {
      const data = await api.getOrg(orgId);
      setOrg(data);
      populateOrg(data);
    } catch (err) {
      setOrgError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setOrgLoading(false);
    }
  }, [orgId, t.common.error, populateOrg]);

  // --- Fetch profile ---
  const populateProfile = useCallback((data: MemberProfile) => {
    setProfileName(data.name ?? "");
    setProfileNameNe(data.name_ne ?? "");
    setProfileEmail(data.email ?? "");
    setDateOfBirth(data.date_of_birth ?? "");
    setGender(data.gender ?? "");
    setAvatarUrl(data.avatar_url ?? "");
    setEmergencyName(data.emergency_contact_name ?? "");
    setEmergencyPhone(data.emergency_contact_phone ?? "");
    if (data.privacy_settings) {
      setShowEmailPref(data.privacy_settings.show_email);
      setShowPhonePref(data.privacy_settings.show_phone);
      setShowProfilePref(data.privacy_settings.show_profile);
      setShowAttendancePref(data.privacy_settings.show_attendance);
      setShowLeaderboardPref(data.privacy_settings.show_on_leaderboard);
    }
  }, []);

  const fetchProfile = useCallback(async () => {
    setProfileLoading(true);
    try {
      const data = await api.getMyProfile();
      setProfile(data);
      populateProfile(data);
    } catch {
      // Profile fetch may fail for new users
    } finally {
      setProfileLoading(false);
    }
  }, [populateProfile]);

  useEffect(() => {
    fetchOrg();
    fetchProfile();
  }, [fetchOrg, fetchProfile]);

  // --- Save org ---
  const handleSaveOrg = async () => {
    if (!orgId) return;
    setOrgSaving(true);
    setOrgSaved(false);
    setOrgSaveError(null);

    const data: UpdateOrganizationRequest = {
      name: orgName,
      name_ne: orgNameNe,
      description,
      description_ne: descriptionNe,
      logo_url: logoUrl,
      phone: orgPhone,
      email: orgEmail,
      address,
      address_ne: addressNe,
      latitude: latitude ? parseFloat(latitude) : undefined,
      longitude: longitude ? parseFloat(longitude) : undefined,
    };

    try {
      const updated = await api.updateOrg(orgId, data);
      setOrg(updated);
      populateOrg(updated);
      setOrgSaved(true);
      setTimeout(() => setOrgSaved(false), 3000);
    } catch (err) {
      setOrgSaveError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setOrgSaving(false);
    }
  };

  // --- Save profile ---
  const handleSaveProfile = async () => {
    setProfileSaving(true);
    setProfileSaved(false);
    setProfileSaveError(null);

    const data: ProfileUpdateRequest = {
      name: profileName || undefined,
      name_ne: profileNameNe || undefined,
      email: profileEmail || undefined,
      date_of_birth: dateOfBirth || undefined,
      gender: gender || undefined,
      avatar_url: avatarUrl || undefined,
      emergency_contact_name: emergencyName || undefined,
      emergency_contact_phone: emergencyPhone || undefined,
    };

    try {
      await api.updateMyProfile(data);
      setProfileSaved(true);
      setTimeout(() => setProfileSaved(false), 3000);
    } catch (err) {
      setProfileSaveError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setProfileSaving(false);
    }
  };

  // --- Save privacy ---
  const handleSavePrivacy = async () => {
    setPrivacySaving(true);
    setPrivacySaved(false);

    const data: PrivacySettingsUpdate = {
      show_email: showEmailPref,
      show_phone: showPhonePref,
      show_profile: showProfilePref,
      show_attendance: showAttendancePref,
      show_on_leaderboard: showLeaderboardPref,
    };

    try {
      await api.updateMyPrivacy(data);
      setPrivacySaved(true);
      setTimeout(() => setPrivacySaved(false), 3000);
    } catch {
      // silent
    } finally {
      setPrivacySaving(false);
    }
  };

  const loading = orgLoading || profileLoading;

  const genderOptions = [
    { value: "", label: "\u2014" },
    { value: "male", label: t.settings.male },
    { value: "female", label: t.settings.female },
    { value: "other", label: t.settings.other },
    { value: "prefer_not_say", label: t.settings.preferNotSay },
  ];

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
          <Settings className="w-5 h-5 text-accent" />
          {t.settings.title}
        </h1>
      </div>

      {loading ? (
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card">
          <div className="flex items-center justify-center py-16">
            <Loader2 className="w-6 h-6 text-accent animate-spin" />
          </div>
        </div>
      ) : (
        <>
          {/* ===================== MY PROFILE ===================== */}

          {/* Your Account — read-only */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card overflow-hidden">
            <div className="px-5 py-4 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg flex items-center gap-2">
                <Shield className="w-4 h-4 text-fg-muted" />
                {t.settings.yourAccount}
              </h2>
            </div>
            <div className="p-5 space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-sm text-fg-muted flex items-center gap-2">
                  <Phone className="w-4 h-4 text-fg-muted" />
                  {t.settings.phone}
                </span>
                <span className="text-sm text-fg font-mono">
                  {user?.phone ?? "\u2014"}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-fg-muted flex items-center gap-2">
                  <Shield className="w-4 h-4 text-fg-muted" />
                  {t.settings.role}
                </span>
                <span className="text-sm text-fg font-mono capitalize">
                  {user?.role ?? "\u2014"}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-fg-muted flex items-center gap-2">
                  <Building2 className="w-4 h-4 text-fg-muted" />
                  {t.settings.orgId}
                </span>
                <span className="text-sm text-fg font-mono truncate max-w-[200px]">
                  {user?.org_id ?? "\u2014"}
                </span>
              </div>
            </div>
          </div>

          {/* Personal Info — editable by everyone */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card overflow-hidden">
            <div className="px-5 py-4 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg flex items-center gap-2">
                <User className="w-4 h-4 text-fg-muted" />
                {t.settings.myProfile}
              </h2>
            </div>
            <div className="p-5 space-y-4">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className={labelClass}>{t.settings.name}</label>
                  <input
                    type="text"
                    value={profileName}
                    onChange={(e) => setProfileName(e.target.value)}
                    className={inputClass}
                  />
                </div>
                <div>
                  <label className={labelClass}>{t.settings.nameNe}</label>
                  <input
                    type="text"
                    value={profileNameNe}
                    onChange={(e) => setProfileNameNe(e.target.value)}
                    className={inputClass}
                  />
                </div>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className={labelClass}>{t.settings.email}</label>
                  <input
                    type="email"
                    value={profileEmail}
                    onChange={(e) => setProfileEmail(e.target.value)}
                    className={inputClass}
                  />
                </div>
                <div>
                  <label className={labelClass}>{t.settings.gender}</label>
                  <select
                    value={gender}
                    onChange={(e) => setGender(e.target.value)}
                    className={inputClass}
                  >
                    {genderOptions.map((opt) => (
                      <option key={opt.value} value={opt.value}>
                        {opt.label}
                      </option>
                    ))}
                  </select>
                </div>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className={labelClass}>
                    <Calendar className="w-3 h-3 inline mr-1" />
                    {t.settings.dateOfBirth}
                  </label>
                  <input
                    type="date"
                    value={dateOfBirth}
                    onChange={(e) => setDateOfBirth(e.target.value)}
                    className={inputClass}
                  />
                </div>
                <div>
                  <label className={labelClass}>{t.settings.avatarUrl}</label>
                  <input
                    type="text"
                    value={avatarUrl}
                    onChange={(e) => setAvatarUrl(e.target.value)}
                    placeholder="https://..."
                    className={inputClass}
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Emergency Contact */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card overflow-hidden">
            <div className="px-5 py-4 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg flex items-center gap-2">
                <Heart className="w-4 h-4 text-fg-muted" />
                {t.settings.emergencyInfo}
              </h2>
            </div>
            <div className="p-5 space-y-4">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className={labelClass}>
                    {t.settings.emergencyContactName}
                  </label>
                  <input
                    type="text"
                    value={emergencyName}
                    onChange={(e) => setEmergencyName(e.target.value)}
                    className={inputClass}
                  />
                </div>
                <div>
                  <label className={labelClass}>
                    {t.settings.emergencyContactPhone}
                  </label>
                  <input
                    type="tel"
                    value={emergencyPhone}
                    onChange={(e) => setEmergencyPhone(e.target.value)}
                    placeholder="98XXXXXXXX"
                    className={inputClass}
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Save Profile Button */}
          <div className="flex items-center gap-3">
            <button
              onClick={handleSaveProfile}
              disabled={profileSaving}
              className="flex items-center gap-2 px-6 py-2.5 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors shadow-accent-glow disabled:opacity-50"
            >
              {profileSaving ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : profileSaved ? (
                <Check className="w-4 h-4" />
              ) : null}
              {profileSaving
                ? t.settings.saving
                : profileSaved
                  ? t.settings.profileSaved
                  : t.settings.saveChanges}
            </button>
            {profileSaveError && (
              <span className="text-sm text-red-400">{profileSaveError}</span>
            )}
          </div>

          {/* Privacy Settings */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card overflow-hidden">
            <div className="px-5 py-4 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg flex items-center gap-2">
                <Eye className="w-4 h-4 text-fg-muted" />
                {t.settings.privacySettings}
              </h2>
            </div>
            <div className="p-5 space-y-3">
              {[
                {
                  label: t.settings.showEmail,
                  value: showEmailPref,
                  setter: setShowEmailPref,
                },
                {
                  label: t.settings.showPhone,
                  value: showPhonePref,
                  setter: setShowPhonePref,
                },
                {
                  label: t.settings.showProfile,
                  value: showProfilePref,
                  setter: setShowProfilePref,
                },
                {
                  label: t.settings.showAttendance,
                  value: showAttendancePref,
                  setter: setShowAttendancePref,
                },
                {
                  label: t.settings.showOnLeaderboard,
                  value: showLeaderboardPref,
                  setter: setShowLeaderboardPref,
                },
              ].map((item) => (
                <label
                  key={item.label}
                  className="flex items-center justify-between cursor-pointer group"
                >
                  <span className="text-sm text-fg-muted group-hover:text-fg transition-colors">
                    {item.label}
                  </span>
                  <button
                    type="button"
                    role="switch"
                    aria-checked={item.value}
                    onClick={() => item.setter(!item.value)}
                    className={`relative w-10 h-5 rounded-full transition-colors ${
                      item.value ? "bg-accent" : "bg-surface border border-white/[0.1]"
                    }`}
                  >
                    <span
                      className={`absolute top-0.5 left-0.5 w-4 h-4 rounded-full bg-white transition-transform ${
                        item.value ? "translate-x-5" : ""
                      }`}
                    />
                  </button>
                </label>
              ))}
            </div>
          </div>

          {/* Save Privacy Button */}
          <div className="flex items-center gap-3">
            <button
              onClick={handleSavePrivacy}
              disabled={privacySaving}
              className="flex items-center gap-2 px-6 py-2.5 bg-surface border border-white/[0.06] text-fg font-medium text-sm rounded-xl hover:bg-surface-hover transition-colors disabled:opacity-50"
            >
              {privacySaving ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : privacySaved ? (
                <Check className="w-4 h-4" />
              ) : null}
              {privacySaving
                ? t.settings.saving
                : privacySaved
                  ? t.settings.privacySaved
                  : t.settings.saveChanges}
            </button>
          </div>

          {/* ===================== ORG SETTINGS ===================== */}

          {/* Error */}
          {orgError && (
            <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
              {orgError}
            </div>
          )}

          {/* General Information Card */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card overflow-hidden">
            <div className="px-5 py-4 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg flex items-center gap-2">
                <Building2 className="w-4 h-4 text-fg-muted" />
                {t.settings.generalInfo}
              </h2>
            </div>
            <div className="p-5 space-y-4">
              {canEditOrg ? (
                <>
                  <div>
                    <label className={labelClass}>{t.settings.orgName}</label>
                    <input
                      type="text"
                      value={orgName}
                      onChange={(e) => setOrgName(e.target.value)}
                      className={inputClass}
                    />
                  </div>
                  <div>
                    <label className={labelClass}>
                      {t.settings.orgNameNe}
                    </label>
                    <input
                      type="text"
                      value={orgNameNe}
                      onChange={(e) => setOrgNameNe(e.target.value)}
                      className={inputClass}
                    />
                  </div>
                  <div>
                    <label className={labelClass}>
                      {t.settings.description}
                    </label>
                    <textarea
                      value={description}
                      onChange={(e) => setDescription(e.target.value)}
                      rows={3}
                      className={inputClass + " resize-none"}
                    />
                  </div>
                  <div>
                    <label className={labelClass}>
                      {t.settings.descriptionNe}
                    </label>
                    <textarea
                      value={descriptionNe}
                      onChange={(e) => setDescriptionNe(e.target.value)}
                      rows={3}
                      className={inputClass + " resize-none"}
                    />
                  </div>
                  <div>
                    <label className={labelClass}>{t.settings.logo}</label>
                    <input
                      type="text"
                      value={logoUrl}
                      onChange={(e) => setLogoUrl(e.target.value)}
                      className={inputClass}
                    />
                  </div>
                </>
              ) : (
                <>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-fg-muted">
                      {t.settings.orgName}
                    </span>
                    <span className="text-sm text-fg">
                      {readOnlyValue(org?.name)}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-fg-muted">
                      {t.settings.orgNameNe}
                    </span>
                    <span className="text-sm text-fg">
                      {readOnlyValue(org?.name_ne)}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-fg-muted">
                      {t.settings.description}
                    </span>
                    <span className="text-sm text-fg text-right max-w-[60%]">
                      {readOnlyValue(org?.description)}
                    </span>
                  </div>
                </>
              )}
            </div>
          </div>

          {/* Contact Information Card */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card overflow-hidden">
            <div className="px-5 py-4 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg flex items-center gap-2">
                <Mail className="w-4 h-4 text-fg-muted" />
                {t.settings.contactInfo}
              </h2>
            </div>
            <div className="p-5 space-y-4">
              {canEditOrg ? (
                <>
                  <div>
                    <label className={labelClass}>{t.settings.phone}</label>
                    <input
                      type="text"
                      value={orgPhone}
                      onChange={(e) => setOrgPhone(e.target.value)}
                      className={inputClass}
                    />
                  </div>
                  <div>
                    <label className={labelClass}>{t.settings.email}</label>
                    <input
                      type="email"
                      value={orgEmail}
                      onChange={(e) => setOrgEmail(e.target.value)}
                      className={inputClass}
                    />
                  </div>
                </>
              ) : (
                <>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-fg-muted">
                      {t.settings.phone}
                    </span>
                    <span className="text-sm text-fg font-mono">
                      {readOnlyValue(org?.phone)}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-fg-muted">
                      {t.settings.email}
                    </span>
                    <span className="text-sm text-fg">
                      {readOnlyValue(org?.email)}
                    </span>
                  </div>
                </>
              )}
            </div>
          </div>

          {/* Location Card */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card overflow-hidden">
            <div className="px-5 py-4 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg flex items-center gap-2">
                <MapPin className="w-4 h-4 text-fg-muted" />
                {t.settings.location}
              </h2>
            </div>
            <div className="p-5 space-y-4">
              {canEditOrg ? (
                <>
                  <div>
                    <label className={labelClass}>{t.settings.address}</label>
                    <input
                      type="text"
                      value={address}
                      onChange={(e) => setAddress(e.target.value)}
                      className={inputClass}
                    />
                  </div>
                  <div>
                    <label className={labelClass}>
                      {t.settings.addressNe}
                    </label>
                    <input
                      type="text"
                      value={addressNe}
                      onChange={(e) => setAddressNe(e.target.value)}
                      className={inputClass}
                    />
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className={labelClass}>
                        {t.settings.latitude}
                      </label>
                      <input
                        type="number"
                        step="any"
                        value={latitude}
                        onChange={(e) => setLatitude(e.target.value)}
                        className={inputClass}
                      />
                    </div>
                    <div>
                      <label className={labelClass}>
                        {t.settings.longitude}
                      </label>
                      <input
                        type="number"
                        step="any"
                        value={longitude}
                        onChange={(e) => setLongitude(e.target.value)}
                        className={inputClass}
                      />
                    </div>
                  </div>
                </>
              ) : (
                <>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-fg-muted">
                      {t.settings.address}
                    </span>
                    <span className="text-sm text-fg text-right max-w-[60%]">
                      {readOnlyValue(org?.address)}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-fg-muted">
                      {t.settings.latitude} / {t.settings.longitude}
                    </span>
                    <span className="text-sm text-fg font-mono">
                      {org?.latitude != null && org?.longitude != null
                        ? `${org.latitude}, ${org.longitude}`
                        : "\u2014"}
                    </span>
                  </div>
                </>
              )}
            </div>
          </div>

          {/* Save Org Button — only for admin/staff */}
          {canEditOrg && (
            <div className="flex items-center gap-3">
              <button
                onClick={handleSaveOrg}
                disabled={orgSaving}
                className="flex items-center gap-2 px-6 py-2.5 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors shadow-accent-glow disabled:opacity-50"
              >
                {orgSaving ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : orgSaved ? (
                  <Check className="w-4 h-4" />
                ) : null}
                {orgSaving
                  ? t.settings.saving
                  : orgSaved
                    ? t.settings.saved
                    : t.settings.saveChanges}
              </button>
              {orgSaveError && (
                <span className="text-sm text-red-400">{orgSaveError}</span>
              )}
            </div>
          )}
        </>
      )}
    </div>
  );
}
