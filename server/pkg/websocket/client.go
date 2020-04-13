package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	_ "sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
	Pool     *Pool
}

type Message struct {
	ID       string
	Count    int64    `json:"count"`
	Type     string   `json:"type"`
	Body     string   `json:"body"`
	Username string   `json:"username"`
	Body2    []string `json:"body2"`
	Body3    string   `json:"body3"`
	Table    string   `json:"table"`
	Receiver []string `json:"receiver"`
	Deleted  bool     `json:"deleted"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		message := Message{}
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		json.Unmarshal([]byte(string(p)), &message)
		log.Println("message type is ", message.Type)
		message.ID = c.ID
		if message.Type == "chat" {
			c.Pool.Broadcast <- message
			fmt.Printf("Message Received: %+v\n", message)
		} else {
			c.Pool.Broadcast <- message
		}
	}
}
