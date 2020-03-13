package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	_ "sync"

	"github.com/gorilla/websocket"
	pb "github.com/nvhai245/go-chat-authservice/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	ID string
	Type     string `json:"type"`
	Body     string `json:"body"`
	Username string `json:"username"`
}

func Signup(message Message) (success bool) {
	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	client := pb.NewAuthClient(conn)
	var header metadata.MD
	token, err := client.Register(context.Background(), &pb.RegisterRequest{Username: message.Username, Password: message.Body}, grpc.Header(&header))
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("Recieved from Auth service: ", token)
	log.Println("JWT header: ", header["authorization"])
	return true
}

func Login(message Message) (isAuthenticated bool) {
	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	client := pb.NewAuthClient(conn)
	var header metadata.MD
	result, err := client.Login(context.Background(), &pb.LoginRequest{Username: message.Username, Password: message.Body}, grpc.Header(&header))
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("Recieved from Auth service: ", result)
	if result.Authorized == true {
		return true
	}
	return false
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
		errorMessage := Message{ID: c.ID, Type: "authfail", Body: "", Username: ""}
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		json.Unmarshal([]byte(string(p)), &message)
		log.Println("message type is ", message.Type)
		message.ID = c.ID
		if message.Type == "register" {
			success := Signup(message)
			if success == false {
				log.Println("Authentication failed")
				c.Pool.Broadcast <- errorMessage
			} else {
				c.Pool.Broadcast <- message
			}
		} else if message.Type == "login" {
			success := Login(message)
			if success == false {
				log.Println("Login failed")
				c.Pool.Broadcast <- errorMessage
			} else {
				c.Pool.Broadcast <- message
			}
		} else if message.Type == "chat" {
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
