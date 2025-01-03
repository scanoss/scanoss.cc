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

import { useQueryClient } from '@tanstack/react-query';
import { File, FolderOpen } from 'lucide-react';
import { useEffect } from 'react';

import useSelectedResult from '@/hooks/useSelectedResult';
import { withErrorHandling } from '@/lib/errors';
import useResultsStore from '@/modules/results/stores/useResultsStore';
import useConfigStore from '@/stores/useConfigStore';

import { SelectDirectory, SelectFile } from '../../wailsjs/go/main/App';
import { toast } from './ui/use-toast';

export default function StatusBar() {
  const queryClient = useQueryClient();
  const selectedResult = useSelectedResult();

  const fetchResults = useResultsStore((state) => state.fetchResults);

  const scanRoot = useConfigStore((state) => state.scanRoot);
  const resultsFile = useConfigStore((state) => state.resultsFile);
  const settingsFile = useConfigStore((state) => state.settingsFile);

  const setScanRoot = useConfigStore((state) => state.setScanRoot);
  const setResultsFile = useConfigStore((state) => state.setResultsFile);
  const setSettingsFile = useConfigStore((state) => state.setSettingsFile);

  const getInitialConfig = useConfigStore((state) => state.getInitialConfig);

  const handleSelectScanRoot = withErrorHandling({
    asyncFn: async () => {
      const selectedDir = await SelectDirectory(scanRoot ?? '.');
      if (selectedDir) {
        await setScanRoot(selectedDir);
        await queryClient.invalidateQueries({
          queryKey: ['localFileContent', selectedResult?.path],
        });
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

  const handleSelectResultsFile = withErrorHandling({
    asyncFn: async () => {
      const file = await SelectFile(scanRoot ?? '.');
      if (file) {
        await setResultsFile(file);
        await fetchResults();
      }
    },
    onError: () => {
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
        await setSettingsFile(file);
        await fetchResults();
      }
    },
    onError: () => {
      toast({
        variant: 'destructive',
        title: 'Error',
        description: 'An error occurred while selecting the settings file. Please try again.',
      });
    },
  });

  useEffect(() => {
    getInitialConfig();
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
