import { useConfirm } from '@/hooks/useConfirm';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';

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

export default function ConfirmDialog() {
  const { isAsking, message, deny, confirm } = useConfirm();

  const ref = useKeyboardShortcut(
    KEYBOARD_SHORTCUTS.confirm.keys,
    confirm,
    {
      enableOnFormTags: true,
    },
    [isAsking]
  );

  return (
    <AlertDialog open={isAsking} onOpenChange={deny}>
      <AlertDialogContent ref={ref} tabIndex={-1}>
        <AlertDialogHeader>
          <AlertDialogTitle>Are you sure?</AlertDialogTitle>
          <AlertDialogDescription asChild>{message}</AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={deny}>Cancel</AlertDialogCancel>
          <AlertDialogAction onClick={confirm}>
            Confirm <span className="ml-2 rounded-sm bg-card p-1 text-[8px] leading-none">⌘ + Enter</span>
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
