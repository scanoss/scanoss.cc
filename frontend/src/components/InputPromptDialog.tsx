import { useState } from 'react';

import { useInputPrompt } from '@/hooks/useInputPrompt';

import { Button } from './ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from './ui/dialog';
import { Textarea } from './ui/textarea';

export default function InputPromptDialog() {
  const { isPrompting, options, confirm, cancel } = useInputPrompt();
  const [inputValue, setInputValue] = useState(
    options?.input.defaultValue ?? ''
  );

  if (!isPrompting || !options) return null;

  const { title, description, cancelText, confirmText } = options;

  const handleConfirm = () => {
    confirm(inputValue);
    setInputValue('');
  };

  const handleCancel = () => {
    cancel();
    setInputValue('');
  };

  return (
    <Dialog open={isPrompting} onOpenChange={cancel}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>

        {options.input.type === 'textarea' && (
          <Textarea
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
          />
        )}

        <DialogFooter>
          <Button variant="ghost" onClick={handleCancel}>
            {cancelText ?? 'Cancel'}
          </Button>
          <Button onClick={handleConfirm}>{confirmText ?? 'Confirm'}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
