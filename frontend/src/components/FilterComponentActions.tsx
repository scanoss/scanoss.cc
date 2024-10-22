import { Check, PackageMinus, Replace } from 'lucide-react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { withErrorHandling } from '@/lib/errors';
import { FilterAction } from '@/modules/components/domain';
import useComponentFilterStore, { OnFilterComponentArgs } from '@/modules/components/stores/useComponentFilterStore';

import FilterActionButton from './FilterActionButton';
import ReplaceComponentDialog from './ReplaceComponentDialog';
import { useToast } from './ui/use-toast';

export default function FilterComponentActions() {
  const { toast } = useToast();
  const navigate = useNavigate();

  const [showReplaceComponentDialog, setShowReplaceComponentDialog] = useState(false);

  const onFilterComponent = useComponentFilterStore((state) => state.onFilterComponent);

  const handleFilterComponent = withErrorHandling({
    asyncFn: async (args?: OnFilterComponentArgs) => {
      const nextResultRoute = await onFilterComponent(args);

      if (nextResultRoute) {
        navigate(nextResultRoute);
      }
    },
    onError: (error) => {
      console.error(error);
      toast({
        title: 'Error',
        description: 'An error occurred while filtering the component. Please try again.',
        variant: 'destructive',
      });
    },
  });

  const handleReplaceComponent = withErrorHandling({
    asyncFn: async (args: OnFilterComponentArgs) => {
      await handleFilterComponent(args);
      setShowReplaceComponentDialog(false);
    },
    onError: (error) => {
      console.error(error);
      toast({
        title: 'Error',
        description: 'An error occurred while replacing the component. Please try again.',
        variant: 'destructive',
      });
    },
  });

  return (
    <>
      <div className="flex gap-2 md:justify-center">
        <FilterActionButton
          action={FilterAction.Include}
          icon={<Check className="h-5 w-5 stroke-green-500" />}
          description="By including a file/component, you force the engine to consider it with priority in future scans."
          onAdd={handleFilterComponent}
          shortcutKeyByFile="i f"
          shortcutKeyByPurl="i c"
        />
        <FilterActionButton
          action={FilterAction.Remove}
          description="Dismissing a file/component will exclude it from future scan results."
          icon={<PackageMinus className="h-5 w-5 stroke-red-500" />}
          onAdd={handleFilterComponent}
          shortcutKeyByFile="d f"
          shortcutKeyByPurl="d c"
        />
        <FilterActionButton
          action={FilterAction.Replace}
          description="Replace detected components with another one."
          icon={<Replace className="h-5 w-5 stroke-yellow-500" />}
          onAdd={() => setShowReplaceComponentDialog(true)}
          shortcutKeyByFile="r f"
          shortcutKeyByPurl="r c"
        />
      </div>
      {showReplaceComponentDialog && (
        <ReplaceComponentDialog
          onOpenChange={() => setShowReplaceComponentDialog((prev) => !prev)}
          onReplaceComponent={handleReplaceComponent}
        />
      )}
    </>
  );
}
