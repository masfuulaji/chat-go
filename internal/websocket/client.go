package websocket

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
    Room string
	mu   sync.Mutex
}

type Message struct {
    Type int `json:"type"`
    Body string `json:"body"`
    Sender string `json:"sender"`
    MessageType int `json:"message_type"`
    Room string `json:"room"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
        message := Message{Type: messageType, Body: string(p), Sender: c.ID, MessageType: 1}
        c.Pool.Broadcast <- message
        fmt.Printf("Message received: %s\n", message.Body)
	}
}
