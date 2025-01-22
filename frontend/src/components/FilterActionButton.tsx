// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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
import useEnvironment from '@/hooks/useEnvironment';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import useSelectedResult from '@/hooks/useSelectedResult';
import { getShortcutDisplay } from '@/lib/shortcuts';
import { FilterAction, filterActionLabelMap } from '@/modules/components/domain';
import { OnFilterComponentArgs } from '@/modules/components/stores/useComponentFilterStore';
import { useDialog } from '@/providers/DialogProvider';

import AddCommentTextarea from './AddCommentTextarea';
import FilterByPurlList from './FilterByPurlList';
import SelectLicenseList from './SelectLicenseList';

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
  const { showDialog } = useDialog();
  const selectedResult = useSelectedResult();
  const isCompletedResult = selectedResult?.workflow_state === 'completed';
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const { modifierKey } = useEnvironment();

  const isReplaceAction = action === FilterAction.Replace;

  const handleFilterByFileWithComments = async () => {
    let comment: string | undefined;

    if (!isReplaceAction) {
      const confirm = await showDialog({
        title: 'Add comments',
        content: <AddCommentTextarea onChange={(value) => (comment = value)} />,
      });

      if (!confirm) return;
    }

    return onAdd({
      action,
      comment,
      withComment: true,
      filterBy: 'by_file',
    });
  };

  const handleFilterByPurlWithComments = async () => {
    let comment: string | undefined;

    if (!isReplaceAction) {
      const confirm = await showDialog({
        title: 'Add comments',
        content: (
          <div className="flex flex-col gap-4">
            <AddCommentTextarea onChange={(value) => (comment = value)} />
            <FilterByPurlList action={action} />
          </div>
        ),
      });

      if (!confirm) return;
    }

    return onAdd({
      action,
      comment,
      withComment: true,
      filterBy: 'by_purl',
    });
  };

  const handleFilterByFileWithoutComments = async () => {
    return onAdd({
      action,
      filterBy: 'by_file',
    });
  };

  const handleFilterByPurlWithoutComments = async () => {
    if (!isReplaceAction) {
      const confirmed = await showDialog({
        title: 'Are you sure?',
        content: <FilterByPurlList action={action} />,
      });

      if (!confirmed) return;
    }

    return onAdd({
      action,
      filterBy: 'by_purl',
    });
  };

  const handleFilterByFileAndDifferentLicense = async () => {
    let license: string | undefined;
    let comment: string | undefined;

    if (!isReplaceAction) {
      const confirmed = await showDialog({
        title: 'Select License',
        description: 'Select the license you want to apply to the selected files. You can also add a comment.',
        content: (
          <div className="flex flex-col gap-4">
            <SelectLicenseList onSelect={(selectedLicense) => (license = selectedLicense)} />
            <AddCommentTextarea onChange={(value) => (comment = value)} />
          </div>
        ),
      });

      if (!confirmed) return;
    }

    return onAdd({
      action,
      comment,
      license,
      withComment: !!comment,
      filterBy: 'by_file',
    });
  };

  const handleFilterByPurlAndDifferentLicense = async () => {
    let license: string | undefined;
    let comment: string | undefined;

    if (!isReplaceAction) {
      const confirmed = await showDialog({
        title: 'Select License',
        description: 'Select the license you want to apply to the selected files. You can also add a comment.',
        content: (
          <div className="flex flex-col gap-4">
            <SelectLicenseList onSelect={(selectedLicense) => (license = selectedLicense)} />
            <AddCommentTextarea onChange={(value) => (comment = value)} />
            <FilterByPurlList action={action} />
          </div>
        ),
      });

      if (!confirmed) return;
    }

    return onAdd({
      action,
      comment,
      license,
      withComment: !!comment,
      filterBy: 'by_purl',
    });
  };

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
          <button className="flex h-full w-4 items-center outline-none transition-colors enabled:hover:bg-accent" disabled={isCompletedResult}>
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
              <DropdownMenuSubContent className="min-w-[420px]">
                <DropdownMenuItem onClick={handleFilterByFileWithoutComments}>
                  <span className="first-letter:uppercase">{`${action} without comments`}</span>
                  <DropdownMenuShortcut>{getShortcutDisplay(shortcutKeysByFileWithoutComments, modifierKey.label)[0]}</DropdownMenuShortcut>
                </DropdownMenuItem>
                <DropdownMenuItem onClick={handleFilterByFileWithComments}>
                  <span className="first-letter:uppercase">{`${action} with comments`}</span>
                  <DropdownMenuShortcut>{getShortcutDisplay(shortcutKeysByFileWithComments, modifierKey.label)[0]}</DropdownMenuShortcut>
                </DropdownMenuItem>
                {!isReplaceAction && (
                  <DropdownMenuItem onClick={handleFilterByFileAndDifferentLicense}>
                    <span className="first-letter:uppercase">{`${action} with a different license`}</span>
                  </DropdownMenuItem>
                )}
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>
        </DropdownMenuGroup>
        <DropdownMenuGroup>
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>Component</DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent className="min-w-[420px]">
                <DropdownMenuItem onClick={handleFilterByPurlWithoutComments}>
                  <span className="first-letter:uppercase">{`${action} without comments`}</span>
                  <DropdownMenuShortcut>{getShortcutDisplay(shortcutKeysByComponentWithoutComments, modifierKey.label)[0]}</DropdownMenuShortcut>
                </DropdownMenuItem>
                <DropdownMenuItem onClick={handleFilterByPurlWithComments}>
                  <span className="first-letter:uppercase">{`${action} with Comments`}</span>
                  <DropdownMenuShortcut>{getShortcutDisplay(shortcutKeysByComponentWithComments, modifierKey.label)[0]}</DropdownMenuShortcut>
                </DropdownMenuItem>
                {!isReplaceAction && (
                  <DropdownMenuItem onClick={handleFilterByPurlAndDifferentLicense}>
                    <span className="first-letter:uppercase">{`${action} with a different license`}</span>
                  </DropdownMenuItem>
                )}
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
