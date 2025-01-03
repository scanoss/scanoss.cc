// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
  const removedStyles = isResultRemoved || replaced ? 'line-through opacity-70 text-muted-foreground' : '';

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <div className={clsx('cursor-pointer', removedStyles)}>
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
              <p className={clsx('text-muted-foreground')}>{component.version}</p>
            </div>
          )}
          {component.licenses?.length ? (
            <div>
              <p className="font-medium">LICENSE</p>
              <p className={clsx('text-muted-foreground')}>{component.licenses?.[0].name}</p>
            </div>
          ) : null}
          {component.url && (
            <div>
              <p className="font-medium">URL</p>
              <Link to={component.url as string}>{component.url}</Link>
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
          <div className={clsx('text-lg font-bold leading-tight', matchPresentation.accent)}>{result?.concluded_name || component.component}</div>
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
              <Link to={result?.concluded_purl_url as string}>{result?.concluded_purl_url}</Link>
            </div>
          )}
        </div>
      </TooltipContent>
    </Tooltip>
  );
}
