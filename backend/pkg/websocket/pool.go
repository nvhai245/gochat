package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	pb2 "github.com/nvhai245/go-chat-synchronizer/proto"
	"github.com/nvhai245/go-websocket-chat/pkg/syncer"
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

func existIn(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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
			difference := syncer.GetDifference(0, "all", syncer.GrpcClient2)
			for cl, _ := range pool.Clients {
				if cl.Username == client.Username {
					if err := cl.Conn.WriteJSON(Message{Type: "update", Count: difference, Body: "", Username: "admin", Table: "all"}); err != nil {
						fmt.Println(err)
						return
					}
				}
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
				if message.Receiver[0] == "all" {
					for client, _ := range pool.Clients {
						if err := client.Conn.WriteJSON(message); err != nil {
							fmt.Println(err)
							return
						}
					}
				} else {
					for client, _ := range pool.Clients {
						if existIn(client.Username, message.Receiver) {
							if err := client.Conn.WriteJSON(message); err != nil {
								fmt.Println(err)
								return
							}
						}
					}
				}
				writeData := []*pb2.WriteRequest{{Count: message.Count, Author: message.Username, Message: message.Body, Table: message.Table}}
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
			if message.Type == "readdb" {
				first, err := strconv.ParseInt(message.Body, 10, 64)
				if err != nil {
					log.Println(err)
				}
				last, err2 := strconv.ParseInt(message.Body3, 10, 64)
				if err2 != nil {
					log.Println(err2)
				}
				messages := syncer.Read(first, last, message.Table, syncer.GrpcClient2)
				for client, _ := range pool.Clients {
					if client.ID == message.ID {
						for _, msg := range messages {
							if err := client.Conn.WriteJSON(msg); err != nil {
								fmt.Println(err)
								return
							}
						}
					}
				}
			}
			if message.Type == "writedb" {
				writeData := []*pb2.WriteRequest{}
				for _, msg := range message.Body2 {
					parsedMsg := Message{}
					err := json.Unmarshal([]byte(msg), &parsedMsg)
					if err != nil {
						log.Println(err)
					}
					writeData = append(writeData, &pb2.WriteRequest{Count: parsedMsg.Count, Author: parsedMsg.Username, Message: parsedMsg.Body, Table: parsedMsg.Table})
				}
				success := syncer.Write(writeData, syncer.GrpcClient2)
				if success == true {
					log.Println("written to db: ", writeData)
				}
			}
			if message.Type == "checkExist" {
				table := syncer.CheckExist(message.Body, message.Body3, syncer.GrpcClient2)
				for client, _ := range pool.Clients {
					if client.ID == message.ID {
						if err := client.Conn.WriteJSON(Message{Type: "checkExist", Table: table, Username: "admin"}); err != nil {
							fmt.Println(err)
							return
						}
					}
				}
			}
			if message.Type == "getDifference" {
				difference := syncer.GetDifference(0, message.Table, syncer.GrpcClient2)
				for client, _ := range pool.Clients {
					if client.ID == message.ID {
						if err := client.Conn.WriteJSON(Message{Count: difference, Type: "update", Table: message.Table, Username: "admin"}); err != nil {
							fmt.Println(err)
							return
						}
					}
				}
			}
		}
	}
}
