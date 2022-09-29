package restapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
)

func storeDevice(t *testing.T, device apiDevice, url string, expectedStatus int) apiDevice {
	bytes, err := json.Marshal(device)
	if err != nil {
		t.Fatalf("Got error marshalling device: %v", err)
	}

	reader := strings.NewReader(string(bytes))

	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		t.Fatalf("Could not POST device to %s: %v", url, err)
	}
	if resp.StatusCode != expectedStatus {
		t.Fatalf("POST successfully to %s but response code was %d, not %d", url, resp.StatusCode, expectedStatus)
	}

	ret := apiDevice{}
	if expectedStatus == http.StatusCreated {
		buffer, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Couldn't read response body from %s: %v", url, err)
		}
		if err := json.Unmarshal(buffer, &ret); err != nil {
			t.Fatalf("Couldn't unmarshal device from %s: %v", url, err)
		}
		devaddr, err := protocol.DevAddrFromString(ret.DevAddr)
		if err != nil {
			t.Fatal("Error decoding devaddr: ", err)
		}
		if devaddr.ToUint32() == 0 {
			t.Fatal("DevAddr == 0. LMiC hates this")
		}
	}
	return ret
}

func TestDeviceDataEndpoint(t *testing.T) {
	h := createTestServer(noAuthConfig)
	h.Start()
	defer h.Shutdown()

	appsURL := h.loopbackURL() + "/applications"
	application := storeApplication(t, apiApplication{}, appsURL, http.StatusCreated)

	appURL := appsURL + "/" + application.ApplicationEUI
	device := storeDevice(t, apiDevice{}, appURL+"/devices", http.StatusCreated)

	eui, _ := protocol.EUIFromString(device.DeviceEUI)
	for i := 0; i < 10; i++ {
		err := h.context.Storage.CreateUpstreamData(eui, application.eui, model.UpstreamMessage{
			DeviceEUI:  eui,
			AppEUI:     application.eui,
			Timestamp:  int64(i),
			Data:       []byte{0, 1, 2, 3, 4, 5},
			GatewayEUI: eui,
		})
		if err != nil {
			t.Fatal("Got error storing device data: ", err)
		}
	}

	deviceURL := appURL + "/devices/" + device.DeviceEUI

	invalidPosts := map[string]int{
		// No POST on this endpoint
	}

	urlTemplate := h.loopbackURL() + "/applications/%s/devices/%s/data"
	const invalidEUI string = "00-00-00-00-00-00-00-00"
	const improperEUI string = "00-00"

	validURL := fmt.Sprintf(urlTemplate, application.ApplicationEUI, device.DeviceEUI)
	invalidDeviceEUI := fmt.Sprintf(urlTemplate, application.ApplicationEUI, invalidEUI)
	incorrectDeviceEUI := fmt.Sprintf(urlTemplate, application.ApplicationEUI, improperEUI)
	invalidAppEUI := fmt.Sprintf(urlTemplate, invalidEUI, device.DeviceEUI)
	incorrectAppEUI := fmt.Sprintf(urlTemplate, improperEUI, device.DeviceEUI)
	invalidGets := map[string]int{
		validURL:           http.StatusOK,
		invalidDeviceEUI:   http.StatusNotFound,
		incorrectDeviceEUI: http.StatusBadRequest,
		invalidAppEUI:      http.StatusNotFound,
		incorrectAppEUI:    http.StatusBadRequest,
	}
	invalidMethods := []string{
		"HEAD",
		"PATCH",
		"PUT",
		"DELETE",
		"POST",
	}
	genericEndpointTest(t, deviceURL+"/data", invalidGets, invalidPosts, invalidMethods)
}

func TestDeviceMessageInput(t *testing.T) {
	h := createTestServer(noAuthConfig)
	h.Start()
	defer h.Shutdown()

	application := storeApplication(t, apiApplication{}, h.loopbackURL()+"/applications", http.StatusCreated)

	appURL := h.loopbackURL() + "/applications/" + application.ApplicationEUI
	device := storeDevice(t, apiDevice{}, appURL+"/devices", http.StatusCreated)

	eui, _ := protocol.EUIFromString(device.DeviceEUI)
	for i := 0; i < 10; i++ {
		err := h.context.Storage.CreateUpstreamData(eui, application.eui, model.UpstreamMessage{
			DeviceEUI:  eui,
			Timestamp:  int64(i),
			Data:       []byte{0, 1, 2, 3, 4, 5},
			GatewayEUI: eui,
		})
		if err != nil {
			t.Fatal("Got error storing device data: ", err)
		}
	}

	deviceURL := appURL + "/devices/" + device.DeviceEUI

	invalidPosts := map[string]int{
		`x`:                         http.StatusBadRequest,
		`{}`:                        http.StatusBadRequest,
		`{"port": -1}`:              http.StatusBadRequest,
		`{"port": 254}`:             http.StatusBadRequest,
		`{"port": 999}`:             http.StatusBadRequest,
		`{"port": 1, "data": "zy"}`: http.StatusBadRequest,
		`{"port": 1, "data": ""}`:   http.StatusBadRequest,
	}

	invalidGets := map[string]int{}

	invalidMethods := []string{
		"HEAD",
		"PATCH",
		"PUT",
	}

	genericEndpointTest(t, deviceURL+"/message", invalidGets, invalidPosts, invalidMethods)

	// Reset output buffer
	h.context.Storage.DeleteDownstreamData(eui)

	// Post a single message and ensure the output buffer is set.
	reader := strings.NewReader(`{"port": 1, "data": "01AA02BB03CC04DD", "ack": true}`)
	resp, _ := http.Post(deviceURL+"/message", "application/json", reader)
	if resp.StatusCode != http.StatusCreated {
		buf, _ := io.ReadAll(resp.Body)
		t.Fatalf("Got status %d with body: %s posting to %s", resp.StatusCode, string(buf), deviceURL+"/message")
	}
	msg, err := h.context.Storage.GetDownstreamData(eui)
	if err != nil {
		t.Fatalf("Expected a new message to be created but got error: %v", err)
	}

	if msg.Data != "01AA02BB03CC04DD" {
		t.Fatalf("Not the expected payload: %v", msg.Data)
	}

	// Retrieve an upstream message
	newMsg := model.NewDownstreamMessage(eui, 100)
	newMsg.Data = "000102030405"
	newMsg.Ack = true

	url := h.loopbackURL() + "/applications/" + application.ApplicationEUI + "/devices/" + device.DeviceEUI + "/message"
	h.context.Storage.DeleteDownstreamData(eui)
	resp, _ = http.Get(url)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404 NOT FOUND but got %d", resp.StatusCode)
	}
	h.context.Storage.CreateDownstreamData(eui, newMsg)

	resp, _ = http.Get(url)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Did not get 200 OK from downstream message but got %d", resp.StatusCode)
	}

	downMsg := apiDownstreamMessage{}
	if err := json.NewDecoder(resp.Body).Decode(&downMsg); err != nil {
		t.Fatalf("Got error decoding upstream response: %v", err)
	}
	if !downMsg.Ack || downMsg.Data != newMsg.Data || downMsg.Port != newMsg.Port {
		t.Fatalf("Got different response from what was created. Got %v, expected %v", downMsg, newMsg)
	}

	// Call DELETE on the resource. Should return 204 NO Content, even when there's no
	// downstream message
	testDelete(t, map[string]int{
		url: http.StatusNoContent,
		url: http.StatusNoContent,
		h.loopbackURL() + "/applications/" + application.ApplicationEUI + "/devices/00/message": http.StatusBadRequest,
	})
}

func TestAutomaticDownstreamMessageRemoval(t *testing.T) {
	h := createTestServer(noAuthConfig)
	h.Start()
	defer h.Shutdown()

	application := storeApplication(t, apiApplication{}, h.loopbackURL()+"/applications", http.StatusCreated)
	deviceURL := fmt.Sprintf("%s/applications/%s/devices", h.loopbackURL(), application.ApplicationEUI)
	device := storeDevice(t, apiDevice{}, deviceURL, http.StatusCreated)
	messageURL := fmt.Sprintf("%s/applications/%s/devices/%s/message", h.loopbackURL(), application.ApplicationEUI, device.DeviceEUI)

	createMessage := func(msgData string, expectedStatus int) {
		resp, err := http.Post(messageURL, "application/json", strings.NewReader(msgData))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != expectedStatus {
			t.Fatalf("Expected %d but got %d response", expectedStatus, resp.StatusCode)
		}
	}

	updateSentAck := func(sent, ack int64) {
		eui, _ := protocol.EUIFromString(device.DeviceEUI)
		if err := h.context.Storage.UpdateDownstreamData(eui, sent, ack); err != nil {
			t.Fatal(err)
		}
	}
	// Schedule a new downstream message. Should succeed.
	createMessage(`{"port": 100, "data": "aabbccdd", "ack": false}`, http.StatusCreated)

	// Schedule another message. Should fail.
	createMessage(`{"port": 101, "data": "aabbccdd", "ack": false}`, http.StatusConflict)

	// Mimic sent status by updating sent field.
	updateSentAck(12, 0)

	// Schedule another message. Should succeed since the message is sent.
	createMessage(`{"port": 102, "data": "aabbccdd", "ack": true}`, http.StatusCreated)

	// Mimic sent status
	updateSentAck(13, 0)

	// Schedule another. Should fail since the message isn't acked
	createMessage(`{"port": 103, "data": "aabbccdd", "ack": false}`, http.StatusConflict)

	// Mimic ack status
	updateSentAck(14, 15)

	// Schedule another. Should succeed since the message is acked.
	createMessage(`{"port": 104, "data": "aabbccdd", "ack": false}`, http.StatusCreated)
}
