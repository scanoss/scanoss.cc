import React from 'react';

import { Skeleton } from '@/components/ui/skeleton';

export default function MatchSkeleton() {
  return (
    <div className="w-full h-full flex flex-col gap-2">
      <Skeleton className="w-full h-[50px]" />
      <div className="flex flex-1 gap-8">
        <div className="flex flex-col flex-1 gap-2">
          <Skeleton className="w-full h-[50px]" />
          <Skeleton className="w-full h-full" />
        </div>
        <div className="flex flex-col flex-1 gap-2">
          <Skeleton className="w-full h-[50px]" />
          <Skeleton className="w-full h-full" />
        </div>
      </div>
    </div>
  );
}
