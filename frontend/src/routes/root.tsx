import React from 'react';
import { Outlet } from 'react-router-dom';

import Sidebar from '@/components/Sidebar';

export default function Root() {
  return (
    <div className="flex h-screen w-full overflow-hidden">
      <div className="h-full w-[330px]">
        <Sidebar />
      </div>
      <div className="flex flex-1 flex-col p-6">
        <Outlet />
      </div>
    </div>
  );
}
