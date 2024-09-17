import './style.css';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { createHashRouter, RouterProvider } from 'react-router-dom';

import ConfirmDialog from './components/ConfirmDialog';
import { Toaster } from './components/ui/toaster';
import { TooltipProvider } from './components/ui/tooltip';
import { ResultsProvider } from './modules/results/providers/ResultsProvider';
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
      <TooltipProvider skipDelayDuration={0}>
        <QueryClientProvider client={queryClient}>
          <ConfirmDialogProvider>
            <ResultsProvider>
              <RouterProvider router={router} />
              <ConfirmDialog />
              <Toaster />
            </ResultsProvider>
          </ConfirmDialogProvider>
        </QueryClientProvider>
      </TooltipProvider>
    </StrictMode>
  );
}
