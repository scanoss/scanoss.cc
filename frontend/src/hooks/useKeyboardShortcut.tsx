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

import { Options, useHotkeys } from 'react-hotkeys-hook';
import { HotkeysEvent } from 'react-hotkeys-hook/dist/types';

import useEditorFocusStore from '@/stores/useEditorFocusStore';

// System shortcuts that should not be intercepted
const SYSTEM_SHORTCUTS = ['mod+a', 'mod+c', 'mod+v', 'mod+x'];

// We re export this to use it always with default configurations already setup
export default function useKeyboardShortcut(
  keys: string | string[],
  callback: (event: KeyboardEvent, handler: HotkeysEvent) => void,
  options?: Options,
  deps: unknown[] = []
) {
  const keysArray = Array.isArray(keys) ? keys : [keys];
  const hasSystemShortcut = keysArray.some((key) => SYSTEM_SHORTCUTS.includes(key.toLowerCase().trim()));

  // Don't prevent default for system shortcuts to allow copy, paste, cut, etc.
  const shouldPreventDefault = hasSystemShortcut ? false : true;

  // Allow shortcuts on form elements when the Monaco editor has focus,
  // since Monaco uses an internal textarea for keyboard capture
  const editorHasFocus = useEditorFocusStore((s) => s.hasFocus);

  return useHotkeys(
    keys,
    callback,
    {
      preventDefault: shouldPreventDefault,
      ...options,
      enableOnFormTags: options?.enableOnFormTags || editorHasFocus || false,
    },
    [...deps, editorHasFocus]
  );
}
