package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	webSocket "github.com/gorilla/websocket"
)

type BasicLogData struct {
	Message string `json:"message"`
}

var jwtSecret = []byte(globalConfig.JwtSecret)
var upGrader = webSocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow any origin for WebSocket connections
	},
}

var defaultAllowedUsers = map[string]string{
	"shiv": "P@ssw0rd",
}

type KeyboardEvent struct {
	Key string `json:"key"`
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := defaultAllowedUsers[creds.Username]

	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		response := struct {
			InvalidLogin bool
		}{true}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)
		return
	}
	token, err := generateToken(creds.Username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonToken := struct{ Token string }{token}
	marshalJsonToken, err := json.Marshal(jsonToken)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Write(marshalJsonToken)
}

func generateToken(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() // 7 days expiration

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	data := BasicLogData{
		Message: "Lazy Panda running on port " + globalConfig.Server.Port,
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

func authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/login" || r.URL.Path == "/" || r.URL.Path == "/ws" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized Request", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized Request", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	httpToWsConnection, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer httpToWsConnection.Close() //Closing the http connection

	for {
		messageType, message, err := httpToWsConnection.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(message), messageType)

		go mouseHandler(messageType, string(message), httpToWsConnection)
	}
}
