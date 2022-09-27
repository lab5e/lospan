package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/ExploratoryEngineering/logging"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/storage"
)

var config = server.NewDefaultConfig()

func init() {
	flag.IntVar(&config.GatewayPort, "gwport", server.DefaultGatewayPort, "Port for gateway listener")
	flag.IntVar(&config.HTTPServerPort, "http", server.DefaultHTTPPort, "HTTP port to listen on")
	flag.UintVar(&config.NetworkID, "netid", server.DefaultNetworkID, "The Network ID to use")
	flag.StringVar(&config.MA, "ma", server.DefaultMA, "MA to use when generating new EUIs")
	flag.StringVar(&config.DBConnectionString, "connectionstring", "", "Database connection string")
	flag.BoolVar(&config.PrintSchema, "printschema", false, "Print schema definition")
	flag.BoolVar(&config.Syslog, "syslog", false, "Send logs to syslog")
	flag.BoolVar(&config.DisableGatewayChecks, "disablegwcheck", false, "Disable ALL gateway checks")
	flag.BoolVar(&config.UseSecureCookie, "securecookie", false, "Set the secure flag for the auth cookie")
	flag.UintVar(&config.LogLevel, "loglevel", server.DefaultLogLevel, "Log level to use (0 = debug, 1 = info, 2 = warning, 3 = error)")
	flag.BoolVar(&config.PlainLog, "plainlog", false, "Use plain-text stderr logs")
	flag.BoolVar(&config.MemoryDB, "memorydb", true, "Use in-memory database for storage (for testing)")
	flag.IntVar(&config.DBMaxConnections, "db-max-connections", server.DefaultMaxConns, "Maximum DB connections")
	flag.IntVar(&config.DBIdleConnections, "db-max-idle-connections", server.DefaultIdleConns, "Maximum idle DB connections")
	flag.DurationVar(&config.DBConnLifetime, "db-max-lifetime-connections", server.DefaultConnLifetime, "Maximum life time of DB connections")
	flag.Parse()
}

func main() {
	if config.PrintSchema {
		fmt.Print(storage.DBSchema)
		return
	}
	logging.SetLogLevel(config.LogLevel)
	congress, err := NewServer(config)
	if err != nil {
		return
	}

	terminator := make(chan bool)

	if err := congress.Start(); err != nil {
		logging.Error("Congress did not start: %v", err)
		return
	}
	defer func() {
		logging.Info("Congress is shutting down...")
		congress.Shutdown()
		logging.Info("Congress has shut down")
	}()

	sigch := make(chan os.Signal, 2)
	signal.Notify(sigch, os.Interrupt, os.Kill)
	go func() {
		sig := <-sigch
		logging.Debug("Caught signal '%v'", sig)
		terminator <- true
	}()

	<-terminator

}
