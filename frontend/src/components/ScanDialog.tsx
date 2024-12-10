'use client';

import { ExternalLink, Folder, Loader2 } from 'lucide-react';
import { useState } from 'react';

import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';

import { SelectDirectory } from '../../wailsjs/go/main/App';
import { Scan } from '../../wailsjs/go/service/ScanServicePythonImpl';
import Link from './Link';
import { Label } from './ui/label';
import { useToast } from './ui/use-toast';

const DEFAULT_SCAN_ARGS = '--quiet -o .scanoss/results.json --no-wfp-output';

export default function ScanModal({ open, onOpenChange }: { open: boolean; onOpenChange: () => void }) {
  const { toast } = useToast();
  const [directory, setDirectory] = useState('');
  const [args, setArgs] = useState(DEFAULT_SCAN_ARGS);
  const [output, setOutput] = useState('');
  const [isScanning, setIsScanning] = useState(false);

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
    setIsScanning(true);
    try {
      const result = await Scan([directory, args]);
      console.log('Scan result:', result);
    } catch (error) {
      console.error('Error scanning:', error);
      toast({
        title: 'Error',
        description: 'An error occurred while scanning. Please try again.',
        variant: 'destructive',
      });
    } finally {
      setIsScanning(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Scan With Options</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="directory">Directory to scan:</Label>
            <div className="flex">
              <Input id="directory" value={directory} readOnly placeholder="Select a directory" className="flex-grow font-mono text-sm" />
              <Button type="button" onClick={handleSelectDirectory} className="ml-2 px-3" variant="secondary">
                <Folder className="h-4 w-4" />
              </Button>
            </div>
          </div>
          <div className="grid gap-2">
            <label htmlFor="args" className="text-sm font-medium">
              Arguments:
            </label>
            <Input id="args" value={args} onChange={(e) => setArgs(e.target.value)} className="font-mono text-sm" />
            <p className="text-sm text-gray-500">
              Need help with arguments?{' '}
              <Link to="https://scanoss.readthedocs.io/projects/scanoss-py/en/latest/#commands-and-arguments" className="inline-flex items-center">
                View documentation
                <ExternalLink className="ml-1 h-3 w-3" />
              </Link>
            </p>
          </div>
          <Button onClick={handleScan} disabled={isScanning || !directory} className="w-full bg-black text-white hover:bg-black/90">
            {isScanning ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Scanning...
              </>
            ) : (
              'Start Scan'
            )}
          </Button>
          <Accordion type="single" collapsible className="w-full" defaultValue="output">
            <AccordionItem value="output">
              <AccordionTrigger className="text-sm font-medium" disabled={!output}>
                Output
              </AccordionTrigger>
              <AccordionContent>
                {output ? (
                  <pre className="overflow-x-auto whitespace-pre-wrap rounded-md bg-zinc-950 p-4 font-mono text-sm text-white">{output}</pre>
                ) : (
                  <p className="text-sm italic text-gray-500">No output available. Start a scan to see results.</p>
                )}
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        </div>
      </DialogContent>
    </Dialog>
  );
}
