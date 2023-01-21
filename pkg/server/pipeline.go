package server

import (
	"time"

	"github.com/lab5e/lospan/pkg/band"
	"github.com/lab5e/lospan/pkg/events/gwevents"
	"github.com/lab5e/lospan/pkg/keys"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/storage"
)

// Pipeline data structures used by the server

// Context is the request/response context. It is passed along with the packets in various states.
type Context struct {
	Storage       *storage.Storage                             // The storage layer
	Terminator    chan bool                                    // Terminator channel. Throw something on this to terminate the processes.
	FrameOutput   *FrameOutputBuffer                           // Device aggregator instance. Common instance for processors.
	Config        *Parameters                                  // Main configuration
	KeyGenerator  *keys.KeyGenerator                           // Key generator for server
	GwEventRouter *EventRouter[protocol.EUI, gwevents.GwEvent] // Router for GW events
	AppRouter     *EventRouter[protocol.EUI, *PayloadMessage]  // Router for app data
}

// RadioContext - metadata for radio stats and settings
type RadioContext struct {
	Channel   uint8              // The channel used
	RFChain   uint8              // The RF chain the packet was received on
	Frequency float32            // Frequency - set by GW IF
	DataRate  string             // DataRate (f.e. "SF7BW125") - set by GW IF
	Band      band.FrequencyPlan // Band used
	RX1Delay  uint8              // RX1Delay - set during decoding
	RX2Delay  uint8              // RX2Delay - set during decoding
	RSSI      int32              // RSSI for device - set by GW IF
	SNR       float32            // SNR for device - set by GW IF
}

// GatewayContext - metadata for gateway; used when responding
type GatewayContext struct {
	GatewayEUI      protocol.EUI // The reported EUI
	GatewayHost     string       // The originating host
	GatewayPort     int          // The originating port
	GatewayClock    uint32       // Clock ticks reported by gateway
	ProtocolVersion uint8        // Protocol version (wrt packet forwarder)
}

// FrameContext is the context for each frame received (frequency, encoding, data rate rx1 offset and so on)
type FrameContext struct {
	Device         model.Device      // The decoded Device. Nil if it haven't been decoded yet.
	Application    model.Application // The decoded application. Nil if it haven't been resolved yet.
	GatewayContext GatewayPacket     // Context for gateway'
}

// GatewayPacket contains a byte buffer plus radio statistics.
type GatewayPacket struct {
	RawMessage []byte
	Radio      RadioContext
	Gateway    GatewayContext
	ReceivedAt time.Time
	Deadline   float64 // Send deadline for packet (in seconds)
}

// LoRaMessage contains the decoded LoRa message
type LoRaMessage struct {
	Payload      protocol.PHYPayload // PHYPayload decoded from GatewayPacket bytes.
	FrameContext FrameContext        // Frame context; set for each frame that arrives
}

// PayloadMessage contains the decrypted and verified payload
type PayloadMessage struct {
	Payload      []byte                // Unencrypted from the PHYPayload struct
	Device       model.Device          // The device that the payload was received from (or will be sent to)
	Application  model.Application     // The device's application.
	MACCommands  []protocol.MACCommand // MAC Commands received from/sent to the device
	FrameContext FrameContext          // The context the packet is received in
}
