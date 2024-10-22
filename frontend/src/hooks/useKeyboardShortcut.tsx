import { useEffect } from 'react';

export default function useKeyboardShortcut(key: string, callback: () => void): void {
  const isPressingModifierKey = (event: KeyboardEvent): boolean => {
    return event.metaKey || event.ctrlKey;
  };

  const isPressingShortcutKey = (event: KeyboardEvent): boolean => {
    return event.key === key;
  };

  useEffect(() => {
    document.addEventListener('keydown', (event) => {
      event.preventDefault();

      if (!isPressingModifierKey(event)) return;
      if (!isPressingShortcutKey(event)) return;

      callback();
    });
    return () => {
      document.removeEventListener('keydown', () => {});
    };
  }, [key, callback]);
}
