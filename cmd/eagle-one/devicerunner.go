package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/lab5e/lospan/pkg/lg"
	"github.com/lab5e/lospan/pkg/pb/lospan"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/server"
)

// GWMessage is the gateway message
type GWMessage struct {
	PHYPayload protocol.PHYPayload
	Buffer     []byte
}

// DeviceRunner processing. Create
type DeviceRunner struct {
	Config           EagleConfig
	Application      *lospan.Application
	devices          []*lospan.Device
	OutgoingMessages chan string
	Publisher        *server.EventRouter[protocol.DevAddr, GWMessage]
}

// Prepare prepares the processing
func (b *DeviceRunner) Prepare(client lospan.LospanClient, app *lospan.Application, gw *lospan.Gateway) error {
	b.devices = make([]*lospan.Device, 0)

	randomizer := NewRandomizer(b.Config.OTAA)

	ctx, done := context.WithTimeout(context.Background(), time.Minute*2)
	defer done()

	if b.Config.Mode == "create" {
		for i := 0; i < int(b.Config.DeviceCount); i++ {
			t := lospan.DeviceState_ABP
			randomizer.Maybe(func() {
				t = lospan.DeviceState_OTAA
			})
			eui := app.GetEui()
			newDevice := &lospan.Device{
				ApplicationEui: &eui,
				State:          &t,
			}
			dev, err := client.CreateDevice(ctx, newDevice)
			if err != nil {
				return fmt.Errorf("unable to create device in Congress: %v", err)
			}
			b.devices = append(b.devices, dev)
			lg.Info("Created device with EUI %s and DevAddr %08x", dev.GetEui(), dev.GetDevAddr())
		}
	} else {
		if b.Config.ApplicationEUI == "" {
			return errors.New("application EUI required when querying for devices")
		}
		devices, err := client.ListDevices(ctx, &lospan.ListDeviceRequest{
			ApplicationEui: b.Config.ApplicationEUI,
		})
		if err != nil {
			lg.Error("Error querying devices: %v", err)
			return err
		}
		b.devices = append(b.devices, devices.Devices...)
	}
	lg.Info("# devices: %d", b.Config.DeviceCount)
	lg.Info("# messages: %d (total: %d)", b.Config.DeviceMessages, b.Config.DeviceCount*b.Config.DeviceMessages)
	return nil
}

// Cleanup resources after use. Remove devices if required.
func (b *DeviceRunner) Cleanup(client lospan.LospanClient, app *lospan.Application, gw *lospan.Gateway) {

}

// The number of join attempts before giving up
const joinAttempts = 5

func (b *DeviceRunner) launchDevice(device *lospan.Device, wg *sync.WaitGroup) {
	keys, err := NewDeviceKeys(device)
	if err != nil {
		lg.Warning("Got error converting lassie data into proper types: %v", err)
	}

	generator := NewMessageGenerator(b.Config)
	remoteDevice := NewEmulatedDevice(
		b.Config,
		&keys,
		b.OutgoingMessages,
		b.Publisher)

	defer wg.Done()
	// Join if needed
	if device.GetState() == lospan.DeviceState_OTAA {
		if err := remoteDevice.Join(joinAttempts); err != nil {
			lg.Warning("Device %s couldn't join after %d attempts", device.GetEui(), joinAttempts)
			return
		}
	}

	lg.Debug("Device %s is now ready to send messages", device.GetEui())
	for i := 0; i < b.Config.DeviceMessages; i++ {
		randomOffset := rand.Int63n(int64(b.Config.TransmissionDelay/5)) - int64(b.Config.TransmissionDelay)/10
		time.Sleep(b.Config.TransmissionDelay + time.Duration(randomOffset))
		if err := remoteDevice.SendMessageWithGenerator(generator); err != nil {
			lg.Warning("Device %s got error sending message #%d", device.GetEui(), i)
		}
		lg.Debug("Device %s has sent message %d of %d", device.GetEui(), i, b.Config.DeviceMessages)
	}
	lg.Info("Device %s has completed", device.GetEui())
}

// Run the device emulation
func (b *DeviceRunner) Run(outgoingMessages chan string, publisher *server.EventRouter[protocol.DevAddr, GWMessage], app *lospan.Application, gw *lospan.Gateway) {
	b.OutgoingMessages = outgoingMessages
	b.Publisher = publisher
	b.Application = app

	// Power up our simulated devices
	completeWg := &sync.WaitGroup{}
	completeWg.Add(len(b.devices))
	for _, dev := range b.devices {
		time.Sleep(4 * time.Second / time.Duration(b.Config.DeviceCount))
		go b.launchDevice(dev, completeWg)
	}
	lg.Info("....waiting for %d devices to complete", len(b.devices))
	completeWg.Wait()
}

// Failed returns true if the mode has failed
func (b *DeviceRunner) Failed() bool {
	return false
}
