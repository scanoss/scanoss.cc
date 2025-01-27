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

import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Separator } from '@/components/ui/separator';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { withErrorHandling } from '@/lib/errors';
import useResultsStore from '@/modules/results/stores/useResultsStore';
import useConfigStore from '@/stores/useConfigStore';

import { GetWorkingDir, SelectDirectory } from '../../wailsjs/go/main/App';
import { ScanStream } from '../../wailsjs/go/service/ScanServicePythonImpl';
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

interface ScanOptions {
  output: string;
  settings: string;
  debug: boolean;
  trace: boolean;
  quiet: boolean;
  showAdvanced: boolean;
  // Advanced options
  threads: number;
  timeout: number;
  retry: number;
  postSize: number;
  flags: number;
  dependencies: boolean;
  dependenciesOnly: boolean;
  noWfpOutput: boolean;
}

interface ScanDialogProps {
  onOpenChange: () => void;
}

export default function ScanDialog({ onOpenChange }: ScanDialogProps) {
  const { toast } = useToast();

  const setScanRoot = useConfigStore((state) => state.setScanRoot);

  const [output, setOutput] = useState<OutputLine[]>([]);
  const [scanStatus, setScanStatus] = useState<ScanStatus>('idle');

  const [directory, setDirectory] = useState('');
  const [options, setOptions] = useState<ScanOptions>({
    output: '',
    settings: '',
    debug: false,
    trace: false,
    quiet: true,
    showAdvanced: false,
    threads: 5,
    timeout: 180,
    retry: 5,
    postSize: 32,
    flags: 0,
    dependencies: false,
    dependenciesOnly: false,
    noWfpOutput: false,
  });

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

      const scanArgs: string[] = [directory];

      // Add core options
      if (options.output) scanArgs.push('--output', options.output);
      if (options.settings) scanArgs.push('--settings', options.settings);
      if (options.debug) scanArgs.push('--debug');
      if (options.trace) scanArgs.push('--trace');
      if (options.quiet) scanArgs.push('--quiet');

      // Add advanced options
      if (options.threads !== 5) scanArgs.push('--threads', options.threads.toString());
      if (options.timeout !== 180) scanArgs.push('--timeout', options.timeout.toString());
      if (options.retry !== 5) scanArgs.push('--retry', options.retry.toString());
      if (options.postSize !== 32) scanArgs.push('--post-size', options.postSize.toString());
      if (options.flags !== 0) scanArgs.push('--flags', options.flags.toString());
      if (options.dependencies) scanArgs.push('--dependencies');
      if (options.dependenciesOnly) scanArgs.push('--dependencies-only');
      if (options.noWfpOutput) scanArgs.push('--no-wfp-output');

      await ScanStream(scanArgs);
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
      // Set default output path based on directory
      setOptions((prev) => ({
        ...prev,
        output: `${dir}/.scanoss/results.json`,
        settings: `${dir}/scanoss.json`,
      }));
    }
    fetchDefaultDirectory();
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

  const handleSelectOutput = withErrorHandling({
    asyncFn: async () => {
      const selectedFile = await SelectDirectory();
      if (selectedFile) {
        setOptions((prev) => ({ ...prev, output: selectedFile }));
      }
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'An error occurred while selecting the output file. Please try again.',
        variant: 'destructive',
      });
    },
  });

  const handleSelectSettings = withErrorHandling({
    asyncFn: async () => {
      const selectedFile = await SelectDirectory();
      if (selectedFile) {
        setOptions((prev) => ({ ...prev, settings: selectedFile }));
      }
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'An error occurred while selecting the settings file. Please try again.',
        variant: 'destructive',
      });
    },
  });

  const isScanning = scanStatus === 'scanning';

  return (
    <Dialog open onOpenChange={onOpenChange}>
      <DialogContent className="flex max-h-[90vh] max-w-3xl flex-col">
        <DialogHeader>
          <DialogTitle>Scan With Options</DialogTitle>
          <DialogDescription>Run a scan on the selected directory with the provided arguments.</DialogDescription>
        </DialogHeader>
        <ScrollArea className="flex-1">
          <div className="flex flex-col gap-6 py-4 pr-4">
            {/* Core Settings */}
            <div className="space-y-6">
              <div className="space-y-2">
                <Label htmlFor="directory">Directory to scan</Label>
                <div className="flex gap-2">
                  <Input id="directory" value={directory} readOnly placeholder="Select a directory" className="text-sm" />
                  <Button type="button" onClick={handleSelectDirectory} variant="secondary" size="icon">
                    <Folder className="h-4 w-4" />
                  </Button>
                </div>
              </div>

              <div className="grid grid-cols-1 gap-4">
                <div className="space-y-2">
                  <div className="flex items-center space-x-2">
                    <Label htmlFor="output">Output Path</Label>
                  </div>
                  <div className="flex gap-2">
                    <Input
                      id="output"
                      value={options.output}
                      onChange={(e) => setOptions((prev) => ({ ...prev, output: e.target.value }))}
                      placeholder="[root]/.scanoss/results.json"
                      className="text-sm"
                    />
                    <Button type="button" onClick={handleSelectOutput} variant="secondary" size="icon">
                      <Folder className="h-4 w-4" />
                    </Button>
                  </div>
                </div>

                <div className="space-y-2">
                  <div className="flex items-center space-x-2">
                    <Label htmlFor="settings">Settings File</Label>
                  </div>
                  <div className="flex gap-2">
                    <Input
                      id="settings"
                      value={options.settings}
                      onChange={(e) => setOptions((prev) => ({ ...prev, settings: e.target.value }))}
                      placeholder="[root]/scanoss.json"
                      className="text-sm"
                    />
                    <Button type="button" onClick={handleSelectSettings} variant="secondary" size="icon">
                      <Folder className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </div>

              {/* Core Flags */}
              <div className="space-y-2">
                <Label>Basic Options</Label>
                <div className="grid grid-cols-3 gap-4">
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <div className="flex items-center space-x-2">
                        <Checkbox
                          id="debug"
                          checked={options.debug}
                          onCheckedChange={(checked) => setOptions((prev) => ({ ...prev, debug: checked === true }))}
                        />
                        <Label htmlFor="debug" className="cursor-help">
                          Debug
                        </Label>
                      </div>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>Enable debug messages for detailed logging</p>
                    </TooltipContent>
                  </Tooltip>

                  <Tooltip>
                    <TooltipTrigger asChild>
                      <div className="flex items-center space-x-2">
                        <Checkbox
                          id="trace"
                          checked={options.trace}
                          onCheckedChange={(checked) => setOptions((prev) => ({ ...prev, trace: checked === true }))}
                        />
                        <Label htmlFor="trace" className="cursor-help">
                          Trace
                        </Label>
                      </div>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>Enable trace messages, including API posts</p>
                    </TooltipContent>
                  </Tooltip>

                  <Tooltip>
                    <TooltipTrigger asChild>
                      <div className="flex items-center space-x-2">
                        <Checkbox
                          id="quiet"
                          checked={options.quiet}
                          onCheckedChange={(checked) => setOptions((prev) => ({ ...prev, quiet: checked === true }))}
                        />
                        <Label htmlFor="quiet" className="cursor-help">
                          Quiet
                        </Label>
                      </div>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>Enable quiet mode (reduce output verbosity)</p>
                    </TooltipContent>
                  </Tooltip>
                </div>
              </div>
            </div>

            {/* Advanced Settings */}
            <Accordion type="single" collapsible className="w-full">
              <AccordionItem value="advanced-settings">
                <AccordionTrigger className="text-sm font-medium">Advanced Settings</AccordionTrigger>
                <AccordionContent>
                  <div className="space-y-6 pt-2">
                    {/* Performance Settings */}
                    <div className="space-y-4">
                      <Label className="text-sm font-medium">Performance</Label>
                      <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                          <Label htmlFor="threads" className="text-sm text-muted-foreground">
                            Threads
                          </Label>
                          <Input
                            id="threads"
                            type="number"
                            min="1"
                            max="32"
                            value={options.threads}
                            onChange={(e) => {
                              const value = parseInt(e.target.value);
                              if (!isNaN(value) && value >= 1 && value <= 32) {
                                setOptions((prev) => ({ ...prev, threads: value }));
                              }
                            }}
                            className="font-mono text-sm"
                          />
                        </div>
                        <div className="space-y-2">
                          <Label htmlFor="timeout" className="text-sm text-muted-foreground">
                            Timeout (s)
                          </Label>
                          <Input
                            id="timeout"
                            type="number"
                            min="1"
                            max="3600"
                            value={options.timeout}
                            onChange={(e) => {
                              const value = parseInt(e.target.value);
                              if (!isNaN(value) && value >= 1 && value <= 3600) {
                                setOptions((prev) => ({ ...prev, timeout: value }));
                              }
                            }}
                            className="font-mono text-sm"
                          />
                        </div>
                      </div>
                    </div>

                    <Separator />

                    {/* Additional Options */}
                    <div className="space-y-4">
                      <Label className="text-sm font-medium">Additional Options</Label>
                      <div className="grid grid-cols-2 gap-4">
                        <div className="flex items-center space-x-2">
                          <Checkbox
                            id="dependencies"
                            checked={options.dependencies}
                            onCheckedChange={(checked) => setOptions((prev) => ({ ...prev, dependencies: checked === true }))}
                          />
                          <Label htmlFor="dependencies" className="text-sm text-muted-foreground">
                            Dependencies
                          </Label>
                        </div>
                        <div className="flex items-center space-x-2">
                          <Checkbox
                            id="dependenciesOnly"
                            checked={options.dependenciesOnly}
                            onCheckedChange={(checked) => setOptions((prev) => ({ ...prev, dependenciesOnly: checked === true }))}
                          />
                          <Label htmlFor="dependenciesOnly" className="text-sm text-muted-foreground">
                            Dependencies Only
                          </Label>
                        </div>
                        <div className="flex items-center space-x-2">
                          <Checkbox
                            id="noWfpOutput"
                            checked={options.noWfpOutput}
                            onCheckedChange={(checked) => setOptions((prev) => ({ ...prev, noWfpOutput: checked === true }))}
                          />
                          <Label htmlFor="noWfpOutput" className="text-sm text-muted-foreground">
                            Skip WFP Output
                          </Label>
                        </div>
                      </div>
                    </div>
                  </div>
                </AccordionContent>
              </AccordionItem>
            </Accordion>

            {/* Documentation Link */}
            <p className="text-sm text-muted-foreground">
              Need help with options?{' '}
              <Link
                to="https://scanoss.readthedocs.io/projects/scanoss-py/en/latest/#commands-and-arguments"
                className="inline-flex items-center hover:text-primary"
              >
                View documentation
                <ExternalLink className="ml-1 h-3 w-3" />
              </Link>
            </p>
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

            {(isScanning || output.length > 0) && (
              <div className="space-y-2">
                <Label>Console Output:</Label>
                <TerminalOutput lines={output} />
              </div>
            )}
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
