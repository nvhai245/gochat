package main

import (
	"context"
	_ "database/sql"
	"flag"
	"fmt"
	_ "io"
	"log"
	"net"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	pb "github.com/nvhai245/gochat/services/notification/proto"
	"google.golang.org/grpc"
)

const (
	port = ":9092"
)

const connStr = "postgres://rneveabk:PCH8f9ePpFOjPgXTr-9yNV6PL-CqLCoQ@satao.db.elephantsql.com:5432/rneveabk"

// const connStr = "postgres://postgres:Harin245@localhost:5432/gochat?sslmode=disable"

var db, dbErr = sqlx.Open("postgres", connStr)

type server struct {
	pb.UnimplementedNotificationServer
}

type Tables struct {
	Table string `db:"table"`
}

func (s *server) Create(ctx context.Context, user *pb.CreateRequest) (*pb.CreateResponse, error) {
	quoted := pq.QuoteIdentifier(user.Username + "notification")
	sqlStatement := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s ("table" TEXT NOT NULL UNIQUE)`, quoted)
	_, err := db.Exec(sqlStatement)
	if err != nil {
		log.Println(err)
		return &pb.CreateResponse{Success: false}, nil
	}
	return &pb.CreateResponse{Success: true}, nil
}
func (s *server) Add(ctx context.Context, message *pb.AddRequest) (*pb.AddResponse, error) {
	quoted := pq.QuoteIdentifier(message.Username + "notification")
	sqlStatement := fmt.Sprintf(`
	INSERT INTO %s ("table")
	VALUES ($1) ON CONFLICT ("table") DO NOTHING`, quoted)
	rows, err := db.Query(sqlStatement, message.Table)
	if err != nil {
		log.Println(err)
		return &pb.AddResponse{Success: false}, nil
	}
	defer rows.Close()
	return &pb.AddResponse{Success: true}, nil
}
func (s *server) Remove(ctx context.Context, message *pb.RemoveRequest) (*pb.RemoveResponse, error) {
	quoted := pq.QuoteIdentifier(message.Username + "notification")
	sqlStatement := fmt.Sprintf(`
	DELETE FROM %s 
	WHERE "table" = $1`, quoted)
	rows, err := db.Query(sqlStatement, message.Table)
	if err != nil {
		log.Println(err)
		return &pb.RemoveResponse{Success: false}, nil
	}
	defer rows.Close()
	return &pb.RemoveResponse{}, nil
}

func (s *server) Get(user *pb.GetRequest, stream pb.Notification_GetServer) error {
	var t = []Tables{}
	quoted := pq.QuoteIdentifier(user.Username + "notification")
	sqlStatement := fmt.Sprintf(`SELECT * FROM %s`, quoted)
	db.Select(&t, sqlStatement)
	for _, tb := range t {
		quoted2 := pq.QuoteIdentifier(tb.Table)
		sqlStatement2 := fmt.Sprintf(`SELECT count, author, message, deleted FROM %s ORDER BY count DESC LIMIT 1`, quoted2)
		rows, err := db.Query(sqlStatement2)
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()
		message := &pb.GetResponse{}
		for rows.Next() {
			if err := rows.Scan(&message.Count, &message.Author, &message.Message, &message.Deleted); err != nil {
				log.Println(err)
			}
			message.Table = tb.Table
		}
		if err := stream.Send(message); err != nil {
			return err
		}
	}
	return nil
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
	pb.RegisterNotificationServer(grpcServer, &server{})
	// determine whether to use TLS

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
