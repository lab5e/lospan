package apiserver

import (
	"context"
	"encoding/hex"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/pb/lospan"
	"github.com/lab5e/lospan/pkg/protocol"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toState(s *lospan.DeviceState) (model.DeviceState, error) {
	ret := model.OverTheAirDevice
	if s != nil {
		switch *s {
		case lospan.DeviceState_ABP:
			ret = model.PersonalizedDevice
		case lospan.DeviceState_OTAA:
			ret = model.OverTheAirDevice
		case lospan.DeviceState_DISABLED:
			ret = model.DisabledDevice
		default:
			return ret, status.Error(codes.InvalidArgument, "Device state must be OTAA, ABP or disabled")
		}
	}
	return ret, nil
}

func (a *apiServer) CreateDevice(ctx context.Context, req *lospan.Device) (*lospan.Device, error) {
	var eui protocol.EUI
	var err error
	if req.Eui == nil {
		eui, err = a.keyGen.NewDeviceEUI()
		if err != nil {
			return nil, toProtoErr(err)
		}
	}
	if req.Eui != nil {
		eui, err = protocol.EUIFromString(req.GetEui())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Invalid EUI")
		}
	}
	if req.ApplicationEui == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing application EUI")
	}
	d := model.NewDevice()
	d.AppEUI, err = protocol.EUIFromString(req.GetEui())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Application EUI")
	}
	d.DeviceEUI = eui
	d.KeyWarning = false
	d.State, err = toState(req.State)
	if err != nil {
		return nil, err
	}

	if req.DevAddr != nil {
		d.DevAddr = protocol.DevAddrFromUint32(req.GetDevAddr())
	}

	if req.AppKey != nil {
		d.AppKey, err = protocol.AESKeyFromString(hex.EncodeToString(req.AppKey))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Invalid App Key")
		}
	}
	if req.AppSessionKey != nil {
		d.AppSKey, err = protocol.AESKeyFromString(hex.EncodeToString(req.AppSessionKey))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Invalid App Session Key")
		}
	}
	if req.NetworkSessionKey != nil {
		d.NwkSKey, err = protocol.AESKeyFromString(hex.EncodeToString(req.NetworkSessionKey))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Invalid Network Session Key")
		}
	}
	d.RelaxedCounter = false
	if req.RelaxedCounter != nil {
		d.RelaxedCounter = req.GetRelaxedCounter()
	}
	if req.KeyWarning != nil {
		d.KeyWarning = req.GetKeyWarning()
	}
	d.FCntDn = 0
	if req.FrameCountDown != nil {
		d.FCntDn = uint16(req.GetFrameCountDown())
	}
	d.FCntUp = 0
	if req.FrameCountUp != nil {
		d.FCntUp = uint16(req.GetFrameCountUp())
	}
	if d.State == model.OverTheAirDevice && (len(d.AppSKey.Key) > 0 || d.DevAddr.ToUint32() != 0 || len(d.NwkSKey.Key) > 0) {
		return nil, status.Error(codes.InvalidArgument, "DevAddr, AppSKey and NwkSKey can only be specified for ABP devices")
	}
	if d.State == model.PersonalizedDevice && len(d.AppKey.Key) > 0 {
		return nil, status.Error(codes.InvalidArgument, "AppKey can only be specified for OTAA devices")
	}

	if d.State == model.OverTheAirDevice {
		if len(d.AppKey.Key) == 0 {
			d.AppKey, err = protocol.NewAESKey()
			if err != nil {
				return nil, toProtoErr(err)
			}
		}
	}
	if d.State == model.PersonalizedDevice {
		if len(d.AppSKey.Key) == 0 {
			d.AppSKey, err = protocol.NewAESKey()
			if err != nil {
				return nil, toProtoErr(err)
			}
		}
		if len(d.NwkSKey.Key) == 0 {
			d.NwkSKey, err = protocol.NewAESKey()
			if err != nil {
				return nil, toProtoErr(err)
			}
		}
		if d.DevAddr.ToUint32() == 0 {
			d.DevAddr = protocol.NewDevAddr()
		}
	}

	if err := a.store.CreateDevice(d, d.AppEUI); err != nil {
		return nil, toProtoErr(err)
	}

	return toAPIDevice(d), nil
}

func (a *apiServer) GetDevice(ctx context.Context, req *lospan.GetDeviceRequest) (*lospan.Device, error) {
	eui, err := protocol.EUIFromString(req.GetEui())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid EUI")
	}

	d, err := a.store.GetDeviceByEUI(eui)
	if err != nil {
		return nil, toProtoErr(err)
	}
	return toAPIDevice(d), nil
}

func (a *apiServer) DeleteDevice(ctx context.Context, req *lospan.DeleteDeviceRequest) (*lospan.Device, error) {
	eui, err := protocol.EUIFromString(req.GetEui())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid EUI")
	}

	d, err := a.store.GetDeviceByEUI(eui)
	if err != nil {
		return nil, toProtoErr(err)
	}
	if err := a.store.DeleteDevice(d.DeviceEUI); err != nil {
		return nil, toProtoErr(err)
	}
	return toAPIDevice(d), nil
}

func (a *apiServer) UpdateDevice(ctx context.Context, req *lospan.Device) (*lospan.Device, error) {
	eui, err := protocol.EUIFromString(req.GetEui())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid EUI")
	}

	d, err := a.store.GetDeviceByEUI(eui)
	if err != nil {
		return nil, toProtoErr(err)
	}
	if req.State != nil {
		d.State, err = toState(req.State)
		if err != nil {
			return nil, err
		}
	}
	if req.DevAddr != nil {
		d.DevAddr = protocol.DevAddrFromUint32(req.GetDevAddr())
	}
	if len(req.AppKey) > 0 {
		d.AppKey, err = protocol.AESKeyFromString(hex.EncodeToString(req.AppKey))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Invalid App Key")
		}
	}
	if len(req.AppSessionKey) > 0 {
		d.AppSKey, err = protocol.AESKeyFromString(hex.EncodeToString(req.AppSessionKey))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Invalid App Session Key")
		}
	}
	if len(req.NetworkSessionKey) > 0 {
		d.NwkSKey, err = protocol.AESKeyFromString(hex.EncodeToString(req.NetworkSessionKey))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Invalid App Session Key")
		}
	}
	if req.FrameCountUp != nil {
		d.FCntUp = uint16(req.GetFrameCountUp())
	}
	if req.FrameCountDown != nil {
		d.FCntDn = uint16(req.GetFrameCountDown())
	}
	if req.KeyWarning != nil {
		d.KeyWarning = req.GetKeyWarning()
	}
	// Reset key warning if app session key and network session key is set
	if len(req.AppSessionKey) > 0 && len(req.NetworkSessionKey) > 0 {
		d.KeyWarning = false
	}
	if req.State != nil {
		switch d.State {
		case model.OverTheAirDevice:
			if len(req.AppKey) == 0 {
				return nil, status.Error(codes.InvalidArgument, "Must specify app key when changing device type to OTAA")
			}
		case model.PersonalizedDevice:
			if len(req.AppSessionKey) == 0 || len(req.NetworkSessionKey) == 0 {
				return nil, status.Error(codes.InvalidArgument, "Must specify app session key and network session key when changing device type to ABP")
			}
		default:
			// no checks
		}
	}
	if err := a.store.UpdateDevice(d); err != nil {
		return nil, toProtoErr(err)
	}
	return toAPIDevice(d), nil
}

func (a *apiServer) ListDevices(ctx context.Context, req *lospan.ListDeviceRequest) (*lospan.ListDeviceResponse, error) {
	eui, err := protocol.EUIFromString(req.ApplicationEui)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid application EUI")
	}
	devs, err := a.store.GetDevicesByApplicationEUI(eui)
	if err != nil {
		return nil, toProtoErr(err)
	}
	ret := &lospan.ListDeviceResponse{
		Devices: make([]*lospan.Device, 0),
	}
	for _, d := range devs {
		ret.Devices = append(ret.Devices, toAPIDevice(d))
	}
	return ret, nil
}
