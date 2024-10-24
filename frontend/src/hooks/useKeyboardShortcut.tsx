import { useEffect, useRef } from 'react';

type KeySequence = string[];
type ShortcutCallback = () => void;

interface KeyboardShortcutOptions {
  resetDelay?: number;
  preventRegistering?: boolean;
}

interface Shortcut {
  id: string;
  keys: KeySequence;
  callback: ShortcutCallback;
  pressedKeys: string[];
  resetDelay: number;
  resetTimeout?: NodeJS.Timeout;
}

class ShortcutManager {
  private static instance: ShortcutManager;
  private shortcuts: Map<string, Shortcut> = new Map();
  private listenerAttached = false;

  private constructor() {}

  static getInstance(): ShortcutManager {
    if (!ShortcutManager.instance) {
      ShortcutManager.instance = new ShortcutManager();
    }
    return ShortcutManager.instance;
  }

  private setupListener() {
    if (this.listenerAttached) return;

    document.addEventListener('keydown', this.handleKeyDown);
    this.listenerAttached = true;
  }

  private handleKeyDown = (event: KeyboardEvent) => {
    const currentKey = event.key.toLowerCase();

    this.shortcuts.forEach((shortcut) => {
      const currentIndex = shortcut.pressedKeys.length;
      const expectedKey = shortcut.keys[currentIndex]?.toLowerCase();

      if (currentKey !== expectedKey) {
        this.resetShortcut(shortcut.id);
        if (currentKey === shortcut.keys[0].toLowerCase()) {
          shortcut.pressedKeys.push(currentKey);
        }
        return;
      }

      if (shortcut.keys.some((k) => k.toLowerCase() === currentKey)) {
        event.preventDefault();
      }

      shortcut.pressedKeys.push(currentKey);

      if (shortcut.resetTimeout) {
        clearTimeout(shortcut.resetTimeout);
      }

      shortcut.resetTimeout = setTimeout(() => {
        this.resetShortcut(shortcut.id);
      }, shortcut.resetDelay);

      if (shortcut.pressedKeys.length === shortcut.keys.length) {
        shortcut.callback();
        this.resetShortcut(shortcut.id);
      }
    });
  };

  private resetShortcut(id: string) {
    const shortcut = this.shortcuts.get(id);
    if (shortcut) {
      shortcut.pressedKeys = [];
      if (shortcut.resetTimeout) {
        clearTimeout(shortcut.resetTimeout);
      }
    }
  }

  registerShortcut(shortcut: Shortcut): void {
    this.shortcuts.set(shortcut.id, shortcut);
    this.setupListener();
  }

  unregisterShortcut(id: string): void {
    const shortcut = this.shortcuts.get(id);
    if (shortcut?.resetTimeout) {
      clearTimeout(shortcut.resetTimeout);
    }
    this.shortcuts.delete(id);

    if (this.shortcuts.size === 0 && this.listenerAttached) {
      document.removeEventListener('keydown', this.handleKeyDown);
      this.listenerAttached = false;
    }
  }
}

let shortcutCounter = 0;

export default function useKeyboardShortcut(
  keys: KeySequence,
  callback: ShortcutCallback,
  options: KeyboardShortcutOptions = {}
): void {
  const { resetDelay = 1000, preventRegistering = false } = options;
  const shortcutId = useRef(`shortcut-${++shortcutCounter}`);
  const manager = ShortcutManager.getInstance();

  useEffect(() => {
    if (preventRegistering) return;

    const shortcut: Shortcut = {
      id: shortcutId.current,
      keys,
      callback,
      pressedKeys: [],
      resetDelay,
    };

    manager.registerShortcut(shortcut);

    return () => {
      manager.unregisterShortcut(shortcutId.current);
    };
  }, [keys, callback, resetDelay, preventRegistering]);
}
