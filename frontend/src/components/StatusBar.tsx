import { File, FolderOpen } from 'lucide-react';
import { useEffect, useState } from 'react';

import { withErrorHandling } from '@/lib/errors';

import {
  GetResultFilePath,
  GetScanRoot,
  GetScanSettingsFilePath,
  SelectDirectory,
  SelectFile,
  SetResultFilePath,
  SetScanRoot,
  SetScanSettingsFilePath,
} from '../../wailsjs/go/main/App';
import { toast } from './ui/use-toast';

export default function StatusBar() {
  const [scanRoot, setScanRoot] = useState<string>('');
  const [resultsFile, setResultsFile] = useState<string>('');
  const [settingsFile, setSettingsFile] = useState<string>('');

  const handleSelectScanRoot = withErrorHandling({
    asyncFn: async () => {
      const selectedDir = await SelectDirectory(scanRoot ?? '.');
      if (selectedDir) {
        await SetScanRoot(selectedDir);
      }
    },
    onError: (error) => {
      console.error('Error selecting scan root:', error);
      toast({
        variant: 'destructive',
        title: 'Error',
        description: 'An error occurred while selecting the scan root. Please try again.',
      });
    },
  });

  const handleSelectResultsFile = withErrorHandling({
    asyncFn: async () => {
      const file = await SelectFile(scanRoot ?? '.');
      if (file) {
        await SetResultFilePath(file);
      }
    },
    onError: (error) => {
      console.error('Error selecting results file:', error);
      toast({
        variant: 'destructive',
        title: 'Error',
        description: 'An error occurred while selecting the results file. Please try again.',
      });
    },
  });

  const handleSelectSettingsFile = withErrorHandling({
    asyncFn: async () => {
      const file = await SelectFile(scanRoot ?? '.');
      if (file) {
        await SetScanSettingsFilePath(file);
      }
    },
    onError: (error) => {
      console.error('Error selecting settings file:', error);
      toast({
        variant: 'destructive',
        title: 'Error',
        description: 'An error occurred while selecting the settings file. Please try again.',
      });
    },
  });

  useEffect(() => {
    GetScanRoot().then(setScanRoot);

    GetResultFilePath().then(setResultsFile);

    GetScanSettingsFilePath().then(setSettingsFile);
  }, []);

  return (
    <div className="w-full bg-background px-4 text-xs text-muted-foreground">
      <div className="flex items-center gap-4">
        <div className="flex items-center gap-2">
          <span>Scan Root:</span>
          <div className="flex h-full cursor-pointer items-center p-1 hover:bg-accent hover:text-accent-foreground" onClick={handleSelectScanRoot}>
            <FolderOpen className="mr-1 h-3 w-3" />
            {scanRoot || 'Select Directory'}
          </div>
        </div>
        <div className="flex items-center gap-2">
          <span>Results File:</span>
          <div className="flex h-full cursor-pointer items-center p-1 hover:bg-accent hover:text-accent-foreground" onClick={handleSelectResultsFile}>
            <File className="mr-1 h-3 w-3" />
            {resultsFile || 'Select File'}
          </div>
        </div>
        <div className="flex items-center gap-2">
          <span>Scanoss Settings File:</span>
          <div
            className="flex h-full cursor-pointer items-center p-1 hover:bg-accent hover:text-accent-foreground"
            onClick={handleSelectSettingsFile}
          >
            <File className="mr-1 h-3 w-3" />
            {settingsFile || 'Select File'}
          </div>
        </div>
      </div>
    </div>
  );
}
