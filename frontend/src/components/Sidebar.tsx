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
import { Braces, File } from 'lucide-react';
import { ReactNode, useRef } from 'react';
import { entities } from 'wailsjs/go/models';

import ResultSearchBar from '@/components/ResultSearchBar';
import { useDialogState } from '@/contexts/DialogStateContext';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { useResults } from '@/hooks/useResults';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import { getDirectory, getFileName } from '@/lib/utils';
import { FilterAction } from '@/modules/components/domain';
import { MatchType, stateInfoPresentation } from '@/modules/results/domain';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import Loading from './Loading';
import MatchTypeSelector from './MatchTypeSelector';
import SelectScanRoot from './SelectScanRoot';
import SortSelector from './SortSelector';
import { ScrollArea } from './ui/scroll-area';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';

export default function Sidebar() {
  const searchInputRef = useRef<HTMLInputElement>(null);
  const resultsListRef = useRef<HTMLDivElement>(null);
  const { isKeyboardNavigationBlocked } = useDialogState();

  const setSelectedResults = useResultsStore((state) => state.setSelectedResults);
  const selectResultRange = useResultsStore((state) => state.selectResultRange);
  const toggleResultSelection = useResultsStore((state) => state.toggleResultSelection);
  const setLastSelectedIndex = useResultsStore((state) => state.setLastSelectedIndex);
  const setLastSelectionType = useResultsStore((state) => state.setLastSelectionType);
  const moveToNextResult = useResultsStore((state) => state.moveToNextResult);
  const moveToPreviousResult = useResultsStore((state) => state.moveToPreviousResult);

  const { data: results, isLoading: isLoadingResults } = useResults();

  const { pendingResults = [], completedResults = [] } = results ?? {};

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

  const moveFocusToResults = () => {
    if (resultsListRef.current) {
      resultsListRef.current.setAttribute('tabindex', '-1');
      resultsListRef.current.focus();
    }
  };

  const handleSingleSelection = (result: entities.ResultDTO, selectionType: 'pending' | 'completed') => {
    const resultsOfType = selectionType === 'pending' ? pendingResults : completedResults;
    const lastSelectedIndex = resultsOfType.findIndex((r) => r.path === result.path);

    setLastSelectionType(result.workflow_state as 'pending' | 'completed');
    setLastSelectedIndex(lastSelectedIndex);
    setSelectedResults([result]);
  };

  useKeyboardShortcut(KEYBOARD_SHORTCUTS.moveUp.keys, moveToPreviousResult, {
    enableOnFormTags: false,
    enabled: !isKeyboardNavigationBlocked,
  });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.moveDown.keys, () => moveToNextResult(), {
    enableOnFormTags: false,
    enabled: !isKeyboardNavigationBlocked,
  });
  useKeyboardShortcut('enter', moveFocusToResults, {
    enableOnFormTags: true,
    enabled: !isKeyboardNavigationBlocked,
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
          ) : isLoadingResults ? (
            <Loading text="Loading results..." size="w-3 h-3" />
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

      <div className="min-h-0 flex-1">
        <div className="flex h-full min-h-0 flex-1 flex-col outline-none" tabIndex={-1} ref={resultsListRef}>
          <div className="flex min-h-0 flex-1 flex-col gap-2">
            <ResultSection
              title="Pending files"
              results={pendingResults}
              onSelect={handleSelectFiles}
              selectionType="pending"
              isLoading={isLoadingResults}
            />
            <ResultSection
              title="Completed files"
              results={completedResults}
              onSelect={handleSelectFiles}
              selectionType="completed"
              isLoading={isLoadingResults}
            />
          </div>
        </div>
      </div>
    </aside>
  );
}

interface ResultSectionProps {
  title: string;
  results: entities.ResultDTO[];
  onSelect: (e: React.MouseEvent, result: entities.ResultDTO, selectionType: 'pending' | 'completed') => void;
  selectionType: 'pending' | 'completed';
  isLoading: boolean;
}

function ResultSection({ title, results, onSelect, selectionType, isLoading }: ResultSectionProps) {
  const parentRef = useRef<HTMLDivElement>(null);
  const ITEM_HEIGHT = 32;

  const virtualizer = useVirtualizer({
    count: results.length,
    getScrollElement: () => parentRef.current?.querySelector('[data-radix-scroll-area-viewport]') as Element,
    estimateSize: () => ITEM_HEIGHT,
    overscan: 5,
  });

  return (
    <div className="flex h-1/2 flex-col">
      <div className="flex flex-shrink-0 items-center gap-1 border-b border-border px-3 pb-1">
        <span className="text-sm text-muted-foreground">
          {title} {results.length > 0 ? <span className="text-xs">({results.length})</span> : null}
        </span>
      </div>
      <div className="min-h-0 flex-1" ref={parentRef}>
        <ScrollArea className="h-full">
          {isLoading ? (
            <div className="flex h-full w-full items-center justify-center">
              <Loading text="Loading results..." size="w-3 h-3" />
            </div>
          ) : (
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
          )}
        </ScrollArea>
      </div>
    </div>
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
