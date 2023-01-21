package main

import (
	"github.com/alecthomas/kong"
	"github.com/lab5e/lospan/pkg/congress"
	"github.com/lab5e/lospan/pkg/lg"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/utils"
)

type params struct {
	LoRa server.Parameters `kong:"embed,prefix='lora-'"`
}

func main() {
	var config params
	kong.Parse(&config)

	s, err := congress.NewLoRaServer(&config.LoRa)
	if err != nil {
		return
	}

	if err := s.Start(); err != nil {
		lg.Error("Congress did not start: %v", err)
		return
	}
	defer func() {
		lg.Info("Congress is shutting down...")
		s.Shutdown()
		lg.Info("Congress has shut down")
	}()

	utils.WaitForSignal()
}
