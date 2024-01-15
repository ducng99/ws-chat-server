package main

import (
	"net/http"
	"strings"
	"sync"
)

var ChatServer = &Server{
	channels: make(map[string]*Channel),
	users:    make(map[string]*Client),
}

type Server struct {
	channels      map[string]*Channel
	channelsMutex sync.RWMutex

	users      map[string]*Client
	usersMutex sync.RWMutex
}

// getChannel returns a multi channel with the given name.
// If the channel does not exist, it creates a new one.
func (s *Server) getChannel(name string) (*Channel, error) {
	s.channelsMutex.RLock()
	channel, ok := s.channels[name]
	s.channelsMutex.RUnlock()

	if !ok {
		if err := IsChannelNameValid(name); err != nil {
			return nil, err
		}

		s.channelsMutex.Lock()
		channel = NewMultiChannel(name)
		s.channels[name] = channel
		s.channelsMutex.Unlock()
	}

	return channel, nil
}

func (s *Server) removeChannel(name string) {
	s.channelsMutex.Lock()
	delete(s.channels, name)
	s.channelsMutex.Unlock()
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

func (s *Server) getUser(username string) (*Client, bool) {
	s.usersMutex.RLock()
	user, ok := s.users[strings.ToLower(username)]
	s.usersMutex.RUnlock()

	return user, ok
}

func (s *Server) handleWebSocket(username string, w http.ResponseWriter, r *http.Request) {
	channel, err := s.getChannel("welcome")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	serveWs(username, channel, w, r)
}
