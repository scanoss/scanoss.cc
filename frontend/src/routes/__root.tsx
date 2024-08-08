import { createRootRoute, Outlet } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/router-devtools';
import React from 'react';

import Header from '@/components/Header';
import Sidebar from '@/components/Sidebar';
import FileService from '@/modules/files/infra/service';

export const Route = createRootRoute({
  loader: () => FileService.getAll(),
  component: () => (
    <div className="grid h-screen w-full grid-cols-[330px_1fr] grid-rows-[60px_1fr] overflow-hidden">
      <Header />
      <Sidebar />
      <div className="flex flex-col p-6">
        <Outlet />
        <TanStackRouterDevtools />
      </div>
    </div>
  ),
});
