package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")

// Handle simple authentication with username in query string
func handleAuthentication(w http.ResponseWriter, r *http.Request) {
	// Checks CORS
	if r.Header.Get("Origin") != "http://localhost:5173" && r.Header.Get("Origin") != "https://static.ducng.dev" {
		http.Error(w, "Invalid origin", http.StatusBadRequest)
		return
	} else {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "GET")
	}

	username := r.URL.Query().Get("username")

	if strings.ToLower(username) == "anonymous" {
		username = CreateAnonymousUsername()
	} else if err := IsUsernameValid(username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := CreateToken(username)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"username": username,
		"token":    token,
	}
	responseJSON, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	username, ok := GetUsernameFromToken(token)
	if !ok {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	ChatServer.handleWebSocket(username, w, r)
}

func main() {
	flag.Parse()

	http.HandleFunc("/login", handleAuthentication)
	http.HandleFunc("/ws", handleWebsocket)

	httpServer := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
