package syncer

import (
	"context"
	"io"
	"log"

	pb2 "github.com/nvhai245/gochat/services/sync/proto"
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
	Table string `json:"table"`
	Receiver []string `json:"receiver"`
	Deleted bool `json:"deleted"`
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

func GetDifference(localCount int64, table string, client pb2.SyncClient) (difference int64) {
	d, err := client.GetDifference(context.Background(), &pb2.GetDifferenceRequest{Local: localCount, Table: table})
	if err != nil {
		log.Println(err)
		return 0
	}
	log.Println("Get count difference successfull: ", d.Difference)
	return d.Difference
}

func Read(localCount int64, dbCount int64, table string, client pb2.SyncClient) (messages []Message) {
	readRequest := &pb2.ReadRequest{First: localCount + 1, Last: dbCount, Table: table}
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
		msgs = append(msgs, Message{ID: "", Type: "readdb", Count: message.Count, Body: message.Message, Username: message.Author, Table: table, Deleted: message.Deleted})
	}
	return msgs
}

func CheckExist(table1 string, table2 string, client pb2.SyncClient) (table string) {
	checkRequest := &pb2.CheckRequest{Table1: table1, Table2: table2}
	rightTable, err := client.CheckExist(context.Background(), checkRequest)
	if err != nil {
		log.Println(err)
		return "all"
	}
	return rightTable.Table
}

func Delete(count int64, table string, client pb2.SyncClient) (success bool) {
	deleteRequest := &pb2.DeleteRequest{Count: count, Table: table}
	s, err := client.Delete(context.Background(), deleteRequest)
	if err != nil {
		log.Println(err)
		return false
	}
	if s.Success == false {
		return false
	}
	return true
}

func Restore(count int64, table string, client pb2.SyncClient) (success bool) {
	restoreRequest := &pb2.RestoreRequest{Count: count, Table: table}
	s, err := client.Restore(context.Background(), restoreRequest)
	if err != nil {
		log.Println(err)
		return false
	}
	if s.Success == false {
		return false
	}
	return true
}
