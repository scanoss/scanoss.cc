import { Editor } from '@monaco-editor/react';
import * as monaco from 'monaco-editor';
import { useRef } from 'react';

import { getHighlightLineRanges } from '@/modules/results/utils';

import { Skeleton } from './ui/skeleton';

interface CodeViewerProps {
  content: string | undefined;
  editorType: 'local' | 'remote';
  height?: string;
  highlightLines?: string;
  isError: boolean;
  isLoading: boolean;
  language: string | null | undefined;
  width?: string;
}

export default function CodeViewer({
  content,
  editorType,
  height = '100%',
  highlightLines,
  isError,
  isLoading,
  language,
  width = '100%',
}: CodeViewerProps) {
  const highlightAll = highlightLines === 'all';
  const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);

  const handleEditorMount = async (
    editor: monaco.editor.IStandaloneCodeEditor
  ) => {
    editorRef.current = editor;

    if (highlightAll) {
      const totalLines = editor.getModel()?.getLineCount();

      if (!totalLines) return;

      const decorations: monaco.editor.IModelDeltaDecoration[] = [
        {
          range: new monaco.Range(1, 1, totalLines, 1),
          options: {
            isWholeLine: true,
            className:
              editorType === 'local'
                ? 'bg-highlight-local-line'
                : 'bg-highlight-remote-line',
          },
        },
      ];

      editorRef.current?.createDecorationsCollection(decorations);

      return;
    }

    if (!highlightLines) return;

    const ranges = getHighlightLineRanges(highlightLines);

    const decorations: monaco.editor.IModelDeltaDecoration[] = ranges.map(
      ({ start, end }) => ({
        range: new monaco.Range(start, 1, end, 1),
        options: {
          isWholeLine: true,
          className:
            editorType === 'local'
              ? 'bg-highlight-local-line'
              : 'bg-highlight-remote-line',
          inlineClassName:
            editorType === 'local'
              ? 'bg-highlight-local-inline'
              : 'bg-highlight-remote-inline',
        },
      })
    );

    editorRef.current?.createDecorationsCollection(decorations);
  };

  if (isLoading || !highlightLines) {
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
      className={`editor-${editorType}`}
      height={height}
      loading={<Skeleton className="h-full w-full" />}
      onMount={handleEditorMount}
      theme="vs-dark"
      value={content}
      width={width}
      {...(language ? { language } : {})}
      options={{
        minimap: { enabled: false },
        readOnly: true,
        wordWrap: 'on',
        smoothScrolling: true,
      }}
    />
  );
}
