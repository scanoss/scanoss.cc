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

import { Search } from 'lucide-react';

import { Input } from '@/components/ui/input';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { useMenuEvent } from '@/hooks/useMenuEvent';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import { entities } from '../../wailsjs/go/models';

export default function ResultSearchBar({ searchInputRef }: { searchInputRef: React.RefObject<HTMLInputElement> }) {
  const query = useResultsStore((state) => state.query);
  const setQuery = useResultsStore((state) => state.setQuery);
  const setSelectedResults = useResultsStore((state) => state.setSelectedResults);

  useKeyboardShortcut(KEYBOARD_SHORTCUTS.focusSearch.keys, () => searchInputRef.current?.focus());

  useMenuEvent(entities.Action.FocusSearch, () => searchInputRef.current?.focus());

  return (
    <div className="relative grid w-full items-center gap-1.5">
      <Search className="absolute left-3 top-1/2 h-3 w-3 -translate-y-1/2 transform text-muted-foreground" />
      <Input
        ref={searchInputRef}
        className="pl-7"
        name="q"
        onChange={(e) => {
          setSelectedResults([]);
          setQuery(e.target.value);
        }}
        placeholder="Search"
        value={query}
        type="text"
      />
      <span className="absolute right-3 text-xs leading-none text-muted-foreground">⌘ + F</span>
    </div>
  );
}
