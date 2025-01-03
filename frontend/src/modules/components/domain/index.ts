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

export enum FilterAction {
  Ignore = 'ignore',
  Include = 'include',
  Remove = 'remove',
  Replace = 'replace',
}

export type FilterBy = 'path' | 'purl';

export const filterActionLabelMap: Record<FilterAction, string> = {
  [FilterAction.Ignore]: 'Omit / Skip',
  [FilterAction.Include]: 'Include',
  [FilterAction.Remove]: 'Dismiss',
  [FilterAction.Replace]: 'Replace',
};

export interface OnAddFilterArgs {
  action: FilterAction;
  filterBy: 'by_file' | 'by_purl';
  withComment?: boolean;
}

export const VALID_PURL_REGEX =
  /^pkg:[a-zA-Z0-9.\-_]+(?:\/[a-zA-Z0-9.\-_]+)?\/[a-zA-Z0-9.\-_]+(?:@[a-zA-Z0-9.\-_]+)?(?:\?[a-zA-Z0-9.\-_=&]+)?(?:#[a-zA-Z0-9.\-_/]+)?$/;
