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
