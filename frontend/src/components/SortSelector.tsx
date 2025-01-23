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

import { ArrowDownWideNarrow, FileText, Package, Percent, Scale } from 'lucide-react';
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
  {
    value: 'component_name',
    label: 'Component Name',
    icon: <Package className="h-4 w-4" />,
    description: 'Sort by component name',
  },
  {
    value: 'license',
    label: 'License',
    icon: <Scale className="h-4 w-4" />,
    description: 'Sort by license type',
  },
];

export default function SortSelector() {
  const sortBy = useResultsStore((state) => state.sortBy);
  const setSortBy = useResultsStore((state) => state.setSortBy);

  const selectedOption = sortOptions.find((option) => option.value === sortBy) || sortOptions[0];

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" className="w-full justify-start gap-2">
          <ArrowDownWideNarrow className="h-4 w-4" />
          <span>Sort by: {selectedOption.label}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" className="w-[200px]">
        {sortOptions.map((option) => (
          <DropdownMenuItem key={option.value} onClick={() => setSortBy(option.value)} className="gap-2">
            {option.icon}
            <div className="flex flex-col">
              <span>{option.label}</span>
              <span className="text-xs text-muted-foreground">{option.description}</span>
            </div>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
