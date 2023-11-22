package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

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
	response := struct{ IsActionSuccess bool }{true}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
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

var systrayQuit chan struct{}

func onReady() {
	systray.SetTitle("Lazy Panda")
	systray.SetTooltip("Lazy Panda Server")
	mQuit := systray.AddMenuItem("Quit Lazy Panda Server", "Quit the application")

	iconPath := "panda.ico"
	iconBytes, err := os.ReadFile(iconPath)
	if err != nil {
		fmt.Printf("Error reading icon file: %v\n", err)
		return
	}

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				close(systrayQuit)
				os.Exit(0)
				return
			}
		}
	}()

	systray.SetIcon(iconBytes)
}

func onExit() {
	fmt.Println("Exit cleanup complete.")
}

func handleSignals() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Wait for a signal to be received
	select {
	case <-sigCh:
		// Handle termination signals
		fmt.Println("Received termination signal. Cleaning up...")
		systray.Quit()
		close(systrayQuit)
		os.Exit(0)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api", DefaultHandler).Methods("GET")
	router.HandleFunc("/api/keyboard-event", KeyboardEventHandler).Methods("POST")

	port := "3010"
	ipAddress := "192.168.1.6"
	addr := ipAddress + ":" + port

	fmt.Printf("Server is running on %s\n", addr)

	systrayQuit = make(chan struct{})

	go systray.Run(onReady, onExit)
	go handleSignals()

	err := http.ListenAndServe(addr, router)
	if err != nil {
		fmt.Printf("Error starting the server: %v\n", err)
		os.Exit(1)
	}

	<-systrayQuit // Wait for systray to be closed before exiting
}
