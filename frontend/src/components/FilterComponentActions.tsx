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
import { useCallback, useState } from 'react';

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

  // === Include handlers ===
  const handleIncludeFile = useCallback(() => {
    if (!isCompletedResult && selectedResult) {
      handleFilterComponent({
        action: FilterAction.Include,
        filterBy: 'by_file',
        purl: selectedResult.detected_purl ?? '',
      });
    }
  }, [isCompletedResult, selectedResult, handleFilterComponent]);

  const handleIncludeFolder = useCallback(() => {
    if (!isCompletedResult && selectedResult) {
      setFilterModalAction(FilterAction.Include);
      setFilterModalInitialSelection('folder');
      setFilterModalOpen(true);
    }
  }, [isCompletedResult, selectedResult]);

  const handleIncludeComponent = useCallback(() => {
    if (!isCompletedResult && selectedResult) {
      setFilterModalAction(FilterAction.Include);
      setFilterModalInitialSelection('component');
      setFilterModalOpen(true);
    }
  }, [isCompletedResult, selectedResult]);

  // === Dismiss handlers ===
  const handleDismissFile = useCallback(() => {
    if (!isCompletedResult && selectedResult) {
      handleFilterComponent({
        action: FilterAction.Remove,
        filterBy: 'by_file',
        purl: selectedResult.detected_purl ?? '',
      });
    }
  }, [isCompletedResult, selectedResult, handleFilterComponent]);

  const handleDismissFolder = useCallback(() => {
    if (!isCompletedResult && selectedResult) {
      setFilterModalAction(FilterAction.Remove);
      setFilterModalInitialSelection('folder');
      setFilterModalOpen(true);
    }
  }, [isCompletedResult, selectedResult]);

  const handleDismissComponent = useCallback(() => {
    if (!isCompletedResult && selectedResult) {
      setFilterModalAction(FilterAction.Remove);
      setFilterModalInitialSelection('component');
      setFilterModalOpen(true);
    }
  }, [isCompletedResult, selectedResult]);

  // === Replace handlers ===
  const handleReplaceFile = useCallback(() => {
    if (!isCompletedResult && selectedResult) {
      setFilterModalAction(FilterAction.Replace);
      setFilterModalInitialSelection('file');
      setFilterModalOpen(true);
    }
  }, [isCompletedResult, selectedResult]);

  const handleReplaceFolder = useCallback(() => {
    if (!isCompletedResult && selectedResult) {
      setFilterModalAction(FilterAction.Replace);
      setFilterModalInitialSelection('folder');
      setFilterModalOpen(true);
    }
  }, [isCompletedResult, selectedResult]);

  const handleReplaceComponent = useCallback(() => {
    if (!isCompletedResult && selectedResult) {
      setFilterModalAction(FilterAction.Replace);
      setFilterModalInitialSelection('component');
      setFilterModalOpen(true);
    }
  }, [isCompletedResult, selectedResult]);

  // === Skip handlers ===
  const handleSkipFile = useCallback(() => {
    if (selectedResult) {
      setSkipModalInitialSelection('file');
      setSkipModalOpen(true);
    }
  }, [selectedResult]);

  const handleSkipFolder = useCallback(() => {
    if (selectedResult) {
      setSkipModalInitialSelection('folder');
      setSkipModalOpen(true);
    }
  }, [selectedResult]);

  const handleSkipExtension = useCallback(() => {
    if (selectedResult) {
      setSkipModalInitialSelection('extension');
      setSkipModalOpen(true);
    }
  }, [selectedResult]);

  // === Keyboard shortcuts ===
  // Include
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.include.keys, handleIncludeFile, {
    enabled: !isCompletedResult && !!selectedResult,
  });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.includeWithModal.keys, () => {
    if (!isCompletedResult && selectedResult) {
      setFilterModalAction(FilterAction.Include);
      setFilterModalInitialSelection('file');
      setFilterModalOpen(true);
    }
  }, { enabled: !isCompletedResult && !!selectedResult });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.includeFolder?.keys ?? '', handleIncludeFolder, {
    enabled: !isCompletedResult && !!selectedResult,
  });

  // Dismiss
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.dismiss.keys, handleDismissFile, {
    enabled: !isCompletedResult && !!selectedResult,
  });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.dismissWithModal.keys, () => {
    if (!isCompletedResult && selectedResult) {
      setFilterModalAction(FilterAction.Remove);
      setFilterModalInitialSelection('file');
      setFilterModalOpen(true);
    }
  }, { enabled: !isCompletedResult && !!selectedResult });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.dismissFolder?.keys ?? '', handleDismissFolder, {
    enabled: !isCompletedResult && !!selectedResult,
  });

  // Replace
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.replace.keys, handleReplaceFile, {
    enabled: !isCompletedResult && !!selectedResult,
  });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.replaceFolder?.keys ?? '', handleReplaceFolder, {
    enabled: !isCompletedResult && !!selectedResult,
  });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.replaceComponent?.keys ?? '', handleReplaceComponent, {
    enabled: !isCompletedResult && !!selectedResult,
  });

  // Skip
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.skip.keys, handleSkipFile, {
    enabled: !!selectedResult,
  });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.skipFolder?.keys ?? '', handleSkipFolder, {
    enabled: !!selectedResult,
  });
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.skipExtension?.keys ?? '', handleSkipExtension, {
    enabled: !!selectedResult,
  });

  // === Menu bar events ===
  useMenuEvents({
    // Include
    [entities.Action.Include]: handleIncludeFile,
    [entities.Action.IncludeWithModal]: () => {
      if (!isCompletedResult && selectedResult) {
        setFilterModalAction(FilterAction.Include);
        setFilterModalInitialSelection('file');
        setFilterModalOpen(true);
      }
    },
    [entities.Action.IncludeFolder]: handleIncludeFolder,
    // Dismiss
    [entities.Action.Dismiss]: handleDismissFile,
    [entities.Action.DismissWithModal]: () => {
      if (!isCompletedResult && selectedResult) {
        setFilterModalAction(FilterAction.Remove);
        setFilterModalInitialSelection('file');
        setFilterModalOpen(true);
      }
    },
    [entities.Action.DismissFolder]: handleDismissFolder,
    // Replace
    [entities.Action.Replace]: handleReplaceFile,
    [entities.Action.ReplaceFolder]: handleReplaceFolder,
    [entities.Action.ReplaceComponent]: handleReplaceComponent,
    // Skip
    [entities.Action.Skip]: handleSkipFile,
    [entities.Action.SkipFolder]: handleSkipFolder,
    [entities.Action.SkipExtension]: handleSkipExtension,
  });

  const handleFilterConfirm = async (args: OnFilterComponentArgs) => {
    await handleFilterComponent(args);
  };

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
            <MenubarItem onSelect={handleIncludeFile}>
              File
              <MenubarShortcut>I</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handleIncludeFolder}>
              Folder
              <MenubarShortcut>Alt+Shift+I</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handleIncludeComponent}>
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
            <MenubarItem onSelect={handleDismissFile}>
              File
              <MenubarShortcut>D</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handleDismissFolder}>
              Folder
              <MenubarShortcut>Alt+Shift+D</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handleDismissComponent}>
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
            <MenubarItem onSelect={handleReplaceFile}>
              File
              <MenubarShortcut>R</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handleReplaceFolder}>
              Folder
              <MenubarShortcut>Alt+Shift+R</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handleReplaceComponent}>
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
            <MenubarItem onSelect={handleSkipFile}>
              File
              <MenubarShortcut>S</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handleSkipFolder}>
              Folder
              <MenubarShortcut>Alt+Shift+S</MenubarShortcut>
            </MenubarItem>
            <MenubarItem onSelect={handleSkipExtension}>
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
          onConfirm={handleFilterConfirm}
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
