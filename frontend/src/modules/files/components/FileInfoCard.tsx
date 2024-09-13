import clsx from 'clsx';
import React from 'react';

interface FileInfoCardProps {
  title: string;
  subtitle: string | undefined;
  noMatch?: boolean;
}

export default function FileInfoCard({
  title,
  subtitle,
  noMatch,
}: FileInfoCardProps) {
  return (
    <div
      className={clsx(
        'flex flex-col rounded-sm border border-border bg-card p-3 text-sm',
        noMatch && 'border-dashed'
      )}
    >
      <p className="font-semibold">{noMatch ? 'No match found' : title}</p>
      <p className="text-muted-foreground">
        {noMatch ? "This file doesn't have a match" : subtitle}
      </p>
    </div>
  );
}
