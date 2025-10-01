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
import { Check, ChevronsUpDown, X } from 'lucide-react';
import { useCallback, useState } from 'react';

import { cn } from '@/lib/utils';

import { entities } from '../../wailsjs/go/models';
import { GetAll as GetAllLicenses } from '../../wailsjs/go/service/LicenseServiceImpl';
import { Button } from './ui/button';
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from './ui/command';
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover';

interface SelectLicenseListProps {
  matchedLicenses?: entities.License[];
  onSelect: (value: string | undefined) => void;
}

function LicenseItem({ license, value }: { license: entities.License; value: string }) {
  return (
    <>
      <Check className={cn('mr-2 h-4 w-4', license.licenseId === value ? 'opacity-100' : 'opacity-0')} />
      <div>
        <p>{license.licenseId}</p>
        <p className="text-xs text-muted-foreground">{license.name}</p>
      </div>
    </>
  );
}

export default function SelectLicenseList({ matchedLicenses, onSelect }: SelectLicenseListProps) {
  const [popoverOpen, setPopoverOpen] = useState(false);
  const [value, setValue] = useState<string>();
  const { data: licenses } = useQuery({
    queryKey: ['licenses'],
    queryFn: GetAllLicenses,
  });

  const handleSelect = useCallback((value: string) => {
    setValue(value);
    onSelect(value);
    setPopoverOpen(false);
  }, []);

  const handleClearSelection = useCallback(() => {
    setValue(undefined);
    onSelect(undefined);
  }, []);

  return (
    <Popover open={popoverOpen} onOpenChange={setPopoverOpen} modal>
      <div className="flex w-full gap-2">
        <PopoverTrigger asChild className="flex-1">
          <Button variant="outline" role="combobox" className={cn('justify-between', !value && 'text-muted-foreground')}>
            {value ? licenses?.find((license) => license.licenseId === value)?.licenseId : 'Select license'}
            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        {value && (
          <Button size="icon" variant="ghost" onClick={handleClearSelection}>
            <X className="h-4 w-4" />
          </Button>
        )}
      </div>
      <PopoverContent className="min-w-[420px] p-0">
        <Command>
          <CommandInput placeholder="Search license..." />
          <CommandList>
            <CommandEmpty>No licenses found.</CommandEmpty>
            {matchedLicenses && matchedLicenses.length > 0 ? (
              <CommandGroup heading="Matched">
                {matchedLicenses.map((license) => (
                  <CommandItem value={license.licenseId} key={license.licenseId} onSelect={handleSelect}>
                    <LicenseItem license={license} value={value || ''} />
                  </CommandItem>
                ))}
              </CommandGroup>
            ) : null}
            <CommandGroup heading="Catalogued">
              {licenses?.map((license) => (
                <CommandItem value={license.licenseId} key={license.licenseId} onSelect={handleSelect}>
                  <LicenseItem license={license} value={value || ''} />
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
