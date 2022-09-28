package apiserver

import (
	"context"
	"errors"

	"github.com/lab5e/lospan/pkg/pb/lospan"
	"github.com/lab5e/lospan/pkg/storage"
)

type apiServer struct {
	store *storage.Storage
	lospan.UnimplementedLospanServer
}

// New creates a new API server
func New(store *storage.Storage) (lospan.LospanServer, error) {
	return &apiServer{
		store: store,
	}, nil
}

func (a *apiServer) ListApplications(ctx context.Context, req *lospan.ListApplicationsRequest) (*lospan.ListApplicationsResponse, error) {
	return nil, errors.New("not implemented")
}
