// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import clsx from 'clsx';
import { MoveVertical } from 'lucide-react';
import { useState } from 'react';

import useEnvironment from '@/hooks/useEnvironment';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { MonacoManager } from '@/lib/editor';
import { getShortcutDisplay, KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';

import { EventsEmit } from '../../wailsjs/runtime/runtime';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';

export default function EditorToolbar() {
  const { modifierKey } = useEnvironment();
  const monacoManager = MonacoManager.getInstance();
  const [scrollEnabled, setScrollEnabled] = useState(monacoManager.getScrollSyncEnabled());

  const handleToggleSyncScroll = () => {
    monacoManager.toggleSyncScroll();
    setScrollEnabled(monacoManager.getScrollSyncEnabled());
    EventsEmit('editor:syncScrollPosition', monacoManager.getScrollSyncEnabled());
  };

  useKeyboardShortcut(KEYBOARD_SHORTCUTS.toggleSyncScrollPosition.keys, handleToggleSyncScroll);

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
