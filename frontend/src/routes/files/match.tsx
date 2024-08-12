import React from 'react';
import { useParams, useSearchParams } from 'react-router-dom';

import CodeViewer from '@/components/CodeViewer';
import EmptyState from '@/components/EmptyState';
import { decodeFilePath } from '@/lib/utils';
import FileInfoCard from '@/modules/files/components/FileInfoCard';
import MatchInfoCard from '@/modules/files/components/MatchInfoCard';
import { MatchType } from '@/modules/results/domain';

const testContent = `
export default function Root() {
  return (
    <div className="flex h-screen w-full overflow-hidden">
      <div className="w-[330px] h-full">
        <Sidebar />
      </div>
      <div className="flex flex-col p-6 flex-1 bg-muted">
        <Outlet />
      </div>
    </div>
  );
}
`;

export default function FileMatchRoute() {
  const { filePath } = useParams();
  const [searchParams] = useSearchParams();
  const matchType = searchParams.get('matchType');

  const decodedFilePath = decodeFilePath(filePath ?? '');

  if (!matchType) {
    return <div>Error loading match</div>;
  }

  if (matchType === MatchType.None) {
    return (
      <div className="w-full h-full flex justify-center items-center">
        <EmptyState
          title="No matches"
          subtitle="There are no matches found for this file, please select another one."
        />
      </div>
    );
  }

  if (matchType === MatchType.File) {
    return (
      <Wrapper localFilePath={decodedFilePath} remoteFilePath={decodedFilePath}>
        <CodeViewer content={testContent} language="javascript" />
      </Wrapper>
    );
  }

  if (matchType === MatchType.Snippet) {
    return (
      <Wrapper localFilePath={decodedFilePath} remoteFilePath={decodedFilePath}>
        <div className="flex flex-1 gap-4">
          <CodeViewer content={testContent} language="javascript" />
          <CodeViewer content={testContent} language="javascript" />
        </div>
      </Wrapper>
    );
  }
}

interface WrapperProps {
  children: React.ReactNode;
  localFilePath: string;
  remoteFilePath: string;
}

function Wrapper({ children, localFilePath, remoteFilePath }: WrapperProps) {
  return (
    <div className="w-full h-full flex flex-col gap-4 overflow-auto">
      <MatchInfoCard />
      <div className="flex flex-1 flex-col h-full">
        <div className="flex gap-4">
          <div className="flex-1">
            <FileInfoCard title="Source file" subtitle={localFilePath} />
          </div>
          <div className="flex-1">
            <FileInfoCard title="Component file" subtitle={remoteFilePath} />
          </div>
        </div>
        {children}
      </div>
    </div>
  );
}
