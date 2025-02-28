// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2025 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

import { FolderOpen, Search } from 'lucide-react';
import { useState } from 'react';

import { withErrorHandling } from '@/lib/errors';
import useConfigStore from '@/stores/useConfigStore';

import { SelectDirectory } from '../../wailsjs/go/main/App';
import { entities } from '../../wailsjs/go/models';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import ScanDialog from './ScanDialog';
import { Button } from './ui/button';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { toast } from './ui/use-toast';

export default function WelcomeScreen() {
  const setScanRoot = useConfigStore((state) => state.setScanRoot);
  const recentScanRoots = useConfigStore((state) => state.recentScanRoots);
  const [scanModal, setScanModal] = useState(false);

  const handleSelectScanRoot = withErrorHandling({
    asyncFn: async () => {
      const selectedDir = await SelectDirectory();
      if (selectedDir) {
        await setScanRoot(selectedDir);
      }
    },
    onError: () => {
      toast({
        variant: 'destructive',
        title: 'Error',
        description: 'An error occurred while selecting the scan root. Please try again.',
      });
    },
  });

  const handleSelectRecentFolder = withErrorHandling({
    asyncFn: async (path: string) => {
      await setScanRoot(path);
    },
    onError: () => {
      toast({
        variant: 'destructive',
        title: 'Error',
        description: 'An error occurred while selecting the recent folder. Please try again.',
      });
    },
  });

  const handleCloseScanModal = () => setScanModal(false);
  const handleShowScanModal = () => setScanModal(true);

  EventsOn(entities.Action.ScanWithOptions, () => {
    handleShowScanModal();
  });

  return (
    <div className="flex h-screen w-screen flex-col items-center justify-center space-y-8 p-8">
      <div className="space-y-2 text-center">
        <h1 className="text-3xl font-semibold tracking-tight">SCANOSS Code Compare</h1>
        <p className="text-lg text-muted-foreground">Select a folder to start</p>
      </div>

      <div className="grid w-full max-w-4xl gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg font-medium">Start</CardTitle>
          </CardHeader>
          <CardContent className="grid gap-4">
            <Button variant="ghost" className="w-full justify-start gap-2 px-2" onClick={handleSelectScanRoot}>
              <FolderOpen className="h-5 w-5" />
              Open Folder...
            </Button>
            <Button variant="ghost" className="w-full justify-start gap-2 px-2" onClick={handleShowScanModal}>
              <Search className="h-5 w-5" />
              Run a new scan...
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-lg font-medium">Recent</CardTitle>
          </CardHeader>
          <CardContent>
            {recentScanRoots.length > 0 ? (
              <div className="grid gap-2">
                {recentScanRoots.map((path) => (
                  <Button
                    key={path}
                    variant="ghost"
                    className="w-full justify-start gap-2 px-2 text-sm"
                    onClick={() => handleSelectRecentFolder(path)}
                  >
                    <FolderOpen className="h-4 w-4" />
                    {path}
                  </Button>
                ))}
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">No recent folders</p>
            )}
          </CardContent>
        </Card>
      </div>

      <ScanDialog open={scanModal} onOpenChange={handleCloseScanModal} />
    </div>
  );
}
