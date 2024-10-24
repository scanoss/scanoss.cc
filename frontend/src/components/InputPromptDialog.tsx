import { useState } from 'react';

import useEnvironment from '@/hooks/useEnvironment';
import { useInputPrompt } from '@/hooks/useInputPrompt';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';

import { Button } from './ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from './ui/dialog';
import { Textarea } from './ui/textarea';

export default function InputPromptDialog() {
  const { modifierKey } = useEnvironment();
  const { isPrompting, options, confirm, cancel } = useInputPrompt();
  const [inputValue, setInputValue] = useState(options?.input.defaultValue ?? '');

  const handleConfirm = () => {
    confirm(inputValue);
    setInputValue('');
  };

  const handleCancel = () => {
    cancel();
    setInputValue('');
  };

  useKeyboardShortcut([modifierKey.keyCode, 'Enter'], handleConfirm, {
    preventRegistering: !inputValue,
  });

  if (!isPrompting || !options) return null;

  const { title, description, cancelText, confirmText } = options;

  return (
    <Dialog open={isPrompting} onOpenChange={cancel}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description && <p className="mb-1">{description}</p>}</DialogDescription>
        </DialogHeader>

        {options.input.type === 'textarea' && (
          <div>
            <Textarea value={inputValue} onChange={(e) => setInputValue(e.target.value)} />
            <p className="mt-2 text-xs text-muted-foreground">
              Use <span className="rounded bg-primary px-1.5">{modifierKey.label} + return</span> to confirm
            </p>
          </div>
        )}

        <DialogFooter>
          <Button variant="ghost" onClick={handleCancel}>
            {cancelText ?? 'Cancel'}
          </Button>
          <Button onClick={handleConfirm} disabled={!inputValue}>
            {confirmText ?? 'Confirm'}{' '}
            <span className="ml-2 rounded-sm bg-card p-1 text-[8px] leading-none">⌘ + Enter</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
