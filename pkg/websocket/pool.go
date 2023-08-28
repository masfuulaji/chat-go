package websocket

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Pool struct {
	Rooms      map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	mu         sync.Mutex
}

func NewPool() *Pool {
	return &Pool{
		Rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			p.mu.Lock()
			if _, ok := p.Rooms[client.Room]; !ok {
				p.Rooms[client.Room] = make(map[*Client]bool)
			}
			p.Rooms[client.Room][client] = true
			p.mu.Unlock()
			welcomeMessage := Message{Type: 1, Body: "Welcome!", Sender: client.ID, MessageType: 3, Room: client.Room}
			client.Conn.WriteJSON(welcomeMessage)
			for client := range p.Rooms[client.Room] {
				client.Conn.WriteJSON(Message{Type: 1, Body: "New client connected", MessageType: 2})
			}
			break
		case client := <-p.Unregister:
			p.mu.Lock()
			delete(p.Rooms[client.Room], client)
			if len(p.Rooms[client.Room]) == 0 {
				delete(p.Rooms, client.Room)
			}
			p.mu.Unlock()
			for client := range p.Rooms[client.Room] {
				client.Conn.WriteJSON(Message{Type: 1, Body: "Client disconnected", MessageType: 2})
			}
			break
		case message := <-p.Broadcast:
			p.mu.Lock()
			var messageData struct {
				Room    string `json:"room"`
				Message string `json:"message"`
			}
			err := json.Unmarshal([]byte(message.Body), &messageData)
			if err != nil {
				fmt.Println(err)
			}
			p.mu.Unlock()

			message.Body = messageData.Message
			for client := range p.Rooms[messageData.Room] {
				client.Conn.WriteJSON(message)
			}
		}
	}
}
