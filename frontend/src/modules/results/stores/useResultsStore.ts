import { entities } from 'wailsjs/go/models';
import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

import { MatchType } from '../domain';
import ResultService from '../infra/service';

interface ResultsState {
  completedResults: entities.ResultDTO[];
  error: string | null;
  isLoading: boolean;
  lastSelectedIndex: number;
  lastSelectionType: 'pending' | 'completed' | null;
  pendingResults: entities.ResultDTO[];
  selectedResults: entities.ResultDTO[];
}

interface ResultsActions {
  fetchResults: (matchType?: MatchType, query?: string) => Promise<void>;
  getNextResult: () => entities.ResultDTO;
  getPreviousResult: () => entities.ResultDTO;
  selectResultRange: (endResult: entities.ResultDTO, selectionType: 'pending' | 'completed') => void;
  setLastSelectedIndex: (index: number) => void;
  setLastSelectionType: (type: 'pending' | 'completed') => void;
  setSelectedResults: (selectedResults: entities.ResultDTO[]) => void;
  toggleResultSelection: (result: entities.ResultDTO, selectionType: 'pending' | 'completed') => void;
}

type ResultsStore = ResultsState & ResultsActions;

const useResultsStore = create<ResultsStore>()(
  devtools((set, get) => ({
    pendingResults: [],
    completedResults: [],
    isLoading: false,
    error: null,
    lastSelectedIndex: -1,
    selectedResults: [],
    lastSelectionType: null,

    setSelectedResults: (selectedResults) => set({ selectedResults }, false, 'SET_SELECTED_RESULTS'),

    setLastSelectedIndex: (index) => set({ lastSelectedIndex: index }, false, 'SET_LAST_SELECTED_INDEX'),

    setLastSelectionType: (type: 'pending' | 'completed') =>
      set({ lastSelectionType: type }, false, 'SET_LAST_SELECTION_TYPE'),

    fetchResults: async (matchType, query) => {
      set({ isLoading: true, error: null }, false, 'FETCH_RESULTS');
      try {
        const results = await ResultService.getAll(matchType, query);
        const pendingResults = results.filter((r) => r.workflow_state === 'pending');
        const completedResults = results.filter((r) => r.workflow_state === 'completed');
        set({ pendingResults, completedResults });
      } catch (error) {
        set({
          error: error instanceof Error ? error.message : 'An error occurred while fetching results',
        });
      } finally {
        set({ isLoading: false }, false, 'FETCH_RESULTS');
      }
    },

    toggleResultSelection: (result, selectionType) =>
      set((state) => {
        const { selectedResults, lastSelectionType } = state;
        const resultType = result.workflow_state as 'pending' | 'completed';

        if (selectionType !== lastSelectionType) {
          return { selectedResults: [result], lastSelectionType: resultType };
        }

        const index = selectedResults.findIndex((r) => r.path === result.path);

        if (index !== -1) {
          const newSelectedResults = selectedResults.filter((_, i) => i !== index);
          return { selectedResults: newSelectedResults };
        } else {
          return {
            selectedResults: [...selectedResults, result],
          };
        }
      }),

    selectResultRange: (endResult, selectionType) =>
      set((state) => {
        const { selectedResults, lastSelectedIndex, lastSelectionType, pendingResults, completedResults } = state;

        const resultType = endResult.workflow_state as 'pending' | 'completed';

        if (selectionType !== lastSelectionType) {
          return {
            selectedResults: [endResult],
            lastSelectedIndex: -1,
            lastSelectionType: resultType,
          };
        }

        const resultsOfType = lastSelectionType === 'pending' ? pendingResults : completedResults;

        const startIndex =
          lastSelectedIndex !== -1
            ? lastSelectedIndex
            : resultsOfType.findIndex((r) => r.path === selectedResults[0]?.path);
        const endIndex = resultsOfType.findIndex((r) => r.path === endResult.path);

        if (startIndex === -1 || endIndex === -1) return state;

        const [start, end] = startIndex < endIndex ? [startIndex, endIndex] : [endIndex, startIndex];

        const newSelectedResults = resultsOfType
          .slice(start, end + 1)
          .filter((r) => r.workflow_state === selectionType);

        return {
          selectedResults: newSelectedResults,
          lastSelectedIndex: startIndex,
          lastSelectionType: selectionType,
        };
      }),

    getNextResult: () => {
      const { pendingResults, completedResults, lastSelectionType, selectedResults } = get();

      const resultsOfType = lastSelectionType === 'pending' ? pendingResults : completedResults;
      const lastResultInSelection = selectedResults[selectedResults.length - 1]; // To handle multi select as well
      const newResultIndex = resultsOfType.findIndex((r) => r.path === lastResultInSelection.path) + 1;
      let nextResult: entities.ResultDTO;

      if (newResultIndex >= resultsOfType.length) {
        const oppositeResults = lastSelectionType === 'pending' ? completedResults : pendingResults;
        nextResult = oppositeResults[0];

        return nextResult;
      }

      nextResult = resultsOfType[newResultIndex];

      return nextResult;
    },
    getPreviousResult: () => {
      const { pendingResults, completedResults, lastSelectionType, selectedResults } = get();

      const resultsOfType = lastSelectionType === 'pending' ? pendingResults : completedResults;
      const newResultIndex = resultsOfType.findIndex((r) => r.path === selectedResults[0]?.path) - 1;
      let previousResult: entities.ResultDTO;

      if (newResultIndex === -1) {
        const oppositeResults = lastSelectionType === 'pending' ? completedResults : pendingResults;
        previousResult = oppositeResults[oppositeResults.length - 1];

        return previousResult;
      }

      previousResult = resultsOfType[newResultIndex];

      return previousResult;
    },
  }))
);

export default useResultsStore;
