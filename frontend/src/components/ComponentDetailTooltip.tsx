// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

import clsx from 'clsx';
import { ArrowRight, ExternalLink } from 'lucide-react';
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
          <Link className="flex items-center gap-1" to={result?.detected_purl_url as string}>
            {result?.detected_purl}
            <ExternalLink className="h-3 w-3" />
          </Link>
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
