import { Editor } from '@monaco-editor/react';
import * as monaco from 'monaco-editor';
import { useRef } from 'react';

import { MonacoManager } from '@/lib/editor';
import { getHighlightLineRanges } from '@/modules/results/utils';

import { Skeleton } from './ui/skeleton';

interface CodeViewerProps {
  content: string | undefined;
  editorId: string;
  editorType: 'local' | 'remote';
  height?: string;
  highlightLines?: string;
  isError: boolean;
  isLoading: boolean;
  language: string | null | undefined;
  width?: string;
}

export default function CodeViewer({
  editorId,
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
  const monacoManager = MonacoManager.getInstance();

  const handleEditorMount = (editor: monaco.editor.IStandaloneCodeEditor) => {
    console.log('CodeViewer handleEditorMount');
    editorRef.current = editor;
    monacoManager.addEditor(editorId, editor);

    const className = editorType === 'local' ? 'bg-highlight-local-line' : 'bg-highlight-remote-line';

    if (highlightAll) {
      const totalLines = editor.getModel()?.getLineCount();
      if (totalLines) {
        monacoManager.highlightLines(editorId, [{ start: 1, end: totalLines }], className);
      }
      return;
    }

    if (highlightLines) {
      const ranges = getHighlightLineRanges(highlightLines);

      monacoManager.highlightLines(editorId, ranges, className);
      setTimeout(() => {
        monacoManager.scrollToLine(editorId, ranges[0].start);
      }, 200);
    }
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
      }}
    />
  );
}
