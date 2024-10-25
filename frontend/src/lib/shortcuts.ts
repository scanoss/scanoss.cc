export enum KeyboardShortcutAction {
  UNDO = 'undo',
  REDO = 'redo',
  SAVE = 'save',
  CONFIRM = 'confirm',

  FOCUS_SEARCH = 'focusSearch',
  SELECT_ALL = 'selectAll',

  MOVE_UP = 'moveUp',
  MOVE_DOWN = 'moveDown',

  INCLUDE_FILE_WITHOUT_COMMENTS = 'includeFileWithoutComments',
  INCLUDE_FILE_WITH_COMMENTS = 'includeFileWithComments',
  INCLUDE_COMPONENT_WITHOUT_COMMENTS = 'includeComponentWithoutComments',
  INCLUDE_COMPONENT_WITH_COMMENTS = 'includeComponentWithComments',

  DISMISS_FILE_WITHOUT_COMMENTS = 'dismissFileWithoutComments',
  DISMISS_FILE_WITH_COMMENTS = 'dismissFileWithComments',
  DISMISS_COMPONENT_WITHOUT_COMMENTS = 'dismissComponentWithoutComments',
  DISMISS_COMPONENT_WITH_COMMENTS = 'dismissComponentWithComments',

  REPLACE_FILE_WITHOUT_COMMENTS = 'replaceFileWithoutComments',
  REPLACE_FILE_WITH_COMMENTS = 'replaceFileWithComments',
  REPLACE_COMPONENT_WITHOUT_COMMENTS = 'replaceComponentWithoutComments',
  REPLACE_COMPONENT_WITH_COMMENTS = 'replaceComponentWithComments',
}

interface KeyboardShortcut {
  name: string;
  description: string;
  keys: string;
}

export const KEYBOARD_SHORTCUTS: Record<KeyboardShortcutAction, KeyboardShortcut> = {
  // Global
  [KeyboardShortcutAction.UNDO]: {
    name: 'Undo',
    description: 'Undo the last action',
    keys: 'mod+z',
  },
  [KeyboardShortcutAction.REDO]: {
    name: 'Redo',
    description: 'Redo the last action',
    keys: 'mod+shift+z',
  },
  [KeyboardShortcutAction.SAVE]: {
    name: 'Save',
    description: 'Save the current changes',
    keys: 'mod+s',
  },
  [KeyboardShortcutAction.CONFIRM]: {
    name: 'Confirm',
    description: 'Confirm the current action',
    keys: 'mod+enter',
  },
  [KeyboardShortcutAction.FOCUS_SEARCH]: {
    name: 'Focus Search',
    description: 'Focus the search bar',
    keys: 'mod+f',
  },
  [KeyboardShortcutAction.SELECT_ALL]: {
    name: 'Select All',
    description: 'Select all text',
    keys: 'mod+a',
  },

  // Navigation
  [KeyboardShortcutAction.MOVE_UP]: {
    name: 'Move Up',
    description: 'Move the selection up',
    keys: 'up, k',
  },
  [KeyboardShortcutAction.MOVE_DOWN]: {
    name: 'Move Down',
    description: 'Move the selection down',
    keys: 'down, j',
  },

  // Actions
  [KeyboardShortcutAction.INCLUDE_FILE_WITHOUT_COMMENTS]: {
    name: 'Include File Without Comments',
    description: 'Include the file without comments',
    keys: 'i, f1',
  },
  [KeyboardShortcutAction.INCLUDE_FILE_WITH_COMMENTS]: {
    name: 'Include File With Comments',
    description: 'Include the file with comments',
    keys: 'shift+i, shift+f1',
  },
  [KeyboardShortcutAction.INCLUDE_COMPONENT_WITHOUT_COMMENTS]: {
    name: 'Include Component Without Comments',
    description: 'Include the component without comments',
    keys: 'c, mod+f1',
  },
  [KeyboardShortcutAction.INCLUDE_COMPONENT_WITH_COMMENTS]: {
    name: 'Include Component With Comments',
    description: 'Include the component with comments',
    keys: 'shift+c, shift+mod+f1',
  },

  [KeyboardShortcutAction.DISMISS_FILE_WITHOUT_COMMENTS]: {
    name: 'Dismiss File Without Comments',
    description: 'Dismiss the file without comments',
    keys: 'd, f2',
  },
  [KeyboardShortcutAction.DISMISS_FILE_WITH_COMMENTS]: {
    name: 'Dismiss File With Comments',
    description: 'Dismiss the file with comments',
    keys: 'shift+d, shift+f2',
  },
  [KeyboardShortcutAction.DISMISS_COMPONENT_WITHOUT_COMMENTS]: {
    name: 'Dismiss Component Without Comments',
    description: 'Dismiss the component without comments',
    keys: 'x, mod+f2',
  },
  [KeyboardShortcutAction.DISMISS_COMPONENT_WITH_COMMENTS]: {
    name: 'Dismiss Component With Comments',
    description: 'Dismiss the component with comments',
    keys: 'shift+x, shift+mod+f2',
  },

  [KeyboardShortcutAction.REPLACE_FILE_WITHOUT_COMMENTS]: {
    name: 'Replace File Without Comments',
    description: 'Replace the file without comments',
    keys: 'r, f3',
  },
  [KeyboardShortcutAction.REPLACE_FILE_WITH_COMMENTS]: {
    name: 'Replace File With Comments',
    description: 'Replace the file with comments',
    keys: 'shift+r, shift+f3',
  },
  [KeyboardShortcutAction.REPLACE_COMPONENT_WITHOUT_COMMENTS]: {
    name: 'Replace Component Without Comments',
    description: 'Replace the component without comments',
    keys: 'e, mod+f3',
  },
  [KeyboardShortcutAction.REPLACE_COMPONENT_WITH_COMMENTS]: {
    name: 'Replace Component With Comments',
    description: 'Replace the component with comments',
    keys: 'shift+e, shift+mod+f3',
  },
};
