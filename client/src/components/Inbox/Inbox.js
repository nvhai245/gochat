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
import InboxInput from '../InboxInput';
import { db } from '../../db';
import { sendMsg } from "../../api";

const useStyles = makeStyles(theme => ({
    typography: {
        padding: theme.spacing(2),
    },
}));

export default function Inbox(props) {
    const send = event => {
        if (event.keyCode === 13 && event.target.value && event.target.value !== "" && !event.shiftKey) {
            event.preventDefault();
            console.log("count here", db.get(table + "count").value() + 1);
            let newMsg;
            newMsg = { count: db.get(table + "count").value() + 1, type: "chat", body: event.target.value, username: props.currentUser, receiver: [props.currentUser, props.user], table: table }
            sendMsg(JSON.stringify(newMsg));
            event.target.value = "";
            event.target.setAttribute('style', '');
        }
    }
    const sendMessage = () => {

    }
    const [table, setTable] = useState("");
    const [chatHistory, setChatHistory] = useState([]);
    useEffect(() => {
        let table1 = props.currentUser + "inboxto" + props.user;
        let table2 = props.user + "inboxto" + props.currentUser;
        if (db.has(table1).value()) {
            setTable(table1);
            let newMsg = { type: "getDifference", table: table1, username: props.currentUser };
            sendMsg(JSON.stringify(newMsg));
        } else if (db.has(table2).value()) {
            setTable(table2);
            let newMsg = { type: "getDifference", table: table2, username: props.currentUser };
            sendMsg(JSON.stringify(newMsg));
        } else {
            let newMsg = { type: "checkExist", body: table1, body3: table2, username: props.currentUser };
            sendMsg(JSON.stringify(newMsg));
        }
    }, []);
    useEffect(() => {
        let table1 = props.currentUser + "inboxto" + props.user;
        let table2 = props.user + "inboxto" + props.currentUser;
        if (props.newlyCreatedTable.includes(table1)) {
            setTable(table1);
        }
        if (props.newlyCreatedTable.includes(table2)) {
            setTable(table2);
        }
    }, [props.newlyCreatedTable]);
    useEffect(() => {
        let table1 = props.currentUser + "inboxto" + props.user;
        let table2 = props.user + "inboxto" + props.currentUser;
        if (props.mostRecentMsg.type === "chat" && props.mostRecentMsg.table === table) {
            if (db.has(props.mostRecentMsg.table).value()) {
                db.get(props.mostRecentMsg.table).push(props.mostRecentMsg).write();
                db.update(props.mostRecentMsg.table + 'count', n => n + 1).write();
                setChatHistory(prevState => ([...prevState, props.mostRecentMsg]));
            }
        }
        if (props.mostRecentMsg.type === "readdb" && props.mostRecentMsg.table === table) {
            db.get(props.mostRecentMsg.table).push(props.mostRecentMsg).write();
            db.update(props.mostRecentMsg.table + 'count', n => n + 1).write();
            if (chatHistory) {
                setChatHistory(prevState => ([...prevState, props.mostRecentMsg]));
            } else {
                setChatHistory([props.mostRecentMsg]);
            }
        }
        if (props.mostRecentMsg.type === "checkExist" && (props.mostRecentMsg.table === table1 || props.mostRecentMsg.table === table2)) {
            if (!db.has(props.mostRecentMsg.table).value()) {
                db.set(props.mostRecentMsg.table, []).write();
              db.set(props.mostRecentMsg.table + "count", 0).write();
              setTable(props.mostRecentMsg.table);
              let newMsg = { type: "getDifference", table: props.mostRecentMsg.table, username: props.currentUser };
            sendMsg(JSON.stringify(newMsg));
            }
        }
    }, [props.mostRecentMsg]);
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
                <InboxInput authorizedUser={props.currentUser} send={send} />
            </div>
        </div>
    )
}
