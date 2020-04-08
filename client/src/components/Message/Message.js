import React, { useEffect, useState, useRef } from 'react';
import "./Message.scss";
import IconButton from '@material-ui/core/IconButton';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import {sendMsg} from '../../api';

export default function Message(props) {
  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);
  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };
  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleDelete = () => {
    let newMsg;
    newMsg = { count: props.message.count, type: "delete", username: props.currentUser, receiver: props.message.receiver, table: props.message.table }
    sendMsg(JSON.stringify(newMsg));
    setAnchorEl(null);
  }

  const handleRestore = () => {
    let newMsg;
    newMsg = { count: props.message.count, type: "restore", username: props.currentUser, receiver: props.message.receiver, table: props.message.table }
    sendMsg(JSON.stringify(newMsg));
    setAnchorEl(null);
    setAnchorEl(null);
  }

    const mounted = useRef();               
    const [message, setMessage] = useState({});
    useEffect(() => {
        let temp = props.message;
        setMessage(temp);
    }, [])
    useEffect(() => {
        if (!mounted.current) {
            console.log("not mounted yet");
          } else {
            // do componentDidUpate logic
            const element = document.getElementById("chat-history-scroll")
            if (element && element.scroll) {
                element.scrollTop = element.scrollHeight;
            }
            props.mounted.current.scrollTop = props.mounted.current.scrollHeight;
          }
    }, [mounted.current]);
    return (
        message.username === props.currentUser ?
            <div ref={mounted} className="MessageBox" style={{ display: "flex", justifyContent: "flex-end" }}>
                <div className="messageContainer" style={{ color: "white", display: "flex", justifyContent: "flex-end", maxWidth: "80%"}}>
                    <div className={message.deleted ? "DeletedMessage" : "UserMessage"}>
                    { message.deleted ? <i>Message has been deleted</i> : message.body }
                    </div>
                    <div className="messageOptions">
      <IconButton
        aria-label="more"
        aria-controls="message-menu"
        aria-haspopup="true"
        onClick={handleClick}
        size="small"
      >
        <MoreVertIcon />
      </IconButton>
      <Menu
        id="long-menu"
        anchorEl={anchorEl}
        keepMounted
        open={open}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'center',
          horizontal: 'right',
        }}
      >
          {message.deleted ? 
          <MenuItem key={"restoreMessage" + message.count + message.table} onClick={handleRestore}>
          Restore
          </MenuItem>
          :
          <MenuItem key={"deleteMessage" + message.count + message.table} onClick={handleDelete}>
          Delete
          </MenuItem>
          }
      </Menu>
    </div>
                </div>
            </div> :
            <div ref={mounted} className="MessageBox" style={{ justifyContent: "flex-start" }}>
                <div className="messageContainer" style={{ color: "white", maxWidth: "80%", display: "flex", alignItems: "flex-end" }}>
                    {message.username !== "admin" && <div className="User" style={{ backgroundColor: "#D6DEBD", borderRadius: "50%", height: "40px", width: "40px", marginRight: "0.2rem", padding: "auto auto" }}>{message.username + ":"}</div>}
                    {message.username === "" && <div className="User">unknown: </div>}
                    <div className={message.deleted ? "DeletedMessage" : "GuestMessage"}>
                        { message.deleted ? <i>Message has been deleted</i> : message.body }
                    </div>
                </div>
            </div>

    )
}
