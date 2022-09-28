package apiserver

import (
	"context"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/pb/lospan"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type apiServer struct {
	store  *storage.Storage
	keyGen server.KeyGenerator
	lospan.UnimplementedLospanServer
}

// New creates a new API server
func New(store *storage.Storage, keyGen server.KeyGenerator) (lospan.LospanServer, error) {
	return &apiServer{
		store:  store,
		keyGen: keyGen,
	}, nil
}

func (a *apiServer) ListApplications(ctx context.Context, req *lospan.ListApplicationsRequest) (*lospan.ListApplicationsResponse, error) {
	apps, err := a.store.ListApplications()
	if err != nil {
		return nil, toProtoErr(err)
	}
	ret := &lospan.ListApplicationsResponse{
		Applications: make([]*lospan.Application, 0),
	}
	for _, app := range apps {
		ret.Applications = append(ret.Applications, toAPIApplication(app))
	}
	return ret, nil
}

func (a *apiServer) CreateApplication(ctx context.Context, req *lospan.CreateApplicationRequest) (*lospan.Application, error) {
	newApp := model.NewApplication()
	var err error
	newApp.AppEUI, err = a.keyGen.NewAppEUI()
	if err != nil {
		return nil, toProtoErr(err)
	}
	if err := a.store.CreateApplication(newApp); err != nil {
		return nil, toProtoErr(err)
	}
	return toAPIApplication(newApp), nil
}

func (a *apiServer) DeleteApplication(ctx context.Context, req *lospan.DeleteApplicationRequest) (*lospan.Application, error) {
	eui, err := protocol.EUIFromString(req.Eui)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid EUI")
	}
	existingApp, err := a.store.GetApplicationByEUI(eui)
	if err != nil {
		return nil, toProtoErr(err)
	}

	if err := a.store.DeleteApplication(existingApp.AppEUI); err != nil {
		return nil, toProtoErr(err)
	}

	return toAPIApplication(existingApp), nil
}
