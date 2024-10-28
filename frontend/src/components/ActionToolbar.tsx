import { ReloadIcon } from '@radix-ui/react-icons';
import { useMutation } from '@tanstack/react-query';
import { RotateCcw, RotateCw, Save } from 'lucide-react';

import useDebounce from '@/hooks/useDebounce';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import useQueryState from '@/hooks/useQueryState';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import useComponentFilterStore from '@/modules/components/stores/useComponentFilterStore';
import { DEBOUNCE_QUERY_MS } from '@/modules/results/constants';
import { MatchType } from '@/modules/results/domain';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import { Save as SaveBomChanges } from '../../wailsjs/go/service/ScanossSettingsServiceImp';
import { Button } from './ui/button';
import { Separator } from './ui/separator';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';
import { useToast } from './ui/use-toast';

export default function ActionToolbar() {
  const { toast } = useToast();

  const undo = useComponentFilterStore((state) => state.undo);
  const redo = useComponentFilterStore((state) => state.redo);
  const canUndo = useComponentFilterStore((state) => state.canUndo);
  const canRedo = useComponentFilterStore((state) => state.canRedo);

  const fetchResults = useResultsStore((state) => state.fetchResults);

  const [filterByMatchType] = useQueryState<MatchType | 'all'>('matchType', 'all');
  const [query] = useQueryState<string>('q', '');
  const debouncedQuery = useDebounce<string>(query, DEBOUNCE_QUERY_MS);

  const { mutate: saveChanges, isPending } = useMutation({
    mutationFn: SaveBomChanges,
    onSuccess: async () => {
      toast({
        title: 'Success',
        description: `Your changes have been successfully saved.`,
      });
      await fetchResults(filterByMatchType === 'all' ? undefined : filterByMatchType, debouncedQuery);
    },
    onError: (e) => {
      toast({
        title: 'Error',
        variant: 'destructive',
        description: e instanceof Error ? e.message : 'An error occurred while trying to save changes.',
      });
    },
  });

  useKeyboardShortcut(KEYBOARD_SHORTCUTS.undo.keys, undo);
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.redo.keys, redo);
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.save.keys, () => saveChanges());

  return (
    <div className="flex justify-end gap-4">
      <div className="flex gap-2">
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="ghost" className="h-full w-14 rounded-none" onClick={() => undo()} disabled={!canUndo}>
              <div className="flex flex-col items-center gap-1">
                <span className="text-xs">Undo</span>
                <RotateCcw className="h-4 w-4" />
              </div>
            </Button>
          </TooltipTrigger>
          <TooltipContent>Undo (⌘ + Z)</TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="ghost" className="h-full w-14 rounded-none" onClick={() => redo()} disabled={!canRedo}>
              <div className="flex flex-col items-center gap-1">
                <span className="text-xs">Redo</span>
                <RotateCw className="h-4 w-4" />
              </div>
            </Button>
          </TooltipTrigger>
          <TooltipContent>Redo (⌘ ⇧ Z)</TooltipContent>
        </Tooltip>
      </div>
      <Separator orientation="vertical" className="h-1/2 self-center" />
      <Tooltip>
        <TooltipTrigger asChild>
          <Button onClick={() => saveChanges()} disabled={isPending} className="self-center" size="icon">
            {isPending ? <ReloadIcon className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
          </Button>
        </TooltipTrigger>
        <TooltipContent>Save Changes (⌘ + S)</TooltipContent>
      </Tooltip>
    </div>
  );
}
