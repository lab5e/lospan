package restapi

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
)

func TestGatewayRoutes(t *testing.T) {
	h := createTestServer(noAuthConfig)
	h.Start()
	defer h.Shutdown()

	// Create two gateways in backend storage
	gwEUI1, _ := protocol.EUIFromString("11-02-03-04-05-06-07-08")
	gwEUI2, _ := protocol.EUIFromString("11-22-03-04-05-06-07-08")

	gw1 := model.Gateway{
		GatewayEUI: gwEUI1,
		IP:         net.ParseIP("127.0.0.1"),
		StrictIP:   true,
		Latitude:   1,
		Longitude:  2,
		Altitude:   3,
	}
	err := h.context.Storage.CreateGateway(gw1)
	if err != nil {
		t.Fatal("Error writing gateway 1 to storage: ", err)
	}

	gw2 := model.Gateway{
		GatewayEUI: gwEUI2,
		IP:         net.ParseIP("127.0.0.2"),
		StrictIP:   true,
		Latitude:   4,
		Longitude:  5,
		Altitude:   6,
	}

	err = h.context.Storage.CreateGateway(gw2)
	if err != nil {
		t.Fatal("Error writing gateway 2 to storage: ", err)
	}

	// Retrieve list of gateways from endpoint. Both should be there
	resp, err := http.Get(h.loopbackURL() + "/gateways")
	if err != nil {
		t.Fatal("Got error querying list of gateways: ", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Didn't get 200 OK when querying list. Got ", resp.StatusCode)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Error reading body: ", err)
	}

	gwList := newGatewayList()
	if err := json.Unmarshal(buf, &gwList); err != nil {
		t.Fatal("Error unmarshaling list: ", err)
	}
	var foundOne, foundTwo bool
	for _, v := range gwList.Gateways {
		modelGw := v.ToModel()
		if modelGw.Equals(gw1) {
			foundOne = true
		}
		if modelGw.Equals(gw2) {
			foundTwo = true
		}
	}

	if !foundOne || !foundTwo {
		t.Fatalf("Couldn't find both stored gateways in response 1:%t, 2:%t", foundOne, foundTwo)
	}

	singleEUI, _ := protocol.EUIFromString("11-22-33-44-55-66-77-88")

	// Store a gateway through the http interface. It should be in the gw Store
	apiGw := apiGateway{
		GatewayEUI: singleEUI.String(),
		IP:         "127.0.1.1",
		StrictIP:   true,
		Altitude:   9,
		Latitude:   8,
		Longitude:  7,
	}

	buf, err = json.Marshal(&apiGw)
	if err != nil {
		t.Fatal("Got error marshaling gateway")
	}
	reader := strings.NewReader(string(buf))

	resp, err = http.Post(h.loopbackURL()+"/gateways", "application/json", reader)
	if err != nil {
		t.Fatal("Got error POSTing: ", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatal("Expected 201 CREATED but got ", resp.StatusCode)
	}

	gw, err := h.context.Storage.GetGateway(singleEUI)
	if err != nil {
		t.Fatal("Could not find POSTed gateway in store: ", err)
	}
	if !gw.Equals(apiGw.ToModel()) {
		t.Fatalf("Stored and POSTed gateway does not match (%v != %v)", gw, apiGw.ToModel())
	}
}

func TestGatewayListEndpoint(t *testing.T) {
	h := createTestServer(noAuthConfig)
	h.Start()
	defer h.Shutdown()

	const duplicateEUI = "aa-bb-aa-bb-aa-bb-aa-bb"
	dupEUI, _ := protocol.EUIFromString(duplicateEUI)
	h.context.Storage.CreateGateway(model.Gateway{GatewayEUI: dupEUI, IP: net.ParseIP("127.0.0.1"), StrictIP: false})

	rootURL := h.loopbackURL() + "/gateways"

	invalidPosts := map[string]int{
		`{}`:                                  http.StatusBadRequest,
		`{EUI: ""}`:                           http.StatusBadRequest,
		`{"gatewayEUI": ""}`:                  http.StatusBadRequest,
		`{"gatewayEUI": "01-02", "ip": "12"}`: http.StatusBadRequest,
		`{"gatewayEUI": "01-02-03-04-05-06-07-08"}`:                                                           http.StatusBadRequest,
		`{"gatewayEUI": "01-02-03-04-05-06-07-08", "ip": "something"}`:                                        http.StatusBadRequest,
		`{"gatewayEUI": "aa-02-03-04-05-06-07-08", "ip": "127.0.0.1"}`:                                        http.StatusCreated,
		`{"gatewayEUI": "` + duplicateEUI + `", "ip": "127.0.0.1"}`:                                           http.StatusConflict,
		`{"gatewayEUI": "01-02-03-04-05-06-07-09", "ip": "127.0.0.1", "latitude": 90.0, "longitude": 180.0}`:  http.StatusCreated,
		`{"gatewayEUI": "01-02-03-04-05-06-07-10", "ip": "127.0.0.1", "latitude": 900.0, "longitude": 180.0}`: http.StatusBadRequest,
		`{"gatewayEUI": "01-02-03-04-05-06-07-11", "ip": "127.0.0.1", "latitude": 90.0, "longitude": 1800.0}`: http.StatusBadRequest,
	}

	invalidGets := map[string]int{
		// All parameters are ignored and no parameters in path
	}

	invalidMethods := []string{
		"HEAD",
		"PATCH",
		"PUT",
		"DELETE",
	}

	genericEndpointTest(t, rootURL, invalidGets, invalidPosts, invalidMethods)

}

func TestGatewayInfoEndpoint(t *testing.T) {
	h := createTestServer(noAuthConfig)
	h.Start()
	defer h.Shutdown()

	eui, _ := protocol.EUIFromString("01-23-45-67-89-AB-CD-EF")
	gw := model.Gateway{
		GatewayEUI: eui,
		IP:         net.ParseIP("127.0.0.1"),
	}
	if err := h.context.Storage.CreateGateway(gw); err != nil {
		t.Fatal("Couldn't create gw: ", err)
	}

	rootURL := h.loopbackURL() + "/gateways/" + eui.String()

	invalidPosts := map[string]int{
		// No posts here
	}

	invalidGets := map[string]int{
		h.loopbackURL() + "/gateways/01-02":                   http.StatusBadRequest,
		h.loopbackURL() + "/gateways/01-02-03-04-01-02-03-04": http.StatusNotFound,
		h.loopbackURL() + "/gateways/" + eui.String():         http.StatusOK,
	}

	invalidMethods := []string{
		"HEAD",
		"PATCH",
		"POST",
	}

	genericPutRequest(t, rootURL, map[string]interface{}{
		"ip":        "10.10.10.10",
		"altitude":  1,
		"latitude":  2,
		"longitude": 3,
		"strictIp":  false,
	}, http.StatusOK)
	genericPutRequest(t, rootURL, map[string]interface{}{
		"ip": "10.10x10.10",
	}, http.StatusBadRequest)
	genericPutRequest(t, rootURL, map[string]interface{}{
		"tags": map[string]string{"name": "value"},
	}, http.StatusOK)
	genericPutRequest(t, rootURL, map[string]interface{}{
		"tags": map[string]string{"name": "alert('Hello');"},
	}, http.StatusBadRequest)
	genericPutRequest(t, rootURL, map[string]interface{}{
		"tags": map[string]interface{}{"name": true, "value": 12},
	}, http.StatusBadRequest)

	genericEndpointTest(t, rootURL, invalidGets, invalidPosts, invalidMethods)
	testDelete(t, map[string]int{
		h.loopbackURL() + "/gateways/01-02":                   http.StatusBadRequest,
		h.loopbackURL() + "/gateways/" + eui.String():         http.StatusNoContent,
		h.loopbackURL() + "/gateways/01-02-03-04-05-06-07-aa": http.StatusNotFound,
	})

}
