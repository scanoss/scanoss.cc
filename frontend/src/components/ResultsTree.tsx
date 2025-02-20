// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2025 SCANOSS.COM
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

import { useQuery } from '@tanstack/react-query';
import clsx from 'clsx';
import { ChevronDown, ChevronRight, File, Folder, Loader2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { NodeRendererProps, Tree } from 'react-arborist';

import { Card } from '@/components/ui/card';
import useConfigStore from '@/stores/useConfigStore';

import { GetTree } from '../../wailsjs/go/service/TreeServiceImpl';
import { ScrollArea } from './ui/scroll-area';
import { useToast } from './ui/use-toast';

interface TreeNode {
  id: string;
  name: string;
  path: string;
  isFolder: boolean;
  workflowState: string;
  children?: TreeNode[];
}

export default function ResultsTree() {
  const { toast } = useToast();
  const [selectedNode, setSelectedNode] = useState<string | null>(null);
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set());

  const scanRoot = useConfigStore((state) => state.scanRoot);

  const {
    data: tree,
    error,
    isLoading,
  } = useQuery({
    queryKey: ['tree', scanRoot],
    queryFn: () => GetTree(scanRoot),
    enabled: !!scanRoot,
  });

  useEffect(() => {
    if (error) {
      console.log(error);
      toast({
        title: 'Error loading file tree',
        description: error.message,
        variant: 'destructive',
      });
    }
  }, [error]);

  if (isLoading) {
    return (
      <div className="flex flex-1 items-center justify-center p-4">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error) {
    return (
      <Card className="flexflex-1 items-center justify-center p-4">
        <div className="text-center">
          <p className="text-sm font-medium text-destructive">Failed to load file tree</p>
          <p className="mt-1 text-sm text-muted-foreground">{error.message}</p>
        </div>
      </Card>
    );
  }

  const Node = ({ node, style, dragHandle, tree }: NodeRendererProps<TreeNode>) => {
    const isSelected = node.id === selectedNode;
    const isExpanded = expandedNodes.has(node.id);
    const hasChildren = node.data.children && node.data.children.length > 0;

    const handleClick = (e: React.MouseEvent) => {
      e.stopPropagation();
      setSelectedNode(node.id);

      if (hasChildren) {
        const newExpandedNodes = new Set(expandedNodes);
        if (isExpanded) {
          newExpandedNodes.delete(node.id);
        } else {
          newExpandedNodes.add(node.id);
        }
        setExpandedNodes(newExpandedNodes);
        tree.toggle(node.id);
      }
    };

    return (
      <div
        style={style}
        ref={dragHandle}
        className={clsx(
          'flex items-center gap-2 rounded-sm px-2 py-1',
          isSelected ? 'bg-primary text-primary-foreground' : 'hover:bg-accent',
          'cursor-pointer select-none'
        )}
        onClick={handleClick}
      >
        {hasChildren && (
          <div className="flex h-4 w-4 items-center justify-center">
            {isExpanded ? <ChevronDown className="h-3 w-3" /> : <ChevronRight className="h-3 w-3" />}
          </div>
        )}
        {node.data.isFolder ? <Folder className="h-4 w-4 shrink-0" /> : <File className="h-4 w-4 shrink-0" />}
        <span className="truncate text-sm">{node.data.name}</span>
      </div>
    );
  };

  return (
    <div className="flex flex-col gap-4 p-4">
      <ScrollArea className="flex-1">
        <Tree<TreeNode>
          data={tree}
          openByDefault={false}
          width={400}
          height={400}
          indent={24}
          rowHeight={32}
          overscanCount={5}
          paddingTop={16}
          paddingBottom={16}
          padding={16}
          disableDrag
          disableDrop
        >
          {Node}
        </Tree>
      </ScrollArea>
    </div>
  );
}
