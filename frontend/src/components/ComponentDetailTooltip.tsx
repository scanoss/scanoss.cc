import clsx from 'clsx';
import { entities } from 'wailsjs/go/models';

import Link from '@/components/Link';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import useSelectedResult from '@/hooks/useSelectedResult';
import { MatchType, matchTypePresentation } from '@/modules/results/domain';

interface ComponentDetailTooltipProps {
  component: entities.ComponentDTO;
}

export default function ComponentDetailTooltip({ component }: ComponentDetailTooltipProps) {
  const result = useSelectedResult();
  const matchPresentation = matchTypePresentation[component.id as MatchType];

  const isPurlReplaced = !!result?.purl?.concluded;
  const currentPurl = isPurlReplaced ? result?.purl?.concluded : result?.purl?.detected;
  const originallyDetectedPurl = result?.purl?.detected;
  const purlTooltipLabel = isPurlReplaced ? 'Detected PURL' : 'PURL';
  const shouldShowUrl = isPurlReplaced ? !!result.purl?.concluded_purl_url : !!component.url;
  const urlToShow = isPurlReplaced ? result.purl?.concluded_purl_url : component.url;

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <div className="cursor-pointer">
          <div className={clsx('text-lg font-bold leading-tight', matchPresentation.accent)}>{component.component}</div>
          <div>{currentPurl}</div>
        </div>
      </TooltipTrigger>
      <TooltipContent side="bottom" align="start" className="p-4">
        <div className="flex flex-col gap-4">
          {component.purl?.length ? (
            <div>
              <p className="font-medium">{purlTooltipLabel}</p>
              <p className="text-muted-foreground">{isPurlReplaced ? originallyDetectedPurl : currentPurl}</p>
            </div>
          ) : null}
          <div>
            <p className="font-medium">VERSION</p>
            <p className="text-muted-foreground">{component.version}</p>
          </div>
          {!isPurlReplaced && component.licenses?.length ? (
            <div>
              <p className="font-medium">LICENSE</p>
              <p className="text-muted-foreground">{component.licenses?.[0].name}</p>
            </div>
          ) : null}
          {shouldShowUrl && (
            <div>
              <p className="font-medium">URL</p>
              <Link to={urlToShow as string} />
            </div>
          )}
        </div>
      </TooltipContent>
    </Tooltip>
  );
}
