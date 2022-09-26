package restapi

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
	"net/http"
	"testing"

	"github.com/ExploratoryEngineering/congress/model"
)

func TestTokenChecker(t *testing.T) {
	server := createTestServer(noAuthConfig)
	// Set up memory storage with a token
	newToken, err := model.NewAPIToken("007", "/something", false)
	if err != nil {
		t.Fatal("Got error creating token: ", err)
	}
	if err := server.context.Storage.Token.Put(newToken, model.SystemUserID); err != nil {
		t.Fatal("Got error storing token: ", err)
	}

	testToken := func(token string, method string, path string, expectedStatus int) {
		success, message, status, _ := server.isValidToken(token, method, path)
		if status != expectedStatus {
			t.Errorf("Got %d but expected %d when requesting %s %s (msg: %s)",
				status, expectedStatus, method, path, message)
		}
		if status == http.StatusOK && !success {
			t.Errorf("Got false return with %d status code", status)
		}
	}

	// Should be OK
	testToken(newToken.Token, "GET", "/something", http.StatusOK)
	testToken(newToken.Token, "OPTIONS", "/something", http.StatusOK)
	testToken(newToken.Token, "HEAD", "/something", http.StatusOK)
	testToken(newToken.Token, "GET", "/something/else", http.StatusOK)

	// Invalid resource
	testToken(newToken.Token, "HEAD", "/other", http.StatusForbidden)
	testToken(newToken.Token, "GET", "/other/else", http.StatusForbidden)

	// Invalid token
	testToken("newToken.Token", "GET", "/something", http.StatusUnauthorized)
	testToken("newToken.Token", "GET", "/something/else", http.StatusUnauthorized)

	// No token header.
	testToken("", "GET", "/something", http.StatusUnauthorized)
	testToken("", "GET", "/something/else", http.StatusUnauthorized)

	// Illegal methods, won't get access even with the proper token
	testToken(newToken.Token, "DELETE", "/something", http.StatusForbidden)
	testToken(newToken.Token, "POST", "/something", http.StatusForbidden)
	testToken(newToken.Token, "PATCH", "/something", http.StatusForbidden)

}
