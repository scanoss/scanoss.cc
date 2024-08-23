import { Check, FileWarning, PackageMinus, Replace } from 'lucide-react';

import { Component } from '@/modules/results/domain';

import { FilterAction } from '../domain';
import FileActionButton from './FileActionButton';

interface FileActionsMenuProps {
  component: Component;
}

export default function FileActionsMenu({ component }: FileActionsMenuProps) {
  return (
    <div className="flex h-[65px] justify-center border-b border-b-border px-4">
      <div className="flex gap-2">
        <FileActionButton
          action={FilterAction.Include}
          component={component}
          icon={<Check className="h-5 w-5 stroke-green-500" />}
          description="By including a file/component, you force the engine to consider it with priority in future scans."
        />
        <FileActionButton
          action={FilterAction.Remove}
          component={component}
          description="Removing a file/component will exclude it from future scans."
          icon={<PackageMinus className="h-5 w-5 stroke-red-500" />}
        />
        <FileActionButton
          action={FilterAction.Ignore}
          component={component}
          icon={<FileWarning className="h-5 w-5 stroke-amber-500" />}
          description="By ignoring a file/component, you prevent the engine from considering it in future scans."
          isDisabled
        />
        <FileActionButton
          action={FilterAction.Replace}
          component={component}
          icon={<Replace className="h-5 w-5 stroke-cyan-500" />}
          description="Replace a file/component with another one."
          isDisabled
        />
      </div>
    </div>
  );
}
