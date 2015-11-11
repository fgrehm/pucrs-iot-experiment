package main

import (
	"log"

	"golang.org/x/net/websocket"
)

type Connection struct {
	// The websocket connection.
	WS *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan string

	// The hub.
	Hub *Hub
}

func (c *Connection) Writer() {
	for message := range c.Send {
		if err := websocket.Message.Send(c.WS, message); err != nil {
			log.Println("Could not send message to a client", err)
			break
		}
	}
	c.Hub.Unregister <- c
	c.WS.Close()
}
