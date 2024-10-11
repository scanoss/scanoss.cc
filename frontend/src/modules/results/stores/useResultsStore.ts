import { entities } from 'wailsjs/go/models';
import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

import { encodeFilePath } from '@/lib/utils';
import FileService from '@/modules/files/infra/service';

import {
  ComponentFilterCanRedo,
  ComponentFilterCanUndo,
  ComponentFilterRedo,
  ComponentFilterUndo,
} from '../../../../wailsjs/go/handlers/ComponentHandler';
import { FilterAction, MatchType } from '../domain';
import ResultService from '../infra/service';

interface ResultsState {
  anchorIndex: number;
  canRedo: boolean;
  canUndo: boolean;
  error: string | null;
  isLoading: boolean;
  lastSelectedIndex: number;
  results: entities.ResultDTO[];
  selectedFiles: Set<string>;
}

interface ResultsActions {
  handleCompleteResult: (
    args: HandleCompleteResultArgs
  ) => Promise<string | null>;
  fetchResults: (matchType?: MatchType, query?: string) => Promise<void>;
  redo: () => Promise<void>;
  setAnchorIndex: (anchorIndex: number) => void;
  setLastSelectedIndex: (lastSelectedIndex: number) => void;
  setResults: (results: entities.ResultDTO[]) => void;
  setSelectedFiles: (selectedFiles: Set<string>) => void;
  undo: () => Promise<void>;
  updateUndoRedoState: () => Promise<void>;
}

interface HandleCompleteResultArgs {
  path: string | undefined;
  purl: string;
  action: FilterAction;
  comment?: string | undefined;
}

type ResultsStore = ResultsState & ResultsActions;

const useResultsStore = create<ResultsStore>()(
  devtools((set, get) => ({
    results: [],
    canUndo: false,
    canRedo: false,
    isLoading: false,
    error: null,

    selectedFiles: new Set<string>(),
    setSelectedFiles: (selectedFiles) =>
      set({ selectedFiles }, false, 'SET_SELECTED_FILES'),

    lastSelectedIndex: -1,
    setLastSelectedIndex: (lastSelectedIndex) =>
      set({ lastSelectedIndex }, false, 'SET_LAST_SELECTED_INDEX'),

    anchorIndex: -1,
    setAnchorIndex: (anchorIndex) =>
      set({ anchorIndex }, false, 'SET_ANCHOR_INDEX'),

    setResults: (results) =>
      set(
        {
          results,
        },
        false,
        'SET_RESULTS'
      ),

    handleCompleteResult: async ({ path, purl, action, comment }) => {
      await FileService.filterComponentByPath({
        path,
        purl,
        action,
        comment,
      });
      await get().updateUndoRedoState();
      await get().fetchResults();

      const pendingResults = get().results.filter(
        (result) => result.workflow_state === 'pending'
      );

      return getNextPendingResultPathRoute(pendingResults);
    },

    undo: async () => {
      await ComponentFilterUndo();
      await get().updateUndoRedoState();
      await get().fetchResults();
    },

    redo: async () => {
      await ComponentFilterRedo();
      await get().updateUndoRedoState();
      await get().fetchResults();
    },

    updateUndoRedoState: async () => {
      const [canUndo, canRedo] = await Promise.all([
        ComponentFilterCanUndo(),
        ComponentFilterCanRedo(),
      ]);
      set({ canUndo, canRedo }, false, 'UPDATE_UNDO_REDO_STATE');
    },

    fetchResults: async (matchType, query) => {
      set({ isLoading: true, error: null }, false, 'FETCH_RESULTS');
      try {
        const results = await ResultService.getAll(matchType, query);
        get().setResults(results);
      } catch (error) {
        set({
          error:
            error instanceof Error
              ? error.message
              : 'An error occurred while fetching results',
        });
      } finally {
        set({ isLoading: false }, false, 'FETCH_RESULTS');
      }
    },
  }))
);

const getNextPendingResultPathRoute = (
  pendingResults: entities.ResultDTO[]
): string | null => {
  const firstAvailablePendingResult = pendingResults[0];
  if (!firstAvailablePendingResult) return null;

  return `/files/${encodeFilePath(firstAvailablePendingResult.path)}`;
};

export default useResultsStore;
