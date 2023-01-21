package main

import (
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/lab5e/lospan/pkg/lg"
	"github.com/lab5e/lospan/pkg/server"
)

type params struct {
	LoRa server.Parameters `kong:"embed,prefix='lora-'"`
}

func main() {
	var config params
	kong.Parse(&config)

	congress, err := NewServer(&config.LoRa)
	if err != nil {
		return
	}

	terminator := make(chan bool)

	if err := congress.Start(); err != nil {
		lg.Error("Congress did not start: %v", err)
		return
	}
	defer func() {
		lg.Info("Congress is shutting down...")
		congress.Shutdown()
		lg.Info("Congress has shut down")
	}()

	sigch := make(chan os.Signal, 2)
	signal.Notify(sigch, os.Interrupt)
	go func() {
		sig := <-sigch
		lg.Debug("Caught signal '%v'", sig)
		terminator <- true
	}()

	<-terminator

}
