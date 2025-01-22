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

import clsx from 'clsx';
import { File, Github } from 'lucide-react';

import useSelectedResult from '@/hooks/useSelectedResult';
import { FilterAction } from '@/modules/components/domain';
import { stateInfoPresentation } from '@/modules/results/domain';

interface FileInfoCardProps {
  title: string;
  subtitle: string | undefined;
  fileType: 'local' | 'remote';
}

export default function FileInfoCard({ title, subtitle, fileType }: FileInfoCardProps) {
  const result = useSelectedResult();

  const filterConfig = result?.filter_config;

  const presentation = stateInfoPresentation[filterConfig?.action as FilterAction];

  const shouldShowStateInfo =
    (fileType === 'local' && filterConfig?.type === 'by_file') || (fileType === 'remote' && filterConfig?.type === 'by_purl');

  return (
    <div
      className={clsx(
        'flex justify-between border p-3 text-sm',
        shouldShowStateInfo ? presentation?.stateInfoContainerStyles : 'border-l-0 border-t-0',
        fileType === 'remote' && !shouldShowStateInfo && 'border-r-0'
      )}
    >
      <div>
        <div className="flex items-center gap-1">
          {fileType === 'local' ? <File className="h-3 w-3" /> : <Github className="h-3 w-3" />}
          <span className="font-semibold">{title}</span>
        </div>
        <p className="text-muted-foreground">{subtitle ?? '-'}</p>
      </div>
      {shouldShowStateInfo && (
        <div className="text-xs">
          <p className={presentation?.stateInfoTextStyles}>{presentation?.label}</p>
        </div>
      )}
    </div>
  );
}
