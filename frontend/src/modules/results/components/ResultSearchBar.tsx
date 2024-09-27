import { Search } from 'lucide-react';

import { Input } from '@/components/ui/input';
import useQueryState from '@/hooks/useQueryState';

export default function ResultSearchBar() {
  const [query, setQuery] = useQueryState<string>('q', '');

  return (
    <div className="relative grid w-full items-center gap-1.5">
      <Search className="absolute left-3 top-1/2 h-3 w-3 -translate-y-1/2 transform text-muted-foreground" />
      <Input
        className="pl-7"
        name="q"
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search"
        type="search"
        value={query}
      />
    </div>
  );
}
