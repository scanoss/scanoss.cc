import { ReloadIcon } from '@radix-ui/react-icons';
import { useMutation } from '@tanstack/react-query';
import { Check, PackageMinus, RefreshCwOff, Replace, Save } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { useToast } from '@/components/ui/use-toast';
import { Component } from '@/modules/results/domain';

import { FilterAction } from '../domain';
import FileService from '../infra/service';
import FileActionButton from './FileActionButton';

interface FileActionsMenuProps {
  component: Component;
}

export default function FileActionsMenu({ component }: FileActionsMenuProps) {
  const { toast } = useToast();

  const { mutate: saveChanges, isPending } = useMutation({
    mutationFn: () => FileService.saveBomChanges(),
    onSuccess: () => {
      toast({
        title: 'Success',
        description: `Your changes have been successfully saved.`,
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
    <div className="flex h-[65px] justify-center border-b border-b-border px-4">
      <div className="ml-auto flex gap-2">
        <FileActionButton
          action={FilterAction.Include}
          component={component}
          icon={<Check className="h-5 w-5 stroke-green-500" />}
          description="By including a file/component, you force the engine to consider it with priority in future scans."
        />
        <FileActionButton
          action={FilterAction.Remove}
          component={component}
          description="Dismissing a file/component will exclude it from future scan results."
          icon={<PackageMinus className="h-5 w-5 stroke-red-500" />}
        />
        <FileActionButton
          action={FilterAction.Ignore}
          component={component}
          icon={<RefreshCwOff className="h-5 w-5 stroke-amber-500" />}
          description="By ignoring a file/component, you prevent the engine from considering it in future scans."
          isDisabled
        />
        <FileActionButton
          action={FilterAction.Replace}
          component={component}
          icon={<Replace className="h-5 w-5 stroke-cyan-500" />}
          description="Replace a file/component with another one."
          isDisabled
        />
      </div>
      <div className="ml-auto flex items-center">
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
