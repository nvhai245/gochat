import React, { useEffect, useState, useRef } from 'react';
import Paper from '@material-ui/core/Paper';
import './App.css';
import { wsConnect, sendMsg } from "./api";
import Header from './components/Header/Header';
import ChatHistory from './components/ChatHistory';
import ChatInput from './components/ChatInput';
import SignUp from './components/SignUp';
import OnlineList from './components/OnlineList';
import Inbox from './components/Inbox';
import { connect } from 'react-redux';
import { authorize } from './redux/actions/authorize';
import axios from 'axios';
import { db } from './db';

function App(props) {

  const [users, setUsers] = useState([]);
  const [inboxList, setInboxList] = useState([]);
  const [newlyCreatedTable, setNewlyCreatedTable] = useState([]);
  const [mostRecentMsg, setMostRecentMsg] = useState({});
  const [noti, setNoti] = useState([]);

  const addInboxList = (target) => {
    if (inboxList.indexOf(target) < 0) {
      let l = inboxList;
      l.push(target);
      let newInboxList = [...new Set(l)];
      setInboxList(newInboxList);
    } else {
      let newInBoxList = inboxList.filter(user => user !== target);
      setInboxList(newInBoxList);
    }
  }


  const [chatHistory, setChatHistory] = useState([]);
  const [onlineUsers, setOnlineUsers] = useState([]);
  const send = (event) => {
    if (event.keyCode === 13 && !event.shiftKey) {
      event.preventDefault();
    }
    if (event.keyCode === 13 && event.target.value !== "" && event.target.value.indexOf("\n") !== 0 && !event.shiftKey) {
      let newMsg;
      newMsg = { count: db.get('allcount').value() + 1, type: "chat", body: event.target.value, username: props.authorization.username, receiver: ["all"], table: "all", deleted: false }
      sendMsg(JSON.stringify(newMsg));
      event.target.value = "";
      event.target.setAttribute('style', '');
    }
  };
  const sendMessage = event => {
    let textInput = document.getElementById("textMessageInput");
    let newMsg
    newMsg = { count: db.get('allcount').value() + 1, type: "chat", body: textInput.value, username: props.authorization.username, receiver: ["all"], table: "all", deleted: false }
    sendMsg(JSON.stringify(newMsg));
    textInput.value = "";
    textInput.setAttribute('style', '');
  }
  useEffect(() => {
    wsConnect((msg) => {
      let msgData = JSON.parse(msg.data);
      if (msgData.type === "chat" || msgData.type === "readdb" || msgData.type === "checkExist" || msgData.type === "delete" || msgData.type === "restore") {
        setMostRecentMsg(msgData);
      }
      if (msgData.type === "system") {
        setChatHistory(prevState => ([...prevState, msgData]));
      }
      if (msgData.type === "chat" && msgData.table === "all") {
        if (!db.has(msgData.table).value()) {
          db.set(msgData.table, []).write()
          db.set(msgData.table + "count", 0).write()
        }
        db.get(msgData.table).push(msgData).write();
        db.update(msgData.table + 'count', n => n + 1).write();
        setChatHistory(prevState => ([...prevState, msgData]));
      }
      if (msgData.type === "chat" && msgData.username !== props.authorization.username && msgData.receiver[1] === props.authorization.username) {
        setInboxList(prevState => ([...new Set([...prevState, msgData.username])]));
      }
      if (msgData.type === "authfail") {
        alert("wrong username or password");
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
          // Dangerous
          // let newMsg;
          // newMsg = { type: "writedb", body: localCount.toString(), body3: msgData.count.toString(), body2: [], username: props.authorization.username, table: msgData.table };
          // for (let i = msgData.count + 1; i <= localCount; i++) {
          //   let writeMsg = db.get(msgData.table).find({ count: i }).value();
          //   newMsg.body2.push(JSON.stringify(writeMsg));
          // }
          // sendMsg(JSON.stringify(newMsg));
        }
      }
      if (msgData.type === "readdb" && msgData.table === "all") {
        db.get(msgData.table).push(msgData).write();
        db.update(msgData.table + 'count', n => n + 1).write();
        setChatHistory(prevState => ([...prevState, msgData]));
      }
      if (msgData.type === "checkExist") {
        if (!db.has(msgData.table).value()) {
          setNewlyCreatedTable(prevState => ([...prevState, msgData.table]));
        }
      }

      if ((msgData.type === "delete") && msgData.table === "all") {
        db.get('all').find({ count: msgData.count }).assign({ deleted: true }).write();
        setChatHistory(prevState => ([...prevState]));
      }
      if ((msgData.type === "restore") && msgData.table === "all") {
        db.get('all').find({ count: msgData.count }).assign({ deleted: false }).write();
        setChatHistory(prevState => ([...prevState]));
      }
      if (msgData.type === "getnoti") {
        setNoti(prevState => ([...prevState, msgData]));
      }
    });
  }, [props.authorization.username]);
  useEffect(() => {
    function getCookie(name) {
      var value = "; " + document.cookie;
      var parts = value.split("; " + name + "=");
      if (parts.length == 2) return parts.pop().split(";").shift();
    }
    let data = { username: getCookie('authorizedUser'), isAdmin: false };
    if (data.username) {
      props.authorize(data);
    }

  }, []);

  useEffect(() => {
    if (onlineUsers.length > 0) {
      let newMsg;
      newMsg = { type: "getnoti", username: props.authorization.username };
      sendMsg(JSON.stringify(newMsg));
    }
  }, [onlineUsers])

  useEffect(() => {
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
  }, []);

  return (
    <div className="App">
      <Header setNoti={setNoti} addInboxList={addInboxList} noti={noti} username={props.authorization.username} authorize={props.authorize} />
      {props.authorization.username &&
        <div>
          <div className="appContainer">
            <div className="chatContainer">
              <div className="chatHistoryContainer">
                <ChatHistory table="all" currentUser={props.authorization.username} setChatHistory={setChatHistory} chatHistory={chatHistory} />
              </div>
              <ChatInput authorizedUser={props.authorization.username} send={send} sendMessage={sendMessage} />
            </div>
            <OnlineList addInboxList={addInboxList} onlineUsers={onlineUsers} users={users} />
          </div>
          <div className="inboxArea">
            {inboxList.map((user, i) =>
              <Paper className="inboxContainer" elevation={3}>
                <Inbox key={user} mostRecentMsg={mostRecentMsg} onlineUsers={onlineUsers} newlyCreatedTable={newlyCreatedTable} sendMsg={sendMsg} currentUser={props.authorization.username} addInboxList={addInboxList} user={user} />
              </Paper>
            )}
          </div>
        </div>

      }
      {!props.authorization.username &&
        <SignUp formValue={props.form.signup && props.form.signup.values} authorize={props.authorize} />
      }
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
