package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/lab5e/lospan/pkg/lg"
	"github.com/lab5e/lospan/pkg/pb/lospan"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/server"
)

// Eagle1 is the main testing tool. It will manage all of the infrastructure
// with the application,  gateway and packet forwarding. Message routing is done
// through the event router. It will publish events based on the device address.
type Eagle1 struct {
	Client         lospan.LospanClient
	Config         EagleConfig
	Application    *lospan.Application
	Gateway        *lospan.Gateway
	Publisher      *server.EventRouter[protocol.DevAddr, GWMessage]
	GatewayChannel chan string
	forwarder      *SyntheticForwarder
	shutdown       chan bool
}

func (e *Eagle1) newRandomEUI() string {
	octets := make([]byte, 8)
	rand.Read(octets)
	return fmt.Sprintf("%02x-%02x-%02x-%02x-%02x-%02x-%02x-%02x",
		octets[0], octets[1], octets[2], octets[3], octets[4], octets[5], octets[6], octets[7])
}

// Setup runs the setup procedures
func (e *Eagle1) Setup() error {
	ctx, done := context.WithTimeout(context.Background(), time.Minute)
	defer done()

	var err error
	if e.Config.Mode == "create" {
		if e.Application, err = e.Client.CreateApplication(ctx, &lospan.CreateApplicationRequest{}); err != nil {
			return fmt.Errorf("unable to create application in Congress: %v", err)
		}
	}

	strict := false
	lat := float32(50.3672)
	lon := float32(6.932)
	alt := float32(476.0)
	ip := "127.0.0.1"
	newGw := &lospan.Gateway{
		Eui:       e.newRandomEUI(),
		Ip:        &ip,
		StrictIp:  &strict,
		Latitude:  &lat,
		Longitude: &lon,
		Altitude:  &alt,
	}
	if e.Gateway, err = e.Client.CreateGateway(ctx, newGw); err != nil {
		return fmt.Errorf("unable to create gateway in Congress: %v", err)
	}
	lg.Info("Gateway EUI: %s", e.Gateway.GetEui())
	lg.Info("Application EUI: %s", e.Application.GetEui())

	return nil
}

// Teardown does a controlled terardown and removes application and gateway if needed.
func (e *Eagle1) Teardown() {
	e.shutdown <- true
}

// Run runs through the mode (batch/interactive)
func (e *Eagle1) Run(runner DeviceRunner) error {
	defer runner.Cleanup(e.Client, e.Application, e.Gateway)
	if err := runner.Prepare(e.Client, e.Application, e.Gateway); err != nil {
		return err
	}
	runner.Run(e.GatewayChannel, e.Publisher, e.Application, e.Gateway)
	return nil
}

func (e *Eagle1) decodingLoop() {
	for msg := range e.forwarder.OutputChannel() {
		p := protocol.NewPHYPayload(protocol.Proprietary)
		if err := p.UnmarshalBinary(msg); err != nil {
			lg.Warning("Unable to unmarshal message from gateway: %v", err)
			continue
		}
		e.Publisher.Publish(p.MACPayload.FHDR.DevAddr, GWMessage{PHYPayload: p, Buffer: msg})
	}
}

// StartForwarder launches a synthetic packet forwarder
func (e *Eagle1) StartForwarder() {
	e.shutdown = make(chan bool)
	e.forwarder = NewSyntheticForwarder(
		e.GatewayChannel, e.shutdown,
		e.Gateway.GetEui(), e.Config.Hostname,
		e.Config.UDPPort)

	go e.forwarder.Start()
	go e.decodingLoop()
}
