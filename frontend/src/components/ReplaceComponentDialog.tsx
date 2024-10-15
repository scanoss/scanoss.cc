import { Button } from './ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from './ui/dialog';

interface ReplaceComponentDialogProps {
  handleCancel: () => void;
  handleConfirm: () => void;
  onOpenChange: (open: boolean) => void;
  open: boolean;
}

export default function ReplaceComponentDialog({
  handleCancel,
  handleConfirm,
  onOpenChange,
  open,
}: ReplaceComponentDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Replace Component</DialogTitle>
          <DialogDescription>
            Select a component to replace the current selected components
          </DialogDescription>
        </DialogHeader>

        <DialogFooter>
          <Button variant="ghost" onClick={handleCancel}>
            Cancel
          </Button>
          <Button onClick={handleConfirm}>Confirm</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
