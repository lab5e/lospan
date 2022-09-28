package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/lab5e/l5log/pkg/lg"
	"github.com/telenordigital/lassie-go"
)

func main() {
	if err := CommandLineParameters.Valid(); err != nil {
		lg.Error("Invalid configuration: %v", err)
		os.Exit(1)
	}

	var mode E1Mode
	switch CommandLineParameters.Mode {
	case "batch":
		mode = &BatchMode{Config: CommandLineParameters}
	case "interactive":
		mode = &InteractiveMode{Config: CommandLineParameters}
	case "test":
		mode = &TestMode{Config: CommandLineParameters}
	default:
		lg.Error("Unknown mode: " + CommandLineParameters.Mode)
		os.Exit(1)
	}

	congress, err := lassie.New()
	if err != nil {
		lg.Error("Couldn't create the Congress API client: %v", err)
		os.Exit(1)
	}
	u, err := url.Parse(congress.Address())
	if err != nil {
		lg.Error("Invalid Congress URL: %v", err)
		os.Exit(1)
	}
	CommandLineParameters.Hostname = u.Hostname()

	lg.Info("Congress UDP: %s:%d", u.Hostname(), CommandLineParameters.UDPPort)
	lg.Info("Using Congress API at: %s", congress.Address())

	e1 := Eagle1{
		Congress:       congress,
		Config:         CommandLineParameters,
		Publisher:      NewEventRouter(2),
		GatewayChannel: make(chan string, 2),
	}
	if err := e1.Setup(); err != nil {
		lg.Error("Init error: %v", err)
		os.Exit(1)
	}
	defer e1.Teardown()

	e1.StartForwarder()

	if err := e1.Run(mode); err != nil {
		lg.Error("Unable to run %s mode: %v", CommandLineParameters.Mode, err)
		os.Exit(1)
	}

	if mode.Failed() {
		fmt.Println("Exiting with errors")
		os.Exit(1)
	}
	fmt.Println("Successful stop")
}
