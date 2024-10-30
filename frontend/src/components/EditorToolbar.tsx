import clsx from 'clsx';
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
        className={clsx('cursor-pointer border px-2 py-1 text-xs hover:bg-primary', scrollEnabled && 'bg-primary')}
        onClick={handleToggleSyncScroll}
      >
        Sync Cursor
      </div>
    </div>
  );
}
