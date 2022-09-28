package model

import (
	"encoding/hex"
	"time"

	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/protocol"
)

// Data model for LoRaWAN network. The data model is used in parts of the decoding

// DownstreamMessageState is the state of the downstream messages
type DownstreamMessageState uint8

// States for the downstream message
const (
	UnsentState DownstreamMessageState = iota
	SentState
	AcknowledgedState
)

// DownstreamMessage is messages sent downstream (ie to devices from the server).
type DownstreamMessage struct {
	DeviceEUI protocol.EUI
	Data      string
	Port      uint8

	Ack         bool
	CreatedTime int64
	SentTime    int64
	AckTime     int64
}

// NewDownstreamMessage creates a new DownstreamMessage
func NewDownstreamMessage(deviceEUI protocol.EUI, port uint8) DownstreamMessage {
	return DownstreamMessage{deviceEUI, "", port, false, time.Now().Unix(), 0, 0}
}

// State returns the message's state based on the value of the time stamps
func (d *DownstreamMessage) State() DownstreamMessageState {
	// Sent time isn't updated => message is still pending
	if d.SentTime == 0 {
		return UnsentState
	}
	// Message isn't acknowledged but sent time is set => message is sent
	if d.AckTime == 0 {
		return SentState
	}
	// AckTime and SentTime is set => acknowledged
	return AcknowledgedState
}

// Payload returns the payload as a byte array. If there's an error decoding the
// data it will return an empty byte array
func (d *DownstreamMessage) Payload() []byte {
	ret, err := hex.DecodeString(d.Data)
	if err != nil {
		lg.Warning("Unable to decode data to be sent to device %s (data=%s). Ignoring it.", d.DeviceEUI, d.Data)
		return []byte{}
	}
	return ret
}

// IsComplete returns true if the message processing is completed. If the ack
// flag isn't set the message would only have to be sent to the device. If the
// ack flag is set the device must acknowledge the message before it is
// considered completed.
func (d *DownstreamMessage) IsComplete() bool {
	// Message haven't been sent yet
	if d.SentTime == 0 {
		return false
	}
	// Message have been sent but not acknowledged yet
	if d.Ack && d.AckTime == 0 {
		return false
	}
	return true
}
