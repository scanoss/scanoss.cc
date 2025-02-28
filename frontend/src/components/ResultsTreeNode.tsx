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

import { useMutation, useQueryClient } from '@tanstack/react-query';
import clsx from 'clsx';
import { CheckCircle, ChevronDown, ChevronRight, File, FileType, Folder, Loader2, XCircle } from 'lucide-react';
import path from 'path-browserify';
import { memo } from 'react';
import { NodeRendererProps } from 'react-arborist';

import useConfigStore from '@/stores/useConfigStore';

import { AddScanningSkipPattern, RemoveScanningSkipPattern } from '../../wailsjs/go/service/ScanossSettingsServiceImp';
import { ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuSeparator, ContextMenuTrigger } from './ui/context-menu';
import { useToast } from './ui/use-toast';

const validExtensionRegex = /^\.[a-zA-Z0-9]+$/;
const SKIP_STATES = {
  EXCLUDED: 'excluded',
  INCLUDED: 'included',
  MIXED: 'mixed',
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
    const scanRoot = useConfigStore((state) => state.scanRoot);
    const queryClient = useQueryClient();

    const isSelected = node.id === selectedNode;
    const isExpanded = expandedNodes.has(node.id);
    const hasChildren = node.data.children && node.data.children.length > 0;
    const isLoading = loadingNodeId === node.id;

    const handleClickNode = (e: React.MouseEvent) => {
      e.stopPropagation();
      onNodeSelect(node.id);

      if (hasChildren) {
        onNodeToggle(node.id, isExpanded);
        tree.toggle(node.id);
      }
    };

    const { mutate: addScanningSkipPattern } = useMutation({
      mutationFn: async (pattern: string) => {
        setLoadingNodeId(node.id);

        await AddScanningSkipPattern(pattern);
      },
      onSuccess: (_, pattern) => {
        toast({
          title: 'Success',
          description: `Added "${pattern}" to scanning skip patterns`,
        });
        queryClient.invalidateQueries({ queryKey: ['resultsTree', scanRoot] });
        setLoadingNodeId(null);
      },
      onError: (error, pattern) => {
        console.error(`Error adding/saving skip pattern: ${pattern}`, error);
        setLoadingNodeId(null);
      },
    });

    const { mutate: removeScanningSkipPattern } = useMutation({
      mutationFn: async (pattern: string) => {
        setLoadingNodeId(node.id);

        await RemoveScanningSkipPattern(pattern);
      },
      onSuccess: (_, pattern) => {
        toast({
          title: 'Success',
          description: `Removed "${pattern}" from scanning skip patterns`,
        });
        queryClient.invalidateQueries({ queryKey: ['resultsTree', scanRoot] });
        setLoadingNodeId(null);
      },
      onError: (error, pattern) => {
        console.error(`Error removing skip pattern: ${pattern}`, error);

        setLoadingNodeId(null);
      },
    });

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
      if (!extension || !validExtensionRegex.test(extension)) {
        toast({
          title: 'Warning',
          description: 'Invalid or missing file extension',
          variant: 'default',
        });
        return;
      }

      const pattern = `**/*${extension}`;

      if (node.data.scanningSkipState === SKIP_STATES.EXCLUDED) {
        removeScanningSkipPattern(pattern);
      } else {
        addScanningSkipPattern(pattern);
      }
    };

    const skipMenuItemText = () => {
      if (node.data.scanningSkipState === SKIP_STATES.EXCLUDED) {
        return `Include ${node.data.isFolder ? 'folder' : 'file'} for scanning`;
      }

      return `Skip ${node.data.isFolder ? 'folder' : 'file'} for scanning`;
    };

    const skipExtensionMenuItemText = () => {
      if (node.data.scanningSkipState === SKIP_STATES.EXCLUDED) {
        return `Include all files with this extension for scanning`;
      }

      return `Skip all files with this extension`;
    };

    return (
      <ContextMenu>
        <ContextMenuTrigger>
          <div
            style={style}
            className={clsx(
              'relative flex cursor-pointer select-none items-center gap-2 rounded-sm px-2 py-1',
              isSelected ? 'bg-primary text-primary-foreground' : 'hover:bg-accent'
            )}
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
