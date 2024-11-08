import './style.css';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { createHashRouter, RouterProvider } from 'react-router-dom';

import MatchComparison from './components/MatchComparison';
import { Toaster } from './components/ui/toaster';
import { TooltipProvider } from './components/ui/tooltip';
import { DialogProvider } from './providers/DialogProvider';
import Root from './routes/root';

const queryClient = new QueryClient();

const router = createHashRouter([
  {
    path: '/',
    element: <Root />,
    children: [{ index: true, element: <MatchComparison /> }],
  },
]);

const rootElement = document.getElementById('root')!;

if (!rootElement.innerHTML) {
  const root = createRoot(rootElement);

  root.render(
    <StrictMode>
      <TooltipProvider skipDelayDuration={0}>
        <QueryClientProvider client={queryClient}>
          <DialogProvider>
            <RouterProvider router={router} />
            <Toaster />
          </DialogProvider>
        </QueryClientProvider>
      </TooltipProvider>
    </StrictMode>
  );
}
