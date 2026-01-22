// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2026 SCANOSS.COM
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

import { AlertTriangle, Info } from 'lucide-react';
import { useEffect, useMemo, useRef, useState } from 'react';

import { useDialogRegistration } from '@/contexts/DialogStateContext';
import useEnvironment from '@/hooks/useEnvironment';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { getShortcutDisplay, KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import { getPathSegments } from '@/lib/utils';

import {
  AddStagedScanningSkipPattern,
  CommitStagedScanningSkipPatterns,
} from '../../wailsjs/go/service/ScanossSettingsServiceImp';
import PathBreadcrumb from './PathBreadcrumb';
import ShortcutBadge from './ShortcutBadge';
import { Button } from './ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from './ui/dialog';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { useToast } from './ui/use-toast';

function getFileExtension(filename: string): string | null {
  const dotIndex = filename.indexOf('.');
  if (dotIndex === -1 || dotIndex === 0) return null;
  return filename.substring(dotIndex + 1);
}

function buildPath(segments: string[], index: number): string {
  const path = segments.slice(0, index + 1).join('/');
  const isFile = index === segments.length - 1;
  return isFile ? path : path + '/';
}

type InitialSelection = 'file' | 'folder';

interface SkipActionModalProps {
  filePath: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  initialSelection?: InitialSelection;
}

export default function SkipActionModal({
  filePath,
  open,
  onOpenChange,
  initialSelection = 'file',
}: SkipActionModalProps) {
  const { toast } = useToast();
  const { modifierKey } = useEnvironment();
  const submitButtonRef = useRef<HTMLButtonElement>(null);

  // Parse path into segments
  const segments = useMemo(() => getPathSegments(filePath), [filePath]);
  const filename = segments[segments.length - 1] || '';
  const extension = getFileExtension(filename);

  // Tab state
  const [activeTab, setActiveTab] = useState<'path' | 'extension'>('path');

  // Path selection state (for Path tab)
  const [selectedPathIndex, setSelectedPathIndex] = useState(segments.length - 1);

  useDialogRegistration('skip-action', open);

  // Reset state when modal opens
  useEffect(() => {
    if (open) {
      setActiveTab('path');
      if (initialSelection === 'folder') {
        // Select parent folder (one level up from file)
        setSelectedPathIndex(Math.max(0, segments.length - 2));
      } else {
        setSelectedPathIndex(segments.length - 1);
      }
    }
  }, [open, segments.length, initialSelection]);

  // Determine target type for dynamic text
  const getTargetType = () => {
    if (activeTab === 'extension') return 'extension';
    return selectedPathIndex < segments.length - 1 ? 'folder' : 'file';
  };
  const targetType = getTargetType();

  const getSkipPattern = (): string => {
    if (activeTab === 'extension') {
      return `*.${extension}`;
    }
    return buildPath(segments, selectedPathIndex);
  };

  const getWarningMessage = (): string => {
    if (activeTab === 'extension') {
      return `All .${extension} files will be skipped from future scans in this project.`;
    }
    if (selectedPathIndex < segments.length - 1) {
      return 'This folder and all its contents will be skipped from future scans.';
    }
    return 'This file will be skipped from future scans.';
  };

  const handleConfirm = async () => {
    const pattern = getSkipPattern();

    try {
      await AddStagedScanningSkipPattern(pattern);
      await CommitStagedScanningSkipPatterns();

      toast({
        title: 'Pattern added',
        description: `"${pattern}" will be skipped in future scans.`,
      });

      onOpenChange(false);
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (e) {
      toast({
        title: 'Error',
        description: 'Failed to add skip pattern. Please try again.',
        variant: 'destructive',
      });
    }
  };

  const ref = useKeyboardShortcut(KEYBOARD_SHORTCUTS.confirm.keys, () => submitButtonRef.current?.click(), {
    enableOnFormTags: true,
    enabled: open,
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent ref={ref} tabIndex={-1} className="max-w-lg p-4">
        <DialogHeader>
          <DialogTitle>Skip from scanning</DialogTitle>
          <DialogDescription>
            Skip the selected {targetType} from future scans.
          </DialogDescription>
        </DialogHeader>

        <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as 'path' | 'extension')}>
          <TabsList className="mb-2 self-start">
            <TabsTrigger value="path">Path</TabsTrigger>
            {extension && <TabsTrigger value="extension">Extension</TabsTrigger>}
          </TabsList>

          <TabsContent value="path" className="mt-4">
            <div className="flex min-h-[44px] items-center rounded-md border bg-muted/30 px-1.5 py-2">
              <PathBreadcrumb
                filePath={filePath}
                selectedIndex={selectedPathIndex}
                onSelect={setSelectedPathIndex}
              />
            </div>
          </TabsContent>

          {extension && (
            <TabsContent value="extension" className="mt-4">
              <div className="flex min-h-[44px] items-center rounded-md border bg-muted/30 px-3 py-2 font-mono text-sm">
                *.{extension}
              </div>
            </TabsContent>
          )}
        </Tabs>

        {/* Warning note */}
        <div className="flex items-start gap-2 rounded-lg border bg-muted/50 px-3 py-2 text-sm text-muted-foreground">
          <Info className="h-4 w-4 mt-0.5 shrink-0" />
          <span>{getWarningMessage()}</span>
        </div>

        {/* Re-scan required note */}
        <div className="flex items-start gap-2 rounded-lg bg-amber-500/10 px-3 py-2 text-sm text-amber-600 dark:text-amber-500">
          <AlertTriangle className="h-4 w-4 mt-0.5 shrink-0" />
          <span>Changes will take effect on the next scan. Trigger a new scan to see the updated results.</span>
        </div>

        <DialogFooter>
          <Button variant="ghost" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button ref={submitButtonRef} onClick={handleConfirm}>
            Skip <ShortcutBadge shortcut={getShortcutDisplay(KEYBOARD_SHORTCUTS.confirm.keys, modifierKey.label)[0]} />
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
