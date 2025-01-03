// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
