package messages

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
	Timestamp int64             `json:"timestamp"`
}

type msg_ServerMessage struct {
	baseServerMessageStruct
	Message string `json:"message"`
}

func (s msg_ServerMessage) getType() ServerMessageType {
	return s.Type
}

type msg_SwitchedChannel struct {
	baseServerMessageStruct
	Channel string `json:"channel"`
}

func (s msg_SwitchedChannel) getType() ServerMessageType {
	return s.Type
}

type msg_UserMessage struct {
	baseServerMessageStruct
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

func (s msg_UserMessage) getType() ServerMessageType {
	return s.Type
}

func CreateMessages(messages []ServerMessageInterface) ServerMessages {
	_serverMessages := ServerMessages{
		Messages: messages,
	}

	return _serverMessages
}

func CreateServerMessage(message string) *msg_ServerMessage {
	_serverMessage := &msg_ServerMessage{
		baseServerMessageStruct: baseServerMessageStruct{
			Type:      serverMessage,
			Timestamp: time.Now().UnixMilli(),
		},
		Message: message,
	}

	return _serverMessage
}

func CreateSwitchedChannelMessage(channel string) *msg_SwitchedChannel {
	_serverMessage := &msg_SwitchedChannel{
		baseServerMessageStruct: baseServerMessageStruct{
			Type:      clientSwitchedChannel,
			Timestamp: time.Now().UnixMilli(),
		},
		Channel: channel,
	}

	return _serverMessage
}

func CreateUserMessage(sender string, message string) *msg_UserMessage {
	_serverMessage := &msg_UserMessage{
		baseServerMessageStruct: baseServerMessageStruct{
			Type:      userMessage,
			Timestamp: time.Now().UnixMilli(),
		},
		Sender:  sender,
		Message: message,
	}

	return _serverMessage
}
