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

import { useDebounceCallback } from '@/hooks/useDebounceCallback';
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
  const decorationsRef = useRef<string[]>([]);
  const scrollListenerRef = useRef<monaco.IDisposable | null>(null);
  const monacoManager = MonacoManager.getInstance();

  const decorationOptions = useMemo<monaco.editor.IModelDecorationOptions>(
    () => ({
      isWholeLine: true,
      className: 'bg-highlight-line',
    }),
    []
  );

  const updateDecorations = useCallback(
    (editor: monaco.editor.IStandaloneCodeEditor) => {
      const model = editor.getModel();
      if (!model) return;

      const visibleRanges = editor.getVisibleRanges();
      if (!visibleRanges.length) return;

      const BUFFER_LINES = 50;
      const firstVisibleLine = Math.max(1, visibleRanges[0].startLineNumber - BUFFER_LINES);
      const lastVisibleLine = Math.min(model.getLineCount(), visibleRanges[visibleRanges.length - 1].endLineNumber + BUFFER_LINES);

      let highlightRanges: HighlightRange[] = [];

      if (highlightAll) {
        // For highlight all, only create decorations for visible lines
        highlightRanges.push({
          start: firstVisibleLine,
          end: lastVisibleLine,
        });
      } else if (highlightLines) {
        const allRanges = getHighlightLineRanges(highlightLines);
        // Filter ranges to only those that intersect with visible lines
        highlightRanges = allRanges.filter((range) => range.start <= lastVisibleLine && range.end >= firstVisibleLine);
      }

      // Create decorations only for visible ranges
      const decorations: monaco.editor.IModelDeltaDecoration[] = highlightRanges.map(({ start, end }) => ({
        range: new monaco.Range(Math.max(start, firstVisibleLine), 1, Math.min(end, lastVisibleLine), 1),
        options: decorationOptions,
      }));

      decorationsRef.current = editor.deltaDecorations(decorationsRef.current, decorations);
    },
    [highlightAll, highlightLines, decorationOptions]
  );

  const debouncedUpdate = useDebounceCallback(updateDecorations, 16);

  const handleEditorMount = useCallback(
    (editor: monaco.editor.IStandaloneCodeEditor) => {
      editorRef.current = editor;

      updateDecorations(editor);

      scrollListenerRef.current = editor.onDidScrollChange(() => {
        debouncedUpdate(editor);
      });

      const highlightRanges = highlightLines ? getHighlightLineRanges(highlightLines) : [];
      monacoManager.addEditor(editorId, editor, {
        revealLine: highlightRanges[0]?.start,
      });
    },
    [editorId, highlightLines, monacoManager, updateDecorations, debouncedUpdate]
  );

  useEffect(() => {
    return () => {
      if (editorRef.current) {
        decorationsRef.current = editorRef.current.deltaDecorations(decorationsRef.current, []);
      }

      if (scrollListenerRef.current) {
        scrollListenerRef.current.dispose();
        scrollListenerRef.current = null;
      }

      monacoManager.dispose();
    };
  }, [monacoManager, highlightLines]);

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
        wordWrap: 'on',
      }}
    />
  );
}
