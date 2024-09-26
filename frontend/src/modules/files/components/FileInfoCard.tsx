import clsx from 'clsx';
import { File, Github } from 'lucide-react';

import { FilterAction } from '@/modules/results/domain';
import { useResults } from '@/modules/results/providers/ResultsProvider';

import useLocalFilePath from '../hooks/useLocalFilePath';

interface FileInfoCardProps {
  title: string;
  subtitle: string | undefined;
  fileType: 'local' | 'remote';
}

export default function FileInfoCard({
  title,
  subtitle,
  fileType,
}: FileInfoCardProps) {
  const { results } = useResults();
  const localFilePath = useLocalFilePath();
  const result = results.find((result) => result.path === localFilePath);

  const filterConfig = result?.filter_config;

  const isResultDismissed = filterConfig?.action === FilterAction.Remove;
  const isResultIncluded = filterConfig?.action === FilterAction.Include;

  const shouldShowStateInfo =
    (fileType === 'local' && filterConfig?.type === 'by_file') ||
    (fileType === 'remote' && filterConfig?.type === 'by_purl');

  return (
    <div
      className={clsx(
        'flex justify-between rounded-sm border border-border bg-card p-3 text-sm',
        shouldShowStateInfo && {
          'border-l-4 border-green-600 border-l-green-600 bg-green-950':
            isResultIncluded,
          'border-l-4 border-red-600 border-l-red-600 bg-red-950':
            isResultDismissed,
        }
      )}
    >
      <div>
        <div className="flex items-center gap-1">
          {fileType === 'local' ? (
            <File className="h-3 w-3" />
          ) : (
            <Github className="h-3 w-3" />
          )}
          <span className="font-semibold">{title}</span>
        </div>
        <p className="text-muted-foreground">{subtitle ?? '-'}</p>
      </div>
      {shouldShowStateInfo && (
        <div>
          {isResultIncluded && (
            <p className="text-xs text-green-600">Included</p>
          )}
          {isResultDismissed && (
            <p className="text-xs text-red-500">Dismissed</p>
          )}
        </div>
      )}
    </div>
  );
}
