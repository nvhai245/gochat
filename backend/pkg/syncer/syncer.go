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
