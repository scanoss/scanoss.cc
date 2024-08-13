import React from 'react';

interface EmptyStateProps {
  title: string;
  subtitle: string;
}

export default function EmptyState({ title, subtitle }: EmptyStateProps) {
  return (
    <div className="rounded-sm border border-border bg-black/20 p-4 backdrop-blur-sm">
      <div className="flex flex-col gap-3">
        <h3 className="text-xl font-bold">{title}</h3>
        <p className="text-muted-foreground">{subtitle}</p>
      </div>
    </div>
  );
}
