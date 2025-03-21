import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

export interface TreeNode {
  id: string;
  name: string;
  path: string;
  isFolder: boolean;
  workflowState: string;
  scanningSkipState: string;
  children?: TreeNode[];
}

interface TreeState {
  loadingNodeId: string | null;
  selectedNode: string | null;
  expandedNodes: Set<string>;
}

interface TreeActions {
  onNodeSelect: (nodeId: string) => void;
  onNodeToggle: (nodeId: string, isExpanded: boolean) => void;
  setLoadingNodeId: (nodeId: string | null) => void;
}

type TreeStore = TreeState & TreeActions;

export default create<TreeStore>()(
  devtools((set) => ({
    loadingNodeId: null,
    selectedNode: null,
    expandedNodes: new Set(),

    onNodeSelect: (nodeId: string) => {
      set({ selectedNode: nodeId });
    },
    onNodeToggle: (nodeId: string, isExpanded: boolean) => {
      set((state) => {
        const newExpandedNodes = new Set(state.expandedNodes);
        if (isExpanded) {
          newExpandedNodes.delete(nodeId);
        } else {
          newExpandedNodes.add(nodeId);
        }
        return { expandedNodes: newExpandedNodes };
      });
    },
    setLoadingNodeId: (nodeId: string | null) => {
      set({ loadingNodeId: nodeId });
    },
  }))
);
