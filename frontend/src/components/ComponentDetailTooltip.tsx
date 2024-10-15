import clsx from 'clsx';
import { entities } from 'wailsjs/go/models';

import Link from '@/components/Link';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { MatchType, matchTypePresentation } from '@/modules/results/domain';

interface ComponentDetailTooltipProps {
  component: entities.ComponentDTO;
}

export default function ComponentDetailTooltip({
  component,
}: ComponentDetailTooltipProps) {
  const matchPresentation = matchTypePresentation[component.id as MatchType];

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <div className="cursor-pointer">
          <div
            className={clsx(
              'text-lg font-bold leading-tight',
              matchPresentation.accent
            )}
          >
            {component.component}
          </div>
          <div>{component.purl?.[0]}</div>
        </div>
      </TooltipTrigger>
      <TooltipContent side="bottom" align="start" className="p-4">
        <div className="flex flex-col gap-4">
          {component.purl?.length ? (
            <div>
              <p className="font-medium">PURL</p>
              <p className="text-muted-foreground">{component.purl?.[0]}</p>
            </div>
          ) : null}
          <div>
            <p className="font-medium">VERSION</p>
            <p className="text-muted-foreground">{component.version}</p>
          </div>
          {component.licenses?.length ? (
            <div>
              <p className="font-medium">LICENSE</p>
              <p className="text-muted-foreground">
                {component.licenses?.[0].name}
              </p>
            </div>
          ) : null}
          {component.url && (
            <div>
              <p className="font-medium">URL</p>
              <Link to={component.url} />
            </div>
          )}
        </div>
      </TooltipContent>
    </Tooltip>
  );
}
