import { useQuery } from '@tanstack/react-query';

import CodeViewer from '@/components/CodeViewer';
import { getFileName } from '@/lib/utils';
import ResultService from '@/modules/results/infra/service';

import useLocalFilePath from '../hooks/useLocalFilePath';
import FileService from '../infra/service';
import FileActionsMenu from './FileActionsMenu';
import FileInfoCard from './FileInfoCard';
import MatchInfoCard from './MatchInfoCard';

export default function MatchComparison() {
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

  const { data: component } = useQuery({
    queryKey: ['component', localFilePath],
    queryFn: () => ResultService.getComponent(localFilePath),
  });

  return (
    <div className="flex h-full flex-col">
      <header className="h-[65px] border-b border-b-border px-6">
        <FileActionsMenu />
      </header>
      <main className="flex-1 p-6">
        <div className="grid h-full grid-cols-2 grid-rows-[auto_auto_1fr] gap-4">
          <div className="col-span-2">
            <MatchInfoCard />
          </div>
          <FileInfoCard
            title="Local file"
            subtitle={getFileName(localFilePath)}
            fileType="local"
          />
          <FileInfoCard
            title="Remote file"
            subtitle={component?.file}
            fileType="remote"
          />
          <CodeViewer
            content={localFileContent?.content}
            isError={isErrorLocalFileContent}
            isLoading={isLoadingLocalFileContent}
            language={localFileContent?.language}
            highlightLines={component?.lines}
            editorType="local"
          />
          <CodeViewer
            content={remoteFileContent?.content}
            isError={isErrorRemoteFileContent}
            isLoading={isLoadingRemoteFileContent}
            language={remoteFileContent?.language}
            highlightLines={component?.oss_lines}
            editorType="remote"
          />
        </div>
      </main>
    </div>
  );
}
