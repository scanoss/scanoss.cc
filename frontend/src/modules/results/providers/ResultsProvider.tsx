import {
  createContext,
  Dispatch,
  ReactNode,
  SetStateAction,
  useContext,
  useMemo,
  useState,
} from 'react';

import { encodeFilePath } from '@/lib/utils';
import {
  FilterAction,
  FilterBy,
  MatchType,
  Result,
} from '@/modules/results/domain';

export interface ResultsContext {
  handleStageResult: (
    path: string,
    action: FilterAction,
    filterBy: FilterBy
  ) => string | null;
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

  const stagedResults = groupedResultsByState.staged ?? [];
  const unstagedResults = groupedResultsByState.unstaged ?? [];

  const handleStageResult = (
    path: string,
    action: FilterAction,
    filterBy: FilterBy
  ): string | null => {
    setResults((results) =>
      results.map((result) =>
        result.path === path
          ? {
              ...result,
              state: 'staged',
              bomState: {
                action,
                filterBy,
              },
            }
          : result
      )
    );
    return getNextUnstagedResultPathRoute(path);
  };

  const getNextUnstagedResultPathRoute = (path: string): string | null => {
    const currentSelectedIndex = unstagedResults.findIndex(
      (r) => r.path === path
    );

    const isLast = currentSelectedIndex === unstagedResults.length - 1;

    if (currentSelectedIndex === -1) return null;

    if (isLast && unstagedResults.length > 0) {
      return `/files/${encodeFilePath(unstagedResults[0].path)}`;
    }

    const nextUnstagedResult = unstagedResults[currentSelectedIndex + 1];

    return `/files/${encodeFilePath(nextUnstagedResult.path)}`;
  };

  return (
    <ResultsContext.Provider
      value={{
        handleStageResult,
        results: orderedResults,
        setResults,
        stagedResults,
        unstagedResults,
      }}
    >
      {children}
    </ResultsContext.Provider>
  );
};

export const useResults = () => useContext(ResultsContext);
