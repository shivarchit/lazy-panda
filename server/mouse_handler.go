package main

import (
	"encoding/json"
	"log"
	"sync"

	robotGo "github.com/go-vgo/robotgo"
	webSocket "github.com/gorilla/websocket"
)

type JSONStruct struct {
	X int `json:"x"`
	Y int `json:"y"`
}

var mutex sync.Mutex

func mouseHandler(messageType int, message string, wsConnection *webSocket.Conn) {
	var jsonResponse JSONStruct
	mutex.Lock()
	err := json.Unmarshal([]byte(message), &jsonResponse)
	if err != nil {
		log.Fatal(err)
		defer mutex.Unlock()
	}
	robotGo.MoveSmoothRelative(jsonResponse.X, jsonResponse.Y)
	defer mutex.Unlock()
}
