// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
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
