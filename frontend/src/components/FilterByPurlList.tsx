import { FilterAction } from '@/modules/components/domain';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import { ScrollArea } from './ui/scroll-area';

interface FilterByPurlListProps {
  action: FilterAction;
}

export default function FilterByPurlList({ action }: FilterByPurlListProps) {
  const selectedResults = useResultsStore((state) => state.selectedResults);

  return (
    <div>
      <p>
        This action will {action} all matches with{' '}
        {selectedResults.length > 1 ? `the following PURLs: ` : `the same PURL:`}
      </p>

      {selectedResults.length > 1 ? (
        <ScrollArea className="py-2">
          <ul className="max-h-[200px] list-disc pl-6">
            {selectedResults.map((result) => (
              <li key={result.path}>{result.purl?.detected}</li>
            ))}
          </ul>
        </ScrollArea>
      ) : (
        <p>{selectedResults[0]?.purl?.detected}</p>
      )}
    </div>
  );
}
