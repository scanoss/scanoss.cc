'use client';

import { ExternalLink, Folder, Loader2 } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';

import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';

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

export default function ScanModal({ open, onOpenChange }: { open: boolean; onOpenChange: () => void }) {
  const { toast } = useToast();

  const outputRef = useRef<HTMLDivElement>(null);

  const [directory, setDirectory] = useState('');
  const [args, setArgs] = useState<string>('');
  const [output, setOutput] = useState<OutputLine[]>([]);
  const [scanStatus, setScanStatus] = useState<ScanStatus>('idle');

  const handleSelectDirectory = async () => {
    try {
      const selectedDir = await SelectDirectory();
      setDirectory(selectedDir);
    } catch (error) {
      console.error('Error selecting directory:', error);
      toast({
        title: 'Error',
        description: 'An error occurred while selecting the directory. Please try again.',
        variant: 'destructive',
      });
    }
  };

  const handleScan = async () => {
    setScanStatus('scanning');
    setOutput([]);
    try {
      await ScanStream([directory, ...args.split(' ')]);
    } catch (error) {
      console.error('Error scanning:', error);
      toast({
        title: 'Error',
        description: 'An error occurred while scanning. Please try again.',
        variant: 'destructive',
      });
    } finally {
      setScanStatus('idle');
    }
  };

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
    if (output && outputRef.current) {
      outputRef.current.scrollTo({ behavior: 'smooth', top: outputRef.current.scrollHeight });
    }
  }, [output]);

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
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Scan With Options</DialogTitle>
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
                <Button type="button" onClick={handleSelectDirectory} className="ml-2 px-3" variant="secondary">
                  <Folder className="h-4 w-4" />
                </Button>
              </div>
            </div>
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
                <Link to="https://scanoss.readthedocs.io/projects/scanoss-py/en/latest/#commands-and-arguments" className="inline-flex items-center">
                  View documentation
                  <ExternalLink className="ml-1 h-3 w-3" />
                </Link>
              </p>
            </div>
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
