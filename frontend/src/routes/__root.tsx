import { createRootRoute, Outlet } from '@tanstack/react-router';
import React from 'react';

import { TooltipProvider } from '@/components/ui/tooltip';

export const Route = createRootRoute({
  component: Root,
});

function Root() {
  return (
    <TooltipProvider>
      <Outlet />
    </TooltipProvider>
  );
}
