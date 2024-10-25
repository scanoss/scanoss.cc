import clsx from 'clsx';
import { ChevronDown } from 'lucide-react';
import { useCallback, useState } from 'react';

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
import { FilterAction, filterActionLabelMap } from '@/modules/components/domain';
import useComponentFilterStore from '@/modules/components/stores/useComponentFilterStore';

import FilterByPurlList from './FilterByPurlList';

interface FilterActionProps {
  action: FilterAction;
  description: string;
  icon: React.ReactNode;
  onAdd: () => Promise<void> | void;
  shortcutKeysByFile: string;
  shortcutKeysByPurl: string;
}

export default function FilterActionButton({
  action,
  description,
  icon,
  onAdd,
  shortcutKeysByFile,
  shortcutKeysByPurl,
}: FilterActionProps) {
  const [dropdownOpen, setDropdownOpen] = useState(false);

  const { ask } = useConfirm();
  const { prompt } = useInputPrompt();

  const setAction = useComponentFilterStore((state) => state.setAction);
  const setComment = useComponentFilterStore((state) => state.setComment);
  const setFilterBy = useComponentFilterStore((state) => state.setFilterBy);
  const setWithComment = useComponentFilterStore((state) => state.setWithComment);

  const actionsThatShouldPromptForCommentOrConfirmation = [FilterAction.Include, FilterAction.Remove];

  const maybePromptForCommentOrConfirmation = async (
    filterBy: 'by_file' | 'by_purl',
    withComment: boolean
  ): Promise<boolean> => {
    let comment: string | undefined;

    if (withComment) {
      comment = await handleGetComment();

      if (!comment) return false;
    }

    if (filterBy === 'by_purl') {
      const confirm = await handleConfirmByPurl();

      if (!confirm) return false;
    }

    setComment(comment);

    return true;
  };

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

  const onSelectOption = async (action: FilterAction, filterBy: 'by_file' | 'by_purl', withComment: boolean) => {
    setAction(action);
    setFilterBy(filterBy);
    setWithComment(withComment);

    if (actionsThatShouldPromptForCommentOrConfirmation.includes(action)) {
      const shouldContinue = await maybePromptForCommentOrConfirmation(filterBy, withComment);

      if (!shouldContinue) return;
    }

    onAdd();
  };

  const handleFilterByFileWithComments = () => onSelectOption(action, 'by_file', true);
  const handleFilterByPurlWithComments = () => onSelectOption(action, 'by_purl', true);
  const handleFilterByFileWithoutComments = () => onSelectOption(action, 'by_file', false);
  const handleFilterByPurlWithoutComments = () => onSelectOption(action, 'by_purl', false);

  const keysByFileWithComments = `alt+${shortcutKeysByFile}`;
  const keysByPurlWithComments = `alt+${shortcutKeysByPurl}`;

  useKeyboardShortcut(shortcutKeysByFile, handleFilterByFileWithoutComments);
  useKeyboardShortcut(keysByFileWithComments, handleFilterByFileWithComments);
  useKeyboardShortcut(shortcutKeysByPurl, handleFilterByPurlWithoutComments);
  useKeyboardShortcut(keysByPurlWithComments, handleFilterByPurlWithComments);

  const getShortcutLabel = useCallback((keys: string) => {
    const parts = keys.split('+');
    return parts
      .map((p) => {
        if (p === 'shift') return '⇧';
        if (p === 'mod') return '⌘';
        return p;
      })
      .join('  ');
  }, []);

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
                  <DropdownMenuShortcut className="uppercase">
                    {getShortcutLabel(shortcutKeysByFile)}
                  </DropdownMenuShortcut>
                </DropdownMenuItem>
                <DropdownMenuItem onClick={handleFilterByFileWithComments}>
                  <span className="first-letter:uppercase">{`${action} with Comments`}</span>
                  <DropdownMenuShortcut className="uppercase">
                    {getShortcutLabel(keysByFileWithComments)}
                  </DropdownMenuShortcut>
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
                  <DropdownMenuShortcut className="uppercase">
                    {getShortcutLabel(shortcutKeysByPurl)}
                  </DropdownMenuShortcut>
                </DropdownMenuItem>
                <DropdownMenuItem onClick={handleFilterByPurlWithComments}>
                  <span className="first-letter:uppercase">{`${action} with Comments`}</span>
                  <DropdownMenuShortcut className="uppercase">
                    {getShortcutLabel(keysByPurlWithComments)}
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
