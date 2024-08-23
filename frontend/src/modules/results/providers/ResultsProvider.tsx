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
      Object.groupBy(results ?? [], (result) => {
        return result.matchType;
      }),
    [results]
  );

  const orderedResults = resultsPriorityOrder.flatMap(
    (matchType) => groupedResultsByMatchType[matchType] ?? []
  );

  const groupedResultsByState = useMemo(
    () =>
      Object.groupBy(orderedResults, (result) => {
        return result.state;
      }),
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
