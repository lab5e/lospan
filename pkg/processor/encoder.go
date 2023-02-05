package processor

import (
	"time"

	"github.com/lab5e/lospan/pkg/lg"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/storage"
)

// Encoder receives LoRaMessage data structures on a channel, encodes into a
// binary buffer and sends the buffer as a GatewayPacket instance on a new
// channel.
type Encoder struct {
	input   <-chan server.LoRaMessage
	output  chan<- server.GatewayPacket
	context *server.Context
}

func (e *Encoder) processMessage(packet server.LoRaMessage) {
	var buffer []byte
	var err error

	switch packet.Payload.MHDR.MType {

	case protocol.JoinRequest:
		lg.Warning("Unsupported encoding: JoinRequest (context=%v)", packet.FrameContext)

	case protocol.UnconfirmedDataUp:
		lg.Warning("Unsupported encoding: UnconfirmedDataUp (context=%v)", packet.FrameContext)

	case protocol.ConfirmedDataUp:
		lg.Warning("Unsupported encoding: ConfirmedDataUp (context=%v)", packet.FrameContext)

	case protocol.RFU:
		lg.Warning("Unsupported encoding: RFU(context=%v)", packet.FrameContext)

	case protocol.Proprietary:
		lg.Warning("Unsupported encoding: Proprietary message (context=%v)", packet.FrameContext)

	case protocol.JoinAccept:
		// Reset frame counter for both
		packet.FrameContext.Device.FCntDn = 0
		packet.FrameContext.Device.FCntUp = 0
		if err := e.context.Storage.UpdateDeviceState(packet.FrameContext.Device); err != nil {
			lg.Warning("Unable to update frame counters for device with EUI %s: %v. Ignoring JoinRequest.", packet.FrameContext.Device.DeviceEUI, err)
			return
		}

		buffer, err = packet.Payload.EncodeJoinAccept(packet.FrameContext.Device.AppKey)
		if err != nil {
			lg.Warning("Unable to encode JoinAccept message for device with EUI %s (DevAddr=%s): %v",
				packet.FrameContext.Device.DeviceEUI,
				packet.FrameContext.Device.DevAddr,
				err)
			return
		}
		packet.FrameContext.GatewayContext.Radio.RX1Delay = 5
		packet.FrameContext.GatewayContext.Deadline = 5

	default:
		packet.Payload.MACPayload.FHDR.FCnt = packet.FrameContext.Device.FCntDn
		buffer, err = packet.Payload.EncodeMessage(packet.FrameContext.Device.NwkSKey, packet.FrameContext.Device.AppSKey)
		if err != nil {
			lg.Error("Unable to encode message for device with EUI %s: %v. (DevAddr=%s)",
				packet.FrameContext.Device.DeviceEUI,
				err,
				packet.FrameContext.Device.DevAddr)
			return
		}

		// Update the sent state for the device. The message might be confirmed or unconfirmed at this point
		// but we don't care. We just send it and set the sent time. The downstream frame counter is updated
		// at this time so it will refer to the current frame counter
		if err := e.context.Storage.SetMessageSentTime(
			packet.FrameContext.Device.DeviceEUI,
			packet.FrameContext.PayloadCreate,
			time.Now().UnixNano(),
			packet.FrameContext.Device.FCntUp); err != nil && err != storage.ErrNotFound {
			lg.Warning("Unable to update downstream message for device %s: %v", packet.FrameContext.Device.DeviceEUI, err)
		}

		// Increase the frame counter after the message is sent. New devices will get 0,1,2...
		packet.FrameContext.Device.FCntDn++
		if err := e.context.Storage.UpdateDeviceState(packet.FrameContext.Device); err != nil {
			lg.Error("Unable to update frame counter for downstream message to device with EUI %s: %v",
				packet.FrameContext.Device.DeviceEUI,
				err)
		}
		packet.FrameContext.GatewayContext.Radio.RX1Delay = 1
		packet.FrameContext.GatewayContext.Deadline = 1
	}

	if len(buffer) == 0 {
		return
	}

	// Copy relevant data to the outgoing packet.
	e.output <- server.GatewayPacket{
		RawMessage: buffer,
		Radio:      packet.FrameContext.GatewayContext.Radio,
		Gateway:    packet.FrameContext.GatewayContext.Gateway,
		ReceivedAt: packet.FrameContext.GatewayContext.ReceivedAt,
		Deadline:   packet.FrameContext.GatewayContext.Deadline,
	}
}

// Start starts the Encoder instance. It will terminate when the input channel
// is closed. The output channel is closed when the method stops. The input channel
// receives messages due to be sent to gateways a short time before the messages
// must be sent from the gateway.
func (e *Encoder) Start() {
	for packet := range e.input {
		go e.processMessage(packet)

	}
	lg.Debug("Input channel for Encoder closed. Terminating")
}

// NewEncoder creates a new Encoder instance.
func NewEncoder(context *server.Context, input <-chan server.LoRaMessage, output chan<- server.GatewayPacket) *Encoder {
	return &Encoder{
		context: context,
		input:   input,
		output:  output,
	}
}
