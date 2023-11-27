package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	sysTray "github.com/getlantern/systray"
)

var sysTrayQuit chan struct{}

func onReady() {
	sysTray.SetTitle("Lazy Panda")
	sysTray.SetTooltip("Lazy Panda Server")
	mQuit := sysTray.AddMenuItem("Quit Lazy Panda Server", "Quit the application")

	iconPath := globalConfig.SysTrayIconPath
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
	go sysTray.Run(onReady, onExit)
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
