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
