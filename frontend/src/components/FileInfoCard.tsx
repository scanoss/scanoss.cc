// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
    (fileType === 'local' && filterConfig?.type === 'by_file') || (fileType === 'remote' && filterConfig?.type === 'by_purl');

  return (
    <div
      className={clsx(
        'flex justify-between border p-3 text-sm',
        shouldShowStateInfo ? presentation?.stateInfoContainerStyles : 'border-l-0 border-t-0',
        fileType === 'remote' && !shouldShowStateInfo && 'border-r-0'
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
