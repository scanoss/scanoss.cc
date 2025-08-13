// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2025 SCANOSS.COM
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

import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';

interface DialogStateContextType {
  isAnyDialogOpen: boolean;
  isKeyboardNavigationBlocked: boolean;
  registerDialog: (id: string, blocksKeyboardNavigation?: boolean) => void;
  unregisterDialog: (id: string) => void;
}

const DialogStateContext = createContext<DialogStateContextType | undefined>(undefined);

interface DialogInfo {
  blocksKeyboardNavigation: boolean;
}

export const DialogStateProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [openDialogs, setOpenDialogs] = useState<Map<string, DialogInfo>>(new Map());

  const registerDialog = useCallback((id: string, blocksKeyboardNavigation = false) => {
    setOpenDialogs((prev) => {
      const next = new Map(prev);
      next.set(id, { blocksKeyboardNavigation });
      return next;
    });
  }, []);

  const unregisterDialog = useCallback((id: string) => {
    setOpenDialogs((prev) => {
      const next = new Map(prev);
      next.delete(id);
      return next;
    });
  }, []);

  const isAnyDialogOpen = openDialogs.size > 0;
  const isKeyboardNavigationBlocked = Array.from(openDialogs.values()).some(
    (dialog) => dialog.blocksKeyboardNavigation
  );

  return (
    <DialogStateContext.Provider
      value={{
        isAnyDialogOpen,
        isKeyboardNavigationBlocked,
        registerDialog,
        unregisterDialog,
      }}
    >
      {children}
    </DialogStateContext.Provider>
  );
};

export const useDialogState = () => {
  const context = useContext(DialogStateContext);
  if (!context) {
    throw new Error('useDialogState must be used within a DialogStateProvider');
  }
  return context;
};

export const useDialogRegistration = (dialogId: string, blocksKeyboardNavigation = false) => {
  const { registerDialog, unregisterDialog } = useDialogState();

  useEffect(() => {
    registerDialog(dialogId, blocksKeyboardNavigation);
    return () => {
      unregisterDialog(dialogId);
    };
  }, [dialogId, blocksKeyboardNavigation, registerDialog, unregisterDialog]);
};