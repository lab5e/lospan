package model

import (
	"bytes"

	"github.com/lab5e/lospan/pkg/protocol"
)

// UpstreamMessage contains a single transmission from an end-device.
type UpstreamMessage struct {
	DeviceEUI  protocol.EUI     // Device address used
	Timestamp  int64            // Timestamp for message. Data type might change.
	Data       []byte           // The data the end-device sent
	GatewayEUI protocol.EUI     // The gateway the message was received from.
	RSSI       int32            // Radio stats; RSSI
	SNR        float32          // Radio; SNR
	Frequency  float32          // Radio; Frequency
	DataRate   string           // Data rate (ie "SF7BW125" or similar)
	DevAddr    protocol.DevAddr // The reported DevAddr (at the time)
}

// Equals compares two DeviceData instances
func (d *UpstreamMessage) Equals(other UpstreamMessage) bool {
	return bytes.Equal(d.Data, other.Data) &&
		d.DeviceEUI.ToInt64() == other.DeviceEUI.ToInt64() &&
		d.Timestamp == other.Timestamp &&
		d.GatewayEUI.String() == other.GatewayEUI.String() &&
		d.RSSI == other.RSSI &&
		d.SNR == other.SNR &&
		d.Frequency == other.Frequency &&
		d.DataRate == other.DataRate &&
		d.DevAddr == other.DevAddr

}
