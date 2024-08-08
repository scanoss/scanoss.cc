import React from 'react';
import { Outlet } from 'react-router-dom';

import Sidebar from '@/components/Sidebar';

export default function Root() {
  return (
    <div className="flex h-screen w-full overflow-hidden">
      <div className="w-[330px] h-full">
        <Sidebar />
      </div>
      <div className="flex flex-col p-6 flex-1 bg-muted">
        <Outlet />
      </div>
    </div>
  );
}
