import { useState } from 'react';

import { useConfirm } from '@/hooks/useConfirm';
import useEnvironment from '@/hooks/useEnvironment';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from './ui/alert-dialog';
import { Checkbox } from './ui/checkbox';

interface ConfirmDialogProps {
  showPersistDecision?: boolean;
}

export default function ConfirmDialog({ showPersistDecision = false }: ConfirmDialogProps) {
  const { modifierKey } = useEnvironment();
  const { isAsking, message, deny, confirm, onPersistDecision } = useConfirm();
  const [persistDecision, setPersistDecision] = useState(false);

  const handleConfirm = () => {
    if (persistDecision) {
      onPersistDecision?.();
    }
    confirm();
  };

  useKeyboardShortcut([modifierKey.keyCode, 'enter'], handleConfirm);

  return (
    <AlertDialog open={isAsking} onOpenChange={deny}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Are you sure?</AlertDialogTitle>
          <AlertDialogDescription asChild>{message}</AlertDialogDescription>
          {showPersistDecision && (
            <div className="items-top !my-5 flex space-x-2">
              <Checkbox
                id="persistDecision"
                checked={persistDecision}
                onCheckedChange={(checked) => (typeof checked === 'boolean' ? setPersistDecision(checked) : undefined)}
              />
              <div className="grid gap-1.5 leading-none">
                <label
                  htmlFor="persistDecision"
                  className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  Don&apos;t ask me again
                </label>
                <p className="text-sm text-muted-foreground">
                  This action will be performed without asking for confirmation in the future
                </p>
              </div>
            </div>
          )}
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={deny}>Cancel</AlertDialogCancel>
          <AlertDialogAction onClick={handleConfirm}>
            Confirm <span className="ml-2 rounded-sm bg-card p-1 text-[8px] leading-none">⌘ + Enter</span>
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
