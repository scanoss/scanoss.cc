import { createFileRoute } from '@tanstack/react-router';
import React from 'react';

export const Route = createFileRoute('/main')({
  component: () => <div>Hello /main!</div>,
});
