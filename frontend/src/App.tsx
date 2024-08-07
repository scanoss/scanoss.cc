import React, { useRef, useState } from 'react';
import './App.css';
import {
    Greet,
    GetFilesToBeCommited,
    FileGetLocalContent,
    FileGetRemoteContent,
    ResultGetAll, ComponentGet


} from "../wailsjs/go/main/App";

function App() {
  const [fileName, setFileName] = useState('scanner.c');
  const leftEditor = useRef<editor.IStandaloneCodeEditor | null>(null);
  const rightEditor = useRef<editor.IStandaloneCodeEditor | null>(null);
  const file = files[fileName];

    function greet() {
        GetFilesToBeCommited().then(m => console.log(m))
        Greet(name).then(updateResultText);
        FileGetLocalContent('main/pkg/file/adapter/controller.go').then((f)=> console.log("Local file content",f.content))
        FileGetRemoteContent('main/pkg/file/adapter/controller.go').then((f)=> console.log("Remote File content", f.content))
       // Filter on matchType can be applied . example : { matchType: "file" } will return all the files with matchType equal to file
        ResultGetAll( {} ).then((r)=> console.log("Results ", r))

        ComponentGet( '/external/inc/json.h' ).then((c)=> console.log("Result with matched component ", c))
        ComponentGet( '/external/inc/crc32c.h' ).then((c)=> console.log("Result without match ", c))


    }

  return (
    <div className="flex flex-row bg-slate-500 w-screen h-screen">
      {/* TODO: Move to another component to isolate the monaco dependency */}
      <div className="bg-red-50 w-1/4">Sidebar</div>
      <div className="w-full h-full bg-red-500">
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
          // TODO: detect language from file extension
          language="c"
        />
      </div>
    </div>
  );
}

export default App;
