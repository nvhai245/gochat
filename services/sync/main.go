package main

import (
	"context"
	_ "database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	pb "github.com/nvhai245/gochat/services/sync/proto"
	"google.golang.org/grpc"
)

const (
	port = ":9091"
)

const connStr = "postgres://rneveabk:PCH8f9ePpFOjPgXTr-9yNV6PL-CqLCoQ@satao.db.elephantsql.com:5432/rneveabk"

// const connStr = "postgres://postgres:Harin245@localhost:5432/gochat?sslmode=disable"

var db, dbErr = sqlx.Open("postgres", connStr)

type server struct {
	pb.UnimplementedSyncServer
}

func (s *server) Write(stream pb.Sync_WriteServer) error {
	var allMessages = []*pb.WriteRequest{}
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		allMessages = append(allMessages, message)
	}
	for _, message := range allMessages {
		quoted := pq.QuoteIdentifier(message.Table)
		sqlStatement := fmt.Sprintf(`
	INSERT INTO %s (count, author, message)
	VALUES ($1, $2, $3)`, quoted)
		rows, err := db.Query(sqlStatement, message.Count, message.Author, message.Message)
		if err != nil {
			log.Println(err)
		}
		defer rows.Close();
	}
	return stream.SendAndClose(&pb.WriteResponse{
		Success: true,
	})
}
func (s *server) GetDifference(ctx context.Context, localCount *pb.GetDifferenceRequest) (*pb.GetDifferenceResponse, error) {
	quoted := pq.QuoteIdentifier(localCount.Table)
	sqlStatement := fmt.Sprintf(`SELECT * FROM %s ORDER BY count DESC LIMIT 1`, quoted)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	var (
		count   int64
		author  string
		message string
	)

	for rows.Next() {
		if err := rows.Scan(&count, &author, &message); err != nil {
			log.Fatal(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	difference := count
	return &pb.GetDifferenceResponse{Difference: difference}, nil
}
func (s *server) Read(messagesRange *pb.ReadRequest, stream pb.Sync_ReadServer) error {
	for i := messagesRange.First; i <= messagesRange.Last; i++ {
		quoted := pq.QuoteIdentifier(messagesRange.Table)
		sqlStatement := fmt.Sprintf(`SELECT count, author, message FROM %s WHERE count=$1`, quoted)
		rows, err := db.Query(sqlStatement, i)
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()
		message := &pb.ReadResponse{}
		for rows.Next() {
			if err := rows.Scan(&message.Count, &message.Author, &message.Message); err != nil {
				log.Fatal(err)
			}
		}
		if err := stream.Send(message); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) CheckExist(ctx context.Context, tables *pb.CheckRequest) (*pb.CheckResponse, error) {
	quoted1 := pq.QuoteIdentifier(tables.Table1)
	quoted2 := pq.QuoteIdentifier(tables.Table2)
	var (
		exists1 bool
		exists2 bool
	)
	sqlStatement1 := fmt.Sprintf(`SELECT 1 FROM %s LIMIT 1`, quoted1)
	sqlStatement2 := fmt.Sprintf(`SELECT 1 FROM %s LIMIT 1`, quoted2)
	rows1, err1 := db.Query(sqlStatement1)
	if err1 != nil {
		log.Println(tables.Table1, " does not exist")
	} else {
		exists1 = true
		log.Println(err1)
		defer rows1.Close();
	}
	rows2, err2 := db.Query(sqlStatement2)
	if err2 != nil {
		log.Println(tables.Table2, " does not exist")
	} else {
		exists2 = true
		log.Println(err2)
		defer rows2.Close();
	}
	log.Println(exists1)
	log.Println(exists2)
	if exists1 {
		return &pb.CheckResponse{Table: tables.Table1}, nil
	} else if exists2 {
		return &pb.CheckResponse{Table: tables.Table2}, nil
	} else {
		quoted := pq.QuoteIdentifier(tables.Table1)
		sqlStatement := fmt.Sprintf(`CREATE TABLE %s (count INT NOT NULL UNIQUE, author TEXT NOT NULL, message TEXT NOT NULL)`, quoted)
		_, err := db.Exec(sqlStatement)
		if err != nil {
			log.Println(err)
		}
		sqlStatement = fmt.Sprintf(`
	INSERT INTO %s (count, author, message)
	VALUES ($1, $2, $3)`, quoted)
		rows4, err := db.Query(sqlStatement, 0, "admin", "Say hi")
		if err != nil {
			log.Println(err)
		}
		defer rows4.Close();
		return &pb.CheckResponse{Table: tables.Table1}, nil
	}
}

func main() {
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	if dbErr != nil {
		log.Println(dbErr)
	}
	// defer db.Close()

	dbErr = db.Ping()
	if dbErr != nil {
		log.Println(dbErr)
	}
	log.Println("DB connected!")

	flag.Parse()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSyncServer(grpcServer, &server{})
	// determine whether to use TLS

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
