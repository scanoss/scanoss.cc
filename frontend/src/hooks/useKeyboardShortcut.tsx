import { useEffect, useRef } from 'react';

interface KeyboardShortcutOptions {
  // Time in ms to reset the key sequence
  resetDelay?: number;
}

export default function useKeyboardShortcut(
  keys: string[],
  callback: () => void,
  options: KeyboardShortcutOptions = {}
): void {
  const { resetDelay = 1000 } = options;
  const pressedKeys = useRef<string[]>([]);
  const modifierPressed = useRef<boolean>(false);
  const resetTimeout = useRef<NodeJS.Timeout>();

  const isPressingModifierKey = (event: KeyboardEvent): boolean => {
    return event.metaKey || event.ctrlKey;
  };

  const isAnyAllowedKeyPressed = (event: KeyboardEvent) => {
    return keys.some((k) => k.toLowerCase() === event.key.toLowerCase());
  };

  const resetKeySequence = () => {
    pressedKeys.current = [];
    modifierPressed.current = false;
  };

  const handleKeyDown = (event: KeyboardEvent) => {
    const currentKey = event.key.toLowerCase();

    if (isPressingModifierKey(event)) {
      modifierPressed.current = true;
      return;
    }

    if (!modifierPressed.current) return;

    if (isAnyAllowedKeyPressed(event)) {
      event.preventDefault();
    }

    pressedKeys.current.push(currentKey);

    if (resetTimeout.current) {
      clearTimeout(resetTimeout.current);
    }

    resetTimeout.current = setTimeout(resetKeySequence, resetDelay);

    const currentSequence = pressedKeys.current;
    if (currentSequence.length === keys.length && currentSequence.every((k, i) => k === keys[i].toLowerCase())) {
      callback();
      resetKeySequence();
    }
  };

  const handleKeyUp = (event: KeyboardEvent) => {
    if (isPressingModifierKey(event)) {
      modifierPressed.current = false;
      pressedKeys.current = [];
    }
  };

  useEffect(() => {
    document.addEventListener('keydown', handleKeyDown);
    document.addEventListener('keyup', handleKeyUp);

    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      document.removeEventListener('keyup', handleKeyUp);
      if (resetTimeout.current) {
        clearTimeout(resetTimeout.current);
      }
    };
  }, [keys, callback, resetDelay]);
}
