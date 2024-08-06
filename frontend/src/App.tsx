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
    <div
      id="App"
      style={{
        display: 'flex',
        alignItems: 'center',
        width: '100vw',
        height: '100vh',
      }}
    >
      <Button>alskdajsd</Button>
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
        language="c"
      />
    </div>
  );
}

export default App;
