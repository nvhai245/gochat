package websocket

import (
	"log"
	"fmt"
	"github.com/nvhai245/go-websocket-chat/pkg/syncer"
	pb2 "github.com/nvhai245/go-chat-synchronizer/proto"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			var onlineUsers []string
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for cl, _ := range pool.Clients {
				onlineUsers = append(onlineUsers, cl.Username)
			}
			for cl, _ := range pool.Clients {
				fmt.Println(cl)
				cl.Conn.WriteJSON(Message{Type: "system", Body: client.Username + " has joined the chat...", Username: "admin"})
				cl.Conn.WriteJSON(Message{Type: "online", Body2: onlineUsers, Username: "admin"})
			}
			break
		case client := <-pool.Unregister:
			var onlineUsers []string
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for cl, _ := range pool.Clients {
				onlineUsers = append(onlineUsers, cl.Username)
			}
			for cl, _ := range pool.Clients {
				cl.Conn.WriteJSON(Message{Type: "offline", Body2: onlineUsers, Username: "admin"})
			}
			break
		case message := <-pool.Broadcast:
			if message.Type == "chat" {
				fmt.Println("Sending message to all clients in Pool")
				for client, _ := range pool.Clients {
					if err := client.Conn.WriteJSON(message); err != nil {
						fmt.Println(err)
						return
					}
				}
				writeData := []*pb2.WriteRequest{{Count: message.Count, Author: message.Username, Message: message.Body}}
				success := syncer.Write(writeData, syncer.GrpcClient2)
				if success == true {
					log.Println("written to db: ", writeData)
				}
			}
			if message.Type == "login" || message.Type == "register" {
				fmt.Println("Responding to newly login user")
				for client, _ := range pool.Clients {
					if client.ID == message.ID {
						if err := client.Conn.WriteJSON(message); err != nil {
							fmt.Println(err)
							return
						}
					}
				}
			}
			if message.Type == "logout" {
				fmt.Println("An user left")
				for client, _ := range pool.Clients {
					if client.ID == message.ID {
						if err := client.Conn.WriteJSON(message); err != nil {
							fmt.Println(err)
							return
						}
					}
				}
				for client, _ := range pool.Clients {
					fmt.Println(client)
					client.Conn.WriteJSON(Message{Type: "system", Body: message.Username + " has left the chat...", Username: "admin"})
				}
			}
		}
	}
}
