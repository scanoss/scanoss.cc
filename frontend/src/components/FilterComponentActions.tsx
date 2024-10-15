import { Check, PackageMinus, Replace } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

import { useConfirm } from '@/hooks/useConfirm';
import { useInputPrompt } from '@/hooks/useInputPrompt';
import { FilterAction, OnAddFilterArgs } from '@/modules/components/domain';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import FilterActionButton from './FilterActionButton';
import { ScrollArea } from './ui/scroll-area';
import { useToast } from './ui/use-toast';

export default function FilterComponentActions() {
  const { ask } = useConfirm();
  const { prompt } = useInputPrompt();
  const { toast } = useToast();
  const navigate = useNavigate();

  const selectedResults = useResultsStore((state) => state.selectedResults);
  const handleCompleteResult = useResultsStore(
    (state) => state.handleCompleteResult
  );

  const handleAddFilter = async (args: OnAddFilterArgs) => {
    const { filterBy, withComment } = args;

    try {
      let comment: string | undefined;

      if (withComment) {
        comment = await prompt({
          title: 'Add comments',
          input: {
            defaultValue: '',
            type: 'textarea',
          },
        });
        if (!comment) return;
      }

      return filterActions[filterBy](args);
    } catch (e) {
      console.error(e);
      toast({
        title: 'Error',
        variant: 'destructive',
        description:
          'An error ocurred while adding comments. Please try again.',
      });
    }
  };

  const handleFilterComponentByPurl = async (args: OnAddFilterArgs) => {
    const confirm = await ask(
      <div>
        <p>
          This action will {args.action} all matches with{' '}
          {selectedResults.length > 1
            ? `the following PURLs: `
            : `the same PURL:`}
        </p>

        {selectedResults.length > 1 ? (
          <ScrollArea className="py-2">
            <ul className="max-h-[200px] list-disc pl-6">
              {selectedResults.map((result) => (
                <li key={result.path}>{result.purl}</li>
              ))}
            </ul>
          </ScrollArea>
        ) : (
          <p>{selectedResults[0]?.purl}</p>
        )}
      </div>
    );

    if (confirm) {
      return handleFilterComponent(args);
    }
  };

  const handleFilterComponentByFile = async (args: OnAddFilterArgs) => {
    return handleFilterComponent(args);
  };

  const handleFilterComponent = async (args: OnAddFilterArgs) => {
    const nextResultRoute = await handleCompleteResult(args);

    if (nextResultRoute) {
      navigate(nextResultRoute);
    }
  };

  const filterActions: Record<
    'by_file' | 'by_purl',
    (args: OnAddFilterArgs) => Promise<void>
  > = {
    by_file: handleFilterComponentByFile,
    by_purl: handleFilterComponentByPurl,
  };

  return (
    <div className="flex gap-2 md:justify-center">
      <FilterActionButton
        action={FilterAction.Include}
        icon={<Check className="h-5 w-5 stroke-green-500" />}
        description="By including a file/component, you force the engine to consider it with priority in future scans."
        onAdd={handleAddFilter}
      />
      <FilterActionButton
        action={FilterAction.Remove}
        description="Dismissing a file/component will exclude it from future scan results."
        icon={<PackageMinus className="h-5 w-5 stroke-red-500" />}
        onAdd={handleAddFilter}
      />
      <FilterActionButton
        action={FilterAction.Replace}
        description="Dismissing a file/component will exclude it from future scan results."
        icon={<Replace className="h-5 w-5 stroke-yellow-500" />}
        onAdd={handleAddFilter}
      />
    </div>
  );
}
