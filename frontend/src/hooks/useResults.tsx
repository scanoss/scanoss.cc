import { useQuery } from '@tanstack/react-query';

import { DEBOUNCE_QUERY_MS } from '@/modules/results/constants';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import useDebounce from './useDebounce';

export const useResults = () => {
  const filterByMatchType = useResultsStore((state) => state.filterByMatchType);
  const query = useResultsStore((state) => state.query);
  const sort = useResultsStore((state) => state.sort);
  const fetchResults = useResultsStore((state) => state.fetchResults);

  const debouncedQuery = useDebounce<string>(query, DEBOUNCE_QUERY_MS);

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['results', debouncedQuery, filterByMatchType, sort] as const,
    queryFn: () => fetchResults(),
  });

  const reset = (): void => {
    refetch();
  };

  return {
    data,
    isLoading,
    reset,
  };
};
