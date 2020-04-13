package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/nvhai245/gochat/server/pkg/auth"
	"github.com/nvhai245/gochat/server/pkg/notification"
	"github.com/nvhai245/gochat/server/pkg/syncer"
	pb2 "github.com/nvhai245/gochat/services/sync/proto"
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
			success := notification.Create(client.Username, notification.GrpcClient3)
			if success == true {
				log.Println("notification table is ready for user ", client.Username)
			}
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
				writeData := []*pb2.WriteRequest{{Count: message.Count, Author: message.Username, Message: message.Body, Table: message.Table, Deleted: message.Deleted}}
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
					writeData = append(writeData, &pb2.WriteRequest{Count: parsedMsg.Count, Author: parsedMsg.Username, Message: parsedMsg.Body, Table: parsedMsg.Table, Deleted: parsedMsg.Deleted})
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
			if message.Type == "delete" {
				success := syncer.Delete(message.Count, message.Table, syncer.GrpcClient2)
				if success == true {
					if message.Table == "all" {
						for client, _ := range pool.Clients {
							if err := client.Conn.WriteJSON(Message{Type: "delete", Count: message.Count, Table: message.Table, Username: "admin", Deleted: true}); err != nil {
								fmt.Println(err)
								return
							}
						}
					} else {
						for client, _ := range pool.Clients {
							if existIn(client.Username, message.Receiver) {
								if err := client.Conn.WriteJSON(Message{Type: "delete", Count: message.Count, Table: message.Table, Username: "admin", Deleted: true}); err != nil {
									fmt.Println(err)
									return
								}
							}
						}
					}
				}
			}
			if message.Type == "restore" {
				success := syncer.Restore(message.Count, message.Table, syncer.GrpcClient2)
				if success == true {
					if message.Table == "all" {
						for client, _ := range pool.Clients {
							if err := client.Conn.WriteJSON(Message{Type: "restore", Count: message.Count, Table: message.Table, Username: "admin", Deleted: false}); err != nil {
								fmt.Println(err)
								return
							}
						}
					} else {
						for client, _ := range pool.Clients {
							if existIn(client.Username, message.Receiver) {
								if err := client.Conn.WriteJSON(Message{Type: "restore", Count: message.Count, Table: message.Table, Username: "admin", Deleted: false}); err != nil {
									fmt.Println(err)
									return
								}
							}
						}
					}
				}
			}
			if message.Type == "addnoti" {
				success := notification.Create(message.Receiver[1], notification.GrpcClient3)
				if success == true {
					log.Println("notification table is ready for user ", message.Receiver[1])
				}
				success = notification.Add(message.Receiver[1], message.Table, notification.GrpcClient3)
				if success == true {
					log.Println("added unread message for user ", message.Receiver[1])
				}
			}
			if message.Type == "removenoti" {
				success := notification.Remove(message.Username, message.Table, notification.GrpcClient3)
				if success == true {
					log.Println("removed unread message for user ", message.Username)
				}
			}
			if message.Type == "getnoti" {
				messages := notification.Get(message.Username, notification.GrpcClient3)
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
			if message.Type == "getUser" {
				success, user := auth.GetUser(message.Username, auth.GrpcClient)
				if success {
					userBytes, err := json.Marshal(user)
					if err != nil {
						log.Println("error:", err)
					}
					for client, _ := range pool.Clients {
						if client.ID == message.ID {
							if err := client.Conn.WriteJSON(Message{Type: "getUser", Username: message.Username, Body: string(userBytes)}); err != nil {
								fmt.Println(err)
								return
							}
						}
					}
				}
			}
			if message.Type == "updateEmail" {
				success, email := auth.UpdateEmail(message.Username, message.Body, auth.GrpcClient)
				if success {
					for client, _ := range pool.Clients {
						if client.ID == message.ID {
							if err := client.Conn.WriteJSON(Message{Type: "updateEmail", Username: message.Username, Body: email}); err != nil {
								log.Println(err)
								return
							}
						}
					}
				}
			}
			if message.Type == "updatePhone" {
				success, phone := auth.UpdatePhone(message.Username, message.Body, auth.GrpcClient)
				if success {
					for client, _ := range pool.Clients {
						if client.ID == message.ID {
							if err := client.Conn.WriteJSON(Message{Type: "updatePhone", Username: message.Username, Body: phone}); err != nil {
								log.Println(err)
								return
							}
						}
					}
				}
			}
			if message.Type == "updateBirthday" {
				success, birthday := auth.UpdateBirthday(message.Username, message.Body, auth.GrpcClient)
				if success {
					for client, _ := range pool.Clients {
						if client.ID == message.ID {
							if err := client.Conn.WriteJSON(Message{Type: "updateBirthday", Username: message.Username, Body: birthday}); err != nil {
								log.Println(err)
								return
							}
						}
					}
				}
			}
			if message.Type == "updateFb" {
				success, fb := auth.UpdateFb(message.Username, message.Body, auth.GrpcClient)
				if success {
					for client, _ := range pool.Clients {
						if client.ID == message.ID {
							if err := client.Conn.WriteJSON(Message{Type: "updateFb", Username: message.Username, Body: fb}); err != nil {
								log.Println(err)
								return
							}
						}
					}
				}
			}
			if message.Type == "updateInsta" {
				success, insta := auth.UpdateInsta(message.Username, message.Body, auth.GrpcClient)
				if success {
					for client, _ := range pool.Clients {
						if client.ID == message.ID {
							if err := client.Conn.WriteJSON(Message{Type: "updateInsta", Username: message.Username, Body: insta}); err != nil {
								log.Println(err)
								return
							}
						}
					}
				}
			}
			if message.Type == "updateAvatar" {
				success, avatar := auth.UpdateAvatar(message.Username, message.Body, auth.GrpcClient)
				if success {
					for client, _ := range pool.Clients {
						if client.ID == message.ID {
							if err := client.Conn.WriteJSON(Message{Type: "updateAvatar", Username: message.Username, Body: avatar}); err != nil {
								log.Println(err)
								return
							}
						}
					}
				}
			}
		}
	}
}
