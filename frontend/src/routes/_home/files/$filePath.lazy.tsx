import { createLazyFileRoute } from '@tanstack/react-router';
import React from 'react';

export const Route = createLazyFileRoute('/_home/files/$filePath')({
  component: Matches,
});

function Matches() {
  const { filePath } = Route.useParams();

  return <div>Hello world from /matches/{filePath}</div>;
}
