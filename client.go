package main

import (
	"log"
	"net/http"
	"time"
	"ws-chat-server/messages"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize * 2,
	WriteBufferSize: maxMessageSize * 2,
	CheckOrigin:     checkWSOrigin,
}

func checkWSOrigin(r *http.Request) bool {
	return r.Header.Get("Origin") == "http://localhost:5173" || r.Header.Get("Origin") == "https://static.ducng.dev"
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// Username
	username string

	// Current chatting channel
	channel *Channel

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan messages.ServerMessageInterface
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.channel.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		chatMessage := &messages.ClientMessage{}

		if err := c.conn.ReadJSON(chatMessage); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		handleClientMessage(c, chatMessage)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			serverMessages := messages.CreateMessages([]messages.ServerMessageInterface{message})

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				message, ok := <-c.send
				if ok {
					serverMessages.Messages = append(serverMessages.Messages, message)
				}
			}

			if err := c.conn.WriteJSON(serverMessages); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) switchChannel(channel *Channel) {
	c.channel.unregister <- c
	c.channel = channel
	c.channel.register <- c
}

func handleClientMessage(c *Client, clientMessage *messages.ClientMessage) {
	switch clientMessage.Type {
	case "sendMessage":
		message := messages.CreateUserMessage(c.username, clientMessage.Data)
		c.channel.broadcast <- message
		break
	case "switchChannel":
		channel, err := ChatServer.getChannel(clientMessage.Data)

		if err != nil {
			message := err.Error()
			serverMessage := messages.CreateServerMessage(message)
			c.send <- serverMessage
		} else {
			c.switchChannel(channel)
		}
		break
	case "ping":
		// Handle ping type
		// You can add your logic here
		break
	default:
		message := "Unknown message type: " + clientMessage.Type
		serverMessage := messages.CreateServerMessage(message)

		c.send <- serverMessage
		break
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(username string, channel *Channel, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		username: username,
		channel:  channel,
		conn:     conn,
		send:     make(chan messages.ServerMessageInterface),
	}
	client.channel.register <- client

	ChatServer.addUser(client.username, client)

	conn.SetCloseHandler(func(code int, text string) error {
		close(client.send)
		ChatServer.removeUser(client.username)
		return nil
	})

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
