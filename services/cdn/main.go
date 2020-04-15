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
	pb "github.com/nvhai245/gochat/services/cdn/proto"
	"google.golang.org/grpc"
)

const (
	port = ":9093"
)

const connStr = "postgres://rneveabk:PCH8f9ePpFOjPgXTr-9yNV6PL-CqLCoQ@satao.db.elephantsql.com:5432/rneveabk"

// const connStr = "postgres://postgres:Harin245@localhost:5432/gochat?sslmode=disable"

var db, dbErr = sqlx.Open("postgres", connStr)

type server struct {
	pb.UnimplementedCdnServer
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
	pb.RegisterCdnServer(grpcServer, &server{})
	// determine whether to use TLS

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}