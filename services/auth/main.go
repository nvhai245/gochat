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
	_ "github.com/nvhai245/gochat/services/auth/model"
	pb "github.com/nvhai245/gochat/services/auth//proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	// PostgreSQL driver
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	port = ":9090"
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

func (s *server) Register(ctx context.Context, registerData *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	hash, error := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if error != nil {
		log.Println(error)
	}
	isAdmin := false
	if registerData.Username == "nvhai245" {
		isAdmin = true
	}
	sqlStatement := `
	INSERT INTO users (username, isAdmin, hash)
	VALUES ($1, $2, $3)`
	rows, err := db.Query(sqlStatement, registerData.Username, isAdmin, string(hash))
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	claims := &JwtCustomClaims{
		registerData.Username,
		false,
		"abcd@gmail.com",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	if registerData.Username == "xmancat" {
		claims.IsAdmin = true
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("go-chat"))
	if err != nil {
		log.Println(err)
	}

	header := metadata.New(map[string]string{"authorization": "Bearer " + t})
	grpc.SendHeader(ctx, header)
	return &pb.RegisterResponse{Success: true, Token: t}, nil
}

func (s *server) Login(ctx context.Context, loginData *pb.LoginRequest) (*pb.LoginResponse, error) {
	sqlStatement := "SELECT isAdmin, hash FROM users WHERE username=$1"
	rows, err := db.Query(sqlStatement, loginData.Username)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var (
		hash    []byte
		isAdmin bool
	)

	for rows.Next() {
		if err := rows.Scan(&isAdmin, &hash); err != nil {
			log.Fatal(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	error := bcrypt.CompareHashAndPassword(hash, []byte(loginData.Password))
	if error != nil {
		log.Println(error)
		return &pb.LoginResponse{Success: false, Token: ""}, nil
	}

	claims := &JwtCustomClaims{
		loginData.Username,
		false,
		"abcd@gmail.com",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	if loginData.Username == "nvhai245" {
		claims.IsAdmin = true
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("go-chat"))
	if err != nil {
		log.Println(err)
	}

	header := metadata.New(map[string]string{"authorization": "Bearer " + t})
	grpc.SendHeader(ctx, header)
	return &pb.LoginResponse{Success: true, Token: t}, nil
}

func (s *server) Check(ctx context.Context, checkData *pb.CheckRequest) (*pb.CheckResponse, error) {
	tokenString := checkData.Token
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("go-chat"), nil
	})

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		log.Printf("%v %v", claims.Username + " authorized, jwt is valid. Expires in: ", claims.StandardClaims.ExpiresAt)
		return &pb.CheckResponse{Valid: true, Username: claims.Username}, nil
	}
	log.Println(err)
	return &pb.CheckResponse{Valid: false}, nil
}

func (s *server) GetAllUser(admin *pb.GetAllUserRequest, stream pb.Auth_GetAllUserServer) error {
	type User struct {
		Username string `db:"username"`
		Hash []byte `db:"hash"`
		IsAdmin bool `db:"isadmin"`
	}
	rows := []User{}
	sqlStatement := "SELECT * FROM users"
	err := db.Select(&rows, sqlStatement)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range rows {
		if err := stream.Send(&pb.User{Username: user.Username, PasswordHash: "", IsAdmin: user.IsAdmin}); err != nil {
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
	pb.RegisterAuthServer(grpcServer, &server{})
	// determine whether to use TLS

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
