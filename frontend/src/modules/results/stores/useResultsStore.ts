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
  canRedo: boolean;
  canUndo: boolean;
  error: string | null;
  isLoading: boolean;
  lastSelectedIndex: number;
  results: entities.ResultDTO[];
  selectedResults: entities.ResultDTO[];
}

interface ResultsActions {
  handleCompleteResult: (
    args: HandleCompleteResultArgs
  ) => Promise<string | null>;
  fetchResults: (matchType?: MatchType, query?: string) => Promise<void>;
  redo: () => Promise<void>;
  selectResultRange: (
    endResult: entities.ResultDTO,
    selectionType: 'pending' | 'completed'
  ) => void;
  setLastSelectedIndex: (index: number) => void;
  setResults: (results: entities.ResultDTO[]) => void;
  setSelectedResults: (selectedResults: entities.ResultDTO[]) => void;
  toggleResultSelection: (
    result: entities.ResultDTO,
    selectionType: 'pending' | 'completed'
  ) => void;
  undo: () => Promise<void>;
  updateUndoRedoState: () => Promise<void>;
}

interface HandleCompleteResultArgs {
  filterBy: 'by_file' | 'by_purl';
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
    lastSelectedIndex: -1,
    anchorIndex: -1,
    selectedResults: [],

    setSelectedResults: (selectedResults) =>
      set({ selectedResults }, false, 'SET_SELECTED_RESULTS'),

    setLastSelectedIndex: (index) =>
      set({ lastSelectedIndex: index }, false, 'SET_LAST_SELECTED_INDEX'),

    setResults: (results) =>
      set(
        {
          results,
        },
        false,
        'SET_RESULTS'
      ),

    handleCompleteResult: async (args: HandleCompleteResultArgs) => {
      const { filterBy, action, comment } = args;
      const selectedResults = get().selectedResults;

      const finalDto = selectedResults.map((result) => ({
        purl: result.purl,
        action: action,
        comment: comment,
        ...(filterBy === 'by_file' && { path: result.path }),
      }));

      await FileService.filterComponents(finalDto);

      await get().updateUndoRedoState();
      await get().fetchResults();

      // Clear the selection after completing the action
      set({ selectedResults: [] }, false, 'CLEAR_SELECTED_RESULTS');

      const allResults = get().results;
      const pendingResults = allResults.filter(
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

    toggleResultSelection: (result, selectionType) =>
      set((state) => {
        const { selectedResults } = state;
        const resultType = result.workflow_state as 'pending' | 'completed';

        // Prevent selection across different types
        if (selectionType !== resultType) return state;

        const index = selectedResults.findIndex((r) => r.path === result.path);

        if (index !== -1) {
          // Deselect if already selected
          const newSelectedResults = selectedResults.filter(
            (_, i) => i !== index
          );
          return { selectedResults: newSelectedResults };
        } else {
          // Add new selection
          return {
            selectedResults: [...selectedResults, result],
          };
        }
      }),

    selectResultRange: (endResult, selectionType) =>
      set((state) => {
        const { results, selectedResults, lastSelectedIndex } = state;

        const resultType = endResult.workflow_state as 'pending' | 'completed';

        if (selectionType !== resultType) return state;

        const startIndex =
          lastSelectedIndex !== -1
            ? lastSelectedIndex
            : results.findIndex((r) => r.path === selectedResults[0]?.path);
        const endIndex = results.findIndex((r) => r.path === endResult.path);

        if (startIndex === -1 || endIndex === -1) return state;

        const [start, end] =
          startIndex < endIndex
            ? [startIndex, endIndex]
            : [endIndex, startIndex];

        const newSelectedResults = results
          .slice(start, end + 1)
          .filter((r) => r.workflow_state === selectionType);

        return {
          selectedResults: newSelectedResults,
          lastSelectedIndex: endIndex,
        };
      }),
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
