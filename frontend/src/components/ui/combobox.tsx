'use client';

import clsx from 'clsx';
import { Check, ChevronsUpDown, Loader } from 'lucide-react';
import { useState } from 'react';

import { Button } from '@/components/ui/button';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { cn } from '@/lib/utils';

interface ComboboxProps {
  buttonClassName?: string;
  emptyText: string;
  isLoading?: boolean;
  options: { value: string; label: string }[];
  placeholder: string;
  popoverClassName?: string;
}

export function Combobox({
  buttonClassName,
  emptyText,
  isLoading,
  options,
  placeholder,
  popoverClassName,
}: ComboboxProps) {
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState('');

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild disabled={isLoading}>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className={clsx('w-full justify-between', buttonClassName)}
        >
          {value
            ? options.find((option) => option.value === value)?.label
            : placeholder}
          {isLoading ? (
            <Loader className="h-4 w-4 animate-spin" />
          ) : (
            <ChevronsUpDown className="h-4 w-4" />
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent className={clsx('w-full p-0', popoverClassName)}>
        <Command>
          <CommandInput placeholder={placeholder} />
          <CommandList>
            <CommandEmpty>{emptyText}</CommandEmpty>
            <CommandGroup>
              {options.map((option) => (
                <CommandItem
                  key={option.value}
                  value={option.value}
                  onSelect={(currentValue) => {
                    setValue(currentValue === value ? '' : currentValue);
                    setOpen(false);
                  }}
                >
                  <Check
                    className={cn(
                      'mr-2 h-4 w-4',
                      value === option.value ? 'opacity-100' : 'opacity-0'
                    )}
                  />
                  {option.label}
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
