import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

import FileService from '@/modules/files/infra/service';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import {
  ComponentFilterCanRedo,
  ComponentFilterCanUndo,
  ComponentFilterRedo,
  ComponentFilterUndo,
} from '../../../../wailsjs/go/handlers/ComponentHandler';
import { entities } from '../../../../wailsjs/go/models';
import { FilterAction } from '../domain';

interface ComponentFilterState {
  action: FilterAction | undefined;
  canRedo: boolean;
  canUndo: boolean;
  comment: string | undefined;
  filterBy: 'by_file' | 'by_purl' | undefined;
  withComment: boolean;
}

export interface OnFilterComponentArgs {
  replaceWith?: {
    purl: string;
    name: string;
  };
}

interface ComponentFilterActions {
  onFilterComponent: (args?: OnFilterComponentArgs) => Promise<entities.ResultDTO | null>;
  redo: () => Promise<void>;
  setAction: (action: FilterAction) => void;
  setFilterBy: (filterBy: 'by_file' | 'by_purl') => void;
  setComment: (comment: string | undefined) => void;
  setWithComment: (withComment: boolean) => void;
  undo: () => Promise<void>;
  updateUndoRedoState: () => Promise<void>;
}

type ComponentFilterStore = ComponentFilterState & ComponentFilterActions;

const useComponentFilterStore = create<ComponentFilterStore>()(
  devtools((set, get) => ({
    canUndo: false,
    canRedo: false,
    action: undefined,
    filterBy: undefined,
    comment: undefined,
    withComment: false,

    setAction: (action) => set({ action }),
    setFilterBy: (filterBy) => set({ filterBy }),
    setComment: (comment) => set({ comment }),
    setWithComment: (withComment) => set({ withComment }),

    onFilterComponent: async (args?: OnFilterComponentArgs) => {
      const { replaceWith } = args || {};
      const { filterBy, action, comment } = get();

      if (!action || !filterBy) {
        throw new Error('There was an error filtering the component. Please try again.');
      }

      if (action === FilterAction.Replace && !replaceWith) {
        throw new Error('There was an error replacing the component. Please try again.');
      }

      const selectedResults = useResultsStore.getState().selectedResults;

      const finalDto = selectedResults.map((result) => ({
        action,
        comment,
        purl: result.detected_purl ?? '',
        ...(filterBy === 'by_file' && { path: result.path }),
        ...(replaceWith && {
          replace_with_purl: replaceWith.purl,
          replace_with_name: replaceWith.name,
        }),
      }));

      const nextResult = useResultsStore.getState().getNextResult();

      await FileService.filterComponents(finalDto);

      await get().updateUndoRedoState();
      await useResultsStore.getState().fetchResults();

      useResultsStore.setState({ selectedResults: [] }, false, 'CLEAR_SELECTED_RESULTS');

      return nextResult;
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
      const [canUndo, canRedo] = await Promise.all([ComponentFilterCanUndo(), ComponentFilterCanRedo()]);
      set({ canUndo, canRedo }, false, 'UPDATE_UNDO_REDO_STATE');
    },
  }))
);

export default useComponentFilterStore;
