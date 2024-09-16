import { useQuery } from '@tanstack/react-query';
import clsx from 'clsx';
import { Braces, ChevronRight, File } from 'lucide-react';
import { ReactNode, useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';

import { decodeFilePath, encodeFilePath } from '@/lib/utils';
import { MatchType, Result } from '@/modules/results/domain';
import ResultService from '@/modules/results/infra/service';
import { useResults } from '@/modules/results/providers/ResultsProvider';

import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from './ui/collapsible';
import { ScrollArea } from './ui/scroll-area';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from './ui/select';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';

export default function Sidebar() {
  const { setResults, results, stagedResults, unstagedResults } = useResults();
  const [filterByMatchType, setFilterByMatchType] = useState<MatchType | 'all'>(
    'all'
  );

  const { data } = useQuery({
    queryKey: ['results', filterByMatchType],
    queryFn: () =>
      ResultService.getAll(
        filterByMatchType === 'all' ? undefined : filterByMatchType
      ),
  });

  useEffect(() => {
    if (data) {
      setResults(data);
    }
  }, [data]);

  return (
    <aside className="z-10 flex h-full flex-col border-r border-border bg-black/20 backdrop-blur-md">
      <div className="flex h-[65px] items-center border-b border-b-border px-4">
        <h2 className="text-sm font-semibold">
          {results?.length
            ? `${results.length} decision${results.length > 1 ? 's' : ''} to make in working directory`
            : 'You have no decisions to make in working directory'}
        </h2>
      </div>

      <div className="flex flex-col gap-1 px-4 py-6">
        <span className="text-xs font-semibold">Filter by match type</span>
        <Select
          onValueChange={(value) => setFilterByMatchType(value as MatchType)}
          defaultValue="all"
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
      <ScrollArea>
        <div className="flex flex-1 flex-col gap-2">
          <Collapsible defaultOpen className="flex-1">
            <CollapsibleTrigger className="group w-full">
              <div className="flex items-center gap-1 border-b border-border px-3 pb-1">
                <ChevronRight className="group-data-[state=open]:stroke-text-foreground h-3 w-3 transform stroke-muted-foreground group-data-[state=open]:rotate-90" />
                <span className="text-sm text-muted-foreground">
                  Pending files{' '}
                  <span className="text-xs">({unstagedResults?.length})</span>
                </span>
              </div>
            </CollapsibleTrigger>
            {unstagedResults.length ? (
              <CollapsibleContent className="flex flex-col gap-1 overflow-y-auto py-2">
                {unstagedResults.map((result) => {
                  return <SidebarItem key={result.path} result={result} />;
                })}
              </CollapsibleContent>
            ) : null}
          </Collapsible>
          <Collapsible defaultOpen className="flex-1">
            <CollapsibleTrigger className="group w-full">
              <div className="flex items-center gap-1 border-b border-border px-3 pb-1">
                <ChevronRight className="group-data-[state=open]:stroke-text-foreground h-3 w-3 transform stroke-muted-foreground group-data-[state=open]:rotate-90" />
                <span className="text-sm text-muted-foreground">
                  Completed files{' '}
                  <span className="text-xs">({stagedResults?.length})</span>
                </span>
              </div>
            </CollapsibleTrigger>
            {stagedResults.length ? (
              <CollapsibleContent className="flex flex-col gap-1 overflow-y-auto py-2">
                {stagedResults.map((result) => {
                  return <SidebarItem key={result.path} result={result} />;
                })}
              </CollapsibleContent>
            ) : null}
          </Collapsible>
        </div>
      </ScrollArea>
    </aside>
  );
}

function SidebarItem({ result }: { result: Result }) {
  const { filePath } = useParams();
  const isActive = decodeFilePath(filePath ?? '') === result.path;
  const encodedFilePath = encodeFilePath(result.path);

  return (
    <Tooltip key={result.path}>
      <TooltipTrigger asChild>
        <Link
          className={clsx(
            'flex w-full items-center gap-2 px-4 py-1 text-sm text-muted-foreground transition-all',
            isActive
              ? 'border-r-2 border-primary-foreground bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground'
              : 'hover:bg-primary/10'
          )}
          to={`files/${encodedFilePath}?matchType=${result.matchType}`}
        >
          {matchTypeIconMap[result.matchType]}
          <span className="max-w-[80%] overflow-hidden text-ellipsis">
            {result.path}
          </span>
        </Link>
      </TooltipTrigger>
      <TooltipContent side="right" sideOffset={15}>
        {result.path}
      </TooltipContent>
    </Tooltip>
  );
}

export const matchTypeIconMap: Record<MatchType, ReactNode> = {
  [MatchType.File]: <File className="h-3 w-3" />,
  [MatchType.Snippet]: <Braces className="h-3 w-3" />,
};
