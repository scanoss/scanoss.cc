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

import clsx from 'clsx';
import { ArrowDownNarrowWide, ArrowUpNarrowWide, FileText, Percent } from 'lucide-react';
import { ReactNode } from 'react';

import useResultsStore from '@/modules/results/stores/useResultsStore';

import { Button } from './ui/button';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from './ui/dropdown-menu';

type SortOption = {
  value: string;
  label: string;
  icon: ReactNode;
  description: string;
};

const sortOptions: SortOption[] = [
  {
    value: 'match_percentage',
    label: 'Match Percentage',
    icon: <Percent className="h-4 w-4" />,
    description: 'Sort by highest match percentage',
  },
  {
    value: 'path',
    label: 'File Path',
    icon: <FileText className="h-4 w-4" />,
    description: 'Sort alphabetically by file path',
  },
];

export default function SortSelector() {
  const sort = useResultsStore((state) => state.sort);
  const setSort = useResultsStore((state) => state.setSort);

  const selectedOption = sortOptions.find((option) => option.value === sort.option) || sortOptions[0];

  return (
    <div className="flex flex-wrap items-center gap-2">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button size="sm" variant="outline" className="justify-start gap-2 whitespace-normal">
            {selectedOption.icon}
            <span className="text-left">Sort by: {selectedOption.label}</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="start">
          {sortOptions.map((option) => (
            <DropdownMenuItem
              key={option.value}
              onClick={() => setSort(option.value, sort.order)}
              className={clsx('gap-2', {
                'bg-primary text-primary-foreground': sort.option === option.value,
              })}
            >
              {option.icon}
              <div className="flex flex-col">
                <span>{option.label}</span>
                <span className="text-xs text-muted-foreground">{option.description}</span>
              </div>
            </DropdownMenuItem>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
      <Button size="sm" variant="outline" className="gap-2" onClick={() => setSort(sort.option, sort.order === 'asc' ? 'desc' : 'asc')}>
        {sort.order === 'asc' ? <ArrowUpNarrowWide className="h-4 w-4" /> : <ArrowDownNarrowWide className="h-4 w-4" />}
      </Button>
    </div>
  );
}
