import clsx from 'clsx';
import { File, Github } from 'lucide-react';

import useSelectedResult from '@/hooks/useSelectedResult';
import { FilterAction } from '@/modules/components/domain';
import { stateInfoPresentation } from '@/modules/results/domain';

interface FileInfoCardProps {
  title: string;
  subtitle: string | undefined;
  fileType: 'local' | 'remote';
}

export default function FileInfoCard({ title, subtitle, fileType }: FileInfoCardProps) {
  const result = useSelectedResult();

  const filterConfig = result?.filter_config;

  const presentation = stateInfoPresentation[filterConfig?.action as FilterAction];

  const shouldShowStateInfo =
    (fileType === 'local' && filterConfig?.type === 'by_file') ||
    (fileType === 'remote' && filterConfig?.type === 'by_purl');

  return (
    <div
      className={clsx(
        'flex justify-between border border-l-0 border-t-0 p-3 text-sm',
        shouldShowStateInfo && presentation?.stateInfoContainerStyles,
        fileType === 'remote' && 'border-r-0'
      )}
    >
      <div>
        <div className="flex items-center gap-1">
          {fileType === 'local' ? <File className="h-3 w-3" /> : <Github className="h-3 w-3" />}
          <span className="font-semibold">{title}</span>
        </div>
        <p className="text-muted-foreground">{subtitle ?? '-'}</p>
      </div>
      {shouldShowStateInfo && (
        <div className="text-xs">
          <p className={presentation?.stateInfoTextStyles}>{presentation?.label}</p>
        </div>
      )}
    </div>
  );
}
