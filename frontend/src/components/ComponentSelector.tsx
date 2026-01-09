// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2024 SCANOSS.COM
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

import { Check, ChevronsUpDown, Plus, Search } from 'lucide-react';
import { useState } from 'react';

import useEnvironment from '@/hooks/useEnvironment';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { getShortcutDisplay, KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import { cn } from '@/lib/utils';

import { entities } from '../../wailsjs/go/models';
import ShortcutBadge from './ShortcutBadge';
import { Button } from './ui/button';
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList, CommandSeparator } from './ui/command';
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover';
import { ScrollArea } from './ui/scroll-area';

interface ComponentSelectorProps {
  components: entities.DeclaredComponent[];
  selectedComponent: entities.DeclaredComponent | null;
  onSelect: (component: entities.DeclaredComponent) => void;
  onAddNew: () => void;
  onSearchOnline: (searchTerm: string) => void;
}

export default function ComponentSelector({
  components,
  selectedComponent,
  onSelect,
  onAddNew,
  onSearchOnline,
}: ComponentSelectorProps) {
  const { modifierKey } = useEnvironment();
  const [open, setOpen] = useState(false);
  const [searchValue, setSearchValue] = useState('');

  const handleSearchOnline = () => {
    onSearchOnline(searchValue);
    setOpen(false);
  };

  useKeyboardShortcut(KEYBOARD_SHORTCUTS.confirm.keys, handleSearchOnline, {
    enableOnFormTags: true,
    enabled: open && !!searchValue,
  });

  return (
    <Popover open={open} onOpenChange={setOpen} modal>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          className={cn('justify-between', !selectedComponent && 'text-muted-foreground')}
        >
          {selectedComponent ? selectedComponent.name : 'Select component'}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-0">
        <Command>
          <CommandInput placeholder="Search component..." value={searchValue} onValueChange={setSearchValue} />
          <CommandList>
            <CommandEmpty>No components found.</CommandEmpty>
            <CommandGroup>
              <CommandItem
                onSelect={() => {
                  onAddNew();
                  setOpen(false);
                }}
              >
                <Plus className="mr-2 h-4 w-4" />
                Add new component
              </CommandItem>
              {searchValue && (
                <CommandItem onSelect={handleSearchOnline} className="flex items-center justify-between gap-2">
                  <div className="flex items-center gap-2">
                    <Search className="h-3 w-3" />
                    Search &ldquo;{searchValue}&rdquo; online
                  </div>
                  <ShortcutBadge shortcut={getShortcutDisplay(KEYBOARD_SHORTCUTS.confirm.keys, modifierKey.label)[0]} />
                </CommandItem>
              )}
            </CommandGroup>
            <CommandSeparator />
            <CommandGroup>
              <ScrollArea className="h-52">
                {components?.map((component) => (
                  <CommandItem
                    key={component.purl}
                    onSelect={() => {
                      onSelect(component);
                      setOpen(false);
                    }}
                  >
                    <Check
                      className={cn('mr-2 h-4 w-4', component.purl === selectedComponent?.purl ? 'opacity-100' : 'opacity-0')}
                    />
                    <div>
                      <p>{component.name}</p>
                      <p className="text-xs text-muted-foreground">{component.purl}</p>
                    </div>
                  </CommandItem>
                ))}
              </ScrollArea>
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
