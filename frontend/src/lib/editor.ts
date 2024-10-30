import * as monaco from 'monaco-editor';

export interface EditorManager {
  addEditor(id: string, editor: monaco.editor.IStandaloneCodeEditor): void;
  scrollToLine(id: string, line: number): void;
  highlightLines(id: string, ranges: { start: number; end: number }[], className: string): void;
  syncCursor(id: string): void;
  syncScroll(id: string): void;
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
  private isProgrammaticScroll: boolean = false;
  private isProgrammaticCursor: boolean = false;

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

    if (options?.revealLine) {
      this.scrollToLine(id, options.revealLine);
    }

    if (options?.highlight) {
      this.highlightLines(id, options.highlight.ranges, options.highlight.className);
    }

    this.syncScroll(id);
    this.syncCursor(id);
  }

  public getEditor(id: string): monaco.editor.IStandaloneCodeEditor | null {
    return this.editors.find((e) => e.id === id)?.editor || null;
  }

  public scrollToLine(id: string, line: number) {
    const editor = this.getEditor(id);
    if (!editor) return;

    this.isProgrammaticScroll = true;
    editor.revealLineInCenterIfOutsideViewport(line, monaco.editor.ScrollType.Smooth);
    setTimeout(() => (this.isProgrammaticScroll = false), 100); // Adjust timeout as needed
  }

  public highlightLines(id: string, ranges: HighlightRange[], className: string) {
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
    if (!editor || this.scrollSyncListeners[id]) return;

    this.scrollSyncListeners[id] = editor.onDidScrollChange(() => {
      if (this.isProgrammaticScroll) return;

      const scrollTop = editor.getScrollTop();
      this.editors.forEach(({ id: otherId, editor: otherEditor }) => {
        if (otherId !== id) otherEditor.setScrollTop(scrollTop);
      });
    });
  }

  public syncCursor(id: string) {
    const editor = this.getEditor(id);
    if (!editor || this.cursorSyncListeners[id]) return;

    this.cursorSyncListeners[id] = editor.onDidChangeCursorPosition(() => {
      if (this.isProgrammaticCursor) return;

      const position = editor.getPosition();
      if (!position) return;

      this.isProgrammaticCursor = true;
      this.editors.forEach(({ id: otherId, editor: otherEditor }) => {
        if (otherId !== id) otherEditor.setPosition(position);
      });
      setTimeout(() => (this.isProgrammaticCursor = false), 100); // Adjust timeout as needed
    });
  }
}
