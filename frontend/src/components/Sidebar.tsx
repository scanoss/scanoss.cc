import clsx from 'clsx';
import { Braces, ChevronRight, File } from 'lucide-react';
import { ReactNode, useEffect, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { entities } from 'wailsjs/go/models';

import useDebounce from '@/hooks/useDebounce';
import useQueryState from '@/hooks/useQueryState';
import { encodeFilePath, getDirectory, getFileName } from '@/lib/utils';
import useLocalFilePath from '@/modules/files/hooks/useLocalFilePath';
import ResultSearchBar from '@/modules/results/components/ResultSearchBar';
import { FilterAction, MatchType } from '@/modules/results/domain';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import MatchTypeSelector from '../modules/results/components/MatchTypeSelector';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from './ui/collapsible';
import { ScrollArea } from './ui/scroll-area';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';

export default function Sidebar() {
  const currentPath = useLocalFilePath();
  const navigate = useNavigate();

  const fetchResults = useResultsStore((state) => state.fetchResults);
  const anchorIndex = useResultsStore((state) => state.anchorIndex);
  const setAnchorIndex = useResultsStore((state) => state.setAnchorIndex);
  const setLastSelectedIndex = useResultsStore(
    (state) => state.setLastSelectedIndex
  );
  const selectedResults = useResultsStore((state) => state.selectedResults);
  const setSelectedResults = useResultsStore(
    (state) => state.setSelectedResults
  );
  const results = useResultsStore((state) => state.results);
  const pendingResults = useMemo(
    () => results.filter((r) => r.workflow_state === 'pending'),
    [results]
  );
  const completedResults = useMemo(
    () => results.filter((r) => r.workflow_state === 'completed'),
    [results]
  );

  const [filterByMatchType] = useQueryState<MatchType | 'all'>(
    'matchType',
    'all'
  );

  const [query] = useQueryState<string>('q', '');
  const debouncedQuery = useDebounce<string>(query, 300);

  const handleSelectFiles = (e: React.MouseEvent, path: string) => {
    e.preventDefault();

    const index = results.findIndex((result) => result.path === path);

    if (e.shiftKey) {
      handleSelectRange(path, index);
    } else if (e.metaKey) {
      handleSelectSingle(path);
      setAnchorIndex(index);
    } else {
      setSelectedResults(new Set([path]));
      setLastSelectedIndex(index);
      setAnchorIndex(index);
      return navigate({
        pathname: `/files/${encodeFilePath(path)}`,
        search: `?matchType=${filterByMatchType}&q=${query}`,
      });
    }
  };

  const handleSelectSingle = (selectedPath: string) => {
    const newSelectedResults = new Set(selectedResults);

    if (newSelectedResults.has(selectedPath)) {
      newSelectedResults.delete(selectedPath);
    } else {
      newSelectedResults.add(selectedPath);
    }

    setSelectedResults(newSelectedResults);
  };

  const handleSelectRange = (path: string, currentIndex: number) => {
    if (anchorIndex === -1) {
      setAnchorIndex(currentIndex);
      setSelectedResults(new Set([path]));
      setLastSelectedIndex(currentIndex);
      return;
    }

    const startIndex = Math.min(anchorIndex, currentIndex);
    const endIndex = Math.max(anchorIndex, currentIndex);

    const resultsToSelect = new Set<string>();

    for (let i = startIndex; i <= endIndex; i++) {
      resultsToSelect.add(results[i].path);
    }

    setSelectedResults(resultsToSelect);
    setLastSelectedIndex(currentIndex);
  };

  useEffect(() => {
    if (currentPath) {
      setSelectedResults(new Set([currentPath]));
      const index = results.findIndex((result) => result.path === currentPath);
      setLastSelectedIndex(index !== -1 ? index : -1);
      setAnchorIndex(index !== -1 ? index : -1);
    }
  }, [currentPath, results]);

  useEffect(() => {
    fetchResults(
      filterByMatchType === 'all' ? undefined : filterByMatchType,
      debouncedQuery
    );
  }, [filterByMatchType, debouncedQuery]);

  return (
    <aside className="flex h-full flex-col border-r border-border bg-black/20 backdrop-blur-md">
      <div className="flex h-[65px] items-center border-b border-b-border px-4">
        <h2 className="text-sm font-semibold">
          {pendingResults?.length
            ? `${pendingResults.length} decision${pendingResults.length > 1 ? 's' : ''} to make in working directory`
            : 'You have no decisions to make in working directory'}
        </h2>
      </div>

      <div className="flex flex-col gap-4 px-4 py-6">
        <MatchTypeSelector />
        <ResultSearchBar />
      </div>

      <ScrollArea>
        <div className="flex flex-1 flex-col gap-2">
          <Collapsible defaultOpen className="flex-1">
            <CollapsibleTrigger className="group w-full">
              <div className="flex items-center gap-1 border-b border-border px-3 pb-1">
                <ChevronRight className="group-data-[state=open]:stroke-text-foreground h-3 w-3 transform stroke-muted-foreground group-data-[state=open]:rotate-90" />
                <span className="text-sm text-muted-foreground">
                  Pending files{' '}
                  <span className="text-xs">({pendingResults?.length})</span>
                </span>
              </div>
            </CollapsibleTrigger>
            {pendingResults.length ? (
              <CollapsibleContent className="flex flex-col gap-1 overflow-y-auto py-2">
                {pendingResults.map((result) => {
                  return (
                    <SidebarItem
                      key={result.path}
                      result={result}
                      onSelect={handleSelectFiles}
                    />
                  );
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
                  <span className="text-xs">({completedResults?.length})</span>
                </span>
              </div>
            </CollapsibleTrigger>
            {completedResults.length ? (
              <CollapsibleContent className="flex flex-col gap-1 overflow-y-auto py-2">
                {completedResults.map((result) => {
                  return (
                    <SidebarItem
                      key={result.path}
                      result={result}
                      onSelect={handleSelectFiles}
                    />
                  );
                })}
              </CollapsibleContent>
            ) : null}
          </Collapsible>
        </div>
      </ScrollArea>
    </aside>
  );
}

function SidebarItem({
  result,
  onSelect,
}: {
  result: entities.ResultDTO;
  onSelect: (e: React.MouseEvent, path: string) => void;
}) {
  const selectedResults = useResultsStore((state) => state.selectedResults);
  const isActive = selectedResults.has(result.path);

  const isResultDismissed =
    result.filter_config?.action === FilterAction.Remove;
  const isResultIncluded =
    result.filter_config?.action === FilterAction.Include;

  const fileName = getFileName(result.path);
  const directory = getDirectory(result.path);

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <div
          className={clsx(
            'flex max-w-full cursor-pointer select-none items-center space-x-2 overflow-hidden px-4 py-1 text-sm text-muted-foreground',
            isActive
              ? 'border-r-2 border-primary-foreground bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground'
              : 'hover:bg-primary/30'
          )}
          onClick={(e) => onSelect(e, result.path)}
        >
          <span className="relative">
            {matchTypeIconMap[result.match_type as MatchType]}
            <span
              className={clsx(
                'absolute bottom-0 right-0 h-1 w-1 rounded-full',
                isResultDismissed && 'bg-red-600',
                isResultIncluded && 'bg-green-600'
              )}
            ></span>
          </span>
          <div className="flex min-w-0 items-center">
            {directory && <span className="truncate">{directory}</span>}
            <span
              className={clsx(
                'whitespace-nowrap',
                !isActive && 'text-foreground'
              )}
            >
              {directory ? `/${fileName}` : fileName}
            </span>
          </div>
        </div>
      </TooltipTrigger>
      <TooltipContent side="right" sideOffset={15}>
        {result.path}
      </TooltipContent>
    </Tooltip>
  );
}

export const matchTypeIconMap: Record<MatchType, ReactNode> = {
  [MatchType.File]: <File className="h-3 w-3 flex-shrink-0" />,
  [MatchType.Snippet]: <Braces className="h-3 w-3 flex-shrink-0" />,
};
