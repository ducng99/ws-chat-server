package main

import (
	"net/http"
	"strings"
	"sync"
)

var ChatServer = newServer()

type Server struct {
	channels      map[string]*Channel
	channelsMutex sync.RWMutex

	users      map[string]*Client
	usersMutex sync.RWMutex
}

func newServer() *Server {
	return &Server{
		channels: make(map[string]*Channel),
		users:    make(map[string]*Client),
	}
}

// getChannel returns a channel with the given name.
// If the channel does not exist, it creates a new one.
func (s *Server) getChannel(name string) *Channel {
	s.channelsMutex.RLock()
	channel, ok := s.channels[name]
	s.channelsMutex.RUnlock()

	if !ok {
		s.channelsMutex.Lock()
		channel = NewChannel(name)
		s.channels[name] = channel
		s.channelsMutex.Unlock()

		go channel.run()
	}

	return channel
}

func (s *Server) addUser(username string, client *Client) {
	s.usersMutex.Lock()
	s.users[username] = client
	s.usersMutex.Unlock()
}

func (s *Server) removeUser(username string) {
	s.usersMutex.Lock()
	delete(s.users, username)
	s.usersMutex.Unlock()
}

func (s *Server) hasUser(username string) bool {
	s.usersMutex.RLock()
	_, ok := s.users[strings.ToLower(username)]
	s.usersMutex.RUnlock()

	return ok
}

func (s *Server) handleWebSocket(username string, w http.ResponseWriter, r *http.Request) {
	channel := s.getChannel("welcome")

	serveWs(username, channel, w, r)
}
