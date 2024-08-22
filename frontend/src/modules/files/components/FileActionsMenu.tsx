import { Ban, Check } from 'lucide-react';

import { Component } from '@/modules/results/domain';

import { FilterAction } from '../domain';
import FileActionButton from './FileActionButton';

interface FileActionsMenuProps {
  component: Component;
}

export default function FileActionsMenu({ component }: FileActionsMenuProps) {
  return (
    <div className="flex h-[65px] justify-center border-b border-b-border px-4">
      <div className="flex gap-1">
        <FileActionButton
          action={FilterAction.Include}
          component={component}
          icon={<Check className="h-5 w-5 stroke-green-500" />}
        />
        <FileActionButton
          action={FilterAction.Remove}
          component={component}
          icon={<Ban className="h-5 w-4 stroke-red-500" />}
        />
      </div>
    </div>
  );
}
