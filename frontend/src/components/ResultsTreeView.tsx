import { Braces, ChevronRight, File, Folder } from 'lucide-react';
import { NodeApi, Tree } from 'react-arborist';

import { cn } from '@/lib/utils';
import { transformToTreeData, TreeNode } from '@/modules/results/utils/tree-utils';

import { entities } from '../../wailsjs/go/models';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';

interface ResultsTreeViewProps {
  results: entities.ResultDTO[];
  onSelect: (e: React.MouseEvent, result: entities.ResultDTO, selectionType: 'pending' | 'completed') => void;
  selectedResults: entities.ResultDTO[];
}

export default function ResultsTreeView({ results, onSelect, selectedResults }: ResultsTreeViewProps) {
  const treeData = transformToTreeData(results);

  return (
    <Tree data={treeData} width="100%" height={500} indent={24} rowHeight={28} overscanCount={10} paddingTop={8} paddingBottom={8}>
      {(props) => <TreeNodeComponent {...props} onSelect={onSelect} selectedResults={selectedResults} />}
    </Tree>
  );
}

interface TreeNodeComponentProps {
  node: NodeApi<TreeNode>;
  style: React.CSSProperties;
  dragHandle?: React.Ref<HTMLDivElement>;
  onSelect: (e: React.MouseEvent, result: entities.ResultDTO, selectionType: 'pending' | 'completed') => void;
  selectedResults: entities.ResultDTO[];
}

function TreeNodeComponent({ node, style, dragHandle, onSelect, selectedResults }: TreeNodeComponentProps) {
  const data = node.data;
  const isSelected = data.data ? selectedResults.some((r) => r.path === data.data?.path) : false;

  const handleClick = (e: React.MouseEvent) => {
    if (data.data) {
      onSelect(e, data.data, data.state as 'pending' | 'completed');
    } else if (data.isFolder) {
      toggleNode();
    }
  };

  const handleChevronClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    toggleNode();
  };

  const toggleNode = () => {
    if (node.isOpen) {
      node.close();
    } else {
      node.open();
    }
  };

  const getStateIndicatorStyle = (state?: 'pending' | 'completed' | 'mixed') => {
    switch (state) {
      case 'pending':
        return 'bg-yellow-500';
      case 'completed':
        return 'bg-green-500';
      case 'mixed':
        return 'bg-orange-500';
      default:
        return 'bg-transparent';
    }
  };

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <div
          ref={dragHandle}
          style={style}
          className={cn(
            'flex items-center space-x-2 overflow-hidden px-2 py-1 text-sm text-muted-foreground',
            'cursor-pointer select-none',
            isSelected
              ? 'border-r-2 border-primary-foreground bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground'
              : 'hover:bg-primary/30'
          )}
          onClick={handleClick}
        >
          {data.isFolder && (
            <ChevronRight className={cn('h-3 w-3 transform transition-transform', node.isOpen && 'rotate-90')} onClick={handleChevronClick} />
          )}
          <span className="relative">
            {data.isFolder ? (
              <Folder className="h-3 w-3 flex-shrink-0" />
            ) : data.data?.match_type === 'snippet' ? (
              <Braces className="h-3 w-3 flex-shrink-0" />
            ) : (
              <File className="h-3 w-3 flex-shrink-0" />
            )}
            <span className={cn('absolute bottom-0 right-0 h-1 w-1 rounded-full', getStateIndicatorStyle(data.state))} />
          </span>
          <span className="truncate">{data.name}</span>
        </div>
      </TooltipTrigger>
      <TooltipContent side="right" sideOffset={15}>
        {data.path}
      </TooltipContent>
    </Tooltip>
  );
}
