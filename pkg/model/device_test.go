package model

import (
	"testing"
	"time"
)

func TestDeviceStateConversion(t *testing.T) {
	states := []DeviceState{
		OverTheAirDevice,
		PersonalizedDevice,
		DisabledDevice,
	}
	for _, v := range states {
		if val, err := DeviceStateFromString(v.String()); val != v || err != nil {
			t.Errorf("Coudldn't convert %v to and from string (error is %v)", v, err)
		}
	}

	if _, err := DeviceStateFromString("unknown state"); err == nil {
		t.Error("Expected error when using unknown string format")
	}
}

func TestRXWindows(t *testing.T) {
	// These values are hard coded. The *real* test will use the device's settings
	device := Device{}
	if device.GetRX1Window() != (time.Second * 1) {
		t.Error("Someone must have fixed the GetRX1Window func but not the test")
	}
	if device.GetRX2Window() != (time.Second * 2) {
		t.Error("Someone must have fixed the GetRX2Window func but not the test")
	}
}

func TestDevNonce(t *testing.T) {
	d := Device{
		DevNonceHistory: []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}
	if !d.HasDevNonce(1) || !d.HasDevNonce(5) || !d.HasDevNonce(9) {
		t.Fatal("Expected 1, 5 and 9 to be in nonce history")
	}
	if d.HasDevNonce(0) || d.HasDevNonce(10) {
		t.Fatal("Didn't expect 0 or 10 to be in nonce history")
	}
}
