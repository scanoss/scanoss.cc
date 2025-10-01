import clsx from 'clsx';
import { useEffect, useRef, useState } from 'react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from './table';

interface Column<T> {
  key?: keyof T;
  label: string;
  render?: (row: T, rowIndex: number) => React.ReactNode;
  className?: string;
}

interface KeyboardTableProps<T> {
  columns: Column<T>[];
  rows: T[];
  onRowSelect: (row: T) => void;
  getRowId?: (row: T, index: number) => string | number;
  rowClassName?: (row: T, index: number, focused: boolean) => string;
  initialFocusedIndex?: number;
}

export function KeyboardTable<T>({ columns, rows, onRowSelect, getRowId, rowClassName, initialFocusedIndex = 0 }: KeyboardTableProps<T>) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [focusedIndex, setFocusedIndex] = useState<number>(initialFocusedIndex);

  const move = (next: number) => {
    const count = containerRef.current?.querySelectorAll<HTMLElement>('[data-row-index]').length ?? 0;
    if (!count) return;
    const last = count - 1;
    setFocusedIndex(Math.max(0, Math.min(last, next)));
  };

  const onKeyDown: React.KeyboardEventHandler<HTMLDivElement> = (e) => {
    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        e.stopPropagation();
        move(focusedIndex + 1);
        break;
      case 'ArrowUp':
        e.preventDefault();
        e.stopPropagation();
        move(focusedIndex - 1);
        break;
      case 'Home':
        e.preventDefault();
        e.stopPropagation();
        move(0);
        break;
      case 'End':
        e.preventDefault();
        e.stopPropagation();
        move(99999);
        break;
      case 'PageDown':
        e.preventDefault();
        e.stopPropagation();
        move(focusedIndex + 10);
        break;
      case 'PageUp':
        e.preventDefault();
        e.stopPropagation();
        move(focusedIndex - 10);
        break;
      case 'Enter': {
        e.preventDefault();
        e.stopPropagation();
        if (focusedIndex >= 0 && focusedIndex < rows.length) {
          onRowSelect(rows[focusedIndex]);
        }
        break;
      }
    }
  };

  useEffect(() => {
    const element = containerRef.current?.querySelector<HTMLElement>(`[data-row-index="${focusedIndex}"]`);

    element?.scrollIntoView({ block: 'nearest' });
  }, [focusedIndex]);

  return (
    <div ref={containerRef} tabIndex={0} role="grid" aria-activedescendant={`row-${focusedIndex}`} onKeyDown={onKeyDown}>
      <Table>
        <TableHeader>
          <TableRow>
            {columns.map((column, i) => (
              <TableHead key={i}>{column.label}</TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.map((row, i) => {
            const isFocused = focusedIndex === i;
            const defaultClassName = clsx('cursor-pointer', isFocused && 'rounded bg-muted/60 ring-2 ring-ring');
            const customClassName = rowClassName ? rowClassName(row, i, isFocused) : defaultClassName;

            return (
              <TableRow
                key={getRowId ? getRowId(row, i) : i}
                id={`row-${i}`}
                role="row"
                data-row-index={i}
                onClick={() => onRowSelect(row)}
                className={customClassName}
                onMouseEnter={() => setFocusedIndex(i)}
              >
                {columns.map((column, colIndex) => (
                  <TableCell key={colIndex} className={column.className}>
                    {column.render ? column.render(row, i) : column.key ? String(row[column.key] ?? '') : ''}
                  </TableCell>
                ))}
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </div>
  );
}
