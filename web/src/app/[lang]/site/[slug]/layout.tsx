import type { ReactNode } from "react";
import { notFound } from "next/navigation";
import type { Locale } from "@/types";
import { fetchSite } from "@/lib/site/api";
import { resolveTheme, themeToCssVars } from "@/lib/site/themes";

/**
 * Tenant site chrome.
 *
 * The theme arrives as CSS custom properties on a wrapper element, which is
 * what lets five visually distinct templates share one set of section
 * renderers — no renderer knows which template it is inside.
 *
 * This route is only ever reached through the middleware rewrite from
 * <slug>.nepalgym.xyz; visitors never see the /site/<slug> path.
 */
export default async function SiteLayout({
  children,
  params,
}: {
  children: ReactNode;
  params: { lang: string; slug: string };
}) {
  const site = await fetchSite(params.slug);
  if (!site) notFound();

  const theme = resolveTheme(site.template, site.theme);
  const cssVars = themeToCssVars(theme);

  return (
    <div
      style={{
        ...cssVars,
        backgroundColor: "var(--site-bg)",
        color: "var(--site-fg)",
        fontFamily: "var(--site-font-body)",
        minHeight: "100vh",
      }}
      data-site-slug={site.slug}
      data-site-template={site.template}
    >
      {children}
    </div>
  );
}

export async function generateMetadata({
  params,
}: {
  params: { lang: string; slug: string };
}) {
  const site = await fetchSite(params.slug);
  if (!site) return { title: "Not found" };

  const locale = params.lang as Locale;
  const name = locale === "ne" && site.org_name_ne ? site.org_name_ne : site.org_name;

  return {
    title: { default: name, template: `%s · ${name}` },
  };
}
