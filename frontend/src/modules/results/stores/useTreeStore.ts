import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

import { entities } from '../../../../wailsjs/go/models';
import { GetChildren, GetNode, GetRootNodes, SearchNodes } from '../../../../wailsjs/go/service/TreeServiceImpl';

interface TreeState {
  nodes: TreeNode[];
  childrenCache: Record<string, TreeNode[]>;
  totalCount: number;
  isLoading: boolean;
  error: string | null;
  expandedPaths: Set<string>;
}

export type TreeNodeState = 'pending' | 'completed' | 'mixed';

export interface TreeNode {
  id: string;
  name: string;
  path: string;
  isFolder: boolean;
  state: TreeNodeState;
  data?: entities.ResultDTO;
  childCount: number;
  children?: TreeNode[];
}

interface TreeActions {
  getChildren: (parentPath: string, offset?: number, limit?: number) => Promise<void>;
  getRootNodes: (offset?: number, limit?: number) => Promise<void>;
  searchNodes: (query: string, offset?: number, limit?: number) => Promise<void>;
  getNode: (path: string) => Promise<TreeNode | null>;
  setExpandedPath: (path: string, isExpanded: boolean) => void;
}

type TreeStore = TreeState & TreeActions;

const DEFAULT_PAGE_SIZE = 100;

const useTreeStore = create<TreeStore>()(
  devtools((set) => ({
    nodes: [],
    childrenCache: {},
    totalCount: 0,
    isLoading: false,
    error: null,
    expandedPaths: new Set(),

    getChildren: async (parentPath: string, offset = 0, limit = DEFAULT_PAGE_SIZE) => {
      try {
        set({ isLoading: true, error: null });
        const response = await GetChildren(parentPath, offset, limit);
        const children = response.nodes as TreeNode[];

        set((state) => ({
          childrenCache: {
            ...state.childrenCache,
            [parentPath]: children,
          },
        }));
      } catch (error) {
        set({ error: error instanceof Error ? error.message : 'Failed to load children' });
      } finally {
        set({ isLoading: false });
      }
    },

    getRootNodes: async (offset = 0, limit = DEFAULT_PAGE_SIZE) => {
      try {
        set({ isLoading: true, error: null });
        const response = await GetRootNodes(offset, limit);
        set({
          nodes: response.nodes as TreeNode[],
          totalCount: response.totalCount,
          childrenCache: {}, // Reset cache when loading root nodes
        });
      } catch (error) {
        set({ error: error instanceof Error ? error.message : 'Failed to load root nodes' });
      } finally {
        set({ isLoading: false });
      }
    },

    searchNodes: async (query: string, offset = 0, limit = DEFAULT_PAGE_SIZE) => {
      try {
        set({ isLoading: true, error: null });
        const response = await SearchNodes(query, offset, limit);
        set({
          nodes: response.nodes as TreeNode[],
          totalCount: response.totalCount,
          childrenCache: {}, // Reset cache when searching
        });
      } catch (error) {
        set({ error: error instanceof Error ? error.message : 'Failed to search nodes' });
      } finally {
        set({ isLoading: false });
      }
    },

    getNode: async (path: string) => {
      try {
        set({ isLoading: true, error: null });
        const node = await GetNode(path);
        return node as TreeNode;
      } catch (error) {
        set({ error: error instanceof Error ? error.message : 'Failed to get node' });
        return null;
      } finally {
        set({ isLoading: false });
      }
    },

    setExpandedPath: (path: string, isExpanded: boolean) => {
      set((state) => {
        const newExpandedPaths = new Set(state.expandedPaths);
        if (isExpanded) {
          newExpandedPaths.add(path);
        } else {
          newExpandedPaths.delete(path);
        }
        return { expandedPaths: newExpandedPaths };
      });
    },
  }))
);

export default useTreeStore;
