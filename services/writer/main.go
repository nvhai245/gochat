package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"
	_ "database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/dgrijalva/jwt-go"
	pb "github.com/nvhai245/gochat/services/writer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	// PostgreSQL driver
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	port = ":9093"
)

const connStr = "postgres://rneveabk:PCH8f9ePpFOjPgXTr-9yNV6PL-CqLCoQ@satao.db.elephantsql.com:5432/rneveabk"
// const connStr = "postgres://postgres:Harin245@localhost:5432/gochat?sslmode=disable"
var	db, dbErr = sqlx.Open("postgres", connStr)

type server struct {
	pb.UnimplementedAuthServer
}

type JwtCustomClaims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	Email    string `json:"email"`
	jwt.StandardClaims
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
	pb.RegisterAuthServer(grpcServer, &server{})
	// determine whether to use TLS

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
