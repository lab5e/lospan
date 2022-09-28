package storage

import (
	"net"
	"testing"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/stretchr/testify/require"
)

// GatewayStorageTest is a generic test for gateway storage
func testGatewayStorage(gwStorage *Storage, t *testing.T) {
	assert := require.New(t)

	// Retrieve the empty list
	gateways, err := gwStorage.GetGatewayList()
	assert.NoError(err, "Should not get error when list is 0")
	assert.Len(gateways, 0, "List should be empty")

	// Create a new gateway
	gw1EUI, _ := protocol.EUIFromString("00-01-02-03-04-05-06-07")
	gateway1 := model.Gateway{
		GatewayEUI: gw1EUI,
		IP:         net.ParseIP("127.0.0.1"),
		StrictIP:   false,
		Latitude:   63.0,
		Longitude:  10.0,
		Altitude:   50.0,
	}

	assert.NoError(gwStorage.CreateGateway(gateway1), "No error when storing gateway")

	assert.Error(gwStorage.CreateGateway(gateway1), "Should get error when gateway exists")

	// ...and another one
	gw2EUI, _ := protocol.EUIFromString("aa-01-02-03-04-05-06-07")
	gateway2 := model.Gateway{
		GatewayEUI: gw2EUI,
		IP:         net.ParseIP("127.0.0.2"),
		StrictIP:   true,
		Latitude:   -63.0,
		Longitude:  -10.0,
		Altitude:   0.0,
	}

	assert.NoError(gwStorage.CreateGateway(gateway2), "Gateway 1 should be stored")

	// Retrieve the list
	gateways, err = gwStorage.GetGatewayList()
	assert.NoError(err)
	assert.Len(gateways, 2, "Should have 2 gateways")
	assert.Contains(gateways, gateway1)
	assert.Contains(gateways, gateway2)

	// Try adding the same gateway twice. Should yield error
	assert.Error(gwStorage.CreateGateway(gateway1), "Gateway already exists")

	// Retrieve just the first gateway. It should - of course - be the same.
	first, err := gwStorage.GetGateway(gateway1.GatewayEUI)
	assert.NoError(err)
	assert.Equal(gateway1, first)

	// Retrieving gateway that doesn't exist should yield error
	nonEUI, _ := protocol.EUIFromString("00-00-00-00-00-00-00-00")
	_, err = gwStorage.GetGateway(nonEUI)
	assert.Error(err)

	assert.NoError(gwStorage.UpdateGateway(gateway1), "Should update gateway")

	// Update fields
	gateway1.Altitude = 111
	gateway1.Latitude = 222
	gateway1.Longitude = 333
	gateway1.IP = net.ParseIP("10.10.10.10")
	gateway1.StrictIP = true
	assert.NoError(gwStorage.UpdateGateway(gateway1), "Should update gateway")

	updatedGW, err := gwStorage.GetGateway(gateway1.GatewayEUI)
	assert.NoError(err)
	assert.Equal(gateway1, updatedGW)

	// Remove both
	assert.NoError(gwStorage.DeleteGateway(gateway1.GatewayEUI))
	assert.NoError(gwStorage.DeleteGateway(gateway2.GatewayEUI))

	// Remove one that isn't supposed to exist in the list
	assert.Error(gwStorage.DeleteGateway(gateway1.GatewayEUI), "Gateway does not exist")

	// Ensure list is empty again
	// Retrieve the empty list
	gateways, err = gwStorage.GetGatewayList()
	assert.NoError(err)
	assert.Len(gateways, 0)

}
