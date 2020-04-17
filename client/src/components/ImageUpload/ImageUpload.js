import React, { useState } from 'react';
import './ImageUpload.scss';
import axios from 'axios';
import IconButton from '@material-ui/core/IconButton';
import ImageIcon from '@material-ui/icons/Image';

export default function ImageUpload(props) {
    const handleChange = e => {
        const reader = new FileReader(); // eslint-disable-line
        const file = e.target.files[0];
        reader.onloadend = () => {
            const formData = new FormData(); // eslint-disable-line
            formData.append('file', file);
            formData.append('tags', 'go-chat');
            formData.append('upload_preset', 'ehfxpga3');
            formData.append('api_key', '557364961884294');
            formData.append('timestamp', (Date.now() / 1000) | 0); // eslint-disable-line
            axios
                .post('https://api.cloudinary.com/v1_1/nvhai245/image/upload', formData, {
                    onUploadProgress: ProgressEvent => {
                        const load = (ProgressEvent.loaded / ProgressEvent.total) * 100;
                        console.log(load);
                        props.setLoading((ProgressEvent.loaded / ProgressEvent.total) * 100);
                    },
                    headers: { 'X-Requested-With': 'XMLHttpRequest' },
                })
                .then(res => {
                    props.setImg(res.data.secure_url);
                })
                .catch(error => console.log(error));
            console.log(reader.result);
        };
        if (file) {
            reader.readAsDataURL(file);
        }
    };
    return (
        <div style={{ display: "flex", marginLeft: props.mgl ? props.mgl : "0px" }}>
            <input
                accept="image/*"
                id={"image-" + props.table}
                multiple
                type="file"
                style={{ display: "none" }}
                onChange={handleChange}
            />
            <label htmlFor={"image-" + props.table}>
                <IconButton edge="start" raised component="span">
                    <ImageIcon />
                </IconButton>
            </label>
        </div>
    )
}
