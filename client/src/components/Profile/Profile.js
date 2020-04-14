import React, { useState, useEffect } from 'react';
import './Profile.scss';
import { makeStyles } from '@material-ui/core/styles';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemAvatar from '@material-ui/core/ListItemAvatar';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import Avatar from '@material-ui/core/Avatar';
import IconButton from '@material-ui/core/IconButton';
import FormGroup from '@material-ui/core/FormGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';
import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';
import MailOutlineIcon from '@material-ui/icons/MailOutline';
import EditIcon from '@material-ui/icons/Edit';
import TextField from '@material-ui/core/TextField';
import PhoneIphoneIcon from '@material-ui/icons/PhoneIphone';
import CakeIcon from '@material-ui/icons/Cake';
import FacebookIcon from '@material-ui/icons/Facebook';
import InstagramIcon from '@material-ui/icons/Instagram';
import { Paper } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import { sendMsg } from '../../api';
import 'date-fns';
import DateFnsUtils from '@date-io/date-fns';
import {
    MuiPickersUtilsProvider,
    KeyboardDatePicker,
} from '@material-ui/pickers';

const useStyles = makeStyles(theme => ({
    root: {
        flexGrow: 1,
        maxWidth: 752,
    },
    demo: {
        backgroundColor: theme.palette.background.paper,
    },
    title: {
        margin: theme.spacing(4, 0, 2),
    },
    username: {
        marginLeft: "10%",
        textTransform: "uppercase",
        marginBottom: "0px",
    }
}));

export default function Profile(props) {
    const [emailEditing, setEmailEditing] = useState(false);
    const [phoneEditing, setPhoneEditing] = useState(false);
    const [avatarEditing, setAvatarEditing] = useState(false);
    const [fbEditing, setFbEditing] = useState(false);
    const [instaEditing, setInstaEditing] = useState(false);
    const [birthdayEditing, setBirthdayEditing] = useState(false);
    const [selectedDate, setSelectedDate] = useState("");

    const enableEmailEdit = () => {
        setEmailEditing(true);
    }
    const enablePhoneEdit = () => {
        setPhoneEditing(true);
    }
    const enableAvatarEdit = () => {
        setAvatarEditing(true);
    }
    const enableFbEdit = () => {
        setFbEditing(true);
    }
    const enableInstaEdit = () => {
        setInstaEditing(true);
    }
    const enableBirthdayEdit = () => {
        setBirthdayEditing(true);
    }

    const editEmail = () => {
        let newEmail = document.getElementById("emailInput").value;
        let newMsg = { type: "updateEmail", username: props.user.username, body: newEmail };
        sendMsg(JSON.stringify(newMsg));
        props.setOpenBackdrop(true);
    }

    const editPhone = () => {
        let newPhone = document.getElementById("phoneInput").value;
        let newMsg = { type: "updatePhone", username: props.user.username, body: newPhone };
        sendMsg(JSON.stringify(newMsg));
        props.setOpenBackdrop(true);
    }

    const editAvatar = () => {
        let newAvatar = document.getElementById("avatarInput").value;
        let newMsg = { type: "updateAvatar", username: props.user.username, body: newAvatar };
        sendMsg(JSON.stringify(newMsg));
        props.setOpenBackdrop(true);
    }

    const editBirthday = () => {
        let newBirthday = document.getElementById("birthdayInput").value;
        let newMsg = { type: "updateBirthday", username: props.user.username, body: newBirthday };
        sendMsg(JSON.stringify(newMsg));
        props.setOpenBackdrop(true);
    }

    const editFb = () => {
        let newFb = document.getElementById("fbInput").value;
        let newMsg = { type: "updateFb", username: props.user.username, body: newFb };
        sendMsg(JSON.stringify(newMsg));
        props.setOpenBackdrop(true);
    }

    const editInsta = () => {
        let newInsta = document.getElementById("instaInput").value;
        let newMsg = { type: "updateInsta", username: props.user.username, body: newInsta };
        sendMsg(JSON.stringify(newMsg));
        props.setOpenBackdrop(true);
    }

    const handleDateChange = (date) => {
        setSelectedDate(date);
    }

    useEffect(() => {
        let element = document.getElementById("emailInput");
        element.focus();
    }, [emailEditing]);

    useEffect(() => {
        setSelectedDate(props.user.birthday)
    }, [props.user.birthday]);

    useEffect(() => {
        let element = document.getElementById("phoneInput");
        element.focus();
    }, [phoneEditing]);

    useEffect(() => {
        let element = document.getElementById("avatarInput");
        element.focus();
    }, [avatarEditing]);

    useEffect(() => {
        let element = document.getElementById("birthdayInput");
        element.focus();
    }, [birthdayEditing]);

    useEffect(() => {
        let element = document.getElementById("fbInput");
        element.focus();
    }, [fbEditing]);

    useEffect(() => {
        let element = document.getElementById("instaInput");
        element.focus();
    }, [instaEditing]);

    useEffect(() => {
        setEmailEditing(false);
        setPhoneEditing(false);
        setFbEditing(false);
        setInstaEditing(false);
        setBirthdayEditing(false);
        setAvatarEditing(false);
    }, [props.user]);

    const classes = useStyles();
    return (
        <div className="profileContainer">
            <div className="profile">
                <div className={classes.root}>
                    <Grid container spacing={2} style={{ justifyContent: "center" }}>
                        <Grid item xs={12} md={10}>
                            <Paper>
                                <div style={{ position: "relative", width: "3vw", height: "3vw" }}>
                                    <Avatar src={props.user.avatar} style={{ position: "absolute", width: "10vw", height: "10vw", top: "-5vw", left: "-5vw" }} />
                                </div>
                                <div className={classes.username}>
                                    <h2 style={{ marginBlockEnd: "unset", marginBlockStart: "unset" }}>{props.user.username}</h2>
                                </div>
                                <List>
                                    <ListItem>
                                        <ListItemAvatar>
                                            <Avatar>
                                                <MailOutlineIcon />
                                            </Avatar>
                                        </ListItemAvatar>
                                        <ListItemText>
                                            <TextField id="emailInput" label="Email" defaultValue={props.user.email} variant="outlined" fullWidth disabled={!emailEditing} />
                                        </ListItemText>
                                        {emailEditing ?
                                            <Button onClick={editEmail} style={{ marginLeft: "1rem" }} color="primary" variant="contained">
                                                Save
                                        </Button>
                                            :
                                            <IconButton edge="end" aria-label="edit" onClick={enableEmailEdit} >
                                                <EditIcon />
                                            </IconButton>
                                        }
                                    </ListItem>
                                    <ListItem>
                                        <ListItemAvatar>
                                            <Avatar src={props.user.avatar} />
                                        </ListItemAvatar>
                                        <ListItemText>
                                            <TextField id="avatarInput" label="Avatar" defaultValue={props.user.avatar} variant="outlined" fullWidth disabled={!avatarEditing} />
                                        </ListItemText>
                                        {avatarEditing ?
                                            <Button onClick={editAvatar} style={{ marginLeft: "1rem" }} color="primary" variant="contained">
                                                Save
                                    </Button>
                                            :
                                            <IconButton edge="end" aria-label="edit" onClick={enableAvatarEdit} >
                                                <EditIcon />
                                            </IconButton>
                                        }
                                    </ListItem>
                                    <ListItem>
                                        <ListItemAvatar>
                                            <Avatar>
                                                <PhoneIphoneIcon />
                                            </Avatar>
                                        </ListItemAvatar>
                                        <ListItemText>
                                            <TextField id="phoneInput" label="Phone" defaultValue={props.user.phone} variant="outlined" fullWidth disabled={!phoneEditing} />
                                        </ListItemText>
                                        {phoneEditing ?
                                            <Button onClick={editPhone} style={{ marginLeft: "1rem" }} color="primary" variant="contained">
                                                Save
                                    </Button>
                                            :
                                            <IconButton edge="end" aria-label="edit" onClick={enablePhoneEdit} >
                                                <EditIcon />
                                            </IconButton>
                                        }
                                    </ListItem>
                                    <ListItem>
                                        <ListItemAvatar>
                                            <Avatar>
                                                <CakeIcon />
                                            </Avatar>
                                        </ListItemAvatar>
                                        <ListItemText>
                                            <MuiPickersUtilsProvider utils={DateFnsUtils}>
                                                <KeyboardDatePicker
                                                    margin="normal"
                                                    id="birthdayInput"
                                                    label="Birthday"
                                                    format="MM/dd/yyyy"
                                                    value={selectedDate}
                                                    onChange={handleDateChange}
                                                    inputVariant="outlined"
                                                    fullWidth 
                                                    disabled={!birthdayEditing} 
                                                    KeyboardButtonProps={{
                                                        'aria-label': 'change date',
                                                    }}
                                                />
                                            </MuiPickersUtilsProvider>
                                        </ListItemText>
                                        {birthdayEditing ?
                                            <Button onClick={editBirthday} style={{ marginLeft: "1rem" }} color="primary" variant="contained">
                                                Save
                                    </Button>
                                            :
                                            <IconButton edge="end" aria-label="edit" onClick={enableBirthdayEdit} >
                                                <EditIcon />
                                            </IconButton>
                                        }
                                    </ListItem>
                                    <ListItem>
                                        <ListItemAvatar>
                                            <Avatar>
                                                <FacebookIcon />
                                            </Avatar>
                                        </ListItemAvatar>
                                        <ListItemText>
                                            <TextField id="fbInput" label="Facebook" defaultValue={props.user.fb} variant="outlined" fullWidth disabled={!fbEditing} />
                                        </ListItemText>
                                        {fbEditing ?
                                            <Button onClick={editFb} style={{ marginLeft: "1rem" }} color="primary" variant="contained">
                                                Save
                                    </Button>
                                            :
                                            <IconButton edge="end" aria-label="edit" onClick={enableFbEdit} >
                                                <EditIcon />
                                            </IconButton>
                                        }
                                    </ListItem><ListItem>
                                        <ListItemAvatar>
                                            <Avatar>
                                                <InstagramIcon />
                                            </Avatar>
                                        </ListItemAvatar>
                                        <ListItemText>
                                            <TextField id="instaInput" label="Instagram" defaultValue={props.user.insta} variant="outlined" fullWidth disabled={!instaEditing} />
                                        </ListItemText>
                                        {instaEditing ?
                                            <Button onClick={editInsta} style={{ marginLeft: "1rem" }} color="primary" variant="contained">
                                                Save
                                    </Button>
                                            :
                                            <IconButton edge="end" aria-label="edit" onClick={enableInstaEdit} >
                                                <EditIcon />
                                            </IconButton>
                                        }
                                    </ListItem>
                                </List>
                            </Paper>
                        </Grid>
                    </Grid>
                </div>
            </div>
        </div>
    )
}
