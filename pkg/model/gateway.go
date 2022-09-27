package model

import (
	"net"

	"github.com/lab5e/lospan/pkg/protocol"
)

// Gateway represents - you guessed it - a gateway.
type Gateway struct {
	GatewayEUI protocol.EUI // EUI of gateway.
	IP         net.IP       // IP address of gateway. This might not be fixed.
	StrictIP   bool         // Strict IP address check
	Latitude   float32      // Latitude, in decimal degrees, positive N <-90-90>
	Longitude  float32      // Longitude, in decimal degrees, positive E [-180-180>
	Altitude   float32      // Altitude, meters
}

// NewGateway creates a new gateway
func NewGateway() Gateway {
	return Gateway{}
}

// Equals checks gateways for equality
func (g *Gateway) Equals(other Gateway) bool {
	return g.Altitude == other.Altitude &&
		g.GatewayEUI.String() == other.GatewayEUI.String() &&
		g.IP.Equal(other.IP) &&
		g.Latitude == other.Latitude &&
		g.Longitude == other.Longitude &&
		g.StrictIP == other.StrictIP
}
