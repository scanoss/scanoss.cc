import { useQuery } from '@tanstack/react-query';
import React from 'react';
import { useParams } from 'react-router-dom';

import { decodeFilePath } from '@/lib/utils';
import MatchComparison from '@/modules/files/components/MatchComparison';
import ResultService from '@/modules/results/infra/service';

export default function FileMatchRoute() {
  const { filePath } = useParams();
  const localFilePath = decodeFilePath(filePath ?? '');

  const { data: component, error } = useQuery({
    queryKey: ['component', localFilePath],
    queryFn: () => ResultService.getComponent(localFilePath),
  });

  // TODO: Add proper loading states
  if (!component) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  return (
    <MatchComparison localFilePath={localFilePath} component={component} />
  );
}
