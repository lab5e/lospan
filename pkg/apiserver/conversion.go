package apiserver

import (
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/pb/lospan"
)

// Convert model.Application -> lospan.Application

func toAPIApplication(app model.Application) *lospan.Application {
	return &lospan.Application{
		Eui: app.AppEUI.String(),
	}
}

func newPtr[T int | bool | float32 | string](v T) *T {
	ret := new(T)
	*ret = v
	return ret
}

func toAPIGateway(gw model.Gateway) *lospan.Gateway {
	return &lospan.Gateway{
		Eui:       gw.GatewayEUI.String(),
		Ip:        newPtr(gw.IP.String()),
		Altitude:  newPtr(gw.Altitude),
		Longitude: newPtr(gw.Longitude),
		Latitude:  newPtr(gw.Latitude),
		StrictIp:  newPtr(gw.StrictIP),
	}
}
