package processor

//
//Copyright 2018 Telenor Digital AS
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
import (
	"github.com/lab5e/lospan/pkg/frequency"
	"github.com/lab5e/lospan/pkg/lg"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/server"
)

// Process the join request. Returns false if it failed.
func (d *Decrypter) processJoinRequest(decoded server.LoRaMessage) bool {
	joinRequest := &decoded.Payload.JoinRequestPayload

	device, err := d.context.Storage.GetDeviceByEUI(joinRequest.DevEUI)
	if err != nil {
		lg.Info("Unknown device attempting JoinRequest: %s", joinRequest.DevEUI)
		return false
	}

	if device.AppEUI != joinRequest.AppEUI {
		lg.Warning("Mismatch between stored device's AppEUI and the AppEUI sent in the JoinRequest message. Stored AppEUI = %s, JoinRequest AppEUI = %s", device.AppEUI, joinRequest.AppEUI)
		return false
	}

	// Check if DevNonce have been used by the device in an earlier request.
	// If so the request should be ignored. [6.2.4].
	if !d.context.Config.DisableNonceCheck && device.HasDevNonce(joinRequest.DevNonce) {
		lg.Warning("Device %s has already used nonce 0x%04x. Ignoring it.",
			joinRequest.DevEUI, joinRequest.DevNonce)
		return false
	}
	if d.context.Config.DisableGatewayChecks && device.HasDevNonce(joinRequest.DevNonce) {
		lg.Warning("Device %s is re-using a nonce (0x%04x) but nonce checks are disabled", joinRequest.DevEUI, joinRequest.DevNonce)
	}

	// Retrieve the application
	app, err := d.context.Storage.GetApplicationByEUI(joinRequest.AppEUI)
	if err != nil {
		lg.Warning("Unable to retrieve application with EUI %s. Ignoring JoinRequest from device with EUI %s",
			joinRequest.AppEUI, joinRequest.DevEUI)
		return false
	}

	// Invariant: Application is OK, device is OK, network is OK. Generate response
	decoded.FrameContext.Application = app
	decoded.FrameContext.Device = device

	// DevAddr is already assigned to the device. It is a function of the EUI.

	// Update the device with new keys and DevNonce
	if !d.context.Config.DisableNonceCheck {
		if err := d.context.Storage.AddDevNonce(device, joinRequest.DevNonce); err != nil {
			lg.Warning("Unable to update DevNonce on device with EUI: %s: %v",
				device.DeviceEUI, err)
		}
	}

	// Generate app nonce, generate keys, store keys
	appNonce, err := app.GenerateAppNonce()
	if err != nil {
		lg.Warning("Unable to generate app nonce: %v (devEUI: %s, appEUI: %s). Ignoring JoinRequest",
			err, joinRequest.DevEUI, joinRequest.AppEUI)
		return false
	}
	nwkSKey, err := protocol.NwkSKeyFromNonces(device.AppKey, appNonce, uint32(d.context.Config.NetworkID), joinRequest.DevNonce)
	if err != nil {
		lg.Error("Unable to generate NwkSKey for device with EUI %s: %v", device.DeviceEUI, err)
		return false
	}
	appSKey, err := protocol.AppSKeyFromNonces(device.AppKey, appNonce, uint32(d.context.Config.NetworkID), joinRequest.DevNonce)
	if err != nil {
		lg.Error("Unable to generate AppSKey for device with EUI %s: %v", device.DeviceEUI, err)
		return false
	}
	device.NwkSKey = nwkSKey
	device.AppSKey = appSKey
	device.FCntDn = 0
	device.FCntUp = 0
	if device.DevAddr.ToUint32() == 0 {
		// Set device address if it isn't set
		device.DevAddr = protocol.NewDevAddr()
	}
	if err := d.context.Storage.UpdateDevice(device); err != nil {
		lg.Error("Unable to update device with EUI %s: %v", device.DeviceEUI, err)
		return false
	}

	// Invariant. Everything is OK - make a JoinAccept response and schedule
	// the output.
	joinAccept := protocol.JoinAcceptPayload{
		AppNonce:   appNonce,
		NetID:      uint32(d.context.Config.NetworkID),
		DevAddr:    device.DevAddr,
		DLSettings: frequency.GetDLSettingsOTAA(),
		RxDelay:    frequency.GetRxDelayOTAA(),
		CFList:     frequency.GetCFListOTAA(),
	}

	d.context.FrameOutput.SetJoinAcceptPayload(device.DeviceEUI, joinAccept)

	lg.Debug("JoinAccept sent to %s. DevAddr=%s", device.DeviceEUI, joinAccept.DevAddr)

	// The incoming message doesn't have a DevAddr set but schedule an empty
	// message for it. TODO (stalehd): this is butt ugly. Needs redesign.
	decoded.Payload.MACPayload.FHDR.DevAddr = joinAccept.DevAddr

	d.macOutput <- decoded
	return true
}
