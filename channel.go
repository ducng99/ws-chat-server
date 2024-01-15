package main

import (
	"errors"
	"ws-chat-server/messages"
)

type ChannelType int

const (
	CHANNEL_TYPE_MULTI ChannelType = iota
	CHANNEL_TYPE_DIRECT
)

// Channel maintains the set of active clients and broadcasts messages to the
// clients.
type Channel struct {
	// Channel type
	channelType ChannelType

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

	// Additional direct channel data (1 to 1)
	ChannelDirect
}

type ChannelDirect struct {
	// Allows only 2 clients
	allowedClients [2]*Client
}

func NewMultiChannel(name string) *Channel {
	channel := &Channel{
		channelType: CHANNEL_TYPE_MULTI,
		name:        name,
		broadcast:   make(chan messages.ServerMessageInterface),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
	}

	go channel.run()

	return channel
}

func NewDirectChannel(client1 *Client, client2 *Client) *Channel {
	channel := &Channel{
		channelType: CHANNEL_TYPE_DIRECT,
		name:        "",
		broadcast:   make(chan messages.ServerMessageInterface),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		ChannelDirect: ChannelDirect{
			allowedClients: [2]*Client{client1, client2},
		},
	}

	go channel.run()

	return channel
}

// Gets channel name based on the client requesting
// Channel name will be prefixed with @ if the channel is a direct channel
// or # if the channel is a public or private channel.
// If the channel is a direct channel, the name of the other client is returned
// If the client is nil, the channel name is returned
func (c *Channel) GetName(client *Client) string {
	if client != nil && c.channelType == CHANNEL_TYPE_DIRECT {
		if c.ChannelDirect.allowedClients[0] == client {
			return "@" + c.ChannelDirect.allowedClients[1].username
		} else {
			return "@" + c.ChannelDirect.allowedClients[0].username
		}
	}

	return "#" + c.name
}

func (c *Channel) Close() {
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

			message := messages.CreateSwitchedChannelMessage(c.GetName(client))
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
