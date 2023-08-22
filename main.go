package main

import (
	"fmt"
	"net/http"

	"github.com/masfuulaji/go-chat/pkg/websocket"
	"github.com/rs/cors"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("websocket enpoint reached")

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}
	pool.Register <- client

	client.Read()
}

func setupRoute(mux *http.ServeMux) {
	pool := websocket.NewPool()
	go pool.Start()

    mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {
	mux := http.NewServeMux()
    mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	setupRoute(mux)
    handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080", handler)
}
