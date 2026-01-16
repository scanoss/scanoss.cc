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
import {
  AddStagedScanningSkipPattern,
  CommitStagedScanningSkipPatterns,
} from '../../wailsjs/go/service/ScanossSettingsServiceImp';
import SkipActionModal from './SkipActionModal';
import { Button } from './ui/button';
import { useToast } from './ui/use-toast';

export default function SkipActionButton() {
  const { toast } = useToast();
  const selectedResult = useSelectedResult();
  const [modalOpen, setModalOpen] = useState(false);

  const handleOpenModal = useCallback(() => setModalOpen(true), []);

  const handleDirectSkip = useCallback(async () => {
    if (!selectedResult) return;

    try {
      await AddStagedScanningSkipPattern(selectedResult.path);
      await CommitStagedScanningSkipPatterns();
      toast({
        title: 'File skipped',
        description: `"${selectedResult.path}" will be skipped in future scans.`,
      });
    } catch {
      toast({
        title: 'Error',
        description: 'Failed to add skip pattern. Please try again.',
        variant: 'destructive',
      });
    }
  }, [selectedResult, toast]);

  useKeyboardShortcut(KEYBOARD_SHORTCUTS.skip.keys, handleDirectSkip, {
    enabled: !!selectedResult,
  });

  useKeyboardShortcut(KEYBOARD_SHORTCUTS.skipWithModal.keys, handleOpenModal, {
    enabled: !!selectedResult,
  });

  useMenuEvents({
    [entities.Action.Skip]: handleDirectSkip,
    [entities.Action.SkipWithModal]: handleOpenModal,
  });

  return (
    <>
      <Button
        variant="ghost"
        size="lg"
        className="h-full w-14 rounded-none enabled:hover:bg-accent enabled:hover:text-accent-foreground"
        disabled={!selectedResult}
        onClick={handleOpenModal}
        title="Skip file/folder/extension from scanning"
      >
        <div className="flex flex-col items-center justify-center gap-1">
          <span className="text-xs">Skip</span>
          <EyeOff className="h-5 w-5 stroke-orange-500" />
        </div>
      </Button>

      {selectedResult && (
        <SkipActionModal
          filePath={selectedResult.path}
          open={modalOpen}
          onOpenChange={setModalOpen}
        />
      )}
    </>
  );
}
