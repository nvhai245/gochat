import React, { useEffect, useState } from 'react';
import logo from './logo.svg';
import './App.css';
import { connect, sendMsg } from "./api";
import Header from './components/Header/Header';
import ChatHistory from './components/ChatHistory/ChatHistory';
import ChatInput from './components/ChatInput/ChatInput';

function App() {
  const [chatHistory, setChatHistory] = useState([]);
  const [username, setUsername] = useState("");
  const updateUsername = (event) => {
    setUsername(event.target.value);
  }
  const send = (event) => {
    if (event.keyCode === 13) {
      let newMsg;
      newMsg = { type: 1, body: event.target.value, username: username }
      newMsg.username = username;
      sendMsg(JSON.stringify(newMsg));
      event.target.value = "";
    }
  };
  useEffect(() => {
    connect((msg) => {
      console.log("New Message")
      setChatHistory(prevState => ([...prevState, msg]));
    })
  }, []);

  return (
    <div className="App">
      <Header />
      <div className="chatContainer">
        <ChatHistory chatHistory={chatHistory} />
        <ChatInput send={send} updateUsername={updateUsername} />
      </div>
    </div>
  );
}

export default App;
