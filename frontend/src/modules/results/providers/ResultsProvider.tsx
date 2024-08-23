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
  results: Result[];
  setResults: Dispatch<SetStateAction<Result[]>>;
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

  return (
    <ResultsContext.Provider
      value={{
        results: orderedResults,
        setResults,
      }}
    >
      {children}
    </ResultsContext.Provider>
  );
};

export const useResults = () => useContext(ResultsContext);
