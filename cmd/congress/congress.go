package main

import (
	"errors"
	"net"
	"os"

	"github.com/lab5e/lospan/pkg/apiserver"
	"github.com/lab5e/lospan/pkg/events/gwevents"
	"github.com/lab5e/lospan/pkg/gateway"
	"github.com/lab5e/lospan/pkg/keys"
	"github.com/lab5e/lospan/pkg/lg"
	"github.com/lab5e/lospan/pkg/pb/lospan"
	"github.com/lab5e/lospan/pkg/processor"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/storage"
	"google.golang.org/grpc"
)

// Server is the main Congress server process. It will launch several
// endpoints and a processing pipeline.
type Server struct {
	config     *server.Parameters
	context    *server.Context
	forwarder  processor.GwForwarder
	pipeline   *processor.Pipeline
	terminator chan bool
}

func (c *Server) checkConfig() error {
	if err := c.config.Validate(); err != nil {
		lg.Error("Invalid configuration: %v Exiting", err)
		return errors.New("invalid configuration")
	}
	return nil
}

// NewServer creates a new server with the given configuration. The configuration
// is checked before the server is created, logging is initialized
func NewServer(config *server.Parameters) (*Server, error) {
	c := &Server{config: config, terminator: make(chan bool)}

	if err := c.checkConfig(); err != nil {
		return nil, err
	}
	lg.Info("This is the Congress server")

	var datastore *storage.Storage
	var err error
	if c.config.ConnectionString != "" {
		lg.Info("Using PostgreSQL as backend storage")
		datastore, err = storage.CreateStorage(config.ConnectionString)
		if err != nil {
			lg.Error("Couldn't connect to database: %v", err)
			return nil, err
		}
	}

	keyGenerator, err := keys.NewEUIKeyGenerator(config.RootMA(), uint32(config.NetworkID), datastore)
	if err != nil {
		lg.Error("Could not create key generator: %v. Terminating.", err)
		return nil, errors.New("unable to create key generator")
	}
	frameOutput := server.NewFrameOutputBuffer()

	appRouter := server.NewEventRouter[protocol.EUI, *server.PayloadMessage](5)
	gwEventRouter := server.NewEventRouter[protocol.EUI, gwevents.GwEvent](5)
	c.context = &server.Context{
		Storage:       datastore,
		Terminator:    make(chan bool),
		FrameOutput:   &frameOutput,
		Config:        config,
		KeyGenerator:  &keyGenerator,
		GwEventRouter: &gwEventRouter,
		AppRouter:     &appRouter,
	}

	lg.Info("Launching generic packet forwarder on port %d...", config.GatewayPort)
	c.forwarder = gateway.NewGenericPacketForwarder(c.config.GatewayPort, datastore, c.context)
	c.pipeline = processor.NewPipeline(c.context, c.forwarder)

	go func() {
		listener, err := net.Listen("tcp", ":4711")
		if err != nil {
			lg.Error("Error creating listener: %v", err)
			os.Exit(1)
		}
		lospanSvc, err := apiserver.New(c.context.Storage, c.context.KeyGenerator)
		if err != nil {
			lg.Error("Error creatig lospan service: %v", err)
			os.Exit(1)
		}
		lg.Info("Listening on %s", listener.Addr().String())
		server := grpc.NewServer()
		lospan.RegisterLospanServer(server, lospanSvc)
		if err := server.Serve(listener); err != nil {
			lg.Error("Error serving gRPC: %v", err)
			os.Exit(2)
		}
	}()
	return c, nil
}

// Start Starts the congress server
func (c *Server) Start() error {
	lg.Debug("Starting pipeline")
	c.pipeline.Start()
	lg.Debug("Starting forwarder")
	go c.forwarder.Start()

	return nil
}

// Shutdown stops the Congress server.
func (c *Server) Shutdown() error {
	c.forwarder.Stop()
	c.context.Storage.Close()

	return nil
}
