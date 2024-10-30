import { useEffect, useState } from 'react';

import useEnvironment from '@/hooks/useEnvironment';
import { useInputPrompt } from '@/hooks/useInputPrompt';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';

import { Button } from './ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from './ui/dialog';
import { Textarea } from './ui/textarea';

export default function InputPromptDialog() {
  const { modifierKey } = useEnvironment();
  const { isPrompting, options, confirm, cancel } = useInputPrompt();
  const [inputValue, setInputValue] = useState('');

  useEffect(() => {
    if (!isPrompting) setInputValue('');
  }, [isPrompting]);

  const ref = useKeyboardShortcut(
    KEYBOARD_SHORTCUTS.confirm.keys,
    () => confirm(inputValue),
    {
      enableOnFormTags: true,
    },
    [inputValue]
  );

  return (
    <Dialog open={isPrompting} onOpenChange={cancel}>
      <DialogContent ref={ref} tabIndex={-1}>
        <DialogHeader>
          <DialogTitle>{options?.title}</DialogTitle>
          <DialogDescription>
            {options?.description && <p className="mb-1">{options?.description}</p>}
          </DialogDescription>
        </DialogHeader>

        {options?.input && 'type' in options.input && options.input.type === 'textarea' && (
          <div>
            <Textarea value={inputValue} onChange={(e) => setInputValue(e.target.value)} />
            <p className="mt-2 text-xs text-muted-foreground">
              Use <span className="rounded bg-primary px-1.5">{modifierKey.label} + return</span> to confirm
            </p>
          </div>
        )}

        {options?.input && 'component' in options.input && options.input.component}

        <DialogFooter>
          <Button variant="ghost" onClick={cancel}>
            {options?.cancelText ?? 'Cancel'}
          </Button>
          <Button onClick={() => confirm(inputValue)}>
            {options?.confirmText ?? 'Confirm'}{' '}
            <span className="ml-2 rounded-sm bg-card p-1 text-[8px] leading-none">⌘ + Enter</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
