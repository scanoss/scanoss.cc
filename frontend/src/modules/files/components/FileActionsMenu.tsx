import { ReloadIcon } from '@radix-ui/react-icons';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Check, PackageMinus, RotateCcw, RotateCw, Save } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { useToast } from '@/components/ui/use-toast';
import useDebounce from '@/hooks/useDebounce';
import useQueryState from '@/hooks/useQueryState';
import { FilterAction, MatchType } from '@/modules/results/domain';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import FileService from '../infra/service';
import FileActionButton from './FileActionButton';

export default function FileActionsMenu() {
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const undo = useResultsStore((state) => state.undo);
  const redo = useResultsStore((state) => state.redo);
  const canUndo = useResultsStore((state) => state.canUndo);
  const canRedo = useResultsStore((state) => state.canRedo);

  const [filterByMatchType] = useQueryState<MatchType | 'all'>(
    'matchType',
    'all'
  );

  const [query] = useQueryState<string>('q', '');
  const debouncedQuery = useDebounce<string>(query, 300);

  const { mutate: saveChanges, isPending } = useMutation({
    mutationFn: () => FileService.saveBomChanges(),
    onSuccess: async () => {
      toast({
        title: 'Success',
        description: `Your changes have been successfully saved.`,
      });
      await queryClient.refetchQueries({
        queryKey: ['results', filterByMatchType, debouncedQuery],
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
    <div className="grid h-full grid-cols-2 md:grid-cols-3">
      <div className="hidden md:block"></div>
      <div className="flex gap-2 md:justify-center">
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
      <div className="flex justify-end gap-4">
        <div className="flex gap-2">
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                className="h-full w-14 rounded-none"
                onClick={() => undo()}
                disabled={!canUndo}
              >
                <div className="flex flex-col items-center gap-1">
                  <span className="text-xs">Undo</span>
                  <RotateCcw className="h-4 w-4" />
                </div>
              </Button>
            </TooltipTrigger>
            <TooltipContent>Undo</TooltipContent>
          </Tooltip>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                className="h-full w-14 rounded-none"
                onClick={() => redo()}
                disabled={!canRedo}
              >
                <div className="flex flex-col items-center gap-1">
                  <span className="text-xs">Redo</span>
                  <RotateCw className="h-4 w-4" />
                </div>
              </Button>
            </TooltipTrigger>
            <TooltipContent>Redo</TooltipContent>
          </Tooltip>
        </div>
        <Separator orientation="vertical" className="h-1/2 self-center" />
        <Button
          onClick={() => saveChanges()}
          disabled={isPending}
          className="self-center"
        >
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
