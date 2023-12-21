package main

import (
	"github.com/valyala/fastjson"
)

type ClientMessage struct {
	Type      string `json:"type"`
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

type ServerMessage struct {
	Type      string `json:"type"`
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

type ServerMessageType string

const (
	serverMessage         ServerMessageType = "serverMessage"
	clientSwitchedChannel ServerMessageType = "clientSwitchedChannel"
)

type ServerMessageRoot struct {
	Messages []ServerMessage `json:"messages"`
}

func parseMessage(message []byte) (*ClientMessage, error) {
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
		c.send <- []byte(message)
		break
	}
}
