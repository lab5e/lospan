package storage

import (
	"testing"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/stretchr/testify/require"
)

func testDataStorage(storage *Storage, t *testing.T) {

	app := model.Application{
		AppEUI: makeRandomEUI(),
	}

	assert := require.New(t)

	assert.NoError(storage.CreateApplication(app))

	device := model.Device{
		DeviceEUI: makeRandomEUI(),
		AppEUI:    app.AppEUI,
		DevAddr: protocol.DevAddr{
			NwkID:   1,
			NwkAddr: 0x400004,
		},
		FCntUp: 4,
	}

	assert.NoError(storage.CreateDevice(device, app.AppEUI))

	data1 := makeRandomData()
	data2 := makeRandomData()

	deviceData1 := model.DeviceData{Timestamp: 1, Data: data1, DeviceEUI: device.DeviceEUI, AppEUI: device.AppEUI, Frequency: 1.0}
	deviceData2 := model.DeviceData{Timestamp: 2, Data: data2, DeviceEUI: device.DeviceEUI, AppEUI: device.AppEUI, Frequency: 2.0}

	assert.NoError(storage.CreateUpstreamData(device.DeviceEUI, device.AppEUI, deviceData1), "Message 1 stored successfully")

	assert.NoError(storage.CreateUpstreamData(device.DeviceEUI, device.AppEUI, deviceData2), "Message 2 stored successfully")

	// Storing it a 2nd time won't work
	assert.Error(storage.CreateUpstreamData(device.DeviceEUI, device.AppEUI, deviceData1), "May only store message 1 once")

	assert.Error(storage.CreateUpstreamData(device.DeviceEUI, device.AppEUI, deviceData2), "May only store message 2 once")

	// Test retrieval
	data, err := storage.GetUpstreamDataByDeviceEUI(device.DeviceEUI, 2)
	assert.NoError(err, "No error when retrieving data")

	assert.Contains(data, deviceData1, "Message 1 returned")
	assert.Contains(data, deviceData2, "Message 2 returned")

	// Try retrieving from device with no data.
	data, err = storage.GetUpstreamDataByDeviceEUI(makeRandomEUI(), 2)
	assert.NoError(err, "No device => no error (and no data)")
	assert.Len(data, 0)

	// Read application device data. Should be the same as the device data.
	appData, err := storage.GetDownstreamDataByApplicationEUI(app.AppEUI, 10)
	assert.NoError(err)
	assert.Len(appData, 2)
}
