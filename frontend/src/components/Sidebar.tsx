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
    <div className="flex flex-col border-r bg-background h-full gap-4">
      <h2 className="text-sm font-semibold px-4 pt-6">
        {files?.length
          ? `${files.length} change${files.length > 1 ? 's' : ''} in working directory`
          : 'You have no changes in working directory'}
      </h2>
      <hr />
      <div className="px-4 flex flex-col gap-1">
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
            <SelectItem value={MatchType.None}>No match</SelectItem>
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
                    'flex w-full items-center gap-2 text-sm text-muted-foreground px-4 py-1 transition-all',
                    isActive
                      ? 'hover:bg-primary hover:text-primary-foreground bg-primary text-primary-foreground'
                      : 'hover:bg-primary/10'
                  )}
                  to={`files/${encodedFilePath}?matchType=${file.matchType}`}
                >
                  <div className="w-1.5 h-1.5 rounded-full bg-red-500"></div>
                  <span className="overflow-hidden text-ellipsis max-w-[80%]">
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
    </div>
  );
}
