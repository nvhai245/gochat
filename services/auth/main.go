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
	"github.com/nvhai245/gochat/services/auth/model"
	pb "github.com/nvhai245/gochat/services/auth/proto"
	userController "github.com/nvhai245/gochat/services/auth/controller/user"
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
	rows := []model.AuthorizedUser{}
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
func (s *server) GetUser(ctx context.Context, user *pb.GetUserRequest) (*pb.AuthorizedUser, error) {
	success, authorizedUser := userController.GetUser(user.Username, db)
	if success == false {
		return &pb.AuthorizedUser{}, nil
	}
	log.Println("getUser successful", authorizedUser)
	return &pb.AuthorizedUser{
		Username: authorizedUser.Username,
		IsAdmin: authorizedUser.IsAdmin,
		Email: authorizedUser.Email,
		Avatar: authorizedUser.Avatar,
		Phone: authorizedUser.Phone,
		Birthday: authorizedUser.Birthday.String(),
		Fb: authorizedUser.Fb,
		Insta: authorizedUser.Insta,
		}, nil
}
func (s *server) UpdateEmail(ctx context.Context, user *pb.UpdateEmailRequest) (*pb.UpdateEmailResponse, error) {
	success, updated := userController.UpdateEmail(user.Username, user.Email, db)
	if success == false {
		return &pb.UpdateEmailResponse{}, nil
	}
	return &pb.UpdateEmailResponse{Email: updated}, nil
}
func (s *server) UpdateAvatar(ctx context.Context, user *pb.UpdateAvatarRequest) (*pb.UpdateAvatarResponse, error) {
	success, updated := userController.UpdateAvatar(user.Username, user.Avatar, db)
	if success == false {
		return &pb.UpdateAvatarResponse{}, nil
	}
	return &pb.UpdateAvatarResponse{Avatar: updated}, nil
}
func (s *server) UpdatePhone(ctx context.Context, user *pb.UpdatePhoneRequest) (*pb.UpdatePhoneResponse, error) {
	success, updated := userController.UpdatePhone(user.Username, user.Phone, db)
	if success == false {
		return &pb.UpdatePhoneResponse{}, nil
	}
	return &pb.UpdatePhoneResponse{Phone: updated}, nil
}
func (s *server) UpdateBirthday(ctx context.Context, user *pb.UpdateBirthdayRequest) (*pb.UpdateBirthdayResponse, error) {
	t, err := time.Parse("20060102T150405", user.Birthday)
	if err != nil {
		log.Println(err)
	}
	success, updated := userController.UpdateBirthday(user.Username, t, db)
	if success == false {
		return &pb.UpdateBirthdayResponse{}, nil
	}
	return &pb.UpdateBirthdayResponse{Birthday: updated.String()}, nil
}
func (s *server) UpdateFb(ctx context.Context, user *pb.UpdateFbRequest) (*pb.UpdateFbResponse, error) {
	success, updated := userController.UpdateFb(user.Username, user.Fb, db)
	if success == false {
		return &pb.UpdateFbResponse{}, nil
	}
	return &pb.UpdateFbResponse{Fb: updated}, nil
}
func (s *server) UpdateInsta(ctx context.Context, user *pb.UpdateInstaRequest) (*pb.UpdateInstaResponse, error) {
	success, updated := userController.UpdateInsta(user.Username, user.Insta, db)
	if success == false {
		return &pb.UpdateInstaResponse{}, nil
	}
	return &pb.UpdateInstaResponse{Insta: updated}, nil
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
