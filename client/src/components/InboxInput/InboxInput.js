import React, { useState, useEffect, useRef } from 'react';
import './InboxInput.scss';
import TextareaAutosize from '@material-ui/core/TextareaAutosize';
import SendIcon from '@material-ui/icons/Send';
import IconButton from '@material-ui/core/IconButton';
import Picker from 'emoji-picker-react';
import Popover from '@material-ui/core/Popover';
import InsertEmoticonIcon from '@material-ui/icons/InsertEmoticon';

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
            <div>
                <TextareaAutosize ref={mounted} autoFocus={true} id="inboxMessageInput" className="textInboxInput" onKeyDown={props.send} />
            </div>
            <div style={{ marginLeft: "0.5rem", minHeight: "1rem" }}>
                <IconButton size="small" edge="start" aria-describedby={id} variant="contained" color="primary" onClick={handleClick}>
                    <InsertEmoticonIcon />
                </IconButton>
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
