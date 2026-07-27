import { NextRequest, NextResponse } from "next/server";
import { locales, defaultLocale } from "@/lib/i18n";
import {
  resolveSiteRewrite,
  siteRewritePath,
  isAppOnlyPath,
  apexUrl,
} from "@/lib/subdomain";

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // Skip Next.js internals and static assets before anything else — these must
  // resolve identically on the apex and on every tenant subdomain.
  if (
    pathname.startsWith("/_next") ||
    pathname.startsWith("/api") ||
    pathname.includes(".")
  ) {
    return NextResponse.next();
  }

  // --- Tenant gym sites: <slug>.nepalgym.xyz ---
  //
  // Runs before the locale redirect, because a tenant URL carries its locale in
  // the same position but rewrites to a different internal route.
  const host = request.headers.get("host");
  const site = resolveSiteRewrite(host, pathname);

  if (site) {
    // The authenticated app lives on the apex only. Someone who follows a
    // dashboard link while on a gym's subdomain gets sent to the real app
    // rather than a confusing 404.
    if (isAppOnlyPath(pathname)) {
      return NextResponse.redirect(
        apexUrl(pathname, request.nextUrl.protocol)
      );
    }

    const url = request.nextUrl.clone();
    url.pathname = siteRewritePath(site);
    const response = NextResponse.rewrite(url);
    response.headers.set("x-site-slug", site.slug);
    return response;
  }

  // --- Apex host: unchanged locale handling ---

  const hasLocale = locales.some(
    (locale) => pathname.startsWith(`/${locale}/`) || pathname === `/${locale}`
  );

  if (hasLocale) return NextResponse.next();

  const url = request.nextUrl.clone();
  url.pathname = `/${defaultLocale}${pathname}`;
  return NextResponse.redirect(url);
}

export const config = {
  matcher: ["/((?!_next|api|favicon.ico).*)"],
};
