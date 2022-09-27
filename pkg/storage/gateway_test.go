package storage

import (
	"net"
	"testing"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
)

// GatewayStorageTest is a generic test for gateway storage
func testGatewayStorage(gwStorage *Storage, t *testing.T) {
	// Retrieve the empty list
	gwChan, err := gwStorage.GetGatewayList()
	if err != nil {
		t.Fatal("Got error retrieving empty list: ", err)
	}
	count := 0
	for range gwChan {
		count++
	}
	if count > 0 {
		t.Fatalf("Expected 0 elements but got %d", count)
	}

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

	if err := gwStorage.CreateGateway(gateway1); err != nil {
		t.Fatal("Got error storing gateway: ", err)
	}

	if err := gwStorage.CreateGateway(gateway1); err != ErrAlreadyExists {
		t.Fatal("Should get ErrAlreadyExists when storing same gateway twice")
	}
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

	if err := gwStorage.CreateGateway(gateway2); err != nil {
		t.Fatal("Got error storing gateway: ", err)
	}

	// Retrieve the list
	gwChan, err = gwStorage.GetGatewayList()
	if err != nil {
		t.Fatal("Got error retrieving list: ", err)
	}

	var foundOne, foundTwo bool
	for val := range gwChan {
		if gateway1.Equals(val) {
			foundOne = true
		}
		if gateway2.Equals(val) {
			foundTwo = true
		}
	}

	if !foundOne || !foundTwo {
		t.Fatal("One or both gateways missing from list")
	}

	// Try adding the same gateway twice. Should yield error
	if err := gwStorage.CreateGateway(gateway1); err == nil {
		t.Fatal("Expected error when adding gateway twice")
	}

	// Retrieve just the first gateway. It should - of course - be the same.
	first, err := gwStorage.GetGateway(gateway1.GatewayEUI)
	if err != nil {
		t.Fatal("Did not expect an error")
	}
	if !gateway1.Equals(first) {
		t.Fatalf("Gateway did not match %v != %v", gateway1, first)
	}

	// Retrieving gateway that doesn't exist should yield error
	nonEUI, _ := protocol.EUIFromString("00-00-00-00-00-00-00-00")
	_, err = gwStorage.GetGateway(nonEUI)
	if err == nil {
		t.Fatal("Expected error when retrieving gw that doesn't exist")
	}

	if err := gwStorage.UpdateGateway(gateway1); err != nil {
		t.Fatalf("Got error storing tags on gateway: %v", err)
	}

	// Update fields
	gateway1.Altitude = 111
	gateway1.Latitude = 222
	gateway1.Longitude = 333
	gateway1.IP = net.ParseIP("10.10.10.10")
	gateway1.StrictIP = true
	if err := gwStorage.UpdateGateway(gateway1); err != nil {
		t.Fatalf("Got error updating gateway: %v", err)
	}
	updatedGw, _ := gwStorage.GetGateway(gateway1.GatewayEUI)
	if updatedGw.Altitude != gateway1.Altitude || updatedGw.Longitude != gateway1.Longitude || updatedGw.IP.String() != gateway1.IP.String() || updatedGw.StrictIP != gateway1.StrictIP {
		t.Fatalf("Gateways doesn't match! %v != %v", updatedGw, gateway1)
	}
	// Remove both
	if err := gwStorage.DeleteGateway(gateway1.GatewayEUI); err != nil {
		t.Fatalf("Got error removing gateway #1: %v", err)
	}
	if err := gwStorage.DeleteGateway(gateway2.GatewayEUI); err != nil {
		t.Fatalf("Got error removing gateway #2: %v", err)
	}
	// Remove one that isn't supposed to exist in the list
	if err := gwStorage.DeleteGateway(gateway1.GatewayEUI); err == nil {
		t.Fatal("Expected error when deleting gateway a second time")
	}

	// Ensure list is empty again
	// Retrieve the empty list
	gwChan, err = gwStorage.GetGatewayList()
	if err != nil {
		t.Fatal("Got error retrieving empty list: ", err)
	}
	count = 0
	for range gwChan {
		count++
	}
	if count > 0 {
		t.Fatalf("Got more than 0 elements (got %d)", count)
	}

}
