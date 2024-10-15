import { Check, CircleDotDashed } from 'lucide-react';
import { ReactNode } from 'react';

export enum MatchType {
  File = 'file',
  Snippet = 'snippet',
}

interface ResultStatusPresentation {
  badgeStyles: string;
  label: string;
  icon: ReactNode;
}

export const resultStatusPresentation: Record<
  string,
  ResultStatusPresentation
> = {
  pending: {
    badgeStyles:
      'bg-[color:hsl(200,40%,20%)] text-white border-[color:hsl(200,40%,30%)] hover:bg-[color:hsl(200,40%,20%)]',
    label: 'Pending',
    icon: <CircleDotDashed className="h-3 w-3" />,
  },
  completed: {
    badgeStyles:
      'bg-[color:hsl(120,50%,10%)] text-white border-[color:hsl(120,45%,25%)] hover:bg-[color:hsl(120,50%,10%)]',
    label: 'Completed',
    icon: <Check className="h-3 w-3" />,
  },
};

interface MatchTypePresentation {
  background: string;
  foreground: string;
  accent: string;
  muted: string;
  label: string;
}

export const matchTypePresentation: Record<MatchType, MatchTypePresentation> = {
  [MatchType.File]: {
    background: 'bg-file border border-file-border',
    foreground: 'text-file-foreground',
    accent: 'text-file-accent',
    muted: 'text-muted-foreground',
    label: 'File',
  },
  [MatchType.Snippet]: {
    background: 'bg-snippet border border-snippet-border',
    foreground: 'text-snippet-foreground',
    accent: 'text-snippet-accent',
    muted: 'text-muted-foreground',
    label: 'Snippet',
  },
};
