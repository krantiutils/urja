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
  Lock,
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
  Receipt,
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

  // Tax / PAN form. Without a PAN no bill can be issued, so this is kept
  // separate from the rest of the org form: a gym can fix just this without
  // touching its name, address, etc.
  const [pan, setPan] = useState("");
  const [taxLegalName, setTaxLegalName] = useState("");
  const [taxAddress, setTaxAddress] = useState("");
  const [panError, setPanError] = useState<string | null>(null);
  const [taxSaving, setTaxSaving] = useState(false);
  const [taxSaved, setTaxSaved] = useState(false);
  const [taxSaveError, setTaxSaveError] = useState<string | null>(null);

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

  // Password. Login by password shipped without any way to set one, so the
  // only route in was an API call.
  const [passwordSet, setPasswordSet] = useState<boolean | null>(null);
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [passwordSaving, setPasswordSaving] = useState(false);
  const [passwordSaved, setPasswordSaved] = useState(false);
  const [passwordError, setPasswordError] = useState<string | null>(null);

  // Privacy save
  const [privacySaving, setPrivacySaving] = useState(false);
  const [privacySaved, setPrivacySaved] = useState(false);

  const orgId = user?.org_id;
  // The JWT `role` claim is vestigial — it is "member" on every token, including
  // an org admin's, so gating on it made this whole page read-only for everybody.
  // org_role is the real one, resolved from actual memberships. Admin, not staff:
  // the API's org update rejects staff with 403, so offering them an editable
  // form would only produce a failed save.
  const canEditOrg = user?.org_role === "admin" || Boolean(user?.is_super_admin);

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
    setPan(data.pan_number ?? "");
    setTaxLegalName(data.tax_legal_name ?? "");
    setTaxAddress(data.tax_address ?? "");
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

  // --- Save tax details ---
  const handleSaveTax = async (e: React.FormEvent) => {
    e.preventDefault();
    setPanError(null);
    setTaxSaveError(null);

    // Match the server's own rule exactly: 9 digits, or empty. An empty PAN
    // is how a gym clears a wrong one — it must not be rejected here.
    if (pan && !/^\d{9}$/.test(pan)) {
      setPanError(t.invoices.panDigits);
      return;
    }

    if (!orgId) return;
    setTaxSaving(true);
    setTaxSaved(false);
    try {
      const updated = await api.updateOrg(orgId, {
        pan_number: pan,
        tax_legal_name: taxLegalName,
        tax_address: taxAddress,
      });
      setOrg(updated);
      populateOrg(updated);
      setTaxSaved(true);
      setTimeout(() => setTaxSaved(false), 3000);
    } catch (err) {
      setTaxSaveError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setTaxSaving(false);
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
  useEffect(() => {
    api
      .passwordStatus()
      .then((r) => setPasswordSet(r.password_set))
      .catch(() => setPasswordSet(null));
  }, []);

  const handleSavePassword = async () => {
    setPasswordError(null);
    if (newPassword.length < 8) {
      setPasswordError(t.settings.passwordTooShort);
      return;
    }
    if (newPassword !== confirmPassword) {
      setPasswordError(t.settings.passwordMismatch);
      return;
    }

    setPasswordSaving(true);
    try {
      await api.setPassword(newPassword);
      setPasswordSet(true);
      setNewPassword("");
      setConfirmPassword("");
      setPasswordSaved(true);
      setTimeout(() => setPasswordSaved(false), 3000);
    } catch (err) {
      setPasswordError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setPasswordSaving(false);
    }
  };

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
                  {user?.org_role ?? "\u2014"}
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

          {/* Password */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card overflow-hidden">
            <div className="px-5 py-4 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg flex items-center gap-2">
                <Lock className="w-4 h-4 text-fg-muted" />
                {t.settings.password}
              </h2>
            </div>
            <div className="p-5">
              <p className="text-sm text-fg-muted mb-1">{t.settings.passwordIntro}</p>
              {passwordSet !== null && (
                <p className="text-xs text-fg-muted mb-4">
                  {passwordSet ? t.settings.passwordIsSet : t.settings.passwordNotSet}
                </p>
              )}

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div>
                  <label htmlFor="new-password" className="block text-xs text-fg-muted mb-1.5">
                    {t.settings.newPassword}
                  </label>
                  <input
                    id="new-password"
                    type="password"
                    autoComplete="new-password"
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg focus:outline-none focus:border-accent/50"
                  />
                </div>
                <div>
                  <label htmlFor="confirm-password" className="block text-xs text-fg-muted mb-1.5">
                    {t.settings.confirmPassword}
                  </label>
                  <input
                    id="confirm-password"
                    type="password"
                    autoComplete="new-password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg focus:outline-none focus:border-accent/50"
                  />
                </div>
              </div>

              {passwordError && (
                <p className="mt-3 text-sm text-red-400" role="alert">
                  {passwordError}
                </p>
              )}

              <div className="mt-4 flex items-center gap-3">
                <button
                  type="button"
                  onClick={handleSavePassword}
                  disabled={passwordSaving || !newPassword}
                  className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-accent text-bg-base text-sm font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
                >
                  {passwordSaving && <Loader2 className="w-4 h-4 animate-spin" />}
                  {t.settings.savePassword}
                </button>
                {passwordSaved && (
                  <span className="text-sm text-accent">{t.settings.passwordSaved}</span>
                )}
              </div>
            </div>
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

          {/* Tax Details Card — no bill can be issued without a PAN here */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card overflow-hidden">
            <div className="px-5 py-4 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg flex items-center gap-2">
                <Receipt className="w-4 h-4 text-fg-muted" />
                {t.settings.taxDetails}
              </h2>
            </div>
            <div className="p-5 space-y-4">
              {canEditOrg ? (
                <form onSubmit={handleSaveTax} className="space-y-4">
                  <p className="text-sm text-fg-muted">{t.settings.taxDetailsIntro}</p>
                  <div>
                    <label htmlFor="tax-pan" className={labelClass}>
                      {t.settings.panNumber}
                    </label>
                    <input
                      id="tax-pan"
                      type="text"
                      inputMode="numeric"
                      value={pan}
                      onChange={(e) => setPan(e.target.value)}
                      placeholder="123456789"
                      className={inputClass}
                    />
                  </div>
                  <div>
                    <label htmlFor="tax-legal-name" className={labelClass}>
                      {t.settings.taxLegalName}
                    </label>
                    <input
                      id="tax-legal-name"
                      type="text"
                      value={taxLegalName}
                      onChange={(e) => setTaxLegalName(e.target.value)}
                      maxLength={255}
                      className={inputClass}
                    />
                  </div>
                  <div>
                    <label htmlFor="tax-address" className={labelClass}>
                      {t.settings.taxAddress}
                    </label>
                    <input
                      id="tax-address"
                      type="text"
                      value={taxAddress}
                      onChange={(e) => setTaxAddress(e.target.value)}
                      maxLength={500}
                      className={inputClass}
                    />
                  </div>

                  {panError && (
                    <p className="text-sm text-red-400" role="alert">
                      {panError}
                    </p>
                  )}

                  <div className="flex items-center gap-3">
                    <button
                      type="submit"
                      disabled={taxSaving}
                      className="flex items-center gap-2 px-6 py-2.5 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors shadow-accent-glow disabled:opacity-50"
                    >
                      {taxSaving ? (
                        <Loader2 className="w-4 h-4 animate-spin" />
                      ) : taxSaved ? (
                        <Check className="w-4 h-4" />
                      ) : null}
                      {taxSaving
                        ? t.settings.saving
                        : taxSaved
                          ? t.settings.saved
                          : t.settings.saveChanges}
                    </button>
                    {taxSaveError && (
                      <span className="text-sm text-red-400">{taxSaveError}</span>
                    )}
                  </div>
                </form>
              ) : (
                <>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-fg-muted">
                      {t.settings.panNumber}
                    </span>
                    <span className="text-sm text-fg font-mono">
                      {readOnlyValue(org?.pan_number)}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-fg-muted">
                      {t.settings.taxLegalName}
                    </span>
                    <span className="text-sm text-fg">
                      {readOnlyValue(org?.tax_legal_name)}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-fg-muted">
                      {t.settings.taxAddress}
                    </span>
                    <span className="text-sm text-fg text-right max-w-[60%]">
                      {readOnlyValue(org?.tax_address)}
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
