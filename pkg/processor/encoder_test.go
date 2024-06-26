package processor

import (
	"testing"
	"time"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/server"
)

func TestEncoderEncoding(t *testing.T) {
	s := NewStorageTestContext()
	d := model.Device{DevAddr: protocol.DevAddr{NwkID: 1, NwkAddr: 2}}
	s.CreateDevice(d, TestAppEUI)
	context := server.Context{
		Storage: s,
	}
	input := make(chan server.LoRaMessage)
	output := make(chan server.GatewayPacket)

	encoder := NewEncoder(&context, input, output)

	go encoder.Start()

	payload := protocol.NewPHYPayload(protocol.ConfirmedDataDown)
	input <- server.LoRaMessage{
		Payload: payload,
		FrameContext: server.FrameContext{
			Device:         d,
			GatewayContext: server.GatewayPacket{},
		},
	}

	select {
	case <-output:
	// This is OK
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Got timeout reading output channel")
	}

	select {
	case <-output:
		t.Fatal("Expected just a single response")
	case <-time.After(100 * time.Millisecond):
		// This is OK
	}

	close(input)
}

// Ensure the encoder doesn't encode upstream messages
func TestEncoderUnknownPackets(t *testing.T) {
	s := NewStorageTestContext()
	d := model.Device{DevAddr: protocol.DevAddr{NwkID: 1, NwkAddr: 2}}
	s.CreateDevice(d, TestAppEUI)
	context := server.Context{
		Storage: s,
	}
	input := make(chan server.LoRaMessage)
	output := make(chan server.GatewayPacket)

	encoder := NewEncoder(&context, input, output)

	go encoder.Start()

	// Build a simple message (just the message type is set)
	makeMessage := func(messageType protocol.MType) server.LoRaMessage {
		payload := protocol.NewPHYPayload(messageType)
		return server.LoRaMessage{
			Payload: payload,
			FrameContext: server.FrameContext{
				Device:         d,
				GatewayContext: server.GatewayPacket{},
			},
		}
	}
	// Ensure no output is received on the output channel
	ensureNoOutput := func() {
		select {
		case msg := <-output:
			t.Fatalf("Got message on output. Did not expect that! (message: %v)", msg)
		case <-time.After(10 * time.Millisecond):
			// OK - no message received
		}
	}

	input <- makeMessage(protocol.JoinRequest)
	ensureNoOutput()

	input <- makeMessage(protocol.UnconfirmedDataUp)
	ensureNoOutput()

	input <- makeMessage(protocol.ConfirmedDataUp)
	ensureNoOutput()

	input <- makeMessage(protocol.RFU)
	ensureNoOutput()

	input <- makeMessage(protocol.Proprietary)
	ensureNoOutput()
}

// Do a simple JoinAccept through the encoder
func TestJoinAcceptEncoder(t *testing.T) {
	s := NewStorageTestContext()
	a := model.NewApplication()
	d := model.NewDevice()
	d.DevAddr = protocol.DevAddr{NwkID: 1, NwkAddr: 2}
	s.CreateApplication(a)
	s.CreateDevice(d, TestAppEUI)

	context := server.Context{
		Storage: s,
	}
	input := make(chan server.LoRaMessage)
	output := make(chan server.GatewayPacket)

	encoder := NewEncoder(&context, input, output)

	go encoder.Start()

	payload := protocol.NewPHYPayload(protocol.JoinAccept)
	input <- server.LoRaMessage{
		Payload: payload,
		FrameContext: server.FrameContext{
			Device:         d,
			Application:    a,
			GatewayContext: server.GatewayPacket{},
		},
	}

	select {
	case <-output:
		// OK - no message received
	case <-time.After(10 * time.Millisecond):
		t.Fatalf("Message timed out")
	}
}
