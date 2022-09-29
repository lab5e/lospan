package restapi

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/storage"
)

// Get and check EUIs for network, app, device. Returns false if one of the
// EUIs are malformed
func (s *Server) getDevice(w http.ResponseWriter, r *http.Request) (
	protocol.EUI, *model.Device) {

	app := s.getApplication(w, r)
	if app == nil {
		return protocol.EUI{}, nil
	}
	appEUI := app.AppEUI

	deviceEUI, err := euiFromPathParameter(r, "deui")
	if err != nil {
		http.Error(w, "Malformed Device EUI", http.StatusBadRequest)
		return appEUI, nil
	}

	device, err := s.context.Storage.GetDeviceByEUI(deviceEUI)
	if err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return appEUI, nil
	}

	if device.AppEUI != appEUI {
		http.Error(w, "Application not found", http.StatusNotFound)
		return appEUI, nil
	}

	return appEUI, &device
}

// Remove downstream message if exists and completed (sent and/or acked).
// Returns false if there's an error. The existing message is left as is if
// it should be kept (ie it isn't sent yet or not acked yet). This might be
// a bit counter-intuitive but it will fail on the PutDownstream call later on
// this elminates a few obvious concurrency issues but not all.
func (s *Server) removeDownstreamIfComplete(w http.ResponseWriter, deviceEUI protocol.EUI) bool {
	existingMessage, err := s.context.Storage.GetDownstreamData(deviceEUI)
	if err == storage.ErrNotFound {
		return true
	}
	if err != nil {
		http.Error(w, "unable to verify if there's a scheduled message", http.StatusInternalServerError)
		return false
	}

	if !existingMessage.IsComplete() {
		http.Error(w, "a message is already scheduled for output", http.StatusConflict)
		return false
	}

	if err := s.context.Storage.DeleteDownstreamData(deviceEUI); err != nil {
		http.Error(w, "unable to remove scheduled message", http.StatusInternalServerError)
		return false
	}

	return true
}

func (s *Server) createDownstream(device *model.Device, w http.ResponseWriter, r *http.Request) {
	//
	// Read body, decode message
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		lg.Warning("Unable to read request body for device %s: %v", device.DeviceEUI, err)
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	outMessage := make(map[string]interface{})
	if err = json.Unmarshal(buf, &outMessage); err != nil {
		lg.Info("Unable to marshal JSON for message to device with EUI %s: %v", device.DeviceEUI, err)
		http.Error(w, "Can't grok JSON", http.StatusBadRequest)
		return
	}
	port, ok := outMessage["port"].(float64)
	if !ok {
		http.Error(w, "port must be set for messages", http.StatusBadRequest)
		return
	}
	if port < 1 || port > 223 {
		http.Error(w, "port must be between 1 and 223", http.StatusBadRequest)
		return
	}

	data, ok := outMessage["data"].(string)
	if !ok {
		http.Error(w, "data field must be set", http.StatusBadRequest)
	}
	payload, err := hex.DecodeString(data)
	if err != nil {
		http.Error(w, "Invalid data encoding. data should be encoded as a hex string", http.StatusBadRequest)
		return
	}
	if len(payload) == 0 {
		http.Error(w, "data cannot be zero bytes", http.StatusBadRequest)
		return
	}

	if !s.removeDownstreamIfComplete(w, device.DeviceEUI) {
		return
	}

	downstreamMsg := model.NewDownstreamMessage(device.DeviceEUI, uint8(port))
	downstreamMsg.Data = data

	ack, ok := outMessage["ack"].(bool)
	if ok {
		downstreamMsg.Ack = ack
	}
	if err := s.context.Storage.CreateDownstreamData(device.DeviceEUI, downstreamMsg); err != nil {
		lg.Warning("Unable to store downstream message: %v", err)
		http.Error(w, "unable to schedule downstream message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newDownstreamMessageFromModel(downstreamMsg)); err != nil {
		lg.Warning("Unable to marshal downstream message for device with EUI %s into JSON: %v", device.DeviceEUI, err)
	}

}

func (s *Server) getDownstream(device *model.Device, w http.ResponseWriter, r *http.Request) {
	msg, err := s.context.Storage.GetDownstreamData(device.DeviceEUI)
	if err == storage.ErrNotFound {
		http.Error(w, "No downstream message scheduled for device", http.StatusNotFound)
		return
	}
	if err != nil {
		lg.Warning("Unable to retrieve downstream message: %v", err)
		http.Error(w, "Unable to retrieve downstream message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(newDownstreamMessageFromModel(msg)); err != nil {
		lg.Warning("Unable to marshal downstream message for device %s into JSON: %v", device.DeviceEUI, err)
	}
}

func (s *Server) deleteDownstream(device *model.Device, w http.ResponseWriter, r *http.Request) {
	if err := s.context.Storage.DeleteDownstreamData(device.DeviceEUI); err != nil && err != storage.ErrNotFound {
		lg.Warning("Unable to remove downstream message: %v", err)
		http.Error(w, "Unable to remove downstream message", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) deviceSendHandler(w http.ResponseWriter, r *http.Request) {
	_, device := s.getDevice(w, r)
	if device == nil {
		return
	}

	switch r.Method {
	case http.MethodPost:
		s.createDownstream(device, w, r)

	case http.MethodGet:
		s.getDownstream(device, w, r)

	case http.MethodDelete:
		s.deleteDownstream(device, w, r)

	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}
