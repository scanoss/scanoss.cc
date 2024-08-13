import React from 'react';

import { Skeleton } from '@/components/ui/skeleton';

export default function MatchSkeleton() {
  return (
    <div className="flex h-full w-full flex-col gap-2">
      <Skeleton className="h-[50px] w-full" />
      <div className="flex flex-1 gap-8">
        <div className="flex flex-1 flex-col gap-2">
          <Skeleton className="h-[50px] w-full" />
          <Skeleton className="h-full w-full" />
        </div>
        <div className="flex flex-1 flex-col gap-2">
          <Skeleton className="h-[50px] w-full" />
          <Skeleton className="h-full w-full" />
        </div>
      </div>
    </div>
  );
}
