import React from 'react';
import "./ChatInput.scss";

export default function ChatInput(props) {
    return (
        <div className="ChatInput">
            <input disabled value={props.authorizedUser} className="userInput" placeholder="username"/>
            <input className="textInput" onKeyDown={props.send} />
        </div>
    )
}
