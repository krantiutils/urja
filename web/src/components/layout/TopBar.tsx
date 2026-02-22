"use client";

import { useState, useRef, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import type { Dictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
import {
  Menu,
  QrCode,
  ChevronDown,
  User,
  LogOut,
  Globe,
  Download,
  Printer,
  X,
  Dumbbell,
  Smartphone,
} from "lucide-react";

interface TopBarProps {
  t: Dictionary;
  locale: Locale;
  onMenuToggle: () => void;
}

export function TopBar({ t, locale, onMenuToggle }: TopBarProps) {
  const { user, logout } = useAuth();
  const router = useRouter();
  const [profileOpen, setProfileOpen] = useState(false);
  const [qrOpen, setQrOpen] = useState(false);
  const [qrBlobUrl, setQrBlobUrl] = useState<string | null>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const cardRef = useRef<HTMLDivElement>(null);

  // Close dropdown on outside click
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(e.target as Node)
      ) {
        setProfileOpen(false);
      }
    };
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, []);

  // Fetch QR image with auth header when modal opens
  useEffect(() => {
    if (!qrOpen || !user?.org_id) return;
    let revoked = false;
    const url = api.getQrCodeUrl(user.org_id);
    const token = localStorage.getItem("access_token");
    fetch(url, {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    })
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.blob();
      })
      .then((blob) => {
        if (!revoked) setQrBlobUrl(URL.createObjectURL(blob));
      })
      .catch((err) => console.error("Failed to load QR code:", err));
    return () => {
      revoked = true;
      setQrBlobUrl((prev) => {
        if (prev) URL.revokeObjectURL(prev);
        return null;
      });
    };
  }, [qrOpen, user?.org_id]);

  const handleLogout = async () => {
    setProfileOpen(false);
    await logout();
    router.replace(`/${locale}/login`);
  };

  const toggleLocale = () => {
    const newLocale = locale === "en" ? "ne" : "en";
    const currentPath = window.location.pathname;
    const newPath = currentPath.replace(`/${locale}/`, `/${newLocale}/`);
    router.push(newPath);
  };

  const handleDownloadQr = async () => {
    if (!qrBlobUrl) return;

    const scale = 3; // 3x for crisp print quality
    const W = 280, H = 440;
    const canvas = document.createElement("canvas");
    canvas.width = W * scale;
    canvas.height = H * scale;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;
    ctx.scale(scale, scale);

    // Helper: rounded rect
    const rr = (x: number, y: number, w: number, h: number, r: number) => {
      ctx.beginPath();
      ctx.moveTo(x + r, y);
      ctx.arcTo(x + w, y, x + w, y + h, r);
      ctx.arcTo(x + w, y + h, x, y + h, r);
      ctx.arcTo(x, y + h, x, y, r);
      ctx.arcTo(x, y, x + w, y, r);
      ctx.closePath();
    };

    // Card background
    rr(0, 0, W, H, 16);
    ctx.fillStyle = "#ffffff";
    ctx.fill();
    ctx.strokeStyle = "#e5e7eb";
    ctx.lineWidth = 1.5;
    ctx.stroke();

    // Green icon background
    rr(W / 2 - 22, 28, 44, 44, 12);
    ctx.fillStyle = "#84cc16";
    ctx.fill();

    // Dumbbell icon (simplified cross shape)
    ctx.strokeStyle = "#ffffff";
    ctx.lineWidth = 2.5;
    ctx.lineCap = "round";
    ctx.beginPath();
    ctx.moveTo(W / 2 - 10, 42); ctx.lineTo(W / 2 + 10, 58);
    ctx.moveTo(W / 2 - 12, 44); ctx.lineTo(W / 2 - 8, 40);
    ctx.moveTo(W / 2 + 8, 60); ctx.lineTo(W / 2 + 12, 56);
    ctx.moveTo(W / 2 - 7, 47); ctx.lineTo(W / 2 - 3, 43);
    ctx.moveTo(W / 2 + 3, 57); ctx.lineTo(W / 2 + 7, 53);
    ctx.stroke();

    // Gym name
    ctx.fillStyle = "#111827";
    ctx.font = "bold 18px -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif";
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    ctx.fillText(gymName, W / 2, 96);

    // QR image
    const qrImg = new Image();
    qrImg.src = qrBlobUrl;
    await new Promise<void>((resolve) => { qrImg.onload = () => resolve(); qrImg.onerror = () => resolve(); });
    // QR frame background
    rr(34, 118, 212, 212, 14);
    ctx.fillStyle = "#fafafa";
    ctx.fill();
    ctx.strokeStyle = "#f0f0f0";
    ctx.lineWidth = 1.5;
    ctx.stroke();
    ctx.drawImage(qrImg, 50, 134, 180, 180);

    // "SCAN TO CHECK IN"
    ctx.fillStyle = "#374151";
    ctx.font = "600 13px -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif";
    ctx.textAlign = "center";
    ctx.fillText("SCAN TO CHECK IN", W / 2 + 4, 358);
    // Phone icon (simple rect)
    rr(W / 2 - 70, 349, 14, 18, 3);
    ctx.strokeStyle = "#6b7280";
    ctx.lineWidth = 1.5;
    ctx.stroke();

    // Divider
    ctx.beginPath();
    ctx.moveTo(W / 2 - 24, 386);
    ctx.lineTo(W / 2 + 24, 386);
    ctx.strokeStyle = "#e5e7eb";
    ctx.lineWidth = 1;
    ctx.stroke();

    // "POWERED BY URJA"
    ctx.font = "500 10px -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif";
    ctx.fillStyle = "#9ca3af";
    ctx.textAlign = "center";
    const pwText = "POWERED BY ";
    const pwWidth = ctx.measureText(pwText).width;
    const urjaWidth = ctx.measureText("URJA").width;
    const totalWidth = pwWidth + urjaWidth;
    const startX = W / 2 - totalWidth / 2;
    ctx.textAlign = "left";
    ctx.fillText(pwText, startX, 406);
    ctx.fillStyle = "#84cc16";
    ctx.font = "bold 10px -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif";
    ctx.fillText("URJA", startX + pwWidth, 406);

    // Download
    canvas.toBlob((blob) => {
      if (!blob) return;
      const url = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = `${gymName.replace(/\s+/g, "-").toLowerCase()}-qr-card.png`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
    }, "image/png");
  };

  const gymName = user?.org_name || "My Gym";

  const handlePrintQr = () => {
    if (!user?.org_id) return;
    const printWindow = window.open("", "_blank");
    if (!printWindow) return;
    printWindow.document.write(`
      <!DOCTYPE html>
      <html>
        <head>
          <title>${gymName} - QR Check-in Card</title>
          <style>
            @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');

            * { margin: 0; padding: 0; box-sizing: border-box; }

            html, body {
              width: 100%;
              height: 100%;
              display: flex;
              justify-content: center;
              align-items: center;
              background: #fff;
              font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
              -webkit-print-color-adjust: exact;
              print-color-adjust: exact;
            }

            .card {
              width: 280px;
              padding: 32px 24px 20px;
              background: #ffffff;
              border: 1.5px solid #e5e7eb;
              border-radius: 16px;
              display: flex;
              flex-direction: column;
              align-items: center;
              gap: 0;
            }

            .gym-icon {
              width: 44px;
              height: 44px;
              background: linear-gradient(135deg, #84cc16, #65a30d);
              border-radius: 12px;
              display: flex;
              align-items: center;
              justify-content: center;
              margin-bottom: 12px;
            }

            .gym-icon svg {
              width: 24px;
              height: 24px;
              color: #fff;
            }

            .gym-name {
              font-size: 18px;
              font-weight: 700;
              color: #111827;
              text-align: center;
              letter-spacing: -0.02em;
              line-height: 1.3;
              margin-bottom: 24px;
            }

            .qr-frame {
              background: #fafafa;
              border: 1.5px solid #f0f0f0;
              border-radius: 14px;
              padding: 16px;
              margin-bottom: 20px;
            }

            .qr-frame img {
              width: 180px;
              height: 180px;
              display: block;
            }

            .scan-row {
              display: flex;
              align-items: center;
              gap: 8px;
              margin-bottom: 24px;
            }

            .scan-row svg {
              width: 16px;
              height: 16px;
              color: #6b7280;
              flex-shrink: 0;
            }

            .scan-text {
              font-size: 13px;
              font-weight: 600;
              color: #374151;
              letter-spacing: 0.03em;
              text-transform: uppercase;
            }

            .divider {
              width: 48px;
              height: 1px;
              background: #e5e7eb;
              margin-bottom: 14px;
            }

            .powered-by {
              font-size: 10px;
              font-weight: 500;
              color: #9ca3af;
              letter-spacing: 0.06em;
              text-transform: uppercase;
            }

            .powered-by span {
              color: #84cc16;
              font-weight: 700;
            }

            @media print {
              html, body {
                background: #fff;
              }
              .card {
                border: 1.5px solid #d1d5db;
                box-shadow: none;
              }
            }
          </style>
        </head>
        <body>
          <div class="card">
            <div class="gym-icon">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="m6.5 6.5 11 11"/><path d="m21 21-1-1"/><path d="m3 3 1 1"/><path d="m18 22 4-4"/><path d="m2 6 4-4"/><path d="m3 10 7-7"/><path d="m14 21 7-7"/>
              </svg>
            </div>
            <div class="gym-name">${gymName}</div>
            <div class="qr-frame">
              <img src="${qrBlobUrl ?? ""}" alt="QR Code" />
            </div>
            <div class="scan-row">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect width="14" height="20" x="5" y="2" rx="2" ry="2"/><path d="M12 18h.01"/>
              </svg>
              <span class="scan-text">${t.topbar.scanToCheckIn}</span>
            </div>
            <div class="divider"></div>
            <div class="powered-by">${t.topbar.poweredBy}</div>
          </div>
          <script>
            document.querySelector('img').onload = function() {
              setTimeout(function() { window.print(); window.close(); }, 200);
            };
          <\/script>
        </body>
      </html>
    `);
    printWindow.document.close();
  };

  return (
    <>
    <header className="sticky top-0 z-30 h-14 bg-bg-elevated/80 backdrop-blur-xl border-b border-white/[0.06] flex items-center justify-between px-4 lg:px-6">
      {/* Left: menu toggle (mobile) + page title area */}
      <div className="flex items-center gap-3">
        <button
          onClick={onMenuToggle}
          className="lg:hidden p-2 rounded-lg hover:bg-surface transition-colors"
          aria-label="Toggle menu"
        >
          <Menu className="w-5 h-5 text-fg-muted" />
        </button>
      </div>

      {/* Right: actions */}
      <div className="flex items-center gap-2">
        {/* Language toggle */}
        <button
          onClick={toggleLocale}
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm text-fg-muted hover:bg-surface hover:text-fg transition-colors"
          title={locale === "en" ? "नेपालीमा हेर्नुहोस्" : "View in English"}
        >
          <Globe className="w-4 h-4" />
          <span className="text-xs font-mono uppercase">
            {locale === "en" ? "NE" : "EN"}
          </span>
        </button>

        {/* QR Code button */}
        <button
          onClick={() => setQrOpen(true)}
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm text-fg-muted hover:bg-surface hover:text-fg transition-colors"
          title={t.topbar.orgQr}
        >
          <QrCode className="w-4 h-4" />
          <span className="hidden sm:inline">{t.topbar.orgQr}</span>
        </button>

        {/* Profile dropdown */}
        <div className="relative" ref={dropdownRef}>
          <button
            onClick={() => setProfileOpen(!profileOpen)}
            className="flex items-center gap-2 px-3 py-1.5 rounded-lg hover:bg-surface transition-colors"
          >
            <div className="w-7 h-7 rounded-full bg-accent/20 flex items-center justify-center">
              <User className="w-3.5 h-3.5 text-accent" />
            </div>
            <span className="hidden sm:inline text-sm text-fg-muted">
              {user?.phone ?? ""}
            </span>
            <ChevronDown
              className={`w-3 h-3 text-fg-muted transition-transform duration-200 ${
                profileOpen ? "rotate-180" : ""
              }`}
            />
          </button>

          {profileOpen && (
            <div className="absolute right-0 top-full mt-1 w-48 bg-bg-elevated border border-white/[0.06] rounded-xl shadow-card overflow-hidden animate-fade-in">
              <div className="px-3 py-2 border-b border-white/[0.06]">
                <p className="text-sm text-fg font-medium truncate">
                  {user?.phone}
                </p>
                <p className="text-xs text-fg-muted capitalize">
                  {user?.role}
                </p>
              </div>
              <div className="py-1">
                <button
                  onClick={handleLogout}
                  className="flex items-center gap-2 w-full px-3 py-2 text-sm text-red-400 hover:bg-surface transition-colors"
                >
                  <LogOut className="w-4 h-4" />
                  {t.topbar.logout}
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </header>

    {/* QR Card Modal */}
    {qrOpen && (
      <div
        className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4"
        onClick={() => setQrOpen(false)}
      >
        <div
          className="relative flex flex-col items-center gap-5 animate-fade-in"
          onClick={(e) => e.stopPropagation()}
        >
          {/* Close button */}
          <button
            onClick={() => setQrOpen(false)}
            className="absolute -top-2 -right-2 z-10 p-1.5 rounded-full bg-bg-elevated border border-white/[0.08] text-fg-muted hover:text-fg hover:bg-surface transition-colors shadow-lg"
          >
            <X className="w-4 h-4" />
          </button>

          {/* Printable card preview */}
          <div
            ref={cardRef}
            className="w-[280px] bg-white rounded-2xl shadow-[0_8px_60px_rgba(0,0,0,0.5),0_0_0_1px_rgba(255,255,255,0.08)] flex flex-col items-center px-6 pt-8 pb-5"
          >
            {/* Gym icon */}
            <div className="w-11 h-11 rounded-xl bg-gradient-to-br from-lime-500 to-lime-600 flex items-center justify-center mb-3 shadow-[0_2px_8px_rgba(132,204,22,0.3)]">
              <Dumbbell className="w-6 h-6 text-white" />
            </div>

            {/* Gym name */}
            <h2 className="text-lg font-bold text-gray-900 text-center tracking-tight leading-tight mb-6">
              {gymName}
            </h2>

            {/* QR code in subtle frame */}
            {user?.org_id && (
              <div className="bg-gray-50 border border-gray-100 rounded-xl p-4 mb-5">
                {qrBlobUrl ? (
                  /* eslint-disable-next-line @next/next/no-img-element */
                  <img
                    src={qrBlobUrl}
                    alt="QR Code"
                    className="w-[180px] h-[180px] block"
                  />
                ) : (
                  <div className="w-[180px] h-[180px] flex items-center justify-center">
                    <div className="w-6 h-6 border-2 border-gray-300 border-t-lime-500 rounded-full animate-spin" />
                  </div>
                )}
              </div>
            )}

            {/* Scan to check in */}
            <div className="flex items-center gap-2 mb-6">
              <Smartphone className="w-4 h-4 text-gray-500" />
              <span className="text-[13px] font-semibold text-gray-700 uppercase tracking-wider">
                {t.topbar.scanToCheckIn}
              </span>
            </div>

            {/* Divider */}
            <div className="w-12 h-px bg-gray-200 mb-3.5" />

            {/* Powered by */}
            <p className="text-[10px] font-medium text-gray-400 uppercase tracking-widest">
              {t.topbar.poweredBy.replace("Urja", "").trim()}{" "}
              <span className="text-lime-500 font-bold">Urja</span>
            </p>
          </div>

          {/* Action buttons below card */}
          <div className="flex gap-3 w-full">
            <button
              onClick={handleDownloadQr}
              className="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors"
            >
              <Download className="w-4 h-4" />
              {t.topbar.download}
            </button>
            <button
              onClick={handlePrintQr}
              className="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 bg-surface border border-white/[0.06] text-fg text-sm rounded-xl hover:bg-surface-hover transition-colors"
            >
              <Printer className="w-4 h-4" />
              {t.topbar.print}
            </button>
          </div>
        </div>
      </div>
    )}
    </>
  );
}
