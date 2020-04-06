package notification

import (
	"context"
	"log"
	"io"

	pb3 "github.com/nvhai245/gochat/services/notification/proto"
	"google.golang.org/grpc"
)

var GrpcConn3, GrpcErr3 = grpc.Dial(":9092", grpc.WithInsecure())
var GrpcClient3 = pb3.NewNotificationClient(GrpcConn3)

type Message struct {
	ID string
	Count int64 `json:"count"`
	Type     string `json:"type"`
	Body     string `json:"body"`
	Username string `json:"username"`
	Body2 []string `json:"body2"`
	Body3 string `json:"body3"`
	Table string `json:"table"`
	Receiver []string `json:"receiver"`
	Deleted bool `json:"deleted"`
}

func Create(username string, client pb3.NotificationClient) (success bool) {
	s, err := client.Create(context.Background(), &pb3.CreateRequest{Username: username})
	if err != nil {
		log.Println(err)
		return false
	}
	if s.Success == false {
		return false
	}
	return true
}

func Add(username string, table string, client pb3.NotificationClient) (success bool) {
	s, err := client.Add(context.Background(), &pb3.AddRequest{Username: username, Table: table})
	if err != nil {
		log.Println(err)
		return false
	}
	if s.Success == false {
		return false
	}
	return true
}

func Remove(username string, table string, client pb3.NotificationClient) (success bool) {
	s, err := client.Remove(context.Background(), &pb3.RemoveRequest{Username: username, Table: table})
	if err != nil {
		log.Println(err)
		return false
	}
	if s.Success == false {
		return false
	}
	return true
}

func Get(username string, client pb3.NotificationClient) (messages []Message) {
	stream, err := client.Get(context.Background(), &pb3.GetRequest{Username: username})
	if err != nil {
		log.Println(err)
		return []Message{}
	}
	msgs := []Message{}
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.Read(_) = _, %v", client, err)
		}
		log.Println(message)
		msgs = append(msgs, Message{ID: "", Type: "getnoti", Count: message.Count, Body: message.Message, Username: message.Author, Table: message.Table, Deleted: message.Deleted})
	}
	return msgs
}