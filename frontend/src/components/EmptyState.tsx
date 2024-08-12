import React from 'react';

interface EmptyStateProps {
  title: string;
  subtitle: string;
}

export default function EmptyState({ title, subtitle }: EmptyStateProps) {
  return (
    <div className="backdrop-blur-sm border-border border bg-black/20 p-4 rounded-sm">
      <div className="flex flex-col gap-3">
        <h3 className="font-bold text-xl">{title}</h3>
        <p className="text-muted-foreground">{subtitle}</p>
      </div>
    </div>
  );
}
