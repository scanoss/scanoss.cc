import { useEffect, useState } from 'react';

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
  const [inputValue, setInputValue] = useState('');

  useEffect(() => {
    return () => {
      console.log('destroyed...');
    };
  }, []);

  if (!isPrompting || !options) return null;

  const { title, description, cancelText, confirmText } = options;

  const handleConfirm = () => {
    confirm(inputValue);
  };

  console.log('PromptDialog', options);
  console.log('inputValue', inputValue);

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
          <Button variant="ghost" onClick={cancel}>
            {cancelText ?? 'Cancel'}
          </Button>
          <Button onClick={handleConfirm}>{confirmText ?? 'Confirm'}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
