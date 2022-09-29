package restapi

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/storage"
)

var ma = protocol.MA{Prefix: [5]byte{0, 1, 3, 4, 5}, Size: protocol.MALarge}

// Configuration without authentication
var noAuthConfig = server.Configuration{HTTPServerPort: 0, MemoryDB: true}

func createTestServer(config server.Configuration) *Server {
	ma, _ := protocol.NewMA([]byte{1, 2, 3})

	// NetID
	netID := uint32(0)

	store := storage.NewMemoryStorage()

	keygen, _ := server.NewEUIKeyGenerator(ma, netID, store)

	fob := server.NewFrameOutputBuffer()

	appRouter := server.NewEventRouter[protocol.EUI, *server.PayloadMessage](5)
	context := &server.Context{
		Storage:      store,
		FrameOutput:  &fob,
		KeyGenerator: &keygen,
		AppRouter:    &appRouter,
		Config:       &config,
	}

	server, _ := NewServer(true, context, &config)
	return server
}

// Enable authentication and ensure it starts as expected (and authenticates)

func checkContentType(t *testing.T, url string, resp *http.Response) {
	if resp.StatusCode >= 200 && resp.StatusCode <= 300 {
		if resp.Header.Get("Content-Type") != "application/json" {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected content-header for %s to say application/json but it says %s with body %s", url, resp.Header.Get("Content-Type"), string(body))
		}
	}
}

// This is a generic endpoint test. It will test invalid URLs (if the URLs
// contains path parameters), invalid POSTs (to the root URL) and verify that
// invalid methods return the appropriate status code (405)
func genericEndpointTest(t *testing.T, rootURL string, invalidGets map[string]int,
	invalidPosts map[string]int, invalidMethods []string) {

	for url, expectedValue := range invalidGets {
		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("Got error %v when GETting %s", err, url)
		}
		if resp.StatusCode != expectedValue {
			body, err := io.ReadAll(resp.Body)
			var output = "<no output>"
			if err == nil {
				output = string(body)
			}
			t.Fatalf("Got status code %d but expected %d for URL %s (body is %s)",
				resp.StatusCode, expectedValue, url, strings.TrimSpace(output))
		}

		checkContentType(t, url, resp)
	}

	// Test POST against the root URL
	for body, expectedValue := range invalidPosts {
		reader := strings.NewReader(body)
		resp, err := http.Post(rootURL, "application/json", reader)
		if err != nil {
			t.Fatalf("Error retrieving URL %s: %v", rootURL, err)
		}
		if resp.StatusCode != expectedValue {
			buf, _ := io.ReadAll(resp.Body)
			t.Errorf("Got %d but expected %d when POSTing %s (response body is %s)", resp.StatusCode, expectedValue, body, string(buf))
		}
		checkContentType(t, rootURL, resp)
	}

	for _, method := range invalidMethods {
		emptyBody := strings.NewReader("")
		req, err := http.NewRequest(method, rootURL, emptyBody)
		if err != nil {
			t.Fatal("Got error creating request: ", err)
		}
		client := &http.Client{}
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal("Got error performing request: ", err)
		}

		if resp.StatusCode != http.StatusMethodNotAllowed {
			buf := make([]byte, 100)
			resp.Body.Read(buf)
			t.Fatalf("Didn't get 405 Method Not Allowed when doing %s (got %d with body '%s')", method, resp.StatusCode, string(buf))
		}
	}
}

func testDelete(t *testing.T, testData map[string]int) {
	for url, expected := range testData {
		client := http.Client{}
		emptyBody := strings.NewReader("")
		req, err := http.NewRequest("DELETE", url, emptyBody)
		if err != nil {
			t.Fatal("Got error creating request: ", err)
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Got error calling DELETE on resource %s: %v", url, err)
		}
		if resp.StatusCode != expected {
			t.Fatalf("Expected %d but got %d when calling DELETE on %s", expected, resp.StatusCode, url)
		}
	}
}
