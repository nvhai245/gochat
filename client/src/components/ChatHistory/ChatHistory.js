import React, {useEffect, useRef} from 'react';
import "./ChatHistory.scss";
import Message from '../Message/Message';
import { db } from '../../db';

export default function ChatHistory(props) {
    if (props.table !== "all") {
        console.log(props.chatHistory);
    }
    const mounted = useRef();
    useEffect(() => {
        if (props.table !== "") {
            props.setChatHistory(db.get(props.table).value());
        }
    }, [props.table]);
    useEffect(() => {
        if (props.table !== "") {
            props.setChatHistory(db.get(props.table).value());
        }
    }, [mounted.current]);
    const messages = props.chatHistory ? props.chatHistory.map(msgData => <Message mounted={mounted} currentUser={props.currentUser} message={msgData} />) : ""
    return (
            <div ref={mounted} className="ChatHistory" id="chat-history-scroll">
                {messages}
            </div>
    )
}
