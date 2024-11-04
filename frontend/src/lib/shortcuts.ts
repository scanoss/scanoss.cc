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
  [entities.Action.IncludeFileWithoutComments]: {
    name: 'Include File Without Comments',
    description: 'Include the file without comments',
    keys: 'i, f1',
  },
  [entities.Action.IncludeFileWithComments]: {
    name: 'Include File With Comments',
    description: 'Include the file with comments',
    keys: 'shift+i, shift+f1',
  },
  [entities.Action.IncludeComponentWithoutComments]: {
    name: 'Include Component Without Comments',
    description: 'Include the component without comments',
    keys: 'c, mod+f1',
  },
  [entities.Action.IncludeComponentWithComments]: {
    name: 'Include Component With Comments',
    description: 'Include the component with comments',
    keys: 'shift+c, shift+mod+f1',
  },

  [entities.Action.DismissFileWithoutComments]: {
    name: 'Dismiss File Without Comments',
    description: 'Dismiss the file without comments',
    keys: 'd, f2',
  },
  [entities.Action.DismissFileWithComments]: {
    name: 'Dismiss File With Comments',
    description: 'Dismiss the file with comments',
    keys: 'shift+d, shift+f2',
  },
  [entities.Action.DismissComponentWithoutComments]: {
    name: 'Dismiss Component Without Comments',
    description: 'Dismiss the component without comments',
    keys: 'x, mod+f2',
  },
  [entities.Action.DismissComponentWithComments]: {
    name: 'Dismiss Component With Comments',
    description: 'Dismiss the component with comments',
    keys: 'shift+x, shift+mod+f2',
  },

  [entities.Action.ReplaceFileWithoutComments]: {
    name: 'Replace File Without Comments',
    description: 'Replace the file without comments',
    keys: 'r, f3',
  },
  [entities.Action.ReplaceFileWithComments]: {
    name: 'Replace File With Comments',
    description: 'Replace the file with comments',
    keys: 'shift+r, shift+f3',
  },
  [entities.Action.ReplaceComponentWithoutComments]: {
    name: 'Replace Component Without Comments',
    description: 'Replace the component without comments',
    keys: 'e, mod+f3',
  },
  [entities.Action.ReplaceComponentWithComments]: {
    name: 'Replace Component With Comments',
    description: 'Replace the component with comments',
    keys: 'shift+e, shift+mod+f3',
  },

  // View
  [entities.Action.ToggleSyncScrollPosition]: {
    name: 'Toggle Sync Scroll Position',
    description: 'Sync the scroll position of the editors',
    keys: 'shift+mod+e',
  },
};

export const getShortcutDisplay = (keys: string): string[] => {
  return keys.split(',').map((keyCombo) => {
    const parts = keyCombo.trim().split('+');

    return parts
      .map((p) => {
        if (p.toLowerCase() === 'mod') return '⌘';
        if (p.toLowerCase() === 'up') return '↑';
        if (p.toLowerCase() === 'down') return '↓';

        return p.toUpperCase();
      })
      .join('+');
  });
};
