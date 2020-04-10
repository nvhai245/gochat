package auth

import (
	"context"
	"log"
	"io"

	pb "github.com/nvhai245/gochat/services/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type User struct {
	IsAdmin  bool      `json:"isadmin"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Avatar   string    `json:"avatar"`
	Phone    string    `json:"phone"`
	Birthday time.Time `json:"birthday"`
	Fb       string    `json:"fb"`
	Insta    string    `json:"insta"`
}

func Signup(user User, client pb.AuthClient) (success bool, jwt string) {
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

func Login(user User, client pb.AuthClient) (isAuthenticated bool, jwt string) {
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

func Check(token string, client pb.AuthClient) (valid bool, username string) {
	var header metadata.MD
	response, err := client.Check(context.Background(), &pb.CheckRequest{Token: token}, grpc.Header(&header))
	if err != nil {
		log.Println(err)
		return false, ""
	}
	log.Println("Recieved from Auth service: ", response.Valid)
	if response.Valid == true {
		return true, response.Username
	}
	return false, ""
}

func GetAllUser(client pb.AuthClient) (allUser []string) {
	adminRequest := &pb.GetAllUserRequest{IsAdmin: true}
	stream, err := client.GetAllUser(context.Background(), adminRequest)
	if err != nil {
		log.Println(err)
	}
	for {
		user, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetAllUser(_) = _, %v", client, err)
		}
		log.Println(user)
		allUser = append(allUser, user.Username)
	}
	return allUser
}

func GetUser(username string, client pb.AuthClient) (user User) {
	
}