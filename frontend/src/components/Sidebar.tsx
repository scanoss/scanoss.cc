import clsx from 'clsx';
import { Braces, ChevronRight, File, Folder, FolderOpen } from 'lucide-react';
import { ReactNode, useEffect, useRef } from 'react';
import TreeView, { flattenTree, INode } from 'react-accessible-treeview';
import { IFlatMetadata } from 'react-accessible-treeview/dist/TreeView/utils';

import ResultSearchBar from '@/components/ResultSearchBar';
import useDebounce from '@/hooks/useDebounce';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import { getDirectory, getFileName } from '@/lib/utils';
import { FilterAction } from '@/modules/components/domain';
import { DEBOUNCE_QUERY_MS } from '@/modules/results/constants';
import { MatchType, stateInfoPresentation } from '@/modules/results/domain';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import { entities } from '../../wailsjs/go/models';
import MatchTypeSelector from './MatchTypeSelector';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from './ui/collapsible';
import { ScrollArea } from './ui/scroll-area';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';
import ViewTypeSelector from './ViewTypeSelector';

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

  const resultsTree = useResultsStore((state) => state.resultsTree);
  const fetchResultsTree = useResultsStore((state) => state.fetchResultsTree);

  const viewType = useResultsStore((state) => state.viewType);

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
    if (viewType === 'flat') {
      fetchResults();
    }
    if (viewType === 'tree') {
      fetchResultsTree();
    }
  }, [filterByMatchType, debouncedQuery, viewType]);

  useKeyboardShortcut(KEYBOARD_SHORTCUTS.moveUp.keys, moveToPreviousResult, {
    enableOnFormTags: false,
    enabled: viewType === 'flat',
  });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.moveDown.keys, moveToNextResult, {
    enableOnFormTags: false,
    enabled: viewType === 'flat',
  });
  useKeyboardShortcut('enter', moveFocusToResults, {
    enableOnFormTags: true,
    enabled: viewType === 'flat',
  });

  return (
    <aside className="flex h-full flex-col border-r border-border bg-black/20 backdrop-blur-md">
      <div className="flex h-16 items-center border-b border-b-border px-4">
        <h2 className="text-sm font-semibold">
          {pendingResults?.length
            ? `${pendingResults.length} decision${pendingResults.length > 1 ? 's' : ''} to make in working directory`
            : 'You have no decisions to make in working directory'}
        </h2>
      </div>

      <div className="flex flex-col gap-4 p-4">
        <MatchTypeSelector />
        <ResultSearchBar searchInputRef={searchInputRef} />
        <ViewTypeSelector />
      </div>

      <ScrollArea>
        {viewType === 'flat' ? (
          <div className="flex flex-1 flex-col gap-2 outline-none" tabIndex={-1} ref={resultsListRef}>
            <ResultSection title="Pending files" results={pendingResults} onSelect={handleSelectFiles} selectionType="pending" />
            <ResultSection title="Completed files" results={completedResults} onSelect={handleSelectFiles} selectionType="completed" />
          </div>
        ) : (
          <ResultTree tree={resultsTree} />
        )}
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
          {results.map((result) => (
            <SidebarItem key={result.path} result={result} onSelect={onSelect} selectionType={selectionType} />
          ))}
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

function ResultTree({ tree }: { tree: entities.ResultTreeDTO[] }) {
  if (!tree) {
    return null;
  }

  return (
    <div className="pb-6">
      {tree ? (
        <TreeView
          // @ts-expect-error TODO: fix this
          data={flattenTree(convertTreeStructure({ tree }))}
          aria-label="directory tree"
          togglableSelect
          clickAction="EXCLUSIVE_SELECT"
          multiSelect
          nodeRenderer={({ element, isBranch, getNodeProps, level, isExpanded, isSelected }) => {
            return (
              <div {...getNodeProps()}>
                <div
                  className={clsx(
                    'flex cursor-pointer items-center gap-2 py-1 text-sm text-muted-foreground',
                    isSelected ? 'bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground' : 'hover:bg-primary/30',
                    `pl-${4 * (level - 1)}`
                  )}
                >
                  <div className={clsx('flex', level === 1 && 'pl-1.5')}>
                    {isBranch ? <FolderIcon isOpen={isExpanded} /> : <FileIcon element={element} />}
                  </div>
                  <div className="flex min-w-0">
                    <span className="truncate">{element.name}</span>
                  </div>
                </div>
              </div>
            );
          }}
        />
      ) : null}
    </div>
  );
}

const FolderIcon = ({ isOpen }: { isOpen: boolean }) => (isOpen ? <FolderOpen className="h-3.5 w-3.5" /> : <Folder className="h-3.5 w-3.5" />);
const FileIcon = ({ element }: { element: INode<IFlatMetadata> }) => {
  // @ts-expect-error TODO: fix this
  const presentation = stateInfoPresentation[element.metadata?.filter_config?.action as FilterAction];

  return (
    <span className="relative">
      {matchTypeIconMap[element.metadata?.match_type as MatchType]}
      <span
        className={clsx('absolute bottom-0 right-0 h-1 w-1 rounded-full', presentation?.stateInfoSidebarIndicatorStyles ?? 'bg-transparent')}
      ></span>
    </span>
  );
};

export const matchTypeIconMap: Record<MatchType, ReactNode> = {
  [MatchType.File]: <File className="h-3.5 w-3.5 flex-shrink-0" />,
  [MatchType.Snippet]: <Braces className="h-3.5 w-3.5 flex-shrink-0" />,
};

function convertTreeStructure(inputData: { tree: entities.ResultTreeDTO[] }): OutputNode {
  function convertNode(node: entities.ResultTreeDTO): OutputNode {
    const result: OutputNode = {
      id: node.id,
      name: node.name,
      children: [],
      parent: node.parent ?? null,
      metadata: node.result,
    };

    if (node.children) {
      result.children = Object.values(node.children).map((child) => convertNode(child));

      if (result.children.length > 0) {
        result.children.sort((a, b) => a.name.localeCompare(b.name));
      }
    }

    return result;
  }

  const root: OutputNode = {
    id: '',
    name: '',
    children: inputData.tree.map((node) => convertNode(node)),
    parent: null,
    metadata: undefined,
  };

  if (root.children) {
    root.children.sort((a, b) => a.name.localeCompare(b.name));
  }

  return root;
}

interface OutputNode {
  id: string;
  name: string;
  children: OutputNode[];
  parent: string | null;
  metadata: entities.ResultDTO | undefined;
}
