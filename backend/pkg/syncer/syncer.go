package syncer

import (
	"context"
	"io"
	"log"

	pb2 "github.com/nvhai245/go-chat-synchronizer/proto"
	"google.golang.org/grpc"
)

var GrpcConn2, GrpcErr2 = grpc.Dial(":9091", grpc.WithInsecure())
var GrpcClient2 = pb2.NewSyncClient(GrpcConn2)

type ChatMessage struct {
	Count   int64  `json:"count"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

type Message struct {
	ID string
	Count int64 `json:"count"`
	Type     string `json:"type"`
	Body     string `json:"body"`
	Username string `json:"username"`
	Body2 []string `json:"body2"`
	Body3 string `json:"body3"`
}

// Write message to db
func Write(messages []*pb2.WriteRequest, client pb2.SyncClient) (success bool) {
	stream, err := client.Write(context.Background())
	if err != nil {
		log.Fatalf("%v.Write(_) = _, %v", client, err)
	}
	for _, message := range messages {
		if err := stream.Send(message); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("%v.Send(%v) = %v", stream, message, err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Write result: %v", reply)
	if reply.Success == true {
		return true
	}
	return false
}

func GetDifference(localCount int64, client pb2.SyncClient) (difference int64) {
	d, err := client.GetDifference(context.Background(), &pb2.GetDifferenceRequest{Local: localCount})
	if err != nil {
		log.Println(err)
		return 0
	}
	log.Println("Get count difference successfull: ", d.Difference)
	return d.Difference
}

func Read(localCount int64, dbCount int64, client pb2.SyncClient) (messages []Message) {
	readRequest := &pb2.ReadRequest{First: localCount + 1, Last: dbCount}
	stream, err := client.Read(context.Background(), readRequest)
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
		msgs = append(msgs, Message{ID: "", Type: "readdb", Count: message.Count, Body: message.Message, Username: message.Author})
	}
	return msgs
}
