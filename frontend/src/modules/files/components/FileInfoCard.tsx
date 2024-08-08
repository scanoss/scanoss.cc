import React from 'react';

interface FileInfoCardProps {
  title: string;
  subtitle: string;
}

export default function FileInfoCard({ title, subtitle }: FileInfoCardProps) {
  return (
    <div className="flex flex-col bg-background p-3 rounded-md border mb-4">
      <p className="font-semibold">{title}</p>
      <p className="text-muted-foreground">{subtitle}</p>
    </div>
  );
}
