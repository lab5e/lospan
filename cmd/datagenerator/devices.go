package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/lab5e/lospan/pkg/keys"
	"github.com/lab5e/lospan/pkg/lg"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/storage"
)

func generateDevices(count int, app model.Application, datastore *storage.Storage, keyGen *keys.KeyGenerator, callback func(createdDevice model.Device)) {
	for i := 0; i < count; i++ {
		d := model.NewDevice()
		d.AppEUI = app.AppEUI
		if rand.Intn(100) < 50 {
			d.State = model.OverTheAirDevice
		} else {
			d.State = model.PersonalizedDevice
		}
		d.AppKey = randomAesKey()
		d.DevAddr = randomDevAddr()
		d.AppSKey = randomAesKey()
		d.NwkSKey = randomAesKey()
		d.FCntDn = uint16(rand.Intn(4096))
		d.FCntDn = uint16(rand.Intn(4096))
		d.DeviceEUI, _ = keyGen.NewDeviceEUI()
		d.RelaxedCounter = false
		if err := datastore.CreateDevice(d, app.AppEUI); err != nil {
			lg.Error("Unable to store device: %v", err)
		} else {
			callback(d)
		}
	}
}

const maxPayloadLength = 128 // TODO: increase to 220 when bug is fixed

func makeRandomPayload() []byte {
	buf := make([]byte, maxPayloadLength)
	rand.Read(buf)
	return buf[0 : 1+rand.Intn(maxPayloadLength-1)]
}

func randomDataRate() string {
	return fmt.Sprintf("DR%d", rand.Intn(7))
}

func randomFrequency() float32 {
	switch rand.Intn(7) {
	case 0:
		return 868.1
	case 1:
		return 868.3
	case 2:
		return 868.5
	case 3:
		return 867.1
	case 4:
		return 867.3
	case 5:
		return 867.5
	default:
		return 867.7
	}
}

func randomGateway(gws []model.Gateway) protocol.EUI {
	return gws[rand.Intn(len(gws))].GatewayEUI
}

func generateDeviceData(device model.Device, count int, gateways []model.Gateway, datastore *storage.Storage) {
	emulatedTime := time.Now().Add(-time.Duration(count) * time.Minute)
	for i := 0; i < count; i++ {
		dd := model.UpstreamMessage{}
		dd.Data = makeRandomPayload()
		dd.DataRate = randomDataRate()
		dd.DevAddr = device.DevAddr
		dd.DeviceEUI = device.DeviceEUI
		dd.Frequency = randomFrequency()
		dd.GatewayEUI = randomGateway(gateways)
		dd.RSSI = -int32(rand.Intn(120))
		dd.SNR = float32(rand.Intn(20))
		dd.Timestamp = emulatedTime.UnixNano()
		emulatedTime = emulatedTime.Add(time.Minute)
		if err := datastore.CreateUpstreamMessage(device.DeviceEUI, dd); err != nil {
			lg.Error("Unable to store device data: %v", err)
		}
	}
}

func generateDownstreamMessage(device model.Device, datastore *storage.Storage) {
	// About 1 in 2 have a downstream message waiting
	if rand.Intn(2) == 0 {
		dm := model.NewDownstreamMessage(device.DeviceEUI, uint8(1+rand.Intn(222)))
		dm.Ack = rand.Intn(2) == 1
		dm.Data = hex.EncodeToString(makeRandomPayload())
		if err := datastore.CreateDownstreamMessage(device.DeviceEUI, dm); err != nil {
			lg.Error("Unable to store downstream message: %v", err)
		}
	}
}

func randomNonce() uint16 {
	return uint16(rand.Uint32() & 0xFFFF)
}

func generateNonces(device model.Device, count int, datastore *storage.Storage) {
	if device.State == model.OverTheAirDevice {
		for i := 0; i < count; i++ {
			if err := datastore.AddDevNonce(device, randomNonce()); err != nil && err != storage.ErrAlreadyExists {
				lg.Warning("Unable to add nonce: %v", err)
			}
		}
	}
}
