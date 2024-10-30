import './style.css';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { createHashRouter, RouterProvider } from 'react-router-dom';

import ConfirmDialog from './components/ConfirmDialog';
import InputPromptDialog from './components/InputPromptDialog';
import MatchComparison from './components/MatchComparison';
import { Toaster } from './components/ui/toaster';
import { TooltipProvider } from './components/ui/tooltip';
import { ConfirmDialogProvider } from './providers/ConfirmDialogProvider';
import { InputPromptDialogProvider } from './providers/InputPromptDialogProvider';
import Index from './routes';
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
        element: <MatchComparison />,
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
            <InputPromptDialogProvider>
              <RouterProvider router={router} />
              <ConfirmDialog />
              <InputPromptDialog />
              <Toaster />
            </InputPromptDialogProvider>
          </ConfirmDialogProvider>
        </QueryClientProvider>
      </TooltipProvider>
    </StrictMode>
  );
}
