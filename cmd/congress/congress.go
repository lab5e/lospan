package main

import (
	"errors"

	"github.com/ExploratoryEngineering/pubsub"
	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/gateway"
	"github.com/lab5e/lospan/pkg/processor"
	"github.com/lab5e/lospan/pkg/restapi"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/storage"
)

// Server is the main Congress server process. It will launch several
// endpoints and a processing pipeline.
type Server struct {
	config     *server.Configuration
	context    *server.Context
	forwarder  processor.GwForwarder
	pipeline   *processor.Pipeline
	restapi    *restapi.Server
	terminator chan bool
}

func (c *Server) setupLogging() {
	lg.InitLogs("lospan", config.Log)
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
func NewServer(config *server.Configuration) (*Server, error) {
	c := &Server{config: config, terminator: make(chan bool)}
	c.setupLogging()

	if err := c.checkConfig(); err != nil {
		return nil, err
	}
	lg.Info("This is the Congress server")

	var datastore *storage.Storage
	var err error
	if c.config.DBConnectionString != "" {
		lg.Info("Using PostgreSQL as backend storage")
		datastore, err = storage.CreateStorage(config.DBConnectionString,
			config.DBMaxConnections, config.DBIdleConnections, config.DBConnLifetime)
		if err != nil {
			lg.Error("Couldn't connect to database: %v", err)
			return nil, err
		}
	}

	keyGenerator, err := server.NewEUIKeyGenerator(config.RootMA(), uint32(config.NetworkID), datastore)
	if err != nil {
		lg.Error("Could not create key generator: %v. Terminating.", err)
		return nil, errors.New("unable to create key generator")
	}
	frameOutput := server.NewFrameOutputBuffer()

	appRouter := pubsub.NewEventRouter(5)
	gwEventRouter := pubsub.NewEventRouter(5)
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
	c.restapi, err = restapi.NewServer(config.OnlyLoopback, c.context, c.config)
	if err != nil {
		lg.Error("Unable to create REST API endpoint: %v", err)
		return nil, err
	}

	return c, nil
}

// Start Starts the congress server
func (c *Server) Start() error {
	lg.Debug("Starting pipeline")
	c.pipeline.Start()
	lg.Debug("Starting forwarder")
	go c.forwarder.Start()

	lg.Debug("Launching http server")
	if err := c.restapi.Start(); err != nil {
		lg.Error("Unable to start REST API endpoint: %v", err)
		return err
	}
	lg.Info("Server is ready and serving HTTP on port %d", c.config.HTTPServerPort)

	return nil
}

// Shutdown stops the Congress server.
func (c *Server) Shutdown() error {
	c.forwarder.Stop()
	c.restapi.Shutdown()
	c.context.Storage.Close()

	return nil
}
