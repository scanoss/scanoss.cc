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

import { Check, CircleDotDashed } from 'lucide-react';
import { ReactNode } from 'react';

import { FilterAction } from '@/modules/components/domain';

export enum MatchType {
  File = 'file',
  Snippet = 'snippet',
}

interface ResultStatusPresentation {
  badgeStyles: string;
  label: string;
  icon: ReactNode;
}

export const resultStatusPresentation: Record<string, ResultStatusPresentation> = {
  pending: {
    badgeStyles: 'bg-[color:hsl(200,40%,20%)] text-white border-[color:hsl(200,40%,30%)] hover:bg-[color:hsl(200,40%,20%)]',
    label: 'Pending',
    icon: <CircleDotDashed className="h-3 w-3" />,
  },
  completed: {
    badgeStyles: 'bg-[color:hsl(120,50%,10%)] text-white border-[color:hsl(120,45%,25%)] hover:bg-[color:hsl(120,50%,10%)]',
    label: 'Completed',
    icon: <Check className="h-3 w-3" />,
  },
};

interface MatchTypePresentation {
  background: string;
  foreground: string;
  accent: string;
  muted: string;
  label: string;
}

export const matchTypePresentation: Record<MatchType, MatchTypePresentation> = {
  [MatchType.File]: {
    background: 'bg-file',
    foreground: 'text-file-foreground',
    accent: 'text-file-accent',
    muted: 'text-muted-foreground',
    label: 'File',
  },
  [MatchType.Snippet]: {
    background: 'bg-snippet',
    foreground: 'text-snippet-foreground',
    accent: 'text-snippet-accent',
    muted: 'text-muted-foreground',
    label: 'Snippet',
  },
};

interface StateInfoPresentation {
  label: string;
  stateInfoContainerStyles: string;
  stateInfoSidebarIndicatorStyles: string;
  stateInfoTextStyles: string;
}

export const stateInfoPresentation: Record<FilterAction, StateInfoPresentation> = {
  [FilterAction.Include]: {
    label: 'Included',
    stateInfoContainerStyles: 'border-l-4 border-green-600 border-l-green-600 bg-green-950',
    stateInfoSidebarIndicatorStyles: 'bg-green-600',
    stateInfoTextStyles: 'text-green-600',
  },
  [FilterAction.Remove]: {
    label: 'Removed',
    stateInfoContainerStyles: 'border-l-4 border-green-600 border-l-green-600 bg-green-950',
    stateInfoSidebarIndicatorStyles: 'bg-red-600',
    stateInfoTextStyles: 'text-green-600',
  },
  [FilterAction.Replace]: {
    label: 'Replaced',
    stateInfoContainerStyles: 'border-l-4 border-green-600 border-l-green-600 bg-green-950',
    stateInfoSidebarIndicatorStyles: 'bg-yellow-600',
    stateInfoTextStyles: 'text-green-600',
  },
  [FilterAction.Ignore]: {
    label: 'Ignored',
    stateInfoContainerStyles: 'border-l-4 border-gray-600 border-l-gray-600 bg-gray-950',
    stateInfoSidebarIndicatorStyles: 'bg-gray-600',
    stateInfoTextStyles: 'text-gray-600',
  },
};
