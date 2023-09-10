package services

import (
	"fmt"
	"net/http"

	"github.com/masfuulaji/go-chat/internal/utils"
	"github.com/masfuulaji/go-chat/internal/websocket"
)

func ServeWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
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

    clientID := utils.GenerateRandomID(6) 

	client := &websocket.Client{
        ID:   clientID,
		Conn: conn,
		Pool: pool,
        Room: room,
	}
	pool.Register <- client

	client.Read()
}
