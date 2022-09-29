package main

import (
	"fmt"
	"os"

	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/pb/lospan"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := CommandLineParameters.Valid(); err != nil {
		lg.Error("Invalid configuration: %v", err)
		os.Exit(1)
	}

	mode := &BatchMode{Config: CommandLineParameters}

	conn, err := grpc.Dial("127.0.0.1:4711", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}
	client := lospan.NewLospanClient(conn)

	CommandLineParameters.Hostname = "127.0.0.1"
	lg.Info("gRPC API UDP: %s:%d", CommandLineParameters.Hostname, CommandLineParameters.UDPPort)

	e1 := Eagle1{
		Client:         client,
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
