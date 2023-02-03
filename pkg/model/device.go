package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/lab5e/lospan/pkg/lg"
	"github.com/lab5e/lospan/pkg/protocol"
)

// DeviceState represents the device state
type DeviceState uint8

// Types of devices. A device can either be OTAA, ABP or disabled.
const (
	OverTheAirDevice   DeviceState = 1
	PersonalizedDevice DeviceState = 8 // Note: 8 is for backwards compatibility
	DisabledDevice     DeviceState = 0
)

// String converts the device state into a human-readable string representation.
func (d DeviceState) String() string {
	switch d {
	case OverTheAirDevice:
		return "OTAA"
	case PersonalizedDevice:
		return "ABP"
	case DisabledDevice:
		return "Disabled"
	default:
		lg.Warning("Unknown device state: %d", d)
		return "Disabled"
	}
}

// DeviceStateFromString converts a string representation of DeviceState into
// a DeviceState value. Unknown strings returns the DisabledDevice state.
// Conversion is not case sensitive. White space is trimmed.
func DeviceStateFromString(str string) (DeviceState, error) {
	switch strings.TrimSpace(strings.ToUpper(str)) {
	case "OTAA":
		return OverTheAirDevice, nil
	case "ABP":
		return PersonalizedDevice, nil
	case "DISABLED":
		return DisabledDevice, nil
	default:
		return DisabledDevice, fmt.Errorf("unknown device state: %s", str)
	}
}

// Device represents a device. Devices are associated with one and only one Application
type Device struct {
	DeviceEUI       protocol.EUI     // EUI for device
	DevAddr         protocol.DevAddr // Device address
	AppKey          protocol.AESKey  // AES key for application
	AppSKey         protocol.AESKey  // Application session key
	NwkSKey         protocol.AESKey  // Network session key
	AppEUI          protocol.EUI     // The application associated with the device. Set by storage backend
	State           DeviceState      // Current state of the device
	FCntUp          uint16           // Frame counter up (from device)
	FCntDn          uint16           // Frame counter down (to device)
	RelaxedCounter  bool             // Relaxed frame count checks
	DevNonceHistory []uint16         // Log of DevNonces sent from the device
	KeyWarning      bool             // Duplicate key warning flag
	Tag             string           // Tag data (for external refs)
}

// NewDevice creates a new device
func NewDevice() Device {
	return Device{}
}

// GetRX1Window returns the 1st receive window for the device
// BUG(stlaehd): Returns constant. Should be set based on device settings and frequency plan.
func (d *Device) GetRX1Window() time.Duration {
	return time.Second * 1
}

// GetRX2Window returns the 2nd receive window for the device
// BUG(stalehd): Returns a constant. Should be set based on frequency plan (EU, US, CN)
func (d *Device) GetRX2Window() time.Duration {
	return time.Second * 2
}

// HasDevNonce returns true if the specified nonce exists in the nonce history
func (d *Device) HasDevNonce(devNonce uint16) bool {
	for _, v := range d.DevNonceHistory {
		if v == devNonce {
			return true
		}
	}
	return false
}
