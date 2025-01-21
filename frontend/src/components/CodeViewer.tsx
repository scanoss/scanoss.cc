// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import { Editor } from '@monaco-editor/react';
import * as monaco from 'monaco-editor';
import { useRef } from 'react';

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
  const monacoManager = MonacoManager.getInstance();

  const handleEditorMount = (editor: monaco.editor.IStandaloneCodeEditor) => {
    editorRef.current = editor;

    const highlightRanges: HighlightRange[] = [];

    if (highlightAll) {
      const totalLines = editor.getModel()?.getLineCount();
      if (totalLines) {
        highlightRanges.push({ start: 1, end: totalLines });
      }
    }

    if (highlightLines) {
      const ranges = getHighlightLineRanges(highlightLines);
      highlightRanges.push(...ranges);
    }

    monacoManager.addEditor(editorId, editor, {
      highlight: { ranges: highlightRanges },
      revealLine: highlightRanges[0]?.start,
    });
  };

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
