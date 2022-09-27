package storage

import (
	"testing"
	"time"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
)

func testDataStorage(
	storage *Storage,
	t *testing.T) {

	app := model.Application{
		AppEUI: makeRandomEUI(),
	}
	if err := storage.CreateApplication(app); err != nil {
		t.Error("Got error adding application: ", err)
	}

	device := model.Device{
		DeviceEUI: makeRandomEUI(),
		AppEUI:    app.AppEUI,
		DevAddr: protocol.DevAddr{
			NwkID:   1,
			NwkAddr: 0x400004,
		},
		FCntUp: 4,
	}

	err := storage.CreateDevice(device, app.AppEUI)
	if err != nil {
		t.Error("Error putting device: ", err)
	}

	data1 := makeRandomData()
	data2 := makeRandomData()

	deviceData1 := model.DeviceData{Timestamp: 1, Data: data1, DeviceEUI: device.DeviceEUI, Frequency: 1.0}
	deviceData2 := model.DeviceData{Timestamp: 2, Data: data2, DeviceEUI: device.DeviceEUI, Frequency: 2.0}

	if err = storage.CreateUpstreamData(device.DeviceEUI, deviceData1); err != nil {
		t.Error("Could not store data: ", err)
	}

	if err = storage.CreateUpstreamData(device.DeviceEUI, deviceData2); err != nil {
		t.Error("Could not store 2nd data: ", err)
	}

	// Storing it a 2nd time won't work
	if err = storage.CreateUpstreamData(device.DeviceEUI, deviceData1); err == nil {
		t.Error("Shouldn't be able to store data twice (data#1)")
	}

	if err = storage.CreateUpstreamData(device.DeviceEUI, deviceData2); err == nil {
		t.Error("Shouldn't be able to store data twice (data#2)")
	}

	// Test retrieval
	dataChan, err := storage.GetUpstreamDataByDeviceEUI(device.DeviceEUI, 2)
	if err != nil {
		t.Error("Did not expect error when retrieving data")
	}

	var firstData, secondData model.DeviceData
	// Read from channels, time out if there's no data
	timestamps := int64(0)
	select {
	case firstData = <-dataChan:
		timestamps += firstData.Timestamp
	case <-time.After(time.Second * 2):
		t.Error("Timed out waiting for data # 1")
	}

	select {
	case secondData = <-dataChan:
		timestamps += secondData.Timestamp
	case <-time.After(time.Second * 2):
		t.Error("Timed out waiting for data # 2")
	}

	if timestamps != int64(3) {
		t.Error("Did not get the correct data pieces")
	}

	// Try retrieving from device with no data.
	var dataChannel chan model.DeviceData
	if _, err = storage.GetUpstreamDataByDeviceEUI(makeRandomEUI(), 2); err != nil {
		t.Error("Did not expect error when retrieving from non-existing device")
	}
	select {
	case <-dataChannel:
		t.Fatal("Did not expect any data on channel")
	case <-time.After(100 * time.Millisecond):
		// This is OK
	}

	// Read application device data. Should be the same as the device data.
	appChan, err := storage.GetDownstreamByApplicationEUI(app.AppEUI, 10)
	if err != nil {
		t.Fatal("Error retrieving from application: ", err)
	}
	count := 0
	for data := range appChan {
		if data.Equals(firstData) {
			count++
		}
		if data.Equals(secondData) {
			count++
		}
	}
	if count != 2 {
		t.Fatal("Missing data on application channel. Expected 2 got ", count)
	}
}
