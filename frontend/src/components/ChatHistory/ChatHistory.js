import React from 'react';
import "./ChatHistory.scss";
import Message from '../Message/Message';

export default function ChatHistory(props) {
    const messages = props.chatHistory.map(msg => <Message currentUser={props.currentUser} message={msg.data} />);
    return (
        <React.Fragment>
            <div className="ChatHistory" id="chat-history-scroll">
                <h2>Chat Room 01</h2>
                {messages}
            </div>
        </React.Fragment>
    )
}
