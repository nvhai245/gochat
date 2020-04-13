package auth

import (
	"context"
	"io"
	"log"

	pb "github.com/nvhai245/gochat/services/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type User struct {
	IsAdmin  bool      `json:"isadmin"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
	Avatar   string    `json:"avatar"`
	Phone    string    `json:"phone"`
	Birthday string `json:"birthday"`
	Fb       string    `json:"fb"`
	Insta    string    `json:"insta"`
}

var GrpcConn, GrpcErr = grpc.Dial(":9090", grpc.WithInsecure())
var GrpcClient = pb.NewAuthClient(GrpcConn)

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

func GetUser(username string, client pb.AuthClient) (success bool, user User) {
	authorizedUser, err := client.GetUser(context.Background(), &pb.GetUserRequest{Username: username})
	if err != nil {
		log.Println(err)
		return false, User{}
	}
	log.Println("getUser successful", authorizedUser)
	return true, User{
		IsAdmin:  authorizedUser.IsAdmin,
		Username: authorizedUser.Username,
		Email:    authorizedUser.Email,
		Avatar:   authorizedUser.Avatar,
		Phone:    authorizedUser.Phone,
		Birthday: authorizedUser.Birthday,
		Fb:       authorizedUser.Fb,
		Insta:    authorizedUser.Insta,
	}
}

func UpdateEmail(username string, email string, client pb.AuthClient) (success bool, updatedEmail string) {
	updated, err := client.UpdateEmail(context.Background(), &pb.UpdateEmailRequest{Username: username, Email: email})
	if err != nil {
		log.Println(err)
		return false, email
	}
	return true, updated.Email
}
func UpdatePhone(username string, phone string, client pb.AuthClient) (success bool, updatedPhone string) {
	updated, err := client.UpdatePhone(context.Background(), &pb.UpdatePhoneRequest{Username: username, Phone: phone})
	if err != nil {
		log.Println(err)
		return false, phone
	}
	return true, updated.Phone
}
func UpdateAvatar(username string, avatar string, client pb.AuthClient) (success bool, updatedAvatar string) {
	updated, err := client.UpdateAvatar(context.Background(), &pb.UpdateAvatarRequest{Username: username, Avatar: avatar})
	if err != nil {
		log.Println(err)
		return false, avatar
	}
	return true, updated.Avatar
}
func UpdateBirthday(username string, birthday string, client pb.AuthClient) (success bool, updatedBirthday string) {
	updated, err := client.UpdateBirthday(context.Background(), &pb.UpdateBirthdayRequest{Username: username, Birthday: birthday})
	if err != nil {
		log.Println(err)
		return false, birthday
	}
	return true, updated.Birthday
}
func UpdateFb(username string, fb string, client pb.AuthClient) (success bool, updatedFb string) {
	updated, err := client.UpdateFb(context.Background(), &pb.UpdateFbRequest{Username: username, Fb: fb})
	if err != nil {
		log.Println(err)
		return false, fb
	}
	return true, updated.Fb
}
func UpdateInsta(username string, insta string, client pb.AuthClient) (success bool, updatedInsta string) {
	updated, err := client.UpdateInsta(context.Background(), &pb.UpdateInstaRequest{Username: username, Insta: insta})
	if err != nil {
		log.Println(err)
		return false, insta
	}
	return true, updated.Insta
}
