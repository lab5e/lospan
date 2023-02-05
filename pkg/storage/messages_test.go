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
	downstreamMsg.CreatedTime = time.Now().UnixNano()
	assert.NoError(s.CreateDownstreamMessage(testDevice.DeviceEUI, downstreamMsg), "Should be able to store downstream message")
	assert.NoError(s.DeleteDownstreamMessage(testDevice.DeviceEUI, downstreamMsg.CreatedTime))

	newDownstreamMsg := model.NewDownstreamMessage(testDevice.DeviceEUI, 43)
	newDownstreamMsg.Ack = false
	newDownstreamMsg.Data = "aabbccddeeff"
	newDownstreamMsg.FCntUp = 99
	created := time.Now().UnixNano()
	newDownstreamMsg.CreatedTime = created
	assert.NoError(s.CreateDownstreamMessage(testDevice.DeviceEUI, newDownstreamMsg),
		"Shouldn't be able to store another downstream message")

	assert.NoError(s.DeleteDownstreamMessage(testDevice.DeviceEUI, created))

	assert.Equal(ErrNotFound, s.DeleteDownstreamMessage(testDevice.DeviceEUI, created))

	empty, err := s.GetNextUnsentMessage(testDevice.DeviceEUI)
	assert.Equal(int64(0), empty.CreatedTime)
	assert.Equal(ErrNotFound, err)

	assert.NoError(s.CreateDownstreamMessage(testDevice.DeviceEUI, newDownstreamMsg))

	time2 := time.Now().Unix()
	assert.NoError(s.SetMessageSentTime(testDevice.DeviceEUI, created, time2, newDownstreamMsg.FCntUp))

	newDownstreamMsg.SentTime = time2
	_, err = s.GetNextUnsentMessage(testDevice.DeviceEUI)
	assert.Equal(ErrNotFound, err)

	confirmableMessage := model.NewDownstreamMessage(testDevice.DeviceEUI, 99)
	confirmableMessage.CreatedTime = time.Now().UnixNano()
	confirmableMessage.Data = "this is the data"
	confirmableMessage.Ack = true
	assert.NoError(s.CreateDownstreamMessage(confirmableMessage.DeviceEUI, confirmableMessage))
	_, err = s.GetNextUnsentMessage(testDevice.DeviceEUI)
	assert.NoError(err)

	assert.NoError(s.SetMessageSentTime(confirmableMessage.DeviceEUI, confirmableMessage.CreatedTime, time.Now().UnixNano(), 101))

	_, err = s.GetNextUnsentMessage(testDevice.DeviceEUI)
	assert.Equal(ErrNotFound, err)

	// Invalid frame counter
	assert.Error(s.UpdateMessageAckTime(downstreamMsg.DeviceEUI, 199, time.Now().UnixNano()))

	// ok - got frame counter
	assert.NoError(s.UpdateMessageAckTime(downstreamMsg.DeviceEUI, 101, time.Now().UnixNano()))
	// can't ack twice
	assert.Error(s.UpdateMessageAckTime(downstreamMsg.DeviceEUI, 101, time.Now().UnixNano()))

}
