package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"wav-to-flac-service/utils"

	"github.com/gorilla/websocket"
)

// Upgrader to upgrade HTTP connections to WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for testing purposes
	},
}

// WebSocket handler for audio streaming
func handleAudioStream(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a new connection request...")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer func() {
		log.Println("Closing WebSocket connection...")
		conn.Close()
	}()

	log.Println("Client connected for audio streaming")

	for {
		// Start timing
		startTime := time.Now()

		// Read audio data from the WebSocket
		_, audioData, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Calculate and log the file size in bytes
		fileSize := len(audioData)
		log.Printf("Received audio data of length: %d bytes\n", fileSize)

		// Write the received audio data to a temporary WAV file
		tempWavFile, err := ioutil.TempFile("", "audio_*.wav")
		if err != nil {
			log.Println("Error creating temporary WAV file:", err)
			break
		}
		defer os.Remove(tempWavFile.Name()) // Clean up the file after conversion

		if _, err := tempWavFile.Write(audioData); err != nil {
			log.Println("Error writing to temporary WAV file:", err)
			break
		}
		tempWavFile.Close()

		// Create a temporary file for the FLAC output
		tempFlacFile, err := ioutil.TempFile("", "audio_*.flac")
		if err != nil {
			log.Println("Error creating temporary FLAC file:", err)
			break
		}
		defer os.Remove(tempFlacFile.Name()) // Clean up the FLAC file after use

		// Check if the FLAC file already exists and delete it
		if _, err := os.Stat(tempFlacFile.Name()); err == nil {
			// The file exists, delete it
			err := os.Remove(tempFlacFile.Name())
			if err != nil {
				log.Println("Error deleting existing FLAC file:", err)
				break
			}
		}

		// Convert WAV to FLAC
		err = utils.ConvertWAVToFLAC(tempWavFile.Name(), tempFlacFile.Name())
		if err != nil {
			log.Println("Error during WAV to FLAC conversion:", err)
			break
		}

		// Read the converted FLAC data
		flacData, err := ioutil.ReadFile(tempFlacFile.Name())
		if err != nil {
			log.Println("Error reading converted FLAC file:", err)
			break
		}

		// Send the FLAC data back to the client
		log.Println("Sending converted FLAC data back to the client...")
		err = conn.WriteMessage(websocket.BinaryMessage, flacData)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}

		// Measure and log the time taken for transmission
		duration := time.Since(startTime)
		log.Printf("Time taken for transmission: %v\n", duration)
	}

	log.Println("Client disconnected from audio streaming")
}

// Main function to set up the HTTP server and routes
func main() {
	http.HandleFunc("/stream", handleAudioStream) // Set up route for WebSocket connections
	log.Println("Setting up routes...")
	log.Println("Starting server on port :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
