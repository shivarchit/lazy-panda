package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	sysTray "github.com/getlantern/systray"
	"github.com/gorilla/mux"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

var sysTrayQuit chan struct{}

type BasicLogData struct {
	Message string `json:"message"`
}

// DefaultHandler handles GET requests on the default path ("/")
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	data := BasicLogData{
		Message: "3010 Port Lazy Panda running!",
	}

	jsonResponse, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {
	keyPressChannel = make(chan string)
	sysTrayQuit = make(chan struct{})

	os.Setenv("NGROK_AUTHTOKEN", "2YXaL3cvlhjHmnAkimWzhdLfVw1_7XLqj3hHw68EYLfjpqFqf")
	ctx := context.Background()
	listener, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(
			config.WithDomain("starfish-hopeful-spaniel.ngrok-free.app"),
		),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		log.Fatal("Error setting up ngrok:", err)
	}

	router := mux.NewRouter()
	router.Use(authenticateMiddleware)

	router.HandleFunc("/", DefaultHandler).Methods("GET")
	router.HandleFunc("/api/keyboard-event", KeyboardEventHandler).Methods("POST")
	router.HandleFunc("/api/login", LoginHandler).Methods("POST")

	port := "3010"
	ipAddress := "localhost"
	addr := ipAddress + ":" + port

	fmt.Printf("Server is running on %s\n", addr)

	go sysTray.Run(onReady, onExit)
	go handleSignals()

	// Start the goroutine to handle key press requests
	go func() {
		for {
			select {
			case key := <-keyPressChannel:
				SimulateKeyPress(key)
			}
		}
	}()

	go func() {
		err := http.Serve(listener, router)
		if err != nil {
			log.Fatal("Error serving requests:", err)
		}
	}()

	// Open the ngrok URL in the default browser
	log.Println("Default Serving URL: ", listener.URL())
	<-sysTrayQuit
}
