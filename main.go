package main

import (
	// "fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for testing purposes
	},
}

// WebSocket handler for audio streaming
func handleAudioStream(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected for audio streaming")

	for {
		// Read audio data from the WebSocket
		_, audioData, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Here you would process the audio data in WAV format and convert to FLAC.
		// This example just echoes the received data back.
		// You need to implement the WAV to FLAC conversion logic.

		// Send converted audio data back over the WebSocket
		err = conn.WriteMessage(websocket.BinaryMessage, audioData)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}

	log.Println("Client disconnected from audio streaming")
}

func main() {
	http.HandleFunc("/stream", handleAudioStream)

	port := ":8080"
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
