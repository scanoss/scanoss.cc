import { Editor } from '@monaco-editor/react';
import React from 'react';

interface CodeViewerProps {
  content: string;
  language: string | null;
  width?: string;
  height?: string;
}

export default function CodeViewer({
  width = '100%',
  height = '100%',
  content,
  language,
}: CodeViewerProps) {
  // const [fileName] = useState('scanner.c');
  // const leftEditor = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
  // // const rightEditor = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
  // const file = files[fileName];

  // const handleLeftEditorMount = (
  //   editor: monaco.editor.IStandaloneCodeEditor
  // ) => {
  //   leftEditor.current = editor;
  //   const START = 25;
  //   const END = 135;
  //   const decorations = [
  //     {
  //       range: new monaco.Range(START, 1, START, 1),
  //       options: {
  //         isWholeLine: true,
  //       },
  //     },
  //     {
  //       range: new monaco.Range(END, 1, END, 1),
  //       options: {
  //         isWholeLine: true,
  //       },
  //     },
  //     {
  //       range: new monaco.Range(START, 1, END, 1),
  //       options: {
  //         isWholeLine: true,
  //         className: 'lineHighlightDecoration',
  //         linesDecorationsClassName: 'lineRangeDecoration',
  //       },
  //     },
  //   ];
  //   leftEditor.current?.deltaDecorations([], decorations);
  // };

  return (
    <Editor
      height={height}
      width={width}
      defaultValue={content}
      theme="vs-dark"
      // onMount={handleLeftEditorMount}
      {...(language ? { language } : {})}
      options={{
        minimap: { enabled: false },
        readOnly: true,
      }}
    />
  );
}
