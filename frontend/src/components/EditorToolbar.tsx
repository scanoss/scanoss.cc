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
