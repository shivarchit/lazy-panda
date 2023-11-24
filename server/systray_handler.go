package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	sysTray "github.com/getlantern/systray"
)

var keyPressChannel chan string

func onReady() {
	sysTray.SetTitle("Lazy Panda")
	sysTray.SetTooltip("Lazy Panda Server")
	mQuit := sysTray.AddMenuItem("Quit Lazy Panda Server", "Quit the application")

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
				sysTray.Quit()
				close(sysTrayQuit)
				os.Exit(0)
				return
			}
		}
	}()

	sysTray.SetIcon(iconBytes)
}

// KeyboardEventHandler handles incoming keyboard events
func KeyboardEventHandler(w http.ResponseWriter, r *http.Request) {
	var event KeyboardEvent

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&event)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	keyPressChannel <- event.Key

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

func onExit() {
	fmt.Println("Exit cleanup complete.")
}

func handleSignals() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigCh:
		fmt.Println("Received termination signal. Cleaning up...")
		sysTray.Quit()
		close(sysTrayQuit)
		os.Exit(0)
	}
}
