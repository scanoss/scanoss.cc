import { useQuery } from '@tanstack/react-query';

import CodeViewer from '@/components/CodeViewer';
import { Component } from '@/modules/results/domain';

import useLocalFilePath from '../hooks/useLocalFilePath';
import FileService from '../infra/service';
import FileActionsMenu from './FileActionsMenu';
import FileInfoCard from './FileInfoCard';
import MatchInfoCard from './MatchInfoCard';

interface WrapperProps {
  children: React.ReactNode;
  component: Component;
}

function Wrapper({ children, component }: WrapperProps) {
  const localFilePath = useLocalFilePath();

  return (
    <div className="flex h-full w-full flex-col">
      <FileActionsMenu component={component} />

      <div className="flex flex-1 flex-col gap-4 p-6">
        <MatchInfoCard component={component} />
        <div className="flex gap-4">
          <div className="flex-1">
            <FileInfoCard title="Local file" subtitle={localFilePath} />
          </div>
          <div className="flex-1">
            <FileInfoCard
              title="Remote file"
              subtitle={component.file}
              noMatch={false}
            />
          </div>
        </div>
        {children}
      </div>
    </div>
  );
}

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
    <Wrapper component={component}>
      <div className="flex flex-1 gap-4">
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
    </Wrapper>
  );
}
