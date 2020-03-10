import React from 'react';
import "./ChatHistory.scss";
import Message from '../Message/Message';

export default function ChatHistory(props) {
    const messages = props.chatHistory.map(msg => <Message message={msg.data} />);
    return (
        <div className="ChatHistory">
            <h2>Chat Room 01</h2>
            {messages}
        </div>
    )
}
