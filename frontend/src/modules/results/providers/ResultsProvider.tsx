import {
  createContext,
  Dispatch,
  ReactNode,
  SetStateAction,
  useContext,
  useMemo,
  useState,
} from 'react';
import { entities } from 'wailsjs/go/models';

import { encodeFilePath } from '@/lib/utils';
import { MatchType } from '@/modules/results/domain';

export interface ResultsContext {
  handleConfirmResult: (
    path: string,
    filterConfig: entities.FilterConfig
  ) => string | null;
  results: entities.ResultDTO[];
  setResults: Dispatch<SetStateAction<entities.ResultDTO[]>>;
  confirmedResults: entities.ResultDTO[];
  pendingResults: entities.ResultDTO[];
}

export const ResultsContext = createContext<ResultsContext>(
  {} as ResultsContext
);

const resultsPriorityOrder = [MatchType.Snippet, MatchType.File];

export const ResultsProvider = ({ children }: { children: ReactNode }) => {
  const [results, setResults] = useState<entities.ResultDTO[]>([]);

  const groupedResultsByMatchType = useMemo(
    () =>
      results.reduce(
        (acc, result) => {
          if (!acc[result.match_type]) {
            acc[result.match_type] = [];
          }

          acc[result.match_type].push(result);

          return acc;
        },
        {} as Record<string, entities.ResultDTO[]>
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
          if (!result.workflow_state) {
            return acc;
          }

          if (!acc[result.workflow_state]) {
            acc[result.workflow_state] = [];
          }

          acc[result.workflow_state].push(result);

          return acc;
        },
        {} as Record<string, entities.ResultDTO[]>
      ),
    [orderedResults]
  );

  const confirmedResults = groupedResultsByState['confirmed'] ?? [];
  const pendingResults = groupedResultsByState['pending'] ?? [];

  const handleConfirmResult = (
    path: string,
    filterConfig: entities.FilterConfig
  ): string | null => {
    // @ts-expect-error - This is a hack to update the state of the results
    setResults((results) =>
      results.map((result) =>
        result.path === path
          ? {
              ...result,
              workflow_state: 'confirmed',
              filter_config: filterConfig,
            }
          : result
      )
    );

    return getNextPendingResultPathRoute(path);
  };

  const getNextPendingResultPathRoute = (path: string): string | null => {
    const currentSelectedIndex = pendingResults.findIndex(
      (r) => r.path === path
    );

    const isLast = currentSelectedIndex === pendingResults.length - 1;

    if (currentSelectedIndex === -1) return null;

    if (isLast && pendingResults.length > 0) {
      return `/files/${encodeFilePath(pendingResults[0].path)}`;
    }

    const nextPendingResult = pendingResults[currentSelectedIndex + 1];

    return `/files/${encodeFilePath(nextPendingResult.path)}`;
  };

  return (
    <ResultsContext.Provider
      value={{
        handleConfirmResult,
        results: orderedResults,
        setResults,
        confirmedResults,
        pendingResults,
      }}
    >
      {children}
    </ResultsContext.Provider>
  );
};

export const useResults = () => useContext(ResultsContext);
