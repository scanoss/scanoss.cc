import { useQuery } from '@tanstack/react-query';
import React from 'react';
import { useParams } from 'react-router-dom';

import CodeViewer from '@/components/CodeViewer';
import { decodeFilePath } from '@/lib/utils';
import FileInfoCard from '@/modules/files/components/FileInfoCard';
import MatchInfoCard from '@/modules/files/components/MatchInfoCard';
import MatchSkeleton from '@/modules/files/components/MatchSkeleton';
import FileService from '@/modules/files/infra/service';

export default function FileMatchRoute() {
  const { filePath } = useParams();

  const decodedFilePath = decodeFilePath(filePath ?? '');

  const { data: localFile } = useQuery({
    queryKey: ['localFile', decodedFilePath],
    queryFn: () => FileService.getLocalFileContent(decodedFilePath),
    enabled: !!decodedFilePath,
  });

  if (!localFile) {
    return <MatchSkeleton />;
  }

  return (
    <div className="w-full h-full flex flex-col gap-2 bg-red-50 overflow-auto">
      <MatchInfoCard />
      <div className="flex flex-1 gap-8 h-full">
        <div className="flex flex-col flex-1">
          <FileInfoCard title="Source file" subtitle={localFile?.path} />
          <CodeViewer
            language={localFile?.language}
            content={localFile?.content}
          />
        </div>
        <div className="flex flex-col flex-1">
          <FileInfoCard title="Component file" subtitle={localFile?.path} />
          <CodeViewer
            language={localFile?.language}
            content={localFile?.content}
          />
        </div>
      </div>
    </div>
  );
}
