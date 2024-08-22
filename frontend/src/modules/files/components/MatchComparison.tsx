import { useQuery } from '@tanstack/react-query';
import React from 'react';

import CodeViewer from '@/components/CodeViewer';
import EmptyState from '@/components/EmptyState';
import { Component, MatchType } from '@/modules/results/domain';

import FileService from '../infra/service';
import FileActionsMenu from './FileActionsMenu';
import FileInfoCard from './FileInfoCard';
import MatchInfoCard from './MatchInfoCard';

interface WrapperProps {
  children: React.ReactNode;
  component: Component;
  localFilePath: string;
}

function Wrapper({ children, component, localFilePath }: WrapperProps) {
  return (
    <div className="flex h-full w-full flex-col">
      <FileActionsMenu component={component} />

      <div className="flex flex-1 flex-col gap-4 p-6">
        <MatchInfoCard component={component} />
        <div className="flex gap-4">
          <div className="flex-1">
            <FileInfoCard title="Source file" subtitle={localFilePath} />
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
  localFilePath: string;
  component: Component;
}

export default function MatchComparison({
  localFilePath,
  component,
}: MatchComparisonProps) {
  const matchType = component.id;

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

  if (matchType === MatchType.File) {
    return (
      <Wrapper localFilePath={localFilePath} component={component}>
        <CodeViewer
          content={localFileContent?.content}
          isError={isErrorLocalFileContent}
          isLoading={isLoadingLocalFileContent}
          language={localFileContent?.language}
        />
      </Wrapper>
    );
  }

  if (matchType === MatchType.Snippet) {
    return (
      <Wrapper localFilePath={localFilePath} component={component}>
        <div className="flex flex-1 gap-4">
          <CodeViewer
            content={localFileContent?.content}
            isError={isErrorLocalFileContent}
            isLoading={isLoadingLocalFileContent}
            language={localFileContent?.language}
            highlightLines={component.lines}
          />
          <CodeViewer
            content={remoteFileContent?.content}
            isError={isErrorRemoteFileContent}
            isLoading={isLoadingRemoteFileContent}
            language={remoteFileContent?.language}
            highlightLines={component.oss_lines}
          />
        </div>
      </Wrapper>
    );
  }

  return (
    <div className="flex h-full w-full items-center justify-center">
      <EmptyState
        title="No matches"
        subtitle="There are no matches found for this file, please select another one."
      />
    </div>
  );
}
