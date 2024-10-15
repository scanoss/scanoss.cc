import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

import FileService from '@/modules/files/infra/service';
import useResultsStore from '@/modules/results/stores/useResultsStore';
import { getNextPendingResultPathRoute } from '@/modules/results/utils';

import {
  ComponentFilterCanRedo,
  ComponentFilterCanUndo,
  ComponentFilterRedo,
  ComponentFilterUndo,
} from '../../../../wailsjs/go/handlers/ComponentHandler';
import { FilterAction } from '../domain';

interface ComponentFilterState {
  action: FilterAction | null;
  canRedo: boolean;
  canUndo: boolean;
  comment?: string;
  filterBy: 'by_file' | 'by_purl' | null;
  withComment?: boolean;
}

interface ComponentFilterActions {
  onFilterComponent: () => Promise<string | null>;
  redo: () => Promise<void>;
  setAction: (action: FilterAction) => void;
  setFilterBy: (filterBy: 'by_file' | 'by_purl') => void;
  setWithComment: (withComment: boolean) => void;
  setComment: (comment: string) => void;
  undo: () => Promise<void>;
  updateUndoRedoState: () => Promise<void>;
}

type ComponentFilterStore = ComponentFilterState & ComponentFilterActions;

const useComponentFilterStore = create<ComponentFilterStore>()(
  devtools((set, get) => ({
    canUndo: false,
    canRedo: false,
    action: null,
    filterBy: null,
    withComment: false,
    comment: '',

    setAction: (action) => set({ action }),
    setFilterBy: (filterBy) => set({ filterBy }),
    setWithComment: (withComment) => set({ withComment }),
    setComment: (comment) => set({ comment }),

    onFilterComponent: async () => {
      const { filterBy, action, comment } = get();

      if (!action || !filterBy) {
        return null;
      }

      const selectedResults = useResultsStore.getState().selectedResults;

      const finalDto = selectedResults.map((result) => ({
        action,
        comment,
        purl: result.purl,
        ...(filterBy === 'by_file' && { path: result.path }),
      }));

      await FileService.filterComponents(finalDto);

      await get().updateUndoRedoState();
      await useResultsStore.getState().fetchResults();

      // Clear the selection after completing the action
      useResultsStore.setState(
        { selectedResults: [] },
        false,
        'CLEAR_SELECTED_RESULTS'
      );

      const allResults = useResultsStore.getState().results;
      const pendingResults = allResults.filter(
        (result) => result.workflow_state === 'pending'
      );

      return getNextPendingResultPathRoute(pendingResults);
    },

    undo: async () => {
      await ComponentFilterUndo();
      await get().updateUndoRedoState();
      await useResultsStore.getState().fetchResults();
    },

    redo: async () => {
      await ComponentFilterRedo();
      await get().updateUndoRedoState();
      await useResultsStore.getState().fetchResults();
    },

    updateUndoRedoState: async () => {
      const [canUndo, canRedo] = await Promise.all([
        ComponentFilterCanUndo(),
        ComponentFilterCanRedo(),
      ]);
      set({ canUndo, canRedo }, false, 'UPDATE_UNDO_REDO_STATE');
    },
  }))
);

export default useComponentFilterStore;
