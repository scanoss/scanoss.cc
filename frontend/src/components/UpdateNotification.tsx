// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
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

import { Download, RefreshCw } from 'lucide-react';
import { useEffect, useState } from 'react';

import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';

import { entities } from '../../wailsjs/go/models';
import { ApplyUpdate, CheckForUpdate, DownloadUpdate } from '../../wailsjs/go/service/UpdateServiceImpl';
import { useToast } from './ui/use-toast';

export default function UpdateNotification() {
  const { toast } = useToast();

  const [updateInfo, setUpdateInfo] = useState<entities.UpdateInfo | null>(null);
  const [, setChecking] = useState(false);
  const [downloading, setDownloading] = useState(false);
  const [showDialog, setShowDialog] = useState(false);
  const [downloadedPath, setDownloadedPath] = useState<string | null>(null);

  useEffect(() => {
    checkForUpdates();
  }, []);

  const checkForUpdates = async () => {
    try {
      setChecking(true);
      const update = await CheckForUpdate();

      if (update && update.available) {
        setUpdateInfo(update);

        // Show a toast notification about the update
        toast({
          title: 'Update Available',
          description: `Version ${update.version} is available for download.`,
        });
      }
    } catch (error) {
      console.error('Failed to check for updates:', error);
    } finally {
      setChecking(false);
    }
  };

  const handleDownloadUpdate = async () => {
    if (!updateInfo) return;

    try {
      setDownloading(true);

      toast({
        title: 'Downloading Update',
        description: 'Please wait while we download the update...',
      });

      const path = await DownloadUpdate(updateInfo);
      setDownloadedPath(path);

      toast({
        title: 'Download Complete',
        description: 'The update is ready to install.',
      });
    } catch (error) {
      console.error('Failed to download update:', error);
      toast({
        variant: 'destructive',
        title: 'Download Failed',
        description: 'Failed to download the update. Please try again.',
      });
    } finally {
      setDownloading(false);
    }
  };

  const handleInstallUpdate = async () => {
    if (!downloadedPath) return;

    try {
      await ApplyUpdate(downloadedPath);
      // Application will restart automatically
    } catch (error) {
      console.error('Failed to install update:', error);
      toast({
        variant: 'destructive',
        title: 'Installation Failed',
        description: 'Failed to install the update. Please try installing manually.',
      });
    }
  };

  if (!updateInfo?.available) {
    return null;
  }

  return (
    <>
      <Button variant="secondary" size="sm" className="h-6 gap-2 px-2 py-1 text-xs" onClick={() => setShowDialog(true)}>
        <Download className="h-3 w-3" />
        Update Available: v{updateInfo.version}
      </Button>

      <Dialog open={showDialog} onOpenChange={setShowDialog}>
        <DialogContent className="max-w-2xl p-4">
          <DialogHeader>
            <DialogTitle>Update Available</DialogTitle>
            <DialogDescription>A new version of SCANOSS Code Compare is available.</DialogDescription>
          </DialogHeader>

          <div className="space-y-4">
            <div>
              <h4 className="font-semibold">Version {updateInfo.version}</h4>
              <p className="text-sm text-muted-foreground">Published: {new Date(updateInfo.published_at).toLocaleDateString()}</p>
            </div>

            {updateInfo.release_notes && (
              <div>
                <h4 className="mb-2 font-semibold">Release Notes</h4>
                <div className="max-h-64 overflow-y-auto rounded-md border bg-muted p-4">
                  <pre className="whitespace-pre-wrap text-sm">{updateInfo.release_notes}</pre>
                </div>
              </div>
            )}
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowDialog(false)}>
              Later
            </Button>
            {!downloadedPath ? (
              <Button onClick={handleDownloadUpdate} disabled={downloading}>
                {downloading ? (
                  <>
                    <RefreshCw className="mr-2 h-4 w-4 animate-spin" />
                    Downloading...
                  </>
                ) : (
                  <>
                    <Download className="mr-2 h-4 w-4" />
                    Download Update
                  </>
                )}
              </Button>
            ) : (
              <Button onClick={handleInstallUpdate}>
                <RefreshCw className="mr-2 h-4 w-4" />
                Restart & Update
              </Button>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
