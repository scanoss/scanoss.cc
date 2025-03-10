// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2025 SCANOSS.COM
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

import { Dialog, DialogContent } from '@/components/ui/dialog';
import { ScrollArea } from '@/components/ui/scroll-area';
import { SETTINGS_OPTIONS } from '@/lib/settings';

import SettingsSidebar from './SettingsSidebar';

interface SettingsDialogProps {
  open: boolean;
  onOpenChange: () => void;
}

export default function SettingsDialog({ open, onOpenChange }: SettingsDialogProps) {
  const [activeSetting, setActiveSetting] = useState<string>(SETTINGS_OPTIONS[0].id);

  const activeSettingTitle = useMemo(() => SETTINGS_OPTIONS.find((option) => option.id === activeSetting)?.title, [activeSetting]);
  const activeSettingComponent = useMemo(() => SETTINGS_OPTIONS.find((option) => option.id === activeSetting)?.component, [activeSetting]);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="settings-dialog h-[calc(80dvh)] w-[calc(70dvw)] max-w-none">
        <div className="flex h-full overflow-hidden">
          <div className="w-1/4 border-r">
            <SettingsSidebar activeSetting={activeSetting} setActiveSetting={setActiveSetting} />
          </div>
          <div className="flex h-full flex-1 flex-col">
            <div className="flex h-12 w-full items-center border-b bg-background px-4">
              <span className="text-sm">{activeSettingTitle}</span>
            </div>
            <div className="flex-1 overflow-hidden">
              <ScrollArea className="h-full">
                <div className="h-full p-4">{activeSettingComponent}</div>
              </ScrollArea>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
