import { useQuery } from '@tanstack/react-query';
import { Check, ChevronsUpDown, X } from 'lucide-react';
import { useCallback, useState } from 'react';

import { cn } from '@/lib/utils';

import { GetAll as GetAllLicenses } from '../../wailsjs/go/service/LicenseServiceImpl';
import { Button } from './ui/button';
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from './ui/command';
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover';
import { ScrollArea } from './ui/scroll-area';

interface SelectLicenseListProps {
  onSelect: (value: string | undefined) => void;
}

export default function SelectLicenseList({ onSelect }: SelectLicenseListProps) {
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
    <Popover open={popoverOpen} onOpenChange={setPopoverOpen}>
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
            <CommandGroup>
              <ScrollArea className="h-52">
                {licenses?.map((license) => (
                  <CommandItem value={license.licenseId} key={license.licenseId} onSelect={handleSelect}>
                    <Check className={cn('mr-2 h-4 w-4', license.licenseId === value ? 'opacity-100' : 'opacity-0')} />
                    <div>
                      <p>{license.licenseId}</p>
                      <p className="text-xs text-muted-foreground">{license.name}</p>
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
