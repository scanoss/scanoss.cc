import { Label } from '@/components/ui/label';
import useQueryState from '@/hooks/useQueryState';
import { MatchType } from '@/modules/results/domain';

import { matchTypeIconMap } from './Sidebar';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from './ui/select';

export default function MatchTypeSelector() {
  const [filterByMatchType, setFilterByMatchType] = useQueryState<
    MatchType | 'all'
  >('matchType', 'all');

  return (
    <div className="flex flex-col gap-1">
      <Label className="text-xs">Filter by match type</Label>
      <Select
        onValueChange={(value) => setFilterByMatchType(value as MatchType)}
        value={filterByMatchType}
      >
        <SelectTrigger className="w-full">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value={MatchType.File}>
            <div className="flex items-center gap-1">
              <span>{matchTypeIconMap[MatchType.File]}</span>
              <span>File</span>
            </div>
          </SelectItem>
          <SelectItem value={MatchType.Snippet}>
            <div className="flex items-center gap-1">
              <span>{matchTypeIconMap[MatchType.Snippet]}</span>
              <span>Snippet</span>
            </div>
          </SelectItem>
          <SelectItem value="all">All</SelectItem>
        </SelectContent>
      </Select>
    </div>
  );
}
