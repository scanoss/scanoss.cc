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
  selectedResults: Set<string>;
}

interface ResultsActions {
  handleCompleteResult: (
    args: HandleCompleteResultArgs
  ) => Promise<string | null>;
  setResults: (results: entities.ResultDTO[]) => void;
  undo: () => Promise<void>;
  redo: () => Promise<void>;
  updateUndoRedoState: () => Promise<void>;
  fetchResults: (matchType?: MatchType, query?: string) => Promise<void>;
  setSelectedResults: (selectedResults: Set<string>) => void;
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
    pendingResults: [],
    completedResults: [],
    canUndo: false,
    canRedo: false,
    isLoading: false,
    error: null,

    selectedResults: new Set<string>(),
    setSelectedResults: (selectedResults) =>
      set({ selectedResults }, false, 'SET_SELECTED_RESULTS'),

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

    handleCompleteResult: async ({ path, purl, action, comment }) => {
      await FileService.filterComponentByPath({
        path,
        purl,
        action,
        comment,
      });
      await get().updateUndoRedoState();
      await get().fetchResults();
      return getNextPendingResultPathRoute(get().pendingResults);
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
