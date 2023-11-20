package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/gorilla/mux"
)

// KeyboardEvent represents the structure of the JSON request body
type KeyboardEvent struct {
	Key string `json:"key"`
}

// SimulateKeyPress simulates a key press using the appropriate command based on the operating system
func SimulateKeyPress(key string) {
	switch runtime.GOOS {
	case "windows":
		// Simulate key press on Windows
		exec.Command("cmd", "/c", "echo "+key+" | clip").Run()
		exec.Command("cmd", "/c", "echo ^v | clip").Run()
	case "darwin":
		// Simulate key press on macOS
		exec.Command("bash", "-c", "osascript -e 'tell application \"System Events\" to keystroke \""+key+"\"'").Run()
	case "linux":
		// Simulate key press on Linux
		exec.Command("bash", "-c", "xdotool type "+key).Run()
	}
}

// KeyboardEventHandler handles incoming keyboard events
func KeyboardEventHandler(w http.ResponseWriter, r *http.Request) {
	var event KeyboardEvent

	// Decode the JSON request body
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&event)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simulate key press based on the received event
	SimulateKeyPress(event.Key)

	fmt.Printf("Received key event: %s\n", event.Key)
	w.WriteHeader(http.StatusOK)
}

type BasicLogData struct {
	Message string `json:"message"`
}

// DefaultHandler handles GET requests on the default path ("/")
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// Create an instance of your custom data structure
	data := BasicLogData{
		Message: "3010 Port LazyPanda running!",
	}

	// Convert the data to JSON
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the client
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {
	router := mux.NewRouter()

	// Endpoint to return JSON on the default path
	router.HandleFunc("/", DefaultHandler).Methods("GET")

	// Endpoint to receive keyboard events
	router.HandleFunc("/keyboard-event", KeyboardEventHandler).Methods("POST")

	port := "3010"
	fmt.Printf("Server is running on port %s\n", port)
	http.ListenAndServe(":"+port, router)
}
