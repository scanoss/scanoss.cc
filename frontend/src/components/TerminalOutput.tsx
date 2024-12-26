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
    setTimeout(scrollToBottom, 0);
  }, [lines, autoScroll, isUserScrolling]);

  return (
    <div ref={scrollViewportRef} className="max-h-80 w-full overflow-y-auto rounded bg-gray-900 p-4 font-mono text-sm">
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
