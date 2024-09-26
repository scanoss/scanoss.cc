import { ReloadIcon } from '@radix-ui/react-icons';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Check, PackageMinus, Save } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { useToast } from '@/components/ui/use-toast';
import { FilterAction } from '@/modules/results/domain';

import FileService from '../infra/service';
import FileActionButton from './FileActionButton';

export default function FileActionsMenu() {
  const { toast } = useToast();
  const queryClient = useQueryClient();

  const { mutate: saveChanges, isPending } = useMutation({
    mutationFn: () => FileService.saveBomChanges(),
    onSuccess: async () => {
      toast({
        title: 'Success',
        description: `Your changes have been successfully saved.`,
      });
      await queryClient.refetchQueries({
        queryKey: ['results'],
      });
    },
    onError: (e) => {
      toast({
        title: 'Error',
        variant: 'destructive',
        description:
          e instanceof Error
            ? e.message
            : 'An error occurred while trying to save changes.',
      });
    },
  });

  return (
    <div className="grid h-full grid-cols-3">
      <div className=""></div>
      <div className="flex justify-center gap-2">
        <FileActionButton
          action={FilterAction.Include}
          icon={<Check className="h-5 w-5 stroke-green-500" />}
          description="By including a file/component, you force the engine to consider it with priority in future scans."
        />
        <FileActionButton
          action={FilterAction.Remove}
          description="Dismissing a file/component will exclude it from future scan results."
          icon={<PackageMinus className="h-5 w-5 stroke-red-500" />}
        />
      </div>
      <div className="flex items-center justify-end">
        <Button size="sm" onClick={() => saveChanges()} disabled={isPending}>
          {isPending ? (
            <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
          ) : (
            <Save className="mr-2 h-4 w-4" />
          )}
          <span>Save Changes</span>
        </Button>
      </div>
    </div>
  );
}
