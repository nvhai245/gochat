import React, { useState } from 'react';
import Paper from '@material-ui/core/Paper';
import { sendMsg } from '../../api';

export default function SignUp(props) {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const handleSignup = (event) => {
        event.preventDefault();
        let newMsg;
        newMsg = { type: "register", body: password, username: username }
        sendMsg(JSON.stringify(newMsg));
        props.handleClose();
    };
    const handleLogin = (event) => {
        event.preventDefault();
        let newMsg;
        newMsg = { type: "login", body: password, username: username }
        sendMsg(JSON.stringify(newMsg));
        props.handleClose();
    };
    const handleUsernameChange = (event) => {
        setUsername(event.target.value);
        props.setUsername(event.target.value);
    };
    const handlePasswordChange = (event) => {
        setPassword(event.target.value);
    };
    return (
        <Paper style={{ width: "60vw", height: "60vh", backgroundColor: "#EEEEEE", display: "flex", justifyContent: "center", alignItems: "center" }}>
            <form>
                <input
                    style={{ maxWidth: "30%", marginRight: "2%" }}
                    type="text"
                    className="usernameInput"
                    placeholder="Username"
                    onChange={handleUsernameChange}
                />
                <input type="password" className="passwordInput"
                    placeholder="Password"
                    style={{ marginRight: "2%" }}
                    onChange={handlePasswordChange}
                />
                <button onClick={handleLogin} style={{ marginRight: "2%" }} type="submit">Login</button>
                <button onClick={handleSignup} type="submit">Signup</button>
            </form>
        </Paper>
    )
}
