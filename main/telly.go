package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/sanjitsaluja/telly"

	_ "fmt"
	"os"
	"os/signal"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(telly.AppConfig.LoggingLevel)
}

func main() {
	sessionManager := telly.DefaultSessionManager
	// Listen for ^C for clean up
	listenForInterrupt()

	// Starting session manager
	sessionManager.Start()

	// Starting web server
	telly.StartServer()
}

func listenForInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			<-c
			log.Debug("Caught ^C interrupt")
			cleanup()
			os.Exit(1)
		}
	}()
}

func cleanup() {
	telly.DefaultSessionManager.Stop()
}
