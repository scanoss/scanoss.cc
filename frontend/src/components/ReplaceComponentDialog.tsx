import { useState } from 'react';

import { Button } from './ui/button';
import { Combobox } from './ui/combobox';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from './ui/dialog';
import { Input } from './ui/input';
import { Label } from './ui/label';

export default function ReplaceComponentDialog() {
  const [open, setOpen] = useState(true);

  const handleCancel = () => {
    setOpen(false);
  };

  const handleConfirm = () => {
    setOpen(false);
  };

  const placeholderComponentOptions = [
    {
      label: 'test',
      value: 'test',
    },
    {
      label: 'test2',
      value: 'test2',
    },
  ];

  return (
    <Dialog open={open} onOpenChange={() => setOpen(false)}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Replace</DialogTitle>
          <DialogDescription>
            You can search for an existing component or manually enter a PURL
          </DialogDescription>
        </DialogHeader>

        <div className="grid grid-cols-1 gap-4">
          <Combobox
            options={placeholderComponentOptions}
            placeholder="Search component..."
            emptyText="No components found"
          />
          <div>
            <Label>Purl</Label>
            <Input placeholder="pkg:github/scanoss/scanner.c" />
          </div>
        </div>

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
