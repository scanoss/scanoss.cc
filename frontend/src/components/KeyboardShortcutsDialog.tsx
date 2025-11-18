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

import { useQuery } from '@tanstack/react-query';

import useEnvironment from '@/hooks/useEnvironment';
import { getShortcutDisplay } from '@/lib/shortcuts';

import { GetGroupedShortcuts } from '../../wailsjs/go/service/KeyboardServiceInMemoryImpl';
import { CommandDialog, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList, CommandShortcut } from './ui/command';

export default function KeyboardShortcutsDialog({ open, onOpenChange }: { open: boolean; onOpenChange: () => void }) {
  const { modifierKey } = useEnvironment();
  const { data } = useQuery({
    queryKey: ['shortcuts'],
    queryFn: GetGroupedShortcuts,
  });

  if (!data) {
    return null;
  }

  return (
    <CommandDialog open={open} onOpenChange={onOpenChange}>
      <CommandInput placeholder="Filter shortcuts" />
      <CommandList>
        <CommandEmpty>No results found.</CommandEmpty>
        {Object.entries(data).map(([group, shortcuts]) => (
          <CommandGroup key={group} heading={group}>
            {Object.entries(shortcuts).map(([shortcutKey, shortcut]) => {
              return (
                <CommandItem key={shortcutKey}>
                  <span>{shortcut.name}</span>
                  <CommandShortcut className="flex items-center gap-1">
                    {getShortcutDisplay(shortcut.keys, modifierKey.label).map((keyCombo, index) => (
                      <>
                        <span
                          key={index}
                          className="mr-1 flex min-w-5 items-center justify-center rounded-sm bg-muted p-0.5 text-xs capitalize leading-tight text-muted-foreground"
                        >
                          {keyCombo}
                        </span>
                        {index < getShortcutDisplay(shortcut.keys, modifierKey.label).length - 1 ? 'or' : null}
                      </>
                    ))}
                  </CommandShortcut>
                </CommandItem>
              );
            })}
          </CommandGroup>
        ))}
      </CommandList>
    </CommandDialog>
  );
}
