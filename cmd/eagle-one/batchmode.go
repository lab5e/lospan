package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/pb/lospan"
)

// BatchMode processing. Create
type BatchMode struct {
	Config           Params
	Application      *lospan.Application
	devices          []*lospan.Device
	OutgoingMessages chan string
	Publisher        *EventRouter
}

// Prepare prepares the processing
func (b *BatchMode) Prepare(client lospan.LospanClient, app *lospan.Application, gw *lospan.Gateway) error {
	b.devices = make([]*lospan.Device, 0)

	randomizer := NewRandomizer(50)

	ctx, done := context.WithTimeout(context.Background(), time.Minute*2)
	defer done()
	for i := 0; i < int(b.Config.DeviceCount); i++ {
		t := lospan.DeviceState_OTAA
		randomizer.Maybe(func() {
			t = lospan.DeviceState_ABP
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
	}
	lg.Info("# devices: %d", b.Config.DeviceCount)
	lg.Info("# messages: %d (total: %d)", b.Config.DeviceMessages, b.Config.DeviceCount*b.Config.DeviceMessages)
	return nil
}

// Cleanup resources after use. Remove devices if required.
func (b *BatchMode) Cleanup(client lospan.LospanClient, app *lospan.Application, gw *lospan.Gateway) {
	ctx, done := context.WithCancel(context.Background())
	defer done()
	if !b.Config.KeepDevices {
		lg.Info("Removing %d devices", len(b.devices))
		for _, d := range b.devices {
			client.DeleteDevice(ctx, &lospan.DeleteDeviceRequest{
				Eui: *d.Eui,
			})
		}
	}
}

// The number of join attempts before giving up
const joinAttempts = 5

func (b *BatchMode) launchDevice(device *lospan.Device, wg *sync.WaitGroup) {
	keys, err := NewDeviceKeys(device)
	if err != nil {
		lg.Warning("Got error converting lassie data into proper types: %v", err)
	}

	generator := NewMessageGenerator(b.Config)
	remoteDevice := NewEmulatedDevice(
		b.Config,
		keys,
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
		if err := remoteDevice.SendMessageWithGenerator(generator); err != nil {
			lg.Warning("Device %s got error sending message #%d", device.GetEui(), i)
		}
		randomOffset := rand.Intn(b.Config.TransmissionDelay/10) - (b.Config.TransmissionDelay / 5)
		lg.Debug("Device %s has sent message %d of %d", device.GetEui(), i, b.Config.DeviceMessages)
		time.Sleep(time.Duration(b.Config.TransmissionDelay+randomOffset) * time.Millisecond)
	}
	lg.Info("Device %s has completed", device.GetEui())
}

// Run the device emulation
func (b *BatchMode) Run(outgoingMessages chan string, publisher *EventRouter, app *lospan.Application, gw *lospan.Gateway) {
	b.OutgoingMessages = outgoingMessages
	b.Publisher = publisher
	b.Application = app

	// Power up our simulated devices
	completeWg := &sync.WaitGroup{}
	completeWg.Add(len(b.devices))
	for _, dev := range b.devices {
		go b.launchDevice(dev, completeWg)
	}
	lg.Info("....waiting for %d devices to complete", len(b.devices))
	completeWg.Wait()
}

// Failed returns true if the mode has failed
func (b *BatchMode) Failed() bool {
	return false
}
