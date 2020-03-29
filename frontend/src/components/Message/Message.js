import React, { useEffect, useState, useRef } from 'react';
import "./Message.scss";

export default function Message(props) {
    const mounted = useRef();               
    const [message, setMessage] = useState({});
    useEffect(() => {
        let temp = props.message;
        setMessage(temp);
    }, [])
    useEffect(() => {
        if (!mounted.current) {
            console.log("not mounted yet");
          } else {
            // do componentDidUpate logic
            const element = document.getElementById("chat-history-scroll")
            if (element && element.scroll) {
                element.scrollTop = element.scrollHeight;
            }
            props.mounted.current.scrollTop = props.mounted.current.scrollHeight;
          }
    }, [mounted.current]);
    return (
        message.username === props.currentUser ?
            <div ref={mounted} className="MessageBox" style={{ display: "flex", justifyContent: "flex-end" }}>
                <div className="messageContainer" style={{ color: "white", display: "flex", justifyContent: "flex-end", maxWidth: "50%", marginRight: "1rem" }}>
                    <div className="Message" style={{ backgroundColor: "#44A4BE", borderRadius: "1rem", padding: "10px 10px" }}>
                        {message.body}
                    </div>
                </div>
            </div> :
            <div ref={mounted} className="MessageBox" style={{ justifyContent: "flex-start" }}>
                <div className="messageContainer" style={{ color: "white", maxWidth: "50%", display: "flex", alignItems: "flex-end" }}>
                    {message.username !== "admin" && <div className="User" style={{ backgroundColor: "#D6DEBD", borderRadius: "50%", height: "40px", width: "40px", marginRight: "0.2rem", padding: "auto auto" }}>{message.username + ":"}</div>}
                    {message.username === "" && <div className="User">unknown: </div>}
                    <div className="Message" style={{ backgroundColor: "#DDDEEC", color: "black", borderRadius: "1rem", padding: "10px 10px" }}>
                        {message.body}
                    </div>
                </div>
            </div>

    )
}
