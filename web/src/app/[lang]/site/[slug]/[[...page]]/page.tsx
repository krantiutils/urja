import { notFound } from "next/navigation";
import type { Locale } from "@/types";
import { fetchSite, fetchPage, fetchPackages, SITE_REVALIDATE_SECONDS } from "@/lib/site/api";
import { SiteHeader } from "@/components/site/SiteHeader";
import { SiteFooter } from "@/components/site/SiteFooter";
import { SectionRenderer } from "@/components/site/SectionRenderer";
import { BASE_DOMAIN } from "@/lib/subdomain";

export const revalidate = SITE_REVALIDATE_SECONDS;

/**
 * A tenant site page.
 *
 * The optional catch-all resolves to a page slug, defaulting to "home" for the
 * subdomain root. Only published pages of a live site are served — the API
 * returns 404 for anything else, so drafts cannot be discovered by guessing.
 */
export default async function SitePage({
  params,
}: {
  params: { lang: string; slug: string; page?: string[] };
}) {
  const locale = params.lang as Locale;
  const pageSlug = params.page?.join("/") || "home";

  const [site, page] = await Promise.all([
    fetchSite(params.slug),
    fetchPage(params.slug, pageSlug),
  ]);

  if (!site || !page) notFound();

  // membership_plans sections with dataSource "auto" read the gym's real
  // packages. Fetched once here rather than inside a renderer, so a page with
  // two such sections does not issue two requests.
  const needsPackages = page.sections.some(
    (s) => s.type === "membership_plans" && s.content?.dataSource !== "manual"
  );
  const packages = needsPackages ? await fetchPackages(params.slug) : [];

  return (
    <>
      <SiteHeader
        site={site}
        locale={locale}
        currentPage={pageSlug}
        apexUrl={`https://${BASE_DOMAIN}/${locale}/login`}
      />

      <main>
        {page.sections.map((section) => (
          <SectionRenderer
            key={section.id}
            section={section}
            locale={locale}
            siteSlug={site.slug}
            packages={packages}
            sourcePage={pageSlug}
          />
        ))}
      </main>

      <SiteFooter site={site} locale={locale} />
    </>
  );
}

export async function generateMetadata({
  params,
}: {
  params: { lang: string; slug: string; page?: string[] };
}) {
  const pageSlug = params.page?.join("/") || "home";
  const [site, page] = await Promise.all([
    fetchSite(params.slug),
    fetchPage(params.slug, pageSlug),
  ]);

  if (!site || !page) return { title: "Not found" };

  const locale = params.lang as Locale;
  const title = locale === "ne" && page.title_ne ? page.title_ne : page.title;
  const orgName = locale === "ne" && site.org_name_ne ? site.org_name_ne : site.org_name;
  const description = page.seo_description || undefined;

  return {
    title,
    description,
    openGraph: {
      title: `${title} · ${orgName}`,
      description,
      type: "website",
    },
    alternates: {
      canonical: pageSlug === "home" ? "/" : `/${pageSlug}`,
    },
  };
}
