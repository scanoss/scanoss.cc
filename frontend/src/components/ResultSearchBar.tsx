import { Search } from 'lucide-react';
import { useRef } from 'react';

import { Input } from '@/components/ui/input';
import useEnvironment from '@/hooks/useEnvironment';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import useQueryState from '@/hooks/useQueryState';

export default function ResultSearchBar() {
  const { modifierKey } = useEnvironment();
  const inputRef = useRef<HTMLInputElement>(null);
  const [query, setQuery] = useQueryState<string>('q', '');

  useKeyboardShortcut([modifierKey.keyCode, 'f'], () => {
    inputRef.current?.focus();
  });

  return (
    <div className="relative grid w-full items-center gap-1.5">
      <Search className="absolute left-3 top-1/2 h-3 w-3 -translate-y-1/2 transform text-muted-foreground" />
      <Input
        ref={inputRef}
        className="pl-7"
        name="q"
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search"
        type="search"
        value={query}
        autoCapitalize="off"
      />
      <span className="absolute right-3 text-xs leading-none text-muted-foreground">⌘ + F</span>
    </div>
  );
}
