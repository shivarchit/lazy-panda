package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

var keyPressChannel = make(chan string)

func main() {
	// sysTrayQuit = make(chan struct{})

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

	// go handleSignals()

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
	// <-sysTrayQuit
	select {}
}
