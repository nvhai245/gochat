import React, { useState } from 'react';
import Paper from '@material-ui/core/Paper';
import SignupForm from '../SignupForm';
import { sendMsg } from '../../api';
import axios from 'axios';

function SignUp(props) {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const handleSignup = (event) => {
        event.preventDefault();
        var bodyFormData = new FormData();
        bodyFormData.set('username', props.formValue.username);
        bodyFormData.set('password', props.formValue.password);
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
                    let data = { username: props.formValue.username, isAdmin: false };
                    props.authorize(data);
                    props.handleClose();
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
        bodyFormData.set('username', props.formValue.username);
        bodyFormData.set('password', props.formValue.password);
        axios({
            method: 'post',
            url: 'http://localhost:8080/login',
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
                    let data = { username: props.formValue.username, isAdmin: false };
                    props.authorize(data);
                    props.handleClose();
                    window.location.reload();
                }
            })
            .catch(function (response) {
                //handle error
            });
    };
    const handleUsernameChange = (event) => {
        setUsername(event.target.value);
        setUsername(event.target.value);
    };
    const handlePasswordChange = (event) => {
        setPassword(event.target.value);
    };
    return (
        <Paper style={{ width: "60vw", height: "60vh", backgroundColor: "#EEEEEE", display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center" }}>
            <div style={{ display: "flex", flexDirection: "column", alignItems: "flex-start" }}>
                <SignupForm />
                <div>
                    <button onClick={props.formValue && handleLogin} style={{ marginRight: "2%" }} type="submit">Login</button>
                    <button onClick={props.formValue && handleSignup} type="submit">Signup</button>
                </div>
            </div>
        </Paper>
    )
}

export default SignUp
