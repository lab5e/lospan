package processor

import (
	"testing"
	"time"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/storage"
)

func TestOTAAJoinRequestProcessing(t *testing.T) {

	deviceEUI, _ := protocol.EUIFromString("00-01-02-03-04-05-06-07")
	appEUI, _ := protocol.EUIFromString("00-01-02-03-04-05-06-08")

	store := storage.NewMemoryStorage()
	application := model.Application{
		AppEUI: appEUI,
	}
	device := model.Device{
		DeviceEUI:       deviceEUI,
		AppEUI:          appEUI,
		State:           model.OverTheAirDevice,
		FCntUp:          100,
		FCntDn:          100,
		RelaxedCounter:  false,
		DevNonceHistory: make([]uint16, 0),
	}

	store.CreateApplication(application)
	store.CreateDevice(device, appEUI)

	inputChan := make(chan server.LoRaMessage)

	foBuffer := server.NewFrameOutputBuffer()
	decrypter := NewDecrypter(&server.Context{
		Storage:     store,
		FrameOutput: &foBuffer,
		Config:      &server.Parameters{},
	}, inputChan)

	payload := protocol.NewPHYPayload(protocol.JoinRequest)
	payload.JoinRequestPayload = protocol.JoinRequestPayload{
		DevEUI:   deviceEUI,
		AppEUI:   appEUI,
		DevNonce: 0x0102,
	}
	input := server.LoRaMessage{
		Payload: payload,
		FrameContext: server.FrameContext{
			Device:         model.NewDevice(),
			Application:    model.NewApplication(),
			GatewayContext: server.GatewayPacket{},
		},
	}

	close(inputChan)

	// This should result in an output message
	go decrypter.processJoinRequest(input)

	select {
	case <-decrypter.Output():
		// OK
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Did not get output on output channel!")
	}
}
