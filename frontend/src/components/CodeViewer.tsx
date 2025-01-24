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

import { Editor } from '@monaco-editor/react';
import * as monaco from 'monaco-editor';
import { useCallback, useEffect, useMemo, useRef } from 'react';

import { HighlightRange, MonacoManager } from '@/lib/editor';
import { getHighlightLineRanges } from '@/modules/results/utils';

import Loading from './Loading';
import { Skeleton } from './ui/skeleton';

interface CodeViewerProps {
  content: string | undefined;
  editorId: string;
  editorType: 'local' | 'remote';
  error: Error | null;
  height?: string;
  highlightLines?: string;
  isError: boolean;
  isLoading: boolean;
  language: string | null | undefined;
  width?: string;
}

export default function CodeViewer({
  content,
  editorId,
  editorType,
  error,
  height = '100%',
  highlightLines,
  isError,
  isLoading,
  language,
  width = '100%',
}: CodeViewerProps) {
  const highlightAll = highlightLines === 'all';
  const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
  const scrollListenerRef = useRef<monaco.IDisposable | null>(null);
  const monacoManager = MonacoManager.getInstance();

  const decorationOptions = useMemo<monaco.editor.IModelDecorationOptions>(
    () => ({
      isWholeLine: true,
      className: 'bg-highlight-line',
      linesDecorationsClassName: 'line-range-decoration',
    }),
    []
  );

  const applyHighlights = useCallback(
    (editor: monaco.editor.IStandaloneCodeEditor) => {
      const model = editor.getModel();
      if (!model) return;

      let highlightRanges: HighlightRange[] = [];

      if (highlightAll) {
        // For highlight all, highlight all lines
        highlightRanges.push({
          start: 1,
          end: model.getLineCount(),
        });
      } else if (highlightLines) {
        highlightRanges = getHighlightLineRanges(highlightLines);
      }

      // Create and apply decorations
      const decorations: monaco.editor.IModelDeltaDecoration[] = highlightRanges.map(({ start, end }) => ({
        range: new monaco.Range(start, 1, end, 1),
        options: decorationOptions,
      }));

      editor.deltaDecorations([], decorations);
    },
    [highlightAll, highlightLines, decorationOptions]
  );

  const handleEditorMount = useCallback(
    (editor: monaco.editor.IStandaloneCodeEditor) => {
      editorRef.current = editor;

      // Apply highlights once on mount
      applyHighlights(editor);

      const highlightRanges = highlightLines ? getHighlightLineRanges(highlightLines) : [];
      monacoManager.addEditor(editorId, editor, {
        revealLine: highlightRanges[0]?.start,
      });
    },
    [editorId, highlightLines, monacoManager, applyHighlights]
  );

  useEffect(() => {
    return () => {
      if (scrollListenerRef.current) {
        scrollListenerRef.current.dispose();
        scrollListenerRef.current = null;
      }

      monacoManager.dispose();
    };
  }, [monacoManager]);

  if (isLoading || !highlightLines) {
    return (
      <div className="flex flex-1 items-center justify-center text-muted-foreground">
        <Loading text={`Loading ${editorType} file...`} />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="flex flex-1 items-center justify-center text-sm text-muted-foreground">
        <span className="mx-auto max-w-[50%] text-center">{error instanceof Error ? error.message : error}</span>
      </div>
    );
  }

  return (
    <Editor
      height={height}
      loading={<Skeleton className="h-full w-full" />}
      onMount={handleEditorMount}
      value={content}
      width={width}
      theme="vs-dark"
      {...(language ? { language } : {})}
      options={{
        minimap: { enabled: false },
        readOnly: true,
      }}
    />
  );
}
