import { useMutation } from '@tanstack/react-query';
import { ChevronDown } from 'lucide-react';

import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Component } from '@/modules/results/domain';

import { FilterAction } from '../domain';
import useLocalFilePath from '../hooks/useLocalFilePath';
import FileService from '../infra/service';

interface FileActionButtonProps {
  action: FilterAction;
  component: Component;
  icon: React.ReactNode;
}

export default function FileActionButton({
  action,
  component,
  icon,
}: FileActionButtonProps) {
  const localFilePath = useLocalFilePath();
  const purl = component.purl?.[0];

  if (!purl || !localFilePath) {
    return null;
  }

  const { mutate } = useMutation({
    mutationFn: () =>
      FileService.filterComponentByPath({
        action,
        path: localFilePath,
        purl,
      }),
  });

  const handleFilterByFile = () => mutate();

  return (
    <DropdownMenu>
      <div className="group flex h-full">
        <Button
          variant="ghost"
          size="lg"
          className="h-full w-14 rounded-none group-hover:bg-accent group-hover:text-accent-foreground"
        >
          <div className="flex flex-col items-center justify-center gap-1">
            <span className="text-xs first-letter:uppercase">{action}</span>
            {icon}
          </div>
        </Button>
        <DropdownMenuTrigger>
          <div className="flex h-full w-4 cursor-pointer items-center transition-colors hover:bg-accent">
            <ChevronDown className="h-4 w-4 stroke-muted-foreground hover:stroke-accent-foreground" />
          </div>
        </DropdownMenuTrigger>
      </div>
      <DropdownMenuContent className="max-w-[400px]">
        <DropdownMenuLabel>
          <span className="text-xs font-normal text-muted-foreground">
            Here should go a description of the action explaining what it does
          </span>
        </DropdownMenuLabel>
        <DropdownMenuSeparator className="bg-border" />
        <DropdownMenuItem onClick={handleFilterByFile}>
          <div className="flex flex-col">
            <span className="text-sm">File</span>
            <span className="text-xs text-muted-foreground">
              {localFilePath}
            </span>
          </div>
          <DropdownMenuShortcut>⌘⇧F</DropdownMenuShortcut>
        </DropdownMenuItem>
        <DropdownMenuItem>
          <div className="flex flex-col">
            <span className="text-sm">Component</span>
            <span className="text-xs text-muted-foreground">{purl}</span>
          </div>
          <DropdownMenuShortcut>⌘⇧C</DropdownMenuShortcut>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
