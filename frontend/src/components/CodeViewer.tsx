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

import * as monaco from 'monaco-editor';
import { useEffect, useMemo, useRef } from 'react';

import { HighlightRange, MonacoManager } from '@/lib/editor';
import { getHighlightLineRanges } from '@/modules/results/utils';

interface CodeViewerProps {
  content: string | undefined;
  editorId: string;
  error: Error | null;
  height?: string;
  highlightLines?: string;
  isError: boolean;
  language: string | null | undefined;
  width?: string;
}

export default function CodeViewer({
  content,
  editorId,
  error,
  height = '100%',
  highlightLines,
  isError,
  language,
  width = '100%',
}: CodeViewerProps) {
  const highlightAll = highlightLines === 'all';
  const editor = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
  const editorContainer = useRef<HTMLDivElement>(null);
  const monacoManager = MonacoManager.getInstance();

  const decorationOptions = useMemo<monaco.editor.IModelDecorationOptions>(
    () => ({
      isWholeLine: true,
      className: 'bg-highlight-line',
      linesDecorationsClassName: 'line-range-decoration',
    }),
    []
  );

  const editorDefaultOptions: monaco.editor.IStandaloneEditorConstructionOptions = useMemo(
    () => ({
      readOnly: true,
      theme: 'vs-dark',
      minimap: { enabled: false },
      automaticLayout: true,
      fontSize: 12,
    }),
    []
  );

  const initMonaco = () => {
    if (editorContainer.current) {
      editor.current = monaco.editor.create(editorContainer.current, {
        language: language as string,
        model: null,
        ...editorDefaultOptions,
      });

      monacoManager.addEditor(editorId, editor.current);
    }
  };

  const destroyMonaco = () => {
    if (editor.current) {
      editor.current.dispose();
      monacoManager.removeEditor(editorId);
    }
  };

  const updateContent = () => {
    const { current: mEditor } = editor;
    let oldModel: monaco.editor.ITextModel | null = null;

    if (mEditor) {
      oldModel = mEditor.getModel();

      const nModel = monaco.editor.createModel(content as string, language as string);
      mEditor.setModel(nModel);

      if (oldModel) oldModel.dispose();
    }
  };

  const updateHighlight = () => {
    if (!editor.current) return;

    const model = editor.current.getModel();
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

    editor.current.deltaDecorations([], decorations);
  };

  useEffect(() => {
    initMonaco();
    return destroyMonaco;
  }, []);

  useEffect(() => {
    updateContent();
    updateHighlight();
  }, [content, language]);

  if (isError) {
    return (
      <div className="flex flex-1 items-center justify-center text-sm text-muted-foreground">
        <span className="mx-auto max-w-[50%] text-center">{error instanceof Error ? error.message : error}</span>
      </div>
    );
  }

  return <div ref={editorContainer} style={{ width, height }} />;
}
