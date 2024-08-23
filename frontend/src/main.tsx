import './style.css';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React, { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { createHashRouter, RouterProvider } from 'react-router-dom';

import ConfirmDialog from './components/ConfirmDialog';
import { TooltipProvider } from './components/ui/tooltip';
import { ConfirmDialogProvider } from './providers/ConfirmDialogProvider';
import Index from './routes';
import FileComparison from './routes/files/match';
import Root from './routes/root';

const queryClient = new QueryClient();

const router = createHashRouter([
  {
    path: '/',
    element: <Root />,
    children: [
      { index: true, element: <Index /> },
      {
        path: 'files/:filePath',
        element: <FileComparison />,
      },
    ],
  },
]);

const rootElement = document.getElementById('root')!;

if (!rootElement.innerHTML) {
  const root = createRoot(rootElement);

  root.render(
    <StrictMode>
      <TooltipProvider>
        <QueryClientProvider client={queryClient}>
          <ConfirmDialogProvider>
            <ConfirmDialog />
            <RouterProvider router={router} />
          </ConfirmDialogProvider>
        </QueryClientProvider>
      </TooltipProvider>
    </StrictMode>
  );
}
