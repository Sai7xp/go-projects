package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan Message)           // Broadcast channel

// Message struct
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages
	go handleMessages()

	// Start server
	log.Println("Server started on :4444")
	err := http.ListenAndServe(":4444", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Register new client
	clients[ws] = true

	for {
		var msg Message
		// Read in new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send message to every connected client
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
