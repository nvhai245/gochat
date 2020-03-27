import React, { useEffect, useState, useRef } from 'react';
import logo from './logo.svg';
import './App.css';
import { wsConnect, sendMsg } from "./api";
import Header from './components/Header/Header';
import ChatHistory from './components/ChatHistory';
import ChatInput from './components/ChatInput';
import Popover from '@material-ui/core/Popover';
import SignUp from './components/SignUp';
import OnlineList from './components/OnlineList';
import { connect } from 'react-redux';
import { authorize } from './redux/actions/authorize';
import axios from 'axios';
import { db } from './db';

function App(props) {

  const [anchorEl, setAnchorEl] = useState(null);
  const [open, setOpen] = useState(true);
  const [users, setUsers] = useState([])

  const handleClose = () => {
    setAnchorEl(null);
    setOpen(false);
  };


  const [chatHistory, setChatHistory] = useState([]);
  const [onlineUsers, setOnlineUsers] = useState([]);
  const send = (event) => {
    if (event.keyCode === 13 && event.target.value !== "" && !event.shiftKey) {
      event.preventDefault();
      let newMsg;
      newMsg = { count: db.get('allcount').value() + 1, type: "chat", body: event.target.value, username: props.authorization.username, receiver: ["all"], table: "all" }
      sendMsg(JSON.stringify(newMsg));
      event.target.value = "";
      event.target.setAttribute('style','');
    }
  };
  const sendMessage = event => {
    let textInput = document.getElementById("textMessageInput");
    let newMsg
    newMsg = { count: db.get('allcount').value() + 1, type: "chat", body: textInput.value, username: props.authorization.username, receiver: ["all"], table: "all" }
      sendMsg(JSON.stringify(newMsg));
      textInput.value = "";
      textInput.setAttribute('style','');
  }
  useEffect(() => {
    wsConnect((msg) => {
      let msgData = JSON.parse(msg.data);
      if ((msgData.type === "chat" && msgData.receiver[0] === "all") || msgData.type === "system") {
        setChatHistory(prevState => ([...prevState, msgData]));
      }
      if (msgData.type === "chat") {
        if (!db.has(msgData.table).value()) {
          db.set(msgData.table, []).write()
          db.set(msgData.table + "count", 0).write()
        }
        db.get(msgData.table).push(msgData).write();
        db.update(msgData.table + 'count', n => n + 1).write();
      }
      if (msgData.type === "authfail") {
        alert("wrong username or password");
        window.location.reload();
      }
      if (msgData.type === "online" || msgData.type === "offline") {
        setOnlineUsers(msgData.body2);
      }
      if (msgData.type === "update") {
        let localCount = db.get(msgData.table + 'count').value();
        console.log("Difference is ", msgData.count - localCount);
        if (msgData.count - localCount > 0) {
          let newMsg;
          newMsg = { type: "readdb", body: localCount.toString(), body3: msgData.count.toString(), username: props.authorization.username, table: msgData.table };
          sendMsg(JSON.stringify(newMsg));
        }
        if (msgData.count - localCount < 0) {
          let newMsg;
          newMsg = { type: "writedb", body: localCount.toString(), body3: msgData.count.toString(), body2: [], username: props.authorization.username, table: msgData.table };
          for (let i = msgData.count + 1; i <= localCount; i++ ) {
            let writeMsg = db.get(msgData.table).find({count: i}).value();
            newMsg.body2.push(JSON.stringify(writeMsg));
          }
          sendMsg(JSON.stringify(newMsg));
        }
      }
      if (msgData.type === "readdb") {
        db.get(msgData.table).push(msgData).write();
        db.update(msgData.table + 'count', n => n + 1).write();
        setChatHistory(prevState => ([...prevState, msgData]));
      }
    });
    console.log(props.authorization.username)
    if (props.authorization.username !== "") {
      setOpen(false)
    }
    axios({
      method: 'get',
      url: 'http://localhost:8080/users',
      withCredentials: true,
      headers: { 'Content-Type': 'application/json' }
    })
      .then(function (response) {
        if (response.status === 200) {
          console.log(response.data)
          setUsers(response.data.users)
        }
      })
      .catch(function (response) {

      });
  }, [open]);

  return (
    <div className="App">
      <Header username={props.authorization.username} authorize={props.authorize} />
      {props.authorization.username !== "" &&
        <div className="appContainer">
          <div className="chatContainer">
            <ChatHistory currentUser={props.authorization.username} setChatHistory={setChatHistory} chatHistory={chatHistory} />
            <ChatInput authorizedUser={props.authorization.username} send={send} sendMessage={sendMessage} />
          </div>
          <OnlineList onlineUsers={onlineUsers} users={users} />
        </div>
      }
      <Popover
        open={open}
        anchorEl={anchorEl}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'center',
        }}
        marginThreshold={200}
      >
        <SignUp formValue={props.form.signup && props.form.signup.values} authorize={props.authorize} handleClose={handleClose} />
      </Popover>
    </div>
  );
}

const mapStateToProps = state => ({
  ...state
});

const mapDispatchToProps = dispatch => ({
  authorize: data => dispatch(authorize(data))
});

export default connect(mapStateToProps, mapDispatchToProps)(App);
