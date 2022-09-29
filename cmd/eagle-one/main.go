package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/pb/lospan"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// EagleConfig is the configuration structure
type EagleConfig struct {
	DeviceCount        int           `kong:"help='Number of devices to emulate',default=100"`
	DeviceMessages     int           `kong:"help='Number of messages to send before terminating device',default=10"`
	CorruptMIC         int           `kong:"help='Percentage of corrupt MIC messages (0-100)',default=0"`
	CorruptedPayload   int           `kong:"help='Percentage of corrupt payload (0-100)',default=0"`
	DuplicateMessages  int           `kong:"help='Percent of duplicated messages (0-100)',default=2"`
	TransmissionDelay  time.Duration `kong:"help='Transmission delay between messages',default='5s'"`
	UDPPort            int           `kong:"help='UDP port for gateway interface',default=8000"`
	Hostname           string        `kong:"help='Hostname for gateway interface',default='127.0.0.1'"`
	MaxPayloadSize     int           `kong:"help='Maximum payload size',default=222"`
	FrameCounterErrors int           `kong:"help='Frame counter errors (0-100)',default=0"`
	GRPCEndpoint       string        `kong:"help='gRPC API endpoint',default='127.0.0.1:4711'"`
}

var config EagleConfig

func main() {
	kong.Parse(&config)

	mode := &DeviceRunner{Config: config}

	conn, err := grpc.Dial(config.GRPCEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}
	client := lospan.NewLospanClient(conn)

	router := server.NewEventRouter[protocol.DevAddr, GWMessage](2)
	e1 := Eagle1{
		Client:         client,
		Config:         config,
		Publisher:      &router,
		GatewayChannel: make(chan string, 2),
	}
	if err := e1.Setup(); err != nil {
		lg.Error("Init error: %v", err)
		os.Exit(1)
	}
	defer e1.Teardown()

	e1.StartForwarder()

	if err := e1.Run(DeviceRunner{Config: config}); err != nil {
		lg.Error("Error running: %v", err)
		os.Exit(1)
	}

	if mode.Failed() {
		fmt.Println("Exiting with errors")
		os.Exit(1)
	}
	fmt.Println("Successful stop")
}
