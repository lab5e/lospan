package storage

import (
	"reflect"
	"testing"
	"time"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
)

func equalDevices(a model.Device, b model.Device) bool {
	return reflect.DeepEqual(a, b)
}

func testDeviceStorage(
	storage *Storage,
	t *testing.T) {

	app1 := model.Application{
		AppEUI: makeRandomEUI(),
	}
	if err := storage.CreateApplication(app1); err != nil {
		t.Error("Got error adding application 1: ", err)
	}

	app2 := model.Application{
		AppEUI: makeRandomEUI(),
	}
	if err := storage.CreateApplication(app2); err != nil {
		t.Error("Got error adding application 2: ", err)
	}

	deviceA := model.Device{
		DeviceEUI: makeRandomEUI(),
		AppEUI:    app1.AppEUI,
		DevAddr: protocol.DevAddr{
			NwkID:   1,
			NwkAddr: 0x000001,
		},
		FCntUp: 1,
	}
	deviceB := model.Device{
		DeviceEUI: makeRandomEUI(),
		AppEUI:    app1.AppEUI,
		DevAddr: protocol.DevAddr{
			NwkID:   1,
			NwkAddr: 0x000002,
		},
		FCntUp: 2,
	}
	deviceC := model.Device{
		DeviceEUI: makeRandomEUI(),
		AppEUI:    app2.AppEUI,
		DevAddr: protocol.DevAddr{
			NwkID:   1,
			NwkAddr: 0x000003,
		},
		FCntUp: 3,
	}
	deviceD := model.Device{
		DeviceEUI: makeRandomEUI(),
		AppEUI:    app2.AppEUI,
		DevAddr: protocol.DevAddr{
			NwkID:   1,
			NwkAddr: 0x000004,
		},
		FCntUp: 4,
	}

	err := storage.CreateDevice(deviceA, app1.AppEUI)
	if err != nil {
		t.Error("Error putting device 1: ", err)
	}
	err = storage.CreateDevice(deviceB, app1.AppEUI)
	if err != nil {
		t.Error("Error putting device 2: ", err)
	}

	err = storage.CreateDevice(deviceC, app2.AppEUI)
	if err != nil {
		t.Error("Error putting device 3: ", err)
	}
	err = storage.CreateDevice(deviceD, app2.AppEUI)
	if err != nil {
		t.Error("Error putting device 4: ", err)
	}
	// Retrieve one of the stored devices via DevAddr (assume the others work)
	if deviceChan, err := storage.GetDeviceByDevAddr(deviceC.DevAddr); err == nil {
		found := false
		for d := range deviceChan {
			t.Logf("Found device %+v", d)
			if equalDevices(d, deviceC) {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("Couldn't find device 3 in storage")
		}
	} else {
		t.Error("Could not retrieve device: ", err)
	}

	// ...and do the same for a device keyed on EUI
	device, err := storage.GetDeviceByEUI(deviceB.DeviceEUI)
	if err != nil {
		t.Error("Could not retrieve device via EUI: ", err)
	}
	if !equalDevices(device, deviceB) {
		t.Error("Returned device isn't the same")
	}
	deviceChan1, err := storage.GetDevicesByApplicationEUI(app1.AppEUI)
	if err != nil {
		t.Error("Got error retrieving devices: ", err)
	}
	deviceChan2, err := storage.GetDevicesByApplicationEUI(app2.AppEUI)
	if err != nil {
		t.Error("Got error retrieving devices: ", err)
	}
	if deviceChan1 == nil {
		t.Error("No channel 1 returned!")
		return
	}
	if deviceChan2 == nil {
		t.Error("No channel 2 returned!")
		return
	}
	app1Count := uint16(0)
	app2Count := uint16(0)
	for i := 0; i < 4; i++ {
		select {
		case dev := <-deviceChan1:
			app1Count += dev.FCntUp
		case dev := <-deviceChan2:
			app2Count += dev.FCntUp
		case <-time.After(time.Millisecond * 100):
			t.Error("Timeout getting devices")
		}
	}

	_, err = storage.GetDeviceByEUI(protocol.EUIFromInt64(0))
	if err == nil {
		t.Error("Expected error when querying for unknown eui")
	}
	// Try adding the same device twice
	err = storage.CreateDevice(deviceA, app1.AppEUI)
	if err == nil {
		t.Error("Should not be able to add device twice")
	}

	// Store device nonce on device, ensure it is stored
	if deviceA.HasDevNonce(12) {
		t.Error("Device A should not have device nonce")
	}

	if err = storage.AddDevNonce(deviceA, 12); err != nil {
		t.Fatal("Got error storing dev nonce: ", err)
	}
	if err = storage.AddDevNonce(deviceA, 24); err != nil {
		t.Fatal("Got error storing dev nonce: ", err)
	}
	if err = storage.AddDevNonce(deviceA, 48); err != nil {
		t.Fatal("Got error storing dev nonce: ", err)
	}

	if device, err = storage.GetDeviceByEUI(deviceA.DeviceEUI); err != nil {
		t.Fatal("Got error retrieving device: ", err)
	}

	if !device.HasDevNonce(12) {
		t.Fatal("Device did not store device nonce properly. Didn't find 12. Device nonces are: ", device.DevNonceHistory)
	}
	if !device.HasDevNonce(24) {
		t.Fatal("Device did not store device nonce properly. Didn't find 24. Device nonces are: ", device.DevNonceHistory)
	}

	if !device.HasDevNonce(48) {
		t.Fatal("Device did not store device nonce properly. Didn't find 48. Device nonces are: ", device.DevNonceHistory)
	}
	if device.HasDevNonce(96) {
		t.Fatal("Device claims to have a lot of nonces (asked for 96, can't see it). Device nonces are: ", device.DevNonceHistory)
	}

	deviceC.NwkSKey = makeRandomKey()
	deviceC.AppSKey = makeRandomKey()
	if err := storage.UpdateDevice(deviceC); err != nil {
		t.Fatal("Could not store the device keys: ", err)
	}

	if device, err = storage.GetDeviceByEUI(deviceC.DeviceEUI); err != nil {
		t.Fatal("could not retrieve device c via EUI: ", err)
	}

	if device.AppSKey != deviceC.AppSKey {
		t.Errorf("AppSKey for stored version isn't the same that was set %v != %v", device.AppSKey, deviceC.AppSKey)
	}
	if device.NwkSKey != deviceC.NwkSKey {
		t.Errorf("AppSKey for stored version isn't the same that was set %v != %v", device.NwkSKey, deviceC.NwkSKey)
	}

	deviceD.FCntDn = 1001
	deviceD.FCntUp = 2002
	deviceD.KeyWarning = true
	if err = storage.UpdateDeviceState(deviceD); err != nil {
		t.Fatal("Error updating frame counters: ", err)
	}

	updatedDevice, err := storage.GetDeviceByEUI(deviceD.DeviceEUI)
	if err != nil {
		t.Fatal("Could not retrieve updated frame counter device: ", err)
	}
	if updatedDevice.FCntDn != 1001 || updatedDevice.FCntUp != 2002 || !updatedDevice.KeyWarning {
		t.Fatalf("Device did not update the frame counters. Expected 1001, 2002 got %d, %d",
			updatedDevice.FCntDn, updatedDevice.FCntUp)
	}

	updatedDevice.DevAddr = protocol.DevAddrFromUint32(0x01020304)
	updatedDevice.RelaxedCounter = true
	updatedDevice.FCntDn = 99
	updatedDevice.FCntUp = 100
	updatedDevice.AppSKey, _ = protocol.AESKeyFromString("aaaa bbbb cccc dddd eeee ffff 0000 1111")
	updatedDevice.NwkSKey, _ = protocol.AESKeyFromString("1111 bbbb 2222 dddd eeee ffff 0000 1111")
	if err := storage.UpdateDevice(updatedDevice); err != nil {
		t.Fatal("Got error updating device: ", err)
	}
	tmp, _ := storage.GetDeviceByEUI(updatedDevice.DeviceEUI)
	if tmp.AppSKey.String() != updatedDevice.AppSKey.String() || tmp.NwkSKey.String() != updatedDevice.NwkSKey.String() || tmp.DevAddr.String() != updatedDevice.DevAddr.String() || tmp.RelaxedCounter != updatedDevice.RelaxedCounter || tmp.FCntDn != updatedDevice.FCntDn || tmp.FCntUp != updatedDevice.FCntUp {
		t.Fatalf("Device did not update correctly %v != %v", tmp, updatedDevice)
	}

	// Delete the devices, then delete application and network
	if err := storage.DeleteDevice(deviceA.DeviceEUI); err != nil {
		t.Fatalf("Got error deleting device: %v", err)
	}
	if err := storage.DeleteDevice(deviceB.DeviceEUI); err != nil {
		t.Fatalf("Got error deleting device: %v", err)
	}
	if err := storage.DeleteDevice(deviceC.DeviceEUI); err != nil {
		t.Fatalf("Got error deleting device: %v", err)
	}
	if err := storage.DeleteDevice(deviceD.DeviceEUI); err != nil {
		t.Fatalf("Got error deleting device: %v", err)
	}
	if err := storage.DeleteDevice(deviceA.DeviceEUI); err == nil {
		t.Fatal("Expected error deleting device twice")
	}
	if err := storage.DeleteApplication(app1.AppEUI); err != nil {
		t.Fatalf("Couldn't delete app 1: %v", err)
	}
	if err := storage.DeleteApplication(app2.AppEUI); err != nil {
		t.Fatalf("Couldn't delete app 2: %v", err)
	}
}
