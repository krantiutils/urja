import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Privacy Policy - Urja",
};

export default function PrivacyPolicy() {
  return (
    <div
      style={{
        maxWidth: 800,
        margin: "0 auto",
        padding: "40px 20px",
        fontFamily: "system-ui, -apple-system, sans-serif",
        color: "#e0e0e0",
        backgroundColor: "#0a0a0a",
        minHeight: "100vh",
        lineHeight: 1.7,
      }}
    >
      <h1 style={{ color: "#f5a623", marginBottom: 8 }}>Privacy Policy</h1>
      <p style={{ color: "#888", marginBottom: 32 }}>
        Last updated: March 4, 2026
      </p>

      <p>
        Urja (&quot;we&quot;, &quot;our&quot;, or &quot;us&quot;) operates the
        Urja mobile application and web platform (collectively, the
        &quot;Service&quot;). This Privacy Policy explains how we collect, use,
        and protect your personal information when you use our Service.
      </p>

      <h2>1. Information We Collect</h2>

      <h3>Account Information</h3>
      <ul>
        <li>Phone number (used for authentication via OTP)</li>
        <li>Name</li>
        <li>Age, gender</li>
      </ul>

      <h3>Health &amp; Fitness Data</h3>
      <ul>
        <li>Body measurements (weight, height)</li>
        <li>Fitness goals and activity level</li>
        <li>Workout logs and exercise history</li>
        <li>Nutrition logs (food intake, water intake, calorie tracking)</li>
        <li>Weight progress history</li>
        <li>Gym attendance records</li>
      </ul>

      <h3>Gym Membership Data</h3>
      <ul>
        <li>Gym/organization affiliation</li>
        <li>Membership and subscription details</li>
        <li>Package and billing information</li>
      </ul>

      <h3>Automatically Collected Information</h3>
      <ul>
        <li>Device type and operating system</li>
        <li>App version</li>
        <li>Usage data (pages visited, features used)</li>
      </ul>

      <h2>2. How We Use Your Information</h2>
      <p>We use the information we collect to:</p>
      <ul>
        <li>Provide and maintain the Service</li>
        <li>Authenticate your identity via OTP</li>
        <li>Calculate personalized nutrition and fitness targets</li>
        <li>Track your workout progress and attendance</li>
        <li>Connect you with your gym or fitness organization</li>
        <li>Send service-related notifications (e.g., membership reminders)</li>
        <li>Improve and optimize the Service</li>
      </ul>

      <h2>3. Data Sharing</h2>
      <p>
        We do <strong>not</strong> sell your personal information. We may share
        your data only in the following cases:
      </p>
      <ul>
        <li>
          <strong>With your gym/organization:</strong> Your gym administrators
          can view your membership status, attendance, and assigned workout
          plans.
        </li>
        <li>
          <strong>SMS provider:</strong> Your phone number is shared with our
          SMS provider solely to deliver OTP codes for authentication.
        </li>
        <li>
          <strong>Legal requirements:</strong> If required by law, regulation, or
          legal process.
        </li>
      </ul>

      <h2>4. Data Storage &amp; Security</h2>
      <p>
        Your data is stored on secure servers. We use industry-standard security
        measures including encrypted connections (HTTPS/TLS), secure token-based
        authentication (JWT), and database access controls to protect your
        information.
      </p>

      <h2>5. Data Retention</h2>
      <p>
        We retain your personal data for as long as your account is active. If
        you request account deletion, we will delete your personal data within
        30 days, except where retention is required by law.
      </p>

      <h2>6. Your Rights</h2>
      <p>You have the right to:</p>
      <ul>
        <li>Access the personal data we hold about you</li>
        <li>Correct inaccurate information</li>
        <li>Request deletion of your account and data</li>
        <li>Export your data</li>
      </ul>

      <h2>7. Children&apos;s Privacy</h2>
      <p>
        The Service is not intended for children under the age of 16. We do not
        knowingly collect personal information from children under 16. If you
        believe we have collected data from a child, please contact us
        immediately.
      </p>

      <h2>8. Third-Party Services</h2>
      <p>
        The Service does not include third-party analytics, advertising SDKs, or
        social media trackers. We do not use Google Analytics, Facebook SDK, or
        similar tracking services.
      </p>

      <h2>9. Changes to This Policy</h2>
      <p>
        We may update this Privacy Policy from time to time. We will notify you
        of significant changes through the app or via SMS. Your continued use of
        the Service after changes constitutes acceptance of the updated policy.
      </p>

      <h2>10. Contact Us</h2>
      <p>
        If you have questions about this Privacy Policy or wish to exercise your
        data rights, contact us at:
      </p>
      <ul>
        <li>
          Email:{" "}
          <a href="mailto:support@nepalgym.xyz" style={{ color: "#f5a623" }}>
            support@nepalgym.xyz
          </a>
        </li>
      </ul>
    </div>
  );
}
