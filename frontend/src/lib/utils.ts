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

import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function getFileName(filePath: string) {
  const parts = filePath.split('/');
  return parts.pop() || '';
}

export function getDirectory(filePath: string) {
  const parts = filePath.split('/');
  return parts.slice(0, -1).join('/');
}

export function truncatePath(path: string, maxLength: number = 30) {
  if (path.length <= maxLength) return path;
  const parts = path.split('/');
  if (parts.length <= 2) return path;

  const isAbsolute = path.startsWith('/');
  const first = isAbsolute ? `/${parts[1]}` : parts[0];
  const last = parts[parts.length - 1];
  return `${first}/.../${last}`;
}

export const isDefaultPath = (path: string, platform: string | undefined) => {
  switch (platform) {
    case 'darwin':
      return path === '/' || path === '/Users';
    case 'windows':
      return path === 'C:\\' || path === '';
    default:
      return path === '/' || path === '/home';
  }
};
