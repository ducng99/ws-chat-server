package server_message

import (
	"time"
)

type ServerMessageType string

const (
	serverMessage         ServerMessageType = "serverMessage"
	clientSwitchedChannel ServerMessageType = "switchedChannel"
	userMessage           ServerMessageType = "userMessage"
)

// This is just so we can accept multiple structs
type ServerMessageInterface interface {
	getType() ServerMessageType
}

type ServerMessages struct {
	Messages []ServerMessageInterface `json:"messages"`
}

type baseServerMessageStruct struct {
	Type      ServerMessageType `json:"type"`
	Message   string            `json:"message"`
	Timestamp int64             `json:"timestamp"`
}

type serverMessageServerMessage struct {
	baseServerMessageStruct
}

func (s serverMessageServerMessage) getType() ServerMessageType {
	return s.Type
}

type serverMessageSwitchedChannel struct {
	baseServerMessageStruct
	Channel string `json:"channel"`
}

func (s serverMessageSwitchedChannel) getType() ServerMessageType {
	return s.Type
}

type serverMessageUserMessage struct {
	baseServerMessageStruct
	Sender string `json:"sender"`
}

func (s serverMessageUserMessage) getType() ServerMessageType {
	return s.Type
}

func CreateMessages(messages []ServerMessageInterface) ServerMessages {
	_serverMessages := ServerMessages{
		Messages: messages,
	}

	return _serverMessages
}

func CreateServerMessage(message string) *serverMessageServerMessage {
	_serverMessage := &serverMessageServerMessage{
		baseServerMessageStruct: baseServerMessageStruct{
			Type:      serverMessage,
			Message:   message,
			Timestamp: time.Now().UnixMilli(),
		},
	}

	return _serverMessage
}

func CreateSwitchedChannelMessage(channel string) *serverMessageSwitchedChannel {
	_serverMessage := &serverMessageSwitchedChannel{
		baseServerMessageStruct: baseServerMessageStruct{
			Type:      clientSwitchedChannel,
			Message:   "Switched to channel #" + channel,
			Timestamp: time.Now().UnixMilli(),
		},
		Channel: channel,
	}

	return _serverMessage
}

func CreateUserMessage(sender string, message string) *serverMessageUserMessage {
	_serverMessage := &serverMessageUserMessage{
		baseServerMessageStruct: baseServerMessageStruct{
			Type:      userMessage,
			Message:   message,
			Timestamp: time.Now().UnixMilli(),
		},
		Sender: sender,
	}

	return _serverMessage
}
