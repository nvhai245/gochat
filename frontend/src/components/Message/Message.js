import React, { useEffect, useState } from 'react';
import "./Message.scss";

export default function Message(props) {
    console.log(props.message);
    const [message, setMessage] = useState({});
    useEffect(() => {
        let temp = JSON.parse(props.message);
        setMessage(temp)
    }, [])
    return (
        <div className="MessageBox">
                {message.username !== "admin" && <div className="User">{message.username + ":"}</div>}
                {message.username === "" && <div className="User">unknown: </div>}
            <div className="Message">
                {message.body}
            </div>
        </div>
    )
}
