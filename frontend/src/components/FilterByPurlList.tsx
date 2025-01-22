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
