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
    <div className="bg-background rounded-md p-3 cursor-pointer border border-blue-400 flex justify-between items-center">
      <div className="flex gap-8 items-center">
        <div>
          <div className="text-indigo-500 font-bold">scancode-toolkit</div>
          <div className="text-sm">pkg:github/nexb/scancode-toolkit</div>
        </div>
        <div>
          <div className="text-muted-foreground text-sm">Version</div>
          <div className="text-sm">32.1.0</div>
        </div>
        <div>
          <div className="text-muted-foreground text-sm">License</div>
          <div className="text-sm">MIT</div>
        </div>
        <div>
          <div className="text-muted-foreground text-sm">Detected</div>
          <div className="text-sm text-green-500">File</div>
        </div>
        <div>
          <div className="text-muted-foreground text-sm">Match</div>
          <div className="text-sm">100%</div>
        </div>
      </div>
      <div className="flex gap-2">
        <Tooltip>
          <TooltipTrigger asChild>
            <Button size="icon" variant="ghost">
              <Check className="w-6 h-6 stroke-green-500" />
            </Button>
          </TooltipTrigger>
          <TooltipContent side="bottom">Identify</TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button size="icon" variant="ghost">
              <Ban className="w-5 h-5 stroke-muted-foreground" />
            </Button>
          </TooltipTrigger>
          <TooltipContent side="bottom">Mark as original</TooltipContent>
        </Tooltip>
      </div>
    </div>
  );
}
