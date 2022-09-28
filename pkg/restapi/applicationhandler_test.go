package restapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func storeApplication(t *testing.T, application apiApplication, url string, expectedStatus int) apiApplication {
	bytes, err := json.Marshal(application)
	if err != nil {
		t.Fatalf("Got error marshaling application: %v", err)
	}
	reader := strings.NewReader(string(bytes))

	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		t.Fatalf("Could not POST application to %s: %v", url, err)
	}
	if resp.StatusCode != expectedStatus {
		t.Fatalf("POSTed successfully to %s but got %d (expected %d)", url, resp.StatusCode, expectedStatus)
	}

	ret := apiApplication{}
	if expectedStatus == http.StatusCreated {
		buffer, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Couldn't read response body from %s: %v", url, err)
		}
		if err := json.Unmarshal(buffer, &ret); err != nil {
			t.Fatalf("Couldn't unmarshal application: %v", err)
		}
	}

	return ret
}
