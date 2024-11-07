import { Search } from 'lucide-react';

import { Input } from '@/components/ui/input';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import useResultsStore from '@/modules/results/stores/useResultsStore';

export default function ResultSearchBar({ searchInputRef }: { searchInputRef: React.RefObject<HTMLInputElement> }) {
  const query = useResultsStore((state) => state.query);
  const setQuery = useResultsStore((state) => state.setQuery);
  const setSelectedResults = useResultsStore((state) => state.setSelectedResults);

  useKeyboardShortcut(KEYBOARD_SHORTCUTS.focusSearch.keys, () => searchInputRef.current?.focus());
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.selectAll.keys, () => searchInputRef.current?.select());

  return (
    <div className="relative grid w-full items-center gap-1.5">
      <Search className="absolute left-3 top-1/2 h-3 w-3 -translate-y-1/2 transform text-muted-foreground" />
      <Input
        ref={searchInputRef}
        className="pl-7"
        name="q"
        onChange={(e) => {
          e.preventDefault();
          setSelectedResults([]);
          setQuery(e.target.value);
        }}
        placeholder="Search"
        value={query}
        type="text"
      />
      <span className="absolute right-3 text-xs leading-none text-muted-foreground">⌘ + F</span>
    </div>
  );
}
