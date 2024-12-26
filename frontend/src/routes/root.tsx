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
    // Register event listeners for shortcuts
    EventsOn(entities.Action.ShowKeyboardShortcutsModal, () => {
      setShowKeyboardShortcuts(true);
    });
    EventsOn(entities.Action.ScanWithOptions, () => {
      console.log('triggered scan with options');
      handleShowScanModal({ withOptions: true });
    });
    EventsOn(entities.Action.ScanCurrentDirectory, () => {
      console.log('triggered scan without options');
      handleShowScanModal({ withOptions: false });
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
      {scanModal.open && <ScanDialog onOpenChange={handleCloseScanModal} withOptions={scanModal.withOptions} />}
    </div>
  );
}
