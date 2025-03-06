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
import { Loader2, RefreshCcw } from 'lucide-react';
import { useCallback, useState } from 'react';
import { Tree } from 'react-arborist';

import { Card } from '@/components/ui/card';
import useConfigStore from '@/stores/useConfigStore';

import { GetTree } from '../../wailsjs/go/service/TreeServiceImpl';
import { Node, TreeNode } from './ResultsTreeNode';
import SaveButton from './SaveSkipSettingsButton';
import { Button } from './ui/button';

interface ResultsTreeProps {
  width: number | undefined;
  height: number | undefined;
}

export default function ResultsTree({ width, height }: ResultsTreeProps) {
  const [selectedNode, setSelectedNode] = useState<string | null>(null);
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set());
  const [loadingNodeId, setLoadingNodeId] = useState<string | null>(null);

  const scanRoot = useConfigStore((state) => state.scanRoot);

  const {
    data: tree,
    error,
    isLoading,
    refetch,
  } = useQuery({
    queryKey: ['resultsTree', scanRoot],
    queryFn: () => GetTree(scanRoot),
    enabled: !!scanRoot,
  });

  const handleNodeSelect = useCallback((nodeId: string) => {
    setSelectedNode(nodeId);
  }, []);

  const handleNodeToggle = useCallback((nodeId: string, isExpanded: boolean) => {
    setExpandedNodes((prev) => {
      const newExpandedNodes = new Set(prev);
      if (isExpanded) {
        newExpandedNodes.delete(nodeId);
      } else {
        newExpandedNodes.add(nodeId);
      }
      return newExpandedNodes;
    });
  }, []);

  if (isLoading) {
    return (
      <Card className="flex flex-1 items-center justify-center p-4 text-muted-foreground">
        <div className="flex items-center gap-2">
          <Loader2 className="h-3 w-3 animate-spin" />
          <p className="text-sm">Loading file tree...</p>
        </div>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className="flex flex-1 items-center justify-center p-4">
        <div className="text-center text-sm">
          <p className="text-destructive">Failed to load file tree</p>
          <p className="mt-1 text-muted-foreground">{error.message}</p>
          <Button variant="outline" size="sm" onClick={() => refetch()}>
            <RefreshCcw className="mr-2 h-3 w-3" />
            Retry
          </Button>
        </div>
      </Card>
    );
  }

  return (
    <div className="flex flex-col">
      <Tree<TreeNode> data={tree} openByDefault={false} width={width} height={height} disableDrag disableDrop>
        {(props) => (
          <Node
            {...props}
            selectedNode={selectedNode}
            expandedNodes={expandedNodes}
            onNodeSelect={handleNodeSelect}
            onNodeToggle={handleNodeToggle}
            loadingNodeId={loadingNodeId}
            setLoadingNodeId={setLoadingNodeId}
          />
        )}
      </Tree>

      <SaveButton />
    </div>
  );
}
