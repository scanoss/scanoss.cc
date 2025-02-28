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

import { useEffect, useState } from 'react';
import { Outlet } from 'react-router-dom';

import KeyboardShortcutsDialog from '@/components/KeyboardShortcutsDialog';
import ScanDialog from '@/components/ScanDialog';
import Sidebar from '@/components/Sidebar';
import StatusBar from '@/components/StatusBar';
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from '@/components/ui/resizable';

import { entities } from '../../wailsjs/go/models';
import { EventsOn } from '../../wailsjs/runtime/runtime';

export default function Root() {
  const [showKeyboardShortcuts, setShowKeyboardShortcuts] = useState(false);
  const [showScanModal, setShowScanModal] = useState(false);

  useEffect(() => {
    // Register event listeners
    EventsOn(entities.Action.ShowKeyboardShortcutsModal, () => {
      setShowKeyboardShortcuts(true);
    });
    EventsOn(entities.Action.ScanWithOptions, () => {
      setShowScanModal(true);
    });
  }, []);

  return (
    <div className="flex h-screen w-full flex-col bg-background backdrop-blur-lg">
      <div className="min-h-0 flex-1">
        <ResizablePanelGroup direction="horizontal" className="h-full" autoSaveId="panels-layout">
          <ResizablePanel defaultSize={30}>
            <Sidebar />
          </ResizablePanel>
          <ResizableHandle />
          <ResizablePanel>
            <div className="h-full overflow-auto">
              <Outlet />
            </div>
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
      <div className="border-t">
        <StatusBar />
      </div>
      <KeyboardShortcutsDialog open={showKeyboardShortcuts} onOpenChange={() => setShowKeyboardShortcuts(false)} />
      <ScanDialog open={showScanModal} onOpenChange={() => setShowScanModal(false)} />
    </div>
  );
}
