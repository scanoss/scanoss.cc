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
  results: entities.ResultDTO[];
  pendingResults: entities.ResultDTO[];
  completedResults: entities.ResultDTO[];
  canUndo: boolean;
  canRedo: boolean;
  isLoading: boolean;
  error: string | null;
}

interface ResultsActions {
  handleCompleteResult: (
    path: string | undefined,
    purl: string,
    action: FilterAction,
    currentPath: string
  ) => Promise<string | null>;
  setResults: (results: entities.ResultDTO[]) => void;
  undo: () => Promise<void>;
  redo: () => Promise<void>;
  updateUndoRedoState: () => Promise<void>;
  fetchResults: (matchType?: MatchType, query?: string) => Promise<void>;
}

type ResultsStore = ResultsState & ResultsActions;

const useResultsStore = create<ResultsStore>()(
  devtools((set, get) => ({
    results: [],
    pendingResults: [],
    completedResults: [],
    canUndo: false,
    canRedo: false,
    isLoading: false,
    error: null,

    setResults: (results) =>
      set(
        {
          results,
          pendingResults: results.filter((r) => r.workflow_state === 'pending'),
          completedResults: results.filter(
            (r) => r.workflow_state === 'completed'
          ),
        },
        false,
        'SET_RESULTS'
      ),

    handleCompleteResult: async (path, purl, action, currentPath) => {
      await FileService.filterComponentByPath({
        path,
        purl,
        action,
      });
      await get().updateUndoRedoState();
      await get().fetchResults();
      return getNextPendingResultPathRoute(currentPath, get().pendingResults);
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
  path: string,
  pendingResults: entities.ResultDTO[]
): string | null => {
  const currentSelectedIndex = pendingResults.findIndex((r) => r.path === path);
  const isLast = currentSelectedIndex === pendingResults.length - 1;

  if (currentSelectedIndex === -1) return null;

  if (isLast && pendingResults.length > 0) {
    return `/files/${encodeFilePath(pendingResults[0].path)}`;
  }

  const nextPendingResult = pendingResults[currentSelectedIndex + 1];
  return `/files/${encodeFilePath(nextPendingResult.path)}`;
};

export default useResultsStore;
