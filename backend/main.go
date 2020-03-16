package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nvhai245/go-websocket-chat/pkg/auth"
	"github.com/nvhai245/go-websocket-chat/pkg/websocket"
)

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
	valid := auth.Check(cookie.Value)
	if valid == true {
		conn, err := websocket.Upgrade(w, r)
		if err != nil {
			fmt.Fprintf(w, "%+V\n", err)
		}

		client := &websocket.Client{
			ID:   uuid.New().String(),
			Conn: conn,
			Pool: pool,
		}

		pool.Register <- client
		client.Read()
	} else {
		return
	}
}

func main() {
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "POST" {
			username, password := r.FormValue("username"), r.FormValue("password")
			newUser := auth.User{Username: username, Password: password}
			success, token := auth.Login(newUser)
			if success == true {
				http.SetCookie(w, &http.Cookie{
					Name:    "go-chat",
					Value:   token,
					Expires: time.Now().Add(5 * time.Minute),
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
			success, token := auth.Login(newUser)
			if success == true {
				http.SetCookie(w, &http.Cookie{
					Name:    "go-chat",
					Value:   token,
					Expires: time.Now().Add(5 * time.Minute),
					HttpOnly: true,
				})
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	})
	fmt.Println("Chat server running")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
