import React, {useEffect, useRef} from 'react';
import "./ChatHistory.scss";
import Message from '../Message/Message';
import { db } from '../../db';

export default function ChatHistory(props) {
    const mounted = useRef();
    useEffect(() => {
        if (mounted.current && props.table !== "") {
            props.setChatHistory(db.get(props.table).value());
        }
    }, [mounted.current]);
    const messages = props.chatHistory.map(msgData => <Message currentUser={props.currentUser} message={msgData} />)
    return (
            <div ref={mounted} className="ChatHistory" id="chat-history-scroll">
                {messages}
            </div>
    )
}
