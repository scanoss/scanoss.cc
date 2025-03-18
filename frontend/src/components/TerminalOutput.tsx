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

import { useEffect, useRef, useState } from 'react';

interface OutputLine {
  type: 'stdout' | 'stderr' | 'error';
  text: string;
}

interface TerminalOutputProps {
  lines: OutputLine[];
  autoScroll?: boolean;
}

export default function TerminalOutput({ lines, autoScroll = true }: TerminalOutputProps) {
  const scrollViewportRef = useRef<HTMLDivElement | null>(null);
  const [isUserScrolling, setIsUserScrolling] = useState(false);

  useEffect(() => {
    const scrollViewport = scrollViewportRef.current;
    if (!scrollViewport) return;

    const handleScroll = () => {
      const { scrollTop, scrollHeight, clientHeight } = scrollViewport;
      const isAtBottom = Math.abs(scrollHeight - scrollTop - clientHeight) < 50;

      setIsUserScrolling(!isAtBottom);
    };

    scrollViewport.addEventListener('scroll', handleScroll);
    return () => scrollViewport.removeEventListener('scroll', handleScroll);
  }, []);

  useEffect(() => {
    const scrollViewport = scrollViewportRef.current;
    if (!scrollViewport || !autoScroll || isUserScrolling) return;

    const scrollToBottom = () => {
      scrollViewport.scrollTop = scrollViewport.scrollHeight;
    };

    scrollToBottom();
    requestAnimationFrame(scrollToBottom);
    const timeoutId = setTimeout(scrollToBottom, 10);

    return () => clearTimeout(timeoutId);
  }, [lines, autoScroll, isUserScrolling]);

  return (
    <div ref={scrollViewportRef} className="h-full max-h-60 w-full overflow-y-auto rounded bg-gray-900 p-4 font-mono text-sm">
      {lines.map((line, i) => (
        <div
          key={i}
          className={` ${
            line.type === 'error' ? 'text-red-400' : line.type === 'stderr' ? 'text-orange-400' : 'text-gray-200'
          } whitespace-pre-wrap font-mono leading-5`}
        >
          <code style={{ wordBreak: 'break-word', whiteSpace: 'pre-wrap' }}>{`${line.type === 'stderr' ? '❯ ' : ''}${line.text}`}</code>
        </div>
      ))}
    </div>
  );
}
