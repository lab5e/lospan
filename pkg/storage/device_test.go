package storage

import (
	"database/sql"
	"testing"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/stretchr/testify/require"
)

func TestDeviceStorage(t *testing.T) {
	assert := require.New(t)

	connectionString := ":memory:"
	db, err := sql.Open(driverName, connectionString)
	assert.Nil(err, "No error opening database")
	assert.Nil(createSchema(db), "Error creating db in %s", connectionString)
	db.Close()

	storage, err := CreateStorage(connectionString)
	assert.Nil(err, "Did not expect error: %v", err)
	defer storage.Close()

	app1 := model.Application{
		AppEUI: makeRandomEUI(),
	}
	assert.NoError(storage.CreateApplication(app1), "Error adding application 1")

	app2 := model.Application{
		AppEUI: makeRandomEUI(),
	}
	assert.NoError(storage.CreateApplication(app2), "Got error adding application 2")

	deviceA := model.Device{
		DeviceEUI: makeRandomEUI(),
		AppEUI:    app1.AppEUI,
		DevAddr: protocol.DevAddr{
			NwkID:   1,
			NwkAddr: 0x000001,
		},
		FCntUp: 1,
	}
	assert.NoError(storage.CreateDevice(deviceA, app1.AppEUI), "Error creating device A")

	deviceB := model.Device{
		DeviceEUI: makeRandomEUI(),
		AppEUI:    app1.AppEUI,
		DevAddr: protocol.DevAddr{
			NwkID:   1,
			NwkAddr: 0x000002,
		},
		FCntUp: 2,
	}
	assert.NoError(storage.CreateDevice(deviceB, app1.AppEUI), "Error creating device B")

	deviceC := model.Device{
		DeviceEUI: makeRandomEUI(),
		AppEUI:    app2.AppEUI,
		DevAddr: protocol.DevAddr{
			NwkID:   1,
			NwkAddr: 0x000003,
		},
		FCntUp: 3,
	}
	assert.NoError(storage.CreateDevice(deviceC, app2.AppEUI), "Error creating device C")

	deviceD := model.Device{
		DeviceEUI: makeRandomEUI(),
		AppEUI:    app2.AppEUI,
		DevAddr: protocol.DevAddr{
			NwkID:   1,
			NwkAddr: 0x000004,
		},
		FCntUp: 4,
	}
	assert.NoError(storage.CreateDevice(deviceD, app2.AppEUI), "Error creating device D")

	// Retrieve one of the stored devices via DevAddr (assume the others work)
	devices, err := storage.GetDeviceByDevAddr(deviceC.DevAddr)
	assert.NoError(err, "Error retrieving by device address")
	assert.Contains(devices, deviceC, "Device C is not in returned list")

	// ...and do the same for a device keyed on EUI
	device, err := storage.GetDeviceByEUI(deviceB.DeviceEUI)
	assert.NoError(err, "Error retrieving device B")
	assert.Equal(deviceB, device, "Device B is not returned correctly")

	devices1, err := storage.GetDevicesByApplicationEUI(app1.AppEUI)
	assert.NoError(err, "Error retrieving device list for app 1")
	assert.Contains(devices1, deviceA, "Device A is not in list")
	assert.Contains(devices1, deviceB, "Device B is not in list")
	assert.Len(devices1, 2)

	devices2, err := storage.GetDevicesByApplicationEUI(app2.AppEUI)
	assert.NoError(err, "Error retrieving device list for app 2")
	assert.Contains(devices2, deviceC, "Device C is not in list")
	assert.Contains(devices2, deviceD, "Device D is not in list")
	assert.Len(devices2, 2)

	_, err = storage.GetDeviceByEUI(protocol.EUIFromInt64(0))
	assert.Error(err, "Expected error for unknow EUI")

	// Try adding the same device twice
	assert.Error(storage.CreateDevice(deviceA, app1.AppEUI), "Expected error on duplicate device")

	// Store device nonce on device, ensure it is stored
	assert.NoError(storage.AddDevNonce(deviceA, 12), "No error storing nonce")
	assert.NoError(storage.AddDevNonce(deviceA, 24), "No error storing nonce")
	assert.NoError(storage.AddDevNonce(deviceA, 48), "No error storing nonce")

	device, err = storage.GetDeviceByEUI(deviceA.DeviceEUI)
	assert.NoError(err, "No error retrieving device")

	assert.True(device.HasDevNonce(12), "Should have nonce 12")
	assert.True(device.HasDevNonce(24), "Should have nonce 24")
	assert.True(device.HasDevNonce(48), "Should have nonce 48")
	assert.False(device.HasDevNonce(96), "Should not have nonce 96")

	deviceC.NwkSKey = makeRandomKey()
	deviceC.AppSKey = makeRandomKey()
	assert.NoError(storage.UpdateDevice(deviceC), "Update for device C should work")

	device, err = storage.GetDeviceByEUI(deviceC.DeviceEUI)
	assert.NoError(err, "Should be able to read device")
	assert.Equal(deviceC.AppSKey, device.AppSKey)
	assert.Equal(deviceC.NwkSKey, device.NwkSKey)

	deviceD.FCntDn = 1001
	deviceD.FCntUp = 2002
	deviceD.KeyWarning = true
	assert.NoError(storage.UpdateDeviceState(deviceD), "State update for device D should work")

	updatedDevice, err := storage.GetDeviceByEUI(deviceD.DeviceEUI)
	assert.NoError(err, "Retrieve device D should work")
	assert.Equal(deviceD.FCntDn, updatedDevice.FCntDn)
	assert.Equal(deviceD.FCntUp, updatedDevice.FCntUp)
	assert.True(updatedDevice.KeyWarning)

	updatedDevice.DevAddr = protocol.DevAddrFromUint32(0x01020304)
	updatedDevice.RelaxedCounter = true
	updatedDevice.FCntDn = 99
	updatedDevice.FCntUp = 100
	updatedDevice.AppSKey, _ = protocol.AESKeyFromString("aaaa bbbb cccc dddd eeee ffff 0000 1111")
	updatedDevice.NwkSKey, _ = protocol.AESKeyFromString("1111 bbbb 2222 dddd eeee ffff 0000 1111")

	assert.NoError(storage.UpdateDevice(updatedDevice), "Expect no error when updating device with keys and counters")

	newDevice, err := storage.GetDeviceByEUI(updatedDevice.DeviceEUI)
	assert.NoError(err)
	assert.Equal(updatedDevice, newDevice)

	// Delete the devices, then delete application and network
	assert.NoError(storage.DeleteDevice(deviceA.DeviceEUI))
	assert.NoError(storage.DeleteDevice(deviceB.DeviceEUI))
	assert.NoError(storage.DeleteDevice(deviceC.DeviceEUI))
	assert.NoError(storage.DeleteDevice(deviceD.DeviceEUI))

	assert.Error(storage.DeleteDevice(deviceA.DeviceEUI), "Should not be able to delete device twice")

	assert.NoError(storage.DeleteApplication(app1.AppEUI))
	assert.NoError(storage.DeleteApplication(app2.AppEUI))
}
