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

import * as monaco from 'monaco-editor';

export interface EditorManager {
  addEditor(id: string, editor: monaco.editor.IStandaloneCodeEditor): void;
  scrollToLineIfNotVisible(id: string, line: number): void;
  syncScroll(id: string): void;
  getScrollSyncEnabled(): boolean;
}

interface AddEditorOptions {
  revealLine?: number;
}

export interface HighlightRange {
  start: number;
  end: number;
}

export class MonacoManager implements EditorManager {
  private static instance: MonacoManager;
  private editors: { id: string; editor: monaco.editor.IStandaloneCodeEditor }[] = [];
  private scrollSyncListeners: { [id: string]: monaco.IDisposable } = {};
  private scrollSyncEnabled = true;
  private isScrolling = false;
  private scrollThrottleTimeout: number | null = null;
  private readonly SCROLL_THROTTLE_MS = 16; // ~60fps

  private constructor() {}

  private throttle(callback: () => void) {
    if (this.scrollThrottleTimeout !== null) return;

    callback();

    this.scrollThrottleTimeout = window.setTimeout(() => {
      this.scrollThrottleTimeout = null;
    }, this.SCROLL_THROTTLE_MS);
  }

  public static getInstance(): MonacoManager {
    if (!MonacoManager.instance) {
      MonacoManager.instance = new MonacoManager();
    }
    return MonacoManager.instance;
  }

  public addEditor(id: string, editor: monaco.editor.IStandaloneCodeEditor, options?: AddEditorOptions) {
    const existingEditorIndex = this.editors.findIndex((e) => e.id === id);
    if (existingEditorIndex > -1) {
      this.editors[existingEditorIndex] = { id, editor };
    } else {
      this.editors.push({ id, editor });
    }

    if (options?.revealLine) {
      this.scrollToLineIfNotVisible(id, options.revealLine);
    }

    if (this.scrollSyncEnabled) {
      this.syncScroll(id);
    }
  }

  public getScrollSyncEnabled(): boolean {
    return this.scrollSyncEnabled;
  }

  public getEditor(id: string): monaco.editor.IStandaloneCodeEditor | null {
    return this.editors.find((e) => e.id === id)?.editor || null;
  }

  public scrollToLineIfNotVisible(id: string, line: number): void {
    const editor = this.getEditor(id);
    if (!editor) return;

    editor.revealLineInCenterIfOutsideViewport(line, monaco.editor.ScrollType.Smooth);
  }

  public syncScroll(id: string) {
    const editor = this.getEditor(id);
    if (!editor) return;

    let lastScrollTop = editor.getScrollTop();
    let lastScrollLeft = editor.getScrollLeft();

    this.scrollSyncListeners[id] = editor.onDidScrollChange(() => {
      if (this.isScrolling) return;

      this.throttle(() => {
        this.isScrolling = true;

        try {
          const sourceEditor = editor;
          const currentScrollTop = sourceEditor.getScrollTop();
          const currentScrollLeft = sourceEditor.getScrollLeft();

          // Only sync if scroll position actually changed
          if (currentScrollTop === lastScrollTop && currentScrollLeft === lastScrollLeft) {
            return;
          }

          const deltaY = currentScrollTop - lastScrollTop;
          lastScrollTop = currentScrollTop;
          lastScrollLeft = currentScrollLeft;

          const sourceLineHeight = editor.getOption(monaco.editor.EditorOption.lineHeight);
          const linesScrolled = deltaY / sourceLineHeight;

          // Pre-calculate scroll values for better performance
          const scrollUpdates = this.editors
            .filter(({ id: otherId }) => otherId !== id)
            .map(({ editor: otherEditor }) => {
              const targetLineHeight = otherEditor.getOption(monaco.editor.EditorOption.lineHeight);
              const currentOtherScrollTop = otherEditor.getScrollTop();
              const maxScrollTop = otherEditor.getScrollHeight() - otherEditor.getLayoutInfo().height;
              const targetDeltaY = linesScrolled * targetLineHeight;

              return {
                editor: otherEditor,
                scrollTop: Math.max(0, Math.min(currentOtherScrollTop + targetDeltaY, maxScrollTop)),
                scrollLeft: currentScrollLeft,
              };
            });

          // Batch scroll updates in next animation frame
          requestAnimationFrame(() => {
            scrollUpdates.forEach(({ editor: otherEditor, scrollTop, scrollLeft }) => {
              otherEditor.setScrollPosition({ scrollTop, scrollLeft });
            });
          });
        } finally {
          requestAnimationFrame(() => {
            this.isScrolling = false;
          });
        }
      });
    });
  }

  public toggleSyncScroll() {
    this.scrollSyncEnabled = !this.scrollSyncEnabled;

    // Clear any pending throttle timeout
    if (this.scrollThrottleTimeout !== null) {
      clearTimeout(this.scrollThrottleTimeout);
      this.scrollThrottleTimeout = null;
    }

    if (!Object.keys(this.scrollSyncListeners).length) {
      return this.editors.forEach(({ id }) => this.syncScroll(id));
    }

    return this.editors.forEach(({ id }) => {
      if (this.scrollSyncListeners[id]) {
        this.scrollSyncListeners[id].dispose();
        delete this.scrollSyncListeners[id];
      }
    });
  }

  public dispose() {
    if (this.scrollThrottleTimeout !== null) {
      clearTimeout(this.scrollThrottleTimeout);
    }

    Object.values(this.scrollSyncListeners).forEach((listener) => listener.dispose());
    this.scrollSyncListeners = {};
    this.editors = [];
  }
}
