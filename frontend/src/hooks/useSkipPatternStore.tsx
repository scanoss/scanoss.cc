import { create } from 'zustand';

interface SkipPatternStore {
  hasUnsavedChanges: boolean;
  setHasUnsavedChanges: (value: boolean) => void;
  nodesWithUnsavedChanges: Set<string>;
  addNodeWithUnsavedChanges: (nodeId: string) => void;
  removeNodeWithUnsavedChanges: (nodeId: string) => void;
  clearNodesWithUnsavedChanges: () => void;
}

const useSkipPatternStore = create<SkipPatternStore>((set) => ({
  hasUnsavedChanges: false,
  setHasUnsavedChanges: (value) => set({ hasUnsavedChanges: value }),
  nodesWithUnsavedChanges: new Set<string>(),
  addNodeWithUnsavedChanges: (nodeId) =>
    set((state) => {
      const newSet = new Set(state.nodesWithUnsavedChanges);
      newSet.add(nodeId);
      return { nodesWithUnsavedChanges: newSet };
    }),
  removeNodeWithUnsavedChanges: (nodeId) =>
    set((state) => {
      const newSet = new Set(state.nodesWithUnsavedChanges);
      newSet.delete(nodeId);
      return { nodesWithUnsavedChanges: newSet };
    }),
  clearNodesWithUnsavedChanges: () => set({ nodesWithUnsavedChanges: new Set<string>() }),
}));

export default useSkipPatternStore;
