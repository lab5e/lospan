package processor

import (
	"time"

	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/storage"
)

// Decrypter decrypts the decoded LoRa packets into payloads. It will
// receive decoded packets on one channel and submit decrypted payloads on
// another channel.
type Decrypter struct {
	input     <-chan server.LoRaMessage
	macOutput chan server.LoRaMessage
	context   *server.Context
}

func (d *Decrypter) validFrameCounter(device *model.Device, decoded server.LoRaMessage) bool {
	// Ignore frame counter for JoinRequest messages since that will be reset
	// when the device have joined.
	if decoded.Payload.MHDR.MType == protocol.JoinRequest {
		return true
	}
	// Ensure frame counters are valid if it has strict checks
	if !device.RelaxedCounter {
		// Ignore frame counters that are less than the stored value. Bigger ones
		// means that we've lost one or more message from the device.
		if device.FCntUp > decoded.Payload.MACPayload.FHDR.FCnt {
			lg.Info("Frame counter check failed for device %s. Expected %d but got %d. Ignoring message.",
				device.DeviceEUI, device.FCntUp, decoded.Payload.MACPayload.FHDR.FCnt)
			return false
		}
	}

	// Issue debug warning if there's a mismatch between expected and actual frame counter
	// but process the message. This warning will be issued for all mismatchs.
	if device.FCntUp != decoded.Payload.MACPayload.FHDR.FCnt {
		lg.Debug("Frame counter will be adjusted. Expected %d but got %d for device with EUI %s",
			device.FCntUp, decoded.Payload.MACPayload.FHDR.FCnt, device.DeviceEUI)
	}
	return true
}

// processMessage forwards the message to the proper application
func (d *Decrypter) processMessage(device *model.Device, decoded server.LoRaMessage, matchingDevices int) {
	// Frame counters are tricky if there's more than one device since two (or more) devices
	// will send different frame counters. But this will be treated like any other message. With strict checks in place you *will* loose messages.

	// Frame counter checks does not apply for JoinRequest messages
	if !d.validFrameCounter(device, decoded) {
		return
	}

	if matchingDevices > 1 {
		// Set the key warning for the device if there's more than one device
		// that matches DevAddr/NwkSKey
		device.KeyWarning = true
	}

	// Update frame counter with the next expected message.
	if decoded.Payload.MACPayload.FHDR.FCnt >= device.FCntUp {
		device.FCntUp = decoded.Payload.MACPayload.FHDR.FCnt + 1
		if err := d.context.Storage.UpdateDeviceState(*device); err != nil {
			lg.Warning("Unable to update frame counters for device with EUI %s: %v", device.DeviceEUI, err)
		}
	}
	decoded.Payload.Decrypt(device.NwkSKey, device.AppSKey)

	deviceData := model.UpstreamMessage{
		DeviceEUI:  device.DeviceEUI,
		Timestamp:  decoded.FrameContext.GatewayContext.ReceivedAt.UnixNano(),
		Data:       decoded.Payload.MACPayload.FRMPayload,
		GatewayEUI: decoded.FrameContext.GatewayContext.Gateway.GatewayEUI,
		RSSI:       decoded.FrameContext.GatewayContext.Radio.RSSI,
		SNR:        decoded.FrameContext.GatewayContext.Radio.SNR,
		Frequency:  decoded.FrameContext.GatewayContext.Radio.Frequency,
		DataRate:   decoded.FrameContext.GatewayContext.Radio.DataRate,
		DevAddr:    device.DevAddr,
	}

	if err := d.context.Storage.CreateUpstreamData(device.DeviceEUI, device.AppEUI, deviceData); err != nil {
		lg.Warning("Unable to store device  with EUI: %s, error: %v", device.DeviceEUI, err)
		return
	}

	application, err := d.context.Storage.GetApplicationByEUI(device.AppEUI)
	if err != nil {
		lg.Warning("Unable to retrieve application with EUI %s: %v", device.AppEUI, err)
		return
	}

	decoded.FrameContext.Application = application
	decoded.FrameContext.Device = *device

	if decoded.Payload.MHDR.MType == protocol.ConfirmedDataUp {
		d.context.FrameOutput.SetMessageAckFlag(device.DeviceEUI, true)
	}

	msg, err := d.context.Storage.GetDownstreamData(device.DeviceEUI)
	if err == nil {
		// Update state of message -- note that this could cause some inconsistent
		// behaviour if you create a message (with ack) and replaces it with a new
		// message after it has been sent but before it has been acked:
		//
		// t=0  : Create new downstream message with ack flag set
		// t=1  : Device sends upstream message, server responds with downstream message
		// t=1.1: Downstream message is removed and replaced with a new message
		// t=1.2: New downstream message with ack flag set
		// t=2  : Device acknowledges message from t=1
		// Since the message from t=1.2 isn't sent yet the ack will be ignored.
		if decoded.Payload.MACPayload.FHDR.FCtrl.ACK && msg.State() == model.SentState && msg.Ack {
			msg.AckTime = time.Now().Unix()
			if err := d.context.Storage.UpdateDownstreamData(device.DeviceEUI, msg.SentTime, msg.AckTime); err != nil {
				lg.Warning("Unable to update downstream message: %v", err)
			}
		}
		if !msg.IsComplete() {
			lg.Debug("Setting downstream message payload (%v) for device %s", msg.Payload(), device.DeviceEUI)
			d.context.FrameOutput.SetPayload(device.DeviceEUI, msg.Payload(), msg.Port, msg.Ack)
		}
	}

	if err != nil && err != storage.ErrNotFound {
		lg.Warning("Unable to retrieve downstream message: %v", err)
	}

	d.macOutput <- decoded

	d.context.AppRouter.Publish(application.AppEUI, &server.PayloadMessage{
		Payload:      decoded.Payload.MACPayload.FRMPayload,
		Device:       *device,
		Application:  application,
		FrameContext: decoded.FrameContext,
	})

}

func (d *Decrypter) verifyAndDecryptMessage(decoded server.LoRaMessage) {
	lg.Debug("Verifying message from device with DevAddr %s", decoded.Payload.MACPayload.FHDR.DevAddr)
	devices, err := d.context.Storage.GetDeviceByDevAddr(decoded.Payload.MACPayload.FHDR.DevAddr)
	if err != nil {
		lg.Warning("Unable to retrieve device from storage. Network ID: %x, Network address: %x. Error: %v",
			decoded.Payload.MACPayload.FHDR.DevAddr.NwkID,
			decoded.Payload.MACPayload.FHDR.DevAddr.NwkAddr, err)
		return
	}

	rawMessage := decoded.FrameContext.GatewayContext.RawMessage
	if len(rawMessage) < protocol.MinimumMessageSize {
		// Ignore message. There's no MIC here.
		return
	}

	// Find all devices with matching devAddr and correct NwkSKey (ie the MIC is
	// valid) and forward it to the appropriate application. Issue warnings
	// wrt key for devices if there's more than one device with the same key.
	var matchingDevices []model.Device
	checked := 0
	for _, dev := range devices {
		checked++
		lg.Debug("Testing MIC for device %s", dev.DeviceEUI)
		mic, err := decoded.Payload.CalculateMIC(dev.NwkSKey, rawMessage[0:len(rawMessage)-4])
		if err != nil {
			lg.Info("Unable to calculate MIC for payload: %v (payload=%v) ", err, decoded.Payload)
			continue
		}
		if mic == decoded.Payload.MIC {
			matchingDevices = append(matchingDevices, dev)
		}
	}
	if len(matchingDevices) == 0 && checked > 0 {
		lg.Info("MIC validation failed for device with DevAddr: %s", decoded.Payload.MACPayload.FHDR.DevAddr)
		return
	}

	// We now have a list of devices
	for _, dev := range matchingDevices {
		d.processMessage(&dev, decoded, len(matchingDevices))
	}
}

// Start launches the decrypter. It will loop forever until the input
// channel closes. The output channel will be closed upon return.
// BUG(stalehd): Doesn't do what it says -- decrypt
func (d *Decrypter) Start() {
	if d.context.Storage == nil {
		lg.Error("No storage. Unable to proceed.")
		return
	}
	for m := range d.input {
		go func(decoded server.LoRaMessage) {
			if decoded.FrameContext.GatewayContext.RawMessage == nil {
				lg.Error("Missing raw message representation. Unable to proceed.")
				return
			}
			if decoded.Payload.MHDR.MType == protocol.JoinRequest {
				go func() {
					d.processJoinRequest(decoded)
				}()
				return
			}

			d.verifyAndDecryptMessage(decoded)
		}(m)
	}

	lg.Debug("Input channel for Decrypter closed. Terminating")
	close(d.macOutput)
}

// Output returns the MAC notification output from the decrypter. This channel
// will receive a message every time a message is successfully verified and
// decrypted.
func (d *Decrypter) Output() <-chan server.LoRaMessage {
	return d.macOutput
}

// NewDecrypter creates a new decrypter instance.
func NewDecrypter(context *server.Context, input <-chan server.LoRaMessage) *Decrypter {
	return &Decrypter{
		input:     input,
		macOutput: make(chan server.LoRaMessage),
		context:   context,
	}
}
