package main

import (
	"encoding/json"
	"log"
	"sync"

	robotGo "github.com/go-vgo/robotgo"
	webSocket "github.com/gorilla/websocket"
)

type MyData struct {
	X int `json:"x"`
	Y int `json:"y"`
}

var mutex sync.Mutex

func mouseHandler(messageType int, message string, wsConnection *webSocket.Conn) {
	var jsonResponse MyData
	mutex.Lock()
	defer mutex.Unlock()
	err := json.Unmarshal([]byte(message), &jsonResponse)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(jsonResponse)
	robotGo.MoveSmoothRelative(jsonResponse.X, jsonResponse.Y)
}
