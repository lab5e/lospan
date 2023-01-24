package apiserver

import (
	"github.com/lab5e/lospan/pkg/keys"
	"github.com/lab5e/lospan/pkg/pb/lospan"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/storage"
)

type apiServer struct {
	store  *storage.Storage
	keyGen *keys.KeyGenerator
	router *server.EventRouter[protocol.EUI, *server.PayloadMessage]
}

// New creates a new API server
func New(store *storage.Storage, keyGen *keys.KeyGenerator, router *server.EventRouter[protocol.EUI, *server.PayloadMessage]) (lospan.LospanServer, error) {
	return &apiServer{
		store:  store,
		keyGen: keyGen,
		router: router,
	}, nil
}
