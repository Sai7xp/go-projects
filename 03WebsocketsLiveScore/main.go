/*
* Created on 02 March 2024
* @author Sai Sumanth
 */

package main

import (
	"cricscore/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

// score database
var scoreDatabase = make(map[string]models.TeamScore)
var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[*websocket.Conn]bool) // Map of Connections

func setUpRoutes() {
	/// live score websocket
	http.HandleFunc("/ws", liveScoreWebSocketHandler)

	/// User home page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	/// Admin Portal
	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/admin.html")
	})
}

// entry point
func main() {
	fmt.Println("Starting Live Score Web Server")

	setUpRoutes()

	/// start server 🛜
	logger.Println("Server started at port 6060")
	log.Fatal(http.ListenAndServe(":6060", nil))
}

// live score websocket handler /ws 🤝
func liveScoreWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println("New user Opened Live Score Dashboard")

	// Upgrade upgrades the HTTP server connection to the WebSocket protocol
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	/// register new client
	clients[conn] = true
	defer conn.Close()

	/// read live score updates from admin
	for {
		messageType, scoreInBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("read failed:", err)
			break
		}
		/// parse score
		var updatedScore models.TeamScore
		json.Unmarshal(scoreInBytes, &updatedScore)

		// update score in local database
		scoreDatabase[updatedScore.TeamName] = updatedScore
		logger.Printf("Updated new score in database %+v\n", scoreDatabase)

		/// write the updated score back
		newScore := fmt.Sprintf("\n %s\n%d - %d", updatedScore.TeamName,
			updatedScore.TotalScore, updatedScore.TotalWickets)
		newScoresBytes := []byte(newScore)
		broadcastScore(messageType, newScoresBytes)

	}
}

// Broadcast scores to all connections 📣
func broadcastScore(messageType int, newScoreBytes []byte) {
	for eachClient := range clients {
		writeError := eachClient.WriteMessage(messageType, newScoreBytes)
		if writeError != nil {
			logger.Println("write failed:", writeError)
			eachClient.Close()
			// Remove Client
			delete(clients, eachClient)
		}
	}
}
