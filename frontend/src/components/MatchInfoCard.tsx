import { useQuery } from '@tanstack/react-query';
import clsx from 'clsx';
import { MessageSquareText } from 'lucide-react';

import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import useSelectedResult from '@/hooks/useSelectedResult';
import { FilterAction } from '@/modules/components/domain';
import useLocalFilePath from '@/modules/files/hooks/useLocalFilePath';
import {
  MatchType,
  matchTypePresentation,
  resultStatusPresentation,
  stateInfoPresentation,
} from '@/modules/results/domain';
import ResultService from '@/modules/results/infra/service';

import ComponentDetailTooltip from './ComponentDetailTooltip';

export default function MatchInfoCard() {
  const result = useSelectedResult();
  const localFilePath = useLocalFilePath();

  const { data: component } = useQuery({
    queryKey: ['component', localFilePath],
    queryFn: () => ResultService.getComponent(localFilePath),
  });

  if (!component) {
    return <Skeleton className="min-h-[68px] w-full"></Skeleton>;
  }

  const status = result?.workflow_state;
  const matchPresentation = matchTypePresentation[component.id as MatchType];

  const isResultFilteredByFile = result?.filter_config?.type === 'by_file';
  const isResultFilteredByPurl = result?.filter_config?.type === 'by_purl';

  const filterPresentation = stateInfoPresentation[result?.filter_config?.action as FilterAction];

  return (
    <div
      className={clsx(
        'flex items-center justify-between rounded-sm p-3',
        matchPresentation.background,
        matchPresentation.foreground
      )}
    >
      <div className="flex flex-wrap items-center gap-8 text-sm">
        <ComponentDetailTooltip component={component} />
        {component.version && (
          <div>
            <div className={matchPresentation.muted}>Version</div>
            <div>{component.version}</div>
          </div>
        )}
        {component.licenses?.length ? (
          <div>
            <div className={matchPresentation.muted}>License</div>
            <div>{component.licenses?.[0].name}</div>
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
            <Badge
              className={clsx('flex items-center gap-1 font-normal', resultStatusPresentation[status].badgeStyles)}
            >
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
                <Badge className="flex items-center gap-1 font-normal">
                  {result.comment && <MessageSquareText className="h-3 w-3" />}
                  {filterPresentation?.label}
                  {isResultFilteredByFile && ' file'}
                  {isResultFilteredByPurl && ` component`}
                </Badge>
              </div>
            </TooltipTrigger>
            {result.comment && (
              <TooltipContent side="bottom" align="start" className="px-4 py-2">
                <pre className="m-0 whitespace-pre-wrap break-words font-sans text-muted-foreground">
                  {result.comment}
                </pre>
              </TooltipContent>
            )}
          </Tooltip>
        )}
      </div>
    </div>
  );
}
