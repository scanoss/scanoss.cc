import * as monaco from 'monaco-editor';

export interface EditorManager {
  addEditor(id: string, editor: monaco.editor.IStandaloneCodeEditor): void;
  scrollToLineIfNotVisible(id: string, line: number): void;
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

      this.syncCursor(id);
      this.syncScroll(id);
    }, 200);
  }

  public getEditor(id: string): monaco.editor.IStandaloneCodeEditor | null {
    return this.editors.find((e) => e.id === id)?.editor || null;
  }

  public scrollToLineIfNotVisible(id: string, line: number) {
    const editor = this.getEditor(id);
    if (!editor) return;

    editor.revealLineInCenterIfOutsideViewport(line, monaco.editor.ScrollType.Smooth);
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
      const position = editor.getPosition();
      if (!position) return;

      this.editors.forEach(({ id: otherId, editor: otherEditor }) => {
        if (otherId !== id) otherEditor.setPosition(position);
      });
    });
  }
}
