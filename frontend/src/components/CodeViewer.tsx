import { Editor } from '@monaco-editor/react';
import React from 'react';

import { Skeleton } from './ui/skeleton';

interface CodeViewerProps {
  content: string | undefined;
  language: string | null | undefined;
  width?: string;
  height?: string;
  isLoading: boolean;
  isError: boolean;
}

export default function CodeViewer({
  width = '100%',
  height = '100%',
  content,
  language,
  isLoading,
  isError,
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

  if (isLoading) {
    return <Skeleton className="w-full h-full" />;
  }

  if (isError) {
    return (
      <div className="w-full h-full flex items-center justify-center">
        <p className="text-red-500">Error loading file</p>
      </div>
    );
  }

  return (
    <Editor
      height={height}
      width={width}
      value={content}
      theme="vs-dark"
      // onMount={handleLeftEditorMount}
      {...(language ? { language } : {})}
      options={{
        minimap: { enabled: false },
        readOnly: true,
        wordWrap: 'on',
      }}
    />
  );
}
