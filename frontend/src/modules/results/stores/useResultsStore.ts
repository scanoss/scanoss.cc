import { entities } from 'wailsjs/go/models';
import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

import { GetAll } from '../../../../wailsjs/go/service/ResultServiceImpl';
import { MatchType } from '../domain';

interface ResultsState {
  completedResults: entities.ResultDTO[];
  error: string | null;
  isLoading: boolean;
  lastSelectedIndex: number;
  lastSelectionType: 'pending' | 'completed' | null;
  pendingResults: entities.ResultDTO[];
  selectedResults: entities.ResultDTO[];
  query: string;
  filterByMatchType: MatchType | 'all';
}

interface ResultsActions {
  fetchResults: () => Promise<void>;
  moveToNextResult: () => void;
  moveToPreviousResult: () => void;
  selectResultRange: (endResult: entities.ResultDTO, selectionType: 'pending' | 'completed') => void;
  setLastSelectedIndex: (index: number) => void;
  setLastSelectionType: (type: 'pending' | 'completed') => void;
  setSelectedResults: (selectedResults: entities.ResultDTO[]) => void;
  toggleResultSelection: (result: entities.ResultDTO, selectionType: 'pending' | 'completed') => void;
  setQuery: (query: string) => void;
  setFilterByMatchType: (matchType: MatchType | 'all') => void;
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
    query: '',
    filterByMatchType: 'all',

    setSelectedResults: (selectedResults) => set({ selectedResults }, false, 'SET_SELECTED_RESULTS'),

    setLastSelectedIndex: (index) => set({ lastSelectedIndex: index }, false, 'SET_LAST_SELECTED_INDEX'),

    setLastSelectionType: (type: 'pending' | 'completed') => set({ lastSelectionType: type }, false, 'SET_LAST_SELECTION_TYPE'),

    fetchResults: async () => {
      const { selectedResults, filterByMatchType, query } = get();
      set({ isLoading: true, error: null }, false, 'FETCH_RESULTS');
      try {
        const results = await GetAll({
          match_type: filterByMatchType === 'all' ? undefined : filterByMatchType,
          query,
        });
        const pendingResults = results.filter((r) => r.workflow_state === 'pending');
        const completedResults = results.filter((r) => r.workflow_state === 'completed');

        set({ pendingResults, completedResults });

        // When the app first loads or if changing the query, select the first result
        if (!selectedResults.length) {
          const hasPendingResults = pendingResults.length > 0;
          const hasCompletedResults = completedResults.length > 0;
          const firstSelectedResult = hasPendingResults ? pendingResults[0] : hasCompletedResults ? completedResults[0] : null;
          const selectedResults = firstSelectedResult ? [firstSelectedResult] : [];

          set({ selectedResults, lastSelectionType: firstSelectedResult?.workflow_state as 'pending' | 'completed', lastSelectedIndex: 0 });
        }
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

        const startIndex = lastSelectedIndex !== -1 ? lastSelectedIndex : resultsOfType.findIndex((r) => r.path === selectedResults[0]?.path);
        const endIndex = resultsOfType.findIndex((r) => r.path === endResult.path);

        if (startIndex === -1 || endIndex === -1) return state;

        const [start, end] = startIndex < endIndex ? [startIndex, endIndex] : [endIndex, startIndex];

        const newSelectedResults = resultsOfType.slice(start, end + 1).filter((r) => r.workflow_state === selectionType);

        return {
          selectedResults: newSelectedResults,
          lastSelectedIndex: startIndex,
          lastSelectionType: selectionType,
        };
      }),

    moveToNextResult: () => {
      const { pendingResults, completedResults, lastSelectionType, selectedResults } = get();

      const resultsOfType = lastSelectionType === 'pending' ? pendingResults : completedResults;
      const lastResultInSelection = selectedResults[selectedResults.length - 1]; // To handle multi select as well
      const currentResultIndex = resultsOfType.findIndex((r) => r.path === lastResultInSelection.path);
      let nextResult: entities.ResultDTO = lastResultInSelection;
      let nextSelectionType = lastSelectionType;
      let newResultIndex = currentResultIndex + 1;

      // If we are at the end of the results, go to the opposite type
      if (currentResultIndex === resultsOfType.length - 1) {
        const oppositeResults = lastSelectionType === 'pending' ? completedResults : pendingResults;
        nextResult = oppositeResults[0];
        nextSelectionType = lastSelectionType === 'pending' ? 'completed' : 'pending';
        newResultIndex = 0;

        // If there are no opposite results, do nothing, just keep the current selection
        if (oppositeResults.length === 0) return;
      } else {
        nextResult = resultsOfType[newResultIndex];
      }

      set({
        selectedResults: [nextResult],
        lastSelectedIndex: newResultIndex,
        lastSelectionType: nextSelectionType,
      });
    },

    moveToPreviousResult: () => {
      const { pendingResults, completedResults, lastSelectionType, selectedResults } = get();

      const resultsOfType = lastSelectionType === 'pending' ? pendingResults : completedResults;
      const lastResultInSelection = selectedResults[selectedResults.length - 1]; // To handle multi select as well
      const currentResultIndex = resultsOfType.findIndex((r) => r.path === lastResultInSelection.path);
      let newSelectedResult: entities.ResultDTO = lastResultInSelection;
      let newSelectionType = lastSelectionType;
      let newResultIndex;

      // If we are at the beginning of the results
      if (currentResultIndex === 0) {
        const oppositeResults = lastSelectionType === 'pending' ? completedResults : pendingResults;

        if (!oppositeResults.length) return;

        const lastPosition = oppositeResults.length - 1;
        const lastResultFromOpposite = oppositeResults[lastPosition];
        newSelectedResult = lastResultFromOpposite;
        newSelectionType = lastSelectionType === 'pending' ? 'completed' : 'pending';
        newResultIndex = lastPosition;
      } else {
        newSelectedResult = resultsOfType[currentResultIndex - 1];
        newResultIndex = currentResultIndex - 1;
      }

      set({
        selectedResults: [newSelectedResult],
        lastSelectedIndex: newResultIndex,
        lastSelectionType: newSelectionType,
      });
    },
    setQuery: (query) => set({ query }, false, 'SET_QUERY'),
    setFilterByMatchType: (matchType) => set({ filterByMatchType: matchType }, false, 'SET_FILTER_BY_MATCH_TYPE'),
  }))
);

export default useResultsStore;
