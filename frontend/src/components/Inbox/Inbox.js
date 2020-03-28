import React, { useState, useEffect, useRef } from 'react';
import './Inbox.scss';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import Paper from '@material-ui/core/Paper';
import shorid from 'shortid';
import CloseIcon from '@material-ui/icons/Close';
import Divider from '@material-ui/core/Divider';
import FaceIcon from '@material-ui/icons/Face';
import CallIcon from '@material-ui/icons/Call';
import VideocamIcon from '@material-ui/icons/Videocam';
import IconButton from '@material-ui/core/IconButton';
import ChatHistory from '../ChatHistory';
import ChatInput from '../ChatInput';
import { db } from '../../db';

const useStyles = makeStyles(theme => ({
    typography: {
        padding: theme.spacing(2),
    },
}));

export default function Inbox(props) {
    const send = () => {

    }
    const sendMessage = () => {
        
    }
    const [chatHistory, setChatHistory] = useState([]);
    const [table, setTable] = useState("");
    useEffect(() => {
        let table1 = props.currentUser + "inboxto" + props.user;
        let table2 = props.user + "inboxto" + props.currentUser;
        if (db.has(table1).value()) {
            setTable(table1);
            let newMsg = { type: "getDifference", table: table1, username: props.currentUser };
            props.sendMsg(JSON.stringify(newMsg));
        } else if (db.has(table2).value()) {
            setTable(table2);
            let newMsg = { type: "getDifference", table: table2, username: props.currentUser };
            props.sendMsg(JSON.stringify(newMsg));
        } else {
            let newMsg = { type: "checkExist", body: table1, body3: table2, username: props.currentUser };
            props.sendMsg(JSON.stringify(newMsg));
        }
    });
    useEffect(() => {
        let table1 = props.currentUser + "inboxto" + props.user;
        let table2 = props.user + "inboxto" + props.currentUser;
        if (props.newlyCreatedTable.includes(table1)) {
            setTable(table1);
        }
        if (props.newlyCreatedTable.includes(table2)) {
            setTable(table2);
        }
    }, [props.newlyCreatedTable])
    return (
        <div className="inbox">
            <div className="toolbar">
                <IconButton>
                    <FaceIcon fontSize="large" />
                </IconButton>
                <Typography>{props.user}</Typography>
                <IconButton edge="end" style={{ marginLeft: "auto" }} >
                    <CallIcon fontSize="small" />
                </IconButton>
                <IconButton edge="end">
                    <VideocamIcon fontSize="small" />
                </IconButton>
                <IconButton onClick={() => props.addInboxList(props.user)}>
                    <CloseIcon fontSize="small" />
                </IconButton>
            </div>
            <Divider />
            <div className="inboxHistory">
                <ChatHistory table={table} currentUser={props.currentUser} setChatHistory={setChatHistory} chatHistory={chatHistory} />
            </div>
            <Divider />
            <div className="inboxInput">
                <ChatInput authorizedUser={props.currentUser} send={send} sendMessage={sendMessage} />
            </div>
        </div>
    )
}
