// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

import { ReloadIcon } from '@radix-ui/react-icons';
import { useMutation } from '@tanstack/react-query';
import { RotateCcw, RotateCw, Save } from 'lucide-react';

import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import useComponentFilterStore from '@/modules/components/stores/useComponentFilterStore';
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

  const { mutate: saveChanges, isPending } = useMutation({
    mutationFn: SaveBomChanges,
    onSuccess: async () => {
      toast({
        title: 'Success',
        description: `Your changes have been successfully saved.`,
      });
      await fetchResults();
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
