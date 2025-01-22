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

import { Label } from '@/components/ui/label';
import { MatchType } from '@/modules/results/domain';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import { matchTypeIconMap } from './Sidebar';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './ui/select';

export default function MatchTypeSelector() {
  const filterByMatchType = useResultsStore((state) => state.filterByMatchType);
  const setFilterByMatchType = useResultsStore((state) => state.setFilterByMatchType);

  return (
    <div className="flex flex-col gap-1">
      <Label className="text-xs">Filter by match type</Label>
      <Select onValueChange={(value) => setFilterByMatchType(value as MatchType)} value={filterByMatchType}>
        <SelectTrigger className="w-full">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value={MatchType.File}>
            <div className="flex items-center gap-1">
              <span>{matchTypeIconMap[MatchType.File]}</span>
              <span>File</span>
            </div>
          </SelectItem>
          <SelectItem value={MatchType.Snippet}>
            <div className="flex items-center gap-1">
              <span>{matchTypeIconMap[MatchType.Snippet]}</span>
              <span>Snippet</span>
            </div>
          </SelectItem>
          <SelectItem value="all">All</SelectItem>
        </SelectContent>
      </Select>
    </div>
  );
}
