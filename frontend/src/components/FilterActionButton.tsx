import clsx from 'clsx';
import { ChevronDown } from 'lucide-react';
import { useState } from 'react';

import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  FilterAction,
  filterActionLabelMap,
  OnAddFilterArgs,
} from '@/modules/components/domain';

interface FilterActionProps {
  action: FilterAction;
  description: string;
  icon: React.ReactNode;
  onAdd: (args: OnAddFilterArgs) => Promise<void>;
}

export default function FilterActionButton({
  action,
  description,
  icon,
  onAdd,
}: FilterActionProps) {
  const [dropdownOpen, setDropdownOpen] = useState(false);

  return (
    <DropdownMenu onOpenChange={setDropdownOpen}>
      <div className="group flex h-full">
        <DropdownMenuTrigger asChild>
          <Button
            variant="ghost"
            size="lg"
            className="h-full w-14 rounded-none enabled:group-hover:bg-accent enabled:group-hover:text-accent-foreground"
          >
            <div className="flex flex-col items-center justify-center gap-1">
              <span className="text-xs">{filterActionLabelMap[action]}</span>
              {icon}
            </div>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuTrigger asChild>
          <button className="flex h-full w-4 items-center outline-none transition-colors enabled:hover:bg-accent">
            <ChevronDown
              className={clsx(
                'h-4 w-4 stroke-muted-foreground transition-transform enabled:hover:stroke-accent-foreground',
                dropdownOpen && 'rotate-180 transform'
              )}
            />
          </button>
        </DropdownMenuTrigger>
      </div>
      <DropdownMenuContent className="max-w-[400px]">
        <DropdownMenuLabel>
          <span className="text-xs font-normal text-muted-foreground">
            {description}
          </span>
        </DropdownMenuLabel>
        <DropdownMenuSeparator className="bg-border" />
        <DropdownMenuGroup>
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>File</DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent>
                <DropdownMenuItem
                  onClick={() =>
                    onAdd({
                      action,
                      filterBy: 'by_file',
                      withComment: true,
                    })
                  }
                >
                  <span className="first-letter:uppercase">{`${action} with Comments`}</span>
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() =>
                    onAdd({
                      action,
                      filterBy: 'by_file',
                    })
                  }
                >
                  <span className="first-letter:uppercase">{`${action} without Comments`}</span>
                </DropdownMenuItem>
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>
        </DropdownMenuGroup>
        <DropdownMenuGroup>
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>Component</DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent>
                <DropdownMenuItem
                  onClick={() =>
                    onAdd({
                      action,
                      filterBy: 'by_purl',
                      withComment: true,
                    })
                  }
                >
                  <span className="first-letter:uppercase">{`${action} with Comments`}</span>
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() =>
                    onAdd({
                      action,
                      filterBy: 'by_purl',
                    })
                  }
                >
                  <span className="first-letter:uppercase">{`${action} without Comments`}</span>
                </DropdownMenuItem>
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
