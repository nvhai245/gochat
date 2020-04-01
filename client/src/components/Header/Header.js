import React from "react";
import "./Header.scss";
import { sendMsg } from '../../api';
import axios from 'axios';
import Button from '@material-ui/core/Button';

function Header(props) {
  const logout = () => {
    var bodyFormData = new FormData();
    bodyFormData.set('username', props.username);
    axios({
      method: 'post',
      url: 'http://localhost:8080/logout',
      data: bodyFormData,
      withCredentials: true,
      headers: { 'Content-Type': 'multipart/form-data' }
    })
      .then(function (response) {
        //handle success
        if (response.status === 401) {
          alert("Unauthorized")
        }
        if (response.status === 200) {
          let newMsg;
          newMsg = { type: "logout", body: "", username: props.username }
          sendMsg(JSON.stringify(newMsg));
          let data = { username: "", isAdmin: false }
          props.authorize(data);
          window.location.reload();
        }
      })
      .catch(function (response) {
        //handle error
      });
  }
  return (
    <div style={{display: "flex"}} className="header">
      <h2 style={{margin: "0 auto"}}>Realtime Chat App</h2>
      {props.username !== "" && <Button size="large" style={{marginRight: "2%", color: "white"}} onClick={logout}><strong>Logout</strong></Button>}
    </div>
  )
};

export default Header;