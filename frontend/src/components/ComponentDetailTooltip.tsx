import clsx from 'clsx';
import { ArrowRight } from 'lucide-react';
import { entities } from 'wailsjs/go/models';

import Link from '@/components/Link';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import useSelectedResult from '@/hooks/useSelectedResult';
import { FilterAction } from '@/modules/components/domain';
import { MatchType, matchTypePresentation } from '@/modules/results/domain';

interface ComponentDetailTooltipProps {
  component: entities.ComponentDTO;
}

export default function ComponentDetailTooltip({ component }: ComponentDetailTooltipProps) {
  const result = useSelectedResult();

  const isResultReplaced = result?.filter_config?.action === FilterAction.Replace;

  if (isResultReplaced) {
    return (
      <div className="flex items-center gap-4">
        <DetectedPurlTooltip component={component} replaced />
        <ArrowRight className="h-4 w-4 text-primary-foreground" />
        <ConcludedPurlTooltip component={component} />
      </div>
    );
  }

  return <DetectedPurlTooltip component={component} />;
}

function DetectedPurlTooltip({ component, replaced }: { component: entities.ComponentDTO; replaced?: boolean }) {
  const result = useSelectedResult();
  const matchPresentation = matchTypePresentation[component.id as MatchType];

  const isResultRemoved = result?.filter_config?.action === FilterAction.Remove;

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <div
          className={clsx('cursor-pointer', {
            'text-muted-foreground line-through opacity-70': replaced,
          })}
        >
          <div
            className={clsx('text-lg font-bold leading-tight', {
              [matchPresentation.accent]: !replaced && !isResultRemoved,
            })}
          >
            {component.component}
          </div>
          <div>{result?.detected_purl}</div>
        </div>
      </TooltipTrigger>
      <TooltipContent side="bottom" align="start" className="p-4">
        <div className="flex flex-col gap-4">
          <div>
            <p className="font-medium">PURL</p>
            <p className="text-muted-foreground">{result?.detected_purl}</p>
          </div>
          {component.version && (
            <div>
              <p className="font-medium">VERSION</p>
              <p className="text-muted-foreground">{component.version}</p>
            </div>
          )}
          {component.licenses?.length ? (
            <div>
              <p className="font-medium">LICENSE</p>
              <p className="text-muted-foreground">{component.licenses?.[0].name}</p>
            </div>
          ) : null}
          {component.url && (
            <div>
              <p className="font-medium">URL</p>
              <Link to={component.url as string} />
            </div>
          )}
        </div>
      </TooltipContent>
    </Tooltip>
  );
}

function ConcludedPurlTooltip({ component }: { component: entities.ComponentDTO }) {
  const result = useSelectedResult();
  const matchPresentation = matchTypePresentation[component.id as MatchType];

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <div className="cursor-pointer">
          <div className={clsx('text-lg font-bold leading-tight', matchPresentation.accent)}>
            {result?.concluded_name || component.component}
          </div>
          <div>{result?.concluded_purl}</div>
        </div>
      </TooltipTrigger>
      <TooltipContent side="bottom" align="start" className="p-4">
        <div className="flex flex-col gap-4">
          <div>
            <p className="font-medium">PURL</p>
            <p className="text-muted-foreground">{result?.concluded_purl}</p>
          </div>
          {result?.concluded_purl_url && (
            <div>
              <p className="font-medium">URL</p>
              <Link to={result?.concluded_purl_url as string} />
            </div>
          )}
        </div>
      </TooltipContent>
    </Tooltip>
  );
}
