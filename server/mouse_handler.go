package main

import (
	"encoding/json"
	"log"

	robotGo "github.com/go-vgo/robotgo"
	webSocket "github.com/gorilla/websocket"
)

type MyData struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func mouseHandler(messageType int, message string, wsConnection *webSocket.Conn) {
	var jsonResponse MyData
	// log.Println(string(message))
	err := json.Unmarshal([]byte(message), &jsonResponse)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(jsonResponse)
	// if string(message) != "" {
	// 	if err := wsConnection.WriteMessage(messageType, []byte("Message received")); err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// robotGo.Sleep(5)
	robotGo.MoveSmoothRelative(jsonResponse.X, jsonResponse.Y)
	// }
}
