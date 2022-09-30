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

func newPtr[T int | int32 | uint32 | int64 | bool | float32 | string](v T) *T {
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

func toAPIState(s model.DeviceState) *lospan.DeviceState {
	ret := lospan.DeviceState_UNSPECIFIED
	switch s {
	case model.DisabledDevice:
		ret = lospan.DeviceState_DISABLED
	case model.OverTheAirDevice:
		ret = lospan.DeviceState_OTAA
	case model.PersonalizedDevice:
		ret = lospan.DeviceState_ABP
	default:
		ret = lospan.DeviceState_DISABLED
	}
	return &ret
}

func toAPINonces(nonces []uint16) []int32 {
	var ret []int32
	for _, n := range nonces {
		ret = append(ret, int32(n))
	}
	return ret
}

func toAPIDevice(d model.Device) *lospan.Device {
	return &lospan.Device{
		Eui:               newPtr(d.DeviceEUI.String()),
		ApplicationEui:    newPtr(d.AppEUI.String()),
		State:             toAPIState(d.State),
		DevAddr:           newPtr(d.DevAddr.ToUint32()),
		AppKey:            d.AppKey.Key[:],
		AppSessionKey:     d.AppSKey.Key[:],
		NetworkSessionKey: d.NwkSKey.Key[:],
		FrameCountUp:      newPtr(int32(d.FCntUp)),
		FrameCountDown:    newPtr(int32(d.FCntDn)),
		RelaxedCounter:    newPtr(d.RelaxedCounter),
		KeyWarning:        newPtr(d.KeyWarning),
		DevNonces:         toAPINonces(d.DevNonceHistory[:]),
	}
}
