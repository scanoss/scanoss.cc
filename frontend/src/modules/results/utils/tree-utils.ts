import { entities } from '../../../../wailsjs/go/models';

export interface TreeNode {
  id: string;
  name: string;
  path: string;
  children?: TreeNode[];
  isFolder?: boolean;
  data?: entities.ResultDTO;
  state?: 'pending' | 'completed' | 'mixed';
}

export function transformToTreeData(results: entities.ResultDTO[]): TreeNode[] {
  const root: { [key: string]: TreeNode } = {};

  // First pass: create all nodes
  results.forEach((result) => {
    const pathParts = result.path.split('/');

    // Create path nodes
    pathParts.forEach((part, index) => {
      const currentPath = pathParts.slice(0, index + 1).join('/');
      const isLastPart = index === pathParts.length - 1;

      if (!root[currentPath]) {
        root[currentPath] = {
          id: currentPath,
          name: part,
          path: currentPath,
          children: [],
          isFolder: !isLastPart,
          data: isLastPart ? result : undefined,
          state: isLastPart ? (result.workflow_state as 'pending' | 'completed') : undefined,
        };
      }
    });
  });

  // Second pass: build parent-child relationships and compute folder states
  Object.values(root).forEach((node) => {
    const parentPath = node.path.substring(0, node.path.lastIndexOf('/'));
    if (parentPath && root[parentPath]) {
      root[parentPath].children!.push(node);
    }
  });

  // Third pass: compute folder states bottom-up
  function computeFolderState(node: TreeNode): 'pending' | 'completed' | 'mixed' {
    if (!node.isFolder) {
      return node.state!;
    }

    const childStates = node.children!.map((child) => {
      if (!child.state) {
        child.state = computeFolderState(child);
      }
      return child.state;
    });

    const hasPending = childStates.includes('pending');
    const hasCompleted = childStates.includes('completed');

    if (hasPending && hasCompleted) {
      return 'mixed';
    }
    return hasPending ? 'pending' : 'completed';
  }

  // Get root level nodes and compute their states
  const rootNodes = Object.values(root).filter((node) => {
    const parentPath = node.path.substring(0, node.path.lastIndexOf('/'));
    return !parentPath || !root[parentPath];
  });

  rootNodes.forEach((node) => {
    if (node.isFolder && !node.state) {
      node.state = computeFolderState(node);
    }
  });

  return rootNodes;
}
