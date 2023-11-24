package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	sysTray "github.com/getlantern/systray"
	"github.com/gorilla/mux"
	"github.com/skratchdot/open-golang/open"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"

	robotGo "github.com/go-vgo/robotgo"
)

var sysTrayQuit chan struct{}
var keyPressChannel chan string

// KeyboardEvent represents the structure of the JSON request body
type KeyboardEvent struct {
	Key string `json:"key"`
}

// SimulateKeyPress simulates a key press using the appropriate command based on the operating system
func SimulateKeyPress(key string) {
	switch runtime.GOOS {
	case "windows":

		log.Println(key)
		robotGo.Sleep(1)
		robotGo.KeyTap("key")
		robotGo.Sleep(2)
		// cmd := exec.Command("cmd", "/c", "echo "+key+"| clip")
		// if err := cmd.Run(); err != nil {
		// 	log.Println("Error running command:", err)
		// }

		// cmd = exec.Command("cmd", "/c", "echo ^v | clip")
		// if err := cmd.Run(); err != nil {
		// 	log.Println("Error running command:", err)
		// }
	case "darwin":
		// cmd := exec.Command("bash", "-c", "osascript -e 'tell application \"System Events\" to keystroke \""+key+"\"'")
		// if err := cmd.Run(); err != nil {
		// 	log.Println("Error running command:", err)
		// }
	case "linux":
		// cmd := exec.Command("bash", "-c", "xdotool type "+key)
		// if err := cmd.Run(); err != nil {
		// 	log.Println("Error running command:", err)
		// }
	default:
		log.Println("Unsupported operating system:", runtime.GOOS)
	}
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
	fmt.Printf("1")
	keyPressChannel <- event.Key
	fmt.Printf("2")

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

func main() {
	keyPressChannel = make(chan string)
	sysTrayQuit = make(chan struct{})

	os.Setenv("NGROK_AUTHTOKEN", "")
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

	router.HandleFunc("/", DefaultHandler).Methods("GET")
	router.HandleFunc("/api/keyboard-event", KeyboardEventHandler).Methods("POST")

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
	open.Run(listener.URL())

	<-sysTrayQuit
}
