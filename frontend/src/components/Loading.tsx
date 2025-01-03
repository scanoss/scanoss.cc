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

interface Props {
  text?: string;
  size?: string;
  textSize?: string;
}

export default function Loading({ text = 'Loading', size = 'w-4 h-4', textSize = 'text-sm' }: Props) {
  return (
    <div className="inline-flex items-center justify-center leading-none">
      <svg className={`${size} animate-spin`} viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
        <g className="opacity-40">
          <path
            fill="currentColor"
            d="M8 0a8 8 0 0 0-8 8 8 8 0 0 0 8 8 8 8 0 0 0 8-8 8 8 0 0 0-8-8zm0 14.4A6.4 6.4 0 0 1 1.6 8 6.4 6.4 0 0 1 8 1.6a6.4 6.4 0 0 1 6.4 6.4 6.4 6.4 0 0 1-6.4 6.4z"
          />
        </g>
        <path fill="currentColor" d="M8 0a8 8 0 0 1 8 8h-1.6A6.4 6.4 0 0 0 8 1.6V0z" />
      </svg>
      <span className={`ml-2 ${textSize} leading-none`}>{text}</span>
    </div>
  );
}
