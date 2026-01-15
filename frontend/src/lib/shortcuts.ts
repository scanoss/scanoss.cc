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

import { entities } from '../../wailsjs/go/models';

interface KeyboardShortcut {
  name: string;
  description: string;
  keys: string;
}

export const KEYBOARD_SHORTCUTS: Record<entities.Action, KeyboardShortcut> = {
  // Global
  [entities.Action.Undo]: {
    name: 'Undo',
    description: 'Undo the last action',
    keys: 'mod+z',
  },
  [entities.Action.Redo]: {
    name: 'Redo',
    description: 'Redo the last action',
    keys: 'mod+shift+z',
  },
  [entities.Action.Save]: {
    name: 'Save',
    description: 'Save the current changes',
    keys: 'mod+s',
  },
  [entities.Action.Confirm]: {
    name: 'Confirm',
    description: 'Confirm the current action',
    keys: 'mod+enter',
  },
  [entities.Action.FocusSearch]: {
    name: 'Focus Search',
    description: 'Focus the search bar',
    keys: 'mod+f',
  },
  [entities.Action.SelectAll]: {
    name: 'Select All',
    description: 'Select all text',
    keys: 'mod+a',
  },

  // Navigation
  [entities.Action.MoveUp]: {
    name: 'Move Up',
    description: 'Move the selection up',
    keys: 'up, k',
  },
  [entities.Action.MoveDown]: {
    name: 'Move Down',
    description: 'Move the selection down',
    keys: 'down, j',
  },

  // Actions
  [entities.Action.Include]: {
    name: 'Include',
    description: 'Include file directly',
    keys: 'i, f1',
  },
  [entities.Action.IncludeWithModal]: {
    name: 'Include (with options)',
    description: 'Open include dialog for file/folder/component',
    keys: 'shift+i, shift+f1',
  },
  [entities.Action.Dismiss]: {
    name: 'Dismiss',
    description: 'Dismiss file directly',
    keys: 'd, f2',
  },
  [entities.Action.DismissWithModal]: {
    name: 'Dismiss (with options)',
    description: 'Open dismiss dialog for file/folder/component',
    keys: 'shift+d, shift+f2',
  },
  [entities.Action.Replace]: {
    name: 'Replace',
    description: 'Open replace dialog to select replacement component',
    keys: 'r, f3',
  },

  // View
  [entities.Action.ToggleSyncScrollPosition]: {
    name: 'Toggle Sync Scroll Position',
    description: 'Sync the scroll position of the editors',
    keys: 'shift+mod+e',
  },
  [entities.Action.ShowKeyboardShortcutsModal]: {
    name: 'Show Keyboard Shortcuts',
    description: 'Show the keyboard shortcuts modal',
    keys: 'mod+shift+k',
  },
  [entities.Action.OpenSettings]: {
    name: 'Open Settings',
    description: 'Open the app settings',
    keys: 'mod+,',
  },

  // Scan
  [entities.Action.ScanWithOptions]: {
    name: 'Scan With Options',
    description: 'Scan with options',
    keys: 'mod+shift+b',
  },
};

export const getShortcutDisplay = (keys: string, modifierKey: string): string[] => {
  return keys.split(',').map((keyCombo) => {
    const parts = keyCombo.trim().split('+');

    return parts
      .map((p) => {
        if (p.toLowerCase() === 'mod') return modifierKey;
        if (p.toLowerCase() === 'up') return '↑';
        if (p.toLowerCase() === 'down') return '↓';

        return p.toUpperCase();
      })
      .join('+');
  });
};
