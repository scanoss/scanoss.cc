import { Ban, Check } from 'lucide-react';
import React from 'react';

import { Button } from '@/components/ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';

export default function MatchInfoCard() {
  // TODO: Get match info from backend

  return (
    <div className="flex cursor-pointer items-center justify-between rounded-sm border border-primary-foreground bg-primary p-3">
      <div className="flex items-center gap-8">
        <div>
          <div className="font-bold text-primary-foreground">
            scancode-toolkit
          </div>
          <div className="text-sm">pkg:github/nexb/scancode-toolkit</div>
        </div>
        <div>
          <div className="text-sm text-muted-foreground">Version</div>
          <div className="text-sm">32.1.0</div>
        </div>
        <div>
          <div className="text-sm text-muted-foreground">License</div>
          <div className="text-sm">MIT</div>
        </div>
        <div>
          <div className="text-sm text-muted-foreground">Detected</div>
          <div className="text-sm text-green-500">File</div>
        </div>
        <div>
          <div className="text-sm text-muted-foreground">Match</div>
          <div className="text-sm">100%</div>
        </div>
      </div>
      <div className="flex gap-2">
        <Tooltip>
          <TooltipTrigger asChild>
            <Button size="icon" variant="ghost">
              <Check className="h-6 w-6 stroke-green-500" />
            </Button>
          </TooltipTrigger>
          <TooltipContent side="bottom">Identify</TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button size="icon" variant="ghost">
              <Ban className="h-5 w-5 stroke-muted-foreground" />
            </Button>
          </TooltipTrigger>
          <TooltipContent side="bottom">Mark as original</TooltipContent>
        </Tooltip>
      </div>
    </div>
  );
}
