import {
  createContext,
  Dispatch,
  ReactNode,
  SetStateAction,
  useContext,
  useMemo,
  useState,
} from 'react';

import { MatchType, Result } from '@/modules/results/domain';

export interface ResultsContext {
  handleStageResult: (path: string) => void;
  results: Result[];
  setResults: Dispatch<SetStateAction<Result[]>>;
  stagedResults: Result[];
  unstagedResults: Result[];
}

export const ResultsContext = createContext<ResultsContext>(
  {} as ResultsContext
);

const resultsPriorityOrder = [MatchType.Snippet, MatchType.File];

export const ResultsProvider = ({ children }: { children: ReactNode }) => {
  const [results, setResults] = useState<Result[]>([]);

  const groupedResultsByMatchType = useMemo(
    () =>
      results.reduce(
        (acc, result) => {
          if (!acc[result.matchType]) {
            acc[result.matchType] = [];
          }

          acc[result.matchType].push(result);

          return acc;
        },
        {} as Record<MatchType, Result[]>
      ),
    [results]
  );

  const orderedResults = resultsPriorityOrder.flatMap(
    (matchType) => groupedResultsByMatchType[matchType] ?? []
  );

  const groupedResultsByState = useMemo(
    () =>
      orderedResults.reduce(
        (acc, result) => {
          if (!acc[result.state]) {
            acc[result.state] = [];
          }

          acc[result.state].push(result);

          return acc;
        },
        {} as Record<string, Result[]>
      ),
    [orderedResults]
  );

  const handleStageResult = (path: string) => {
    setResults((results) =>
      results.map((r) => (r.path === path ? { ...r, state: 'staged' } : r))
    );
  };

  return (
    <ResultsContext.Provider
      value={{
        handleStageResult,
        results: orderedResults,
        setResults,
        stagedResults: groupedResultsByState.staged ?? [],
        unstagedResults: groupedResultsByState.unstaged ?? [],
      }}
    >
      {children}
    </ResultsContext.Provider>
  );
};

export const useResults = () => useContext(ResultsContext);
