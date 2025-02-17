// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

import { useVirtualizer } from '@tanstack/react-virtual';
import clsx from 'clsx';
import { Braces, ChevronRight, File } from 'lucide-react';
import { ReactNode, useEffect, useRef } from 'react';
import { entities } from 'wailsjs/go/models';

import ResultSearchBar from '@/components/ResultSearchBar';
import useDebounce from '@/hooks/useDebounce';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import { getDirectory, getFileName } from '@/lib/utils';
import { FilterAction } from '@/modules/components/domain';
import { DEBOUNCE_QUERY_MS } from '@/modules/results/constants';
import { MatchType, stateInfoPresentation } from '@/modules/results/domain';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import MatchTypeSelector from './MatchTypeSelector';
import SelectScanRoot from './SelectScanRoot';
import SortSelector from './SortSelector';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from './ui/collapsible';
import { ScrollArea } from './ui/scroll-area';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';

export default function Sidebar() {
  const searchInputRef = useRef<HTMLInputElement>(null);
  const resultsListRef = useRef<HTMLDivElement>(null);

  const fetchResults = useResultsStore((state) => state.fetchResults);
  const setSelectedResults = useResultsStore((state) => state.setSelectedResults);
  const selectResultRange = useResultsStore((state) => state.selectResultRange);
  const toggleResultSelection = useResultsStore((state) => state.toggleResultSelection);
  const pendingResults = useResultsStore((state) => state.pendingResults);
  const completedResults = useResultsStore((state) => state.completedResults);
  const setLastSelectedIndex = useResultsStore((state) => state.setLastSelectedIndex);
  const setLastSelectionType = useResultsStore((state) => state.setLastSelectionType);
  const moveToNextResult = useResultsStore((state) => state.moveToNextResult);
  const moveToPreviousResult = useResultsStore((state) => state.moveToPreviousResult);

  const filterByMatchType = useResultsStore((state) => state.filterByMatchType);
  const query = useResultsStore((state) => state.query);
  const debouncedQuery: string = useDebounce(query, DEBOUNCE_QUERY_MS);

  const handleSelectFiles = (e: React.MouseEvent, result: entities.ResultDTO, selectionType: 'pending' | 'completed') => {
    e.preventDefault();

    if (e.shiftKey) {
      selectResultRange(result, selectionType);
    } else if (e.metaKey || e.ctrlKey) {
      toggleResultSelection(result, selectionType);
    } else {
      handleSingleSelection(result, selectionType);
    }
  };

  const handleSingleSelection = (result: entities.ResultDTO, selectionType: 'pending' | 'completed') => {
    const resultsOfType = selectionType === 'pending' ? pendingResults : completedResults;
    const lastSelectedIndex = resultsOfType.findIndex((r) => r.path === result.path);

    setLastSelectionType(result.workflow_state as 'pending' | 'completed');
    setLastSelectedIndex(lastSelectedIndex);
    setSelectedResults([result]);
  };

  const moveFocusToResults = () => {
    if (resultsListRef.current) {
      resultsListRef.current.setAttribute('tabindex', '-1');
      resultsListRef.current.focus();
    }
  };

  useEffect(() => {
    fetchResults();
  }, [filterByMatchType, debouncedQuery]);

  useKeyboardShortcut(KEYBOARD_SHORTCUTS.moveUp.keys, moveToPreviousResult, {
    enableOnFormTags: false,
  });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.moveDown.keys, moveToNextResult, {
    enableOnFormTags: false,
  });
  useKeyboardShortcut('enter', moveFocusToResults, {
    enableOnFormTags: true,
  });

  return (
    <aside className="flex h-full flex-col border-r border-border bg-black/20 backdrop-blur-md">
      <div className="flex h-auto min-h-[65px] items-center border-b border-b-border px-4">
        <div className="break-all text-xs">
          {pendingResults?.length ? (
            <div className="flex flex-wrap items-center gap-1">
              <span className="font-semibold">
                {pendingResults?.length} decision{pendingResults.length > 1 ? 's' : ''} to make in
              </span>{' '}
              <SelectScanRoot />
            </div>
          ) : (
            'You have no decisions to make in working directory'
          )}
        </div>
      </div>

      <div className="flex flex-col gap-4 px-4 py-6">
        <div className="flex flex-col gap-2">
          <MatchTypeSelector />
          <SortSelector />
        </div>
        <ResultSearchBar searchInputRef={searchInputRef} />
      </div>

      <ScrollArea className="flex-1">
        <div className="flex flex-1 flex-col gap-2 outline-none" tabIndex={-1} ref={resultsListRef}>
          <ResultSection title="Pending files" results={pendingResults} onSelect={handleSelectFiles} selectionType="pending" />
          <ResultSection title="Completed files" results={completedResults} onSelect={handleSelectFiles} selectionType="completed" />
        </div>
      </ScrollArea>
    </aside>
  );
}

interface ResultSectionProps {
  title: string;
  results: entities.ResultDTO[];
  onSelect: (e: React.MouseEvent, result: entities.ResultDTO, selectionType: 'pending' | 'completed') => void;
  selectionType: 'pending' | 'completed';
}

function ResultSection({ title, results, onSelect, selectionType }: ResultSectionProps) {
  const parentRef = useRef<HTMLDivElement>(null);
  const ITEM_HEIGHT = 32; // Height of each result item in pixels

  const virtualizer = useVirtualizer({
    count: results.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => ITEM_HEIGHT,
    overscan: 5, // Number of items to render outside of the visible area
  });

  return (
    <Collapsible defaultOpen className="flex-1">
      <CollapsibleTrigger className="group w-full">
        <div className="flex items-center gap-1 border-b border-border px-3 pb-1">
          <ChevronRight className="group-data-[state=open]:stroke-text-foreground h-3 w-3 transform stroke-muted-foreground group-data-[state=open]:rotate-90" />
          <span className="text-sm text-muted-foreground">
            {title} <span className="text-xs">({results.length})</span>
          </span>
        </div>
      </CollapsibleTrigger>
      {results.length > 0 && (
        <CollapsibleContent className="flex flex-col gap-1 overflow-y-auto py-2">
          <div ref={parentRef} className="h-[400px] overflow-auto">
            <div
              style={{
                height: `${virtualizer.getTotalSize()}px`,
                width: '100%',
                position: 'relative',
              }}
            >
              {virtualizer.getVirtualItems().map((virtualRow) => (
                <div
                  key={virtualRow.key}
                  style={{
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    width: '100%',
                    height: `${virtualRow.size}px`,
                    transform: `translateY(${virtualRow.start}px)`,
                  }}
                >
                  <SidebarItem result={results[virtualRow.index]} onSelect={onSelect} selectionType={selectionType} />
                </div>
              ))}
            </div>
          </div>
        </CollapsibleContent>
      )}
    </Collapsible>
  );
}

interface SidebarItemProps {
  result: entities.ResultDTO;
  onSelect: (e: React.MouseEvent, result: entities.ResultDTO, selectionType: 'pending' | 'completed') => void;
  selectionType: 'pending' | 'completed';
}

function SidebarItem({ result, onSelect, selectionType }: SidebarItemProps) {
  const { selectedResults } = useResultsStore();

  const isSelected = selectedResults.some((r) => r.path === result.path);

  const fileName = getFileName(result.path);
  const directory = getDirectory(result.path);

  const presentation = stateInfoPresentation[result.filter_config?.action as FilterAction];

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <div
          className={clsx(
            'flex max-w-full cursor-pointer select-none items-center space-x-2 overflow-hidden px-4 py-1 text-sm text-muted-foreground',
            isSelected
              ? 'border-r-2 border-primary-foreground bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground'
              : 'hover:bg-primary/30'
          )}
          onClick={(e) => onSelect(e, result, selectionType)}
        >
          <span className="relative">
            {matchTypeIconMap[result.match_type as MatchType]}
            <span
              className={clsx('absolute bottom-0 right-0 h-1 w-1 rounded-full', presentation?.stateInfoSidebarIndicatorStyles ?? 'bg-transparent')}
            ></span>
          </span>
          <div className="flex min-w-0 items-center">
            {directory && <span className="truncate">{directory}</span>}
            <span className={clsx('whitespace-nowrap', !isSelected && 'text-foreground')}>{directory ? `/${fileName}` : fileName}</span>
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
