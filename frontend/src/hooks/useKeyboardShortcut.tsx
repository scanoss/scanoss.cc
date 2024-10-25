import { Options, useHotkeys } from 'react-hotkeys-hook';
import { HotkeysEvent } from 'react-hotkeys-hook/dist/types';

// We re export this to use it always with default configurations already setup
export default function useKeyboardShortcut(
  keys: string | string[],
  callback: (event: KeyboardEvent, handler: HotkeysEvent) => void,
  options: Options = {
    preventDefault: true,
  },
  deps: unknown[] = []
) {
  return useHotkeys(keys, callback, options, deps);
}
