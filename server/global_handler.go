package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Ngrok struct {
		AuthToken string `json:"auth_token"`
		Domain    string `json:"domain"`
	} `json:"ngrok"`

	JwtSecret string `json:"jwt_secret"`

	Server struct {
		Port      string `json:"port"`
		IPAddress string `json:"ipAddress"`
	} `json:"server"`

	SysTrayIconPath string `json:"sysTrayIconPath"`
}

func readConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&globalConfig)
	if err != nil {
		return err
	}

	return nil
}
