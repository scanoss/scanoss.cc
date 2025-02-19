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
import { MoveVertical } from 'lucide-react';
import { useEffect, useState } from 'react';

import useEnvironment from '@/hooks/useEnvironment';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { MonacoManager } from '@/lib/editor';
import { getShortcutDisplay, KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';

import { entities } from '../../wailsjs/go/models';
import { EventsOff, EventsOn } from '../../wailsjs/runtime/runtime';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';

export default function EditorToolbar() {
  const { modifierKey } = useEnvironment();
  const monacoManager = MonacoManager.getInstance();
  const [scrollEnabled, setScrollEnabled] = useState(monacoManager.getScrollSyncEnabled());

  const handleToggleSyncScroll = () => {
    monacoManager.toggleSyncScroll();
    setScrollEnabled(monacoManager.getScrollSyncEnabled());
  };

  useKeyboardShortcut(KEYBOARD_SHORTCUTS.toggleSyncScrollPosition.keys, handleToggleSyncScroll);

  useEffect(() => {
    EventsOn(entities.Action.ToggleSyncScrollPosition, handleToggleSyncScroll);

    return () => {
      EventsOff(entities.Action.ToggleSyncScrollPosition);
    };
  }, []);

  return (
    <div className="flex p-2">
      <Tooltip>
        <TooltipTrigger asChild>
          <div
            className={clsx(
              'flex cursor-pointer items-center border px-2 py-1 text-xs text-muted-foreground hover:text-white',
              scrollEnabled && 'bg-primary text-white'
            )}
            onClick={handleToggleSyncScroll}
          >
            <MoveVertical className="mr-1 inline-block h-3 w-3" />
            Sync Scroll Position
          </div>
        </TooltipTrigger>
        <TooltipContent>({getShortcutDisplay(KEYBOARD_SHORTCUTS.toggleSyncScrollPosition.keys, modifierKey.label)})</TooltipContent>
      </Tooltip>
    </div>
  );
}
