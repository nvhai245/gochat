package main

import (
	"fmt"
	"net/http"
	"github.com/nvhai245/go-websocket-chat/pkg/websocket"
)

func setupRoutes() {
	pool := websocket.NewPool()
    go pool.Start()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(pool, w, r)
    })
}

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r*http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := &websocket.Client{
        Conn: conn,
        Pool: pool,
	}
	
	pool.Register <- client
    client.Read()
}



func main() {
	fmt.Println("Chat server running")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}