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

import { useEffect, useMemo, useRef, useState } from 'react';

import { useDialogRegistration } from '@/contexts/DialogStateContext';
import useEnvironment from '@/hooks/useEnvironment';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { getShortcutDisplay, KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import { getPathSegments } from '@/lib/utils';
import { FilterAction, filterActionLabelMap } from '@/modules/components/domain';
import { OnFilterComponentArgs } from '@/modules/components/stores/useComponentFilterStore';

import { entities } from '../../wailsjs/go/models';
import PathBreadcrumb from './PathBreadcrumb';
import ReplaceComponentSection from './ReplaceComponentSection';
import SelectLicenseList from './SelectLicenseList';
import ShortcutBadge from './ShortcutBadge';
import { Button } from './ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from './ui/dialog';
import { Label } from './ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Textarea } from './ui/textarea';
import { useToast } from './ui/use-toast';

function buildPath(segments: string[], index: number): string {
  const path = segments.slice(0, index + 1).join('/');
  const isFile = index === segments.length - 1;
  return isFile ? path : path + '/';
}

interface FilterActionModalProps {
  action: FilterAction;
  filePath: string;
  purl: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: (args: OnFilterComponentArgs) => void;
}

export default function FilterActionModal({
  action,
  filePath,
  purl,
  open,
  onOpenChange,
  onConfirm,
}: FilterActionModalProps) {
  const { toast } = useToast();
  const { modifierKey } = useEnvironment();
  const submitButtonRef = useRef<HTMLButtonElement>(null);

  // Tab state
  const [activeTab, setActiveTab] = useState<'path' | 'component'>('path');

  // Path selection state
  const segments = useMemo(() => getPathSegments(filePath), [filePath]);
  const [selectedPathIndex, setSelectedPathIndex] = useState(segments.length - 1);

  // Form state
  const [comment, setComment] = useState('');
  const [license, setLicense] = useState<string>();
  const [licenseKey, setLicenseKey] = useState(0);
  const [selectedComponent, setSelectedComponent] = useState<entities.DeclaredComponent | null>(null);

  const isIncludeAction = action === FilterAction.Include;
  const isReplaceAction = action === FilterAction.Replace;

  // Determine target type for dynamic text
  const getTargetType = () => {
    if (activeTab === 'component') return 'component';
    return selectedPathIndex < segments.length - 1 ? 'folder' : 'file';
  };
  const targetType = getTargetType();

  useDialogRegistration('filter-action', open);

  // Reset state when modal opens
  useEffect(() => {
    if (open) {
      setActiveTab('path');
      setSelectedPathIndex(segments.length - 1);
      setComment('');
      setLicense(undefined);
      setLicenseKey((k) => k + 1);
      setSelectedComponent(null);
    }
  }, [open, segments.length]);

  const handleConfirm = () => {
    const selectedPath = buildPath(segments, selectedPathIndex);
    const isFileSelected = selectedPathIndex === segments.length - 1;

    // Determine filterBy based on tab and path selection
    let filterBy: 'by_file' | 'by_purl' | 'by_folder';
    if (activeTab === 'component') {
      filterBy = 'by_purl';
    } else if (isFileSelected) {
      filterBy = 'by_file';
    } else {
      filterBy = 'by_folder';
    }

    // Validate for Replace action
    if (isReplaceAction && !selectedComponent) {
      toast({
        title: 'Error',
        description: 'Please select a component to replace with.',
        variant: 'destructive',
      });
      return;
    }

    onConfirm({
      action,
      filterBy,
      comment: comment || undefined,
      license: license || undefined,
      replaceWith: selectedComponent?.purl,
      folderPath: filterBy === 'by_folder' ? selectedPath : undefined,
      purl,
    });

    onOpenChange(false);
  };

  const ref = useKeyboardShortcut(KEYBOARD_SHORTCUTS.confirm.keys, () => submitButtonRef.current?.click(), {
    enableOnFormTags: true,
    enabled: open,
  });

  const actionLabel = filterActionLabelMap[action];

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent ref={ref} tabIndex={-1} className="max-w-lg p-4">
        <DialogHeader>
          <DialogTitle>{actionLabel}</DialogTitle>
          <DialogDescription>
            {isReplaceAction
              ? `Replace the selected ${targetType} with another component.`
              : `${actionLabel} the selected ${targetType}.`}
          </DialogDescription>
        </DialogHeader>

        <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as 'path' | 'component')}>
          <TabsList className="mb-2 self-start">
            <TabsTrigger value="path">Path</TabsTrigger>
            <TabsTrigger value="component">Component</TabsTrigger>
          </TabsList>
          <TabsContent value="path" className="mt-4">
            <div className="flex min-h-[44px] items-center rounded-md border bg-muted/30 px-1.5 py-2">
              <PathBreadcrumb filePath={filePath} selectedIndex={selectedPathIndex} onSelect={setSelectedPathIndex} />
            </div>
          </TabsContent>
          <TabsContent value="component" className="mt-4">
            <div className="flex min-h-[44px] items-center rounded-md border bg-muted/30 px-3 py-2 font-mono text-sm">
              {purl}
            </div>
          </TabsContent>
        </Tabs>

        {/* License selector - Include action only */}
        {isIncludeAction && (
          <div className="flex flex-col gap-2">
            <Label>
              License <span className="text-xs font-normal text-muted-foreground">(optional)</span>
            </Label>
            <SelectLicenseList key={licenseKey} onSelect={setLicense} />
          </div>
        )}

        {/* Component selector - Replace action only */}
        {isReplaceAction && (
          <ReplaceComponentSection
            licenseKey={licenseKey}
            onComponentChange={setSelectedComponent}
            onLicenseChange={setLicense}
          />
        )}

        {/* Comment textarea - always shown */}
        <div className="flex flex-col gap-2">
          <Label>
            Comment <span className="text-xs font-normal text-muted-foreground">(optional)</span>
          </Label>
          <Textarea value={comment} onChange={(e) => setComment(e.target.value)} placeholder="Add a comment..." />
        </div>

        <DialogFooter>
          <Button variant="ghost" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button ref={submitButtonRef} onClick={handleConfirm}>
            {actionLabel} <ShortcutBadge shortcut={getShortcutDisplay(KEYBOARD_SHORTCUTS.confirm.keys, modifierKey.label)[0]} />
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
