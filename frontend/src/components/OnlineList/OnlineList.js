import React from 'react';
import './OnlineList.scss';

export default function OnlineList(props) {
    console.log("users is : ", props.onlineUsers);
    return (
        <div className="OnlineList">
            <h3>Current users</h3>
            {props.users && props.users.map(user =>
                <div style={{ display: "flex", alignItems: "flex-end", textAlign: "left" }}>
                    {props.onlineUsers.length !== 0 && props.onlineUsers.includes(user) ?
                        <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAABe0lEQVQ4T7WTvUsCcRjHv57h3ZloGiRJgUH5QkGDIbRagk1Baw0tDdIWIv0JIdIWDS0Ntboq+AcEkUOBYAZ1QacpvuTLeZ0ev4sarLMURPrBs/2ez/PhedFgxKcZMR//ALiwGKEb24ai+AFMA8hDo0miLV9ip1LvNVYbfCZTVNQ14/A7bU6rXsfSgtSSMrn7QpZ/SIKQUC9EDTifDDpmF8Jum8suimK3GE3TSL+kucc8F8Fu+fSnhRpwZomteXwBoS6whBBIkvQViqKAGWfE2+e7BPYqW/0BJ6Yrn3fdy7/yVLVRRUfudP8yNEPyJf4a+7XV/oBjU2xxeSmQE/Lsr/HKEKtcOYGD2gCDI0PQaDOHmXm9vaN8V4cCtLINTiqKERw2B/RgA0Z49FF2jvXr3AYrZdDSpCFL7xmhID2JSaRaIcShGmXvIulghgteahM2agUUpkBQRI7cIEXiKCEDoNm/B4AWwAQA+o8VlwG8AWgPAgx9Gv9wC0M6fABigowRmhtEUwAAAABJRU5ErkJggg=="
                        /> :
                        <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAABlklEQVQ4T63TP0gCURwH8N915ulRV0MQVINEcwQOgYr4ry5FqSHaHBqjprakFglyCgRLnBrcoqEGvYQ8BD2hwaWhKcKhgqRBTO/ZXXcvHJKUThB98/t93pcv70fAgIcYcB6GD3i9XspoNC6SJDmPMR4nCOJTUZRnhNADx3Ff3Yk7ErSGDQbDmt/v37BarR6apqcbjca7IAh3qVTqutls3nYjHYDP51sOBAI7DocjiBAa+X2NoiiV5/kkx3HxdDp9/zdFB8CybDAcDp8ghGYxxiBJEsiyDKqqgk6ne41GoweZTCapCbhcrr1IJHJaqVRG63URFOW7fVev18vx+Pk+z/MxTcBisWyHQkfHCEkz3WWJYuMtkYgdFovFC03AbDbb3O7VXbt9ZUtV1XYHGIOazXKXgpA7K5VKBU3AZDIZaJredDrZdZvNY2OYialarfqRz2cLuVzmRhTFq3K53NQEAEAPAAsMwyyRJDlHEMQYxriuKMoLQuhRluUnAKj3AkgAmAQA6p8v3mq0CgBSL6Dv1Rj+LvQb4QdhOqcRY27XNgAAAABJRU5ErkJggg=="
                        />
                    }
                    {user}
                </div>
            )}
        </div>
    )
}
