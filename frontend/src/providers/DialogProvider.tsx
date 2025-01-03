// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import React, { createContext, useCallback, useContext, useState } from 'react';

import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';

type DialogOptions = {
  title: string;
  description?: string;
  content: React.ReactNode;
  showCancel?: boolean;
  confirmText?: string;
  cancelText?: string;
};

type DialogContextType = {
  showDialog: <T>(options: DialogOptions) => Promise<T | undefined>;
};

const DialogContext = createContext<DialogContextType | null>(null);

export function DialogProvider({ children }: { children: React.ReactNode }) {
  const [isOpen, setIsOpen] = useState(false);
  const [dialogContent, setDialogContent] = useState<React.ReactNode | null>(null);
  const [resolveRef, setResolveRef] = useState<((value: unknown) => void) | null>(null);

  const [options, setOptions] = useState<DialogOptions>({
    title: '',
    description: '',
    content: null,
    showCancel: true,
    confirmText: 'Confirm',
    cancelText: 'Cancel',
  });

  const showDialog = useCallback(<T,>(dialogOptions: DialogOptions): Promise<T | undefined> => {
    return new Promise((resolve) => {
      setOptions({
        showCancel: true,
        confirmText: 'Confirm',
        cancelText: 'Cancel',
        ...dialogOptions,
      });
      setDialogContent(dialogOptions.content);
      setResolveRef(() => resolve);
      setIsOpen(true);
    });
  }, []);

  const handleClose = useCallback(() => {
    setIsOpen(false);
    resolveRef?.(undefined);
    setResolveRef(null);
  }, [resolveRef]);

  const handleConfirm = useCallback(
    (value: unknown) => {
      setIsOpen(false);
      resolveRef?.(value);
      setResolveRef(null);
    },
    [resolveRef]
  );

  const ref = useKeyboardShortcut(
    KEYBOARD_SHORTCUTS.confirm.keys,
    handleConfirm,
    {
      enableOnFormTags: true,
    },
    [isOpen]
  );

  return (
    <DialogContext.Provider value={{ showDialog }}>
      {children}
      <Dialog open={isOpen} onOpenChange={(open) => !open && handleClose()}>
        <DialogContent ref={ref} tabIndex={-1}>
          <DialogHeader>
            <DialogTitle>{options.title}</DialogTitle>
            <DialogDescription>{options.description}</DialogDescription>
          </DialogHeader>

          {dialogContent}

          <DialogFooter>
            <Button variant="ghost" onClick={handleClose}>
              {options.cancelText}
            </Button>
            <Button onClick={handleConfirm}>
              {options.confirmText} <span className="ml-2 rounded-sm bg-card p-1 text-[8px] leading-none">⌘ + Enter</span>
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </DialogContext.Provider>
  );
}

export function useDialog() {
  const context = useContext(DialogContext);
  if (!context) {
    throw new Error('useDialog must be used within a DialogProvider');
  }
  return context;
}
