import { useEffect, useRef, useState } from 'react';

import { ScrollArea } from './ui/scroll-area';

interface OutputLine {
  type: 'stdout' | 'stderr' | 'error';
  text: string;
}

interface TerminalOutputProps {
  lines: OutputLine[];
  autoScroll?: boolean;
}

export default function TerminalOutput({ lines, autoScroll = true }: TerminalOutputProps) {
  const terminalRef = useRef<HTMLDivElement>(null);
  const [isUserScrolling, setIsUserScrolling] = useState(false);
  const [isProgrammaticScroll, setIsProgrammaticScroll] = useState(false);

  useEffect(() => {
    const scrollViewport = terminalRef.current?.querySelector('[data-radix-scroll-area-viewport]');

    const handleScroll = () => {
      if (!scrollViewport) return;

      if (isProgrammaticScroll) {
        setIsProgrammaticScroll(false);
        return;
      }

      const isAtBottom = scrollViewport.scrollHeight - scrollViewport.scrollTop <= scrollViewport.clientHeight + 10;
      setIsUserScrolling(!isAtBottom);
    };

    scrollViewport?.addEventListener('scroll', handleScroll);
    return () => scrollViewport?.removeEventListener('scroll', handleScroll);
  }, [isProgrammaticScroll]);

  useEffect(() => {
    const scrollViewport = terminalRef.current?.querySelector('[data-radix-scroll-area-viewport]');

    if (autoScroll && !isUserScrolling && terminalRef.current && scrollViewport) {
      setIsProgrammaticScroll(true);
      scrollViewport.scrollTo({
        top: scrollViewport.scrollHeight,
      });
    }
  }, [lines.length, autoScroll, isUserScrolling]);

  return (
    <ScrollArea ref={terminalRef} className="max-h-80 w-full overflow-y-auto rounded bg-gray-900 p-4 font-mono text-sm">
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
    </ScrollArea>
  );
}
