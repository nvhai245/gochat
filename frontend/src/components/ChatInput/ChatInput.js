import React from 'react';
import "./ChatInput.scss";

export default function ChatInput(props) {
    return (
        <div className="ChatInput">
            <input className="userInput" placeholder="username" onChange={props.updateUsername} />
            <input className="textInput" onKeyDown={props.send} />
        </div>
    )
}
