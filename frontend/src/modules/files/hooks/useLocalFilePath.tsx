import { useParams } from 'react-router-dom';

import { decodeFilePath } from '@/lib/utils';

export default function useLocalFilePath(): string {
  const { filePath } = useParams();
  const localFilePath = decodeFilePath(filePath ?? '');

  return localFilePath;
}
