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

import { useQueryClient } from '@tanstack/react-query';
import { Folder } from 'lucide-react';

import useSelectedResult from '@/hooks/useSelectedResult';
import { withErrorHandling } from '@/lib/errors';
import { truncatePath } from '@/lib/utils';
import useResultsStore from '@/modules/results/stores/useResultsStore';
import useConfigStore from '@/stores/useConfigStore';

import { SelectDirectory } from '../../wailsjs/go/main/App';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';
import { toast } from './ui/use-toast';

export default function SelectScanRoot() {
  const queryClient = useQueryClient();

  const selectedResult = useSelectedResult();

  const scanRoot = useConfigStore((state) => state.scanRoot);
  const setScanRoot = useConfigStore((state) => state.setScanRoot);

  const fetchResults = useResultsStore((state) => state.fetchResults);

  const handleSelectScanRoot = withErrorHandling({
    asyncFn: async () => {
      const selectedDir = await SelectDirectory();
      if (selectedDir) {
        await setScanRoot(selectedDir);
        await fetchResults();
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

  return (
    <Tooltip>
      <TooltipTrigger>
        <span
          className="flex h-full cursor-pointer items-center gap-1 p-1 text-xs hover:bg-accent hover:text-accent-foreground"
          onClick={handleSelectScanRoot}
        >
          <Folder className="h-3 w-3 flex-shrink-0" />
          {truncatePath(scanRoot)}
        </span>
      </TooltipTrigger>
      <TooltipContent>{scanRoot}</TooltipContent>
    </Tooltip>
  );
}
