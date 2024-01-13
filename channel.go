package main

import (
	"errors"
	"ws-chat-server/messages"
)

// Channel maintains the set of active clients and broadcasts messages to the
// clients.
type Channel struct {
	// Channel name
	name string

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan messages.ServerMessageInterface

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewChannel(name string) *Channel {
	return &Channel{
		name:       name,
		broadcast:  make(chan messages.ServerMessageInterface),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (c *Channel) Close() {
	ChatServer.removeChannel(c.name)

	close(c.broadcast)
	close(c.register)
	close(c.unregister)
}

func IsChannelNameValid(name string) error {
	if len(name) < 1 || len(name) > 20 {
		return errors.New("Channel name must be between 1 and 20 characters")
	}

	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '-') {
			return errors.New("Channel name can only contain alphanumeric characters, underscores and dashes")
		}
	}

	return nil
}

func (c *Channel) run() {
	for {
		select {
		case client := <-c.register:
			c.clients[client] = true

			message := messages.CreateSwitchedChannelMessage(c.name)
			client.send <- message
		case client := <-c.unregister:
			if _, ok := c.clients[client]; ok {
				delete(c.clients, client)
			}
		case message := <-c.broadcast:
			for client := range c.clients {
				select {
				case client.send <- message:
				default:
					client.conn.Close()
					delete(c.clients, client)
				}
			}
		}
	}
}
