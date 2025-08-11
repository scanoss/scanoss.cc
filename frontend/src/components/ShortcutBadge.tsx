export default function ShortcutBadge({ shortcut }: { shortcut: string }) {
  return <span className="ml-2 rounded-sm bg-card p-1 text-[8px] leading-none">{shortcut}</span>;
}
