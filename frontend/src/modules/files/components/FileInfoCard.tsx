import React from 'react';

interface FileInfoCardProps {
  title: string;
  subtitle: string;
}

export default function FileInfoCard({ title, subtitle }: FileInfoCardProps) {
  return (
    <div className="mb-4 flex flex-col rounded-sm border border-border p-3 text-sm">
      <p className="font-semibold">{title}</p>
      <p className="text-muted-foreground">{subtitle}</p>
    </div>
  );
}
