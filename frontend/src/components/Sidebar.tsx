import { Link, useParams } from '@tanstack/react-router';
import clsx from 'clsx';
import React from 'react';

import { GitFile } from '@/modules/files/domain';

import { ScrollArea } from './ui/scroll-area';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';

interface SidebarProps {
  files: GitFile[];
}

export default function Sidebar({ files }: SidebarProps) {
  const { filePath } = useParams({ strict: false });

  return (
    <div className="flex flex-col border-r bg-muted/40 h-full p-6">
      <h2 className="text-sm font-semibold">
        {files.length} changes in working directory
      </h2>
      <div className="flex flex-col gap-2 mt-3">
        {files.map((file) => {
          const isActive = filePath === file.path;

          return (
            <Tooltip key={file.path}>
              <TooltipTrigger asChild>
                <Link
                  className={clsx(
                    'text-sm text-muted-foreground hover:bg-primary/10 px-2 py-1 rounded-md transition-all overflow-hidden text-ellipsis',
                    isActive &&
                      'hover:bg-primary hover:text-primary-foreground bg-primary text-primary-foreground'
                  )}
                  to="/files/$filePath"
                  params={{
                    filePath: file.path,
                  }}
                >
                  {file.path}
                </Link>
              </TooltipTrigger>
              <TooltipContent side="right" sideOffset={15}>
                {file.path}
              </TooltipContent>
            </Tooltip>
          );
        })}
      </div>
    </div>
  );
}
