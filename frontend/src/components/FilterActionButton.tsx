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

import { useMemo, useState } from 'react';

import { Button } from '@/components/ui/button';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { useMenuEvents } from '@/hooks/useMenuEvent';
import useSelectedResult from '@/hooks/useSelectedResult';
import { FilterAction, filterActionLabelMap } from '@/modules/components/domain';
import { OnFilterComponentArgs } from '@/modules/components/stores/useComponentFilterStore';

import { entities } from '../../wailsjs/go/models';
import FilterActionModal from './FilterActionModal';

interface FilterActionProps {
  action: FilterAction;
  description: string;
  icon: React.ReactNode;
  onAdd: (args: OnFilterComponentArgs) => Promise<void> | void;
  directShortcutKeys?: string;
  modalShortcutKeys: string;
}

export default function FilterActionButton({
  action,
  description,
  icon,
  onAdd,
  directShortcutKeys,
  modalShortcutKeys,
}: FilterActionProps) {
  const selectedResult = useSelectedResult();
  const isCompletedResult = selectedResult?.workflow_state === 'completed';
  const [modalOpen, setModalOpen] = useState(false);

  const handleOpenModal = () => {
    if (!isCompletedResult && selectedResult) {
      setModalOpen(true);
    }
  };

  const handleDirectAction = () => {
    if (!isCompletedResult && selectedResult) {
      onAdd({
        action,
        filterBy: 'by_file',
        purl: selectedResult.detected_purl ?? '',
      });
    }
  };

  const handleConfirm = async (args: OnFilterComponentArgs) => {
    await onAdd(args);
  };

  // Direct shortcut (F1/i, F2/d) - executes action immediately on file-level
  useKeyboardShortcut(directShortcutKeys ?? '', handleDirectAction, {
    enabled: !isCompletedResult && !!selectedResult && !!directShortcutKeys,
  });

  // Modal shortcut (shift+F1/shift+i, shift+F2/shift+d, F3/r) - opens modal
  useKeyboardShortcut(modalShortcutKeys, handleOpenModal, {
    enabled: !isCompletedResult && !!selectedResult,
  });

  // Listen for menu bar events
  const eventHandlerMap = useMemo(
    () => ({
      // Include actions
      [entities.Action.Include]: action === FilterAction.Include ? handleDirectAction : null,
      [entities.Action.IncludeWithModal]: action === FilterAction.Include ? handleOpenModal : null,
      // Dismiss actions (mapped to FilterAction.Remove)
      [entities.Action.Dismiss]: action === FilterAction.Remove ? handleDirectAction : null,
      [entities.Action.DismissWithModal]: action === FilterAction.Remove ? handleOpenModal : null,
      // Replace action (always opens modal)
      [entities.Action.Replace]: action === FilterAction.Replace ? handleOpenModal : null,
    }),
    [action, handleDirectAction, handleOpenModal]
  );
  useMenuEvents(eventHandlerMap);

  return (
    <>
      <Button
        variant="ghost"
        size="lg"
        className="h-full w-14 rounded-none enabled:hover:bg-accent enabled:hover:text-accent-foreground"
        disabled={isCompletedResult}
        onClick={handleOpenModal}
        title={description}
      >
        <div className="flex flex-col items-center justify-center gap-1">
          <span className="text-xs">{filterActionLabelMap[action]}</span>
          {icon}
        </div>
      </Button>

      {selectedResult && (
        <FilterActionModal
          action={action}
          filePath={selectedResult.path}
          purl={selectedResult.detected_purl ?? ''}
          open={modalOpen}
          onOpenChange={setModalOpen}
          onConfirm={handleConfirm}
        />
      )}
    </>
  );
}
