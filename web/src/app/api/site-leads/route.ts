import { NextResponse } from "next/server";

/**
 * Same-origin proxy for tenant site enquiry forms.
 *
 * The form on ibckirtipur.nepalgym.xyz posts here rather than straight to the
 * API for two reasons: it is the gym's only conversion path, and posting
 * cross-origin to the apex would make it depend on both a correct CORS policy
 * and `NEXT_PUBLIC_API_URL` being baked in at image-build time. Neither fails
 * loudly — a misconfigured build just silently stops capturing leads.
 *
 * This runs server-side with the same internal API base the page renderer
 * already uses, so it works on every tenant host with no per-tenant config.
 */

export const runtime = "nodejs";
/** Leads are writes; nothing here may be cached. */
export const dynamic = "force-dynamic";

const MAX_BODY_BYTES = 16 * 1024;

function apiBase(): string {
  return (
    process.env.INTERNAL_API_URL?.trim() ||
    process.env.NEXT_PUBLIC_API_URL?.trim() ||
    "http://localhost:8080"
  );
}

/**
 * The API rate-limits by client IP, so the real caller must be forwarded —
 * otherwise every gym's enquiries share this server's single bucket and one
 * busy site throttles the rest.
 */
function clientIp(request: Request): string | null {
  const forwarded = request.headers.get("x-forwarded-for");
  if (forwarded) return forwarded.split(",")[0].trim();
  return request.headers.get("x-real-ip");
}

export async function POST(request: Request) {
  const raw = await request.text();
  if (raw.length > MAX_BODY_BYTES) {
    return NextResponse.json({ error: "request too large" }, { status: 413 });
  }

  let body: Record<string, unknown>;
  try {
    body = JSON.parse(raw);
  } catch {
    return NextResponse.json({ error: "invalid request body" }, { status: 400 });
  }

  const slug = typeof body.site === "string" ? body.site : "";
  if (!/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(slug)) {
    return NextResponse.json({ error: "unknown site" }, { status: 400 });
  }

  // `site` addresses the endpoint; it is not part of the lead itself.
  const lead = { ...body };
  delete lead.site;

  const headers: Record<string, string> = { "Content-Type": "application/json" };
  const ip = clientIp(request);
  if (ip) {
    headers["X-Forwarded-For"] = ip;
    headers["X-Real-IP"] = ip;
  }

  try {
    const res = await fetch(
      `${apiBase()}/api/v1/sites/${encodeURIComponent(slug)}/leads`,
      { method: "POST", headers, body: JSON.stringify(lead), cache: "no-store" }
    );

    const text = await res.text();
    const payload = text ? JSON.parse(text) : {};
    return NextResponse.json(payload, { status: res.status });
  } catch {
    return NextResponse.json(
      { error: "Could not send. Please try again." },
      { status: 502 }
    );
  }
}
