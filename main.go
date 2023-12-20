package main

import (
	"flag"
	"log"
	"net/http"
	"sync"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")

type Server struct {
	hubs  map[string]*Hub
	mutex sync.RWMutex
}

func newServer() *Server {
	return &Server{
		hubs: make(map[string]*Hub),
	}
}

func (s *Server) getHub(name string) *Hub {
	s.mutex.RLock()
	hub, ok := s.hubs[name]
	s.mutex.RUnlock()

	if !ok {
		s.mutex.Lock()
		hub = newHub()
		s.hubs[name] = hub
		s.mutex.Unlock()
	}

	return hub
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	channelName := r.URL.Query().Get("channel")
	hub := s.getHub(channelName)

	serveWs(hub, w, r)
}

func main() {
	flag.Parse()
	chatServer := newServer()

	http.HandleFunc("/", chatServer.handleWebSocket)

	httpServer := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}