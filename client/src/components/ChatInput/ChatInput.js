import React, { useState } from 'react';
import "./ChatInput.scss";
import TextareaAutosize from '@material-ui/core/TextareaAutosize';
import SendIcon from '@material-ui/icons/Send';
import IconButton from '@material-ui/core/IconButton';
import ImageUpload from '../ImageUpload';
import HighlightOffIcon from '@material-ui/icons/HighlightOff';

export default function ChatInput(props) {
    return (
        <div className="ChatInput">
            <input disabled value={props.authorizedUser} className="userInput" placeholder="username" />
            <div style={{width: "87%", display: "flex", flexDirection: "column", justifyContents: "flex-start"}}>
            <TextareaAutosize id="textMessageInput" className="textInput" onKeyDown={props.send} />
            {props.loading > 0 &&
                <div style={{ width: '100px', height: '100px', display: "flex", flexDirection: "column", position: "relative", marginTop: "0.5rem" }} className="progress">
                    <img src={props.img} style={{ height: "100%", width: "100%", objectFit: "contain"}} alt="" />
                    <IconButton style={{ position: "absolute", top: "0", right: "0" }} size="small" onClick={props.deleteImage}>
                        <HighlightOffIcon size="small" />
                    </IconButton>
                    <progress style={{ width: "100px" }} id="file" value={`${props.loading}`} max="100"></progress>
                </div>
            }
            </div>
            <ImageUpload table="all" key="all" setImg={props.setImg} setLoading={props.setLoading} />
        </div>
    )
}
