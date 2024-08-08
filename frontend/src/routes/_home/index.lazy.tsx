import { createLazyFileRoute } from '@tanstack/react-router';
import React from 'react';

export const Route = createLazyFileRoute('/_home/')({
  component: Index,
});

function Index() {
  return (
    <div className="flex justify-center items-center w-full h-full text-xl">
      No file selected, please select a file from the sidebar to start.
    </div>
  );
}
