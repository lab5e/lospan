package model

import (
	"net"
	"testing"

	"github.com/lab5e/lospan/pkg/protocol"
)

func TestGatewayCompare(t *testing.T) {
	g1 := Gateway{GatewayEUI: protocol.EUIFromInt64(0), IP: net.ParseIP("127.0.0.1"), Altitude: 1}
	g2 := Gateway{GatewayEUI: protocol.EUIFromInt64(1), IP: net.ParseIP("127.1.0.1"), Altitude: 2}
	g3 := Gateway{GatewayEUI: protocol.EUIFromInt64(0), IP: net.ParseIP("127.0.0.1"), Altitude: 1}

	if g1.Equals(g2) || g2.Equals(g1) {
		t.Fatal("Should not be the same")
	}

	if !g1.Equals(g3) || !g3.Equals(g1) {
		t.Fatal("Should be equal")
	}

	g4 := NewGateway()
	g5 := NewGateway()

	if !g4.Equals(g5) {
		t.Fatal("Empty gateways should be equal")
	}

}
