// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
  const [scanModal, setScanModal] = useState({
    open: false,
    withOptions: false,
  });

  const handleCloseScanModal = () => {
    setScanModal({
      open: false,
      withOptions: false,
    });
  };

  const handleShowScanModal = ({ withOptions }: { withOptions: boolean }) => {
    setScanModal({
      open: true,
      withOptions,
    });
  };

  useEffect(() => {
    // Register event listeners
    EventsOn(entities.Action.ShowKeyboardShortcutsModal, () => {
      setShowKeyboardShortcuts(true);
    });
    EventsOn(entities.Action.ScanWithOptions, () => {
      handleShowScanModal({ withOptions: true });
    });
    EventsOn(entities.Action.ScanCurrentDirectory, () => {
      handleShowScanModal({ withOptions: false });
    });
  }, []);

  return (
    <div className="flex h-screen w-full flex-col bg-background backdrop-blur-lg">
      <div className="min-h-0 flex-1">
        <ResizablePanelGroup direction="horizontal" className="h-full" autoSaveId="panels-layout">
          <ResizablePanel defaultSize={30}>
            <div className="h-full overflow-auto">
              <Sidebar />
            </div>
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
      {scanModal.open && <ScanDialog onOpenChange={handleCloseScanModal} withOptions={scanModal.withOptions} />}
    </div>
  );
}
