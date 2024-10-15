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
import { useConfirm } from '@/hooks/useConfirm';
import { useInputPrompt } from '@/hooks/useInputPrompt';
import {
  FilterAction,
  filterActionLabelMap,
} from '@/modules/components/domain';
import useComponentFilterStore from '@/modules/components/stores/useComponentFilterStore';
import useResultsStore from '@/modules/results/stores/useResultsStore';

import { ScrollArea } from './ui/scroll-area';

interface FilterActionProps {
  action: FilterAction;
  description: string;
  icon: React.ReactNode;
  onAdd: () => Promise<void> | void;
}

export default function FilterActionButton({
  action,
  description,
  icon,
  onAdd,
}: FilterActionProps) {
  const [dropdownOpen, setDropdownOpen] = useState(false);

  const { prompt } = useInputPrompt();
  const { ask } = useConfirm();

  const setAction = useComponentFilterStore((state) => state.setAction);
  const setComment = useComponentFilterStore((state) => state.setComment);
  const setFilterBy = useComponentFilterStore((state) => state.setFilterBy);

  const selectedResults = useResultsStore((state) => state.selectedResults);

  const onSelectOption = async (
    action: FilterAction,
    filterBy: 'by_file' | 'by_purl',
    withComment: boolean
  ) => {
    let comment: string | undefined;

    if (withComment) {
      comment = await handleGetComment();
    }

    if (filterBy === 'by_purl') {
      const confirm = await handleConfirmByPurl();

      if (!confirm) return;
    }

    setAction(action);
    setFilterBy(filterBy);
    setComment(comment);

    onAdd();
  };

  const handleGetComment = async (): Promise<string | undefined> => {
    return prompt({
      title: 'Add comments',
      input: {
        defaultValue: '',
        type: 'textarea',
      },
    });
  };

  const handleConfirmByPurl = async (): Promise<boolean> => {
    return ask(
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
                  onClick={() => onSelectOption(action, 'by_file', true)}
                >
                  <span className="first-letter:uppercase">{`${action} with Comments`}</span>
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => onSelectOption(action, 'by_file', false)}
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
                  onClick={() => onSelectOption(action, 'by_purl', true)}
                >
                  <span className="first-letter:uppercase">{`${action} with Comments`}</span>
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => onSelectOption(action, 'by_purl', false)}
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
