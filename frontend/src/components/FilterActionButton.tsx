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
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useConfirm } from '@/hooks/useConfirm';
import { useInputPrompt } from '@/hooks/useInputPrompt';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import useSelectedResult from '@/hooks/useSelectedResult';
import { getShortcutDisplay } from '@/lib/shortcuts';
import { FilterAction, filterActionLabelMap } from '@/modules/components/domain';
import { OnFilterComponentArgs } from '@/modules/components/stores/useComponentFilterStore';

import FilterByPurlList from './FilterByPurlList';

interface FilterActionProps {
  action: FilterAction;
  description: string;
  icon: React.ReactNode;
  onAdd: (args: OnFilterComponentArgs) => Promise<void> | void;
  shortcutKeysByFileWithComments: string;
  shortcutKeysByFileWithoutComments: string;
  shortcutKeysByComponentWithComments: string;
  shortcutKeysByComponentWithoutComments: string;
}

export default function FilterActionButton({
  action,
  description,
  icon,
  onAdd,
  shortcutKeysByFileWithComments,
  shortcutKeysByFileWithoutComments,
  shortcutKeysByComponentWithComments,
  shortcutKeysByComponentWithoutComments,
}: FilterActionProps) {
  const selectedResult = useSelectedResult();
  const isCompletedResult = selectedResult?.workflow_state === 'completed';
  const [dropdownOpen, setDropdownOpen] = useState(false);

  const { ask } = useConfirm();
  const { prompt } = useInputPrompt();

  const actionsThatShouldPromptForCommentOrConfirmation = [FilterAction.Include, FilterAction.Remove];

  const handleConfirmByPurl = async (): Promise<boolean> => ask(<FilterByPurlList action={action} />);

  const handleGetComment = async (): Promise<string | undefined> => {
    return prompt({
      title: 'Add comments',
      input: {
        defaultValue: '',
        type: 'textarea',
      },
    });
  };

  const onSelectOption = async (filterBy: 'by_file' | 'by_purl', withComment: boolean) => {
    if (isCompletedResult) return;

    let comment: string | undefined;
    if (actionsThatShouldPromptForCommentOrConfirmation.includes(action)) {
      if (withComment) {
        comment = await handleGetComment();

        if (!comment) return;
      }

      if (filterBy === 'by_purl') {
        const confirm = await handleConfirmByPurl();

        if (!confirm) return;
      }
    }

    onAdd({
      action,
      filterBy,
      comment,
      withComment,
    });
  };

  const handleFilterByFileWithComments = () => onSelectOption('by_file', true);
  const handleFilterByPurlWithComments = () => onSelectOption('by_purl', true);
  const handleFilterByFileWithoutComments = () => onSelectOption('by_file', false);
  const handleFilterByPurlWithoutComments = () => onSelectOption('by_purl', false);

  useKeyboardShortcut(shortcutKeysByFileWithoutComments, handleFilterByFileWithoutComments, {
    enabled: !isCompletedResult,
  });
  useKeyboardShortcut(shortcutKeysByComponentWithoutComments, handleFilterByPurlWithoutComments, {
    enabled: !isCompletedResult,
  });
  useKeyboardShortcut(shortcutKeysByFileWithComments, handleFilterByFileWithComments, {
    enabled: !isCompletedResult,
  });
  useKeyboardShortcut(shortcutKeysByComponentWithComments, handleFilterByPurlWithComments, {
    enabled: !isCompletedResult,
  });

  return (
    <DropdownMenu onOpenChange={setDropdownOpen}>
      <div className="group flex h-full">
        <DropdownMenuTrigger asChild disabled={isCompletedResult}>
          <Button
            variant="ghost"
            size="lg"
            className="h-full w-14 rounded-none enabled:group-hover:bg-accent enabled:group-hover:text-accent-foreground"
            disabled={isCompletedResult}
          >
            <div className="flex flex-col items-center justify-center gap-1">
              <span className="text-xs">{filterActionLabelMap[action]}</span>
              {icon}
            </div>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuTrigger asChild disabled={isCompletedResult}>
          <button
            className="flex h-full w-4 items-center outline-none transition-colors enabled:hover:bg-accent"
            disabled={isCompletedResult}
          >
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
          <span className="text-xs font-normal text-muted-foreground">{description}</span>
        </DropdownMenuLabel>
        <DropdownMenuSeparator className="bg-border" />
        <DropdownMenuGroup>
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>File</DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent className="min-w-[300px]">
                <DropdownMenuItem onClick={handleFilterByFileWithoutComments}>
                  <span className="first-letter:uppercase">{`${action} without Comments`}</span>
                  <DropdownMenuShortcut>
                    {getShortcutDisplay(shortcutKeysByFileWithoutComments)[0]}
                  </DropdownMenuShortcut>
                </DropdownMenuItem>
                <DropdownMenuItem onClick={handleFilterByFileWithComments}>
                  <span className="first-letter:uppercase">{`${action} with Comments`}</span>
                  <DropdownMenuShortcut>{getShortcutDisplay(shortcutKeysByFileWithComments)[0]}</DropdownMenuShortcut>
                </DropdownMenuItem>
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>
        </DropdownMenuGroup>
        <DropdownMenuGroup>
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>Component</DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent className="min-w-[300px]">
                <DropdownMenuItem onClick={handleFilterByPurlWithoutComments}>
                  <span className="first-letter:uppercase">{`${action} without Comments`}</span>
                  <DropdownMenuShortcut>
                    {getShortcutDisplay(shortcutKeysByComponentWithoutComments)[0]}
                  </DropdownMenuShortcut>
                </DropdownMenuItem>
                <DropdownMenuItem onClick={handleFilterByPurlWithComments}>
                  <span className="first-letter:uppercase">{`${action} with Comments`}</span>
                  <DropdownMenuShortcut>
                    {getShortcutDisplay(shortcutKeysByComponentWithComments)[0]}
                  </DropdownMenuShortcut>
                </DropdownMenuItem>
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
