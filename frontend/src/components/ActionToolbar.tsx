import { ReloadIcon } from '@radix-ui/react-icons';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { RotateCcw, RotateCw, Save } from 'lucide-react';

import useDebounce from '@/hooks/useDebounce';
import useQueryState from '@/hooks/useQueryState';
import FileService from '@/modules/files/infra/service';
import { MatchType } from '@/modules/results/domain';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import { Button } from './ui/button';
import { Separator } from './ui/separator';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';
import { useToast } from './ui/use-toast';

export default function ActionToolbar() {
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
  );
}
