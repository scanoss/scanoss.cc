import { useQuery } from '@tanstack/react-query';
import { v4 as uuidv4 } from 'uuid';

import CodeViewer from '@/components/CodeViewer';
import { getFileName } from '@/lib/utils';
import useLocalFilePath from '@/modules/files/hooks/useLocalFilePath';

import { GetComponentByPath } from '../../wailsjs/go/service/ComponentServiceImpl';
import { GetLocalFile, GetRemoteFile } from '../../wailsjs/go/service/FileServiceImpl';
import FileInfoCard from './FileInfoCard';
import Header from './Header';
import MatchInfoCard from './MatchInfoCard';

export default function MatchComparison() {
  const localFilePath = useLocalFilePath();

  const {
    data: localFileContent,
    isFetching: isLoadingLocalFileContent,
    isError: isErrorLocalFileContent,
  } = useQuery({
    queryKey: ['localFileContent', localFilePath],
    queryFn: () => GetLocalFile(localFilePath),
  });

  const {
    data: remoteFileContent,
    isFetching: isLoadingRemoteFileContent,
    isError: isErrorRemoteFileContent,
  } = useQuery({
    queryKey: ['remoteFileContent', localFilePath],
    queryFn: () => GetRemoteFile(localFilePath),
  });

  const { data: component } = useQuery({
    queryKey: ['component', localFilePath],
    queryFn: () => GetComponentByPath(localFilePath),
  });

  return (
    <div className="flex h-full flex-col">
      <header className="h-[65px] border-b border-b-border px-6">
        <Header />
      </header>
      <main className="flex-1 p-6">
        <div className="grid h-full grid-cols-2 grid-rows-[auto_auto_1fr] gap-4">
          <div className="col-span-2">
            <MatchInfoCard />
          </div>
          <FileInfoCard title="Local file" subtitle={getFileName(localFilePath)} fileType="local" />
          <FileInfoCard title="Remote file" subtitle={component?.file} fileType="remote" />
          <CodeViewer
            content={localFileContent?.content}
            isError={isErrorLocalFileContent}
            isLoading={isLoadingLocalFileContent}
            language={localFileContent?.language}
            highlightLines={component?.lines}
            editorType="local"
            editorId={uuidv4()}
          />
          <CodeViewer
            content={remoteFileContent?.content}
            isError={isErrorRemoteFileContent}
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
