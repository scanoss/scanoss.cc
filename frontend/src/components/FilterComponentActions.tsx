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

import { Check, EyeOff, PackageMinus, Replace } from 'lucide-react';
import { useCallback, useMemo, useState } from 'react';

import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { useMenuEvents } from '@/hooks/useMenuEvent';
import { useResults } from '@/hooks/useResults';
import useSelectedResult from '@/hooks/useSelectedResult';
import { withErrorHandling } from '@/lib/errors';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import { FilterAction } from '@/modules/components/domain';
import useComponentFilterStore, { OnFilterComponentArgs } from '@/modules/components/stores/useComponentFilterStore';

import { entities } from '../../wailsjs/go/models';
import FilterActionModal from './FilterActionModal';
import SkipActionModal from './SkipActionModal';
import {
  Menubar,
  MenubarContent,
  MenubarItem,
  MenubarMenu,
  MenubarSeparator,
  MenubarShortcut,
  MenubarTrigger,
} from './ui/menubar';
import { useToast } from './ui/use-toast';

type FilterInitialSelection = 'file' | 'folder' | 'component';
type SkipInitialSelection = 'file' | 'folder' | 'extension';

export default function FilterComponentActions() {
  const { toast } = useToast();
  const { reset } = useResults();
  const selectedResult = useSelectedResult();
  const isCompletedResult = selectedResult?.workflow_state === 'completed';
  const onFilterComponent = useComponentFilterStore((state) => state.onFilterComponent);

  // Filter action modal state
  const [filterModalOpen, setFilterModalOpen] = useState(false);
  const [filterModalAction, setFilterModalAction] = useState<FilterAction>(FilterAction.Include);
  const [filterModalInitialSelection, setFilterModalInitialSelection] = useState<FilterInitialSelection>('file');

  // Skip action modal state
  const [skipModalOpen, setSkipModalOpen] = useState(false);
  const [skipModalInitialSelection, setSkipModalInitialSelection] = useState<SkipInitialSelection>('file');

  const handleFilterComponent = withErrorHandling({
    asyncFn: async (args: OnFilterComponentArgs) => {
      await onFilterComponent(args);
      reset();
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'An error occurred while filtering the component. Please try again.',
        variant: 'destructive',
      });
    },
  });

  // Creates handler that opens filter modal with given action and selection
  const createModalActionHandler = useCallback(
    (action: FilterAction, selection: FilterInitialSelection) => () => {
      if (!selectedResult || isCompletedResult) return;
      setFilterModalAction(action);
      setFilterModalInitialSelection(selection);
      setFilterModalOpen(true);
    },
    [selectedResult, isCompletedResult]
  );

  // Creates handler that applies filter action directly to the current file (no modal)
  const createDirectActionHandler = useCallback(
    (action: FilterAction) => () => {
      if (!selectedResult || isCompletedResult) return;
      handleFilterComponent({
        action,
        filterBy: 'by_file',
        purl: selectedResult.detected_purl ?? '',
      });
    },
    [selectedResult, isCompletedResult, handleFilterComponent]
  );

  // Creates handler that opens skip modal with given selection
  const createModalSkipHandler = useCallback(
    (selection: SkipInitialSelection) => () => {
      if (!selectedResult) return;
      setSkipModalInitialSelection(selection);
      setSkipModalOpen(true);
    },
    [selectedResult]
  );

  // Generate all handlers
  const handlers = useMemo(
    () => ({
      // Include: file applies directly, others open modal
      includeFile: createDirectActionHandler(FilterAction.Include),
      includeFolder: createModalActionHandler(FilterAction.Include, 'folder'),
      includeComponent: createModalActionHandler(FilterAction.Include, 'component'),
      includeWithModal: createModalActionHandler(FilterAction.Include, 'file'),

      // Dismiss: file applies directly, others open modal
      dismissFile: createDirectActionHandler(FilterAction.Remove),
      dismissFolder: createModalActionHandler(FilterAction.Remove, 'folder'),
      dismissComponent: createModalActionHandler(FilterAction.Remove, 'component'),
      dismissWithModal: createModalActionHandler(FilterAction.Remove, 'file'),

      // Replace: always opens modal (needs PURL selection)
      replaceFile: createModalActionHandler(FilterAction.Replace, 'file'),
      replaceFolder: createModalActionHandler(FilterAction.Replace, 'folder'),
      replaceComponent: createModalActionHandler(FilterAction.Replace, 'component'),

      // Skip: always opens modal
      skipFile: createModalSkipHandler('file'),
      skipFolder: createModalSkipHandler('folder'),
      skipExtension: createModalSkipHandler('extension'),
    }),
    [createDirectActionHandler, createModalActionHandler, createModalSkipHandler]
  );

  // === Keyboard shortcuts ===
  const filterEnabled = !isCompletedResult && !!selectedResult;
  const skipEnabled = !!selectedResult;

  // Include
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.include.keys, handlers.includeFile, { enabled: filterEnabled });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.includeWithModal.keys, handlers.includeWithModal, { enabled: filterEnabled });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.includeFolder.keys, handlers.includeFolder, { enabled: filterEnabled });

  // Dismiss
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.dismiss.keys, handlers.dismissFile, { enabled: filterEnabled });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.dismissWithModal.keys, handlers.dismissWithModal, { enabled: filterEnabled });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.dismissFolder.keys, handlers.dismissFolder, { enabled: filterEnabled });

  // Replace
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.replace.keys, handlers.replaceFile, { enabled: filterEnabled });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.replaceFolder.keys, handlers.replaceFolder, { enabled: filterEnabled });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.replaceComponent.keys, handlers.replaceComponent, { enabled: filterEnabled });

  // Skip
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.skip.keys, handlers.skipFile, { enabled: skipEnabled });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.skipFolder.keys, handlers.skipFolder, { enabled: skipEnabled });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.skipExtension.keys, handlers.skipExtension, { enabled: skipEnabled });

  // === Menu bar events ===
  useMenuEvents({
    // Include
    [entities.Action.Include]: handlers.includeFile,
    [entities.Action.IncludeWithModal]: handlers.includeWithModal,
    [entities.Action.IncludeFolder]: handlers.includeFolder,
    // Dismiss
    [entities.Action.Dismiss]: handlers.dismissFile,
    [entities.Action.DismissWithModal]: handlers.dismissWithModal,
    [entities.Action.DismissFolder]: handlers.dismissFolder,
    // Replace
    [entities.Action.Replace]: handlers.replaceFile,
    [entities.Action.ReplaceFolder]: handlers.replaceFolder,
    [entities.Action.ReplaceComponent]: handlers.replaceComponent,
    // Skip
    [entities.Action.Skip]: handlers.skipFile,
    [entities.Action.SkipFolder]: handlers.skipFolder,
    [entities.Action.SkipExtension]: handlers.skipExtension,
  });

  const isDisabled = isCompletedResult || !selectedResult;

  return (
    <>
      <Menubar className="h-auto border-0 bg-transparent p-0 shadow-none">
        {/* Include */}
        <MenubarMenu>
          <MenubarTrigger
            disabled={isDisabled}
            className="flex h-full w-14 flex-col items-center justify-center gap-1 rounded-none px-2 py-1 data-[state=open]:bg-accent"
          >
            <span className="text-xs">Include</span>
            <Check className="h-5 w-5 stroke-green-500" />
          </MenubarTrigger>
          <MenubarContent align="start" className="min-w-[180px]">
            <MenubarItem onSelect={handlers.includeFile}>
              File
              <MenubarShortcut>I</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handlers.includeFolder}>
              Folder
              <MenubarShortcut>Alt+Shift+I</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handlers.includeComponent}>
              Component
              <MenubarShortcut>Shift+I</MenubarShortcut>
            </MenubarItem>
          </MenubarContent>
        </MenubarMenu>

        {/* Dismiss */}
        <MenubarMenu>
          <MenubarTrigger
            disabled={isDisabled}
            className="flex h-full w-14 flex-col items-center justify-center gap-1 rounded-none px-2 py-1 data-[state=open]:bg-accent"
          >
            <span className="text-xs">Dismiss</span>
            <PackageMinus className="h-5 w-5 stroke-red-500" />
          </MenubarTrigger>
          <MenubarContent align="start" className="min-w-[180px]">
            <MenubarItem onSelect={handlers.dismissFile}>
              File
              <MenubarShortcut>D</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handlers.dismissFolder}>
              Folder
              <MenubarShortcut>Alt+Shift+D</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handlers.dismissComponent}>
              Component
              <MenubarShortcut>Shift+D</MenubarShortcut>
            </MenubarItem>
          </MenubarContent>
        </MenubarMenu>

        {/* Replace */}
        <MenubarMenu>
          <MenubarTrigger
            disabled={isDisabled}
            className="flex h-full w-14 flex-col items-center justify-center gap-1 rounded-none px-2 py-1 data-[state=open]:bg-accent"
          >
            <span className="text-xs">Replace</span>
            <Replace className="h-5 w-5 stroke-yellow-500" />
          </MenubarTrigger>
          <MenubarContent align="start" className="min-w-[180px]">
            <MenubarItem onSelect={handlers.replaceFile}>
              File
              <MenubarShortcut>R</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handlers.replaceFolder}>
              Folder
              <MenubarShortcut>Alt+Shift+R</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handlers.replaceComponent}>
              Component
              <MenubarShortcut>Shift+R</MenubarShortcut>
            </MenubarItem>
          </MenubarContent>
        </MenubarMenu>

        {/* Separator */}
        <MenubarSeparator className="mx-1 h-8 w-px" />

        {/* Skip */}
        <MenubarMenu>
          <MenubarTrigger
            disabled={!selectedResult}
            className="flex h-full w-14 flex-col items-center justify-center gap-1 rounded-none px-2 py-1 data-[state=open]:bg-accent"
          >
            <span className="text-xs">Skip</span>
            <EyeOff className="h-5 w-5 stroke-orange-500" />
          </MenubarTrigger>
          <MenubarContent align="start" className="min-w-[180px]">
            <MenubarItem onSelect={handlers.skipFile}>
              File
              <MenubarShortcut>S</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handlers.skipFolder}>
              Folder
              <MenubarShortcut>Alt+Shift+S</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handlers.skipExtension}>
              Extension
              <MenubarShortcut>Shift+S</MenubarShortcut>
            </MenubarItem>
          </MenubarContent>
        </MenubarMenu>
      </Menubar>

      {/* Filter Action Modal */}
      {selectedResult && (
        <FilterActionModal
          action={filterModalAction}
          filePath={selectedResult.path}
          purl={selectedResult.detected_purl ?? ''}
          open={filterModalOpen}
          onOpenChange={setFilterModalOpen}
          onConfirm={handleFilterComponent}
          initialSelection={filterModalInitialSelection}
        />
      )}

      {/* Skip Action Modal */}
      {selectedResult && (
        <SkipActionModal
          filePath={selectedResult.path}
          open={skipModalOpen}
          onOpenChange={setSkipModalOpen}
          initialSelection={skipModalInitialSelection}
        />
      )}
    </>
  );
}
