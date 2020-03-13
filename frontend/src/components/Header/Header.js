import React from "react";
import "./Header.scss";
import { sendMsg } from '../../api';

function Header(props) {
  const logout = () => {
    let newMsg;
        newMsg = { type: "logout", body: "", username: props.username }
        sendMsg(JSON.stringify(newMsg));
    let data = {username: "", isAdmin: false}
    props.authorize(data);
    window.location.reload();
  }
  return (
    <div className="header">
      <h2>Realtime Chat App</h2>
      {props.username !== "" && <button onClick={logout}>Logout</button>}
    </div>
  )
};

export default Header;