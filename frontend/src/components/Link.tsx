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
import { ReactNode } from 'react';

import { BrowserOpenURL } from '../../wailsjs/runtime';

export default function Link({ children, to, className, ...props }: { children: ReactNode; to: string; className?: string }) {
  return (
    <a
      onClick={(e) => {
        e.preventDefault();
        BrowserOpenURL(to);
      }}
      className={clsx(className, 'cursor-pointer text-blue-500 underline')}
      {...props}
    >
      {children}
    </a>
  );
}
