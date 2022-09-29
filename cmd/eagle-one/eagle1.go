package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/pb/lospan"
	"github.com/lab5e/lospan/pkg/protocol"
)

// Eagle1 is the main testing tool. It will manage all of the infrastructure
// with the application,  gateway and packet forwarding. Message routing is done
// through the event router. It will publish events based on the device address.
type Eagle1 struct {
	Client         lospan.LospanClient
	Config         Params
	Application    *lospan.Application
	Gateway        *lospan.Gateway
	Publisher      *EventRouter
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
	if e.Config.AppEUI == "" {
		if e.Application, err = e.Client.CreateApplication(ctx, &lospan.CreateApplicationRequest{}); err != nil {
			return fmt.Errorf("unable to create application in Congress: %v", err)
		}
	} else {
		e.Config.KeepApplication = true
		if e.Application, err = e.Client.GetApplication(ctx, &lospan.GetApplicationRequest{Eui: e.Config.AppEUI}); err != nil {
			return fmt.Errorf("couldn't read application %s: %v", e.Config.AppEUI, err)
		}
	}

	if e.Config.GatewayEUI == "" {
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
	} else {
		e.Config.KeepGateway = true
		if e.Gateway, err = e.Client.GetGateway(ctx, &lospan.GetGatewayRequest{Eui: e.Config.GatewayEUI}); err != nil {
			return fmt.Errorf("cannot retrieve a gateway with the EUI %s: %v", e.Config.GatewayEUI, err)
		}
	}
	lg.Info("Gateway EUI: %s", e.Gateway.GetEui())
	lg.Info("Application EUI: %s", e.Application.GetEui())

	return nil
}

// Teardown does a controlled terardown and removes application and gateway if needed.
func (e *Eagle1) Teardown() {
	e.shutdown <- true
	ctx, done := context.WithTimeout(context.Background(), time.Minute)
	defer done()
	if !e.Config.KeepApplication {
		lg.Info("Removing application %s", e.Application.GetEui())
		e.Client.DeleteApplication(ctx, &lospan.DeleteApplicationRequest{Eui: e.Application.GetEui()})
	}
	if !e.Config.KeepGateway {
		lg.Info("Removing gateway %s", e.Gateway.GetEui())
		e.Client.DeleteGateway(ctx, &lospan.DeleteGatewayRequest{Eui: e.Gateway.GetEui()})
	}

}

// Run runs through the mode (batch/interactive)
func (e *Eagle1) Run(mode E1Mode) error {
	defer mode.Cleanup(e.Client, e.Application, e.Gateway)
	if err := mode.Prepare(e.Client, e.Application, e.Gateway); err != nil {
		return err
	}
	mode.Run(e.GatewayChannel, e.Publisher, e.Application, e.Gateway)
	return nil
}

func (e *Eagle1) decodingLoop() {
	for msg := range e.forwarder.OutputChannel() {
		p := protocol.NewPHYPayload(protocol.Proprietary)
		if err := p.UnmarshalBinary(msg); err != nil {
			lg.Warning("Unable to unmarshal message from gateway: %v", err)
			continue
		}
		e.Publisher.Publish(p.MACPayload.FHDR.DevAddr, p, msg)
	}
}

// StartForwarder launches a synthetic packet forwarder
func (e *Eagle1) StartForwarder() {
	e.shutdown = make(chan bool)
	e.forwarder = NewSyntheticForwarder(
		e.GatewayChannel, e.shutdown,
		e.Gateway.GetEui(), e.Config.Hostname,
		e.Config.UDPPort)

	lg.Info("Launching synthetic forwarder")
	go e.forwarder.Start()
	go e.decodingLoop()
}
