import { isValidLocale, locales } from "@/lib/i18n";
import { redirect } from "next/navigation";
import { AuthProvider } from "@/lib/auth";
import { ThemeProvider } from "@/lib/theme";

export function generateStaticParams() {
  return locales.map((lang) => ({ lang }));
}

export default function LangLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: { lang: string };
}) {
  if (!isValidLocale(params.lang)) {
    redirect("/en");
  }

  return (
    <ThemeProvider>
      <AuthProvider>{children}</AuthProvider>
    </ThemeProvider>
  );
}
