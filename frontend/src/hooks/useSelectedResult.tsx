import useResultsStore from '@/modules/results/stores/useResultsStore';

import { entities } from '../../wailsjs/go/models';

export default function useSelectedResult(): entities.ResultDTO | undefined {
  const selectedResults = useResultsStore((state) => state.selectedResults);

  return selectedResults[0];
}
