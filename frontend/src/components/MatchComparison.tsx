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

import { useQuery } from '@tanstack/react-query';
import { memo } from 'react';
import { v4 as uuidv4 } from 'uuid';

import CodeViewer from '@/components/CodeViewer';
import useSelectedResult from '@/hooks/useSelectedResult';
import { getFileName } from '@/lib/utils';

import { GetComponentByPath } from '../../wailsjs/go/service/ComponentServiceImpl';
import { GetLocalFile, GetRemoteFile } from '../../wailsjs/go/service/FileServiceImpl';
import EditorToolbar from './EditorToolbar';
import FileInfoCard from './FileInfoCard';
import Header from './Header';
import MatchInfoCard from './MatchInfoCard';

const MemoizedCodeViewer = memo(CodeViewer);

export default function MatchComparison() {
  const selectedResult = useSelectedResult();

  const {
    data: localFileContent,
    isFetching: isLoadingLocalFileContent,
    isError: isErrorLocalFileContent,
    error: errorLocalFileContent,
  } = useQuery({
    queryKey: ['localFileContent', selectedResult?.path],
    queryFn: () => GetLocalFile(selectedResult?.path as string),
    enabled: !!selectedResult?.path,
  });

  const {
    data: remoteFileContent,
    isFetching: isLoadingRemoteFileContent,
    isError: isErrorRemoteFileContent,
    error: errorRemoteFileContent,
  } = useQuery({
    queryKey: ['remoteFileContent', selectedResult?.path],
    queryFn: () => GetRemoteFile(selectedResult?.path as string),
    enabled: !!selectedResult?.path,
  });

  const { data: component } = useQuery({
    queryKey: ['component', selectedResult?.path],
    queryFn: () => GetComponentByPath(selectedResult?.path as string),
    enabled: !!selectedResult?.path,
  });

  if (!selectedResult) {
    return null;
  }

  return (
    <div className="flex h-full flex-col">
      <header className="h-[65px] border-b border-b-border px-6">
        <Header />
      </header>
      <main className="flex-1">
        <div className="grid h-full grid-cols-2 grid-rows-[auto_auto_auto_1fr]">
          <div className="col-span-2">
            <MatchInfoCard />
          </div>
          <FileInfoCard title="Local file" subtitle={getFileName(selectedResult?.path as string)} fileType="local" />
          <FileInfoCard title="Remote file" subtitle={component?.file} fileType="remote" />
          <div className="col-span-2">
            <EditorToolbar />
          </div>
          <MemoizedCodeViewer
            content={localFileContent?.content}
            isError={isErrorLocalFileContent}
            error={errorLocalFileContent}
            isLoading={isLoadingLocalFileContent}
            language={localFileContent?.language}
            highlightLines={component?.lines}
            editorType="local"
            editorId={uuidv4()}
          />
          <MemoizedCodeViewer
            content={remoteFileContent?.content}
            isError={isErrorRemoteFileContent}
            error={errorRemoteFileContent}
            isLoading={isLoadingRemoteFileContent}
            language={remoteFileContent?.language}
            highlightLines={component?.oss_lines}
            editorType="remote"
            editorId={uuidv4()}
          />
        </div>
      </main>
    </div>
  );
}
