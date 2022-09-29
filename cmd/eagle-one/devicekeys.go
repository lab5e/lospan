package main

import (
	"encoding/hex"

	"github.com/lab5e/lospan/pkg/pb/lospan"
	"github.com/lab5e/lospan/pkg/protocol"
)

// DeviceKeys stores the keys, EUIs and DevAddr in Congress-friendly formats
type DeviceKeys struct {
	AppEUI  protocol.EUI
	DevEUI  protocol.EUI
	AppKey  protocol.AESKey
	AppSKey protocol.AESKey  // App session key
	NwkSKey protocol.AESKey  // Network session key
	DevAddr protocol.DevAddr // Device address
}

// NewDeviceKeys creates a new device key type from a Lassie Device
func NewDeviceKeys(device *lospan.Device) (DeviceKeys, error) {
	d := DeviceKeys{}
	var err error
	if d.AppEUI, err = protocol.EUIFromString(device.GetApplicationEui()); err != nil {
		return d, err
	}
	if d.DevEUI, err = protocol.EUIFromString(device.GetEui()); err != nil {
		return d, err
	}
	if d.AppKey, err = protocol.AESKeyFromString(hex.EncodeToString(device.AppKey)); err != nil {
		return d, err
	}
	d.DevAddr = protocol.DevAddrFromUint32(device.GetDevAddr())
	if d.NwkSKey, err = protocol.AESKeyFromString(hex.EncodeToString(device.NetworkSessionKey)); err != nil {
		return d, err
	}
	if d.AppSKey, err = protocol.AESKeyFromString(hex.EncodeToString(device.AppSessionKey)); err != nil {
		return d, err
	}
	return d, nil
}
