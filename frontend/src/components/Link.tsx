import clsx from 'clsx';
import { ReactNode } from 'react';

import { BrowserOpenURL } from '../../wailsjs/runtime';

export default function Link({ children, to, className, ...props }: { children: ReactNode; to: string; className?: string }) {
  return (
    <a
      onClick={(e) => {
        e.preventDefault();
        BrowserOpenURL(to);
      }}
      className={clsx(className, 'cursor-pointer text-blue-500 underline')}
      {...props}
    >
      {children}
    </a>
  );
}
