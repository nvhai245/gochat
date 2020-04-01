import React from 'react';
import './InboxInput.scss';
import TextareaAutosize from '@material-ui/core/TextareaAutosize';
import SendIcon from '@material-ui/icons/Send';
import IconButton from '@material-ui/core/IconButton';

export default function InboxInput(props) {
    return (
        <div className="InboxInput">
            <TextareaAutosize autoFocus={true} id="inboxMessageInput" className="textInboxInput" onKeyDown={props.send} />
        </div>
    )
}
