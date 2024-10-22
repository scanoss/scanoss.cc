import { entities } from 'wailsjs/go/models';

import { encodeFilePath } from '@/lib/utils';

interface HighlightLineRange {
  start: number;
  end: number;
}

export const getHighlightLineRanges = (lines: string): HighlightLineRange[] => {
  const ranges = lines.split(',').map((range) => range.split('-'));

  return ranges.map((range) => ({
    start: Number(range[0]),
    end: Number(range[1]),
  }));
};

export const getNextPendingResultPathRoute = (
  pendingResults: entities.ResultDTO[]
): string | null => {
  const firstAvailablePendingResult = pendingResults[0];
  if (!firstAvailablePendingResult) return null;

  return `/files/${encodeFilePath(firstAvailablePendingResult.path)}`;
};
