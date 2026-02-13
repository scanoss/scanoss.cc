// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2026 SCANOSS.COM
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

import clsx from 'clsx';
import { ChevronDown, ChevronRight, File, Folder } from 'lucide-react';
import { useCallback, useMemo, useState } from 'react';
import { entities } from 'wailsjs/go/models';

import { FilterAction } from '@/modules/components/domain';
import { stateInfoPresentation } from '@/modules/results/domain';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import { ScrollArea } from './ui/scroll-area';

interface ResultTreeNode {
  name: string;
  path: string;
  isFolder: boolean;
  children: ResultTreeNode[];
  result?: entities.ResultDTO;
}

function buildTree(results: entities.ResultDTO[]): ResultTreeNode[] {
  const root: ResultTreeNode = { name: '', path: '', isFolder: true, children: [] };

  for (const result of results) {
    const parts = result.path.split('/').filter(Boolean);
    let current = root;

    for (let i = 0; i < parts.length; i++) {
      const part = parts[i];
      const isLast = i === parts.length - 1;
      const currentPath = parts.slice(0, i + 1).join('/');

      let child = current.children.find((c) => c.name === part);
      if (!child) {
        child = {
          name: part,
          path: currentPath,
          isFolder: !isLast,
          children: [],
          result: isLast ? result : undefined,
        };
        current.children.push(child);
      }
      current = child;
    }
  }

  // Sort: folders first, then alphabetically
  const sortTree = (nodes: ResultTreeNode[]) => {
    nodes.sort((a, b) => {
      if (a.isFolder && !b.isFolder) return -1;
      if (!a.isFolder && b.isFolder) return 1;
      return a.name.localeCompare(b.name);
    });
    for (const node of nodes) {
      if (node.children.length > 0) sortTree(node.children);
    }
  };
  sortTree(root.children);

  return root.children;
}

interface SidebarTreeViewProps {
  pendingResults: entities.ResultDTO[];
  completedResults: entities.ResultDTO[];
  onSelect: (e: React.MouseEvent, result: entities.ResultDTO, selectionType: 'pending' | 'completed') => void;
}

export default function SidebarTreeView({ pendingResults, completedResults, onSelect }: SidebarTreeViewProps) {
  const tree = useMemo(
    () => buildTree([...pendingResults, ...completedResults]),
    [pendingResults, completedResults]
  );

  return (
    <div className="flex h-full flex-col">
      <div className="flex flex-shrink-0 items-center gap-1 border-b border-border px-3 pb-1">
        <span className="text-sm text-muted-foreground">
          Results tree {pendingResults.length + completedResults.length > 0 && (
            <span className="text-xs">({pendingResults.length + completedResults.length})</span>
          )}
        </span>
      </div>
      <div className="min-h-0 flex-1">
        <ScrollArea className="h-full">
          <div className="py-1">
            {tree.map((node) => (
              <TreeNodeComponent key={node.path} node={node} depth={0} onSelect={onSelect} />
            ))}
          </div>
        </ScrollArea>
      </div>
    </div>
  );
}

interface TreeNodeComponentProps {
  node: ResultTreeNode;
  depth: number;
  onSelect: (e: React.MouseEvent, result: entities.ResultDTO, selectionType: 'pending' | 'completed') => void;
}

function TreeNodeComponent({ node, depth, onSelect }: TreeNodeComponentProps) {
  const [isExpanded, setIsExpanded] = useState(true);
  const selectedResults = useResultsStore((state) => state.selectedResults);

  const isSelected = node.result ? selectedResults.some((r) => r.path === node.result!.path) : false;
  const presentation = node.result
    ? stateInfoPresentation[node.result.filter_config?.action as FilterAction]
    : undefined;

  const handleClick = useCallback(
    (e: React.MouseEvent) => {
      e.stopPropagation();
      if (node.isFolder) {
        setIsExpanded((prev) => !prev);
      } else if (node.result) {
        const selectionType = (node.result.workflow_state as 'pending' | 'completed') ?? 'pending';
        onSelect(e, node.result, selectionType);
      }
    },
    [node, onSelect]
  );

  return (
    <>
      <div
        className={clsx(
          'flex cursor-pointer select-none items-center gap-1 py-0.5 pr-2 text-sm text-muted-foreground hover:bg-primary/30',
          isSelected && 'border-r-2 border-primary-foreground bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground'
        )}
        style={{ paddingLeft: `${depth * 12 + 8}px` }}
        onClick={handleClick}
      >
        {node.isFolder ? (
          <>
            <span className="flex h-4 w-4 items-center justify-center">
              {isExpanded ? <ChevronDown className="h-3 w-3" /> : <ChevronRight className="h-3 w-3" />}
            </span>
            <Folder className="h-3 w-3 shrink-0" />
          </>
        ) : (
          <>
            <span className="h-4 w-4" />
            <span className="relative">
              <File className="h-3 w-3 shrink-0" />
              <span
                className={clsx(
                  'absolute bottom-0 right-0 h-1 w-1 rounded-full',
                  presentation?.stateInfoSidebarIndicatorStyles ?? 'bg-transparent'
                )}
              />
            </span>
          </>
        )}
        <span className={clsx('truncate', !node.isFolder && !isSelected && 'text-foreground')}>{node.name}</span>
      </div>
      {node.isFolder && isExpanded && node.children.map((child) => (
        <TreeNodeComponent key={child.path} node={child} depth={depth + 1} onSelect={onSelect} />
      ))}
    </>
  );
}
