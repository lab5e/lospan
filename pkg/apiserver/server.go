package apiserver

import (
	"github.com/lab5e/lospan/pkg/pb/lospan"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/storage"
)

type apiServer struct {
	store  *storage.Storage
	keyGen *server.KeyGenerator
	lospan.UnimplementedLospanServer
}

// New creates a new API server
func New(store *storage.Storage, keyGen *server.KeyGenerator) (lospan.LospanServer, error) {
	return &apiServer{
		store:  store,
		keyGen: keyGen,
	}, nil
}
