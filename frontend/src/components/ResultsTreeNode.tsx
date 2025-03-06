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

import { useMutation } from '@tanstack/react-query';
import clsx from 'clsx';
import { CheckCircle, ChevronDown, ChevronRight, File, FileType, Folder, Loader2, XCircle } from 'lucide-react';
import path from 'path-browserify';
import { memo } from 'react';
import { NodeRendererProps } from 'react-arborist';

import useSkipPatternStore from '@/hooks/useSkipPatternStore';

import { AddStagedScanningSkipPattern, RemoveStagedScanningSkipPattern } from '../../wailsjs/go/service/ScanossSettingsServiceImp';
import { ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuSeparator, ContextMenuTrigger } from './ui/context-menu';
import { useToast } from './ui/use-toast';

const validExtensionRegex = /^\.[a-zA-Z0-9]+$/;
export const SKIP_STATES = {
  EXCLUDED: 'excluded',
  INCLUDED: 'included',
  MIXED: 'mixed',
};

export const SKIP_STATE_DESCRIPTIONS = {
  [SKIP_STATES.EXCLUDED]: 'Excluded from scanning',
  [SKIP_STATES.INCLUDED]: 'Included in scanning',
  [SKIP_STATES.MIXED]: 'Partially included in scanning',
};

export interface TreeNode {
  id: string;
  name: string;
  path: string;
  isFolder: boolean;
  workflowState: string;
  scanningSkipState: string;
  children?: TreeNode[];
}

interface NodeProps extends NodeRendererProps<TreeNode> {
  selectedNode: string | null;
  expandedNodes: Set<string>;
  onNodeSelect: (nodeId: string) => void;
  onNodeToggle: (nodeId: string, isExpanded: boolean) => void;
  loadingNodeId: string | null;
  setLoadingNodeId: (nodeId: string | null) => void;
}

export const Node = memo(
  ({ node, style, tree, selectedNode, expandedNodes, onNodeSelect, onNodeToggle, loadingNodeId, setLoadingNodeId }: NodeProps) => {
    const { toast } = useToast();
    const isExpanded = expandedNodes.has(node.id);
    const hasChildren = node.data.children && node.data.children.length > 0;
    const isLoading = loadingNodeId === node.id;
    const isSelected = node.id === selectedNode;
    const { setHasUnsavedChanges, addNodeWithUnsavedChanges, nodesWithUnsavedChanges } = useSkipPatternStore();

    const hasUnsavedChanges = nodesWithUnsavedChanges.has(node.id);

    const handleClickNode = (e: React.MouseEvent) => {
      e.stopPropagation();
      onNodeSelect(node.id);

      if (hasChildren) {
        onNodeToggle(node.id, isExpanded);
        tree.toggle(node.id);
      }
    };

    const handleToggleSkipFile = () => {
      if (loadingNodeId !== null) return;
      const pattern = node.data.path;

      if (node.data.scanningSkipState === SKIP_STATES.EXCLUDED) {
        removeScanningSkipPattern(pattern);
      } else {
        addScanningSkipPattern(pattern);
      }
    };

    const handleToggleSkipFolder = () => {
      if (loadingNodeId !== null) return;

      const pattern = `${node.data.path}/`;

      if (node.data.scanningSkipState === SKIP_STATES.EXCLUDED) {
        removeScanningSkipPattern(pattern);
      } else {
        addScanningSkipPattern(pattern);
      }
    };

    const handleToggleSkipExtension = () => {
      if (loadingNodeId !== null) return;

      const extension = path.extname(node.data.path);
      if (!validExtensionRegex.test(extension)) {
        toast({
          title: 'Invalid Extension',
          description: `Cannot skip files with extension "${extension}"`,
          variant: 'destructive',
        });
        return;
      }

      const pattern = `*${extension}`;

      if (node.data.scanningSkipState === SKIP_STATES.EXCLUDED) {
        removeScanningSkipPattern(pattern);
      } else {
        addScanningSkipPattern(pattern);
      }

      addNodeWithUnsavedChanges(node.id);
    };

    const { mutate: addScanningSkipPattern } = useMutation({
      mutationFn: async (pattern: string) => {
        setLoadingNodeId(node.id);
        await AddStagedScanningSkipPattern(pattern);

        setHasUnsavedChanges(true);
        addNodeWithUnsavedChanges(node.id);
      },
      onSuccess: (_, pattern) => {
        toast({
          title: 'Success',
          description: `Added "${pattern}" to scanning skip patterns`,
        });
        setLoadingNodeId(null);
      },
      onError: (error, pattern) => {
        console.error(`Error adding skip pattern: ${pattern}`, error);
        setLoadingNodeId(null);
      },
    });

    const { mutate: removeScanningSkipPattern } = useMutation({
      mutationFn: async (pattern: string) => {
        setLoadingNodeId(node.id);
        await RemoveStagedScanningSkipPattern(pattern);

        setHasUnsavedChanges(true);
        addNodeWithUnsavedChanges(node.id);
      },
      onSuccess: (_, pattern) => {
        toast({
          title: 'Success',
          description: `Removed "${pattern}" from scanning skip patterns`,
        });
        setLoadingNodeId(null);
      },
      onError: (error, pattern) => {
        console.error(`Error removing skip pattern: ${pattern}`, error);
        setLoadingNodeId(null);
      },
    });

    const skipMenuItemText = () => {
      return node.data.scanningSkipState === SKIP_STATES.EXCLUDED
        ? `Include ${node.data.isFolder ? 'folder' : 'file'} in scanning`
        : `Exclude ${node.data.isFolder ? 'folder' : 'file'} from scanning`;
    };

    const skipExtensionMenuItemText = () => {
      const extension = path.extname(node.data.path);
      return node.data.scanningSkipState === SKIP_STATES.EXCLUDED
        ? `Include all ${extension} files in scanning`
        : `Exclude all ${extension} files from scanning`;
    };

    return (
      <ContextMenu>
        <ContextMenuTrigger>
          <div
            className={clsx(
              'flex cursor-pointer items-center gap-1 rounded-md px-2 py-1 hover:bg-accent',
              isSelected && 'bg-accent',
              hasUnsavedChanges && 'border-l-2 border-yellow-500'
            )}
            style={style}
            onClick={handleClickNode}
          >
            {hasChildren && (
              <div className="flex h-4 w-4 items-center justify-center">
                {isExpanded ? <ChevronDown className="h-3 w-3" /> : <ChevronRight className="h-3 w-3" />}
              </div>
            )}
            <div className="relative">
              {node.data.isFolder ? <Folder className="h-4 w-4 shrink-0" /> : <File className="h-4 w-4 shrink-0" />}
              <div
                className={clsx(
                  'absolute bottom-0 right-0 h-1 w-1 rounded-full',
                  node.data.scanningSkipState === SKIP_STATES.EXCLUDED && 'bg-red-500',
                  node.data.scanningSkipState === SKIP_STATES.INCLUDED && 'bg-green-500',
                  node.data.scanningSkipState === SKIP_STATES.MIXED && 'bg-yellow-500'
                )}
                title={`Scanning status: ${node.data.scanningSkipState}`}
              ></div>
              {hasUnsavedChanges && <div className="absolute -right-1 -top-1 h-2 w-2 rounded-full bg-yellow-500" title="Has unsaved changes"></div>}
            </div>
            <span className="truncate text-sm">{node.data.name}</span>
            {isLoading && <Loader2 className="ml-1 h-3 w-3 animate-spin" />}
          </div>
        </ContextMenuTrigger>
        <ContextMenuContent>
          <ContextMenuItem onClick={node.data.isFolder ? handleToggleSkipFolder : handleToggleSkipFile} asChild>
            <div className="flex items-center gap-2">
              {node.data.scanningSkipState === SKIP_STATES.EXCLUDED ? <CheckCircle className="h-4 w-4" /> : <XCircle className="h-4 w-4" />}
              <span>{skipMenuItemText()}</span>
            </div>
          </ContextMenuItem>
          {!node.data.isFolder && (
            <>
              <ContextMenuSeparator />
              <ContextMenuItem onClick={handleToggleSkipExtension} asChild>
                <div className="flex items-center gap-2">
                  <FileType className="h-4 w-4" />
                  <span>{skipExtensionMenuItemText()}</span>
                </div>
              </ContextMenuItem>
            </>
          )}
        </ContextMenuContent>
      </ContextMenu>
    );
  }
);

Node.displayName = 'TreeNode';
