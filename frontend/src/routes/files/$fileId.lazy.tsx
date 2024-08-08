import React from 'react';
import { createLazyFileRoute } from '@tanstack/react-router';

export const Route = createLazyFileRoute('/files/$fileId')({
  component: Matches,
});

function Matches() {
  const { fileId } = Route.useParams();

  return <div>Hello world from /matches/{fileId}</div>;
}
