package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	_ "sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Username string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	ID string
	Type     string `json:"type"`
	Body     string `json:"body"`
	Username string `json:"username"`
}

func Check(token string) (authorized bool) {
	return true
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
			authorized := Check(message.Body)
			if authorized == true {
				c.Pool.Broadcast <- message
				fmt.Printf("Message Received: %+v\n", message)
			} else {
                log.Println("Unauthorized")
            }
		} else {
			c.Pool.Broadcast <- message
		}
	}
}
