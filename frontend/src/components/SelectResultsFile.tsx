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

import { File } from 'lucide-react';

import { useResults } from '@/hooks/useResults';
import { withErrorHandling } from '@/lib/errors';
import { truncatePath } from '@/lib/utils';
import useConfigStore from '@/stores/useConfigStore';

import { SelectFile } from '../../wailsjs/go/main/App';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';
import { toast } from './ui/use-toast';

export default function SelectResultsFile() {
  const scanRoot = useConfigStore((state) => state.scanRoot);
  const resultsFile = useConfigStore((state) => state.resultsFile);
  const setResultsFile = useConfigStore((state) => state.setResultsFile);

  const { reset: resetResults } = useResults();

  const handleSelectResultsFile = withErrorHandling({
    asyncFn: async () => {
      const file = await SelectFile(scanRoot ?? '.');
      if (file) {
        await setResultsFile(file);
        resetResults();
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

  return (
    <Tooltip>
      <TooltipTrigger>
        <span
          className="flex h-full cursor-pointer items-center gap-1 p-1 text-xs hover:bg-accent hover:text-accent-foreground"
          onClick={handleSelectResultsFile}
        >
          <File className="h-3 w-3 flex-shrink-0" />
          {truncatePath(resultsFile)}
        </span>
      </TooltipTrigger>
      <TooltipContent>{resultsFile}</TooltipContent>
    </Tooltip>
  );
}
