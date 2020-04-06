import React, { useState } from 'react';
import Paper from '@material-ui/core/Paper';
import { sendMsg } from '../../api';
import axios from 'axios';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';

function SignUp(props) {
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
            headers: { 'Content-Type': 'multipart/form-data' }
        })
            .then(function (response) {
                //handle success
                if (response.status === 401) {
                    alert("Signup failed")
                }
                if (response.status === 200) {
                    window.location.reload();
                }
            })
            .catch(function (response) {
                //handle error
            });
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
            headers: { 'Content-Type': 'multipart/form-data' }
        })
            .then(function (response) {
                //handle success
                if (response.status === 200) {
                    window.location.reload();
                } else {
                    alert("Wrong username or password");
                }
            })
            .catch(function (response) {
                //handle error
            });
    };
    const handleUsernameChange = (event) => {
        setUsername(event.target.value);
    };
    const handlePasswordChange = (event) => {
        setPassword(event.target.value);
    };
    return (
        <Paper style={{ width: "60vw", height: "60vh", margin: "auto auto", backgroundColor: "#EEEEEE", display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center" }}>
            <div style={{ display: "flex", flexDirection: "column", alignItems: "flex-start"}}>
                        <TextField autoFocus style={{marginBottom: "3%"}} autoComplete="current-username" variant="outlined" label="Username" onChange={handleUsernameChange} name="username" type="text" />
                        <TextField style={{marginBottom: "3%"}} autoComplete="current-password" variant="outlined" label="Password" onChange={handlePasswordChange} name="password" type="password" />
                <div style={{display: "flex"}}>
                    <Button color="primary" variant="contained" onClick={handleLogin} style={{ marginRight: "2%" }}>Login</Button>
                    <Button color="secondary" variant="contained" onClick={handleSignup}>Signup</Button>
                </div>
            </div>
        </Paper>
    )
}

export default SignUp
