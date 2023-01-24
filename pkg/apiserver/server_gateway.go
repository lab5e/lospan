package apiserver

import (
	"context"
	"errors"
	"net"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/pb/lospan"
	"github.com/lab5e/lospan/pkg/protocol"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *apiServer) CreateGateway(ctx context.Context, req *lospan.Gateway) (*lospan.Gateway, error) {
	var err error

	if req.Ip == nil {
		return nil, status.Error(codes.InvalidArgument, "Missing IP address")
	}
	if req.Eui == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing EUI")
	}

	ip := net.ParseIP(req.GetIp())
	if ip == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid IP")
	}
	if req.GetLongitude() < -360 || req.GetLongitude() > 360 {
		return nil, status.Error(codes.InvalidArgument, "Invalid longitude")
	}
	if req.GetLatitude() < -90 || req.GetLatitude() > 90 {
		return nil, status.Error(codes.InvalidArgument, "Invalid latitude")
	}
	eui, err := protocol.EUIFromString(req.Eui)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid EUI")
	}

	newGW := model.NewGateway()
	newGW.GatewayEUI = eui

	newGW.Latitude = req.GetLatitude()
	newGW.Longitude = req.GetLongitude()
	newGW.StrictIP = true
	if req.StrictIp != nil {
		newGW.StrictIP = req.GetStrictIp()
	}
	newGW.IP = ip
	newGW.Altitude = req.GetAltitude()

	if err := a.store.CreateGateway(newGW); err != nil {
		return nil, toProtoErr(err)
	}

	return toAPIGateway(newGW), nil
}

func (a *apiServer) ListGateways(ctx context.Context, req *lospan.ListGatewaysRequest) (*lospan.ListGatewaysResponse, error) {
	gws, err := a.store.GetGatewayList()
	if err != nil {
		return nil, toProtoErr(err)
	}
	ret := &lospan.ListGatewaysResponse{
		Gateways: make([]*lospan.Gateway, 0),
	}
	for _, gw := range gws {
		ret.Gateways = append(ret.Gateways, toAPIGateway(gw))
	}
	return ret, nil
}

func (a *apiServer) GetGateway(ctx context.Context, req *lospan.GetGatewayRequest) (*lospan.Gateway, error) {
	if req.Eui == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing EUI")
	}
	eui, err := protocol.EUIFromString(req.Eui)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid EUI")
	}
	gw, err := a.store.GetGateway(eui)
	if err != nil {
		return nil, toProtoErr(err)
	}
	return toAPIGateway(gw), nil
}

func (a *apiServer) DeleteGateway(ctx context.Context, req *lospan.DeleteGatewayRequest) (*lospan.Gateway, error) {
	if req.Eui == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing EUI")
	}
	eui, err := protocol.EUIFromString(req.Eui)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid EUI")
	}
	gw, err := a.store.GetGateway(eui)
	if err != nil {
		return nil, toProtoErr(err)
	}
	err = a.store.DeleteGateway(eui)
	if err != nil {
		return nil, toProtoErr(err)
	}

	return toAPIGateway(gw), nil
}

func (a *apiServer) UpdateGateway(ctx context.Context, req *lospan.Gateway) (*lospan.Gateway, error) {

	if req.Eui == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing EUI")
	}
	if req.GetLongitude() < -360 || req.GetLongitude() > 360 {
		return nil, status.Error(codes.InvalidArgument, "Invalid longitude")
	}
	if req.GetLatitude() < -90 || req.GetLatitude() > 90 {
		return nil, status.Error(codes.InvalidArgument, "Invalid latitude")
	}
	eui, err := protocol.EUIFromString(req.Eui)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid EUI")
	}

	gw, err := a.store.GetGateway(eui)
	if err != nil {
		return nil, toProtoErr(err)
	}

	if req.Latitude != nil {
		gw.Latitude = req.GetLatitude()
	}
	if req.Longitude != nil {
		gw.Longitude = req.GetLongitude()
	}
	if req.Altitude != nil {
		gw.Altitude = req.GetAltitude()
	}
	if req.StrictIp != nil {
		gw.StrictIP = req.GetStrictIp()
	}
	if err := a.store.UpdateGateway(gw); err != nil {
		return nil, toProtoErr(err)
	}
	return toAPIGateway(gw), nil
}

func (a *apiServer) StreamGateway(req *lospan.StreamGatewayRequest, stream lospan.Lospan_StreamGatewayServer) error {
	return errors.New("not implemented")
}
