import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
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
import useDebounce from '@/hooks/useDebounce';
import useQueryState from '@/hooks/useQueryState';
import {
  FilterAction,
  filterActionLabelMap,
  MatchType,
} from '@/modules/results/domain';
import ResultService from '@/modules/results/infra/service';
import { useResults } from '@/modules/results/providers/ResultsProvider';

import useLocalFilePath from '../hooks/useLocalFilePath';
import FileService from '../infra/service';

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
  const { toast } = useToast();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { ask } = useConfirm();

  const { handleConfirmResult } = useResults();
  const localFilePath = useLocalFilePath();
  const [dropdownOpen, setDropdownOpen] = useState(false);

  const [filterByMatchType] = useQueryState<MatchType | 'all'>(
    'matchType',
    'all'
  );

  const [query] = useQueryState<string>('q', '');
  const debouncedQuery = useDebounce<string>(query, 300);

  const { data: component } = useQuery({
    queryKey: ['component', localFilePath],
    queryFn: () => ResultService.getComponent(localFilePath),
  });

  const { mutate: filterComponent } = useMutation({
    mutationFn: ({ path, purl }: { path: string | undefined; purl: string }) =>
      FileService.filterComponentByPath({
        action,
        path,
        purl,
      }),
    onSuccess: async () => {
      toast({
        title: 'Success',
        description: `Your changes have been successfully saved.`,
      });

      await queryClient.refetchQueries({
        queryKey: ['results', filterByMatchType, debouncedQuery],
      });

      const nextResultRoute = handleConfirmResult(localFilePath);
      if (nextResultRoute) {
        navigate(nextResultRoute);
      }
    },
    onMutate: async ({ path, purl }) => {
      if (path) return;

      const confirm = await ask(
        `This action will ${action} all components with the same purl: ${purl}. Are you sure?`
      );
      if (!confirm) {
        return Promise.reject(new Error('user_cancel'));
      }
    },
    onError: (e) => {
      if (e.message !== 'user_cancel') {
        toast({
          title: 'Error',
          variant: 'destructive',
          description: e instanceof Error ? e.message : 'An error occurred.',
        });
      }
    },
  });

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
        <DropdownMenuItem
          onClick={() => filterComponent({ path: localFilePath, purl: purl! })}
        >
          <div className="flex flex-col">
            <span className="text-sm">File</span>
            <span className="text-xs text-muted-foreground">
              {localFilePath}
            </span>
          </div>
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => filterComponent({ path: undefined, purl: purl! })}
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
