package main

import (
	"github.com/valyala/fastjson"
	"ws-chat-server/server_message"
)

type ClientMessage struct {
	Type      string `json:"type"`
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

func parseClientMessage(message []byte) (*ClientMessage, error) {
	var p fastjson.Parser
	v, err := p.Parse(string(message))
	if err != nil {
		return nil, err
	}

	clientMessage := &ClientMessage{
		Type:      string(v.GetStringBytes("type")),
		Data:      string(v.GetStringBytes("data")),
		Timestamp: v.GetInt64("timestamp"),
	}

	return clientMessage, nil
}

func handleClientMessage(c *Client, clientMessage *ClientMessage) {
	switch clientMessage.Type {
	case "sendMessage":
		// Handle sendMessage type
		// You can add your logic here
		break
	case "switchChannel":
		// Handle switchChannel type
		// You can add your logic here
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
