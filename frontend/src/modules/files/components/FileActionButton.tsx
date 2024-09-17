import { useMutation, useQuery } from '@tanstack/react-query';
import clsx from 'clsx';
import { ChevronDown } from 'lucide-react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useToast } from '@/components/ui/use-toast';
import { useConfirm } from '@/hooks/useConfirm';
import ResultService from '@/modules/results/infra/service';
import { useResults } from '@/modules/results/providers/ResultsProvider';

import { FilterAction, filterActionLabelMap } from '../domain';
import useLocalFilePath from '../hooks/useLocalFilePath';
import FileService from '../infra/service';

const getActionPersistKey = (action: FilterAction) =>
  `filter-action-${action}-dont-ask-again`;
const isFilterActionPersisted = (action: FilterAction) =>
  localStorage.getItem(getActionPersistKey(action)) === 'true';
const persistAction = (action: FilterAction) =>
  localStorage.setItem(getActionPersistKey(action), 'true');

interface FileActionButtonProps {
  action: FilterAction;
  description: string;
  icon: React.ReactNode;
}

export default function FileActionButton({
  action,
  description,
  icon,
}: FileActionButtonProps) {
  const navigate = useNavigate();
  const { toast } = useToast();
  const { ask } = useConfirm();
  const { handleStageResult } = useResults();
  const localFilePath = useLocalFilePath();
  const [dropdownOpen, setDropdownOpen] = useState(false);

  const { data: component } = useQuery({
    queryKey: ['component', localFilePath],
    queryFn: () => ResultService.getComponent(localFilePath),
  });

  const { mutate } = useMutation({
    mutationFn: ({ path, purl }: { path: string; purl: string }) =>
      FileService.filterComponentByPath({
        action,
        path,
        purl,
      }),
    onSuccess: () => {
      const nextResultRoute = handleStageResult(localFilePath);
      if (nextResultRoute) {
        navigate(nextResultRoute);
      }
    },
    onError: (e) => {
      toast({
        title: 'Error',
        variant: 'destructive',
        description: e instanceof Error ? e.message : 'An error occurred.',
      });
    },
  });

  const handleFilter = async ({
    path = '',
    purl = '',
  }: {
    path?: string;
    purl?: string;
  }) => {
    const isPersisted = isFilterActionPersisted(action);
    if (isPersisted) {
      mutate({ path, purl });
    } else {
      const res = await ask(
        // TODO: Define proper descriptions
        'You are about to filter this file for future scans. Are you sure you want to proceed?',
        () => persistAction(action)
      );
      if (res) {
        mutate({ path, purl });
      }
    }
  };

  const purl = component?.purl?.[0];

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
        <DropdownMenuItem onClick={() => handleFilter({ path: localFilePath })}>
          <div className="flex flex-col">
            <span className="text-sm">File</span>
            <span className="text-xs text-muted-foreground">
              {localFilePath}
            </span>
          </div>
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => handleFilter({ purl })}
          disabled={!purl}
        >
          <div className="flex flex-col">
            <span className="text-sm">Component</span>
            <span className="text-xs text-muted-foreground">{purl}</span>
          </div>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
