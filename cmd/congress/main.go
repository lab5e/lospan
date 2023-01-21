package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/storage"
)

type params struct {
	Log  lg.LogParameters  `kong:"embed,prefix='log-'"`
	LoRa server.Parameters `kong:"embed,prefix='lora-'"`
}

var config = server.NewDefaultConfig()

func init() {
	flag.IntVar(&config.GatewayPort, "gwport", server.DefaultGatewayPort, "Port for gateway listener")
	flag.UintVar(&config.NetworkID, "netid", server.DefaultNetworkID, "The Network ID to use")
	flag.StringVar(&config.MA, "ma", server.DefaultMA, "MA to use when generating new EUIs")
	flag.StringVar(&config.ConnectionString, "connectionstring", ":memory:", "Database connection string")
	flag.BoolVar(&config.PrintSchema, "printschema", false, "Print schema definition")
	flag.BoolVar(&config.DisableGatewayChecks, "disablegwcheck", false, "Disable ALL gateway checks")
	flag.Parse()
}

func main() {
	var config params
	kong.Parse(&config)

	if config.LoRa.PrintSchema {
		fmt.Print(storage.DBSchema)
		return
	}
	lg.InitLogs("congress", config.Log)
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
