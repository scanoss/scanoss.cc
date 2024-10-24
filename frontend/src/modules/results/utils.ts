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
