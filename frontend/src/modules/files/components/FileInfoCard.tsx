import React from 'react';

interface FileInfoCardProps {
  title: string;
  subtitle: string;
}

export default function FileInfoCard({ title, subtitle }: FileInfoCardProps) {
  return (
    <div className="flex flex-col p-3 rounded-sm border border-border mb-4 text-sm">
      <p className="font-semibold">{title}</p>
      <p className="text-muted-foreground">{subtitle}</p>
    </div>
  );
}
