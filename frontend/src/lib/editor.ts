import * as monaco from 'monaco-editor';

export interface EditorManager {
  addEditor(id: string, editor: monaco.editor.IStandaloneCodeEditor): void;
  scrollToLineIfNotVisible(id: string, line: number): void;
  highlightLines(id: string, ranges: { start: number; end: number }[], className: string): void;
  syncCursor(id: string): void;
  syncScroll(id: string): void;
  getScrollSyncEnabled(): boolean;
}

interface AddEditorOptions {
  revealLine?: number;
  highlight?: {
    ranges: HighlightRange[];
    className: string;
  };
}

export interface HighlightRange {
  start: number;
  end: number;
}

export class MonacoManager implements EditorManager {
  private static instance: MonacoManager;
  private editors: { id: string; editor: monaco.editor.IStandaloneCodeEditor }[] = [];
  private cursorSyncListeners: { [id: string]: monaco.IDisposable } = {};
  private scrollSyncListeners: { [id: string]: monaco.IDisposable } = {};
  private scrollSyncEnabled = true;
  private isScrolling = false;

  private constructor() {}

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

    if (options?.highlight) {
      this.highlightLines(id, options.highlight.ranges, options.highlight.className);
    }

    setTimeout(() => {
      if (options?.revealLine) {
        this.scrollToLineIfNotVisible(id, options.revealLine);
      }

      if (this.scrollSyncEnabled) {
        this.syncCursor(id);
        this.syncScroll(id);
      }
    }, 200);
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

  public highlightLines(id: string, ranges: HighlightRange[], className: string): void {
    const editor = this.getEditor(id);
    if (!editor) return;

    const decorations: monaco.editor.IModelDeltaDecoration[] = ranges.map(({ start, end }) => ({
      range: new monaco.Range(start, 1, end, 1),
      options: { isWholeLine: true, className },
    }));

    editor.createDecorationsCollection(decorations);
  }

  public syncScroll(id: string) {
    const editor = this.getEditor(id);
    if (!editor) return;

    let lastScrollTop = editor.getScrollTop();

    this.scrollSyncListeners[id] = editor.onDidScrollChange(() => {
      if (this.isScrolling) return;
      this.isScrolling = true;

      try {
        const sourceEditor = editor;
        const currentScrollTop = sourceEditor.getScrollTop();
        const deltaY = currentScrollTop - lastScrollTop;
        lastScrollTop = currentScrollTop;

        const sourceLineHeight = editor.getOption(monaco.editor.EditorOption.lineHeight);

        // This is the number of lines scrolled
        const linesScrolled = deltaY / sourceLineHeight;

        this.editors.forEach(({ id: otherId, editor: otherEditor }) => {
          if (otherId !== id) {
            const targetLineHeight = otherEditor.getOption(monaco.editor.EditorOption.lineHeight);
            const currentOtherScrollTop = otherEditor.getScrollTop();
            const maxScrollTop = otherEditor.getScrollHeight() - otherEditor.getLayoutInfo().height;

            const targetDeltaY = linesScrolled * targetLineHeight;
            const newScrollTop = Math.max(0, Math.min(currentOtherScrollTop + targetDeltaY, maxScrollTop));

            otherEditor.setScrollPosition({
              scrollTop: newScrollTop,
              scrollLeft: sourceEditor.getScrollLeft(),
            });
          }
        });
      } finally {
        requestAnimationFrame(() => {
          this.isScrolling = false;
        });
      }
    });
  }

  public syncCursor(id: string) {
    const editor = this.getEditor(id);
    if (!editor) return;

    this.cursorSyncListeners[id] = editor.onDidChangeCursorPosition(() => {
      const position = editor.getPosition();
      if (!position) return;

      this.editors.forEach(({ id: otherId, editor: otherEditor }) => {
        if (otherId !== id) otherEditor.setPosition(position);
      });
    });
  }

  public toggleSyncScroll() {
    this.scrollSyncEnabled = !this.scrollSyncEnabled;

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
}
