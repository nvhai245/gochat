package websocket

import (
    "fmt"
    "log"
	_ "sync"
	"encoding/json"

    "github.com/gorilla/websocket"
)

type Client struct {
    ID   string
    Conn *websocket.Conn
    Pool *Pool
}

type Message struct {
    Type int    `json:"type"`
	Body string `json:"body"`
	Username string `json:"username"`
}

func (c *Client) Read() {
    defer func() {
        c.Pool.Unregister <- c
        c.Conn.Close()
    }()

    for {
		message := Message{}
        messageType, p, err := c.Conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
		}
		json.Unmarshal([]byte(string(p)), &message)
		message.Type = messageType
        c.Pool.Broadcast <- message
        fmt.Printf("Message Received: %+v\n", message)
    }
}