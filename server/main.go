package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/getlantern/systray"
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
		Message: "3010 Port Lazy Panda running!",
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

func onReady() {
	// Add a simple menu item
	systray.SetTitle("Lazy Panda")
	systray.SetTooltip("Lazy Panda Server")
	mQuit := systray.AddMenuItem("Quit Lazy Panda Server", "Quit the application")

	// Set the icon (assuming icon.ico is in the same directory as main.go)
	iconPath := "panda.ico"
	iconBytes, err := os.ReadFile(iconPath)
	if err != nil {
		fmt.Println("Error reading icon file:", err)
		return
	}

	// Handle menu item clicks
	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
				return
			}
		}
	}()

	systray.SetIcon(iconBytes)
}

func onExit() {
	// Cleanup code when the systray is closed
}

func main() {
	router := mux.NewRouter()

	// Endpoint to return JSON on the default path
	router.HandleFunc("/api", DefaultHandler).Methods("GET")

	// Endpoint to receive keyboard events
	router.HandleFunc("/api/keyboard-event", KeyboardEventHandler).Methods("POST")

	port := "3010"
	fmt.Printf("Server is running on port %s\n", port)
	systray.Run(onReady, onExit)
	http.ListenAndServe(":"+port, router)
}
