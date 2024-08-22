import { decodeFilePath } from '@/lib/utils';
import { useParams } from 'react-router-dom';

export default function useLocalFilePath(): string {
  const { filePath } = useParams();
  const localFilePath = decodeFilePath(filePath ?? '');

  return localFilePath;
}
