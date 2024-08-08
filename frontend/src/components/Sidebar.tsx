import { useQuery } from '@tanstack/react-query';
import clsx from 'clsx';
import React from 'react';
import { Link, useParams } from 'react-router-dom';

import { decodeFilePath, encodeFilePath } from '@/lib/utils';
import FileService from '@/modules/files/infra/service';

import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';

export default function Sidebar() {
  const { filePath } = useParams();

  // TODO: Add loading and error states
  const { data: files } = useQuery({
    queryKey: ['filesToBeCommited'],
    queryFn: () => FileService.getAllToBeCommited(),
  });

  return (
    <div className="flex flex-col border-r bg-background h-full">
      <h2 className="text-sm font-semibold px-4 py-6">
        {files?.length
          ? `${files.length} changes in working directory`
          : 'You have no changes in working directory'}
      </h2>
      <div className="flex flex-col gap-1.5">
        {files?.map((file) => {
          const isActive = decodeFilePath(filePath ?? '') === file.path;
          const encodedFilePath = encodeFilePath(file.path);

          return (
            <Tooltip key={file.path}>
              <TooltipTrigger asChild>
                <Link
                  className={clsx(
                    'text-sm text-muted-foreground px-4 py-1 transition-all overflow-hidden text-ellipsis',
                    isActive
                      ? 'hover:bg-primary hover:text-primary-foreground bg-primary text-primary-foreground'
                      : 'hover:bg-primary/10'
                  )}
                  to={`files/${encodedFilePath}`}
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
