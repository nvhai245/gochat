package auth

import (
	"context"
	"log"

	pb "github.com/nvhai245/go-chat-authservice/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type User struct {
	Username string
	Password string
}

func Signup(user User) (success bool, jwt string) {
	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	client := pb.NewAuthClient(conn)
	var header metadata.MD
	response, err := client.Register(context.Background(), &pb.RegisterRequest{Username: user.Username, Password: user.Password}, grpc.Header(&header))
	if err != nil {
		log.Println(err)
		return false, ""
	}
	log.Println("Recieved from Auth service: ", response.Token)
	log.Println("JWT header: ", header["authorization"])
	return true, response.Token
}

func Login(user User) (isAuthenticated bool, jwt string) {
	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	client := pb.NewAuthClient(conn)
	var header metadata.MD
	response, err := client.Login(context.Background(), &pb.LoginRequest{Username: user.Username, Password: user.Password}, grpc.Header(&header))
	if err != nil {
		log.Println(err)
		return false, ""
	}
	log.Println("Recieved from Auth service: ", response.Token)
	if response.Success == true {
		return true, response.Token
	}
	return false, ""
}

func Check(token string) (valid bool) {
	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	client := pb.NewAuthClient(conn)
	var header metadata.MD
	response, err := client.Check(context.Background(), &pb.CheckRequest{Token: token}, grpc.Header(&header))
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("Recieved from Auth service: ", response.Valid)
	if response.Valid == true {
		return true
	}
	return false
}