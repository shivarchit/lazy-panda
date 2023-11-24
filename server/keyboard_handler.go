package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"

	robotGo "github.com/go-vgo/robotgo"
)

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

func SimulateKeyPress(key string) {
	log.Println("Running for architecture: ", runtime.GOOS, "\n")
	log.Println("Pressing key: ", key, "\n")
	robotGo.Sleep(1)
	robotGo.KeyTap(key)
	robotGo.Sleep(1)
}
