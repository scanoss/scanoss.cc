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
            {/* @ts-expect-error wails type issue */}
            {Object.entries(shortcuts).map(([shortcutKey, shortcut]) => {
              return (
                <CommandItem key={shortcutKey}>
                  {/* @ts-expect-error wails type issue */}
                  <span>{shortcut.name}</span>
                  <CommandShortcut className="flex items-center gap-1">
                    {/* @ts-expect-error wails type issue */}
                    {getShortcutDisplay(shortcut.keys, modifierKey.label).map((keyCombo, index) => (
                      <>
                        <span
                          key={index}
                          className="mr-1 flex min-w-5 items-center justify-center rounded-sm bg-muted p-0.5 text-xs capitalize leading-tight text-muted-foreground"
                        >
                          {keyCombo}
                        </span>
                        {/* @ts-expect-error wails type issue */}
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
