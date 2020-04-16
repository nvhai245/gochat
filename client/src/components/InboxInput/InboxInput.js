import React, { useState, useEffect, useRef } from 'react';
import './InboxInput.scss';
import TextareaAutosize from '@material-ui/core/TextareaAutosize';
import SendIcon from '@material-ui/icons/Send';
import IconButton from '@material-ui/core/IconButton';
import Picker from 'emoji-picker-react';
import Popover from '@material-ui/core/Popover';
import InsertEmoticonIcon from '@material-ui/icons/InsertEmoticon';
import ImageUpload from '../ImageUpload';
import HighlightOffIcon from '@material-ui/icons/HighlightOff';
import VideoUpload from '../VideoUpload';

export default function InboxInput(props) {
    const mounted = useRef();
    const [chosenEmoji, setChosenEmoji] = useState(null);
    const [anchorEl, setAnchorEl] = React.useState(null);

    const handleClick = event => {
        setAnchorEl(props.box.current);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const open = Boolean(anchorEl);
    const id = open ? 'simple-popover' : undefined;

    const onEmojiClick = (event, emojiObject) => {
        setChosenEmoji(emojiObject);
    }
    useEffect(() => {
        if (chosenEmoji) {
            mounted.current.value = mounted.current.value + chosenEmoji.emoji;
        }
    }, [chosenEmoji]);
    return (
        <div className="InboxInput">
            <div style={{display: "flex", alignItems: "center"}}>
            <div>
                <TextareaAutosize ref={mounted} autoFocus={true} id="inboxMessageInput" className="textInboxInput" onKeyDown={props.send} />
            </div>
            <div style={{ marginLeft: "0.5rem", minHeight: "1rem" }}>
                <IconButton size="small" edge="start" aria-describedby={id} variant="contained" color="primary" onClick={handleClick}>
                    <InsertEmoticonIcon />
                </IconButton>
            </div>
            <ImageUpload table={props.table} setImg={props.setImg} setLoading={props.setLoading} />
            <VideoUpload table={props.table} setImg={props.setImg} setLoading={props.setLoading} />
            </div>
            <div style={{marginTop: "0.3rem"}}>
            {props.loading > 0 &&
                <div style={{ width: '100px', height: '100px', display: "flex", flexDirection: "column", position: "relative", marginTop: "0.5rem" }} className="progress">
                    <img src={props.img} style={{ height: "100%", width: "100%", objectFit: "contain"}} alt="" />
                    <IconButton style={{ position: "absolute", top: "0", left: "0" }} size="small" onClick={props.deleteImage}>
                        <HighlightOffIcon size="small" />
                    </IconButton>
                    <progress style={{ width: "100px" }} id="file" value={`${props.loading}`} max="100"></progress>
                </div>
            }
            </div>
            <Popover
                id={id}
                open={open}
                anchorEl={anchorEl}
                onClose={handleClose}
                anchorOrigin={{
                    vertical: 'top',
                    horizontal: 'left',
                }}
                transformOrigin={{
                    vertical: 'top',
                    horizontal: 'center',
                }}
                marginThreshold={100}
            >
                <div style={{ display: "flex" }}>
                    <Picker onEmojiClick={onEmojiClick} />
                </div>
            </Popover>
        </div>
    )
}
