import clsx from 'clsx';
import { ChevronDown } from 'lucide-react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

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
import { ScrollArea } from '@/components/ui/scroll-area';
import { useToast } from '@/components/ui/use-toast';
import { useConfirm } from '@/hooks/useConfirm';
import { useInputPrompt } from '@/hooks/useInputPrompt';
import { FilterAction, filterActionLabelMap } from '@/modules/results/domain';
import useResultsStore from '@/modules/results/stores/useResultsStore';

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
  const { prompt } = useInputPrompt();
  const { toast } = useToast();
  const navigate = useNavigate();

  const [dropdownOpen, setDropdownOpen] = useState(false);

  const selectedResults = useResultsStore((state) => state.selectedResults);
  const handleCompleteResult = useResultsStore(
    (state) => state.handleCompleteResult
  );

  const handleAddFilter = async (
    filterType: 'by_file' | 'by_purl',
    withComment?: boolean
  ) => {
    try {
      let comment: string | undefined;

      if (withComment) {
        comment = await prompt({
          title: 'Add comments',
          input: {
            defaultValue: '',
            type: 'textarea',
          },
        });
        if (!comment) return;
      }

      return filterActions[filterType](comment);
    } catch (e) {
      console.error(e);
      toast({
        title: 'Error',
        variant: 'destructive',
        description:
          'An error ocurred while adding comments. Please try again.',
      });
    }
  };

  const handleFilterComponentByPurl = async (comment?: string) => {
    const confirm = await ask(
      <div>
        <p>
          This action will {action} all matches with{' '}
          {selectedResults.length > 1
            ? `the following PURLs: `
            : `the same PURL:`}
        </p>

        {selectedResults.length > 1 ? (
          <ScrollArea className="py-2">
            <ul className="max-h-[200px] list-disc pl-6">
              {selectedResults.map((result) => (
                <li key={result.path}>{result.purl}</li>
              ))}
            </ul>
          </ScrollArea>
        ) : (
          <p>{selectedResults[0]?.purl}</p>
        )}
      </div>
    );

    if (confirm) {
      return handleFilterComponent('by_purl', comment);
    }
  };

  const handleFilterComponentByFile = async (comment?: string) => {
    return handleFilterComponent('by_file', comment);
  };

  const handleFilterComponent = async (
    filterBy: 'by_file' | 'by_purl',
    comment?: string
  ) => {
    const nextResultRoute = await handleCompleteResult({
      comment,
      action,
      filterBy,
    });

    if (nextResultRoute) {
      navigate(nextResultRoute);
    }
  };

  const filterActions: Record<
    'by_file' | 'by_purl',
    (comments?: string) => Promise<void>
  > = {
    by_file: handleFilterComponentByFile,
    by_purl: handleFilterComponentByPurl,
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
        <DropdownMenuGroup>
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>File</DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent>
                <DropdownMenuItem
                  onClick={() => handleAddFilter('by_file', true)}
                >
                  <span className="first-letter:uppercase">{`${action} with Comments`}</span>
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => handleAddFilter('by_file')}>
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
                  onClick={() => handleAddFilter('by_purl', true)}
                >
                  <span className="first-letter:uppercase">{`${action} with Comments`}</span>
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => handleAddFilter('by_purl')}>
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
