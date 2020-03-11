package websocket

import (
    "fmt"
    "log"
	_ "sync"
    "encoding/json"
    "context"

    "github.com/gorilla/websocket"
    "google.golang.org/grpc"
    pb "github.com/nvhai245/go-chat-authservice/proto"
    "google.golang.org/grpc/metadata"
)

type Client struct {
    ID   string
    Conn *websocket.Conn
    Pool *Pool
}

type Message struct {
    Type string    `json:"type"`
	Body string `json:"body"`
	Username string `json:"username"`
}

func Authenticate() {
	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
    client := pb.NewAuthClient(conn)
    var header metadata.MD
	token, err := client.Login(context.Background(), &pb.LoginRequest{Username: "nvhai245", Password: "12345"}, grpc.Header(&header))
	if err != nil {
		log.Println(err)
	}
    log.Println("Recieved from Auth service: ", token)
    log.Println("JWT header: ", header["authorization"])
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
        if (message.Type == "auth") {
            Authenticate()
        } else {
        c.Pool.Broadcast <- message
        fmt.Printf("Message Received: %+v\n", message)
        }
    }
}