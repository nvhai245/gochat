import React, {useState} from 'react';
import "./ChatInput.scss";
import TextareaAutosize from '@material-ui/core/TextareaAutosize';
import SendIcon from '@material-ui/icons/Send';
import IconButton from '@material-ui/core/IconButton';

export default function ChatInput(props) {
    return (
        <div className="ChatInput">
            <input disabled value={props.authorizedUser} className="userInput" placeholder="username" />
            <TextareaAutosize id="textMessageInput" className="textInput" onKeyDown={props.send} />
            <IconButton onClick={props.sendMessage} iconStyle={{fontSize: "20px"}}>
                <SendIcon style={{color: "#44A4BE", fontSize: "2rem"}} size="large" />
            </IconButton>
        </div>
    )
}
