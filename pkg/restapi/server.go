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
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/ExploratoryEngineering/logging"
	"github.com/ExploratoryEngineering/rest"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/utils"

	"golang.org/x/crypto/acme/autocert"
)

// Server is a type capable of serving the REST API for Congress. It can be started
// and dhut down only once reliably since the port lingers. There is no check if
// the server is running so calling Start() twice will result in errors
type Server struct {
	srv       *http.Server
	mux       *http.ServeMux
	context   *server.Context
	config    *server.Configuration
	port      int
	completed chan bool
}

// NewServer returns a new server instance. if the port is set to 0 it will
// pick a random port. If loopbackOnly is true only the loopback adapter
// will be used.
func NewServer(loopbackOnly bool, scontext *server.Context, config *server.Configuration) (*Server, error) {
	ret := &Server{context: scontext, config: config, completed: make(chan bool)}
	portno := config.HTTPServerPort
	var err error
	if portno == 0 {
		portno, err = utils.FreePort()
		if err != nil {
			return nil, err
		}
	}
	ret.port = portno

	host := ""
	if loopbackOnly {
		host = "localhost"
	}
	ret.mux = http.NewServeMux()

	ret.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, ret.port),
		Handler: ret,
	}

	if config.ACMECert {
		logging.Info("Using Let's Encrypt for certificates")
		// See https://godoc.org/golang.org/x/crypto/acme/autocert#example-Manager
		m := &autocert.Manager{
			Cache:      autocert.DirCache(config.ACMESecretDir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(config.ACMEHost),
		}
		go http.ListenAndServe(":http", m.HTTPHandler(nil))
		ret.srv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
	}

	handler := ret.handler()

	ret.mux.HandleFunc("/", rest.AddCORSHeaders(handler))
	return ret, nil
}

// Start launches the server. The server won't check if it has been started twice
func (h *Server) Start() error {
	logging.Info("HTTP server listening on port %d", h.port)
	go func() {
		if h.config.ACMECert {
			if err := h.srv.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
				logging.Error("ListenAndServeTLS returned error: %v", err)
			}
		} else if h.config.TLSCertFile != "" && h.config.TLSKeyFile != "" {
			if err := h.srv.ListenAndServeTLS(h.config.TLSCertFile, h.config.TLSKeyFile); err != http.ErrServerClosed {
				logging.Error("ListenAndServeTLS returned error: %v", err)
			}
		} else {
			if err := h.srv.ListenAndServe(); err != http.ErrServerClosed {
				logging.Error("ListenAndServe returned error: %v", err)
			}
		}
		h.completed <- true
	}()
	return nil
}

// Shutdown stops the server. There's no check if the server is already running. Run Shutdown() twice at your own risk.
func (h *Server) Shutdown() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()
	if err := h.srv.Shutdown(ctx); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
	case <-h.completed:
	}
	return nil
}

// loopbackURL returns the loopback URL for the server. Used for testing
func (h *Server) loopbackURL() string {
	return fmt.Sprintf("http://localhost:%d", h.port)
}

func (h *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// Handler returns a HandlerFunc with all the routes for the endpoint
func (h *Server) handler() http.HandlerFunc {
	router := rest.NewParameterRouter()
	router.AddRoute("/", h.rootHandler)
	router.AddRoute("/applications", h.applicationListHandler)
	router.AddRoute("/applications/{aeui}", h.applicationInfoHandler)
	router.AddRoute("/applications/{aeui}/devices", h.deviceListHandler)
	router.AddRoute("/applications/{aeui}/devices/{deui}", h.deviceInfoHandler)
	router.AddRoute("/applications/{aeui}/devices/{deui}/message", h.deviceSendHandler)
	router.AddRoute("/gateways", h.gatewayListHandler)
	router.AddRoute("/gateways/{geui}", h.gatewayInfoHandler)

	return func(w http.ResponseWriter, r *http.Request) {
		router.GetHandler(r.RequestURI).ServeHTTP(w, r)
	}
}

// Extract the corresponding storage.UserID from the ID session. If auth is
// disabled return the system user.
func (h *Server) connectUserID(r *http.Request) model.UserID {
	return model.SystemUserID
}

// updateTags updates tags from a JSON struct in a request, Returns false if the
// struct contains invalid tags. This is used by both application, gateway and device
// resources.
func (h *Server) updateTags(tags *model.Tags, values map[string]interface{}) bool {
	updatedTags, ok := values["tags"].(map[string]interface{})
	if !ok {
		return true
	}
	for k, v := range updatedTags {
		val, ok := v.(string)
		if !ok {
			logging.Debug("Not a string type %s=%v(%T)", k, v, v)
			return false
		}
		if err := tags.SetTag(k, val); err != nil {
			logging.Debug("Invalid tag %s=%v: %v", k, v, err)
			return false
		}

	}
	return true
}
