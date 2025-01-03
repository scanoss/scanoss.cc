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

import { FilterAction } from '@/modules/components/domain';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import { ScrollArea } from './ui/scroll-area';

interface FilterByPurlListProps {
  action: FilterAction;
}

export default function FilterByPurlList({ action }: FilterByPurlListProps) {
  const selectedResults = useResultsStore((state) => state.selectedResults);

  return (
    <div>
      <p>
        This action will {action} all matches with {selectedResults.length > 1 ? `the following PURLs: ` : `the same PURL:`}
      </p>

      {selectedResults.length > 1 ? (
        <ScrollArea className="py-2">
          <ul className="max-h-[200px] list-disc pl-6">
            {selectedResults.map((result) => (
              <li key={result.path}>{result.detected_purl}</li>
            ))}
          </ul>
        </ScrollArea>
      ) : (
        <p>{selectedResults[0]?.detected_purl}</p>
      )}
    </div>
  );
}
