"use client";

import { useState } from "react";

export default function DeleteAccount() {
  const [phone, setPhone] = useState("");
  const [submitted, setSubmitted] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!phone.match(/^[9876]\d{9}$/)) {
      setError("Please enter a valid 10-digit phone number.");
      return;
    }

    try {
      // Request OTP for verification
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL || ""}/api/v1/auth/login`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ phone }),
        }
      );

      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        setError(data.error || "Failed to send OTP. Please try again.");
        return;
      }

      setSubmitted(true);
    } catch {
      setError("Network error. Please try again.");
    }
  };

  const [otp, setOtp] = useState("");
  const [deleted, setDeleted] = useState(false);
  const [deleting, setDeleting] = useState(false);

  const handleDelete = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setDeleting(true);

    try {
      // Verify OTP to get token
      const verifyRes = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL || ""}/api/v1/auth/verify`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ phone, otp }),
        }
      );

      if (!verifyRes.ok) {
        const data = await verifyRes.json().catch(() => ({}));
        setError(data.error || "Invalid OTP. Please try again.");
        setDeleting(false);
        return;
      }

      const { token } = await verifyRes.json();

      // Delete account
      const deleteRes = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL || ""}/api/v1/members/me`,
        {
          method: "DELETE",
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      if (!deleteRes.ok) {
        const data = await deleteRes.json().catch(() => ({}));
        setError(data.error || "Failed to delete account. Please try again.");
        setDeleting(false);
        return;
      }

      setDeleted(true);
    } catch {
      setError("Network error. Please try again.");
      setDeleting(false);
    }
  };

  return (
    <div
      style={{
        maxWidth: 500,
        margin: "0 auto",
        padding: "40px 20px",
        fontFamily: "system-ui, -apple-system, sans-serif",
        color: "#e0e0e0",
        backgroundColor: "#0a0a0a",
        minHeight: "100vh",
        lineHeight: 1.7,
      }}
    >
      <h1 style={{ color: "#f5a623", marginBottom: 8 }}>Delete Account</h1>

      {deleted ? (
        <div>
          <p style={{ color: "#4caf50", fontSize: 18, marginBottom: 16 }}>
            Your account has been successfully deleted.
          </p>
          <p>
            All your personal data has been removed. You can uninstall the Urja
            app from your device.
          </p>
          <p style={{ color: "#888", marginTop: 16 }}>
            If you change your mind, you can create a new account at any time by
            logging in with your phone number.
          </p>
        </div>
      ) : !submitted ? (
        <div>
          <p style={{ marginBottom: 24 }}>
            To delete your Urja account, enter the phone number associated with
            your account. We will send you an OTP to verify your identity before
            proceeding.
          </p>

          <div
            style={{
              backgroundColor: "#1a1a1a",
              border: "1px solid #333",
              borderRadius: 8,
              padding: 20,
              marginBottom: 24,
            }}
          >
            <h3 style={{ color: "#ff6b6b", marginTop: 0 }}>
              What happens when you delete your account:
            </h3>
            <ul style={{ paddingLeft: 20 }}>
              <li>Your profile information will be permanently removed</li>
              <li>Your workout history and nutrition logs will be deleted</li>
              <li>Your gym memberships will be terminated</li>
              <li>This action cannot be undone</li>
            </ul>
          </div>

          <form onSubmit={handleSubmit}>
            <label
              htmlFor="phone"
              style={{ display: "block", marginBottom: 8, color: "#ccc" }}
            >
              Phone Number
            </label>
            <input
              id="phone"
              type="tel"
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
              placeholder="98XXXXXXXX"
              maxLength={10}
              style={{
                width: "100%",
                padding: "12px 16px",
                fontSize: 16,
                backgroundColor: "#1a1a1a",
                border: "1px solid #333",
                borderRadius: 8,
                color: "#e0e0e0",
                marginBottom: 16,
                boxSizing: "border-box",
              }}
            />

            {error && (
              <p style={{ color: "#ff6b6b", marginBottom: 16 }}>{error}</p>
            )}

            <button
              type="submit"
              style={{
                width: "100%",
                padding: "12px 24px",
                fontSize: 16,
                backgroundColor: "#dc3545",
                color: "#fff",
                border: "none",
                borderRadius: 8,
                cursor: "pointer",
                fontWeight: 600,
              }}
            >
              Send Verification Code
            </button>
          </form>
        </div>
      ) : (
        <div>
          <p style={{ marginBottom: 24 }}>
            We sent an OTP to <strong>{phone}</strong>. Enter it below to
            confirm account deletion.
          </p>

          <form onSubmit={handleDelete}>
            <label
              htmlFor="otp"
              style={{ display: "block", marginBottom: 8, color: "#ccc" }}
            >
              Verification Code (OTP)
            </label>
            <input
              id="otp"
              type="text"
              value={otp}
              onChange={(e) => setOtp(e.target.value)}
              placeholder="123456"
              maxLength={6}
              style={{
                width: "100%",
                padding: "12px 16px",
                fontSize: 16,
                backgroundColor: "#1a1a1a",
                border: "1px solid #333",
                borderRadius: 8,
                color: "#e0e0e0",
                marginBottom: 16,
                boxSizing: "border-box",
                letterSpacing: 4,
                textAlign: "center",
              }}
            />

            {error && (
              <p style={{ color: "#ff6b6b", marginBottom: 16 }}>{error}</p>
            )}

            <button
              type="submit"
              disabled={deleting}
              style={{
                width: "100%",
                padding: "12px 24px",
                fontSize: 16,
                backgroundColor: deleting ? "#666" : "#dc3545",
                color: "#fff",
                border: "none",
                borderRadius: 8,
                cursor: deleting ? "not-allowed" : "pointer",
                fontWeight: 600,
              }}
            >
              {deleting ? "Deleting..." : "Delete My Account"}
            </button>

            <button
              type="button"
              onClick={() => {
                setSubmitted(false);
                setOtp("");
                setError("");
              }}
              style={{
                width: "100%",
                padding: "12px 24px",
                fontSize: 16,
                backgroundColor: "transparent",
                color: "#888",
                border: "1px solid #333",
                borderRadius: 8,
                cursor: "pointer",
                marginTop: 12,
              }}
            >
              Go Back
            </button>
          </form>
        </div>
      )}

      <p style={{ color: "#666", marginTop: 40, fontSize: 14 }}>
        If you have any questions, contact us at{" "}
        <a href="mailto:support@nepalgym.xyz" style={{ color: "#f5a623" }}>
          support@nepalgym.xyz
        </a>
      </p>
    </div>
  );
}
