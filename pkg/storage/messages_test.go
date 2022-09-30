package storage

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/stretchr/testify/require"
)

func TestUpstreamStorage(t *testing.T) {
	assert := require.New(t)

	connectionString := ":memory:"
	db, err := sql.Open(driverName, connectionString)
	assert.Nil(err, "No error opening database")
	assert.Nil(createSchema(db), "Error creating db in %s", connectionString)
	db.Close()

	storage, err := CreateStorage(connectionString)
	assert.Nil(err, "Did not expect error: %v", err)
	defer storage.Close()

	app := model.Application{
		AppEUI: makeRandomEUI(),
	}

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

	deviceData1 := model.UpstreamMessage{Timestamp: 1, Data: data1, DeviceEUI: device.DeviceEUI, Frequency: 1.0}
	deviceData2 := model.UpstreamMessage{Timestamp: 2, Data: data2, DeviceEUI: device.DeviceEUI, Frequency: 2.0}

	assert.NoError(storage.CreateUpstreamMessage(device.DeviceEUI, deviceData1), "Message 1 stored successfully")

	assert.NoError(storage.CreateUpstreamMessage(device.DeviceEUI, deviceData2), "Message 2 stored successfully")

	// Storing it a 2nd time won't work
	assert.Error(storage.CreateUpstreamMessage(device.DeviceEUI, deviceData1), "May only store message 1 once")

	assert.Error(storage.CreateUpstreamMessage(device.DeviceEUI, deviceData2), "May only store message 2 once")

	// Test retrieval
	data, err := storage.ListUpstreamMessages(device.DeviceEUI, 2)
	assert.NoError(err, "No error when retrieving data")

	assert.Contains(data, deviceData1, "Message 1 returned")
	assert.Contains(data, deviceData2, "Message 2 returned")

	// Try retrieving from device with no data.
	data, err = storage.ListUpstreamMessages(makeRandomEUI(), 2)
	assert.NoError(err, "No device => no error (and no data)")
	assert.Len(data, 0)

}

func TestDownstreamStorage(t *testing.T) {
	assert := require.New(t)

	connectionString := ":memory:"
	db, err := sql.Open(driverName, connectionString)
	assert.Nil(err, "No error opening database")
	assert.Nil(createSchema(db), "Error creating db in %s", connectionString)
	db.Close()

	s, err := CreateStorage(connectionString)
	assert.Nil(err, "Did not expect error: %v", err)
	defer s.Close()

	application := model.NewApplication()
	application.AppEUI = makeRandomEUI()
	s.CreateApplication(application)

	testDevice := model.NewDevice()
	testDevice.AppEUI = application.AppEUI
	testDevice.DeviceEUI = makeRandomEUI()
	testDevice.AppSKey = makeRandomKey()
	testDevice.DevAddr = protocol.DevAddrFromUint32(0x01020304)
	testDevice.NwkSKey = makeRandomKey()
	s.CreateDevice(testDevice, application.AppEUI)

	downstreamMsg := model.NewDownstreamMessage(testDevice.DeviceEUI, 42)
	downstreamMsg.Ack = false
	downstreamMsg.Data = "aabbccddeeff"
	assert.NoError(s.CreateDownstreamMessage(testDevice.DeviceEUI, downstreamMsg), "Should be able to store downstream message")

	newDownstreamMsg := model.NewDownstreamMessage(testDevice.DeviceEUI, 43)
	newDownstreamMsg.Ack = false
	newDownstreamMsg.Data = "aabbccddeeff"

	assert.Error(s.CreateDownstreamMessage(testDevice.DeviceEUI, newDownstreamMsg),
		"Shouldn't be able to store another downstream message")

	assert.NoError(s.DeleteDownstreamMessage(testDevice.DeviceEUI))

	assert.Equal(ErrNotFound, s.DeleteDownstreamMessage(testDevice.DeviceEUI))

	_, err = s.GetNextDownstreamMessage(testDevice.DeviceEUI)
	assert.Equal(ErrNotFound, err)

	assert.NoError(s.CreateDownstreamMessage(testDevice.DeviceEUI, newDownstreamMsg))

	time2 := time.Now().Unix()
	assert.NoError(s.UpdateDownstreamMessage(testDevice.DeviceEUI, time2, 0))

	newDownstreamMsg.SentTime = time2
	stored, err := s.GetNextDownstreamMessage(testDevice.DeviceEUI)
	assert.NoError(err)

	assert.Equal(newDownstreamMsg, stored, "Sent time isn't updated properly")

	time3 := time.Now().Unix()
	assert.NoError(s.UpdateDownstreamMessage(testDevice.DeviceEUI, time2, time3))

	stored, err = s.GetNextDownstreamMessage(testDevice.DeviceEUI)
	assert.NoError(err)

	assert.Equal(time3, stored.AckTime)

	assert.NoError(s.CreateDownstreamMessage(testDevice.DeviceEUI, model.DownstreamMessage{
		DeviceEUI: testDevice.DeviceEUI,
		Data:      "0102030405",
		Port:      2,
		Ack:       false,
	}))
	list, err := s.ListDownstreamMessages(testDevice.DeviceEUI)
	assert.NoError(err)
	assert.Len(list, 2)

	assert.NoError(s.DeleteDownstreamMessage(testDevice.DeviceEUI))

	assert.Equal(ErrNotFound, s.UpdateDownstreamMessage(testDevice.DeviceEUI, 0, 0))

}
