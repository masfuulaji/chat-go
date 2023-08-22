package websocket

import "fmt"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			p.Clients[client] = true
			fmt.Println("Client registered")
			client.Conn.WriteJSON(Message{Type: 1, Body: "Welcome!", Sender: client.ID, MessageType: 3})
			for client := range p.Clients {
				fmt.Println("Sending welcome message")
				client.Conn.WriteJSON(Message{Type: 1, Body: "New client connected", MessageType: 2})
			}
			break
		case client := <-p.Unregister:
			delete(p.Clients, client)
			fmt.Println("Client unregistered")
			for client := range p.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "Client disconnected", MessageType: 2})
			}
			break
		case message := <-p.Broadcast:
			fmt.Println("Broadcasting message")
			for client := range p.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
