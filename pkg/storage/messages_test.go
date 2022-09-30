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

	assert.NoError(storage.CreateUpstreamData(device.DeviceEUI, deviceData1), "Message 1 stored successfully")

	assert.NoError(storage.CreateUpstreamData(device.DeviceEUI, deviceData2), "Message 2 stored successfully")

	// Storing it a 2nd time won't work
	assert.Error(storage.CreateUpstreamData(device.DeviceEUI, deviceData1), "May only store message 1 once")

	assert.Error(storage.CreateUpstreamData(device.DeviceEUI, deviceData2), "May only store message 2 once")

	// Test retrieval
	data, err := storage.GetUpstreamDataByDeviceEUI(device.DeviceEUI, 2)
	assert.NoError(err, "No error when retrieving data")

	assert.Contains(data, deviceData1, "Message 1 returned")
	assert.Contains(data, deviceData2, "Message 2 returned")

	// Try retrieving from device with no data.
	data, err = storage.GetUpstreamDataByDeviceEUI(makeRandomEUI(), 2)
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
	if err := s.CreateDownstreamData(testDevice.DeviceEUI, downstreamMsg); err != nil {
		t.Fatal("Couldn't store downstream message: ", err)
	}

	newDownstreamMsg := model.NewDownstreamMessage(testDevice.DeviceEUI, 43)
	newDownstreamMsg.Ack = false
	newDownstreamMsg.Data = "aabbccddeeff"
	if err := s.CreateDownstreamData(testDevice.DeviceEUI, newDownstreamMsg); err == nil {
		t.Fatal("Shouldn't be able to store another downstream message")
	}

	if err := s.DeleteDownstreamData(testDevice.DeviceEUI); err != nil {
		t.Fatalf("Couldn't remove downstream message: %v", err)
	}

	if err := s.DeleteDownstreamData(testDevice.DeviceEUI); err != ErrNotFound {
		t.Fatalf("Should get ErrNotFound when removing message but got: %v", err)
	}

	if _, err := s.GetDownstreamData(testDevice.DeviceEUI); err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound but got %v", err)
	}

	if err := s.CreateDownstreamData(testDevice.DeviceEUI, newDownstreamMsg); err != nil {
		t.Fatalf("Should be able to store the new downstream message but got %v: ", err)
	}

	time2 := time.Now().Unix()
	if err := s.UpdateDownstreamData(testDevice.DeviceEUI, time2, 0); err != nil {
		t.Fatal("Should be able to update sent time but got error: ", err)
	}

	newDownstreamMsg.SentTime = time2
	stored, err := s.GetDownstreamData(testDevice.DeviceEUI)
	if err != nil {
		t.Fatal("Got error retrieving downstream message: ", err)
	}
	if stored != newDownstreamMsg {
		t.Fatalf("Sent time isn't updated properly. Got %+v but expected %+v", stored, newDownstreamMsg)
	}

	time3 := time.Now().Unix()
	if err := s.UpdateDownstreamData(testDevice.DeviceEUI, 0, time3); err != nil {
		t.Fatal("Got error updating downstream message: ", err)
	}

	stored, err = s.GetDownstreamData(testDevice.DeviceEUI)
	if err != nil {
		t.Fatal("Got error retrieving downstream message: ", err)
	}
	if stored.AckTime != time3 {
		t.Fatalf("Ack time isn't updated properly. Got %d but expected %d", stored.AckTime, time3)
	}

	if err := s.DeleteDownstreamData(testDevice.DeviceEUI); err != nil {
		t.Fatalf("Did not expect error when deleting downstream but got %v", err)
	}

	if err := s.UpdateDownstreamData(testDevice.DeviceEUI, 0, 0); err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound when updating nonexisting message but got %v", err)
	}

}
