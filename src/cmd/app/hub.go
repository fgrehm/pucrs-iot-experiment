package main

import (
	"log"
)

type Hub struct {
	// Registered connections.
	Connections map[*Connection]bool

	// Inbound messages from the connections.
	Broadcast chan string

	// Register requests from the connections.
	Register chan *Connection

	// Unregister requests from connections.
	Unregister chan *Connection
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Connections[c] = true
			log.Println("Number of clients connected:", len(h.Connections))
		case c := <-h.Unregister:
			log.Println("Client disconnected")
			if _, ok := h.Connections[c]; ok {
				delete(h.Connections, c)
				close(c.Send)
			}
		case m := <-h.Broadcast:
			log.Println("Sending messages")
			for c := range h.Connections {
				select {
				case c.Send <- m:
				default:
					delete(h.Connections, c)
					close(c.Send)
				}
			}
		}
	}
}
