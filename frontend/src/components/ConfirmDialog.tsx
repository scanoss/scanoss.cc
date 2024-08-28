import { useState } from 'react';

import { useConfirm } from '@/hooks/useConfirm';

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

export default function ConfirmDialog() {
  const { isAsking, message, deny, confirm, onPersistDecision } = useConfirm();
  const [persistDecision, setPersistDecision] = useState(false);

  const handleConfirm = () => {
    if (persistDecision) {
      onPersistDecision?.();
    }
    confirm();
  };

  return (
    <AlertDialog open={isAsking} onOpenChange={deny}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Confirm</AlertDialogTitle>
          <AlertDialogDescription>{message}</AlertDialogDescription>
          <div className="items-top !my-5 flex space-x-2">
            <Checkbox
              id="persistDecision"
              checked={persistDecision}
              onCheckedChange={(checked) =>
                typeof checked === 'boolean'
                  ? setPersistDecision(checked)
                  : undefined
              }
            />
            <div className="grid gap-1.5 leading-none">
              <label
                htmlFor="persistDecision"
                className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              >
                Don&apos;t ask me again
              </label>
              <p className="text-sm text-muted-foreground">
                This action will be performed without asking for confirmation in
                the future
              </p>
            </div>
          </div>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={deny}>Cancel</AlertDialogCancel>
          <AlertDialogAction onClick={handleConfirm}>Confirm</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}