// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2024 SCANOSS.COM
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

import { EyeOff } from 'lucide-react';
import { useCallback, useState } from 'react';

import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { useMenuEvents } from '@/hooks/useMenuEvent';
import useSelectedResult from '@/hooks/useSelectedResult';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';

import { entities } from '../../wailsjs/go/models';
import SkipActionModal from './SkipActionModal';
import { Button } from './ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from './ui/dropdown-menu';

type InitialSelection = 'file' | 'folder' | 'extension';

export default function SkipActionButton() {
  const selectedResult = useSelectedResult();
  const [modalOpen, setModalOpen] = useState(false);
  const [modalInitialSelection, setModalInitialSelection] = useState<InitialSelection>('file');

  const handleOpenModal = useCallback((initialSelection: InitialSelection = 'file') => {
    setModalInitialSelection(initialSelection);
    setModalOpen(true);
  }, []);

  const handleOpenModalFile = useCallback(() => handleOpenModal('file'), [handleOpenModal]);
  const handleOpenModalFolder = useCallback(() => handleOpenModal('folder'), [handleOpenModal]);
  const handleOpenModalExtension = useCallback(() => handleOpenModal('extension'), [handleOpenModal]);

  // Skip action always opens modal (like Replace)
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.skip.keys, handleOpenModalFile, {
    enabled: !!selectedResult,
  });

  // Folder shortcut (alt+shift+s) - opens modal with folder selected
  useKeyboardShortcut(KEYBOARD_SHORTCUTS.skipFolder?.keys ?? '', handleOpenModalFolder, {
    enabled: !!selectedResult,
  });

  useMenuEvents({
    [entities.Action.Skip]: handleOpenModalFile,
    [entities.Action.SkipFolder]: handleOpenModalFolder,
    [entities.Action.SkipExtension]: handleOpenModalExtension,
  });

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild disabled={!selectedResult}>
          <Button
            variant="ghost"
            size="lg"
            className="h-full w-14 rounded-none enabled:hover:bg-accent enabled:hover:text-accent-foreground data-[state=open]:bg-accent"
            disabled={!selectedResult}
          >
            <div className="flex flex-col items-center justify-center gap-1">
              <span className="text-xs">Skip</span>
              <EyeOff className="h-5 w-5 stroke-orange-500" />
            </div>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="start" className="min-w-[180px]">
          <DropdownMenuItem onClick={handleOpenModalFile}>
            File
            <DropdownMenuShortcut>S</DropdownMenuShortcut>
          </DropdownMenuItem>
          <DropdownMenuItem onClick={handleOpenModalFolder}>
            Folder
            <DropdownMenuShortcut>Alt+Shift+S</DropdownMenuShortcut>
          </DropdownMenuItem>
          <DropdownMenuItem onClick={handleOpenModalExtension}>
            Extension
            <DropdownMenuShortcut>Shift+S</DropdownMenuShortcut>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {selectedResult && (
        <SkipActionModal
          filePath={selectedResult.path}
          open={modalOpen}
          onOpenChange={setModalOpen}
          initialSelection={modalInitialSelection}
        />
      )}
    </>
  );
}
