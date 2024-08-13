import { Ban, Check } from 'lucide-react';
import React from 'react';

import { Button } from '@/components/ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { Component } from '@/modules/results/domain';

interface MatchInfoCardProps {
  component: Component;
}

export default function MatchInfoCard({ component }: MatchInfoCardProps) {
  return (
    <div className="bg-primary rounded-sm p-3 cursor-pointer border border-primary-foreground flex justify-between items-center">
      <div className="flex gap-8 items-center">
        <div>
          <div className="text-primary-foreground font-bold">
            {component.component}
          </div>
          <div className="text-sm">{component.purl?.[0]}</div>
        </div>
        {component.version && (
          <div>
            <div className="text-muted-foreground text-sm">Version</div>
            <div className="text-sm">{component.version}</div>
          </div>
        )}
        {component.licenses?.length ? (
          <div>
            <div className="text-muted-foreground text-sm">License</div>
            <div className="text-sm">{component.licenses?.[0].name}</div>
          </div>
        ) : null}
        <div>
          <div className="text-muted-foreground text-sm">Detected</div>
          <div className="text-sm text-green-500 first-letter:uppercase">
            {component.id}
          </div>
        </div>
        <div>
          <div className="text-muted-foreground text-sm">Match</div>
          <div className="text-sm">{component.matched}</div>
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
