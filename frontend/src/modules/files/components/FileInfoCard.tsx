import clsx from 'clsx';
import React from 'react';

interface FileInfoCardProps {
  title: string;
  subtitle: string | undefined;
}

export default function FileInfoCard({ title, subtitle }: FileInfoCardProps) {
  return (
    <div
      className={clsx(
        'flex flex-col rounded-sm border border-border bg-card p-3 text-sm'
      )}
    >
      <p className="font-semibold">{title}</p>
      <p className="text-muted-foreground">{subtitle}</p>
    </div>
  );
}
