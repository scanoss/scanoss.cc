import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {
    Greet,
    GetFilesToBeCommited,
    FileGetLocalContent,
    FileGetRemoteContent
} from "../wailsjs/go/main/App";

function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const updateName = (e: any) => setName(e.target.value);
    const updateResultText = (result: string) => setResultText(result);

    function greet() {
        GetFilesToBeCommited().then(m => console.log(m))
        Greet(name).then(updateResultText);
        FileGetLocalContent('main/pkg/file/adapter/controller.go').then((f)=> console.log("Local file content",f.content))
        FileGetRemoteContent('main/pkg/file/adapter/controller.go').then((f)=> console.log("Remote File content", f.content))
    }

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <button className="btn" onClick={greet}>Greet</button>
            </div>
        </div>
    )
}

export default App
