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

import { useEffect } from 'react';

import { entities } from '../../wailsjs/go/models';
import { EventsOn } from '../../wailsjs/runtime';

/**
 * Custom hook to subscribe to Wails menu bar events.
 * Handles cleanup automatically when the component unmounts or dependencies change.
 *
 * @param event - The event name to listen for (e.g., entities.Action.Save)
 * @param handler - The callback function to execute when the event fires
 */
export function useMenuEvent(event: entities.Action, handler: () => void) {
  useEffect(() => {
    const unsub = EventsOn(event, handler);
    return () => unsub();
  }, [event, handler]);
}

/**
 * Custom hook to subscribe to multiple Wails menu bar events at once.
 * Useful when a component needs to handle several related events.
 *
 * @param eventHandlerMap - Object mapping event names to their handlers (null handlers are skipped)
 */
export function useMenuEvents(eventHandlerMap: Record<string, (() => void) | null>) {
  useEffect(() => {
    const unsubs: (() => void)[] = [];

    Object.entries(eventHandlerMap).forEach(([event, handler]) => {
      if (handler) {
        unsubs.push(EventsOn(event, handler));
      }
    });

    return () => unsubs.forEach((unsub) => unsub());
  }, [eventHandlerMap]);
}
