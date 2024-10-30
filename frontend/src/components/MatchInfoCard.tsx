import { useQuery } from '@tanstack/react-query';
import clsx from 'clsx';
import { MessageSquareText } from 'lucide-react';

import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import useSelectedResult from '@/hooks/useSelectedResult';
import { FilterAction } from '@/modules/components/domain';
import useLocalFilePath from '@/hooks/useLocalFilePath';
import { MatchType, matchTypePresentation, resultStatusPresentation, stateInfoPresentation } from '@/modules/results/domain';

import { GetComponentByPath } from '../../wailsjs/go/service/ComponentServiceImpl';
import ComponentDetailTooltip from './ComponentDetailTooltip';

export default function MatchInfoCard() {
  const result = useSelectedResult();
  const localFilePath = useLocalFilePath();

  const { data: component } = useQuery({
    queryKey: ['component', localFilePath],
    queryFn: () => GetComponentByPath(localFilePath),
  });

  if (!component) {
    return <Skeleton className="min-h-[68px] w-full"></Skeleton>;
  }

  const status = result?.workflow_state;
  const matchPresentation = matchTypePresentation[component.id as MatchType];

  const isResultFilteredByFile = result?.filter_config?.type === 'by_file';
  const isResultFilteredByPurl = result?.filter_config?.type === 'by_purl';

  const filterPresentation = stateInfoPresentation[result?.filter_config?.action as FilterAction];

  const isResultReplaced = result?.filter_config?.action === FilterAction.Replace;
  const isResultRemoved = result?.filter_config?.action === FilterAction.Remove;

  const shouldShowVersion = !isResultReplaced && component.version;
  const shouldShowLicense = !isResultReplaced && component.licenses?.length;

  const removedStyles = isResultRemoved ? 'line-through opacity-70 text-muted-foreground' : '';

  return (
    <div
      className={clsx(
        'flex items-center justify-between rounded-sm p-3',
        matchPresentation.foreground,
        matchPresentation.background,
        isResultRemoved && '!bg-card'
      )}
    >
      <div className="flex flex-wrap items-center gap-8 text-sm">
        <ComponentDetailTooltip component={component} />
        {shouldShowVersion && (
          <div>
            <div className={matchPresentation.muted}>Version</div>
            <div className={clsx(removedStyles)}>{component.version}</div>
          </div>
        )}
        {shouldShowLicense ? (
          <div>
            <div className={matchPresentation.muted}>License</div>
            <div className={clsx(removedStyles)}>{component.licenses?.[0].name}</div>
          </div>
        ) : null}
        <div>
          <div className={matchPresentation.muted}>Detected</div>
          <div className={matchPresentation.accent}>{matchPresentation.label}</div>
        </div>
        <div>
          <div className={matchPresentation.muted}>Match</div>
          <div>{component.matched}</div>
        </div>
        {status && (
          <div>
            <div className={matchPresentation.muted}>Status</div>
            <Badge className={clsx('flex items-center gap-1 font-normal', resultStatusPresentation[status].badgeStyles)}>
              {resultStatusPresentation[status].icon}
              {resultStatusPresentation[status].label}
            </Badge>
          </div>
        )}
        {result?.filter_config?.action && result.filter_config.type && (
          <Tooltip>
            <TooltipTrigger asChild>
              <div>
                <div className={matchPresentation.muted}>Decision</div>
                <Badge className={clsx('flex items-center gap-1 font-normal')}>
                  {result.comment && <MessageSquareText className="h-3 w-3" />}
                  {filterPresentation?.label}
                  {isResultFilteredByFile && ' file'}
                  {isResultFilteredByPurl && ` component`}
                </Badge>
              </div>
            </TooltipTrigger>
            {result.comment && (
              <TooltipContent side="bottom" align="start" className="px-4 py-2">
                <pre className="m-0 whitespace-pre-wrap break-words font-sans text-muted-foreground">{result.comment}</pre>
              </TooltipContent>
            )}
          </Tooltip>
        )}
      </div>
    </div>
  );
}
