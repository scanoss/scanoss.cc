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
import { FileSearch } from 'lucide-react';
import { memo } from 'react';
import { v4 as uuidv4 } from 'uuid';

import CodeViewer from '@/components/CodeViewer';
import useSelectedResult from '@/hooks/useSelectedResult';
import { getFileName } from '@/lib/utils';

import { entities } from '../../wailsjs/go/models';
import { GetComponentByPath } from '../../wailsjs/go/service/ComponentServiceImpl';
import { GetLocalFile, GetRemoteFile } from '../../wailsjs/go/service/FileServiceImpl';
import { EventsEmit } from '../../wailsjs/runtime/runtime';
import EditorToolbar from './EditorToolbar';
import EmptyState from './EmptyState';
import FileInfoCard from './FileInfoCard';
import Header from './Header';
import Loading from './Loading';
import MatchInfoCard from './MatchInfoCard';

const MemoizedCodeViewer = memo(CodeViewer);

export default function MatchComparison() {
  const selectedResult = useSelectedResult();

  const {
    data: localFileContent,
    isError: isErrorLocalFileContent,
    error: errorLocalFileContent,
    isLoading: isLoadingLocalFileContent,
  } = useQuery({
    queryKey: ['localFileContent', selectedResult?.path],
    queryFn: () => GetLocalFile(selectedResult?.path as string),
    enabled: !!selectedResult?.path,
  });

  const {
    data: remoteFileContent,
    isError: isErrorRemoteFileContent,
    isLoading: isLoadingRemoteFileContent,
    error: errorRemoteFileContent,
  } = useQuery({
    queryKey: ['remoteFileContent', selectedResult?.path],
    queryFn: () => GetRemoteFile(selectedResult?.path as string),
    enabled: !!selectedResult?.path,
  });

  const { data: component, isLoading: isLoadingComponent } = useQuery({
    queryKey: ['component', selectedResult?.path],
    queryFn: () => GetComponentByPath(selectedResult?.path as string),
    enabled: !!selectedResult?.path,
  });

  if (!selectedResult) {
    return (
      <EmptyState
        icon={<FileSearch className="h-12 w-12" />}
        title="No file selected"
        subtitle="Select a file from the results list to view the comparison, or run a new scan to compare files."
        action={{
          label: 'Run a new scan',
          onClick: () => {
            EventsEmit(entities.Action.ScanWithOptions);
          },
        }}
      />
    );
  }

  const isLoadingLocalFile = isLoadingLocalFileContent || isLoadingComponent;
  const isLoadingRemoteFile = isLoadingRemoteFileContent || isLoadingComponent;

  return (
    <div className="flex h-full flex-col overflow-hidden">
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
          <div className="flex flex-1 items-center justify-center text-muted-foreground">
            {isLoadingLocalFile ? (
              <Loading text="Loading local file..." />
            ) : (
              <MemoizedCodeViewer
                content={localFileContent?.content}
                isError={isErrorLocalFileContent}
                error={errorLocalFileContent}
                language={localFileContent?.language}
                highlightLines={component?.lines}
                editorId={uuidv4()}
              />
            )}
          </div>
          <div className="flex flex-1 items-center justify-center text-muted-foreground">
            {isLoadingRemoteFile ? (
              <Loading text="Loading remote file..." />
            ) : (
              <MemoizedCodeViewer
                content={remoteFileContent?.content}
                isError={isErrorRemoteFileContent}
                error={errorRemoteFileContent}
                language={remoteFileContent?.language}
                highlightLines={component?.oss_lines}
                editorId={uuidv4()}
              />
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
