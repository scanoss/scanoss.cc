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

import { ExternalLink, Folder, Loader2 } from 'lucide-react';
import { useEffect, useState } from 'react';

import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { withErrorHandling } from '@/lib/errors';
import useResultsStore from '@/modules/results/stores/useResultsStore';
import useConfigStore from '@/stores/useConfigStore';

import { GetWorkingDir, SelectDirectory } from '../../wailsjs/go/main/App';
import { GetDefaultScanArgs, ScanStream } from '../../wailsjs/go/service/ScanServicePythonImpl';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import Link from './Link';
import TerminalOutput from './TerminalOutput';
import { Label } from './ui/label';
import { ScrollArea } from './ui/scroll-area';
import { useToast } from './ui/use-toast';

interface OutputLine {
  type: 'stdout' | 'stderr' | 'error';
  text: string;
}

type ScanStatus = 'idle' | 'scanning' | 'completed' | 'failed';

interface ScanDialogProps {
  onOpenChange: () => void;
  withOptions: boolean;
}

export default function ScanDialog({ onOpenChange, withOptions }: ScanDialogProps) {
  const { toast } = useToast();

  const setScanRoot = useConfigStore((state) => state.setScanRoot);

  const [directory, setDirectory] = useState('');
  const [args, setArgs] = useState<string>('');
  const [output, setOutput] = useState<OutputLine[]>([]);
  const [scanStatus, setScanStatus] = useState<ScanStatus>('idle');

  const fetchResults = useResultsStore((state) => state.fetchResults);
  const setSelectedResults = useResultsStore((state) => state.setSelectedResults);

  const handleSelectDirectory = withErrorHandling({
    asyncFn: async () => {
      const selectedDir = await SelectDirectory();
      if (selectedDir) {
        setDirectory(selectedDir);
      }
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'An error occurred while selecting the directory. Please try again.',
        variant: 'destructive',
      });
    },
  });

  const handleScan = withErrorHandling({
    asyncFn: async () => {
      setScanStatus('scanning');
      setOutput([]);
      console.log(args);
      await ScanStream([directory, ...args.split(' ')]);
      await setScanRoot(directory);
      setSelectedResults([]);
      await fetchResults();
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'An error occurred while scanning. Please try again.',
        variant: 'destructive',
      });
    },
    onFinish: () => {
      setScanStatus('idle');
    },
  });

  useEffect(() => {
    async function fetchDefaultDirectory() {
      const dir = await GetWorkingDir();
      setDirectory(dir);
    }
    async function fetchDefaultScanArgs() {
      const defaultArgs = await GetDefaultScanArgs();
      setArgs(defaultArgs.join(' '));
    }
    fetchDefaultDirectory();
    fetchDefaultScanArgs();
  }, []);

  useEffect(() => {
    const subs = [
      EventsOn('commandOutput', (data) => {
        setOutput((prev) => [...prev, { type: 'stdout', text: data }]);
      }),
      EventsOn('commandError', (data) => {
        setOutput((prev) => [...prev, { type: 'stderr', text: data }]);
      }),
      EventsOn('scanComplete', () => {
        setScanStatus('completed');
      }),
      EventsOn('scanFailed', (error) => {
        setScanStatus('failed');
        setOutput((prev) => [...prev, { type: 'error', text: error }]);
      }),
    ];

    return () => subs.forEach((unsub) => unsub());
  }, []);

  const isScanning = scanStatus === 'scanning';

  return (
    <Dialog open onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Scan With Options</DialogTitle>
          <DialogDescription>Run a scan on the selected directory with the provided arguments.</DialogDescription>
        </DialogHeader>
        <ScrollArea className="max-h-[600px]">
          <div className="flex flex-col gap-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="directory">Directory to scan:</Label>
              <div className="flex">
                <Input
                  id="directory"
                  value={directory}
                  readOnly
                  placeholder="Select a directory"
                  className="text-left font-mono text-sm [direction:rtl]"
                />
                {withOptions && (
                  <Button type="button" onClick={handleSelectDirectory} className="ml-2 px-3" variant="secondary">
                    <Folder className="h-4 w-4" />
                  </Button>
                )}
              </div>
            </div>
            {withOptions && (
              <div className="space-y-2">
                <Label htmlFor="scan-arguments">Arguments:</Label>
                <Input
                  id="scan-arguments"
                  value={args}
                  onChange={(e) => setArgs(e.target.value)}
                  className="font-mono text-sm"
                  placeholder="Enter scan arguments..."
                />
                <p className="text-sm text-gray-500">
                  Need help with arguments?{' '}
                  <Link
                    to="https://scanoss.readthedocs.io/projects/scanoss-py/en/latest/#commands-and-arguments"
                    className="inline-flex items-center"
                  >
                    View documentation
                    <ExternalLink className="ml-1 h-3 w-3" />
                  </Link>
                </p>
              </div>
            )}
          </div>
        </ScrollArea>
        <DialogFooter>
          <div className="flex w-full flex-col gap-3">
            <Button onClick={handleScan} disabled={isScanning || !directory} className="w-full">
              {isScanning ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Scanning...
                </>
              ) : (
                'Start Scan'
              )}
            </Button>
            <div className="space-y-2">
              <Label>Output:</Label>
              {output.length ? (
                <TerminalOutput lines={output} />
              ) : (
                <p className="text-sm italic text-gray-500">No output available. Start a scan to see results.</p>
              )}
            </div>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
