import React from 'react';
import "./ChatInput.scss";
import TextareaAutosize from '@material-ui/core/TextareaAutosize';

export default function ChatInput(props) {
    return (
        <div className="ChatInput">
            <input disabled value={props.authorizedUser} className="userInput" placeholder="username" />
            <TextareaAutosize className="textInput" onKeyDown={props.send} />
        </div>
    )
}
