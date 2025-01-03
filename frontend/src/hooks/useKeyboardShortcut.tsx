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

import { Options, useHotkeys } from 'react-hotkeys-hook';
import { HotkeysEvent } from 'react-hotkeys-hook/dist/types';

// We re export this to use it always with default configurations already setup
export default function useKeyboardShortcut(
  keys: string | string[],
  callback: (event: KeyboardEvent, handler: HotkeysEvent) => void,
  options?: Options,
  deps: unknown[] = []
) {
  return useHotkeys(
    keys,
    callback,
    {
      preventDefault: true,
      ...options,
    },
    deps
  );
}
