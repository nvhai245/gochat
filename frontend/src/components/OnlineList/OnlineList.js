import React from 'react';
import './OnlineList.scss';

export default function OnlineList(props) {
    console.log("users is : ", props.onlineUsers);
    return (
        <div className="OnlineList" style={{ textAlign: "left" }}>
            {props.users && props.users.map(user =>
                <div style={{ display: "flex", alignItems: "flex-end" }}>
                    {props.onlineUsers.length !== 0 && props.onlineUsers.includes(user) ?
                        <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAABe0lEQVQ4T7WTvUsCcRjHv57h3ZloGiRJgUH5QkGDIbRagk1Baw0tDdIWIv0JIdIWDS0Ntboq+AcEkUOBYAZ1QacpvuTLeZ0ev4sarLMURPrBs/2ez/PhedFgxKcZMR//ALiwGKEb24ai+AFMA8hDo0miLV9ip1LvNVYbfCZTVNQ14/A7bU6rXsfSgtSSMrn7QpZ/SIKQUC9EDTifDDpmF8Jum8suimK3GE3TSL+kucc8F8Fu+fSnhRpwZomteXwBoS6whBBIkvQViqKAGWfE2+e7BPYqW/0BJ6Yrn3fdy7/yVLVRRUfudP8yNEPyJf4a+7XV/oBjU2xxeSmQE/Lsr/HKEKtcOYGD2gCDI0PQaDOHmXm9vaN8V4cCtLINTiqKERw2B/RgA0Z49FF2jvXr3AYrZdDSpCFL7xmhID2JSaRaIcShGmXvIulghgteahM2agUUpkBQRI7cIEXiKCEDoNm/B4AWwAQA+o8VlwG8AWgPAgx9Gv9wC0M6fABigowRmhtEUwAAAABJRU5ErkJggg=="
                        /> :
                        <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAABaklEQVQ4T7WTMUtCYRSGn3sDr14sbEmKCLGGIFCoCMJZMByC1hpaGqItRPoJEdEWDS0NtbYW+AOCSFBcNKgriUUmGZrmp/LdqKHykkFYB77xPN973vMehS5L6bKfvwccQp8NFk0IAoPAnQKxBhwtQdmquE3BW7MK216fL+j1+92a06nVKxVxnUzeG6lUTELECmkDHMDq2MREdMTv91Sr1Y/P7HY7RiKRvclktpZh76uKNsA+HAfC4VCxVHJIKRFCvD/TNOnV9ZerePx0BRY6AnbhLBAKzeQMQ62USrSazU8Vui4f8vnzNZjtCNiB40mfL1TO5RxWs1qK8pJ9fDxd/0nBJqwOu1zRUU3zyEbjg2ECl7VatiDE1sZPHsxB39TbFhyO4LjN5naqqlaRUqTr9XtDiFgcIieWVVqDZOuH8RmYH4JpFQYkFG7hIg4nRUgDzx09AHoAF6B9E/EW8AR8zgb/EOXfHlfXx/QKQrt/EaXk3vAAAAAASUVORK5CYII="
                        />
                    }
                    {user}
                </div>
            )}
        </div>
    )
}
