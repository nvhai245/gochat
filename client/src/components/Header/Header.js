import React from "react";
import "./Header.scss";
import { sendMsg } from '../../api';
import axios from 'axios';
import Button from '@material-ui/core/Button';
import ChatIcon from '@material-ui/icons/Chat';
import IconButton from '@material-ui/core/IconButton';
import Badge from '@material-ui/core/Badge';
import Popover from '@material-ui/core/Popover';
import Typography from '@material-ui/core/Typography';
import Divider from '@material-ui/core/Divider';
import FaceIcon from '@material-ui/icons/Face';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';

function Header(props) {
  const [anchorEl, setAnchorEl] = React.useState(null);

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleListClick = (username) => {
    props.addInboxList(username);
    setAnchorEl(null);
  }

  const open = Boolean(anchorEl);
  const id = open ? 'simple-popover' : undefined;
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
          let data = { username: undefined, isAdmin: false };
          props.authorize(data);
        }
      })
      .catch(function (response) {
        //handle error
      });
  }
  return (
    <div style={{ display: "flex" }} className="header">
      <h2 style={{ margin: "0 auto" }}>Realtime Chat App</h2>
      {props.username &&
        <div style={{ display: "flex", alignItems: "center" }}>
          <IconButton onClick={handleClick}>
            <Badge badgeContent={props.noti.length} color="primary">
              <ChatIcon style={{ color: "white" }} />
            </Badge>
          </IconButton>
          {props.noti.length > 0 &&
            <Popover
            id={id}
            open={open}
            anchorEl={anchorEl}
            onClose={handleClose}
            anchorOrigin={{
              vertical: 'center',
              horizontal: 'right',
            }}
            transformOrigin={{
              vertical: 'top',
              horizontal: 'center',
            }}
          >
              <List component="nav" style={{ display: "flex", flexDirection: "column", padding: "0.5rem 0.5rem" }}>
                {props.noti.map(message =>
                  <div style={{ display: "flex", alignItems: "center" }}>
                    <ListItem button onClick={() => handleListClick(message.username)}>
                      <FaceIcon fontSize="large" />
                      <div style={{ display: "flex", flexDirection: "column" }}>
                        <strong>{message.username}</strong>
                        <i>{message.body}</i>
                      </div>
                    </ListItem>
                    <Divider />
                  </div>
                )}
              </List>
          </Popover>
      }
          <Button size="large" style={{ marginRight: "2%", color: "white" }} onClick={logout}><strong>Logout</strong></Button>
        </div>
      }
    </div>
  )
};

export default Header;