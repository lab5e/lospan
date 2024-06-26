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
	if req.GetApplicationEui() == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing application EUI")
	}
	d := model.NewDevice()
	d.AppEUI, err = protocol.EUIFromString(req.GetApplicationEui())
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
	if d.State == model.OverTheAirDevice && (!d.AppSKey.Empty() || d.DevAddr.ToUint32() != 0 || !d.NwkSKey.Empty()) {
		return nil, status.Error(codes.InvalidArgument, "DevAddr, AppSKey and NwkSKey can only be specified for ABP devices")
	}
	if d.State == model.PersonalizedDevice && !d.AppKey.Empty() {
		return nil, status.Error(codes.InvalidArgument, "AppKey can only be specified for OTAA devices")
	}

	if d.State == model.OverTheAirDevice {
		if d.AppKey.Empty() {
			d.AppKey, err = protocol.NewAESKey()
			if err != nil {
				return nil, toProtoErr(err)
			}
		}
	}
	if d.State == model.PersonalizedDevice {
		if d.AppSKey.Empty() {
			d.AppSKey, err = protocol.NewAESKey()
			if err != nil {
				return nil, toProtoErr(err)
			}
		}
		if d.NwkSKey.Empty() {
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
			if d.AppKey.Empty() {
				// assign a new app key if it isn't set
				d.AppKey, err = protocol.NewAESKey()
				if err != nil {
					return nil, status.Error(codes.Internal, "Could not create application key for device")
				}
			}
		case model.PersonalizedDevice:
			if d.AppSKey.Empty() {
				d.AppSKey, err = protocol.NewAESKey()
				if err != nil {
					return nil, status.Error(codes.Internal, "Could not create application session key for device")
				}
			}
			if d.NwkSKey.Empty() {
				d.NwkSKey, err = protocol.NewAESKey()
				if err != nil {
					return nil, status.Error(codes.Internal, "Could not create network session key for device")
				}
			}
			if d.DevAddr.ToUint32() == 0 {
				d.DevAddr = protocol.NewDevAddr()
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
