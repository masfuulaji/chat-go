package main

import (
	"fmt"
	"net/http"

	"github.com/masfuulaji/go-chat/pkg/websocket"
)

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

func setupRoute() {
    pool := websocket.NewPool() 
    go pool.Start()

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
       serveWs(pool, w, r) 
    })
}
func main() {
    fmt.Println("Hello, World!")
    setupRoute()
    http.ListenAndServe("9000", nil)
}
