package main

import (
	"errors"
	"ws-chat-server/server_message"
)

// Channel maintains the set of active clients and broadcasts messages to the
// clients.
type Channel struct {
	name string
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan server_message.ServerMessageInterface

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewChannel(name string) *Channel {
	return &Channel{
		name:       name,
		broadcast:  make(chan server_message.ServerMessageInterface),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
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

func (h *Channel) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			message := server_message.CreateSwitchedChannelMessage(h.name)
			client.send <- message
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
