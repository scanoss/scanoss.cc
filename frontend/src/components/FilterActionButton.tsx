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
  shortcutKeysByFileWithComments: string;
  shortcutKeysByFileWithoutComments: string;
  shortcutKeysByComponentWithComments: string;
  shortcutKeysByComponentWithoutComments: string;
}

export default function FilterActionButton({
  action,
  description,
  icon,
  onAdd,
  shortcutKeysByFileWithComments,
  shortcutKeysByFileWithoutComments,
  shortcutKeysByComponentWithComments,
  shortcutKeysByComponentWithoutComments,
}: FilterActionProps) {
  const selectedResult = useSelectedResult();
  const isCompletedResult = selectedResult?.workflow_state === 'completed';
  const [modalOpen, setModalOpen] = useState(false);

  const handleOpenModal = () => {
    if (!isCompletedResult) {
      setModalOpen(true);
    }
  };

  const handleConfirm = async (args: OnFilterComponentArgs) => {
    await onAdd(args);
  };

  // Keyboard shortcuts all open the modal
  useKeyboardShortcut(shortcutKeysByFileWithoutComments, handleOpenModal, {
    enabled: !isCompletedResult,
  });
  useKeyboardShortcut(shortcutKeysByComponentWithoutComments, handleOpenModal, {
    enabled: !isCompletedResult,
  });
  useKeyboardShortcut(shortcutKeysByFileWithComments, handleOpenModal, {
    enabled: !isCompletedResult,
  });
  useKeyboardShortcut(shortcutKeysByComponentWithComments, handleOpenModal, {
    enabled: !isCompletedResult,
  });

  // Listen for menu bar events and map them to open the modal
  const eventHandlerMap = useMemo(
    () => ({
      // Include actions
      [entities.Action.IncludeFileWithoutComments]: action === FilterAction.Include ? handleOpenModal : null,
      [entities.Action.IncludeFileWithComments]: action === FilterAction.Include ? handleOpenModal : null,
      [entities.Action.IncludeComponentWithoutComments]: action === FilterAction.Include ? handleOpenModal : null,
      [entities.Action.IncludeComponentWithComments]: action === FilterAction.Include ? handleOpenModal : null,
      // Dismiss actions (mapped to FilterAction.Remove)
      [entities.Action.DismissFileWithoutComments]: action === FilterAction.Remove ? handleOpenModal : null,
      [entities.Action.DismissFileWithComments]: action === FilterAction.Remove ? handleOpenModal : null,
      [entities.Action.DismissComponentWithoutComments]: action === FilterAction.Remove ? handleOpenModal : null,
      [entities.Action.DismissComponentWithComments]: action === FilterAction.Remove ? handleOpenModal : null,
      // Replace actions
      [entities.Action.ReplaceFileWithoutComments]: action === FilterAction.Replace ? handleOpenModal : null,
      [entities.Action.ReplaceFileWithComments]: action === FilterAction.Replace ? handleOpenModal : null,
      [entities.Action.ReplaceComponentWithoutComments]: action === FilterAction.Replace ? handleOpenModal : null,
      [entities.Action.ReplaceComponentWithComments]: action === FilterAction.Replace ? handleOpenModal : null,
    }),
    [action, handleOpenModal]
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
