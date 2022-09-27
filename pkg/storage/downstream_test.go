package storage

import (
	"testing"
	"time"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
)

func testDownstreamStorage(s *Storage, t *testing.T) {
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
