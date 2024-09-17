import { useQuery } from '@tanstack/react-query';
import clsx from 'clsx';

import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import {
  MatchType,
  matchTypePresentation,
  resultStatusPresentation,
} from '@/modules/results/domain';
import ResultService from '@/modules/results/infra/service';
import { useResults } from '@/modules/results/providers/ResultsProvider';

import useLocalFilePath from '../hooks/useLocalFilePath';

export default function MatchInfoCard() {
  const { results } = useResults();
  const localFilePath = useLocalFilePath();

  const { data: component } = useQuery({
    queryKey: ['component', localFilePath],
    queryFn: () => ResultService.getComponent(localFilePath),
  });

  if (!component) {
    return <Skeleton className="min-h-[68px] w-full"></Skeleton>;
  }

  const result = results.find((result) => result.path === localFilePath);

  const status = result?.state;
  const matchPresentation = matchTypePresentation[component.id as MatchType];

  return (
    <div
      className={clsx(
        'flex items-center justify-between rounded-sm p-3',
        matchPresentation.background,
        matchPresentation.foreground
      )}
    >
      <div className="flex flex-wrap items-center gap-8 text-sm">
        <div>
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
          <div className={matchPresentation.accent}>
            {matchPresentation.label}
          </div>
        </div>
        <div>
          <div className={matchPresentation.muted}>Match</div>
          <div>{component.matched}</div>
        </div>
        {status && (
          <div>
            <div className={matchPresentation.muted}>Status</div>
            <Badge
              className={clsx(
                'flex items-center gap-1 font-normal',
                resultStatusPresentation[status].badgeStyles
              )}
            >
              {resultStatusPresentation[status].icon}
              {resultStatusPresentation[status].label}
            </Badge>
          </div>
        )}
      </div>
    </div>
  );
}
