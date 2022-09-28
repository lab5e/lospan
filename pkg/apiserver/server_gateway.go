package apiserver

import (
	"context"
	"errors"

	"github.com/lab5e/lospan/pkg/pb/lospan"
)

func (a *apiServer) CreateGateway(ctx context.Context, req *lospan.Gateway) (*lospan.Gateway, error) {
	return nil, errors.New("no")
}

func (a *apiServer) ListGateways(ctx context.Context, req *lospan.ListGatewaysRequest) (*lospan.ListGatewaysResponse, error) {
	return nil, errors.New("no")
}

func (a *apiServer) GetGateway(ctx context.Context, req *lospan.GetGatewayRequest) (*lospan.Gateway, error) {
	return nil, errors.New("no")
}

func (a *apiServer) DeleteGateway(ctx context.Context, req *lospan.DeleteGatewayRequest) (*lospan.Gateway, error) {
	return nil, errors.New("no")
}

func (a *apiServer) UpdateGateway(ctx context.Context, req *lospan.Gateway) (*lospan.Gateway, error) {
	return nil, errors.New("no")
}
