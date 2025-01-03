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
