package model

import (
	"testing"

	"github.com/lab5e/lospan/pkg/protocol"
)

func TestUpstreamMessageCompare(t *testing.T) {
	d1 := UpstreamMessage{DeviceEUI: protocol.EUIFromInt64(0), Data: []byte{1, 2, 3}, Frequency: 99.0, GatewayEUI: protocol.EUIFromInt64(1)}
	d2 := UpstreamMessage{DeviceEUI: protocol.EUIFromInt64(1), Data: []byte{1, 2, 3}, Frequency: 98.0, GatewayEUI: protocol.EUIFromInt64(1)}
	d3 := UpstreamMessage{DeviceEUI: protocol.EUIFromInt64(0), Data: []byte{1, 2, 3}, Frequency: 99.0, GatewayEUI: protocol.EUIFromInt64(1)}

	if d1.Equals(d2) || d2.Equals(d1) {
		t.Fatal("Should not be the same")
	}

	if !d1.Equals(d3) || !d3.Equals(d1) {
		t.Fatal("Should be equal")
	}
}
