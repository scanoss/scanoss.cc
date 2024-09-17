import { Outlet } from 'react-router-dom';

import Sidebar from '@/components/Sidebar';
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from '@/components/ui/resizable';

export default function Root() {
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
    </div>
  );
}
