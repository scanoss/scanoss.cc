import { useEffect, useRef } from 'react';

type KeySequence = string[];
type ShortcutCallback = () => void;

interface KeyboardShortcutOptions {
  resetDelay?: number; // Time in ms to reset the key sequence
}

export default function useKeyboardShortcut(
  keys: KeySequence,
  callback: ShortcutCallback,
  options: KeyboardShortcutOptions = {}
): void {
  const { resetDelay = 1000 } = options;
  const pressedKeys = useRef<string[]>([]);
  const resetTimeout = useRef<NodeJS.Timeout>();

  const isPressingModifierKey = (event: KeyboardEvent): boolean => {
    return event.metaKey || event.ctrlKey;
  };

  const resetKeySequence = () => {
    pressedKeys.current = [];
  };

  const checkSequence = (currentKey: string): boolean => {
    const currentIndex = pressedKeys.current.length;
    return keys[currentIndex]?.toLowerCase() === currentKey.toLowerCase();
  };

  const handleKeyDown = (event: KeyboardEvent) => {
    if (!isPressingModifierKey(event)) {
      resetKeySequence();
      return;
    }

    const currentKey = event.key.toLowerCase();

    if (currentKey === 'meta' || currentKey === 'control') {
      return;
    }

    if (!checkSequence(currentKey)) {
      resetKeySequence();
      if (currentKey === keys[0].toLowerCase()) {
        pressedKeys.current.push(currentKey);
      }
      return;
    }

    if (keys.some((k) => k.toLowerCase() === currentKey)) {
      event.preventDefault();
    }

    pressedKeys.current.push(currentKey);

    if (resetTimeout.current) {
      clearTimeout(resetTimeout.current);
    }

    resetTimeout.current = setTimeout(resetKeySequence, resetDelay);

    if (pressedKeys.current.length === keys.length) {
      callback();
      resetKeySequence();
    }
  };

  useEffect(() => {
    document.addEventListener('keydown', handleKeyDown);

    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      if (resetTimeout.current) {
        clearTimeout(resetTimeout.current);
      }
    };
  }, [keys, callback, resetDelay]);
}
