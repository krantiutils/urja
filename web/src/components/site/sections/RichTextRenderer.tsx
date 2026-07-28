import type { CSSProperties, ReactNode } from "react";
import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { text, parseMarkdownBlocks, type MarkdownBlock } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";

/**
 * rich_text — standard | two_column | centered
 */

const headingStyle: CSSProperties = {
  fontFamily: "var(--site-font-display)",
  textTransform: "var(--site-display-transform)" as never,
};

function renderBlocks(blocks: MarkdownBlock[]): ReactNode[] {
  const out: ReactNode[] = [];
  let listBuffer: string[] = [];
  let key = 0;

  const flushList = () => {
    if (listBuffer.length === 0) return;
    out.push(
      <ul key={`ul-${key++}`} className="mt-4 list-disc pl-5 space-y-1 first:mt-0">
        {listBuffer.map((t, i) => (
          <li key={i}>{t}</li>
        ))}
      </ul>
    );
    listBuffer = [];
  };

  for (const block of blocks) {
    if (block.kind === "listItem") {
      listBuffer.push(block.text);
      continue;
    }
    flushList();
    if (block.kind === "heading") {
      const Tag = block.level === 2 ? "h3" : "h4";
      out.push(
        <Tag key={key++} className="mt-6 text-lg font-medium first:mt-0" style={headingStyle}>
          {block.text}
        </Tag>
      );
    } else {
      out.push(
        <p key={key++} className="mt-4 leading-relaxed text-[var(--site-fg-muted)] first:mt-0">
          {block.text}
        </p>
      );
    }
  }
  flushList();
  return out;
}

export default function RichTextRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const { content, variant } = section;
  const title = text(content, "title", locale);
  const body = text(content, "body", locale);
  const blocks = parseMarkdownBlocks(body);

  if (!title && blocks.length === 0) return null;

  if (variant === "centered") {
    return (
      <div className="mx-auto max-w-2xl text-center">
        <SectionHeading title={title} align="center" />
        <div>{renderBlocks(blocks)}</div>
      </div>
    );
  }

  if (variant === "two_column") {
    return (
      <>
        <SectionHeading title={title} align={section.style?.align ?? "left"} />
        <div className="columns-1 md:columns-2 gap-10 [&>*]:break-inside-avoid">{renderBlocks(blocks)}</div>
      </>
    );
  }

  // standard
  return (
    <>
      <SectionHeading title={title} align={section.style?.align ?? "left"} />
      <div>{renderBlocks(blocks)}</div>
    </>
  );
}
