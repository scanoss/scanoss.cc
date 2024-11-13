import clsx from 'clsx';
import { FolderTree, Menu } from 'lucide-react';

import useResultsStore from '@/modules/results/stores/useResultsStore';

import { Button } from './ui/button';

export default function ViewTypeSelector() {
  const viewType = useResultsStore((state) => state.viewType);
  const setViewType = useResultsStore((state) => state.setViewType);

  const activeStyles = 'bg-accent text-accent-foreground';

  return (
    <div className="flex w-full items-center justify-center">
      <div className="border">
        <Button
          variant="ghost"
          size="sm"
          className={clsx('h-fit rounded-none px-2 py-1', {
            [activeStyles]: viewType === 'flat',
          })}
          onClick={() => setViewType('flat')}
        >
          <Menu className="mr-1 h-3 w-3" />
          Flat
        </Button>
        <Button
          variant="ghost"
          size="sm"
          className={clsx('h-fit rounded-none px-2 py-1', {
            [activeStyles]: viewType === 'tree',
          })}
          onClick={() => setViewType('tree')}
        >
          <FolderTree className="mr-1 h-3 w-3" />
          Tree
        </Button>
      </div>
    </div>
  );
}
