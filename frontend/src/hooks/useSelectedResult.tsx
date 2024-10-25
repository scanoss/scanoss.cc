import useLocalFilePath from '@/modules/files/hooks/useLocalFilePath';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import { entities } from '../../wailsjs/go/models';

export default function useSelectedResult(): entities.ResultDTO | undefined {
  const pendingResults = useResultsStore((state) => state.pendingResults);
  const completedResults = useResultsStore((state) => state.completedResults);
  const results = [...pendingResults, ...completedResults];
  const localFilePath = useLocalFilePath();
  const result = results.find((result) => result.path === localFilePath);

  return result;
}
