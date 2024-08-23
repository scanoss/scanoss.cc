import { Editor } from '@monaco-editor/react';
import * as monaco from 'monaco-editor';
import React, { useRef } from 'react';

import { getHighlightLineRanges } from '@/modules/results/utils';

import { Skeleton } from './ui/skeleton';

interface CodeViewerProps {
  content: string | undefined;
  height?: string;
  highlightLines?: string;
  isError: boolean;
  isLoading: boolean;
  language: string | null | undefined;
  width?: string;
}

export default function CodeViewer({
  content,
  height = '100%',
  highlightLines,
  isError,
  isLoading,
  language,
  width = '100%',
}: CodeViewerProps) {
  const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);

  const handleEditorMount = (editor: monaco.editor.IStandaloneCodeEditor) => {
    if (!highlightLines) {
      return;
    }

    editorRef.current = editor;

    const ranges = getHighlightLineRanges(highlightLines);

    const decorations = ranges.map(({ start, end }) => ({
      range: new monaco.Range(start, 1, end, 1),
      options: {
        isWholeLine: true,
        className: 'lineHighlightDecoration',
        linesDecorationsClassName: 'lineRangeDecoration',
      },
    }));

    editorRef.current?.deltaDecorations([], decorations);
  };

  if (isLoading) {
    return <Skeleton className="h-full w-full" />;
  }

  if (isError) {
    return (
      <div className="flex h-full w-full items-center justify-center">
        <p className="text-red-500">Error loading file</p>
      </div>
    );
  }

  return (
    <Editor
      height={height}
      onMount={handleEditorMount}
      theme="vs-dark"
      value={content}
      width={width}
      {...(language ? { language } : {})}
      loading={<Skeleton className="h-full w-full" />}
      options={{
        minimap: { enabled: false },
        readOnly: true,
        wordWrap: 'on',
      }}
    />
  );
}
