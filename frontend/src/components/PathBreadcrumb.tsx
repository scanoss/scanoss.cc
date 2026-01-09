// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2026 SCANOSS.COM
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

import { useMemo } from 'react';

import { cn, getPathSegments } from '@/lib/utils';

import { Button } from './ui/button';

interface PathBreadcrumbProps {
  filePath: string;
  selectedIndex: number;
  onSelect: (index: number) => void;
}

export default function PathBreadcrumb({ filePath, selectedIndex, onSelect }: PathBreadcrumbProps) {
  const segments = useMemo(() => getPathSegments(filePath), [filePath]);

  return (
    <div className="flex flex-wrap items-center font-mono text-sm">
      {segments.map((segment, index) => (
        <div key={index} className="flex items-center">
          <Button
            variant="ghost"
            onClick={() => onSelect(index)}
            className={cn(
              'h-auto px-1.5 py-0.5 text-sm',
              selectedIndex === index && 'bg-accent text-accent-foreground font-medium'
            )}
          >
            {segment}
          </Button>
          {index < segments.length - 1 && <span className="text-muted-foreground">/</span>}
        </div>
      ))}
    </div>
  );
}
