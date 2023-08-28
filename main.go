package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/masfuulaji/go-chat/pkg/websocket"
	"github.com/rs/cors"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func generateRandomID(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    seed := rand.NewSource(time.Now().UnixNano())
    random := rand.New(seed)

    id := make([]byte, length)
    for i := range id {
        id[i] = charset[random.Intn(len(charset))]
    }

    return string(id)
}

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("websocket enpoint reached")

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

    room := r.URL.Query().Get("room")
    if room == "" {
        conn.Close()
        return
    }

    clientID := generateRandomID(6) 

	client := &websocket.Client{
        ID:   clientID,
		Conn: conn,
		Pool: pool,
        Room: room,
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
