import clsx from 'clsx';
import { MousePointer, MoveVertical } from 'lucide-react';
import { useState } from 'react';

import { MonacoManager } from '@/lib/editor';

export default function EditorToolbar() {
  const monacoManager = MonacoManager.getInstance();
  const [scrollEnabled, setScrollEnabled] = useState(monacoManager.getScrollEnabled());

  const handleToggleSyncScroll = () => {
    monacoManager.toggleSyncScroll();
    setScrollEnabled(monacoManager.getScrollEnabled());
  };

  return (
    <div className="flex p-2">
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
    </div>
  );
}
