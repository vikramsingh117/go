# Real-Time WAV to FLAC Audio Conversion Service

This project is a real-time audio streaming service that allows users to upload WAV audio files from their browser, stream them to a Go server, convert them to FLAC format using WebSocket communication, and download the converted file back to their device.

---

## Project Structure

```bash
wav-to-flac-service/
│
├── main.go                # Main Go server code
├── utils/
│   └── audio_converter.go  # Utility for WAV to FLAC conversion
├── routes/
│   └── routes.go
├── handlers/
│   └── audio_handler.go
└── client.html            # HTML file for the frontend interface
```

---

## Setup Instructions

### Prerequisites
- **Go (Golang)** installed on your machine.
- **FLAC Encoder** (`flac` command-line tool) installed and available in your system's PATH.
- Basic understanding of WebSocket communication.

---

### Installation

1. **Clone the repository**:
   ```sh
   git clone https://github.com/vikramsingh117/go.git
   cd go
   ```

2. **Install dependencies**:
   Make sure you have the `gorilla/websocket` package installed:
   ```sh
   go get github.com/gorilla/websocket
   ```

3. **Ensure FLAC Encoder is available**:
   - Check if `flac` is installed:
     ```sh
     flac --version
     ```
   - If not, install it using your package manager:
     - **Ubuntu**: `sudo apt-get install flac`
     - **macOS**: `brew install flac`

---

## Running the Server

1. Navigate to your project directory:
   ```sh
   cd path/to/wav-to-flac-service
   ```

2. Run the Go server:
   ```sh
   go run main.go
   ```

3. The server will start and listen on `localhost:8080`.

---

## Client Interface

- Open `client.html` in your web browser to access the frontend interface for uploading and streaming WAV files.

---

## Code Explanation

### `main.go`

- **Imports**:
  ```go
  import (
      "io/ioutil"
      "log"
      "net/http"
      "os"
      "time"
      "wav-to-flac-service/utils"
      "github.com/gorilla/websocket"
  )
  ```
  - **ioutil, os, and time**: File operations and time tracking.
  - **gorilla/websocket**: WebSocket handling.
  - **utils**: Custom package for audio conversion.

- **WebSocket Upgrader**:
  ```go
  var upgrader = websocket.Upgrader{
      CheckOrigin: func(r *http.Request) bool {
          return true
      },
  }
  ```

- **WebSocket Handler**:
  ```go
  func handleAudioStream(w http.ResponseWriter, r *http.Request) {
      conn, err := upgrader.Upgrade(w, r, nil)
      if err != nil {
          log.Println("Failed to upgrade to WebSocket:", err)
          return
      }
      defer conn.Close()

      for {
          startTime := time.Now()

          // Read audio data
          _, audioData, err := conn.ReadMessage()
          if err != nil {
              log.Println("Error reading message:", err)
              break
          }

          // Handle temporary files
          tempWavFile, _ := ioutil.TempFile("", "*.wav")
          defer os.Remove(tempWavFile.Name())
          tempWavFile.Write(audioData)
          tempWavFile.Close()

          tempFlacFile, _ := ioutil.TempFile("", "*.flac")
          defer os.Remove(tempFlacFile.Name())
          tempFlacFile.Close()

          // Convert WAV to FLAC
          err = utils.ConvertWAVToFLAC(tempWavFile.Name(), tempFlacFile.Name())
          if err != nil {
              log.Println("Conversion error:", err)
              break
          }

          flacData, _ := ioutil.ReadFile(tempFlacFile.Name())
          conn.WriteMessage(websocket.BinaryMessage, flacData)

          duration := time.Since(startTime)
          log.Printf("Time taken for transmission: %v\n", duration)
      }
  }
  ```

- **Starting the Server**:
  ```go
  func main() {
      http.HandleFunc("/stream", handleAudioStream)
      log.Println("Starting server on :8080...")
      log.Fatal(http.ListenAndServe(":8080", nil))
  }
  ```

---

### `utils/audio_converter.go`

- **ConvertWAVToFLAC Function**:
  ```go
  package utils

  import (
      "fmt"
      "os/exec"
  )

  func ConvertWAVToFLAC(inputPath, outputPath string) error {
      cmd := exec.Command("flac", inputPath, "-o", outputPath, "-f")
      err := cmd.Run()
      if err != nil {
          return fmt.Errorf("conversion failed: %v", err)
      }
      return nil
  }
  ```
  - Uses the `flac` command-line tool to perform the conversion.
  - `-f` flag forces overwrite of the output file if it exists.

---

## Client Interface (`client.html`)

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>WAV to FLAC Converter</title>
    <style>
        body { font-family: Arial; }
    </style>
</head>
<body>
    <h1>WAV to FLAC Converter</h1>
    <input type="file" id="fileInput" accept=".wav">
    <button id="startStreaming">Start Streaming</button>
    <div id="status">Status: Waiting...</div>

    <script>
        document.getElementById("startStreaming").addEventListener("click", () => {
            const file = document.getElementById("fileInput").files[0];
            if (!file) return alert("Please select a WAV file.");

            const ws = new WebSocket("ws://localhost:8080/stream");
            ws.onopen = () => {
                const reader = new FileReader();
                reader.onload = (e) => ws.send(e.target.result);
                reader.readAsArrayBuffer(file);
            };
            ws.onmessage = (e) => {
                const blob = new Blob([e.data], { type: "audio/flac" });
                const url = URL.createObjectURL(blob);
                const a = document.createElement("a");
                a.href = url;
                a.download = "converted.flac";
                a.click();
            };
            ws.onclose = () => console.log("Connection closed.");
        });
    </script>
</body>
</html>
```

---

## Logging and Debugging

- **Console Logs**: Added throughout the server and client for tracking connection status and data handling.
- **Error Handling**: Graceful error reporting in both the client and server.

---

## Improvements and Future Work

- **Chunked Data Transfer**: Implement chunked file reading for large audio files to reduce latency.
- **Progress Indicators**: Show upload and conversion progress on the client side.
- **Optimizations**: Optimize file I/O and command execution to improve performance.

---

## Conclusion

This service efficiently converts WAV audio files to FLAC in real-time using WebSocket communication between a browser-based client and a Go server. The provided documentation outlines the project setup, code explanation, and potential areas for enhancement.

---

**Author**: Vikram Singh  
**GitHub**: [Vikram Singh](https://github.com/vikramsingh117)

---

Feel free to contribute or suggest improvements! Happy coding!
