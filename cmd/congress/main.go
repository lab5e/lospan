package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/storage"
)

var config = server.NewDefaultConfig()

func init() {
	flag.IntVar(&config.GatewayPort, "gwport", server.DefaultGatewayPort, "Port for gateway listener")
	flag.IntVar(&config.HTTPServerPort, "http", server.DefaultHTTPPort, "HTTP port to listen on")
	flag.UintVar(&config.NetworkID, "netid", server.DefaultNetworkID, "The Network ID to use")
	flag.StringVar(&config.MA, "ma", server.DefaultMA, "MA to use when generating new EUIs")
	flag.StringVar(&config.DBConnectionString, "connectionstring", ":memory:", "Database connection string")
	flag.BoolVar(&config.PrintSchema, "printschema", false, "Print schema definition")
	flag.BoolVar(&config.DisableGatewayChecks, "disablegwcheck", false, "Disable ALL gateway checks")
	flag.BoolVar(&config.UseSecureCookie, "securecookie", false, "Set the secure flag for the auth cookie")
	flag.BoolVar(&config.MemoryDB, "memorydb", true, "Use in-memory database for storage (for testing)")
	flag.Parse()
}

func main() {
	if config.PrintSchema {
		fmt.Print(storage.DBSchema)
		return
	}
	lg.InitLogs("congress", config.Log)
	congress, err := NewServer(config)
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
	signal.Notify(sigch)
	go func() {
		sig := <-sigch
		lg.Debug("Caught signal '%v'", sig)
		terminator <- true
	}()

	<-terminator

}
