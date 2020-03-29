package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	pb "github.com/nvhai245/go-chat-authservice/proto"
	"github.com/nvhai245/go-websocket-chat/pkg/auth"
	"github.com/nvhai245/go-websocket-chat/pkg/syncer"
	"github.com/nvhai245/go-websocket-chat/pkg/websocket"
	"google.golang.org/grpc"
)

type UserObject struct {
	Users []string `json:"users"`
}

var grpcConn, grpcErr = grpc.Dial(":9090", grpc.WithInsecure())
var grpcClient = pb.NewAuthClient(grpcConn)

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	cookie, err := r.Cookie("go-chat")
	if err != nil {
		fmt.Printf("Cant find cookie :/\r\n")
		return
	}
	valid, username := auth.Check(cookie.Value, grpcClient)
	log.Println("client username is: ", username)
	if valid == true {
		conn, err := websocket.Upgrade(w, r)
		if err != nil {
			fmt.Fprintf(w, "%+V\n", err)
		}

		client := &websocket.Client{
			ID:       uuid.New().String(),
			Username: username,
			Conn:     conn,
			Pool:     pool,
		}

		pool.Register <- client
		client.Read()
	} else {
		return
	}
}

func main() {
	if grpcErr != nil {
		log.Println(grpcErr)
	}
	if syncer.GrpcErr2 != nil {
		log.Println(syncer.GrpcErr2)
	}
	defer grpcConn.Close()
	defer syncer.GrpcConn2.Close()
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "POST" {
			username, password := r.FormValue("username"), r.FormValue("password")
			newUser := auth.User{Username: username, Password: password}
			success, token := auth.Signup(newUser, grpcClient)
			if success == true {
				http.SetCookie(w, &http.Cookie{
					Name:     "go-chat",
					Value:    token,
					Expires:  time.Now().Add(60 * time.Minute),
					HttpOnly: true,
				})
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "POST" {
			username, password := r.FormValue("username"), r.FormValue("password")
			newUser := auth.User{Username: username, Password: password}
			success, token := auth.Login(newUser, grpcClient)
			if success == true {
				http.SetCookie(w, &http.Cookie{
					Name:     "go-chat",
					Value:    token,
					Expires:  time.Now().Add(60 * time.Minute),
					HttpOnly: true,
				})
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	})
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "POST" {
			http.SetCookie(w, &http.Cookie{
				Name:     "go-chat",
				Value:    "",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
			})
			w.WriteHeader(http.StatusOK)
			return
		}
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "GET" {
			users := auth.GetAllUser(grpcClient)
			allUsers := UserObject{Users: users}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(allUsers)
		}
	})

	fmt.Println("Chat server running")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
