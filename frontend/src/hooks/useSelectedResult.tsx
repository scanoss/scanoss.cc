import useLocalFilePath from '@/modules/files/hooks/useLocalFilePath';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import { entities } from '../../wailsjs/go/models';

export default function useSelectedResult(): entities.ResultDTO | undefined {
  const results = useResultsStore((state) => state.results);
  const localFilePath = useLocalFilePath();
  const result = results.find((result) => result.path === localFilePath);

  return result;
}
