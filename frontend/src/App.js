import React, { useEffect, useState } from 'react';
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

function App(props) {

  const [anchorEl, setAnchorEl] = useState(null);
  const [open, setOpen] = useState(true);
  const [users, setUsers] = useState([])
  const handleClick = event => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
    setOpen(false);
  };


  const [chatHistory, setChatHistory] = useState([]);
  const [username, setUsername] = useState("");
  const [authorizedUser, setAuthorizedUser] = useState("");
  const [onlineUsers, setOnlineUsers] = useState([]);
  const updateUsername = (event) => {
    setUsername(event.target.value);
  }
  const send = (event) => {
    if (event.keyCode === 13) {
      let newMsg;
      newMsg = { type: "chat", body: event.target.value, username: props.authorization.username }
      sendMsg(JSON.stringify(newMsg));
      event.target.value = "";
    }
  };
  useEffect(() => {
    wsConnect((msg) => {
      let msgData = JSON.parse(msg.data);
      if (msgData.type !== "login" && msgData.type !== "register" && msgData.type !== "authfail") {
        setChatHistory(prevState => ([...prevState, msg]));
      }
      if (msgData.type === "authfail") {
        alert("wrong username or password");
        window.location.reload();
      }
      if(msgData.type === "online" || msgData.type === "offline") {
          setOnlineUsers(msgData.body2);
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
      <div className="appContainer">
        <div className="chatContainer">
          <ChatHistory chatHistory={chatHistory} />
          <ChatInput authorizedUser={props.authorization.username} send={send} />
        </div>
        <OnlineList onlineUsers={onlineUsers} users={users} />
      </div>
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
        <SignUp authorize={props.authorize} handleClose={handleClose} />
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
