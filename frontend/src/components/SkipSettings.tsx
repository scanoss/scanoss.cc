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

import { HelpCircle } from 'lucide-react';
import useResizeObserver from 'use-resize-observer';

import ResultsTree from './ResultsTree';
import { SKIP_STATE_DESCRIPTIONS, SKIP_STATES } from './ResultsTreeNode';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';

export default function SkipSettings() {
  const { ref, width, height } = useResizeObserver();

  return (
    <div className="flex h-full flex-col gap-4">
      <p className="flex items-center gap-1 text-sm text-muted-foreground">
        Configure the files, directories and extensions to skip while scanning{' '}
        <Tooltip>
          <TooltipTrigger asChild>
            <HelpCircle className="h-3 w-3 hover:cursor-pointer hover:text-primary-foreground" />
          </TooltipTrigger>
          <TooltipContent side="bottom" align="start" className="w-60">
            <div className="flex flex-col gap-2 p-1">
              <p className="text-xs font-semibold">Scanning Status Indicators:</p>
              <div className="flex items-center gap-2">
                <div className="h-3 w-3 rounded-full bg-green-500"></div>
                <span className="text-xs">{SKIP_STATE_DESCRIPTIONS[SKIP_STATES.INCLUDED]}</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="h-3 w-3 rounded-full bg-red-500"></div>
                <span className="text-xs">{SKIP_STATE_DESCRIPTIONS[SKIP_STATES.EXCLUDED]}</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="h-3 w-3 rounded-full bg-yellow-500"></div>
                <span className="text-xs">{SKIP_STATE_DESCRIPTIONS[SKIP_STATES.MIXED]}</span>
              </div>
              <div className="mt-1 text-xs text-muted-foreground">
                <p>Right-click on any file or folder to change its scanning status.</p>
              </div>
            </div>
          </TooltipContent>
        </Tooltip>
      </p>

      <div className="flex flex-1 gap-4">
        <div className="flex flex-grow" style={{ minBlockSize: 0 }} ref={ref}>
          <ResultsTree width={width} height={height} />
        </div>
      </div>
    </div>
  );
}
