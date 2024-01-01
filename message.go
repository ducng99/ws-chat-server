package main

import (
	"github.com/valyala/fastjson"
	"ws-chat-server/server_message"
)

type ClientMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func parseClientMessage(message []byte) (*ClientMessage, error) {
	var p fastjson.Parser
	v, err := p.Parse(string(message))
	if err != nil {
		return nil, err
	}

	clientMessage := &ClientMessage{
		Type: string(v.GetStringBytes("type")),
		Data: string(v.GetStringBytes("data")),
	}

	return clientMessage, nil
}

func handleClientMessage(c *Client, clientMessage *ClientMessage) {
	switch clientMessage.Type {
	case "sendMessage":
		message := server_message.CreateUserMessage(c.username, clientMessage.Data)
		c.channel.broadcast <- message
		break
	case "switchChannel":
		channel, err := ChatServer.getChannel(clientMessage.Data)

		if err != nil {
			message := err.Error()
			serverMessage := server_message.CreateServerMessage(message)
			c.send <- serverMessage
		} else {
			c.channel.unregister <- c
			c.channel = channel
			channel.register <- c
		}
		break
	case "ping":
		// Handle ping type
		// You can add your logic here
		break
	default:
		message := "Unknown message type: " + clientMessage.Type
		serverMessage := server_message.CreateServerMessage(message)

		c.send <- serverMessage
		break
	}
}
