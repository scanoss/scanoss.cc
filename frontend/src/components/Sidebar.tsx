import { useQuery } from '@tanstack/react-query';
import clsx from 'clsx';
import React, { useState } from 'react';
import { Link, useParams } from 'react-router-dom';

import { decodeFilePath, encodeFilePath } from '@/lib/utils';
import { MatchType } from '@/modules/results/domain';
import ResultService from '@/modules/results/infra/service';

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from './ui/select';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';

export default function Sidebar() {
  const [filterByMatchType, setFilterByMatchType] = useState<MatchType | 'all'>(
    'all'
  );
  const { filePath } = useParams();

  const { data: files } = useQuery({
    queryKey: ['results', filterByMatchType],
    queryFn: () =>
      ResultService.getAll(
        filterByMatchType === 'all' ? undefined : filterByMatchType
      ),
  });

  return (
    <aside className="z-10 flex h-full flex-col border-r border-border bg-black/20 backdrop-blur-sm">
      <div className="flex h-[65px] items-center border-b border-b-border px-4">
        <h2 className="text-sm font-semibold">
          {files?.length
            ? `${files.length} change${files.length > 1 ? 's' : ''} in working directory`
            : 'You have no changes in working directory'}
        </h2>
      </div>

      <div className="flex flex-col gap-1 px-4 py-6">
        <span className="text-xs font-semibold">Filter by match type</span>
        <Select
          onValueChange={(value) => setFilterByMatchType(value as MatchType)}
          defaultValue="all"
        >
          <SelectTrigger className="w-full">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value={MatchType.File}>File</SelectItem>
            <SelectItem value={MatchType.Snippet}>Snippet</SelectItem>
            <SelectItem value="all">All</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="flex flex-col gap-1.5">
        {files?.map((file) => {
          const isActive = decodeFilePath(filePath ?? '') === file.path;
          const encodedFilePath = encodeFilePath(file.path);

          return (
            <Tooltip key={file.path}>
              <TooltipTrigger asChild>
                <Link
                  className={clsx(
                    'flex w-full items-center gap-2 px-4 py-1 text-sm text-muted-foreground transition-all',
                    isActive
                      ? 'border-r-2 border-primary-foreground bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground'
                      : 'hover:bg-primary/10'
                  )}
                  to={`files/${encodedFilePath}?matchType=${file.matchType}`}
                >
                  <div
                    className={clsx(
                      'h-1.5 w-1.5 rounded-full',
                      matchTypeColors[file.matchType]
                    )}
                  ></div>
                  <span className="max-w-[80%] overflow-hidden text-ellipsis">
                    {file.path}
                  </span>
                </Link>
              </TooltipTrigger>
              <TooltipContent side="right" sideOffset={15}>
                {file.path}
              </TooltipContent>
            </Tooltip>
          );
        })}
      </div>
    </aside>
  );
}

const matchTypeColors: Record<MatchType, string> = {
  [MatchType.File]: 'bg-cyan-500',
  [MatchType.Snippet]: 'bg-yellow-500',
};
