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

import { Check, PackageMinus, Replace } from 'lucide-react';
import { useState } from 'react';

import { withErrorHandling } from '@/lib/errors';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import { FilterAction } from '@/modules/components/domain';
import useComponentFilterStore, { OnFilterComponentArgs } from '@/modules/components/stores/useComponentFilterStore';

import FilterActionButton from './FilterActionButton';
import ReplaceComponentDialog from './ReplaceComponentDialog';
import { useToast } from './ui/use-toast';

export default function FilterComponentActions() {
  const { toast } = useToast();

  const [showReplaceComponentDialog, setShowReplaceComponentDialog] = useState(false);
  const [withComment, setWithComment] = useState(false);
  const [filterBy, setFilterBy] = useState<'by_file' | 'by_purl'>();

  const onFilterComponent = useComponentFilterStore((state) => state.onFilterComponent);

  const handleFilterComponent = withErrorHandling({
    asyncFn: async (args: OnFilterComponentArgs) => {
      await onFilterComponent(args);
    },
    onError: () => {
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
    onError: () => {
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
          shortcutKeysByFileWithComments={KEYBOARD_SHORTCUTS.includeFileWithComments.keys}
          shortcutKeysByFileWithoutComments={KEYBOARD_SHORTCUTS.includeFileWithoutComments.keys}
          shortcutKeysByComponentWithComments={KEYBOARD_SHORTCUTS.includeComponentWithComments.keys}
          shortcutKeysByComponentWithoutComments={KEYBOARD_SHORTCUTS.includeComponentWithoutComments.keys}
        />
        <FilterActionButton
          action={FilterAction.Remove}
          description="Dismissing a file/component will exclude it from future scan results."
          icon={<PackageMinus className="h-5 w-5 stroke-red-500" />}
          onAdd={handleFilterComponent}
          shortcutKeysByFileWithComments={KEYBOARD_SHORTCUTS.dismissFileWithComments.keys}
          shortcutKeysByFileWithoutComments={KEYBOARD_SHORTCUTS.dismissFileWithoutComments.keys}
          shortcutKeysByComponentWithComments={KEYBOARD_SHORTCUTS.dismissComponentWithComments.keys}
          shortcutKeysByComponentWithoutComments={KEYBOARD_SHORTCUTS.dismissComponentWithoutComments.keys}
        />
        <FilterActionButton
          action={FilterAction.Replace}
          description="Replace detected components with another one."
          icon={<Replace className="h-5 w-5 stroke-yellow-500" />}
          onAdd={(args) => {
            setShowReplaceComponentDialog(true);
            setFilterBy(args.filterBy);
            setWithComment(args.withComment ?? false);
          }}
          shortcutKeysByFileWithComments={KEYBOARD_SHORTCUTS.replaceFileWithComments.keys}
          shortcutKeysByFileWithoutComments={KEYBOARD_SHORTCUTS.replaceFileWithoutComments.keys}
          shortcutKeysByComponentWithComments={KEYBOARD_SHORTCUTS.replaceComponentWithComments.keys}
          shortcutKeysByComponentWithoutComments={KEYBOARD_SHORTCUTS.replaceComponentWithoutComments.keys}
        />
      </div>
      {showReplaceComponentDialog && filterBy && (
        <ReplaceComponentDialog
          onOpenChange={() => setShowReplaceComponentDialog((prev) => !prev)}
          onReplaceComponent={handleReplaceComponent}
          withComment={withComment}
          filterBy={filterBy}
        />
      )}
    </>
  );
}
