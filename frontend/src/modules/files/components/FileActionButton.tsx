import { useQuery } from '@tanstack/react-query';
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
import { FilterAction, filterActionLabelMap } from '@/modules/results/domain';
import ResultService from '@/modules/results/infra/service';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import useLocalFilePath from '../hooks/useLocalFilePath';

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
  const { ask } = useConfirm();
  const { toast } = useToast();
  const localFilePath = useLocalFilePath();
  const navigate = useNavigate();

  const [dropdownOpen, setDropdownOpen] = useState(false);

  const handleCompleteResult = useResultsStore(
    (state) => state.handleCompleteResult
  );

  const { data: component } = useQuery({
    queryKey: ['component', localFilePath],
    queryFn: () => ResultService.getComponent(localFilePath),
  });

  const handleFilterComponent = async (
    path: string | undefined,
    purl: string
  ) => {
    return handleCompleteResult(path, purl, action, localFilePath)
      .then((nextResultRoute) => {
        if (nextResultRoute) {
          navigate(nextResultRoute);
        }
      })
      .catch((e) => {
        toast({
          title: 'Error',
          variant: 'destructive',
          description: e instanceof Error ? e.message : 'An error occurred.',
        });
      });
  };

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
        <DropdownMenuItem
          onClick={async () => {
            await handleFilterComponent(
              localFilePath,
              component?.purl?.[0] as string
            );
          }}
        >
          <div className="flex flex-col">
            <span className="text-sm">File</span>
            <span className="text-xs text-muted-foreground">
              {localFilePath}
            </span>
          </div>
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={async () => {
            const confirm = await ask(
              `This action will ${action} all components with the same purl: ${component?.purl?.[0]}. Are you sure?`
            );
            if (confirm) {
              await handleFilterComponent(
                undefined,
                component?.purl?.[0] as string
              );
            }
          }}
          disabled={!component?.purl?.[0]}
        >
          <div className="flex flex-col">
            <span className="text-sm">Component</span>
            <span className="text-xs text-muted-foreground">
              {component?.purl?.[0]}
            </span>
          </div>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
