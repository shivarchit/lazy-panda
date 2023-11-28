package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.ngrok.com/ngrok"

	ngRokConfig "golang.ngrok.com/ngrok/config"
)

var keyPressChannel = make(chan string)

var globalConfig Config

func main() {
	err := readConfig("./config.json")
	if err != nil {
		fmt.Println("Error reading configuration:", err)
		return
	}

	sysTrayQuit = make(chan struct{})

	ctx := context.Background()
	listener, err := ngrok.Listen(ctx,
		ngRokConfig.HTTPEndpoint(
			ngRokConfig.WithDomain(globalConfig.Ngrok.Domain),
		),
		ngrok.WithAuthtoken(globalConfig.Ngrok.AuthToken),
	)
	if err != nil {
		log.Fatal("Error setting up ngrok:", err)
	}

	router := mux.NewRouter()
	router.Use(authenticateMiddleware)

	router.HandleFunc("/", DefaultHandler).Methods("GET")
	router.HandleFunc("/api/keyboard-event", KeyboardEventHandler).Methods("POST")
	router.HandleFunc("/api/login", LoginHandler).Methods("POST")

	port := globalConfig.Server.Port
	ipAddress := globalConfig.Server.IPAddress
	addr := ipAddress + ":" + port

	fmt.Printf("Server is running on %s\n", addr)

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
