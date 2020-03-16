import React, { useState } from 'react';
import Paper from '@material-ui/core/Paper';
import { sendMsg } from '../../api';
import axios from 'axios';

export default function SignUp(props) {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const handleSignup = (event) => {
        event.preventDefault();
        var bodyFormData = new FormData();
        bodyFormData.set('username', username);
        bodyFormData.set('password', password);
        axios({
            method: 'post',
            url: 'http://localhost:8080/signup',
            data: bodyFormData,
            withCredentials: true,
            headers: { 'Content-Type': 'multipart/form-data'}
        })
            .then(function (response) {
                //handle success
                if (response.status === 401) {
                    alert("Signup failed")
                }
                if (response.status === 200) {
                    let data = { username: username, isAdmin: false };
                    props.authorize(data);
                }
                console.log(response);
            })
            .catch(function (response) {
                //handle error
                console.log(response);
            });
        props.handleClose();
    };
    const handleLogin = (event) => {
        event.preventDefault();
        var bodyFormData = new FormData();
        bodyFormData.set('username', username);
        bodyFormData.set('password', password);
        axios({
            method: 'post',
            url: 'http://localhost:8080/login',
            data: bodyFormData,
            withCredentials: true,
            headers: { 'Content-Type': 'multipart/form-data'}
        })
            .then(function (response) {
                //handle success
                if (response.status === 401) {
                    alert("Signup failed")
                }
                if (response.status === 200) {
                    let data = { username: username, isAdmin: false };
                    props.authorize(data);
                }
                console.log(response);
            })
            .catch(function (response) {
                //handle error
                console.log(response);
            });
        props.handleClose();
        window.location.reload();
    };
    const handleUsernameChange = (event) => {
        setUsername(event.target.value);
        setUsername(event.target.value);
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
