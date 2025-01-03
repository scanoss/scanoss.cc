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

import { useQuery } from '@tanstack/react-query';
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
          <CodeViewer
            content={localFileContent?.content}
            isError={isErrorLocalFileContent}
            error={errorLocalFileContent}
            isLoading={isLoadingLocalFileContent}
            language={localFileContent?.language}
            highlightLines={component?.lines}
            editorType="local"
            editorId={uuidv4()}
          />
          <CodeViewer
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
