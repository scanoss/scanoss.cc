import { useQuery } from '@tanstack/react-query';

import CodeViewer from '@/components/CodeViewer';
import { Component } from '@/modules/results/domain';

import useLocalFilePath from '../hooks/useLocalFilePath';
import FileService from '../infra/service';
import FileActionsMenu from './FileActionsMenu';
import FileInfoCard from './FileInfoCard';
import MatchInfoCard from './MatchInfoCard';

interface MatchComparisonProps {
  component: Component;
}

export default function MatchComparison({ component }: MatchComparisonProps) {
  const localFilePath = useLocalFilePath();

  const {
    data: localFileContent,
    isFetching: isLoadingLocalFileContent,
    isError: isErrorLocalFileContent,
  } = useQuery({
    queryKey: ['localFileContent', localFilePath],
    queryFn: () => FileService.getLocalFileContent(localFilePath),
  });

  const {
    data: remoteFileContent,
    isFetching: isLoadingRemoteFileContent,
    isError: isErrorRemoteFileContent,
  } = useQuery({
    queryKey: ['remoteFileContent', localFilePath],
    queryFn: () => FileService.getRemoteFileContent(localFilePath),
  });

  return (
    <div className="flex h-full flex-col">
      <header className="h-[65px] border-b border-b-border px-6">
        <FileActionsMenu component={component} />
      </header>
      <main className="flex-1 p-6">
        <div className="grid h-full grid-cols-2 grid-rows-[auto_auto_1fr] gap-4">
          <div className="col-span-2">
            <MatchInfoCard component={component} />
          </div>
          <FileInfoCard title="Local file" subtitle={localFilePath} />
          <FileInfoCard title="Remote file" subtitle={component.file} />
          <CodeViewer
            content={localFileContent?.content}
            isError={isErrorLocalFileContent}
            isLoading={isLoadingLocalFileContent}
            language={localFileContent?.language}
            highlightLines={component.lines}
            editorType="local"
          />
          <CodeViewer
            content={remoteFileContent?.content}
            isError={isErrorRemoteFileContent}
            isLoading={isLoadingRemoteFileContent}
            language={remoteFileContent?.language}
            highlightLines={component.oss_lines}
            editorType="remote"
          />
        </div>
      </main>
    </div>
  );
}
