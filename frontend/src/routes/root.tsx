import { useEffect, useState } from 'react';
import { Outlet } from 'react-router-dom';

import KeyboardShortcutsDialog from '@/components/KeyboardShortcutsDialog';
import ScanDialog from '@/components/ScanDialog';
import Sidebar from '@/components/Sidebar';
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from '@/components/ui/resizable';

import { entities } from '../../wailsjs/go/models';
import { EventsOn } from '../../wailsjs/runtime/runtime';

export default function Root() {
  const [showKeyboardShortcuts, setShowKeyboardShortcuts] = useState(false);
  const [showScanModal, setShowScanModal] = useState(false);

  useEffect(() => {
    // Register event listeners for shortcuts
    EventsOn(entities.Action.ShowKeyboardShortcutsModal, () => {
      setShowKeyboardShortcuts(true);
    });
    EventsOn(entities.Action.ScanWithOptions, () => {
      setShowScanModal(true);
    });
  }, []);

  return (
    <div className="h-screen w-full bg-background backdrop-blur-lg">
      <ResizablePanelGroup direction="horizontal" autoSaveId="panels-layout">
        <ResizablePanel defaultSize={30}>
          <Sidebar />
        </ResizablePanel>
        <ResizableHandle />
        <ResizablePanel>
          <Outlet />
        </ResizablePanel>
      </ResizablePanelGroup>
      <KeyboardShortcutsDialog open={showKeyboardShortcuts} onOpenChange={() => setShowKeyboardShortcuts(false)} />
      <ScanDialog open={showScanModal} onOpenChange={() => setShowScanModal(false)} />
    </div>
  );
}
