import { redirect } from "next/navigation";

export default function LangRootPage({
  params,
}: {
  params: { lang: string };
}) {
  redirect(`/${params.lang}/dashboard`);
}
