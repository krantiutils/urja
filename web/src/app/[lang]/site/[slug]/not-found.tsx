import Link from "next/link";

/**
 * Shown when a subdomain has no matching gym, the gym's site is not live, or a
 * page slug does not resolve to a published page.
 *
 * All three cases render identically on purpose: an unfinished site should not
 * be distinguishable from one that does not exist.
 */
export default function SiteNotFound() {
  return (
    <div className="min-h-screen flex items-center justify-center px-6 bg-[var(--site-bg,#08080a)] text-[var(--site-fg,#f5f5f7)]">
      <div className="text-center max-w-md">
        <p className="text-sm tracking-widest opacity-50 mb-4">404</p>
        <h1 className="text-2xl sm:text-3xl font-semibold mb-3">
          This page isn&apos;t here
        </h1>
        <p className="text-sm opacity-70 leading-relaxed mb-8">
          The gym you are looking for may not have published its website yet.
        </p>
        <Link
          href="https://nepalgym.xyz"
          className="inline-block px-5 py-2.5 text-sm border border-current/20 rounded-lg hover:opacity-80 transition-opacity"
        >
          Go to nepalgym.xyz
        </Link>
      </div>
    </div>
  );
}
