import { Editor } from '@monaco-editor/react';
import * as monaco from 'monaco-editor';
import React, { useRef, useState } from 'react';

import files from '@/files';

export default function CodeViewer() {
  const [fileName] = useState('scanner.c');
  const leftEditor = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
  // const rightEditor = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
  const file = files[fileName];

  const handleLeftEditorMount = (
    editor: monaco.editor.IStandaloneCodeEditor
  ) => {
    leftEditor.current = editor;
    const START = 25;
    const END = 135;
    const decorations = [
      {
        range: new monaco.Range(START, 1, START, 1),
        options: {
          isWholeLine: true,
        },
      },
      {
        range: new monaco.Range(END, 1, END, 1),
        options: {
          isWholeLine: true,
        },
      },
      {
        range: new monaco.Range(START, 1, END, 1),
        options: {
          isWholeLine: true,
          className: 'lineHighlightDecoration',
          linesDecorationsClassName: 'lineRangeDecoration',
        },
      },
    ];
    leftEditor.current?.deltaDecorations([], decorations);
  };

  return (
    <Editor
      height="100%"
      width="50%"
      options={{
        minimap: { enabled: false },
        readOnly: true,
      }}
      path={file.name}
      defaultLanguage={file.language}
      defaultValue={file.value}
      theme="vs-dark"
      onMount={handleLeftEditorMount}
      // TODO: detect language from file extension
      language="c"
    />
  );
}
