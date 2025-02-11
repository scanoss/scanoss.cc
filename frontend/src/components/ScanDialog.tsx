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

import { AnimatePresence, motion } from 'framer-motion';
import { Check, ExternalLink, Folder, Loader2 } from 'lucide-react';
import { useEffect, useState } from 'react';

import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { withErrorHandling } from '@/lib/errors';
import useResultsStore from '@/modules/results/stores/useResultsStore';
import useConfigStore from '@/stores/useConfigStore';

import { GetScanRoot, SelectDirectory } from '../../wailsjs/go/main/App';
import { entities } from '../../wailsjs/go/models';
import { GetScanArgs, ScanStream } from '../../wailsjs/go/service/ScanServicePythonImpl';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import Link from './Link';
import ScanOption from './ScanOption';
import TerminalOutput from './TerminalOutput';
import { useToast } from './ui/use-toast';

interface OutputLine {
  type: 'stdout' | 'stderr' | 'error';
  text: string;
}

type ScanStatus = 'idle' | 'scanning' | 'completed' | 'failed';

interface ScanDialogProps {
  onOpenChange: () => void;
}

export default function ScanDialog({ onOpenChange }: ScanDialogProps) {
  const { toast } = useToast();

  const setScanRoot = useConfigStore((state) => state.setScanRoot);

  const [output, setOutput] = useState<OutputLine[]>([]);
  const [scanStatus, setScanStatus] = useState<ScanStatus>('idle');
  const [showSuccess, setShowSuccess] = useState(false);

  const [directory, setDirectory] = useState('');
  const [scanArgs, setScanArgs] = useState<entities.ScanArgDef[]>([]);
  const [advancedScanArgs, setAdvancedScanArgs] = useState<string[]>([]);
  const [options, setOptions] = useState<Record<string, string | number | boolean | string[]>>({});

  const fetchResults = useResultsStore((state) => state.fetchResults);
  const setSelectedResults = useResultsStore((state) => state.setSelectedResults);

  const handleSelectDirectory = withErrorHandling({
    asyncFn: async () => {
      const selectedDir = await SelectDirectory();
      if (selectedDir) {
        setDirectory(selectedDir);
        setOptions((prev) => ({ ...prev, output: `${selectedDir}/.scanoss/results.json`, settings: `${selectedDir}/scanoss.json` }));
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

      let cmdArgs: string[] = [directory];

      // We should only add options that have non-default values
      Object.entries(options).forEach(([key, value]) => {
        const arg = scanArgs.find((a) => a.Name === key);

        if (arg && value !== arg.Default) {
          if (arg.Type === 'bool') {
            if (value === true) {
              cmdArgs.push(`--${key}`);
            } else {
              cmdArgs = cmdArgs.filter((arg) => arg !== `--${key}`);
            }
          } else {
            cmdArgs.push(`--${key}`, String(value));
          }
        }
      });

      cmdArgs.push(...advancedScanArgs);

      await ScanStream(cmdArgs);
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

  const handleOptionChange = (name: string, value: string | number | boolean | string[]) => {
    setOptions((prev) => ({ ...prev, [name]: value }));
  };

  const handleFileSelect = (name: string) => {
    return withErrorHandling({
      asyncFn: async () => {
        const selectedFile = await SelectDirectory();
        if (selectedFile) {
          const arg = scanArgs.find((a) => a.Name === name);
          if (arg?.Type === 'string') {
            handleOptionChange(name, selectedFile);
          }
        }
      },
      onError: () => {
        toast({
          title: 'Error',
          description: `An error occurred while selecting the ${name} file. Please try again.`,
          variant: 'destructive',
        });
      },
    });
  };

  useEffect(() => {
    async function initialize() {
      try {
        // Get scan arguments from backend
        const args = await GetScanArgs();
        setScanArgs(args);

        // Initialize options with default values
        const initialOptions: Record<string, string | number | boolean | string[]> = {};
        args.forEach((arg) => {
          initialOptions[arg.Name] = arg.Default;
        });
        setOptions(initialOptions);

        // Set default directory and paths
        const dir = await GetScanRoot();
        setDirectory(dir);
        setOptions((prev) => ({
          ...prev,
          output: `${dir}/.scanoss/results.json`,
          settings: `${dir}/scanoss.json`,
        }));
      } catch (error) {
        console.error('Failed to initialize scan dialog:', error);
        toast({
          title: 'Error',
          description: 'Failed to initialize scan options. Please try again.',
          variant: 'destructive',
        });
      }
    }
    initialize();
  }, [toast]);

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
        setShowSuccess(true);

        setTimeout(() => {
          setShowSuccess(false);
          setScanStatus('idle');
        }, 2000);
      }),
      EventsOn('scanFailed', (error) => {
        setScanStatus('failed');
        setOutput((prev) => [...prev, { type: 'error', text: error }]);
      }),
    ];

    return () => subs.forEach((unsub) => unsub());
  }, []);

  const isScanning = scanStatus === 'scanning';
  const coreOptions = scanArgs.filter((arg) => arg.IsCore);

  return (
    <Dialog open onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Scan With Options</DialogTitle>
          <DialogDescription>Run a scan on the selected directory with the provided arguments.</DialogDescription>
        </DialogHeader>

        <div className="flex flex-col gap-6 py-4">
          <div className="space-y-2">
            <Label htmlFor="directory">Directory to scan</Label>
            <div className="flex gap-2">
              <Input id="directory" value={directory} readOnly placeholder="Select a directory" className="text-sm" />
              <Button type="button" onClick={handleSelectDirectory} variant="secondary" size="icon">
                <Folder className="h-4 w-4" />
              </Button>
            </div>
          </div>

          {/* Core Options */}
          <div className="grid grid-cols-4 gap-6">
            {coreOptions.map((arg) => (
              <div key={arg.Name} className={arg.Type.toLowerCase() === 'bool' ? 'col-span-1' : 'col-span-4'}>
                <ScanOption
                  name={arg.Name}
                  type={arg.Type.toLowerCase() as 'string' | 'int' | 'bool' | 'stringSlice'}
                  value={options[arg.Name]}
                  defaultValue={arg.Default}
                  tooltip={arg.Tooltip}
                  onChange={(value) => handleOptionChange(arg.Name, value)}
                  onSelectFile={arg.IsFileSelector ? () => handleFileSelect(arg.Name)() : undefined}
                  isFileSelector={arg.IsFileSelector}
                />
              </div>
            ))}
          </div>

          {/* Advanced Options */}
          <div className="space-y-2">
            <Label htmlFor="advanced-options">Advanced Options</Label>
            <Input
              id="advanced-options"
              value={advancedScanArgs.join(' ')}
              onChange={(e) => setAdvancedScanArgs(e.target.value.split(' '))}
              placeholder="e.g. --timeout 500 --skip-settings-file --dependencies"
              className="text-sm"
            />
            <p className="text-xs text-muted-foreground">
              You can add advanced options here. These will be passed to the scan command as is.
              <Link to="https://scanoss.readthedocs.io/projects/scanoss-py/en/latest/#commands-and-arguments">
                <div className="inline-flex items-center">
                  View documentation
                  <ExternalLink className="ml-1 h-3 w-3" />
                </div>
              </Link>
            </p>
          </div>
        </div>
        <DialogFooter className="flex flex-1 gap-2 sm:flex-col">
          {output.length > 0 && (
            <div className="space-y-2">
              <Label>Console Output:</Label>
              <TerminalOutput lines={output} />
            </div>
          )}

          <div className="flex items-center justify-end gap-2">
            <Button variant="outline" onClick={onOpenChange}>
              Close
            </Button>

            <motion.div initial={false} animate={scanStatus} className="relative">
              <Button onClick={handleScan} disabled={isScanning || !directory} className="relative min-w-[120px] overflow-hidden">
                <AnimatePresence mode="wait">
                  {isScanning && (
                    <motion.span key="scanning" initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }} className="flex items-center">
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Scanning...
                    </motion.span>
                  )}
                  {showSuccess && (
                    <div className="flex items-center justify-center">
                      <Check className="h-4 w-4" />
                      <span className="ml-2 text-sm">Done!</span>
                    </div>
                  )}

                  {!isScanning && !showSuccess && (
                    <motion.span key="default" initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="flex items-center">
                      {scanStatus === 'completed' ? 'Scan Again' : 'Start Scan'}
                    </motion.span>
                  )}
                </AnimatePresence>
              </Button>
              {isScanning && (
                <>
                  <motion.div
                    className="pointer-events-none absolute inset-0 rounded-md border border-blue-400/30"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                  />
                  <motion.div
                    className="pointer-events-none absolute inset-0 rounded-md border-2 border-blue-500"
                    initial={{ opacity: 0 }}
                    animate={{
                      opacity: [0.15, 0.3, 0.15],
                      scale: [1, 1.05, 1],
                      transition: {
                        repeat: Infinity,
                        duration: 2,
                        ease: 'easeInOut',
                      },
                    }}
                  />
                </>
              )}
            </motion.div>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
