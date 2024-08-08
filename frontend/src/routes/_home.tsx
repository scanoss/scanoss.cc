import { createFileRoute, Outlet } from '@tanstack/react-router';
import React from 'react';

import Sidebar from '@/components/Sidebar';
import FileService from '@/modules/files/infra/service';

export const Route = createFileRoute('/_home')({
  loader: FileService.getAllFilesToBeCommited,
  component: Layout,
});

function Layout() {
  const filesToBeCommited = Route.useLoaderData();

  return (
    <div className="flex h-screen w-full overflow-hidden">
      <div className="w-[330px] h-full">
        <Sidebar files={filesToBeCommited} />
      </div>
      <div className="flex flex-col p-6 flex-1">
        <Outlet />
      </div>
    </div>
  );
}
