package main

import (
	"log"
	"runtime"

	robotGo "github.com/go-vgo/robotgo"
)

func SimulateKeyPress(key string) {
	log.Println("Running for architecture: ", runtime.GOOS, "\n")
	log.Println("Pressing key: ", key, "\n")
	robotGo.Sleep(1)
	robotGo.KeyTap(key)
	robotGo.Sleep(1)
}
