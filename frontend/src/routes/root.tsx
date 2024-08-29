import { Outlet } from 'react-router-dom';

import Sidebar from '@/components/Sidebar';

export default function Root() {
  return (
    <div className="flex h-screen w-full overflow-hidden">
      <div className="h-full w-[330px]">
        <Sidebar />
      </div>
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  );
}
