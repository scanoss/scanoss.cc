import { Braces, ChevronRight, File, Folder } from 'lucide-react';
import React from 'react';
import { NodeApi, Tree } from 'react-arborist';

import { cn } from '@/lib/utils';
import useTreeStore, { TreeNode, TreeNodeState } from '@/modules/results/stores/useTreeStore';

import { entities } from '../../wailsjs/go/models';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';

interface ResultsTreeViewProps {
  onSelect: (e: React.MouseEvent, result: entities.ResultDTO, selectionType: 'pending' | 'completed') => void;
  selectedResults: entities.ResultDTO[];
}

export default function ResultsTreeView({ onSelect, selectedResults }: ResultsTreeViewProps) {
  const { nodes, isLoading, childrenCache } = useTreeStore();

  const treeData = React.useMemo(() => {
    return nodes.map((node) => ({
      ...node,
      children: childrenCache[node.path] || [],
    }));
  }, [nodes, childrenCache]);

  if (isLoading) {
    return <div className="p-4 text-sm text-muted-foreground">Loading...</div>;
  }

  return (
    <Tree<TreeNode>
      data={treeData}
      width="100%"
      height={500}
      indent={24}
      rowHeight={28}
      overscanCount={10}
      paddingTop={8}
      paddingBottom={8}
      openByDefault={false}
    >
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
  const getChildren = useTreeStore((state) => state.getChildren);
  const data = node.data;
  const isSelected = data.data ? selectedResults.some((r) => r.path === data.data?.path) : false;
  const [isLoadingChildren, setIsLoadingChildren] = React.useState(false);

  const handleClick = async (e: React.MouseEvent) => {
    if (data.isFolder) {
      if (!node.isOpen) {
        if (!data.children?.length) {
          setIsLoadingChildren(true);
          try {
            await getChildren(data.path);
          } finally {
            setIsLoadingChildren(false);
          }
        }
        node.open();
      } else {
        node.close();
      }
    } else if (data.data) {
      onSelect(e, data.data, data.state as 'pending' | 'completed');
    }
  };

  const handleChevronClick = async (e: React.MouseEvent) => {
    e.stopPropagation();
    if (data.isFolder) {
      if (!node.isOpen) {
        if (!data.children?.length) {
          setIsLoadingChildren(true);
          try {
            await getChildren(data.path);
          } finally {
            setIsLoadingChildren(false);
          }
        }
        node.open();
      } else {
        node.close();
      }
    }
  };

  const getStateIndicatorStyle = (state: TreeNodeState) => {
    switch (state) {
      case 'pending':
        return 'bg-yellow-500';
      case 'completed':
        return 'bg-green-500';
      case 'mixed':
        return 'bg-orange-500';
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
            <ChevronRight
              className={cn('h-3 w-3 transform transition-transform', node.isOpen && 'rotate-90', isLoadingChildren && 'animate-spin')}
              onClick={handleChevronClick}
            />
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
          {data.isFolder && data.childCount > 0 && <span className="text-xs text-muted-foreground">({data.childCount})</span>}
        </div>
      </TooltipTrigger>
      <TooltipContent side="right" sideOffset={15}>
        {data.path}
      </TooltipContent>
    </Tooltip>
  );
}
