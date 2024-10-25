// import { useEffect } from 'react';
import { useEffect, useState } from 'react';
import { Outlet } from 'react-router-dom';

import KeyboardShortcutsDialog from '@/components/keyboard-shortcuts-dialog';
import Sidebar from '@/components/Sidebar';
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from '@/components/ui/resizable';

import { EventsOn } from '../../wailsjs/runtime';

export default function Root() {
  const [showKeyboardShortcuts, setShowKeyboardShortcuts] = useState(false);

  useEffect(() => {
    EventsOn('showKeyboardShortcuts', () => {
      setShowKeyboardShortcuts(true);
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
    </div>
  );
}
